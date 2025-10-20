// main.go (Versi Gin)
package main

import (
	"log"
	"net/http"
	"thera-api/config"
	"thera-api/controllers"
	"thera-api/repository"
	"thera-api/routes"

	// Ganti import Echo
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file, using system env")
	}

	config.ConnectDatabase()

	tenantUserRepo := &repository.TenantUserRepository{DB: config.DB}
	sessionRepo := &repository.SessionRepository{DB: config.DB}

	authController := &controllers.AuthController{
		TenantUserRepo: tenantUserRepo,
		SessionRepo:    sessionRepo,
	}

	r := routes.SetupRoutes(authController)

	log.Println("🚀 Server berjalan di http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
