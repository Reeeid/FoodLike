# FoodLike

家族や友人との食事の場で、各メンバーの「好き嫌い」を考慮した外食先を提案するアプリ。

好き嫌いは本人にとって言い出しづらい潜在的な情報であるため、直接聞き出さずに配慮できる仕組みを提供する。グループ内メンバー全員の好み(食材・ジャンルなど)を、本人以外には非公開のまま集約し、エリア×ジャンル×予算で外食先候補を検索・提案する。

## 主要機能(予定)

- グループ／メンバー管理(QRコード・IDによる追加)
- 好き嫌い情報の個人単位登録(本人以外には非公開)
- 外部グルメAPIを用いた外食先(店舗)の検索・提案
- 全員OKな店が無い場合の妥協案ランキング
- (MVP後)チャット形式(WebSocket)での好み抽出・店舗絞り込み

## 技術スタック

### 決定済み

| 領域 | 技術 |
| --- | --- |
| フロントエンド | Next.js (App Router, TypeScript) |
| バックエンド | Go / Gin / GORM |
| アーキテクチャ(backend) | Clean Architecture + Hexagonal(Ports & Adapters) |
| 本番DB | TiDB Cloud Starter(MySQL互換) |
| ローカル/テストDB | MySQL 8 (docker-compose) |
| 認証 | Firebase Auth |
| フロントホスティング | Firebase Hosting(PWA) |
| バックエンドデプロイ先 | Cloud Run |

### 検討中(未決定)

- フロントエンド状態管理(候補: Tanstack Query)
- PWA対応の可否
- デザイン方針
- 外部グルメAPI選定(ホットペッパー / Google Places 等)
- 好き嫌い情報の粒度(食材レベル／ジャンルレベル／自由記述)
- メンバーの集合場所最適化ロジック

### backendの設計方針

- ドメインサービスは単一責務で分割(`PreferenceAggregator` / `RestaurantFilter` / `CompromiseRanker` 等)
- `Repository`(DB永続化)と`Gateway`(外部API呼び出し)をポートとして分離
- Ginは`adapter/handler`・`infrastructure/router`層のみに閉じ込め、usecase/domain層はフレームワークに依存しない
- プライバシー設計: 個人の好み詳細はドメインサービス内で集合演算まで潰し、DTO/APIレスポンス層・LLMには一切渡さない

## 構成

- `backend/` — Go (Gin + GORM) / Clean Architecture + Hexagonal(Ports & Adapters)
- `frontend/` — Next.js (App Router, TypeScript)

## ローカル起動(Docker)

```bash
docker compose up --build
```

- frontend: http://localhost:3000
- backend: http://localhost:8080/health
- mysql: localhost:3306 (テスト用、本番はTiDB Cloudを想定)

## ローカル起動(個別)

### backend

```bash
cd backend
cp .env.example .env
go run ./cmd/api
```

### frontend

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

## Lint / Test

```bash
cd backend && make lint   # golangci-lint (Docker経由)
cd backend && make test
cd frontend && npm run lint
```

## MVPの実装状況

グループ作成 → 招待コード参加 → 好き嫌い登録 → 店舗提案の縦のフローは一通り動く。
ただし**提案の絞り込みロジックは学習用の穴埋め課題**として仮実装のまま残してある。

### 🎓 穴埋め課題(backend/internal/domain/service/)

| 課題 | ファイル | 内容 |
| --- | --- | --- |
| 課題1 | `preference_aggregator.go` | 好き嫌いを個人と紐づかない制約条件に集約 |
| 課題2 | `restaurant_filter.go` | 制約条件で店舗候補を絞り込み |
| 課題3 | `compromise_ranker.go` | 全員OKが0件のときの妥協案ランキング |

- 各ファイルの `【課題N】` コメントに仕様と考えるポイントを記載
- 同ディレクトリの `*_test.go` に期待仕様のテストがある。`t.Skip` の行を消して `go test ./...` が通れば完成
- 現在は仮実装(絞り込まない)なので、嫌いなものを登録しても全店舗が「全員OK」と表示される

### 仮実装(将来差し替え)

- **認証**: `X-Member-ID` ヘッダーの仮認証 → Firebase Auth(issue #8)。差し替えポイントは `backend/internal/adapter/handler/auth.go` と `frontend/src/lib/member-store.ts`
- **店舗検索**: 固定データのモックGateway(`backend/internal/adapter/gateway/mock_restaurant_gateway.go`) → 外部グルメAPI(issue #5)

### APIエンドポイント

| Method | Path | 説明 |
| --- | --- | --- |
| POST | `/api/members` | メンバー登録(認証不要) |
| GET | `/api/me` | 自分の情報 |
| GET/PUT | `/api/me/preferences` | 自分の好き嫌い(本人のみ) |
| POST | `/api/groups` | グループ作成 |
| POST | `/api/groups/join` | 招待コードで参加 |
| GET | `/api/groups` | 所属グループ一覧 |
| GET | `/api/groups/:id` | グループ詳細 |
| GET | `/api/groups/:id/suggestions?area=` | 店舗提案 |

認証が必要なAPIは `X-Member-ID: <id>` ヘッダーを付ける(MVPの仮認証)。
