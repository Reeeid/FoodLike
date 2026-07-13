package service

import (
	"testing"

	"foodlike-backend/internal/domain/model"
)

// 課題2を実装したら t.Skip の行を削除して `go test ./...` を通すこと。
func TestRestaurantFilter_Filter(t *testing.T) {

	f := NewRestaurantFilter()

	restaurants := []model.Restaurant{
		{ID: "r1", Name: "四川料理店", Genres: []string{"中華", "辛い物"}},
		{ID: "r2", Name: "寿司屋", Genres: []string{"和食", "寿司"}},
		{ID: "r3", Name: "町中華", Genres: []string{"中華"}},
		{ID: "r4", Name: "イタリアン", Genres: []string{"イタリアン"}},
	}
	c := model.Constraints{ExcludedGenres: map[string]struct{}{
		"辛い物": {},
		"中華":  {},
	}}

	got := f.Filter(restaurants, c)

	wantIDs := []string{"r2", "r4"}
	if len(got) != len(wantIDs) {
		t.Fatalf("件数 = %d, want %d (got: %+v)", len(got), len(wantIDs), got)
	}
	for i, id := range wantIDs {
		if got[i].ID != id {
			t.Errorf("got[%d].ID = %s, want %s", i, got[i].ID, id)
		}
	}
}

func TestRestaurantFilter_Filter_NoConstraints(t *testing.T) {

	f := NewRestaurantFilter()
	restaurants := []model.Restaurant{
		{ID: "r1", Genres: []string{"中華"}},
		{ID: "r2", Genres: []string{"和食"}},
	}

	got := f.Filter(restaurants, model.Constraints{})

	if len(got) != 2 {
		t.Errorf("制約が空なら全店舗を返すこと: got %d件", len(got))
	}
}
