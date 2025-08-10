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
- gin で db へのデータ登録（フォームに入力した情報を POST）
- gin で db からデータ取得（db から GET してページにレンダリング）
- Docker DB コンテナ + Gin データベース連携
  理由: アプリケーションの基盤機能を先に完成させる

  - ローカル開発環境で DB 接続を確立
  - CRUD 操作の実装とテスト
  - データベーススキーマの設計確定

- CI/CD 構築 - 自動テスト
  理由: コード品質を保証する仕組みを早期構築

  - データベース機能のテストを含める
  - 後続作業の品質担保

- CI/CD 構築 - 自動ビルド
  理由: デプロイの前準備として必要

  - Docker イメージの自動生成
  - アーティファクトの管理

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
