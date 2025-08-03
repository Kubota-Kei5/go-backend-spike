package main

import (
	"spike-app/controllers"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	
	// Load HTML templates
	r.LoadHTMLGlob("templates/*")

	// Set up routes
	r.GET("/hello", controllers.Hello)

	// Uncomment the following lines to add recipe routes
	// r.GET("/recipes", controllers.ListRecipes)
	// r.GET("/recipes/new", controllers.NewRecipe)
	// r.POST("/recipes/new", controllers.CreateRecipe)

	return r
}

func main() {
	Router().Run(":8080")
}