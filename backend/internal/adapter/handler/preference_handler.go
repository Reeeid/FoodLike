package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"foodlike-backend/internal/domain/model"
	"foodlike-backend/internal/usecase"
)

type PreferenceHandler struct {
	prefs *usecase.PreferenceUsecase
}

func NewPreferenceHandler(prefs *usecase.PreferenceUsecase) *PreferenceHandler {
	return &PreferenceHandler{prefs: prefs}
}

// ListOwn GET /api/me/preferences
// 本人の好き嫌いのみ返す。他人の好き嫌いを取得するエンドポイントは存在させない。
func (h *PreferenceHandler) ListOwn(c *gin.Context) {
	prefs, err := h.prefs.ListOwn(c.Request.Context(), currentMemberID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list preferences"})
		return
	}
	res := make([]preferenceItem, 0, len(prefs))
	for _, p := range prefs {
		res = append(res, preferenceItem{Kind: string(p.Kind), Category: string(p.Category), Value: p.Value})
	}
	c.JSON(http.StatusOK, res)
}

type replacePreferencesRequest struct {
	Preferences []preferenceItem `json:"preferences" binding:"required,dive"`
}

// Replace PUT /api/me/preferences
func (h *PreferenceHandler) Replace(c *gin.Context) {
	var req replacePreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	memberID := currentMemberID(c)
	prefs := make([]model.Preference, 0, len(req.Preferences))
	for _, p := range req.Preferences {
		prefs = append(prefs, model.Preference{
			MemberID: memberID,
			Kind:     model.PreferenceKind(p.Kind),
			Category: model.PreferenceCategory(p.Category),
			Value:    p.Value,
		})
	}
	if err := h.prefs.Replace(c.Request.Context(), memberID, prefs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save preferences"})
		return
	}
	c.Status(http.StatusNoContent)
}
