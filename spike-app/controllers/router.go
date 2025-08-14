package controllers

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

// テスト時: tests/からtemplates/を探すため../templates/*が必要
// 本番時: spike-app/からtemplates/を探すためtemplates/*が必要
func getTemplatePath() string {
	if _, err := os.Stat("templates"); err == nil {
		return "templates/*"
	}
	return "../templates/*"
}

func SetupRouter() *gin.Engine {
	r := gin.New()  // Use gin.New() instead of gin.Default()
	
	// Custom recovery middleware
	r.Use(func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("PANIC: %v", r)
				c.JSON(500, gin.H{"error": "Internal Server Error", "detail": r})
				c.Abort()
			}
		}()
		c.Next()
	})

	// Debug middleware
	r.Use(func(c *gin.Context) {
		log.Printf("Request: %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
	})

	r.LoadHTMLGlob(getTemplatePath())
	r.GET("/hello", Hello)

	r.GET("/recipes", ListRecipe)
	r.GET("/recipes/new", NewRecipe)
	r.POST("/recipes/new", CreateRecipe)
	r.GET("/recipes/:id", GetRecipe)

	return r
}