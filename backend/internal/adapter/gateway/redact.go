package gateway

import (
	"regexp"
	"sort"
	"strings"

	"foodlike-backend/internal/domain/model"
)

// redactPlaceholder は伏せ字。具体的な品目の代わりにこれを出す。
const redactPlaceholder = "苦手なもの"

// redactedRun は「苦手なもの、苦手なもの、苦手なもの」のような連続を1つに畳む。
var redactedRun = regexp.MustCompile(`(?:` + redactPlaceholder + `)(?:\s*[、,・と/]\s*(?:` + redactPlaceholder + `))+`)

// redactExcluded は回答から「避けたい品目」の具体名を消す。
//
// なぜコードでやるか:
// システム指示で「列挙するな」と頼んでもLLMは確率的に破る。しかしこれは
// プライバシー要件(誰が何を苦手か推測させない)なので、確率ではなく決定的に守る必要がある。
// 特に2人グループでは「集合 - 自分の苦手 = 相方の苦手」と逆算できてしまうため、
// 集約済み(誰のか不明)であっても具体名を出してはいけない。
//
// 制約: 単純な文字列置換なので「きのこ類」→「苦手なもの類」のような粗さは残る。
// 表示の綺麗さより漏らさないことを優先している。
func redactExcluded(text string, c model.Constraints) string {
	terms := append(sortedKeys(c.ExcludedGenres), sortedKeys(c.ExcludedIngredients)...)

	// 長い語から置換する。短い語が先だと「辛い物」の一部だけ潰れる等の取りこぼしが出る。
	sort.Slice(terms, func(i, j int) bool { return len([]rune(terms[i])) > len([]rune(terms[j])) })

	for _, t := range terms {
		if strings.TrimSpace(t) == "" {
			continue
		}
		text = strings.ReplaceAll(text, t, redactPlaceholder)
	}

	// 「苦手なもの、苦手なもの、…を避けて」と並ぶと不自然なので1つに畳む。
	return redactedRun.ReplaceAllString(text, redactPlaceholder)
}
