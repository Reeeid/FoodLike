package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"sync"
	"time"

	"foodlike-backend/internal/domain/model"
)

// hotpepperRestaurantGateway はホットペッパー グルメサーチAPIの実装。
// 無料のリクルートWEBサービス(https://webservice.recruit.co.jp/)のAPIキーが要る。
// RestaurantGateway(店舗検索)とSearchOptionGateway(エリア/予算マスタ)を兼ねる。
//
// 制約: このAPIはジャンル(料理カテゴリ)と予算しか返さず、食材レベルの情報は無い。
// そのため「エビ・パクチーが苦手」のような食材単位の配慮は、実データでは効かない
// (ジャンル単位の配慮のみ有効)。食材レベルの実演はモックデータで行う想定。
type hotpepperRestaurantGateway struct {
	apiKey string
	client *http.Client

	// マスタ(エリア/予算)は静的なので初回取得後にキャッシュする。
	mu      sync.Mutex
	areas   []model.AreaEntry
	budgets []model.BudgetOption
}

func NewHotPepperRestaurantGateway(apiKey string) *hotpepperRestaurantGateway {
	return &hotpepperRestaurantGateway{
		apiKey: apiKey,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

const (
	hotpepperGourmet   = "https://webservice.recruit.co.jp/hotpepper/gourmet/v1/"
	hotpepperSmallArea = "https://webservice.recruit.co.jp/hotpepper/small_area/v1/"
	hotpepperBudget    = "https://webservice.recruit.co.jp/hotpepper/budget/v1/"
)

// getJSON は key/format付きでGETし、レスポンスJSONを out にデコードする。
func (g *hotpepperRestaurantGateway) getJSON(ctx context.Context, endpoint string, extra url.Values, out any) error {
	q := url.Values{}
	q.Set("key", g.apiKey)
	q.Set("format", "json")
	for k, vs := range extra {
		for _, v := range vs {
			if v != "" {
				q.Add(k, v)
			}
		}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"?"+q.Encode(), nil)
	if err != nil {
		return err
	}
	res, err := g.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("hotpepper: unexpected status %d", res.StatusCode)
	}
	return json.NewDecoder(res.Body).Decode(out)
}

// ---- 店舗検索 (RestaurantGateway) ----

type gourmetResponse struct {
	Results struct {
		Shop []struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Genre struct {
				Name string `json:"name"`
			} `json:"genre"`
			SubGenre struct {
				Name string `json:"name"`
			} `json:"sub_genre"`
			Budget struct {
				Average string `json:"average"`
			} `json:"budget"`
			SmallArea struct {
				Name string `json:"name"`
			} `json:"small_area"`
		} `json:"shop"`
	} `json:"results"`
}

func (g *hotpepperRestaurantGateway) Search(ctx context.Context, criteria model.SearchCriteria) ([]model.Restaurant, error) {
	extra := url.Values{}
	extra.Set("count", "20")
	// エリア/予算は指定があるものだけ渡す(空文字はgetJSON側で無視)。
	extra.Set("large_area", criteria.LargeArea)
	extra.Set("middle_area", criteria.MiddleArea)
	extra.Set("small_area", criteria.SmallArea)
	extra.Set("budget", criteria.Budget)

	var body gourmetResponse
	if err := g.getJSON(ctx, hotpepperGourmet, extra, &body); err != nil {
		return nil, err
	}

	restaurants := make([]model.Restaurant, 0, len(body.Results.Shop))
	for _, s := range body.Results.Shop {
		genres := []string{s.Genre.Name}
		if s.SubGenre.Name != "" {
			genres = append(genres, s.SubGenre.Name)
		}
		restaurants = append(restaurants, model.Restaurant{
			ID:   s.ID,
			Name: s.Name,
			Area: s.SmallArea.Name,
			// このAPIはジャンルしか返さないためIngredientsは空。
			// 結果として食材弾きは効かず、ジャンル弾きのみ有効になる。
			Genres: genres,
			Budget: parseBudget(s.Budget.Average),
		})
	}
	return restaurants, nil
}

// ---- マスタ取得 (SearchOptionGateway) ----

type smallAreaResponse struct {
	Results struct {
		SmallArea []struct {
			Code       string `json:"code"`
			Name       string `json:"name"`
			MiddleArea struct {
				Code string `json:"code"`
				Name string `json:"name"`
			} `json:"middle_area"`
			LargeArea struct {
				Code string `json:"code"`
				Name string `json:"name"`
			} `json:"large_area"`
		} `json:"small_area"`
	} `json:"results"`
}

func (g *hotpepperRestaurantGateway) ListAreas(ctx context.Context) ([]model.AreaEntry, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.areas != nil {
		return g.areas, nil
	}
	var body smallAreaResponse
	if err := g.getJSON(ctx, hotpepperSmallArea, url.Values{}, &body); err != nil {
		return nil, err
	}
	areas := make([]model.AreaEntry, 0, len(body.Results.SmallArea))
	for _, s := range body.Results.SmallArea {
		areas = append(areas, model.AreaEntry{
			LargeCode:  s.LargeArea.Code,
			LargeName:  s.LargeArea.Name,
			MiddleCode: s.MiddleArea.Code,
			MiddleName: s.MiddleArea.Name,
			SmallCode:  s.Code,
			SmallName:  s.Name,
		})
	}
	g.areas = areas
	return areas, nil
}

type budgetResponse struct {
	Results struct {
		Budget []struct {
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"budget"`
	} `json:"results"`
}

func (g *hotpepperRestaurantGateway) ListBudgets(ctx context.Context) ([]model.BudgetOption, error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.budgets != nil {
		return g.budgets, nil
	}
	var body budgetResponse
	if err := g.getJSON(ctx, hotpepperBudget, url.Values{}, &body); err != nil {
		return nil, err
	}
	budgets := make([]model.BudgetOption, 0, len(body.Results.Budget))
	for _, b := range body.Results.Budget {
		budgets = append(budgets, model.BudgetOption{Code: b.Code, Name: b.Name})
	}
	g.budgets = budgets
	return budgets, nil
}

var digitsRe = regexp.MustCompile(`\d+`)

// parseBudget は "2001～3000円" のような予算文字列から代表値(最初の数値)を取る。
// 取れなければ0(予算不明)。
func parseBudget(s string) int {
	m := digitsRe.FindString(s)
	if m == "" {
		return 0
	}
	n, _ := strconv.Atoi(m)
	return n
}
