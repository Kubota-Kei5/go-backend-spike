package tests

import (
	"net/http"
	"spike-app/config"
	"spike-app/controllers"
	"spike-app/models"
	testutil "spike-app/tests/utils"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)


func Test_recipeエンドポイントに(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	r := controllers.SetupRouter()
	w := testutil.RouterRequest(r, "GET", "/recipes", "")

	t.Run("アクセスできる", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, w.Code)
	})

	body := w.Body.String()

	t.Run("アクセスするとTitleにSpike Appが表示される", func(t *testing.T) {
		assert.Contains(t, body, "<title>Spike App</title>")
	})
	t.Run("アクセスするとbodyのh1にWelcome to Spike Appが表示される", func(t *testing.T) {
		assert.Contains(t, body, "<h1>Welcome to Spike App</h1>")
	})
	t.Run("アクセスするとbodyのCreate Recipeボタンが表示される", func(t *testing.T) {
		assert.Contains(t, body, "Create Recipe")
	})
}


func Test_recipes_newエンドポイントに(t *testing.T) {
	gin.SetMode(gin.TestMode)

	db, err := config.ConnectDatabase()
	if err != nil {
		t.Fatal("Failed to connect to test database:", err)
	}

	models.SetDB(db)

	r := controllers.SetupRouter()

	testRecipe := `{"Title": "Test Recipe", "Servings": 4, "CookingTime": 30}`

	t.Run("testRecipeをpostすると200が返ってくる", func(t *testing.T) {
		w := testutil.RouterJSONRequest(r, "POST", "/recipes/new", testRecipe)
		assert.Equal(t, http.StatusOK, w.Code)
	})
	t.Run("testRecipeをpostするとDBにレシピが保存される", func(t *testing.T) {
		var recipe models.Recipe
		if err := db.First(&recipe, "title = ?", "Test Recipe").Error; err != nil {
			t.Fatal("Failed to find recipe in database:", err)
		}
		assert.Equal(t, "Test Recipe", recipe.Title)
		assert.Equal(t, 4, recipe.Servings)
		assert.Equal(t, 30, recipe.CookingTime)
	})

	t.Run("testRecipeをgetで取得することができる", func(t *testing.T) {
		w := testutil.RouterRequest(r, "GET", "/recipes/1", "")
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Test Recipe")
		assert.Contains(t, w.Body.String(), "4")
		assert.Contains(t, w.Body.String(), "30")
	})
}