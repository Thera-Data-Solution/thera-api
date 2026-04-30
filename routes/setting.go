package routes

import (
	initpkg "thera-api/init"

	"github.com/gin-gonic/gin"
)

func RegisterSettingRoutes(r *gin.RouterGroup, c *initpkg.Container) {
	// r.Use(middlewares.CORSMiddleware())

	group := r.Group("/settings")

	group.GET("", c.SettingHandler.GetAll)
	group.GET("admin", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.SettingHandler.GetAllAdmin)
	// group.GET("/:id", c.SettingHandler.GetById)
	// group.GET("/tenant/:tenantId", c.SettingHandler.GetByTenantId)
	group.POST("", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.SettingHandler.Upsert)
	// group.DELETE("/:id", c.SettingHandler.Delete)
}
