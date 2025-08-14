package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"spike-app/controllers"
	testutil "spike-app/tests/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Helper function to perform HTTP requests for testing
func performRequest(r *gin.Engine, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func Test_helloエンドポイントに(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := controllers.SetupRouter()
	w := testutil.RouterRequest(r, "GET", "/hello", "")

	t.Run("アクセスできる", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, w.Code)
	})

	body := w.Body.String()

	t.Run("アクセスするとbodyにHello Worldが表示される", func(t *testing.T) {
		assert.Contains(t, body, "<h1>Hello World!</h1>")
	})

}
