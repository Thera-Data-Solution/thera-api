package routes

import (
	initpkg "thera-api/init"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.RouterGroup, c *initpkg.Container) {
	user := router.Group("/users")
	{
		user.GET("", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.TenantUserHandler.GetAllByTenantId)
		user.POST("disable/:id", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.TenantUserHandler.DisableUser)
		user.GET("/admin", c.Middlewares.Handle(), c.OnlySU.Handle(), c.TenantUserHandler.GetAll)
	}
}
