# Go Backend Spike

Nutraze のバックエンドを Go にリプレイスするための技術検証用リポジトリ

## Go へのリプレイス内容

|                  | Before                     | After                            |
| ---------------- | -------------------------- | -------------------------------- |
| レシピの栄養計算 | RESTfulAPI で 1 つずつ実行 | gRPC で並列実行                  |
| API              | Next.js (API Routes)       | gin                              |
| ORM              | Prisma                     | GORM + golang-migrate            |
| Testing          |                            | testify                          |
| デプロイ         | すべて Vercel              | Backend: CloudRun, DB: Cloud SQL |

## 技術検証内容

- ~~gin で静的ページルーティング（ダミーデータを使って testify でテスト）~~
- ~~Cloud Run でデプロイ~~
- ~~マルチステージビルドの開発環境構築~~
- ~~docker で db コンテナを作る~~
- ~~gin で db へのデータ登録（フォームに入力した情報を POST）~~
- ~~gin で db からデータ取得（db から GET してページにレンダリング）~~
- ~~Docker DB コンテナ + Gin データベース連携~~
  ~~理由: アプリケーションの基盤機能を先に完成させる~~

  - ~~ローカル開発環境で DB 接続を確立~~
  - ~~CRUD 操作の実装とテスト~~
  - ~~データベーススキーマの設計確定~~

- ~~CI/CD 構築 - 自動テスト~~
  ~~理由: コード品質を保証する仕組みを早期構築~~

  - ~~formatter, 自動テスト（go test）の実施~~
  - ~~データベース機能のテストを含める~~
  - ~~後続作業の品質担保~~

- ~~CI/CD 構築 - 自動ビルド~~
  ~~理由: デプロイの前準備として必要~~

- Cloud SQL デプロイ
  理由: 本番 DB を準備してからデプロイ環境を整える

  - ローカル DB で動作確認済みの状態で移行
  - 接続設定の調整

- CI/CD 構築 - 自動デプロイ
  理由: 全ての要素が揃ってから最終統合
  - DB、アプリ、テストが全て完成した状態
  - 本番環境への安全なデプロイ
- 現状の TS で実装している栄養価計算を見ながら gin にリプレイスしてみる（ダミーデータを使って testify でテスト）
- gRPC を使った並列実装に変更（ダミーデータを使って testify でテスト）

## SETUP と動作確認

### Docker ビルドの実行

```bash
docker compose build web-test
```

## TEST

### コンテナのシェルに入る

```bash
docker compose run --rm web-test bash
```

### テストの実行

```bash
go test ./tests/ -v
```

※ `-v` オプションで詳細なテスト結果が見られる

## マイグレーション

### 開発環境マイグレーション

```bash
# 開発DBを起動
docker compose up db-dev -d

# 開発環境でマイグレーション実行
docker compose run --rm web-dev go run main.go
```

### 本番環境マイグレーション

事前に`.env`ファイルに以下を設定：

```env
PROD_DB_PASSWORD=<Cloud SQLのパスワード>
PROD_DB_USER=postgres
PROD_DB_NAME=spike-app-1-prod
```

マイグレーション実行：

```bash
# Cloud SQL Proxyを起動
docker compose up cloud-sql -d

# 本番環境でマイグレーション実行（環境変数は自動設定）
docker compose run --rm prod-shell ./main
```

### マイグレーション確認

ホストから以下を実行：

```bash
# テーブル一覧確認
docker run --rm -it --net=host postgres:16 psql "host=127.0.0.1 port=5432 sslmode=disable dbname=spike-app-1-prod user=postgres" -c "\dt"

# テーブル構造確認
docker run --rm -it --net=host postgres:16 psql "host=127.0.0.1 port=5432 sslmode=disable dbname=spike-app-1-prod user=postgres" -c "\d recipes"
```

※`Password for user postgres:` を求められるため入力

#### 作成されるテーブル

- `recipes`: レシピ情報
- `ingredients`: 食材情報
- `recipe_ingredients`: レシピと食材の関連テーブル

## Cloud Run デプロイ

### 事前準備

`.env`ファイルに以下の設定があることを確認：

```env
# データベース設定
PROD_DB_PASSWORD=<Cloud SQLのパスワード>
PROD_DB_USER=postgres
PROD_DB_NAME=spike-app-1-prod

# Cloud Run デプロイ設定
PROJECT_NAME=spike-backend-gin
SERVICE_NAME=spike-app
REGION=asia-northeast1
```

### デプロイ実行

```bash
# Cloud Run に自動デプロイ
./scripts/deploy.sh
```

このスクリプトは以下を自動実行します：

1. **Docker イメージビルド**: 本番用イメージを作成
2. **Container Registry プッシュ**: GCR にイメージをアップロード
3. **Cloud Run デプロイ**: サービスを更新
4. **デプロイ確認**: サービス URL を表示

### デプロイ確認

```bash
# サービス一覧確認
docker compose run --rm gcloud gcloud run services list

# 特定サービス詳細確認
docker compose run --rm gcloud gcloud run services describe spike-app --region=asia-northeast1
```
