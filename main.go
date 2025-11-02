package main

import (
	"os"
	"thera-api/config"
	initpkg "thera-api/init"
	"thera-api/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://staging.theraavickya.com"}, // ganti sesuai origin frontend kamu
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Tenant-Id", "Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	container := initpkg.NewContainer()
	routes.SetupRoutes(r, container)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run("0.0.0.0:" + port)
}
