# recipes/newデータ登録修正計画

## 問題の分析

### 現在の問題
本番環境で`/recipes/new`からフォーム送信すると`null`が返されてDBに登録されない。

### 原因
**HTMLフォームとGoコントローラーのデータ構造・バインド方式の不一致**

#### 1. データ形式の不一致
- **HTMLフォーム**: `application/x-www-form-urlencoded`で送信
- **コントローラー**: `c.ShouldBindJSON(&recipe)`でJSON形式を想定

#### 2. フィールド名の不一致
- **HTMLフォーム**: `name="cook_time"`
- **Go構造体**: `CookingTime`（JSONタグのみ）

#### 3. エラーハンドリング不備
- バインディング失敗時に`null`を返すだけでエラー詳細が不明

## 修正計画

### 1. モデル構造体修正（models/recipe.go）
```go
type Recipe struct {
    ID          uint   `gorm:"primaryKey" json:"ID" form:"id"`
    Title       string `gorm:"not null" json:"Title" form:"title"`
    Servings    int    `gorm:"not null" json:"Servings" form:"servings"`
    CookingTime int    `gorm:"not null" json:"CookingTime" form:"cook_time"`
}
```

**変更点**:
- `form`タグを追加してHTMLフォームのname属性と対応
- `cook_time` → `CookingTime`のマッピング

### 2. コントローラー修正（controllers/recipe_controller.go）
```go
func CreateRecipe(c *gin.Context) {
    var recipe models.Recipe
    
    // フォームデータとJSONの両方に対応
    if err := c.ShouldBind(&recipe); err != nil {
        log.Printf("Failed to bind recipe data: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid form data", 
            "details": err.Error(),
        })
        return
    }
    
    log.Printf("Received recipe data: %+v", recipe)
    
    createdRecipe, err := recipe.Create()
    if err != nil {
        log.Printf("Failed to create recipe: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to create recipe",
            "details": err.Error(),
        })
        return
    }
    
    log.Printf("Successfully created recipe: %+v", createdRecipe)
    c.JSON(http.StatusOK, createdRecipe)
}
```

**変更点**:
- `ShouldBindJSON` → `ShouldBind`（フォーム・JSON両対応）
- 詳細なログ追加
- エラーレスポンスの改善

### 3. テンプレート修正（オプション）
レスポンス後の処理改善：
```html
<form action="/recipes/new" method="post" onsubmit="handleSubmit(event)">
  <!-- 既存のフォーム内容 -->
</form>

<script>
async function handleSubmit(event) {
    event.preventDefault();
    
    const formData = new FormData(event.target);
    const data = Object.fromEntries(formData.entries());
    
    try {
        const response = await fetch('/recipes/new', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                Title: data.title,
                Servings: parseInt(data.servings),
                CookingTime: parseInt(data.cook_time)
            })
        });
        
        const result = await response.json();
        
        if (response.ok) {
            alert('レシピが登録されました！');
            window.location.href = '/recipes';
        } else {
            alert('エラー: ' + result.error);
        }
    } catch (error) {
        alert('通信エラー: ' + error.message);
    }
}
</script>
```

## 修正手順

### ステップ1: モデル修正
1. `spike-app/models/recipe.go`の`Recipe`構造体に`form`タグ追加

### ステップ2: コントローラー修正
1. `spike-app/controllers/recipe_controller.go`の`CreateRecipe`関数を修正
2. `ShouldBindJSON` → `ShouldBind`に変更
3. ログ出力とエラーハンドリング改善

### ステップ3: ローカルテスト
```bash
# 開発環境で動作確認
docker compose up web-dev db-dev -d
docker compose run --rm web-dev go run main.go

# ブラウザでテスト: http://localhost:8000/recipes/new
```

### ステップ4: デプロイ・本番テスト
```bash
# mainブランチにpushしてCD実行
git add .
git commit -m "Fix recipes/new form data binding issue"
git push origin main

# デプロイ後、本番環境でテスト
```

## 検証ポイント

1. **フォーム送信成功**: レシピがDBに正常登録される
2. **エラー表示改善**: 問題発生時に詳細なエラーメッセージ表示
3. **ログ出力**: Cloud Runログでデバッグ情報確認可能
4. **既存機能**: JSON APIとしても引き続き動作

## 追加調査項目（必要に応じて）

1. **Cloud Runログ確認**:
   ```bash
   gcloud run services logs read spike-app --region=asia-northeast1
   ```

2. **データベース接続確認**:
   ```bash
   # Cloud SQL接続テスト
   gcloud sql connect spike-app --user=postgres
   ```

3. **環境変数確認**:
   Cloud Runの環境変数が正しく設定されているかチェック