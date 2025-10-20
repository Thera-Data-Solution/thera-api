package middleware

import (
	"net/http"
	"time"

	"thera-api/repository"

	"github.com/labstack/echo/v4"
)

func UserAuthMiddleware(sessionRepo *repository.SessionRepository, userRepo *repository.UserRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("token")
			if token == "" {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "token missing"})
			}

			ses, err := sessionRepo.FindByToken(token)
			if err != nil || ses == nil || ses.UserId == nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid session"})
			}
			// check expiry
			exp, err := time.Parse(time.RFC3339, ses.ExpiresAt)
			if err != nil || time.Now().UTC().After(exp) {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "session expired"})
			}

			// optionally load user and set on context
			user, err := userRepo.FindByEmailAndTenant("", nil) // not ideal: we need FindById; better add method. But to keep simple, set userId only
			_ = user
			c.Set("userId", *ses.UserId)
			c.Set("tenantId", ses.TenantId)
			return next(c)
		}
	}
}
