package gateway

import (
	"context"

	"foodlike-backend/internal/domain/model"
)

// WebSearcher はWeb検索のポート。AIチャット検索で実在店舗の根拠を得るために使う。
// LLM側のグラウンディング機能(有料枠)の代わりに、外部の検索APIへ差し替えられる。
type WebSearcher interface {
	Search(ctx context.Context, query string) ([]model.WebSearchResult, error)
}
