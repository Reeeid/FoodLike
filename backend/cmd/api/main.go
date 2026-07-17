package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	adaptergw "foodlike-backend/internal/adapter/gateway"
	"foodlike-backend/internal/adapter/handler"
	adapterrepo "foodlike-backend/internal/adapter/repository"
	"foodlike-backend/internal/domain/service"
	"foodlike-backend/internal/port/gateway"
	"foodlike-backend/internal/infrastructure/db"
	"foodlike-backend/internal/infrastructure/firebaseauth"
	"foodlike-backend/internal/infrastructure/router"
	"foodlike-backend/internal/usecase"
)

func main() {
	conn, err := db.NewMySQLConnection()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	if err := db.Migrate(conn); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	// 手動DI。依存が増えてきたらwire等の導入を検討する。
	memberRepo := adapterrepo.NewMemberRepository(conn)
	groupRepo := adapterrepo.NewGroupRepository(conn)
	prefRepo := adapterrepo.NewPreferenceRepository(conn)

	// 外部グルメAPI: HOTPEPPER_API_KEY があればホットペッパー実装、無ければモック。
	// APIエラーは握りつぶさずhandlerで500として表に出す(silent fallbackはしない)。
	// エリア/予算マスタも同じキーで取得する(キー無しなら空=絞り込み無効)。
	restaurantGw := adaptergw.NewMockRestaurantGateway()
	searchOptionGw := adaptergw.NewNullSearchOptionGateway()
	if key := os.Getenv("HOTPEPPER_API_KEY"); key != "" {
		hp := adaptergw.NewHotPepperRestaurantGateway(key)
		restaurantGw = hp
		searchOptionGw = hp
		log.Println("using HotPepper restaurant gateway")
	} else {
		log.Println("WARN: HOTPEPPER_API_KEY not set; using mock restaurant gateway")
	}

	memberUC := usecase.NewMemberUsecase(memberRepo)
	groupUC := usecase.NewGroupUsecase(groupRepo)
	prefUC := usecase.NewPreferenceUsecase(prefRepo)
	suggestionUC := usecase.NewSuggestionUsecase(
		groupRepo,
		prefRepo,
		restaurantGw,
		service.NewPreferenceAggregator(),
		service.NewRestaurantFilter(),
		service.NewCompromiseRanker(),
	)

	searchOptionUC := usecase.NewSearchOptionUsecase(searchOptionGw)

	// チャットのメッセージはGORMで永続化。AI検索の回答生成はAIResponderポート。
	// GEMINI_API_KEY があればGemini実装、無ければ定型文フェイク。
	// Geminiは失敗時に定型文へ自動フォールバックするのでチャットは止まらない。
	messageRepo := adapterrepo.NewMessageRepository(conn)
	var aiResponder gateway.AIResponder = adaptergw.NewFakeAIResponder()
	if key := os.Getenv("GEMINI_API_KEY"); key != "" {
		// 実在店舗の根拠の取り方は2択。TAVILY_API_KEY があれば無料枠のTavily検索を
		// 前段に挟み、無ければ GEMINI_GROUNDING(既定ON)でGeminiのGoogle検索に任せる。
		var searcher gateway.WebSearcher
		grounding := os.Getenv("GEMINI_GROUNDING") != "false"
		if tavilyKey := os.Getenv("TAVILY_API_KEY"); tavilyKey != "" {
			searcher = adaptergw.NewTavilyWebSearcher(tavilyKey)
			grounding = false // Tavilyで根拠を取るのでGoogle検索(課金枠)は使わない
			log.Println("using Tavily web searcher for AI chat search")
		}
		g, err := adaptergw.NewGeminiResponder(context.Background(), adaptergw.GeminiConfig{
			APIKey:    key,
			Model:     os.Getenv("GEMINI_MODEL"),
			Grounding: grounding,
			Searcher:  searcher,
			Fallback:  aiResponder,
		})
		if err != nil {
			log.Fatalf("failed to init Gemini responder: %v", err)
		}
		aiResponder = g
		log.Printf("using Gemini AI responder (grounding=%v, tavily=%v)", grounding, searcher != nil)
	} else {
		log.Println("WARN: GEMINI_API_KEY not set; using template (fake) AI responder")
	}
	chatUC := usecase.NewChatUsecase(groupRepo, messageRepo, aiResponder, suggestionUC)

	r := router.New(conn, router.Handlers{
		Member:       handler.NewMemberHandler(memberUC),
		Group:        handler.NewGroupHandler(groupUC),
		Preference:   handler.NewPreferenceHandler(prefUC),
		Suggestion:   handler.NewSuggestionHandler(suggestionUC),
		Chat:         handler.NewChatHandler(chatUC),
		SearchOption: handler.NewSearchOptionHandler(searchOptionUC),
	}, newAuthMiddleware(memberUC))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

// newAuthMiddleware は認証方式を選ぶ。FIREBASE_PROJECT_IDがあればFirebase認証、
// なければローカル開発用のX-Member-IDモック認証にフォールバックする。
func newAuthMiddleware(memberUC *usecase.MemberUsecase) gin.HandlerFunc {
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	if projectID == "" {
		log.Println("WARN: FIREBASE_PROJECT_ID not set; using mock auth (X-Member-ID header)")
		return handler.MockAuthMiddleware()
	}
	verifier, err := firebaseauth.NewVerifier(context.Background(), projectID)
	if err != nil {
		log.Fatalf("failed to init firebase auth: %v", err)
	}
	return handler.FirebaseAuthMiddleware(verifier, memberUC)
}
