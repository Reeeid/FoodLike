package service

import "foodlike-backend/internal/domain/model"

// RestaurantFilter は制約条件(Constraints)で店舗候補を絞り込むドメインサービス。
type RestaurantFilter struct{}

func NewRestaurantFilter() RestaurantFilter {
	return RestaurantFilter{}
}

// ────────────────────────────────────────────────────────────
// 【課題2】Filter を実装せよ
//
// 仕様:
//   - 「ジャンルで弾く」: Restaurant.Genres のいずれかが ExcludedGenres に
//     含まれる店舗を除外
//   - 「食材で弾く」: Restaurant.Ingredients のいずれかが ExcludedIngredients に
//     含まれる店舗を除外(店舗が食材情報を持つ場合のみ有効)
//   - 制約が空なら全店舗をそのまま返す
//   - 入力スライスを破壊(並び替え・書き換え)しないこと
//
// テスト: restaurant_filter_test.go の t.Skip を外して全部通せばOK
//
// 考えるポイント:
//   - 計算量は? 店舗数N × 除外ジャンルM で O(N*M) になるが、
//     除外ジャンルを map(set) にすれば O(N+M) にできる。
//     MVPの規模で最適化する価値はあるか?(YAGNIとのバランス)
//
// ────────────────────────────────────────────────────────────
func (f RestaurantFilter) Filter(restaurants []model.Restaurant, c model.Constraints) []model.Restaurant {
	filtered := make([]model.Restaurant, 0, len(restaurants))
	for _, rest := range restaurants {
		if countViolations(rest, c) == 0 {
			filtered = append(filtered, rest)
		}
	}
	return filtered
}
