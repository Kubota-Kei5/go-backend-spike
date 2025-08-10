package tests

import (
	"net/http"
	"os"
	"spike-app/config"
	"spike-app/controllers"
	"spike-app/models"
	testutil "spike-app/tests/utils"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)


func Test_recipes_newエンドポイントに(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Setup test environment variables to use db-dev container  
	os.Setenv("POSTGRES_USER", "spike")
	os.Setenv("POSTGRES_PASSWORD", "spike")
	os.Setenv("POSTGRES_DB", "spike_dev")
	os.Setenv("ENV", "development")

	// Setup database connection for testing
	db, err := config.ConnectDatabase()
	if err != nil {
		t.Fatal("Failed to connect to test database:", err)
	}

	// Set the database connection for models
	models.SetDB(db)

	r := controllers.SetupRouter()

	testRecipe := `{"Title": "Test Recipe", "Servings": 4, "CookingTime": 30}`

	t.Run("testRecipeをpostすると200が返ってくる", func(t *testing.T) {
		w := testutil.RouterJSONRequest(r, "POST", "/recipes/new", testRecipe)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}