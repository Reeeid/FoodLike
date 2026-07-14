# FoodLike アーキテクチャ構成図

Mermaid記法。GitHub / VS Code(拡張) でそのままレンダリングされ、レポートにも流用できる。

## 1. システム構成図(デプロイ構成 / 案A: 全フリーティア)

```mermaid
flowchart LR
    subgraph Client["ユーザーのブラウザ / スマホ"]
        UI["Next.js アプリ<br/>(React 19)"]
    end

    subgraph Vercel["Vercel (無料)"]
        FE["フロントエンド<br/>Next.js 16"]
    end

    subgraph Firebase["Firebase (無料)"]
        AUTH["Firebase Authentication<br/>Google / メール / 匿名ゲスト"]
    end

    subgraph GCP["Google Cloud Run (無料枠)"]
        BE["バックエンドAPI<br/>Go + Gin"]
    end

    subgraph TiDB["TiDB Cloud Serverless (無料)"]
        DB[("MySQL互換DB<br/>members / groups /<br/>preferences / messages")]
    end

    UI -->|"静的配信 / SSR"| FE
    UI -->|"① ログイン"| AUTH
    AUTH -->|"② IDトークン(JWT)"| UI
    UI -->|"③ REST + SSE<br/>Authorization: Bearer"| BE
    BE -->|"④ トークン検証<br/>(Admin SDK)"| AUTH
    BE -->|"⑤ GORM (TLS)"| DB
```

- 全構成が無料枠に収まる。AI検索もLLM APIを使わず、ルールベースのフェイク実装(コスト0)で提供する
- Cloud RunをFirebaseと同一GCPプロジェクト(foodlikeauth)に置くことで、サービスアカウント鍵ファイル不要(ADCで検証)

## 2. バックエンドのレイヤードアーキテクチャ

依存の向きは常に外側→内側。ドメイン層は他の層を一切知らない。

```mermaid
flowchart TB
    subgraph Infra["infrastructure 層"]
        ROUTER["router<br/>(ルーティング/CORS)"]
        MYSQL["db<br/>(MySQL接続/Migrate)"]
        FBAUTH["firebaseauth.Verifier<br/>(IDトークン検証)"]
    end

    subgraph Adapter["adapter 層 (ポートの実装)"]
        HANDLER["handler<br/>Gin ハンドラ / DTO / 認証MW"]
        REPO["repository<br/>GORM実装 + InMemoryMessage(仮)"]
        GW["gateway<br/>MockRestaurant / FakeAIResponder"]
    end

    subgraph Usecase["usecase 層 (アプリケーションルール)"]
        UC["Member / Group / Preference /<br/>Suggestion / Chat の各Usecase<br/>※メンバーシップ検査(BOLA対策)はここ"]
    end

    subgraph Domain["domain 層 (ビジネスルール)"]
        MODEL["model: Group / Member /<br/>Preference / Message / Candidate"]
        SVC["service: 好み集約 /<br/>絞り込み / 妥協案ランキング"]
    end

    subgraph Port["port 層 (インターフェース)"]
        PREPO["repository.*Repository"]
        PGW["gateway.RestaurantGateway /<br/>AIResponder / AuthVerifier"]
    end

    ROUTER --> HANDLER
    HANDLER --> UC
    UC --> PREPO
    UC --> PGW
    UC --> MODEL
    UC --> SVC
    REPO -.実装.-> PREPO
    GW -.実装.-> PGW
    FBAUTH -.実装.-> PGW
    MYSQL --> REPO
```

差し替え可能ポイント(ポートの実装を入れ替えるだけで、usecase以下は無変更):

| ポート | 現在の実装 | 将来の実装 |
|---|---|---|
| `MessageRepository` | インメモリ(仮) | GORM + messagesテーブル |
| `RestaurantGateway` | モック店舗データ | ホットペッパーAPI |
| `AIResponder` | ルールベース(テンプレート文生成) | LLM API(必要になったら) |
| `AuthVerifier` | Firebase Admin SDK | (モック認証MWと起動時に切替) |

## 3. 認証フロー(JITメンバー登録)

```mermaid
sequenceDiagram
    participant U as ブラウザ
    participant F as Firebase Auth
    participant B as バックエンド
    participant D as DB

    U->>F: ログイン(Google/メール/匿名)
    F-->>U: IDトークン(JWT)
    U->>B: GET /api/me (Bearer IDトークン)
    B->>F: トークン検証(Admin SDK)
    F-->>B: UID / 表示名
    B->>D: FirebaseUIDでメンバー検索
    alt 初回ログイン
        B->>D: メンバーをJIT作成
    end
    B-->>U: 内部メンバー(id, name)
```

## 4. チャット/AI検索フロー(SSEストリーミング)

```mermaid
sequenceDiagram
    participant U as ブラウザ(ChatCard)
    participant B as バックエンド
    participant D as DB(メッセージ)

    Note over U,B: 通常チャット: POST + 3秒ポーリング
    U->>B: POST /groups/:id/messages
    B->>D: 発言を保存
    loop 3秒ごと(全メンバー)
        U->>B: GET /messages?after_id=N
        B-->>U: 新着のみ返す(未読バッジにも利用)
    end

    Note over U,B: AI検索: SSEでストリーミング
    U->>B: GET /groups/:id/ai-search?q=... (fetchストリーム)
    B->>B: メンバーシップ検査
    B->>D: 質問を発言として保存
    B->>B: 好み集約→候補絞り込み(全員の苦手を考慮)
    loop 回答生成
        B-->>U: event: chunk(数文字ずつ)
    end
    B->>D: 回答をAIメッセージとして保存
    B-->>U: event: done(確定メッセージ)
    Note over U,D: 以降は他メンバーもポーリングで同じ回答を閲覧(API呼び出しは1回だけ)
```
