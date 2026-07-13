package service

import (
	"testing"

	"foodlike-backend/internal/domain/model"
)

// 課題3を実装したら t.Skip の行を削除して `go test ./...` を通すこと。
func TestCompromiseRanker_Rank(t *testing.T) {

	r := NewCompromiseRanker()

	restaurants := []model.Restaurant{
		{ID: "r1", Name: "四川料理店", Genres: []string{"中華", "辛い物"}}, // 違反2
		{ID: "r2", Name: "寿司屋", Genres: []string{"和食", "寿司"}},    // 違反0
		{ID: "r3", Name: "町中華", Genres: []string{"中華"}},          // 違反1
	}
	c := model.Constraints{ExcludedGenres: map[string]struct{}{"辛い物": {}, "中華": {}}}

	got := r.Rank(restaurants, c)

	if len(got) != 3 {
		t.Fatalf("件数 = %d, want 3", len(got))
	}

	// 違反数の昇順
	wantOrder := []struct {
		id         string
		violations int
		matchedAll bool
	}{
		{"r2", 0, true},
		{"r3", 1, false},
		{"r1", 2, false},
	}
	for i, w := range wantOrder {
		if got[i].Restaurant.ID != w.id {
			t.Errorf("got[%d].ID = %s, want %s", i, got[i].Restaurant.ID, w.id)
		}
		if got[i].ViolationCount != w.violations {
			t.Errorf("got[%d].ViolationCount = %d, want %d", i, got[i].ViolationCount, w.violations)
		}
		if got[i].MatchedAll != w.matchedAll {
			t.Errorf("got[%d].MatchedAll = %v, want %v", i, got[i].MatchedAll, w.matchedAll)
		}
	}
}

func TestCompromiseRanker_Rank_StableOrder(t *testing.T) {

	r := NewCompromiseRanker()

	// 全店舗が違反数0 → 元の並び順を維持すること
	restaurants := []model.Restaurant{
		{ID: "r1", Genres: []string{"和食"}},
		{ID: "r2", Genres: []string{"イタリアン"}},
		{ID: "r3", Genres: []string{"寿司"}},
	}

	got := r.Rank(restaurants, model.Constraints{})

	for i, id := range []string{"r1", "r2", "r3"} {
		if got[i].Restaurant.ID != id {
			t.Errorf("同順位は元の並び順を維持すること: got[%d].ID = %s, want %s", i, got[i].Restaurant.ID, id)
		}
	}
}
