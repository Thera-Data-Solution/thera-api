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
	config.ConnectDatabase()

	if config.DB == nil {
		log.Fatal("❌ DB masih nil setelah ConnectDatabase()")
	} else {
		log.Println("✅ DB aktif di main.go")
	}

	tenantUserRepo := &repository.TenantUserRepository{DB: config.DB}
	sessionRepo := &repository.SessionRepository{DB: config.DB}
	userRepo := &repository.UserRepository{DB: config.DB}
	tenantRepo := &repository.TenantRepo{DB: config.DB}

	authController := &controllers.AuthController{
		TenantUserRepo: tenantUserRepo,
		TenantRepo:     tenantRepo,
		SessionRepo:    sessionRepo,
		UserRepo:       userRepo,
	}

	r := routes.SetupRoutes(authController)

	log.Println("🚀 Server berjalan di http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
