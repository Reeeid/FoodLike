package gateway

import "context"

// AuthToken は検証済みIDトークンから取り出した本人情報。
type AuthToken struct {
	UID  string // 認証プロバイダ上のユーザーID(Firebase UID)
	Name string // 表示名(未設定なら空)
}

// AuthVerifier はIDトークン検証のポート。
// 実装(Firebase Admin SDK)はinfrastructure側に置き、
// handler層はこのinterfaceにだけ依存する。
type AuthVerifier interface {
	Verify(ctx context.Context, idToken string) (AuthToken, error)
}
