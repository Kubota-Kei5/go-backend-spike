# CLAUDE.md

このファイルは、Claude Code (claude.ai/code) がこのリポジトリでコードを操作する際のガイダンスを提供します。

## プロジェクト概要

Nutazeのバックエンドを Next.js から Go にリプレイスするための技術検証プロジェクトです。使用技術：
- **フレームワーク**: HTTP ルーティングに Gin
- **データベース**: PostgreSQL + GORM ORM
- **テスト**: testify フレームワーク
- **デプロイ**: Cloud Run + Cloud SQL
- **コンテナ化**: マルチステージビルドの Docker

## アーキテクチャ

クリーンアーキテクチャパターンに従ったコードベース構成：

```
spike-app/
├── main.go              # アプリケーションエントリーポイント
├── config/
│   └── database.go      # データベース接続とマイグレーションロジック
├── controllers/
│   ├── router.go        # ルート定義とミドルウェア
│   └── recipe_controller.go # レシピ関連の HTTP ハンドラー
├── models/
│   └── recipe.go        # GORM モデルとデータベース操作
├── templates/           # ビュー用 HTML テンプレート
└── tests/               # テストファイル
```

### 主要コンポーネント

- **データベース設定**: 環境固有の接続文字列を使用。本番環境では `DATABASE_URL` を優先し、ローカル/テスト環境ではコンポーネントベースの DSN にフォールバック
- **自動マイグレーション**: `config.ConnectDatabase()` 経由で起動時に自動的にデータベーススキーマをマイグレーション
- **モデル層**: CRUD 操作のメソッドを持つ GORM モデル
- **コントローラー層**: HTTP リクエストを処理してレスポンスを返す Gin ハンドラー

## 開発コマンド

### ローカル開発

```bash
# 開発環境のビルドと起動
docker compose build web-dev
docker compose up web-dev -d

# 開発用データベース起動
docker compose up db-dev -d

# 開発モードでアプリケーション実行
docker compose run --rm web-dev go run main.go
```

### テスト

```bash
# テストコンテナのシェルに入る
docker compose run --rm web-test bash

# 全テストを詳細出力で実行
go test ./tests/ -v

# ワーキングディレクトリでテスト実行 (spike-app/ から)
go test -v ./tests/...

# 特定のテストファイルを実行
go test ./tests/hello_test.go -v
```

### データベース操作

```bash
# 開発環境マイグレーション（起動時に自動実行）
docker compose run --rm web-dev go run main.go

# Cloud SQL での本番環境マイグレーション
docker compose up cloud-sql -d
docker compose run --rm prod-shell ./main

# 本番データベースへの接続確認
docker run --rm -it --net=host postgres:16 psql "host=127.0.0.1 port=5432 sslmode=disable dbname=spike-app-1-prod user=postgres" -c "\dt"
```

### ビルドとデプロイ

```bash
# 本番イメージのビルド
docker compose build web-prod

# Cloud Run へのデプロイ（.env ファイルが必要）
./scripts/deploy.sh

# 手動 Cloud Run 操作
docker compose run --rm gcloud gcloud run services list
docker compose run --rm gcloud gcloud run services describe spike-app --region=asia-northeast1
```

## 環境設定

アプリケーションは環境変数を使用して設定：

### 開発/テスト環境
- `ENV`: "development", "dev", "test", または "ci"
- `POSTGRES_USER`: データベースユーザー名
- `POSTGRES_PASSWORD`: データベースパスワード
- `POSTGRES_DB`: データベース名
- `PORT`: アプリケーションポート（デフォルト: 8080）

### 本番環境
- `ENV`: "production" または "prod"
- `DATABASE_URL`: 完全な PostgreSQL 接続文字列（優先）
- `PROD_DB_USER`, `PROD_DB_PASSWORD`, `PROD_DB_NAME`: 個別の DB コンポーネント
- `CLOUD_SQL_CONNECTION_NAME`: GCP Cloud SQL インスタンス識別子

## CI/CD パイプライン

### CI (Continuous Integration)
**トリガー条件：**
- mainブランチとdevelopブランチへのpush
- mainブランチとdevelopブランチへのプルリクエスト

**実行内容：**
1. **環境準備**
   - Ubuntu最新版でジョブ実行
   - PostgreSQL 16サービス起動（テスト用DB）
   - Go 1.24セットアップとモジュールキャッシュ

2. **静的解析**
   - `go fmt`でコードフォーマットチェック
   - `go vet`で静的解析実行

3. **テスト実行**
   - 依存関係インストール（`go mod download`）
   - データベース統合テストを含む全テスト実行（`go test -v ./tests/...`）

### CD (Continuous Deployment)
**トリガー条件：**
- mainブランチへのpush時
- 手動実行（`workflow_dispatch`）
- CIワークフロー成功後の自動実行

**実行フロー：**

1. **認証・環境準備**
   - GCPサービスアカウント認証
   - gcloud CLI設定とDocker認証
   - 環境変数設定（プロジェクト名、サービス名、リージョン等）

