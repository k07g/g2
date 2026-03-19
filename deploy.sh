#!/bin/bash
set -euo pipefail

# ---- 設定 ----
AWS_REGION="${AWS_REGION:-ap-northeast-1}"
AWS_ACCOUNT_ID="${AWS_ACCOUNT_ID:?必須: AWS_ACCOUNT_ID を設定してください}"
ECR_REPO="${ECR_REPO:-g2-lambda}"
LAMBDA_FUNCTION="${LAMBDA_FUNCTION:-g2}"
IMAGE_TAG="${IMAGE_TAG:-latest}"

ECR_URI="${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${ECR_REPO}"

echo "==> ECR ログイン"
aws ecr get-login-password --region "${AWS_REGION}" \
  | docker login --username AWS --password-stdin "${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com"

echo "==> ECR リポジトリ作成 (存在しない場合)"
aws ecr describe-repositories --repository-names "${ECR_REPO}" --region "${AWS_REGION}" > /dev/null 2>&1 \
  || aws ecr create-repository --repository-name "${ECR_REPO}" --region "${AWS_REGION}"

echo "==> Docker ビルド"
docker build --platform linux/amd64 -t "${ECR_REPO}:${IMAGE_TAG}" .

echo "==> イメージにタグ付け"
docker tag "${ECR_REPO}:${IMAGE_TAG}" "${ECR_URI}:${IMAGE_TAG}"

echo "==> ECR へプッシュ"
docker push "${ECR_URI}:${IMAGE_TAG}"

echo "==> Lambda 関数の更新 (存在しない場合は作成)"
if aws lambda get-function --function-name "${LAMBDA_FUNCTION}" --region "${AWS_REGION}" > /dev/null 2>&1; then
  echo "  -> 既存の Lambda を更新"
  aws lambda update-function-code \
    --function-name "${LAMBDA_FUNCTION}" \
    --image-uri "${ECR_URI}:${IMAGE_TAG}" \
    --region "${AWS_REGION}"
else
  echo "  -> 新規 Lambda を作成"
  echo "  ※ LAMBDA_ROLE_ARN を設定してください"
  aws lambda create-function \
    --function-name "${LAMBDA_FUNCTION}" \
    --package-type Image \
    --code "ImageUri=${ECR_URI}:${IMAGE_TAG}" \
    --role "${LAMBDA_ROLE_ARN:?必須: LAMBDA_ROLE_ARN を設定してください}" \
    --region "${AWS_REGION}"
fi

echo "==> デプロイ完了: ${ECR_URI}:${IMAGE_TAG}"
