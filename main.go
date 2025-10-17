package main

import (
	"log"
	"net/http"

	"go-api/config"
	"go-api/routes"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Tidak ada file .env")
	}

	config.ConnectDatabase()
	config.InitMigrate()

	r := routes.SetupRoutes()

	log.Println("🚀 Server berjalan di http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
