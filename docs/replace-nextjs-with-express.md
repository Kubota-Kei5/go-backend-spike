# Nutraze フロントエンド・バックエンドリプレイス手順書

## 0. 前提
- 既存: Next.js でフロント＆バックエンドを構成
- 変更後: フロントエンド (Next.js) + バックエンド (Express+Node.js)
- モノレポ構成は継続

## 1. リポジトリ構成の整理

### 1.1 ディレクトリ構成
```
project-root/
├── apps/
│   ├── frontend/          # Next.js フロントエンド
│   └── backend/           # Express バックエンド
├── packages/
│   └── shared/           # 共通コード・型定義
├── package.json          # ルートパッケージ設定
└── turbo.json           # Turbo設定（モノレポ管理）
```

### 1.2 共通コードの移行
- API型定義を `packages/shared/types` へ移動
- バリデーションスキーマを `packages/shared/schemas` へ移動
- 共通ユーティリティを `packages/shared/utils` へ移動

### 1.3 tsconfig設定
```json
{
  "compilerOptions": {
    "baseUrl": ".",
    "paths": {
      "@shared/*": ["packages/shared/*"],
      "@frontend/*": ["apps/frontend/*"],
      "@backend/*": ["apps/backend/*"]
    }
  }
}
```

## 2. Express バックエンドの実装

### 2.1 初期セットアップ
```bash
# バックエンドディレクトリ作成
mkdir -p apps/backend
cd apps/backend

# 依存関係インストール
npm init -y
npm install express cors helmet morgan
npm install -D @types/express @types/cors typescript ts-node nodemon
```

### 2.2 API ルーティングの移行
既存のNext.js API Routes（`pages/api/` または `app/api/`）をExpressルートに移行：

```typescript
// apps/backend/src/routes/recipes.ts
import express from 'express';
import { PrismaClient } from '@prisma/client';

const router = express.Router();
const prisma = new PrismaClient();

// 旧: pages/api/recipes/index.ts
router.get('/recipes', async (req, res) => {
  const recipes = await prisma.recipe.findMany();
  res.json(recipes);
});

// 旧: pages/api/recipes/[id].ts
router.get('/recipes/:id', async (req, res) => {
  const recipe = await prisma.recipe.findUnique({
    where: { id: parseInt(req.params.id) }
  });
  res.json(recipe);
});

export default router;
```

### 2.3 Expressサーバー設定
```typescript
// apps/backend/src/app.ts
import express from 'express';
import cors from 'cors';
import helmet from 'helmet';
import morgan from 'morgan';
import recipeRoutes from './routes/recipes';

const app = express();
const PORT = process.env.PORT || 4000;

// ミドルウェア
app.use(helmet());
app.use(cors({
  origin: process.env.FRONTEND_URL || 'http://localhost:3000'
}));
app.use(morgan('combined'));
app.use(express.json());

// ルート
app.use('/api', recipeRoutes);

app.listen(PORT, () => {
  console.log(`Backend server running on port ${PORT}`);
});
```

### 2.4 開発用スクリプト
```json
{
  "scripts": {
    "dev": "nodemon --exec ts-node src/app.ts",
    "build": "tsc",
    "start": "node dist/app.js"
  }
}
```

## 3. フロントエンドの更新

### 3.1 API呼び出し先変更
```typescript
// apps/frontend/src/lib/api.ts
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:4000/api';

export const fetchRecipes = async () => {
  const response = await fetch(`${API_BASE_URL}/recipes`);
  return response.json();
};

export const fetchRecipe = async (id: string) => {
  const response = await fetch(`${API_BASE_URL}/recipes/${id}`);
  return response.json();
};
```

### 3.2 環境変数設定
```bash
# apps/frontend/.env.local
NEXT_PUBLIC_API_URL=http://localhost:4000/api

# apps/frontend/.env.production
NEXT_PUBLIC_API_URL=https://your-backend-domain.com/api
```

### 3.3 共通ライブラリ参照
```typescript
// apps/frontend/src/types/recipe.ts
export { Recipe, CreateRecipeInput } from '@shared/types';
```

## 4. Prisma/Postgres の設定変更

### 4.1 Prismaクライアント移行
```bash
# Prismaをバックエンドに移行
mv prisma apps/backend/
cd apps/backend
npm install @prisma/client prisma
```

### 4.2 環境変数整理
```bash
# apps/backend/.env
DATABASE_URL="postgresql://username:password@localhost:5432/nutraze_db"
JWT_SECRET="your-secret-key"
```

### 4.3 マイグレーション実行
```bash
# バックエンドディレクトリで実行
cd apps/backend
npx prisma migrate dev
npx prisma generate
```

### 4.4 セキュリティ設定
- データベース接続をバックエンドのみに制限
- フロントエンドから直接DB接続を削除
- 認証・認可ミドルウェアの実装

## 5. デプロイ構成

### 5.1 フロントエンド（Vercel）
```json
// vercel.json
{
  "builds": [
    {
      "src": "apps/frontend/package.json",
      "use": "@vercel/next"
    }
  ],
  "env": {
    "NEXT_PUBLIC_API_URL": "@api-url"
  }
}
```

### 5.2 バックエンド選択肢

#### Option A: Vercel Functions
```typescript
// api/index.ts (Vercel Functions用)
import app from '../apps/backend/src/app';
export default app;
```

#### Option B: AWS App Runner / Cloud Run
```dockerfile
# apps/backend/Dockerfile
FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY . .
RUN npm run build
EXPOSE 4000
CMD ["npm", "start"]
```

### 5.3 環境変数管理
- フロントエンド: Vercelの環境変数設定
- バックエンド: 各ホスティングサービスの環境変数設定
- 本番用DATABASE_URLの設定

## 6. テストおよびリリース

### 6.1 ユニットテストの更新
```typescript
// apps/backend/src/tests/recipes.test.ts
import request from 'supertest';
import app from '../app';

describe('GET /api/recipes', () => {
  test('should return recipes list', async () => {
    const response = await request(app).get('/api/recipes');
    expect(response.status).toBe(200);
    expect(Array.isArray(response.body)).toBe(true);
  });
});
```

### 6.2 統合テストの実装
- フロントエンド-バックエンド間のAPI通信テスト
- データベース接続テスト
- 認証フローテスト

### 6.3 リリース手順
1. **準備フェーズ**
   - 機能凍結・コードレビュー完了
   - ステージング環境での動作確認
   - パフォーマンステスト実施

2. **デプロイフェーズ**
   - バックエンドを先にデプロイ
   - フロントエンドのAPI URL更新
   - DNS設定・SSL証明書確認

3. **検証フェーズ**
   - 本番環境での動作確認
   - エラー監視・ログ確認
   - パフォーマンス監視

### 6.4 ロールバック計画
- バックエンドAPIの下位互換性確保
- フロントエンドの段階的デプロイ
- データベーススキーマ変更の慎重な管理
- 監視アラートと自動ロールバック設定

## 7. 移行後の運用

### 7.1 監視・ログ
- アプリケーションログの集約
- エラー監視（Sentry等）
- パフォーマンス監視（New Relic等）

### 7.2 スケーリング考慮
- バックエンドの水平スケーリング対応
- データベース接続プールの最適化
- CDN活用によるフロントエンド高速化

### 7.3 セキュリティ強化
- API レート制限の実装
- CORS設定の最適化
- セキュリティヘッダーの強化