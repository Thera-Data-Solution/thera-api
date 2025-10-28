package routes

import (
	initpkg "thera-api/init"

	"github.com/gin-gonic/gin"
)

func RegisterTenantRoutes(router *gin.RouterGroup, c *initpkg.Container) {
	tenant := router.Group("/tenants")

	{
		tenant.GET("/", c.Middlewares.Handle(), c.OnlySU.Handle(), c.TenantHandler.GetAll)
		tenant.POST("/", c.Middlewares.Handle(), c.OnlySU.Handle(), c.TenantHandler.Create)
		tenant.GET("/:id", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.TenantHandler.GetByID)
		tenant.PUT("/:id", c.Middlewares.Handle(), c.OnlySU.Handle(), c.TenantHandler.Update)
		tenant.DELETE("/:id", c.Middlewares.Handle(), c.OnlySU.Handle(), c.TenantHandler.Delete)
	}
}
