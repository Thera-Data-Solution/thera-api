package routes

import (
	initpkg "thera-api/init"

	"github.com/gin-gonic/gin"
)

func RegisterEventsRoutes(router *gin.RouterGroup, c *initpkg.Container) {
	// Public Routes
	events := router.Group("/events")
	{
		events.GET("", c.EventHandler.GetAll)
		events.GET("/:id", c.EventHandler.GetByID)
	}

	adminEvents := router.Group("/eadm")
	{
		adminEvents.GET("", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.EventHandler.GetAllAsAdmin)
		adminEvents.POST("", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.EventHandler.Create)
		adminEvents.PUT("/:id", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.EventHandler.Update)
		adminEvents.DELETE("/:id", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.EventHandler.Delete)
	}
}
