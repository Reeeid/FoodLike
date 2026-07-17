package model

// SearchCriteria は店舗検索の絞り込み条件。すべて任意(空文字なら絞り込まない)。
// エリアはHotPepperのエリアコード(大/中/小)、Budgetは予算コード。
type SearchCriteria struct {
	LargeArea  string
	MiddleArea string
	SmallArea  string
	Budget     string
}

// AreaEntry は小エリア1件と、その親の中/大エリアをまとめた平坦なエントリ。
// フロントはこれを検索可能なエリア候補(react-select)として使う。
type AreaEntry struct {
	LargeCode  string
	LargeName  string
	MiddleCode string
	MiddleName string
	SmallCode  string
	SmallName  string
}

// BudgetOption は予算の選択肢1件(コードと表示名"2001～3000円"等)。
type BudgetOption struct {
	Code string
	Name string
}

// WebSearchResult はWeb検索1件。AIチャット検索で「実在する店の根拠」として
// LLMに渡すために使う(LLMの学習知識だけに頼ると存在しない店を作りがちなため)。
type WebSearchResult struct {
	Title   string
	URL     string
	Content string
}
