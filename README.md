# FoodLike

好き嫌い配慮 外食提案アプリ

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
