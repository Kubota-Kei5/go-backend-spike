package controllers

import (
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
	r := gin.Default()

	r.LoadHTMLGlob(getTemplatePath())
	r.GET("/hello", Hello)

	r.GET("/recipes", ListRecipe)
	r.GET("/recipes/new", NewRecipe)
	r.POST("/recipes/new", TestCreateRecipe)
	r.GET("/recipes/:id", GetRecipe)
	// r.POST("/recipes/new", CreateRecipe)

	return r
}