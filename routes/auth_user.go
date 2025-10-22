package routes

import (
	"thera-api/controllers"
	"thera-api/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine, ac *controllers.AuthController) {
	auth := r.Group("/api/auth")

	{
		auth.POST("/register", ac.UserRegister)
		auth.POST("/login", ac.UserLogin)
	}

	protected := r.Group("")
	protected.Use(middleware.UserAuthMiddleware(ac.SessionRepo, ac.UserRepo))
	{
		protected.GET("/me", ac.AdminMe)
	}
}
