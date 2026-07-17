package gateway

import (
	"context"

	"foodlike-backend/internal/domain/model"
)

// AIResponder はAI検索の回答生成ポート。
// 回答はチャンク(数文字〜数語)単位でonChunkに渡し、戻り値で全文を返す。
// 渡すのは集約済みの制約(誰の苦手かは潰した匿名集合)のみ。実店舗の検索は
// 実装側(Gemini+Google検索)に委ねる。実装を差し替えてもusecase/handlerは変わらない。
type AIResponder interface {
	Respond(ctx context.Context, query string, constraints model.Constraints, onChunk func(text string)) (string, error)
}
