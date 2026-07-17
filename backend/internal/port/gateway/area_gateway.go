package gateway

import (
	"context"

	"foodlike-backend/internal/domain/model"
)

// SearchOptionGateway は外部グルメAPIのマスタ(エリア/予算)を取得するポート。
// エリア・予算の選択肢をフロントへ渡すために使う。
type SearchOptionGateway interface {
	ListAreas(ctx context.Context) ([]model.AreaEntry, error)
	ListBudgets(ctx context.Context) ([]model.BudgetOption, error)
}
