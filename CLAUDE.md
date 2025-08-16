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