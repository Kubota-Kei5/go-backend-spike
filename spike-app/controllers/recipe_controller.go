package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Hello(c *gin.Context) {
	c.HTML(http.StatusOK, "hello.html", gin.H{})
}

func NewRecipe(c *gin.Context) {
	c.HTML(http.StatusOK, "new_recipe.html", gin.H{})
}

func ListRecipe(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}