2. **Cloud SQL デプロイ**
   - インスタンス存在確認（`gcloud sql instances describe`）
   - 存在しない場合：PostgreSQL 16インスタンス自動作成
     - データベース版：POSTGRES_16
     - マシンタイプ：db-f1-micro
     - リージョン：asia-northeast1
     - バックアップ：毎日3:00AM
     - メンテナンス：日曜4:00AM
   - rootパスワード設定
   - 専用DBユーザー作成（存在しない場合）
   - アプリケーション用データベース作成（存在しない場合）

3. **アプリケーションデプロイ**
   - Dockerイメージビルド（本番用マルチステージビルド）
   - Google Container Registryへプッシュ（`gcr.io/{project}/{service}:{commit-sha}`）
   - Cloud Runサービスデプロイ
     - Cloud SQLインスタンス接続設定
     - 本番環境変数注入
     - タイムアウト300秒設定

4. **デプロイ確認**
   - サービスURL取得
   - GitHub Actions サマリー生成
   - デプロイ結果表示

#### 必要なGitHub Secrets
| Secret名 | 説明 | 例 |
|----------|------|-----|
| `GCP_SA_KEY` | Google Cloud サービスアカウントキー（JSON形式） | `{"type": "service_account", ...}` |
| `GCP_PROJECT_ID` | GCPプロジェクトID | `spike-backend-gin` |
| `PROD_DB_USER` | 本番データベースユーザー名 | `spike_user` |
| `PROD_DB_PASSWORD` | 本番データベースパスワード | `secure_password_123` |
| `PROD_DB_NAME` | 本番データベース名 | `spike_prod` |
| `INSTANCE_CONNECTION_NAME` | Cloud SQL接続名 | `spike-backend-gin:asia-northeast1:spike-app` |
| `DATABASE_URL` | 完全なPostgreSQL接続文字列 | `postgresql://user:pass@host/dbname` |
| `PROJECT_NAME` | プロジェクト名（オプション） | `spike-backend-gin` |
| `SERVICE_NAME` | Cloud Runサービス名（オプション） | `spike-app` |
| `REGION` | デプロイリージョン（オプション） | `asia-northeast1` |

#### 必要なサービスアカウント権限
Google Cloudでサービスアカウントに以下のロールを付与：

| ロール名 | 用途 |
|----------|------|
| `Cloud SQL Admin` | Cloud SQLインスタンス、データベース、ユーザーの作成・管理 |
| `Cloud Run Admin` | Cloud Runサービスのデプロイと管理 |
| `Cloud Build Editor` | Dockerイメージのビルド |
| `Storage Admin` | Container Registryへのイメージプッシュ |

#### セットアップ手順

1. **GCPサービスアカウント作成**
   ```bash
   # サービスアカウント作成
   gcloud iam service-accounts create github-actions-sa \
     --description="GitHub Actions CD pipeline" \
     --display-name="GitHub Actions SA"
   
   # 権限付与
   gcloud projects add-iam-policy-binding PROJECT_ID \
     --member="serviceAccount:github-actions-sa@PROJECT_ID.iam.gserviceaccount.com" \
     --role="roles/cloudsql.admin"
   
   gcloud projects add-iam-policy-binding PROJECT_ID \
     --member="serviceAccount:github-actions-sa@PROJECT_ID.iam.gserviceaccount.com" \
     --role="roles/run.admin"
   
   gcloud projects add-iam-policy-binding PROJECT_ID \
     --member="serviceAccount:github-actions-sa@PROJECT_ID.iam.gserviceaccount.com" \
     --role="roles/cloudbuild.builds.editor"
   
   gcloud projects add-iam-policy-binding PROJECT_ID \
     --member="serviceAccount:github-actions-sa@PROJECT_ID.iam.gserviceaccount.com" \
     --role="roles/storage.admin"
   
   # キー生成
   gcloud iam service-accounts keys create key.json \
     --iam-account=github-actions-sa@PROJECT_ID.iam.gserviceaccount.com
   ```

2. **GitHub Secrets設定**
   - リポジトリの Settings > Secrets and variables > Actions
   - 上記表の各Secretを設定

3. **デプロイテスト**
   - mainブランチにpushしてCDパイプライン動作確認

## テスト戦略

- **統合テスト**: CI で PostgreSQL に対してデータベーステストを実行
- **ユニットテスト**: testify を使用したコントローラーとモデルのテスト
- **CI パイプライン**: `go fmt`、`go vet`、完全なテストスイートを実行
- **テスト用データベース**: 分離のため独立したデータベースインスタンスを使用

## データベーススキーマ

コアモデル：
- `Recipe`: ID、Title、Servings、CookingTime を持つレシピ情報
- `Ingredients`: 食材の栄養データ（豊富な栄養素フィールド）
- `RecipeIngredient`: レシピと食材の多対多関係

モデルには CRUD メソッドが含まれ、適切な外部キー制約で GORM アソシエーションを使用。

## コード品質

CI パイプラインで以下を強制：
- `go fmt` によるコードフォーマット
- `go vet` による静的解析
- データベース統合テストを含む全テストの合格
- Go 1.24 でモジュールを有効化