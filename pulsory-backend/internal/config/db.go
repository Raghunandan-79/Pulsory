package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DATABSE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect databse:", err)
	}

	DB = db

	log.Println("Database connected")
}
