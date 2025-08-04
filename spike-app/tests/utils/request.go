package testutil

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func RouterRequest(r *gin.Engine, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}