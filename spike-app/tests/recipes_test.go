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
	db, err := config.ConnectDatabase()
	if err != nil {
		t.Fatal("Failed to connect to test database:", err)
	}

	models.SetDB(db)

	r := controllers.SetupRouter()

	t.Run("アクセスできる", func(t *testing.T) {
		w := testutil.RouterRequest(r, "GET", "/recipes", "")
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("アクセスするとTitleにSpike Appが表示される", func(t *testing.T) {
		w := testutil.RouterRequest(r, "GET", "/recipes", "")
		assert.Contains(t, w.Body.String(), "<title>Spike App</title>")
	})
	t.Run("レシピ一覧というタイトルが表示される", func(t *testing.T) {
		w := testutil.RouterRequest(r, "GET", "/recipes", "")
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "レシピ一覧")
	})
	testRecipe1 := `{"Title": "Test Recipe1", "Servings": 4, "CookingTime": 30}`
	// testRecipe2 := `{"Title": "Test Recipe2", "Servings": 1, "CookingTime": 10}`
	t.Run("postしたrecipe一覧が表示される", func(t *testing.T) {
		testutil.RouterJSONRequest(r, "POST", "/recipes/new", testRecipe1)
		// testutil.RouterJSONRequest(r, "POST", "/recipes", testRecipe2)
		w := testutil.RouterRequest(r, "GET", "/recipes", "")
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Test Recipe1")
		assert.Contains(t, w.Body.String(), "4")
		assert.Contains(t, w.Body.String(), "30")

		// assert.Contains(t, w.Body.String(), "Test Recipe2")
		// assert.Contains(t, w.Body.String(), "1")
		// assert.Contains(t, w.Body.String(), "10")
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