# コントリビューションガイド

コントリビューションに興味を持っていただきありがとうございます。
バグ報告・機能提案・コードの改善など、あらゆる形での貢献を歓迎します。

## 行動規範

参加にあたっては [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) を遵守してください。

## Issue

### バグを報告する

[バグ報告テンプレート](https://github.com/k07g/g2/issues/new?template=bug_report.yml)を使用して Issue を作成してください。

### 機能を提案する

[機能リクエストテンプレート](https://github.com/k07g/g2/issues/new?template=feature_request.yml)を使用して Issue を作成してください。
実装前に Issue で方向性を確認することを推奨します。

## Pull Request

### 開発の流れ

```bash
# 1. リポジトリをフォーク・クローン
git clone https://github.com/<your-name>/g2.git
cd g2

# 2. ブランチを作成
git checkout -b feat/your-feature

# 3. 変更・テスト
make test

# 4. コミット・プッシュ
git push origin feat/your-feature

# 5. PR を作成
```

### ブランチ命名規則

| プレフィックス | 用途 |
|---|---|
| `feat/` | 新機能 |
| `fix/` | バグ修正 |
| `refactor/` | リファクタリング |
| `docs/` | ドキュメントのみの変更 |
| `test/` | テストのみの変更 |

### コミットメッセージ

変更内容が簡潔に伝わるメッセージを英語または日本語で記載してください。

```
feat: タスクのページネーションを追加
fix: 存在しないIDを更新した際のエラー処理を修正
```

### PR 提出前のチェックリスト

- `make build` が通る
- `make test` が通る
- 新しいロジックにはテストを追加している

## 開発環境のセットアップ

**必要なもの**

- Go 1.26 以上
- Docker

```bash
# ビルド
make build

# テスト
make test

# Docker でローカル動作確認
make run
```

## アーキテクチャ

このプロジェクトはクリーンアーキテクチャを採用しています。
依存の方向は常に内側（`domain`）へのみ向けてください。

```
handler → usecase → domain ← infrastructure
```

詳細は [README.md](README.md) を参照してください。
