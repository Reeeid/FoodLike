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

	// TODO(issue #5): 外部グルメAPI選定後、実APIのGateway実装に差し替える。
	restaurantGw := adaptergw.NewMockRestaurantGateway()

	memberUC := usecase.NewMemberUsecase(memberRepo)
	groupUC := usecase.NewGroupUsecase(groupRepo)
	prefUC := usecase.NewPreferenceUsecase(prefRepo)
	suggestionUC := usecase.NewSuggestionUsecase(
		prefRepo,
		restaurantGw,
		service.NewPreferenceAggregator(),
		service.NewRestaurantFilter(),
		service.NewCompromiseRanker(),
	)

	r := router.New(conn, router.Handlers{
		Member:     handler.NewMemberHandler(memberUC),
		Group:      handler.NewGroupHandler(groupUC),
		Preference: handler.NewPreferenceHandler(prefUC),
		Suggestion: handler.NewSuggestionHandler(suggestionUC),
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
