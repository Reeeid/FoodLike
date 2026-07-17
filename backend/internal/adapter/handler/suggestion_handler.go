package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"foodlike-backend/internal/domain/model"
	"foodlike-backend/internal/usecase"
)

type SuggestionHandler struct {
	suggestions *usecase.SuggestionUsecase
}

func NewSuggestionHandler(suggestions *usecase.SuggestionUsecase) *SuggestionHandler {
	return &SuggestionHandler{suggestions: suggestions}
}

// Suggest GET /api/groups/:id/suggestions?large_area=&middle_area=&small_area=&budget=
// エリア/予算は任意(未指定なら絞り込まない)。
// レスポンスは集約後の候補のみ。誰の好みが理由で除外されたかは一切返さない。
func (h *SuggestionHandler) Suggest(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}
	criteria := model.SearchCriteria{
		LargeArea:  c.Query("large_area"),
		MiddleArea: c.Query("middle_area"),
		SmallArea:  c.Query("small_area"),
		Budget:     c.Query("budget"),
	}
	candidates, err := h.suggestions.Suggest(c.Request.Context(), currentMemberID(c), uint(id), criteria)
	if err != nil {
		// 非メンバーはグループの存在を漏らさないため404に統一。
		if errors.Is(err, usecase.ErrNotMember) {
			c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
			return
		}
		// APIエラー等は詳細をログに残し(観測用)、ユーザーには汎用メッセージを返す。
		log.Printf("suggest failed: group=%d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバー側でエラーが発生しました"})
		return
	}
	c.JSON(http.StatusOK, toCandidateResponses(candidates))
}
