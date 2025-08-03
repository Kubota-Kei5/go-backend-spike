package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"spike-app/controllers/router"

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

func TestHelloEndpoint(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	
	// Use the SetupRouter function from router package
	r := router.SetupRouter()
	
	// Perform request
	w := performRequest(r, "GET", "/hello")
	
	// Check status code
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Check response body contains expected HTML
	body := w.Body.String()
	
	// Verify HTML structure
	assert.Contains(t, body, "<html")
	assert.Contains(t, body, "</html>")
	assert.Contains(t, body, "<h1>Hello World!</h1>")
	
	// Check Content-Type header
	assert.Contains(t, w.Header().Get("Content-Type"), "text/html")
}

func TestHelloEndpointResponseStructure(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)
	
	// Use the SetupRouter function from router package
	r := router.SetupRouter()
	
	// Perform request
	w := performRequest(r, "GET", "/hello")
	
	// Check status code
	assert.Equal(t, http.StatusOK, w.Code)
	
	body := w.Body.String()
	
	// Test HTML structure more specifically
	assert.True(t, strings.Contains(body, "<!DOCTYPE html"))
	assert.True(t, strings.Contains(body, "<title>Spike App</title>"))
	assert.True(t, strings.Contains(body, "<body>"))
	assert.True(t, strings.Contains(body, "</body>"))
	
	// Ensure the h1 tag is in the body section
	bodyStart := strings.Index(body, "<body>")
	bodyEnd := strings.Index(body, "</body>")
	assert.True(t, bodyStart < bodyEnd)
	
	bodyContent := body[bodyStart:bodyEnd]
	assert.Contains(t, bodyContent, "<h1>Hello World!</h1>")
}