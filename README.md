# FoodLike

**グループ全員の「苦手なもの」を本人以外に明かさないまま、みんなが食べられる外食先を提案するアプリ。**

「パクチーが苦手」「生魚がダメ」——こうした好き嫌いは本人にとって言い出しづらい情報です。FoodLike は各メンバーの苦手（ジャンル・食材）を**個人と紐づかない形に集約**したうえで外食先を検索し、「誰が何を苦手か」を誰にも見せずに全員が安心できるお店を提案します。

> プライバシーが機能そのものであるアプリ。「集約して個人情報を潰す」「LLMに具体名を出させない」といった配慮を、設計と実装の両面で担保しています。

---

## 主な機能

- **グループ／メンバー管理** — QRコード・招待コードで参加
- **苦手なものの登録**（本人のみ編集可・他人には非公開）— ジャンル／食材レベル
- **外食先の提案** — ホットペッパー グルメサーチAPIでエリア×予算から検索し、全員の苦手を避けて絞り込み
- **妥協案ランキング** — 全員OKな店が0件のとき、苦手への抵触が少ない順に提示
- **グループチャット＋AI検索** — チャットから自然文で相談。AIが実在するお店を検索して、みんなの苦手を避けたおすすめを提案
- **通知** — 未読バッジ／通知ベル（可視時のみポーリング）

## 技術スタック

| 領域                     | 技術                                                                |
| ------------------------ | ------------------------------------------------------------------- |
| フロントエンド           | Next.js 16 (App Router) / React 19 / TypeScript                     |
| バックエンド             | Go 1.25 / Gin / GORM                                                |
| アーキテクチャ (backend) | Clean Architecture + Hexagonal (Ports & Adapters)                   |
| 認証                     | Firebase Authentication                                             |
| 本番DB                   | TiDB Cloud Serverless (MySQL互換)                                   |
| ローカル/テストDB        | MySQL 8 (docker compose)                                            |
| 外部グルメAPI            | ホットペッパー グルメサーチAPI                                      |
| AIチャット検索           | Google Gemini (2.5 Flash Lite) + Tavily Search API                  |
| デプロイ                 | フロント: Vercel ／ バック: Google Cloud Run ／ DB: TiDB Serverless |

## 設計のポイント

**プライバシーを確率ではなくコードで保証**
個人の苦手は `PreferenceAggregator` で「誰のか分からない制約集合」に潰してからしか扱わない。AIチャットでは、LLMへの指示（「品目名を書くな」）は確率的で保証にならないため、**生成文から苦手品目名を決定的に伏せる検閲層**（`redact.go`）を後段に置いている。少人数グループでは「集合 − 自分の苦手 = 相方の苦手」と逆算できてしまうため、集約後であっても具体名は出さない。

**認証と認可の分離**
認証（誰か）はミドルウェア、認可（このグループのメンバーか＝BOLA対策）はusecase内で行い、非メンバーには404を返す。usecase/domain層は `gin.Context` を一切知らない。

**差し替え可能なポート設計**
外部依存（グルメAPI・AI・Web検索・認証・永続化）はすべてポート（インターフェース）の背後に置き、APIキーの有無で実装を切り替える。例：`HOTPEPPER_API_KEY` 未設定ならモック店舗、`GEMINI_API_KEY` 未設定なら定型文へフォールバック。テスト時はフェイクを注入。

**AIチャット検索のコスト設計（¥0で実検索）**
GeminiのGoogle検索グラウンディングは課金枠のため、**実店舗の根拠は無料枠のTavily検索で取得 → その結果をGeminiに渡して制約推論＋日本語生成**、という役割分担にした。Tavily（事実）とGemini（推論）を分けることで、無料枠のまま「実在する店を、みんなの苦手を避けて提案」を成立させている。

**リアルタイムはWebSocketではなくポーリング**
Cloud Runのスケールが厳しいので3秒間のポーリングで月の無料枠に収まるようにしている

## ディレクトリ構成

```
backend/    Go (Gin + GORM) / Clean Architecture + Hexagonal
  ├ internal/domain      ドメインモデル・サービス（集約/絞り込み/ランキング）
  ├ internal/usecase     ユースケース（認可・プライバシー境界）
  ├ internal/port        ポート（Repository / Gateway インターフェース）
  ├ internal/adapter     アダプタ（handler / gateway / repository 実装）
  └ internal/infrastructure  DB接続・ルーティング・Firebase
frontend/   Next.js (App Router, TypeScript)
```

## ローカル起動（Docker）

```bash
docker compose up --build
```

- frontend: http://localhost:3000
- backend: http://localhost:8080
- mysql: localhost:3306（ローカル用。本番はTiDB Serverless）

環境変数はバック／フロントで分離：`backend/.env`（compose の `env_file` で注入）と `frontend/.env.local`（`next build` が読む）。各 `*.example` をコピーして値を設定する。APIキー未設定でもモック／フォールバックで動作する。

## ローカル起動（個別）

```bash
# backend
cd backend && cp .env.example .env && go run ./cmd/api

# frontend
cd frontend && cp .env.local.example .env.local && npm install && npm run dev
```

## Lint / Test

```bash
cd backend && make lint      # golangci-lint
cd backend && make test      # go test ./...
cd frontend && npm run lint
```

## デプロイ

Vercel（フロント）+ Cloud Run（バック）+ TiDB Serverless（DB）の無料枠構成。手順は [docs/deploy.md](docs/deploy.md)。

## API エンドポイント

認証が必要なAPIは `Authorization: Bearer <Firebase IDトークン>`（モック認証時は `X-Member-ID: <id>`）を付ける。

| Method   | Path                          | 説明                                             |
| -------- | ----------------------------- | ------------------------------------------------ |
| POST     | `/api/members`                | メンバー登録（認証不要）                         |
| GET      | `/api/me`                     | 自分の情報                                       |
| GET/PUT  | `/api/me/preferences`         | 自分の苦手なもの（本人のみ）                     |
| POST     | `/api/groups`                 | グループ作成                                     |
| POST     | `/api/groups/join`            | 招待コードで参加                                 |
| GET      | `/api/groups`                 | 所属グループ一覧                                 |
| GET      | `/api/groups/:id`             | グループ詳細                                     |
| GET      | `/api/search-options`         | エリア／予算マスタ                               |
| GET      | `/api/groups/:id/suggestions` | 店舗提案（`middle_area`・`budget` で絞り込み可） |
| GET/POST | `/api/groups/:id/messages`    | チャットメッセージ                               |
| GET      | `/api/groups/:id/ai-search`   | AIチャット検索（SSEストリーミング）              |
