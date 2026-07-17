package gateway

import (
	"context"
	"fmt"
	"time"

	"foodlike-backend/internal/domain/model"
)

// FakeAIResponder はAI検索の代替(フォールバック)実装。GEMINI_API_KEY未設定や
// Gemini呼び出し失敗時に使う。自前で店舗検索はできないので、正直な案内文を
// ストリーミングで流してチャットを止めない。
type FakeAIResponder struct {
	// ChunkInterval はチャンク間の待ち時間。テストでは0にする。
	ChunkInterval time.Duration
}

func NewFakeAIResponder() *FakeAIResponder {
	return &FakeAIResponder{ChunkInterval: 60 * time.Millisecond}
}

func (g *FakeAIResponder) Respond(ctx context.Context, query string, _ model.Constraints, onChunk func(string)) (string, error) {
	return streamInChunks(ctx, buildFakeAnswer(query), g.ChunkInterval, onChunk)
}

func buildFakeAnswer(query string) string {
	return fmt.Sprintf(
		"「%s」ですね!ただ今AI検索につながりませんでした🙏 "+
			"少し時間をおいてもう一度試すか、下の「お店を提案してもらう」から"+
			"エリアや予算を選んで探してみてください。みんなの苦手なものには配慮して探します。",
		query)
}
