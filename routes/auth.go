package routes

import (
	initpkg "thera-api/init"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.RouterGroup, c *initpkg.Container) {
	user := router.Group("/auth/user")
	{
		user.POST("/register", c.UserHandler.Register)
		user.POST("/login", c.UserHandler.Login)
		user.GET("/me", c.Middlewares.Handle(), c.UserHandler.Me)
		user.POST("/forgot-password", c.UserHandler.ForgotPassword)
		user.POST("/reset-password", c.UserHandler.ResetPassword)
		user.PUT("/profile", c.Middlewares.Handle(), c.UserHandler.UpdateProfile)
	}

	admin := router.Group("/auth/admin")
	{
		admin.POST("/register", c.AdminHandler.Register)
		admin.POST("/login", c.AdminHandler.Login)
		admin.GET("/me", c.Middlewares.Handle(), c.UserHandler.Me)
		admin.POST("/forgot-password", c.AdminHandler.ForgotPassword)
		admin.POST("/reset-password", c.AdminHandler.ResetPassword)
		admin.PUT("/profile", c.Middlewares.Handle(), c.AdminHandler.UpdateProfile)
	}
}
