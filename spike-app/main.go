package main

import (
	"spike-app/router"
)

func main() {
	router.SetupRouter().Run(":8080")
}