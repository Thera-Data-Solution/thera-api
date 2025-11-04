package main

import (
	"fmt"
	"os"
	"thera-api/config"
	initpkg "thera-api/init"
	"thera-api/logger"
	"thera-api/migrate"
	"thera-api/routes"
	"time" // Import time untuk MaxAge

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	logger.Init()
	config.ConnectDatabase()
	migrate.RunMigrations()

	r := gin.New()
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false
	r.HandleMethodNotAllowed = true
	r.Use(gin.Recovery(), gin.Logger())

	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"*"}, // Mengaktifkan lagi header * untuk lebih aman
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}))
	fmt.Println(r.Routes())

	// <<< Tambahkan rute test ini >>>
	r.GET("/cors-test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "CORS test successful!"})
	})
	// <<< Akhir penambahan >>>

	container := initpkg.NewContainer()
	routes.SetupRoutes(r, container)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	logger.Log.Info("Server running", zap.String("port", port))
	r.Run("0.0.0.0:" + port)
}
