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

	var dsn string

	// Cloud Run環境では DATABASE_URL を優先使用
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" && (env == "production" || env == "prod") {
		dsn = databaseURL
	} else {
		var host string
		switch env {
		case "production", "prod":
			host = "cloud-sql"
		case "development", "dev":
			host = "db-dev"
		case "test", "ci":
			host = "localhost"
		default:
			host = "cloud-sql"
		}
		dsn = "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	}

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
