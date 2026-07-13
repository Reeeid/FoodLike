package model

// Candidate は提案候補の店舗。
type Candidate struct {
	Restaurant Restaurant
	// MatchedAll は制約条件をすべて満たす(=全員OK)かどうか。
	MatchedAll bool
	// ViolationCount は違反した制約の数。妥協案ランキングで小さい順に並べる。
	ViolationCount int
}
