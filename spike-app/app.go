package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		response := gin.H{
			"message": "Hello, World!",
		}
		c.JSON(http.StatusOK, response)
	})

	r.Run(":8080")
}