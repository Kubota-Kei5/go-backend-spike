package controllers

import (
	"net/http"
	"spike-app/models"

	"github.com/gin-gonic/gin"
)

func Hello(c *gin.Context) {
	c.HTML(http.StatusOK, "hello.html", gin.H{})
}

func NewRecipe(c *gin.Context) {
	c.HTML(http.StatusOK, "new_recipes.html", gin.H{})
}

func ListRecipe(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func TestCreateRecipe(c *gin.Context) {
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	testCreatedRecipe, err := recipe.Create()
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, testCreatedRecipe)
}

func GetRecipe(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe
	if err := models.GetRecipeByID(id, &recipe); err != nil {
		c.JSON(http.StatusNotFound, nil)
		return
	}
	c.JSON(http.StatusOK, recipe)
}