// Package invitecode はグループ招待コードの生成。
// 純粋計算・決定的なロジックはinterface化せず/pkgに共通化する方針に従う。
package invitecode

import (
	"crypto/rand"
	"encoding/hex"
)

// New は16文字のランダムな招待コードを生成する。
func New() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
