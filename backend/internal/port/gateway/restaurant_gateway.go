package gateway

import (
	"context"

	"foodlike-backend/internal/domain/model"
)

// RestaurantGateway は外部グルメAPIのポート。
// ホットペッパー/Google Places等への乗り換え・併用を想定してinterface化している。
// 実装の差し替えは adapter/gateway 側だけで完結すること(issue #5)。
type RestaurantGateway interface {
	Search(ctx context.Context, area string) ([]model.Restaurant, error)
}
