package usecase

import (
	"context"

	"foodlike-backend/internal/domain/model"
	"foodlike-backend/internal/port/gateway"
)

// SearchOptionUsecase はエリア/予算の選択肢(マスタ)を取得する。
type SearchOptionUsecase struct {
	gw gateway.SearchOptionGateway
}

func NewSearchOptionUsecase(gw gateway.SearchOptionGateway) *SearchOptionUsecase {
	return &SearchOptionUsecase{gw: gw}
}

func (u *SearchOptionUsecase) ListAreas(ctx context.Context) ([]model.AreaEntry, error) {
	return u.gw.ListAreas(ctx)
}

func (u *SearchOptionUsecase) ListBudgets(ctx context.Context) ([]model.BudgetOption, error) {
	return u.gw.ListBudgets(ctx)
}
