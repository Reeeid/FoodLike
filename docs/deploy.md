# デプロイ手順(案A: Vercel + Cloud Run + TiDB Serverless / 全部無料枠)

構成は [architecture.md](architecture.md) の構成図①を参照。所要 約1〜1.5時間。

## 事前準備

- Googleアカウント(GCP)・GitHubアカウント・リポジトリをGitHubにpush済みであること
- コード側の必要変更は **DSNのTLS対応の1点だけ**(手順1に記載)

---

## 1. TiDB Cloud Serverless(DB) 約15分

1. https://tidbcloud.com にGitHubでサインアップ → 「Create Cluster」→ **Serverless**(無料)を選択、リージョンは Tokyo (ap-northeast-1)
2. クラスタ作成後、「Connect」→ 接続情報を控える:
   - Host: `gateway01.ap-northeast-1.prod.aws.tidbcloud.com` 系
   - Port: `4000` / User: `xxxx.root` / Password: 生成したもの
3. SQLエディタ(またはmysqlクライアント)で `CREATE DATABASE foodlike;`
4. **コード変更(必須)**: TiDBはTLS必須なので、`backend/internal/infrastructure/db/mysql.go` のDSNをTLS対応にする:

```go
dsn := fmt.Sprintf(
    "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
    ...
)
if os.Getenv("DB_TLS") == "true" {
    dsn += "&tls=true"
}
```

ローカルのdocker composeは`DB_TLS`未設定のまま動き、本番だけ`DB_TLS=true`を渡す。

## 2. Cloud Run(バックエンド) 約30分

**重要: FirebaseプロジェクトfoodlikeauthをそのままGCPプロジェクトとして使う。** 同一プロジェクトならCloud Runのサービスアカウントで自動認証(ADC)されるため、`GOOGLE_APPLICATION_CREDENTIALS`もサービスアカウント鍵ファイルも不要。

```powershell
# gcloud CLI: https://cloud.google.com/sdk/docs/install
gcloud auth login
gcloud config set project foodlikeauth
gcloud services enable run.googleapis.com cloudbuild.googleapis.com artifactregistry.googleapis.com

# backend/ のDockerfileからビルド&デプロイ(ソースデプロイ)
gcloud run deploy foodlike-backend `
  --source backend `
  --region asia-northeast1 `
  --allow-unauthenticated `
  --set-env-vars "FIREBASE_PROJECT_ID=foodlikeauth" `
  --set-env-vars "DB_HOST=<TiDBのHost>,DB_PORT=4000,DB_NAME=foodlike,DB_USER=<User>,DB_TLS=true" `
  --set-env-vars "FRONTEND_ORIGIN=https://<あとで決まるVercelドメイン>" `
  --set-secrets "DB_PASSWORD=foodlike-db-password:latest"
```

- DBパスワードはSecret Manager経由:
  ```powershell
  gcloud services enable secretmanager.googleapis.com
  echo -n "<パスワード>" | gcloud secrets create foodlike-db-password --data-file=-
  # Cloud Runのサービスアカウントに Secret Manager Secret Accessor を付与(初回デプロイ時に案内が出る)
  ```
- デプロイ完了で出るURL(`https://foodlike-backend-xxxx.a.run.app`)を控える
- 動作確認: `curl https://<URL>/health`
- `FRONTEND_ORIGIN`はVercelのドメイン確定後に更新:
  ```powershell
  gcloud run services update foodlike-backend --region asia-northeast1 --update-env-vars "FRONTEND_ORIGIN=https://foodlike.vercel.app"
  ```

補足: `PORT`はCloud Runが自動注入(main.goが読む)。コンテナは0台までスケールインするので無料枠に収まるが、初回アクセスにコールドスタート(数秒)がある。デモ直前に1回アクセスして温めておくと良い。

## 3. Vercel(フロントエンド) 約15分

1. https://vercel.com にGitHubでサインアップ → 「Add New → Project」→ FoodLikeリポジトリをimport
2. **Root Directory を `frontend` に設定**(モノレポなので必須)。FrameworkはNext.jsが自動検出される。`frontend/Dockerfile`は使われない(Vercelネイティブビルド)
3. Environment Variables に以下を設定(値は`frontend/.env.local`と同じ+API URL):

| 変数 | 値 |
|---|---|
| `NEXT_PUBLIC_API_BASE_URL` | Cloud RunのURL |
| `NEXT_PUBLIC_FIREBASE_API_KEY` | .env.localの値 |
| `NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN` | `foodlikeauth.firebaseapp.com` |
| `NEXT_PUBLIC_FIREBASE_PROJECT_ID` | `foodlikeauth` |
| `NEXT_PUBLIC_FIREBASE_APP_ID` | .env.localの値 |

4. Deploy → `https://<プロジェクト名>.vercel.app` が発行される

## 4. 仕上げの相互設定 約10分

1. **Firebaseコンソール** → Authentication → Settings → 承認済みドメイン に Vercelドメイン(`xxx.vercel.app`)を追加(これがないとGoogleログインのポップアップが拒否される)
2. **Cloud Runの`FRONTEND_ORIGIN`** をVercelドメインに更新(手順2の最後のコマンド)
3. 本番で一通り動作確認: ログイン → グループ作成 → QR招待(別端末で読み取り) → チャット → AI検索 → 提案

## トラブルシューティング

| 症状 | 原因/対処 |
|---|---|
| フロントからAPIがCORSエラー | `FRONTEND_ORIGIN`のスキーム/ドメイン不一致(https必須・末尾スラッシュなし) |
| ログインポップアップが閉じる/エラー | Firebase承認済みドメイン未登録 |
| 401 invalid token | Cloud Runのプロジェクトがfoodlikeauthでない(ADC不成立)。`gcloud config get project`確認 |
| DB接続エラー | `DB_TLS=true`漏れ、またはTiDB側のパスワード/Host誤り |
| AI検索が一気に出て流れない | プロキシのバッファリング。Cloud Run直なら発生しない(`X-Accel-Buffering: no`も送信済み) |
