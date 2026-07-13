package service

import (
	"reflect"
	"sort"
	"testing"

	"foodlike-backend/internal/domain/model"
)

// 課題1を実装したら t.Skip の行を削除して `go test ./...` を通すこと。
func TestPreferenceAggregator_Aggregate(t *testing.T) {

	a := NewPreferenceAggregator()

	prefs := []model.Preference{
		{MemberID: 1, Kind: model.PreferenceKindDislike, Category: model.PreferenceCategoryGenre, Value: "辛い物"},
		{MemberID: 2, Kind: model.PreferenceKindDislike, Category: model.PreferenceCategoryGenre, Value: "辛い物"}, // 重複
		{MemberID: 2, Kind: model.PreferenceKindDislike, Category: model.PreferenceCategoryIngredient, Value: "エビ"},
		{MemberID: 3, Kind: model.PreferenceKindLike, Category: model.PreferenceCategoryGenre, Value: "中華"}, // likeは無視
		{MemberID: 3, Kind: model.PreferenceKindDislike, Category: model.PreferenceCategoryGenre, Value: "中華"},
	}

	got := a.Aggregate(prefs)

	wantGenres := []string{"中華", "辛い物"}
	wantIngredients := []string{"エビ"}

	gotGenres := sortedKeys(got.ExcludedGenres)
	gotIngredients := sortedKeys(got.ExcludedIngredients)

	if !reflect.DeepEqual(gotGenres, wantGenres) {
		t.Errorf("ExcludedGenres = %v, want %v", gotGenres, wantGenres)
	}
	if !reflect.DeepEqual(gotIngredients, wantIngredients) {
		t.Errorf("ExcludedIngredients = %v, want %v", gotIngredients, wantIngredients)
	}
}

// sortedKeys は集合(map)のキーをソート済みスライスとして返すテスト用ヘルパー。
func sortedKeys(set map[string]struct{}) []string {
	keys := make([]string, 0, len(set))
	for k := range set {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func TestPreferenceAggregator_Aggregate_Empty(t *testing.T) {

	a := NewPreferenceAggregator()
	got := a.Aggregate(nil)

	if len(got.ExcludedGenres) != 0 || len(got.ExcludedIngredients) != 0 {
		t.Errorf("空入力なら空の制約を返すこと: got %+v", got)
	}
}
