package routes

import (
	"thera-api/controllers"
	"thera-api/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine, ac *controllers.AuthController) {
	auth := r.Group("/api/admin")

	{
		auth.POST("/register", ac.AdminRegister)
		auth.POST("/login", ac.AdminLogin)
	}

	protected := r.Group("")
	protected.Use(middleware.AdminAuthMiddleware(ac.SessionRepo, ac.TenantUserRepo))
	{
		protected.GET("/me", ac.AdminMe)
	}
}
