package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"foodlike-backend/internal/usecase"
)

type SuggestionHandler struct {
	suggestions *usecase.SuggestionUsecase
}

func NewSuggestionHandler(suggestions *usecase.SuggestionUsecase) *SuggestionHandler {
	return &SuggestionHandler{suggestions: suggestions}
}

// Suggest GET /api/groups/:id/suggestions?area=新宿
// レスポンスは集約後の候補のみ。誰の好みが理由で除外されたかは一切返さない。
func (h *SuggestionHandler) Suggest(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}
	candidates, err := h.suggestions.Suggest(c.Request.Context(), uint(id), c.Query("area"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to suggest restaurants"})
		return
	}
	c.JSON(http.StatusOK, toCandidateResponses(candidates))
}
