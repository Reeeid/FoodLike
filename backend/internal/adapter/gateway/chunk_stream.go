package gateway

import (
	"context"
	"strings"
	"time"
)

// streamInChunks は完成済みのテキストを数文字ずつonChunkへ流す(疑似ストリーミング)。
//
// LLMの生成をそのまま垂れ流さずここを通すのは、送信済みの文字は取り消せないため。
// 検閲(redactExcluded)を通した後の安全なテキストだけを流す必要がある。
// 本物のトークンストリームより初回表示は遅いが、プライバシー保証と引き換えにしている。
func streamInChunks(ctx context.Context, text string, interval time.Duration, onChunk func(string)) (string, error) {
	var sent strings.Builder
	runes := []rune(text)
	const step = 4
	for i := 0; i < len(runes); i += step {
		end := i + step
		if end > len(runes) {
			end = len(runes)
		}
		select {
		case <-ctx.Done():
			return sent.String(), ctx.Err()
		default:
		}
		chunk := string(runes[i:end])
		onChunk(chunk)
		sent.WriteString(chunk)
		if interval > 0 {
			time.Sleep(interval)
		}
	}
	return sent.String(), nil
}
