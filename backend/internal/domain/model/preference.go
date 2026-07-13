package model

// PreferenceKind は好き/嫌いの種別。MVPでは嫌い(dislike)のみ提案ロジックで使う。
type PreferenceKind string

const (
	PreferenceKindLike    PreferenceKind = "like"
	PreferenceKindDislike PreferenceKind = "dislike"
)

// PreferenceCategory は好き嫌いの粒度(ジャンル or 食材)。
type PreferenceCategory string

const (
	PreferenceCategoryGenre      PreferenceCategory = "genre"
	PreferenceCategoryIngredient PreferenceCategory = "ingredient"
)

// Preference は個人の好き嫌い1件。
// 本人以外に公開してはならない(APIレスポンスに他人のPreferenceを含めない)。
type Preference struct {
	ID       uint
	MemberID uint
	Kind     PreferenceKind
	Category PreferenceCategory
	// Value はジャンル名("辛い物"、"中華"等)や食材名("エビ"等)。
	Value string
}
