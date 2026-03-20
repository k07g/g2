# g2

Go + クリーンアーキテクチャで実装したタスク管理 Web API です。ECS (Docker コンテナ) へのデプロイを想定しています。

## アーキテクチャ

```
g2/
├── domain/                 # エンティティ・リポジトリインターフェース
│   └── task.go
├── usecase/                # ビジネスロジック
│   └── task_usecase.go
├── handler/                # HTTP ハンドラ（インターフェースアダプタ層）
│   └── task_handler.go
├── infrastructure/
│   └── inmemory/           # インメモリ リポジトリ実装
│       └── task_repository.go
├── main.go
├── Dockerfile
└── Makefile
```

依存の方向は常に内側（domain）へのみ向きます。

```
handler → usecase → domain ← infrastructure
```

## API

| メソッド | パス | 説明 |
|--------|------|------|
| GET | `/tasks` | タスク一覧取得 |
| POST | `/tasks` | タスク作成 |
| GET | `/tasks/{id}` | タスク取得 |
| PUT | `/tasks/{id}` | タスク更新 |
| DELETE | `/tasks/{id}` | タスク削除 |

### タスクのステータス

| 値 | 意味 |
|----|------|
| `todo` | 未着手（デフォルト） |
| `in_progress` | 進行中 |
| `done` | 完了 |

### リクエスト/レスポンス例

**タスク作成**
```bash
curl -XPOST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "買い物", "description": "牛乳を買う"}'
```
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "買い物",
  "description": "牛乳を買う",
  "status": "todo",
  "created_at": "2026-03-20T10:00:00Z",
  "updated_at": "2026-03-20T10:00:00Z"
}
```

**タスク更新**
```bash
curl -XPUT http://localhost:8080/tasks/{id} \
  -H "Content-Type: application/json" \
  -d '{"title": "買い物", "description": "牛乳を買う", "status": "done"}'
```

## ローカル動作確認

Docker が起動している状態で実行します。

```bash
make run
```

別ターミナルでリクエスト:

```bash
# タスク一覧
curl http://localhost:8080/tasks

# タスク作成
curl -XPOST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "タスク1"}'
```

## コマンド

| コマンド | 説明 |
|---------|------|
| `make build` | バイナリをビルド (`./server`) |
| `make test` | テストを実行 |
| `make run` | Docker でローカル起動 (`:8080`) |

## 環境変数

| 変数 | デフォルト | 説明 |
|------|-----------|------|
| `PORT` | `8080` | リッスンするポート番号 |
