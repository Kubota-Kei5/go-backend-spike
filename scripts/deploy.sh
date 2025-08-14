#!/bin/bash

# Cloud Run デプロイスクリプト

set -e

# 環境変数読み込み
if [ -f .env ]; then
    source ./.env
else
    echo "Error: .env file not found"
    exit 1
fi

# デフォルト値設定
PROJECT_NAME=${PROJECT_NAME:-spike-backend-gin}
SERVICE_NAME=${SERVICE_NAME:-spike-app}
REGION=${REGION:-asia-northeast1}
IMAGE_TAG_NAME=gcr.io/${PROJECT_NAME}/${SERVICE_NAME}

# Docker イメージビルド・プッシュ
timeout 600s docker compose run --rm gcloud gcloud builds submit --tag ${IMAGE_TAG_NAME}
BUILD_RESULT=$?

if [ $BUILD_RESULT -eq 124 ]; then
    echo "Error: Build timed out"
    exit 1
elif [ $BUILD_RESULT -ne 0 ]; then
    echo "Error: Build failed"
    exit 1
fi

# Cloud Run デプロイ
timeout 300s docker compose run --rm gcloud gcloud run deploy ${SERVICE_NAME} \
    --image ${IMAGE_TAG_NAME} \
    --platform managed \
    --region ${REGION} \
    --allow-unauthenticated \
    --add-cloudsql-instances ${INSTANCE_CONNECTION_NAME} \
    --set-env-vars="ENV=production,POSTGRES_USER=${PROD_DB_USER},POSTGRES_PASSWORD=${PROD_DB_PASSWORD},POSTGRES_DB=${PROD_DB_NAME},CLOUD_SQL_CONNECTION_NAME=${INSTANCE_CONNECTION_NAME},DATABASE_URL=${DATABASE_URL}" \
    --quiet

DEPLOY_RESULT=$?

if [ $DEPLOY_RESULT -eq 124 ]; then
    echo "Error: Deploy timed out"
    exit 1
elif [ $DEPLOY_RESULT -ne 0 ]; then
    echo "Error: Deploy failed"
    exit 1
fi

# サービスURL取得・表示
SERVICE_URL=$(docker compose run --rm gcloud gcloud run services describe ${SERVICE_NAME} --region=${REGION} --format="value(status.url)" 2>/dev/null | tr -d '\r')
echo "Deployed: ${SERVICE_URL}"