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

## Lint

```bash
cd backend && make lint   # golangci-lint (Docker経由)
cd frontend && npm run lint
```
