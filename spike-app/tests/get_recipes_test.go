package tests

import (
	"net/http"
	"spike-app/controllers"
	testutil "spike-app/tests/utils"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)


func Test_recipeエンドポイントに(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	r := controllers.SetupRouter()
	w := testutil.RouterRequest(r, "GET", "/recipes")

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