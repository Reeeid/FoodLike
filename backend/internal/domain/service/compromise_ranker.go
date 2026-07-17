package service

import (
	"foodlike-backend/internal/domain/model"
	"sort"
)

// CompromiseRanker は「全員OKな店が0件」の場合に、
// 制約違反が少ない順の妥協案ランキングを作るドメインサービス。
type CompromiseRanker struct{}

func NewCompromiseRanker() CompromiseRanker {
	return CompromiseRanker{}
}

// ────────────────────────────────────────────────────────────
// 【課題3】Rank を実装せよ
//
// 仕様:
//   - 各店舗について、Genres と Constraints.ExcludedGenres の
//     重複数を ViolationCount として数える
//   - ViolationCount の昇順(違反が少ない順)に並べた Candidate の
//     スライスを返す。MatchedAll は ViolationCount == 0 のときtrue
//   - ViolationCount が同じ場合の順序は安定(元の並び順を維持)にすること
//     (ヒント: sort.SliceStable)
//   - 入力スライスを破壊しないこと
//
// ────────────────────────────────────────────────────────────
func (r CompromiseRanker) Rank(restaurants []model.Restaurant, c model.Constraints) []model.Candidate {
	candidates := make([]model.Candidate, 0, len(restaurants))
	for _, rest := range restaurants {
		violations := countViolations(rest, c)
		candidates = append(candidates, model.Candidate{
			Restaurant:     rest,
			MatchedAll:     violations == 0,
			ViolationCount: violations,
		})
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].ViolationCount < candidates[j].ViolationCount
	})

	return candidates
}

// countViolations は店舗の制約違反数を数える。
// 「ジャンルで弾く」と「食材で弾く」を明確に分けて集計する。
func countViolations(rest model.Restaurant, c model.Constraints) int {
	violationCount := 0
	// ジャンルで弾く: 実グルメAPIも提供できる粒度。
	for _, genre := range rest.Genres {
		if _, ok := c.ExcludedGenres[genre]; ok {
			violationCount++
		}
	}
	// 食材で弾く: 店舗が食材情報を持つ場合のみ有効。実APIはrest.Ingredientsが
	// 空になるため、ここは自然にno-op(＝実データでは食材弾きが効かない)。
	for _, ingredient := range rest.Ingredients {
		if _, ok := c.ExcludedIngredients[ingredient]; ok {
			violationCount++
		}
	}
	return violationCount
}
