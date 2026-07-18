# 開発でハマったことまとめ

デプロイ〜運用で踏んだ罠と対処。デプロイ手順そのものは [deploy.md](deploy.md) を参照。

---

## フロントエンド(Vercel)

### 全ルートが 404(ビルドは成功しているのに)

- **症状**: `next build` は成功(ログ全部緑)なのに、本番URLの `/` すべてが 404。`curl -sI` すると `X-Vercel-Error: NOT_FOUND` / `text/plain`(Next.js の not-found ページではなく Vercel プラットフォームの404)。
- **原因**: `frontend/next.config.ts` の `output: "standalone"`。これは Docker等で自前 `server.js` を起動する self-host 用。Vercel は自前 server.js を起動せず、独自にサーバーレス関数＋静的アセットへ変換してルーティングを組むため、standalone だと出力がひもづかず全ルート404になる。
- **対処**: Vercel はビルド時に `VERCEL=1` を立てるので、そのときだけ無効化して Docker と両立:
  ```ts
  output: process.env.VERCEL ? undefined : "standalone",
  ```
- **教訓**: 「ビルド成功なのに404」は配信・ルーティング側の問題。まず `next.config` の `output` / `basePath` を疑い、`curl -sI` の `X-Vercel-Error` ヘッダで切り分ける。

### 本番URLの取り違え

- `<project>.vercel.app` はグローバルで早い者勝ち。想定した名前が取られてると別名が割り当たる。正規URLは **Vercel の Project → Settings → Domains** で確認する。
- CLIで明示的に本番デプロイするなら `vercel --prod`(git連携がPreview止まりのときの回避にもなる)。

---

## バックエンド(Render / Cloud Run)

### Cloud Run は Billing(クレカ)必須

- 無料枠内で使う場合でも、`run.googleapis.com` 等の有効化に **Billing アカウントのリンクが必須**。未設定だと `FAILED_PRECONDITION: Billing ... not found` で止まる。
- クレカを使いたくない → **Render** に変更(Docker常駐・無料枠・クレカ不要)。`backend/Dockerfile` をそのまま使える。

### Firebase 認証が起動時に ADC(鍵)を要求する

- **症状**: `firebase.NewApp` が認証情報を探し、ADC が無い環境(Render等)だと初期化に失敗 → `log.Fatalf` でクラッシュ。
- **原因**: IDトークン**検証**自体は Google の公開鍵だけで完結するのに、Go版 Admin SDK はデフォルトで管理者認証(サービスアカウント鍵)を初期化時に要求する。Cloud Run では同一プロジェクトのADCが自動で刺さるので気づかなかった。
- **対処**: 検証しかしないなら `option.WithoutAuthentication()` を渡す。鍵ファイル不要になる(カスタムトークン発行など管理者操作は使えなくなるが、このアプリは検証のみ)。
  ```go
  firebase.NewApp(ctx, &firebase.Config{ProjectID: projectID}, option.WithoutAuthentication())
  ```
- 名前・UID は「トークンの claims」と「自前DB」から取っているので、鍵を消しても機能は落ちない。

### `gcloud run deploy` の落とし穴

- **`PORT` を `--set-env-vars` で渡すとエラー**。Cloud Run の予約語で自動注入される。アプリ側は `os.Getenv("PORT")` を読む。
- **`--set-secrets "KEY=値"` は不可**。`--set-secrets` は `KEY=シークレット名:バージョン` 構文。生値は `--set-env-vars` に入れるか、Secret Manager に登録してから参照する。
- **`--set-env-vars` を複数回書くと最後だけ残る**ことがある。1つにカンマ区切りでまとめるのが安全。
- **値に改行を混ぜない**。PowerShellで複数行の `"..."` にすると改行がキー名に紛れて壊れる。

---

## DB(TiDB Serverless)

- **TLS 必須**。DSN に `&tls=true` が要る。`DB_TLS=true` のときだけ付与する分岐にして、ローカルMySQL(TLSなし)と両立。
- コンテナに `ca-certificates` が入っていないと TLS 検証に失敗する(`backend/Dockerfile` は導入済み)。
- Host の末尾(`...tidbcloud.com`)や DB名(既定 `test`)の取り違えに注意。

---

## CORS / 環境変数

- CORS 許可オリジンは `FRONTEND_ORIGIN`(単一)。**末尾スラッシュを付けると不一致でCORSエラー**(`https://xxx.vercel.app` が正、`.../` は誤り)。
- 未設定時は `http://localhost:3000` にフォールバック。ローカルはローカルのバックエンド、本番は本番を叩くので、通常は片方でよい。ローカルフロント→本番バックを叩きたい場合だけ複数オリジン対応が必要。
- `NEXT_PUBLIC_*` は**ビルド時にJSへ焼き込まれる**公開値。変更したら再デプロイが要る。

---

## ローカル環境

- `gcloud` は scoop の `extras` bucket にある。未追加だと `Couldn't find manifest for 'gcloud'`。→ `scoop bucket add extras` してから `scoop install gcloud`。
