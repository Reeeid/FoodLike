package gateway

import (
	"context"
	"strings"

	"foodlike-backend/internal/domain/model"
	port "foodlike-backend/internal/port/gateway"
)

// mockRestaurantGateway は外部グルメAPI選定(issue #5)までのつなぎの固定データ実装。
// ホットペッパー/Google Places実装に差し替えるときは、このファイルと同じ形で
// port.RestaurantGateway を満たす実装を追加し、main.goの配線を変えるだけでよい。
type mockRestaurantGateway struct{}

func NewMockRestaurantGateway() port.RestaurantGateway {
	return &mockRestaurantGateway{}
}

var mockRestaurants = []model.Restaurant{
	{ID: "m1", Name: "四川菜館 炎", Area: "新宿", Genres: []string{"中華", "辛い物"}, Budget: 3000},
	{ID: "m2", Name: "鮨 さかえ", Area: "新宿", Genres: []string{"和食", "寿司", "生魚"}, Budget: 6000},
	{ID: "m3", Name: "トラットリア・ソーレ", Area: "新宿", Genres: []string{"イタリアン", "チーズ"}, Budget: 4000},
	{ID: "m4", Name: "焼肉 牛蔵", Area: "新宿", Genres: []string{"焼肉"}, Budget: 5000},
	{ID: "m5", Name: "スパイスカリー ガラム", Area: "新宿", Genres: []string{"カレー", "辛い物"}, Budget: 1500},
	{ID: "m6", Name: "海鮮居酒屋 波", Area: "渋谷", Genres: []string{"和食", "生魚", "居酒屋"}, Budget: 4000},
	{ID: "m7", Name: "餃子の東龍", Area: "渋谷", Genres: []string{"中華", "餃子"}, Budget: 2000},
	{ID: "m8", Name: "ビストロ・ルポ", Area: "渋谷", Genres: []string{"フレンチ", "チーズ"}, Budget: 5500},
	{ID: "m9", Name: "うどん処 こむぎ", Area: "渋谷", Genres: []string{"和食", "うどん"}, Budget: 1000},
	{ID: "m10", Name: "タイ食堂 ガパオ", Area: "渋谷", Genres: []string{"タイ料理", "辛い物", "パクチー"}, Budget: 2500},
}

func (g *mockRestaurantGateway) Search(ctx context.Context, area string) ([]model.Restaurant, error) {
	if area == "" {
		return append([]model.Restaurant(nil), mockRestaurants...), nil
	}
	var result []model.Restaurant
	for _, r := range mockRestaurants {
		if strings.Contains(r.Area, area) {
			result = append(result, r)
		}
	}
	return result, nil
}
