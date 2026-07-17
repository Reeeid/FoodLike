package model

// Restaurant は外部グルメAPIから取得する店舗情報。
type Restaurant struct {
	// ID は外部APIの店舗ID(自前DBのIDではないためstring)。
	ID   string
	Name string
	Area string
	// Genres は料理ジャンル(中華/和食/辛い物 等)。実グルメAPIも提供できる粒度で、
	// 「ジャンルで弾く」フィルタに使う。
	Genres []string
	// Ingredients は食材タグ(エビ/生魚/パクチー 等)。「食材で弾く」フィルタに使う。
	// 実グルメAPIは食材を返さないため実データでは空になり、食材弾きは
	// 食材情報を持つデータ(モック等)でのみ有効になる。
	Ingredients []string
	// Budget は1人あたりの予算目安(円)。
	Budget int
}
