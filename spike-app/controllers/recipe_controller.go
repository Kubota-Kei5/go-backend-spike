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
	recipes, err := models.GetAllRecipes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	c.HTML(http.StatusOK, "index.html", gin.H{"recipes": recipes})
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

type FormRecipeInfo struct {
	Title                string  `form:"title" binding:"required"`
	Servings             int     `form:"servings" binding:"required"`
	CookingTime          int     `form:"cooking_time" binding:"required"`
	IngredientNames      []string `form:"ingredient_names" binding:"required"`
	IngredientQuantities []int    `form:"ingredient_quantities" binding:"required"`
}

func CreateRecipe(c *gin.Context) {
	var form FormRecipeInfo
	var recipe models.Recipe
	// var ingredient models.Ingredients
	// var recipeIngredient models.RecipeIngredient

	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	recipe.Title = form.Title
	recipe.Servings = form.Servings
	recipe.CookingTime = form.CookingTime


	// recipeIngredient.Quantities = make([]int, 0, len(form.IngredientQuantities))
	// recipe.Ingredients = make([]models.RecipeIngredient, 0, len(form.IngredientNames))

	createdRecipe, err := recipe.Create()
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, createdRecipe)
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