package gateway

import (
	"strings"
	"testing"

	"foodlike-backend/internal/domain/model"
)

func constraints(genres, ingredients []string) model.Constraints {
	c := model.Constraints{
		ExcludedGenres:      map[string]struct{}{},
		ExcludedIngredients: map[string]struct{}{},
	}
	for _, g := range genres {
		c.ExcludedGenres[g] = struct{}{}
	}
	for _, i := range ingredients {
		c.ExcludedIngredients[i] = struct{}{}
	}
	return c
}

// LLMは指示に従わないことがある前提で、出力に品目名が残らないことを保証する。
func TestRedactExcluded(t *testing.T) {
	tests := []struct {
		name string
		text string
		c    model.Constraints
		want string
	}{
		{
			name: "冒頭の復唱を伏せて連続を畳む",
			text: "大阪駅周辺で、洋食、辛い物、きのこ、野菜を避けられそうなお店を3つご紹介しますね。",
			c:    constraints([]string{"洋食", "辛い物"}, []string{"きのこ", "野菜"}),
			want: "大阪駅周辺で、苦手なものを避けられそうなお店を3つご紹介しますね。",
		},
		{
			name: "単独の言及も伏せる",
			text: "きのこが苦手とのことなので和食にしました。",
			c:    constraints(nil, []string{"きのこ"}),
			want: "苦手なものが苦手とのことなので和食にしました。",
		},
		{
			// 「辛い」が先に置換されると「苦手なもの物」と品目の破片が残る。
			// 長い語から処理することでそれを防ぐ。畳み込みで「と」繋ぎも1つになる。
			name: "長い語を優先して置換し破片を残さない",
			text: "辛い物と辛いラーメンは避けました。",
			c:    constraints([]string{"辛い", "辛い物"}, nil),
			want: "苦手なものラーメンは避けました。",
		},
		{
			name: "避けたい条件が無ければ素通し",
			text: "和食のお店を3つご紹介します。",
			c:    constraints(nil, nil),
			want: "和食のお店を3つご紹介します。",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := redactExcluded(tt.text, tt.c); got != tt.want {
				t.Errorf("redactExcluded()\n got = %q\nwant = %q", got, tt.want)
			}
		})
	}
}

// 不変条件: 検閲後のテキストに品目名が1つも残らないこと。
func TestRedactExcluded_NoTermSurvives(t *testing.T) {
	c := constraints([]string{"洋食", "辛い物"}, []string{"きのこ", "野菜", "パクチー"})
	text := "洋食も辛い物もきのこも野菜もパクチーも避けたお店です。パクチー抜き可。"

	got := redactExcluded(text, c)

	for term := range c.ExcludedGenres {
		if strings.Contains(got, term) {
			t.Errorf("ジャンル %q が出力に残っている: %q", term, got)
		}
	}
	for term := range c.ExcludedIngredients {
		if strings.Contains(got, term) {
			t.Errorf("食材 %q が出力に残っている: %q", term, got)
		}
	}
}
