# CD移行作業手順書

## 概要

既存の`scripts/deploy.sh`による手動デプロイ処理を、GitHub ActionsのCD（Continuous Deployment）ワークフローに移行する手順書です。

## 目次

1. [現状分析](#現状分析)
2. [ブランチ戦略とCI/CD設計](#ブランチ戦略とcicd設計)
3. [移行の目的と効果](#移行の目的と効果)
4. [前提条件・準備](#前提条件準備)
5. [実装手順](#実装手順)
6. [設定値・環境変数](#設定値環境変数)
7. [テスト・検証](#テスト検証)
8. [移行後の運用](#移行後の運用)
9. [トラブルシューティング](#トラブルシューティング)

## 現状分析

### 既存の`scripts/deploy.sh`の処理内容

現在のデプロイスクリプトは以下の処理を実行しています：

1. **環境変数読み込み**
   - `.env`ファイルから環境変数を読み込み
   - デフォルト値の設定（PROJECT_NAME、SERVICE_NAME、REGION）

2. **Docker イメージビルド・プッシュ**
   - Google Cloud Buildを使用してDockerイメージをビルド
   - Container Registry（GCR）にイメージをプッシュ
   - タイムアウト設定（600秒）

3. **Cloud Run デプロイ**
   - Cloud Runサービスのデプロイ
   - 環境変数の設定
   - Cloud SQL接続設定
   - タイムアウト設定（300秒）

4. **デプロイ確認**
   - サービスURLの取得と表示

### 現在のCI/CDの状況

- **CI（`.github/workflows/ci.yml`）**
  - `main`、`develop`ブランチへのpush/PRでテスト実行
  - PostgreSQLサービスを使用したDBテスト
  - Go言語のフォーマット、vet、テストを実行

- **CD**
  - 現在は存在せず、手動デプロイのみ
  - 空の`cd.yaml`ファイルが作成済み

## ブランチ戦略とCI/CD設計

### ブランチ戦略

本プロジェクトでは以下の3ブランチでの開発フローを採用：

```
feature branches → develop → main
     ↓             ↓        ↓
   開発作業      統合テスト  本番リリース
```

#### ブランチの役割

| ブランチ | 役割 | CI実行 | CD実行 |
|----------|------|--------|--------|
| `feature/*` | 機能開発・バグ修正 | - | - |
| `develop` | 統合・テスト環境 | ✓ | - |
| `main` | 本番環境 | ✓ | ✓ |

### CI/CD設計

#### CI（Continuous Integration）
- **トリガー**: `develop`, `main` ブランチへの push/merge
- **実行内容**: 
  - Goコードのフォーマットチェック（`go fmt`）
  - 静的解析（`go vet`）
  - 単体・統合テスト（`go test`）
  - PostgreSQL DBを使った統合テスト

#### CD（Continuous Deployment）
- **トリガー**: `main` ブランチへの push/merge のみ
- **実行内容**:
  - Docker イメージビルド
  - Google Container Registry へのプッシュ
  - Cloud Run への自動デプロイ

### 開発フロー

```mermaid
graph TD
    A[feature branch作成] --> B[開発作業]
    B --> C[develop へ PR作成]
    C --> D[CI実行]
    D --> E{CI成功?}
    E -->|No| F[修正]
    F --> D
    E -->|Yes| G[レビュー・マージ]
    G --> H[develop でCI実行]
    H --> I[main へ PR作成]
    I --> J[CI実行]
    J --> K{CI成功?}
    K -->|No| L[修正]
    L --> J
    K -->|Yes| M[レビュー・マージ]
    M --> N[main でCI + CD実行]
    N --> O[本番デプロイ完了]
```

### 推奨設定

#### GitHub Branch Protection Rules

```yaml
# develop ブランチ保護設定
develop:
  - Require pull request reviews before merging
  - Require status checks to pass before merging
    - CI/test (required)
  - Restrict pushes that create files

# main ブランチ保護設定  
main:
  - Require pull request reviews before merging
  - Require status checks to pass before merging
    - CI/test (required)
  - Restrict pushes that create files
  - Require branches to be up to date before merging
```

## 移行の目的と効果

### 目的
- 手動デプロイによる人的ミスの削減
- デプロイ作業の自動化・効率化
- mainブランチマージ時の自動デプロイ実現

### 効果
- **自動化**: mainブランチへのマージ時に自動デプロイ
- **一貫性**: 毎回同じ手順でのデプロイを保証
- **履歴管理**: GitHub Actionsでデプロイ履歴の管理
- **安全性**: CIテスト通過後のみデプロイ実行

## 前提条件・準備

### 1. Google Cloud設定

#### サービスアカウントの作成と権限設定
```bash
# サービスアカウント作成
gcloud iam service-accounts create github-actions-cd \
    --description="Service account for GitHub Actions CD" \
    --display-name="GitHub Actions CD"

# 必要な権限を付与
gcloud projects add-iam-policy-binding PROJECT_ID \
    --member="serviceAccount:github-actions-cd@PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/cloudbuild.builds.builder"

gcloud projects add-iam-policy-binding PROJECT_ID \
    --member="serviceAccount:github-actions-cd@PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/run.admin"

gcloud projects add-iam-policy-binding PROJECT_ID \
    --member="serviceAccount:github-actions-cd@PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/storage.admin"

gcloud projects add-iam-policy-binding PROJECT_ID \
    --member="serviceAccount:github-actions-cd@PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/cloudsql.client"

# サービスアカウントキーの作成
gcloud iam service-accounts keys create key.json \
    --iam-account=github-actions-cd@PROJECT_ID.iam.gserviceaccount.com
```

### 2. GitHub Branch Protection設定

#### Settings > Branches で保護ルールを設定

##### develop ブランチ
1. **Branch name pattern**: `develop`
2. **Protect matching branches** を有効化
3. 以下の設定を有効化：
   - ✅ Require a pull request before merging
   - ✅ Require status checks to pass before merging
     - Required status checks: `test` (CIワークフローのjob名)
   - ✅ Restrict pushes that create files

##### main ブランチ  
1. **Branch name pattern**: `main`
2. **Protect matching branches** を有効化
3. 以下の設定を有効化：
   - ✅ Require a pull request before merging
   - ✅ Require status checks to pass before merging
     - Required status checks: `test` (CIワークフローのjob名)
   - ✅ Require branches to be up to date before merging
   - ✅ Restrict pushes that create files
   - ✅ Require review from code owners (推奨)

### 3. GitHub Secrets設定

GitHubリポジトリの Settings > Secrets and variables > Actions で以下のSecretsを設定：

#### 必須Secrets
- `GCP_PROJECT_ID`: Google CloudプロジェクトID
- `GCP_SA_KEY`: サービスアカウントキー（JSON形式）
- `PROD_DB_USER`: 本番DB用ユーザー名
- `PROD_DB_PASSWORD`: 本番DB用パスワード
- `PROD_DB_NAME`: 本番DB名
- `INSTANCE_CONNECTION_NAME`: Cloud SQLインスタンス接続名
- `DATABASE_URL`: データベース接続URL

#### デプロイ設定用Secrets
- `PROJECT_NAME`: プロジェクト名（デフォルト: spike-backend-gin）
- `SERVICE_NAME`: サービス名（デフォルト: spike-app）
- `REGION`: リージョン（デフォルト: asia-northeast1）

## セキュリティ考慮事項

### 機密情報の取り扱い

#### 問題点
従来の環境変数直接埋め込み方式は以下のリスクがあります：

```yaml
# ❌ 危険: コマンドラインに機密情報が表示される
--set-env-vars="POSTGRES_PASSWORD=${{ secrets.PASSWORD }}"
```

**リスク**:
- GitHub Actionsログに機密情報が表示される可能性
- プロセス一覧やコマンド履歴に残る可能性
- デバッグ時に意図せず露出する可能性

#### 解決策: 環境変数ファイル方式

```yaml
# ✅ 安全: 一時ファイルを使用して機密情報を保護
- name: Create secure environment variables file
  run: |
    cat << 'EOF' > .env.deploy
    POSTGRES_PASSWORD=${{ secrets.PROD_DB_PASSWORD }}
    EOF
    chmod 600 .env.deploy

- name: Deploy
  run: gcloud run deploy --env-vars-file .env.deploy

- name: Cleanup
  if: always()
  run: rm -f .env.deploy
```

**セキュリティ強化ポイント**:
- ✅ ファイル権限を600に制限（所有者のみ読み書き可能）
- ✅ `if: always()`でジョブ失敗時も必ずクリーンアップ
- ✅ ログに機密情報が表示されない
- ✅ ローカル開発環境に影響しない

### その他のセキュリティベストプラクティス

#### 1. サービスアカウント権限の最小化
```bash
# 必要最小限の権限のみ付与
roles/cloudbuild.builds.builder  # ビルド用
roles/run.admin                  # デプロイ用
roles/storage.admin              # イメージ保存用
roles/cloudsql.client           # DB接続用
```

#### 2. ブランチ保護による安全性確保
- mainブランチへの直接pushを禁止
- CIテスト通過必須
- コードレビュー必須

#### 3. デプロイ対象の限定
```yaml
# mainブランチのみデプロイ対象
on:
  push:
    branches: [main]
```

## 実装手順

### 1. CDワークフローファイル作成

`.github/workflows/cd.yaml`を以下の内容で作成：

```yaml
name: CD

on:
  push:
    branches: [main]
  workflow_dispatch:  # 手動実行も可能

jobs:
  # CIワークフローが成功した場合のみデプロイを実行
  deploy:
    runs-on: ubuntu-latest
    needs: []  # CIジョブが存在する場合は依存関係を設定

    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Set up Google Cloud CLI
        uses: google-github-actions/setup-gcloud@v2
        with:
          service_account_key: ${{ secrets.GCP_SA_KEY }}
          project_id: ${{ secrets.GCP_PROJECT_ID }}

      - name: Configure Docker for GCR
        run: gcloud auth configure-docker

      - name: Set environment variables
        run: |
          echo "PROJECT_NAME=${{ secrets.PROJECT_NAME || 'spike-backend-gin' }}" >> $GITHUB_ENV
          echo "SERVICE_NAME=${{ secrets.SERVICE_NAME || 'spike-app' }}" >> $GITHUB_ENV
          echo "REGION=${{ secrets.REGION || 'asia-northeast1' }}" >> $GITHUB_ENV
          echo "IMAGE_TAG=gcr.io/${{ secrets.GCP_PROJECT_ID }}/${{ secrets.SERVICE_NAME || 'spike-app' }}:${{ github.sha }}" >> $GITHUB_ENV

      - name: Build and push Docker image
        run: |
          timeout 600s gcloud builds submit --tag ${{ env.IMAGE_TAG }}

      - name: Create secure environment variables file
        run: |
          cat << 'EOF' > .env.deploy
          ENV=production
          POSTGRES_USER=${{ secrets.PROD_DB_USER }}
          POSTGRES_PASSWORD=${{ secrets.PROD_DB_PASSWORD }}
          POSTGRES_DB=${{ secrets.PROD_DB_NAME }}
          CLOUD_SQL_CONNECTION_NAME=${{ secrets.INSTANCE_CONNECTION_NAME }}
          DATABASE_URL=${{ secrets.DATABASE_URL }}
          EOF
          chmod 600 .env.deploy

      - name: Deploy to Cloud Run
        run: |
          timeout 300s gcloud run deploy ${{ env.SERVICE_NAME }} \
            --image ${{ env.IMAGE_TAG }} \
            --platform managed \
            --region ${{ env.REGION }} \
            --allow-unauthenticated \
            --add-cloudsql-instances ${{ secrets.INSTANCE_CONNECTION_NAME }} \
            --env-vars-file .env.deploy \
            --quiet

      - name: Cleanup environment file
        if: always()
        run: rm -f .env.deploy

      - name: Get service URL
        id: get-url
        run: |
          SERVICE_URL=$(gcloud run services describe ${{ env.SERVICE_NAME }} --region=${{ env.REGION }} --format="value(status.url)")
          echo "service_url=$SERVICE_URL" >> $GITHUB_OUTPUT
          echo "Deployed to: $SERVICE_URL"

      - name: Deployment summary
        run: |
          echo "## Deployment Summary" >> $GITHUB_STEP_SUMMARY
          echo "- **Service**: ${{ env.SERVICE_NAME }}" >> $GITHUB_STEP_SUMMARY
          echo "- **Region**: ${{ env.REGION }}" >> $GITHUB_STEP_SUMMARY
          echo "- **Image**: ${{ env.IMAGE_TAG }}" >> $GITHUB_STEP_SUMMARY
          echo "- **URL**: ${{ steps.get-url.outputs.service_url }}" >> $GITHUB_STEP_SUMMARY
```

### 2. CI/CDの連携設定（オプション）

CIワークフローが成功した場合のみCDを実行する場合は、以下のように設定：

```yaml
# cd.yamlのjobsセクションを以下に変更
jobs:
  # CIワークフローの完了を待機
  wait-for-ci:
    runs-on: ubuntu-latest
    steps:
      - name: Wait for CI
        uses: lewagon/wait-on-check-action@v1.3.1
        with:
          ref: ${{ github.ref }}
          check-name: 'test'  # ci.ymlのjob名
          repo-token: ${{ secrets.GITHUB_TOKEN }}

  deploy:
    runs-on: ubuntu-latest
    needs: wait-for-ci  # CIの完了を待機
    # 以下は同じ
```

## 設定値・環境変数

### GitHub Secrets一覧

| Secret名 | 説明 | 例 | 必須 |
|----------|------|-----|------|
| `GCP_PROJECT_ID` | Google CloudプロジェクトID | `spike-backend-gin` | ✓ |
| `GCP_SA_KEY` | サービスアカウントキー（JSON） | `{"type": "service_account"...}` | ✓ |
| `PROD_DB_USER` | 本番DB用ユーザー名 | `postgres` | ✓ |
| `PROD_DB_PASSWORD` | 本番DB用パスワード | `your-secure-password` | ✓ |
| `PROD_DB_NAME` | 本番DB名 | `spike-app-1-prod` | ✓ |
| `INSTANCE_CONNECTION_NAME` | Cloud SQLインスタンス接続名 | `spike-backend-gin:asia-northeast1:spike-app` | ✓ |
| `DATABASE_URL` | データベース接続URL | `postgresql+pg8000://postgres:password@...` | ✓ |
| `PROJECT_NAME` | プロジェクト名 | `spike-backend-gin` | △ |
| `SERVICE_NAME` | サービス名 | `spike-app` | △ |
| `REGION` | デプロイリージョン | `asia-northeast1` | △ |

### Cloud Run環境変数

デプロイ時に設定される環境変数：

| 変数名 | 値 | 説明 |
|--------|-----|------|
| `ENV` | `production` | 実行環境 |
| `POSTGRES_USER` | GitHub Secrets | DBユーザー名 |
| `POSTGRES_PASSWORD` | GitHub Secrets | DBパスワード |
| `POSTGRES_DB` | GitHub Secrets | DB名 |
| `CLOUD_SQL_CONNECTION_NAME` | GitHub Secrets | Cloud SQL接続名 |
| `DATABASE_URL` | GitHub Secrets | DB接続URL |

## テスト・検証

### 1. 初回デプロイテスト

1. **手動実行でテスト**
   ```bash
   # GitHubのActionsタブから手動実行
   # workflow_dispatch を使用
   ```

2. **mainブランチへのプッシュでテスト**
   ```bash
   git checkout main
   git merge develop
   git push origin main
   ```

### 2. 動作確認

1. **GitHub Actionsの確認**
   - GitHubのActionsタブでワークフロー実行状況を確認
   - ログでエラーがないことを確認

2. **Cloud Runサービスの確認**
   ```bash
   gcloud run services list
   gcloud run services describe spike-app --region=asia-northeast1
   ```

3. **アプリケーションの動作確認**
   - デプロイされたURLにアクセス
   - アプリケーションが正常に動作することを確認

### 3. ロールバック手順

問題が発生した場合のロールバック：

```bash
# 前のリビジョンにロールバック
gcloud run services update-traffic spike-app \
    --to-revisions=REVISION-NAME=100 \
    --region=asia-northeast1
```

## 移行後の運用

### 1. デプロイフロー

#### 標準的な開発・リリースフロー

1. **機能開発**
   ```bash
   # feature ブランチ作成・開発
   git checkout -b feature/new-function
   # 開発作業...
   git push origin feature/new-function
   ```

2. **develop への統合**
   ```bash
   # develop へのPR作成 → CI実行 → レビュー → マージ
   # CI: フォーマット、vet、テスト実行
   ```

3. **本番リリース** 
   ```bash
   # main へのPR作成 → CI実行 → レビュー → マージ
   # CI + CD: テスト実行 → 自動本番デプロイ
   ```

#### ブランチ保護による品質保証

- **feature → develop**: CI必須、レビュー推奨
- **develop → main**: CI必須、レビュー必須、ブランチ最新化必須

#### 緊急リリースフロー

緊急修正が必要な場合：

```bash
# main から hotfix ブランチ作成
git checkout main
git checkout -b hotfix/critical-bug-fix

# 修正作業後、main に直接PR
# CI + CD で即座にデプロイ
```

### 2. 監視・ログ確認

#### GitHub Actions
- ワークフロー実行履歴: `https://github.com/USER/REPO/actions`
- 実行ログ: 各ワークフロー実行の詳細ページ

#### Cloud Run
```bash
# サービス状況確認
gcloud run services describe spike-app --region=asia-northeast1

# ログ確認
gcloud logs read "resource.type=cloud_run_revision AND resource.labels.service_name=spike-app" --limit=50
```

#### Cloud Build
```bash
# ビルド履歴確認
gcloud builds list --limit=10
```

### 3. 緊急時対応

#### デプロイ停止
```bash
# 手動でサービス停止
gcloud run services update spike-app \
    --no-traffic \
    --region=asia-northeast1
```

#### 手動デプロイ
```bash
# 緊急時は既存のdeploy.shも使用可能
./scripts/deploy.sh
```

## トラブルシューティング

### よくある問題と解決方法

#### 1. 認証エラー
**エラー**: `Error: (gcloud.auth.activate-service-account) Could not read json file`

**解決方法**:
- GitHub SecretsのGCP_SA_KEYが正しく設定されているか確認
- サービスアカウントキーのJSON形式が正しいか確認

#### 2. 権限エラー
**エラー**: `Permission denied`

**解決方法**:
- サービスアカウントに必要な権限が付与されているか確認
- Cloud Build、Cloud Run、Storage Admin権限を確認

#### 3. タイムアウトエラー
**エラー**: `Build timed out`

**解決方法**:
- ビルド時間の調整（600s → 900s等）
- Dockerイメージサイズの最適化検討

#### 4. Cloud SQL接続エラー
**エラー**: `Could not connect to Cloud SQL`

**解決方法**:
- INSTANCE_CONNECTION_NAMEの形式確認
- Cloud SQLインスタンスの状態確認
- ネットワーク設定の確認

### デバッグ方法

#### 1. ローカルでのテスト
```bash
# 環境変数設定
export GCP_PROJECT_ID="your-project-id"
# ... 他の環境変数

# gcloudコマンドのテスト
gcloud builds submit --tag gcr.io/$GCP_PROJECT_ID/spike-app:test
```

#### 2. ステップバイステップ実行
```yaml
# cd.yamlでデバッグステップ追加
- name: Debug environment
  run: |
    echo "Project: $GCP_PROJECT_ID"
    echo "Service: $SERVICE_NAME"
    gcloud config list
```

### サポート情報

- **Google Cloud ドキュメント**: https://cloud.google.com/run/docs
- **GitHub Actions ドキュメント**: https://docs.github.com/en/actions
- **本プロジェクトIssue**: リポジトリのIssueページで質問・バグ報告

---

この手順書に従って作業を進めることで、安全かつ確実にCD環境を構築できます。不明な点がある場合は、まず小規模なテスト環境で試すことをお勧めします。