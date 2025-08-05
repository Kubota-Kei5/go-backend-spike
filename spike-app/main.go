package main

import (
	"os"
	"spike-app/controllers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	controllers.SetupRouter().Run(":" + port)
}