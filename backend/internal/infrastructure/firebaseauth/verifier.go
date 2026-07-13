// Package firebaseauth はFirebase AuthによるIDトークン検証の実装。
package firebaseauth

import (
	"context"
	"strings"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"

	"foodlike-backend/internal/port/gateway"
)

type Verifier struct {
	client *auth.Client
}

// NewVerifier はIDトークン検証用のクライアントを初期化する。
// 検証自体はGoogleの公開鍵で行うが、Go版Admin SDKはクライアント初期化時に
// ADC(GOOGLE_APPLICATION_CREDENTIALS等)を要求する点に注意。
func NewVerifier(ctx context.Context, projectID string) (*Verifier, error) {
	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: projectID})
	if err != nil {
		return nil, err
	}
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}
	return &Verifier{client: client}, nil
}

func (v *Verifier) Verify(ctx context.Context, idToken string) (gateway.AuthToken, error) {
	t, err := v.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return gateway.AuthToken{}, err
	}
	name, _ := t.Claims["name"].(string)
	// メール/パスワード登録は表示名を持たないので、メールのローカル部を初期名にする。
	if name == "" {
		if email, ok := t.Claims["email"].(string); ok {
			name, _, _ = strings.Cut(email, "@")
		}
	}
	return gateway.AuthToken{UID: t.UID, Name: name}, nil
}
