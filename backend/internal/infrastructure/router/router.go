package router

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"foodlike-backend/internal/adapter/handler"
)

type Handlers struct {
	Member       *handler.MemberHandler
	Group        *handler.GroupHandler
	Preference   *handler.PreferenceHandler
	Suggestion   *handler.SuggestionHandler
	Chat         *handler.ChatHandler
	SearchOption *handler.SearchOptionHandler
}

// New はルーティングを組み立てる。認証方式(Firebase/モック)はmain側で
// 決めてauthMiddlewareとして注入する。
func New(db *gorm.DB, h Handlers, authMiddleware gin.HandlerFunc) *gin.Engine {
	r := gin.Default()

	// どのプロキシのX-Forwarded-Forも信頼しない。
	// Cloud Runの前段はGoogleのフロントエンドで信頼するIPレンジを列挙できず、
	// 全信頼(既定)のままだとクライアントが自分でX-Forwarded-Forを送ってIPを詐称できる。
	// 本アプリはClientIPを認可・制限に使わない(Ginのアクセスログのみ)ため、
	// ログのIPが接続元になる代わりに詐称を不可能にする方を取る。
	// 将来IPベースのレート制限を入れるならここの再検討が要る。
	_ = r.SetTrustedProxies(nil) // nil指定ではエラーは返らない

	r.Use(corsMiddleware())

	r.GET("/health", handler.NewHealthHandler(db))

	api := r.Group("/api")
	api.POST("/members", h.Member.Register)

	authed := api.Group("")
	authed.Use(authMiddleware)
	{
		authed.GET("/me", h.Member.Me)
		authed.GET("/me/preferences", h.Preference.ListOwn)
		authed.PUT("/me/preferences", h.Preference.Replace)

		authed.POST("/groups", h.Group.Create)
		authed.POST("/groups/join", h.Group.Join)
		authed.GET("/groups", h.Group.List)
		authed.GET("/groups/:id", h.Group.Get)
		authed.GET("/search-options", h.SearchOption.List)
		authed.GET("/groups/:id/suggestions", h.Suggestion.Suggest)
		authed.GET("/groups/:id/messages", h.Chat.List)
		authed.POST("/groups/:id/messages", h.Chat.Post)
		authed.GET("/groups/:id/ai-search", h.Chat.AISearch)
	}

	return r
}

// corsMiddleware はフロントエンド(別オリジン)からのアクセスを許可する。
// 依存を増やさないため自前実装。許可オリジンはFRONTEND_ORIGINで指定。
func corsMiddleware() gin.HandlerFunc {
	origin := os.Getenv("FRONTEND_ORIGIN")
	if origin == "" {
		origin = "http://localhost:3000"
	}
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Member-ID")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
