# g2

Go で書いた AWS Lambda コンテナデプロイのサンプルです。

## 構成

```
g2/
├── main.go       # Lambda ハンドラ
├── Dockerfile    # Lambda コンテナイメージ (golang:1.26.1-alpine → provided:al2023)
├── deploy.sh     # ECR + Lambda デプロイスクリプト
├── Makefile      # ビルド/テスト/デプロイコマンド
├── go.mod
└── go.sum
```

## API

**リクエスト**
```json
{ "name": "Lambda" }
```

**レスポンス**
```json
{ "message": "Hello, Lambda!", "statusCode": 200 }
```

`name` を省略すると `"Hello, World!"` が返ります。

## ローカル動作確認

Docker が起動している状態で実行します。

```bash
# コンテナをビルドして Lambda Runtime Interface Emulator (RIE) で起動
make run
```

別ターミナルで呼び出し:

```bash
make invoke
# => { "message": "Hello, Lambda!", "statusCode": 200 }
```

## AWS へのデプロイ

### 前提条件

- AWS CLI が設定済み (`aws configure`)
- Docker が起動中
- Lambda 実行用 IAM ロール (`AWSLambdaBasicExecutionRole` 以上)

### デプロイ手順

```bash
export AWS_ACCOUNT_ID=123456789012
export LAMBDA_ROLE_ARN=arn:aws:iam::123456789012:role/lambda-execution-role

# 任意: デフォルト値から変更する場合
export AWS_REGION=ap-northeast-1   # デフォルト: ap-northeast-1
export ECR_REPO=g2-lambda          # デフォルト: g2-lambda
export LAMBDA_FUNCTION=g2          # デフォルト: g2

make deploy
```

`deploy.sh` は以下を自動で行います:

1. ECR へログイン
2. ECR リポジトリを作成 (初回のみ)
3. Docker イメージをビルド (`linux/amd64`)
4. ECR へプッシュ
5. Lambda 関数を更新、または新規作成

### Lambda の動作確認

```bash
aws lambda invoke \
  --function-name g2 \
  --payload '{"name":"Lambda"}' \
  --cli-binary-format raw-in-base64-out \
  response.json && cat response.json
```
