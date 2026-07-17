package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"foodlike-backend/internal/usecase"
)

type SearchOptionHandler struct {
	options *usecase.SearchOptionUsecase
}

func NewSearchOptionHandler(options *usecase.SearchOptionUsecase) *SearchOptionHandler {
	return &SearchOptionHandler{options: options}
}

type areaOption struct {
	LargeCode  string `json:"large_code"`
	LargeName  string `json:"large_name"`
	MiddleCode string `json:"middle_code"`
	MiddleName string `json:"middle_name"`
	SmallCode  string `json:"small_code"`
	SmallName  string `json:"small_name"`
}

type budgetOption struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type searchOptionsResponse struct {
	Areas   []areaOption   `json:"areas"`
	Budgets []budgetOption `json:"budgets"`
}

// List GET /api/search-options
// エリア/予算の選択肢を返す。HOTPEPPER_API_KEY未設定なら空配列(絞り込み無効)。
func (h *SearchOptionHandler) List(c *gin.Context) {
	ctx := c.Request.Context()

	areas, err := h.options.ListAreas(ctx)
	if err != nil {
		log.Printf("list areas failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバー側でエラーが発生しました"})
		return
	}
	budgets, err := h.options.ListBudgets(ctx)
	if err != nil {
		log.Printf("list budgets failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "サーバー側でエラーが発生しました"})
		return
	}

	res := searchOptionsResponse{
		Areas:   make([]areaOption, 0, len(areas)),
		Budgets: make([]budgetOption, 0, len(budgets)),
	}
	for _, a := range areas {
		res.Areas = append(res.Areas, areaOption{
			LargeCode:  a.LargeCode,
			LargeName:  a.LargeName,
			MiddleCode: a.MiddleCode,
			MiddleName: a.MiddleName,
			SmallCode:  a.SmallCode,
			SmallName:  a.SmallName,
		})
	}
	for _, b := range budgets {
		res.Budgets = append(res.Budgets, budgetOption{Code: b.Code, Name: b.Name})
	}
	c.JSON(http.StatusOK, res)
}
