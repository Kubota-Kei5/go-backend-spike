package controllers

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("../templates/*")
	r.GET("/hello", Hello)

	// r.GET("/recipes", ListRecipes)
	// r.GET("/recipes/new", NewRecipe)
	// r.POST("/recipes/new", CreateRecipe)

	return r
}