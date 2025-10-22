package middleware

import (
	"net/http"
	"time"

	"thera-api/repository"

	"github.com/gin-gonic/gin"
)

func UserAuthMiddleware(sessionRepo *repository.SessionRepository, userRepo *repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token missing"})
			c.Abort()
			return
		}

		ses, err := sessionRepo.FindSessionByToken(token)
		if err != nil || ses == nil || ses.UserId == nil {
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

		user, err := userRepo.FindUserById(*ses.UserId)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
			c.Abort()
			return
		}

		c.Set("user", user)

		c.Next()
	}
}
