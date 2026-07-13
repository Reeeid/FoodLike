package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"foodlike-backend/internal/port/gateway"
	"foodlike-backend/internal/usecase"
)

const memberIDKey = "memberID"

// FirebaseAuthMiddleware はAuthorization: Bearer <IDトークン> を検証し、
// Firebase UIDから内部メンバーIDを解決してコンテキストに載せる(issue #8)。
// ハンドラ側は currentMemberID(c) 経由でしかIDを取らないので、
// 認証方式が変わってもハンドラは影響を受けない。
func FirebaseAuthMiddleware(verifier gateway.AuthVerifier, members *usecase.MemberUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := strings.CutPrefix(c.GetHeader("Authorization"), "Bearer ")
		if !ok || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization: Bearer <ID token> required"})
			return
		}
		auth, err := verifier.Verify(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}
		m, err := members.Authenticate(c.Request.Context(), auth.UID, auth.Name)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve member"})
			return
		}
		c.Set(memberIDKey, m.ID)
		c.Next()
	}
}

// MockAuthMiddleware はローカル開発用の仮認証。X-Member-IDヘッダーをそのまま信用する。
// FIREBASE_PROJECT_ID未設定のときだけ使われる(main.goで切替)。
func MockAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.GetHeader("X-Member-ID"), 10, 64)
		if err != nil || id == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "X-Member-ID header required"})
			return
		}
		c.Set(memberIDKey, uint(id))
		c.Next()
	}
}

func currentMemberID(c *gin.Context) uint {
	return c.GetUint(memberIDKey)
}
