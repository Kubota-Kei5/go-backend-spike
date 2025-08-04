package main

import (
	"spike-app/controllers"
)

func main() {
	controllers.SetupRouter().Run(":8080")
}