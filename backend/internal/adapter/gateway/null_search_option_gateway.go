package gateway

import (
	"context"

	"foodlike-backend/internal/domain/model"
	port "foodlike-backend/internal/port/gateway"
)

// nullSearchOptionGateway は HOTPEPPER_API_KEY 未設定時のマスタ取得ポート。
// マスタはHotPepper由来のため、キーが無ければ空を返す(=フロントの絞り込みは無効)。
type nullSearchOptionGateway struct{}

func NewNullSearchOptionGateway() port.SearchOptionGateway {
	return &nullSearchOptionGateway{}
}

func (g *nullSearchOptionGateway) ListAreas(_ context.Context) ([]model.AreaEntry, error) {
	return []model.AreaEntry{}, nil
}

func (g *nullSearchOptionGateway) ListBudgets(_ context.Context) ([]model.BudgetOption, error) {
	return []model.BudgetOption{}, nil
}
