package middlewares

import (
	"net/http"
	"thera-api/models"

	"github.com/gin-gonic/gin"
)

type IsSUMiddleware struct{}

func NewIsSUMiddleware() *IsSUMiddleware {
	return &IsSUMiddleware{}
}

func (m *IsSUMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ambil data auth dari middleware sebelumnya
		authData, exists := c.Get("auth")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "autentikasi belum dilakukan"})
			c.Abort()
			return
		}

		authMap, ok := authData.(gin.H)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "data autentikasi tidak valid"})
			c.Abort()
			return
		}

		userType, _ := authMap["userType"].(string)
		if userType != "tenant" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak"})
			c.Abort()
			return
		}

		userData, ok := authMap["user"].(*models.TenantUser)
		if !ok || userData == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "data tenant user tidak ditemukan"})
			c.Abort()
			return
		}

		// hanya izinkan role tertentu
		switch userData.Role {
		case "SU":
			c.Next()
			return
		default:
			c.JSON(http.StatusForbidden, gin.H{"error": "Akses ditolak"})
			c.Abort()
			return
		}
	}
}
