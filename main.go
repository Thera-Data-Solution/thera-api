package main

import (
	"os"
	"thera-api/config"
	initpkg "thera-api/init"
	"thera-api/logger"
	"thera-api/migrate"
	"thera-api/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	logger.Init()
	config.ConnectDatabase()
	migrate.RunMigrations()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://staging.theravickya.com", "https://staging.admin.theravickya.com", "https://theravickya.com", "https://admin.theravickya.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "X-Tenant-Id", "Token"},
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
	logger.Log.Info("Server running", zap.String("port", port))
	r.Run("0.0.0.0:" + port)
}
