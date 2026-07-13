package model

// Restaurant は外部グルメAPIから取得する店舗情報。
type Restaurant struct {
	// ID は外部APIの店舗ID(自前DBのIDではないためstring)。
	ID     string
	Name   string
	Area   string
	Genres []string
	// Budget は1人あたりの予算目安(円)。
	Budget int
}
