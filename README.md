# Go Backend Spike

Nutraze のバックエンドを Go にリプレイスするための技術検証用リポジトリ

## Go へのリプレイス内容

|                  | Before                     | After                            |
| ---------------- | -------------------------- | -------------------------------- |
| レシピの栄養計算 | RESTfulAPI で 1 つずつ実行 | gRPC で並列実行                  |
| API              | Next.js (API Routes)       | gin                              |
| ORM              | Prisma                     | GORM                             |
| Testing          |                            | testify                          |
| デプロイ         | すべて Vercel              | Backend: CloudRun, DB: Cloud SQL |

## 技術検証内容

- ~~gin で静的ページルーティング（ダミーデータを使って testify でテスト）~~
- ~~Cloud Run でデプロイ~~
- ~~マルチステージビルドの開発環境構築~~
- docker で db コンテナを作り、gin で db へのデータ登録（フォームページの情報を POST）とデータ取得（db から GET してページにレンダリング）
- Cloud SQL で db デプロイ
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
