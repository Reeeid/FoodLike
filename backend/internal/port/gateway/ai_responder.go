package gateway

import (
	"context"

	"foodlike-backend/internal/domain/model"
)

// AIResponder はAI検索の回答生成ポート。
// 回答はチャンク(数文字〜数語)単位でonChunkに渡し、戻り値で全文を返す。
// 実装をLLM API(SSE)に差し替えてもusecase/handlerは変わらない。
type AIResponder interface {
	Respond(ctx context.Context, query string, candidates []model.Candidate, onChunk func(text string)) (string, error)
}
