package main

import (
	"log"
	"os"
	"spike-app/config"
	"spike-app/controllers"
)

func main() {
	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get database instance:", err)
	}
	defer sqlDB.Close()

	log.Println("Database connected and migrated successfully")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	controllers.SetupRouter().Run(":" + port)
}