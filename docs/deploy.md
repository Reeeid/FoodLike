# デプロイ手順(案B: Vercel + Render + TiDB Serverless / 全部無料枠・クレカ不要)

構成は [architecture.md](architecture.md) の構成図①を参照。所要 約1時間。

> 案A(Cloud Run)は Billing(クレカ)登録が必須のため、クレカ不要で常駐サーバーを
> 動かせる **Render** に変更。バックエンドは `backend/Dockerfile` をそのまま使う。

## 事前準備

- Googleアカウント(Firebase)・GitHubアカウント・リポジトリをGitHubにpush済みであること
- コード変更は不要(TiDBのTLS対応・Firebase認証のADC不要化はコミット済み)

---

## 1. TiDB Cloud Serverless(DB) 約15分

1. https://tidbcloud.com にGitHubでサインアップ → 「Create Cluster」→ **Serverless**(無料)を選択、リージョンは Tokyo (ap-northeast-1)
2. クラスタ作成後、「Connect」→ 接続情報を控える:
   - Host: `gateway01.ap-northeast-1.prod.aws.tidbcloud.com` 系
   - Port: `4000` / User: `xxxx.root` / Password: 生成したもの
3. SQLエディタ(またはmysqlクライアント)で `CREATE DATABASE foodlike;`(既定の `test` DBを使うならスキップ可)

補足: TiDBはTLS必須。`backend/internal/infrastructure/db/mysql.go` が `DB_TLS=true` のときDSNに `&tls=true` を付与する(コミット済み)。ローカルの docker compose は `DB_TLS` 未設定のまま動き、本番だけ `DB_TLS=true` を渡す。

## 2. Render(バックエンド) 約20分

Firebaseの認証情報(サービスアカウント鍵)は**不要**。`firebaseauth.NewVerifier` が `option.WithoutAuthentication()` で初期化し、IDトークン検証はプロジェクトIDと公開鍵だけで完結する(コミット済み)。

1. https://render.com にGitHubでサインアップ(**クレカ不要**)
2. **New → Web Service** → FoodLikeリポジトリを接続
3. 設定:
   - **Root Directory**: `backend`
   - **Runtime**: Docker(`backend/Dockerfile` を自動検出)
   - **Dockerfile Path**: `./Dockerfile`(見つからない場合は `./backend/Dockerfile`)
   - **Region**: Singapore(日本から最も近い無料リージョン)
   - **Instance Type**: Free
4. **Environment** に環境変数を設定(`PORT` は入れない=Renderが自動注入し `main.go` が読む):

   | キー | 値 |
   |---|---|
   | `FIREBASE_PROJECT_ID` | `foodlikeauth` |
   | `GEMINI_API_KEY` | (取得したキー) |
   | `GEMINI_MODEL` | `gemini-2.5-flash-lite` |
   | `GEMINI_GROUNDING` | `false` |
   | `TAVILY_API_KEY` | (取得したキー) |
   | `HOTPEPPER_API_KEY` | (ホットペッパー検索を使う場合) |
   | `DB_HOST` | `gateway01.ap-northeast-1.prod.aws.tidbcloud.com` |
   | `DB_PORT` | `4000` |
   | `DB_NAME` | `foodlike`(または `test`) |
   | `DB_USER` | `xxxx.root` |
   | `DB_PASSWORD` | (TiDBのパスワード) |
   | `DB_TLS` | `true` |
   | `FRONTEND_ORIGIN` | `https://<Vercelドメイン>`(**末尾スラッシュなし**) |

5. Deploy → 発行URL(`https://foodlike-backend-xxxx.onrender.com`)を控える
6. 動作確認: `curl https://<URL>/health`

補足: Render無料枠は **15分無操作でスリープ**し、次アクセスでコールドスタート(30〜50秒)。デモ直前に1回アクセスして温めておくと良い。`FRONTEND_ORIGIN` はVercelドメイン確定後に更新する(手順4)。

## 3. Vercel(フロントエンド) 約15分

1. https://vercel.com にGitHubでサインアップ → 「Add New → Project」→ FoodLikeリポジトリをimport
2. **Root Directory を `frontend` に設定**(モノレポなので必須)。FrameworkはNext.jsが自動検出される。`frontend/Dockerfile`は使われない(Vercelネイティブビルド)
3. Environment Variables に以下を設定(値は`frontend/.env.local`と同じ+API URL):

| 変数 | 値 |
|---|---|
| `NEXT_PUBLIC_API_BASE_URL` | RenderのURL(手順2) |
| `NEXT_PUBLIC_FIREBASE_API_KEY` | .env.localの値 |
| `NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN` | `foodlikeauth.firebaseapp.com` |
| `NEXT_PUBLIC_FIREBASE_PROJECT_ID` | `foodlikeauth` |
| `NEXT_PUBLIC_FIREBASE_APP_ID` | .env.localの値 |

4. Deploy → `https://<プロジェクト名>.vercel.app` が発行される

## 4. 仕上げの相互設定 約10分

1. **Firebaseコンソール** → Authentication → Settings → 承認済みドメイン に Vercelドメイン(`xxx.vercel.app`)を追加(これがないとGoogleログインのポップアップが拒否される)
2. **Renderの`FRONTEND_ORIGIN`** をVercelドメインに更新(Environment を編集して保存 → 自動で再デプロイ)
3. 本番で一通り動作確認: ログイン → グループ作成 → QR招待(別端末で読み取り) → チャット → AI検索 → 提案

## トラブルシューティング

| 症状 | 原因/対処 |
|---|---|
| フロントからAPIがCORSエラー | `FRONTEND_ORIGIN`のスキーム/ドメイン不一致(https必須・末尾スラッシュなし) |
| ログインポップアップが閉じる/エラー | Firebase承認済みドメイン未登録 |
| 401 invalid token | `FIREBASE_PROJECT_ID` がフロントのプロジェクトと不一致 |
| 起動時にfirebase初期化エラー | 古いコードだとADCを要求する。`option.WithoutAuthentication()` 適用済みか確認 |
| DB接続エラー | `DB_TLS=true`漏れ、またはTiDB側のパスワード/Host誤り |
| 初回アクセスが遅い | Render無料枠のコールドスタート(スリープ復帰)。デモ前に温める |
| AI検索が一気に出て流れない | プロキシのバッファリング(`X-Accel-Buffering: no`送信済み)。通常Renderでは発生しない |
