package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"foodlike-backend/internal/domain/model"
)

const tavilySearchEndpoint = "https://api.tavily.com/search"

// TavilyWebSearcher は Tavily(LLM向け検索API)による WebSearcher 実装。
// GeminiのGoogle検索グラウンディングが有料枠のため、無料枠(月1,000検索)で
// 実在店舗の根拠を得る代替として使う。標準のnet/httpのみで依存を増やさない。
type TavilyWebSearcher struct {
	apiKey string
	client *http.Client
}

func NewTavilyWebSearcher(apiKey string) *TavilyWebSearcher {
	return &TavilyWebSearcher{
		apiKey: apiKey,
		// 検索はチャット応答の前段。長く待たせないよう短めに切る。
		client: &http.Client{Timeout: 8 * time.Second},
	}
}

type tavilyRequest struct {
	Query       string `json:"query"`
	SearchDepth string `json:"search_depth"`
	MaxResults  int    `json:"max_results"`
}

type tavilyResponse struct {
	Results []struct {
		Title   string `json:"title"`
		URL     string `json:"url"`
		Content string `json:"content"`
	} `json:"results"`
}

func (g *TavilyWebSearcher) Search(ctx context.Context, query string) ([]model.WebSearchResult, error) {
	body, err := json.Marshal(tavilyRequest{
		Query:       query,
		SearchDepth: "basic", // 無料枠を節約(1クレジット/回)
		MaxResults:  5,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tavilySearchEndpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+g.apiKey)
	req.Header.Set("Content-Type", "application/json")

	res, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tavily search: unexpected status %d", res.StatusCode)
	}

	var parsed tavilyResponse
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	results := make([]model.WebSearchResult, 0, len(parsed.Results))
	for _, r := range parsed.Results {
		results = append(results, model.WebSearchResult{
			Title:   r.Title,
			URL:     r.URL,
			Content: r.Content,
		})
	}
	return results, nil
}
