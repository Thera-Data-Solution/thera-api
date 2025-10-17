package middleware

import (
	"fmt"
	"go-api/config"
	"go-api/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware memeriksa token dari header Authorization
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("token")
		fmt.Println(authHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan"})
			c.Abort()
			return
		}

		var session models.Session
		if err := config.DB.Where("token = ?", authHeader).First(&session).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			c.Abort()
			return
		}
		fmt.Println(session)

		if session.ExpiresAt.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token sudah kadaluarsa"})
			c.Abort()
			return
		}

		var user models.User
		if err := config.DB.Where("id = ?", session.UserID).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan"})
			c.Abort()
			return
		}

		// Simpan user ke context
		c.Set("user", user)
		c.Next()
	}
}
