package model

// Constraints はグループ全員の好き嫌いを「個人と紐づかない集合」に潰した制約条件。
//
// プライバシー設計の要:
// PreferenceAggregatorでこの形に変換した後は、誰がどれを嫌いかという
// 対応関係は失われる。DTO/APIレスポンス層やLLMに渡してよいのはこの形まで。
type Constraints struct {
	// ExcludedGenres は誰か1人でも嫌いなジャンルの集合(重複なし)。
	ExcludedGenres map[string]struct{}
	// ExcludedIngredients は誰か1人でも嫌いな食材の集合(重複なし)。
	// 「食材で弾く」フィルタに使う。店舗が食材情報を持つ場合のみ有効
	// (実グルメAPIは食材を返さないため、実データでは効かない)。
	ExcludedIngredients map[string]struct{}
}
