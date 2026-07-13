package service

import "foodlike-backend/internal/domain/model"

// PreferenceAggregator はグループ内メンバーの好き嫌いを
// 「個人と紐づかない集合の制約条件(Constraints)」に潰すドメインサービス。
//
// プライバシー設計の第一関門: ここから先の層(Gateway/DTO/LLM)には
// 個人を特定できる情報を一切渡さない。
type PreferenceAggregator struct{}

func NewPreferenceAggregator() PreferenceAggregator {
	return PreferenceAggregator{}
}

// ────────────────────────────────────────────────────────────
// 【課題1】Aggregate を実装せよ
//
// 仕様:
//   - kind = dislike の Preference のみを集約する(likeは無視)
//   - category = genre      → Constraints.ExcludedGenres へ
//   - category = ingredient → Constraints.ExcludedIngredients へ
//   - 同じ値の重複は除去する(複数人が同じものを嫌いでも1つにまとめる)
//   - 出力に MemberID など個人を特定できる情報を含めてはならない
//   - 入力の順序に依存しない安定した結果になっているとテストしやすい
//     (ヒント: mapで重複除去 → sortする、など)
//
// テスト: preference_aggregator_test.go の t.Skip を外して全部通せばOK
//
// 考えるポイント:
//   - なぜこの変換を usecase層 ではなく domain/service に置くのか?
//   - 「好みの対応関係を潰す」処理としてこの実装で十分か?
//     (例: グループに1人しかいない場合、集合を見れば個人が特定できてしまう。
//     MVPでは許容するとして、将来どう対策するか考えてみる)
//
// ────────────────────────────────────────────────────────────
func (a PreferenceAggregator) Aggregate(prefs []model.Preference) model.Constraints {

	c := model.Constraints{
		ExcludedGenres:      map[string]struct{}{},
		ExcludedIngredients: map[string]struct{}{},
	}
	for _, p := range prefs {
		if p.Kind == model.PreferenceKindDislike {
			switch p.Category {
			case model.PreferenceCategoryGenre:
				c.ExcludedGenres[p.Value] = struct{}{}
			case model.PreferenceCategoryIngredient:
				c.ExcludedIngredients[p.Value] = struct{}{}
			}
		}
	}
	return c
}
