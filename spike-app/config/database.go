package config

import (
	"os"

	"spike-app/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() (*gorm.DB, error) {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	env := os.Getenv("ENV")
	
	var host string
	switch env {
	case "production", "prod":
		host = "db-prod"
	case "development", "dev":
		host = "db-dev"
	default:
		host = "db-dev"
	}

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := migrateDatabase(db); err != nil {
		return nil, err
	}
	return db, nil
}

func migrateDatabase(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Recipe{},
		&models.Ingredients{},
		&models.RecipeIngredient{},
	)
}