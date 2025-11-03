package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Coba load .env secara opsional
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, using environment variables from system")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("❌ DATABASE_URL not found in environment")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ failed to connect database:", err)
	}

	DB = db
	log.Println("✅ Database connected successfully")
}
