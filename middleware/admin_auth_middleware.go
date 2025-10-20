package middleware

import (
	"net/http"
	"time"

	"thera-api/repository"

	"github.com/gin-gonic/gin"
)

func AdminAuthMiddleware(sessionRepo *repository.SessionRepository, tenantRepo *repository.TenantUserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token missing"})
			c.Abort()
			return
		}

		ses, err := sessionRepo.FindByToken(token)
		if err != nil || ses == nil || ses.TenantUserId == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
			c.Abort()
			return
		}

		exp, err := time.Parse(time.RFC3339, ses.ExpiresAt)
		if err != nil || time.Now().UTC().After(exp) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "session expired"})
			c.Abort()
			return
		}

		// var user *models.TenantUser
		user, err := tenantRepo.FindTenantById(*ses.TenantUserId)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			c.Abort()
			return
		}

		c.Set("user", user)

		c.Next()
	}
}
