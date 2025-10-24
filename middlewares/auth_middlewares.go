package middlewares

import (
	"net/http"
	"thera-api/repositories"

	"github.com/gin-gonic/gin"
)

type IsAuthMiddleware struct {
	SessionRepo *repositories.SessionRepository
	UserRepo    *repositories.UserRepository
	AdminRepo   *repositories.TenantUserRepository
}

func NewAuthMiddleware(
	sessionRepo *repositories.SessionRepository,
	userRepo *repositories.UserRepository,
	adminRepo *repositories.TenantUserRepository,
) *IsAuthMiddleware {
	return &IsAuthMiddleware{
		SessionRepo: sessionRepo,
		UserRepo:    userRepo,
		AdminRepo:   adminRepo,
	}
}

func (m *IsAuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token tidak ditemukan"})
			c.Abort()
			return
		}

		session, err := m.SessionRepo.FindByToken(token)
		if err != nil || session == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "session tidak valid atau tidak ditemukan"})
			c.Abort()
			return
		}

		var (
			userType string
			user     any
		)
		if session.TenantUserId != nil {
			u, err := m.AdminRepo.FindByID(*session.TenantUserId)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "pengguna tenant tidak ditemukan"})
				return
			}
			userType = "tenant"
			user = u
		} else if session.UserId != nil {
			u, err := m.UserRepo.FindByID(*session.UserId)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "pengguna tidak ditemukan"})
				return
			}
			userType = "user"
			user = u
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "data pengguna tidak valid"})
			c.Abort()
			return
		}
		c.Set("auth", gin.H{
			"user":         user,
			"session":      session,
			"userType":     userType,
			"tenantId":     session.TenantId,
			"userId":       session.UserId,
			"tenantUserId": session.TenantUserId,
		})
		c.Set("session", session)
		c.Next()
	}
}
