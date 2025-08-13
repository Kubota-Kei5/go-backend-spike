#!/bin/bash

# Cloud Run デプロイスクリプト
# 使用方法: ./scripts/deploy.sh

set -e  # エラー時に停止

# .env ファイルから環境変数を読み込む
if [ -f .env ]; then
    source ./.env
else
    echo "Error: .env file not found. Please create .env file with required variables."
    exit 1
fi

# 必要な環境変数の設定（デフォルト値含む）
PROJECT_NAME=${PROJECT_NAME:-spike-backend-gin}
SERVICE_NAME=${SERVICE_NAME:-spike-app}
REGION=${REGION:-asia-northeast1}

# Container Registry のイメージタグ
IMAGE_TAG_NAME=gcr.io/${PROJECT_NAME}/${SERVICE_NAME}

echo "================================"
echo "Cloud Run Deployment Script"
echo "================================"
echo "Project: ${PROJECT_NAME}"
echo "Service: ${SERVICE_NAME}"
echo "Region: ${REGION}"
echo "Image: ${IMAGE_TAG_NAME}"
echo "================================"

# ステップ1: Docker イメージをビルドしてContainer Registryにプッシュ
echo "Step 1: Building and pushing Docker image..."
docker compose run --rm gcloud gcloud builds submit --tag ${IMAGE_TAG_NAME}

if [ $? -ne 0 ]; then
    echo "Error: Failed to build and push Docker image"
    exit 1
fi

echo "✅ Docker image built and pushed successfully"

# ステップ2: Cloud Run サービスをデプロイ
echo "Step 2: Deploying to Cloud Run..."
docker compose run --rm gcloud gcloud run deploy ${SERVICE_NAME} \
    --image ${IMAGE_TAG_NAME} \
    --platform managed \
    --region ${REGION} \
    --allow-unauthenticated \
    --add-cloudsql-instances ${INSTANCE_CONNECTION_NAME} \
    --set-env-vars="ENV=production,POSTGRES_USER=${PROD_DB_USER},POSTGRES_PASSWORD=${PROD_DB_PASSWORD},POSTGRES_DB=${PROD_DB_NAME},CLOUD_SQL_CONNECTION_NAME=${INSTANCE_CONNECTION_NAME},DATABASE_URL=${DATABASE_URL}"

if [ $? -ne 0 ]; then
    echo "Error: Failed to deploy to Cloud Run"
    exit 1
fi

echo "✅ Cloud Run deployment completed successfully"

# ステップ3: デプロイ状況確認
echo "Step 3: Verifying deployment..."
docker compose run --rm gcloud gcloud run services describe ${SERVICE_NAME} --region=${REGION} --format="value(status.url)"

echo "================================"
echo "Deployment completed! 🎉"
echo "Your service is now available at the URL above."
echo "================================"