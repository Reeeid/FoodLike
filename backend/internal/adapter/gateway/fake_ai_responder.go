package gateway

import (
	"context"
	"fmt"
	"strings"
	"time"

	"foodlike-backend/internal/domain/model"
)

// FakeAIResponder はAI検索の仮実装。候補リストから定型文を組み立てて
// 少しずつチャンクで流し、ストリーミングUIの動作確認に使う。
// TODO: LLM API(ストリーミング)のGateway実装に差し替える。
type FakeAIResponder struct {
	// ChunkInterval はチャンク間の待ち時間。テストでは0にする。
	ChunkInterval time.Duration
}

func NewFakeAIResponder() *FakeAIResponder {
	return &FakeAIResponder{ChunkInterval: 60 * time.Millisecond}
}

func (g *FakeAIResponder) Respond(ctx context.Context, query string, candidates []model.Candidate, onChunk func(string)) (string, error) {
	full := buildFakeAnswer(query, candidates)

	// 「単語/文字数個」単位で区切って流し、LLMのストリーミングっぽく見せる。
	var sent strings.Builder
	runes := []rune(full)
	const step = 4
	for i := 0; i < len(runes); i += step {
		end := i + step
		if end > len(runes) {
			end = len(runes)
		}
		chunk := string(runes[i:end])
		select {
		case <-ctx.Done():
			return sent.String(), ctx.Err()
		default:
		}
		onChunk(chunk)
		sent.WriteString(chunk)
		if g.ChunkInterval > 0 {
			time.Sleep(g.ChunkInterval)
		}
	}
	return sent.String(), nil
}

func buildFakeAnswer(query string, candidates []model.Candidate) string {
	var b strings.Builder
	fmt.Fprintf(&b, "「%s」ですね!みんなの苦手なものを踏まえて探してみました🍽\n\n", query)
	if len(candidates) == 0 {
		b.WriteString("条件に合うお店が見つかりませんでした…。エリアや条件を変えて、もう一度聞いてみてください。")
		return b.String()
	}
	limit := len(candidates)
	if limit > 3 {
		limit = 3
	}
	for i := 0; i < limit; i++ {
		c := candidates[i]
		mark := "全員が安心して食べられそうです。"
		if !c.MatchedAll {
			mark = "一部苦手なものに触れるかもしれませんが、有力な妥協案です。"
		}
		fmt.Fprintf(&b, "%d. **%s**(%s / ~¥%d / %s) — %s\n",
			i+1, c.Restaurant.Name, c.Restaurant.Area, c.Restaurant.Budget,
			strings.Join(c.Restaurant.Genres, "・"), mark)
	}
	b.WriteString("\n気になるお店はありましたか?エリアや予算を添えて聞いてもらえれば絞り込みます!")
	return b.String()
}
