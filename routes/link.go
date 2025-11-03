package routes

import (
	initpkg "thera-api/init"

	"github.com/gin-gonic/gin"
)

func RegisterLinkRoutes(router *gin.RouterGroup, c *initpkg.Container) {
	links := router.Group("/links")
	{
		links.GET("/", c.LinkHandler.GetAll)
		links.GET("/:id", c.LinkHandler.GetByID)
		links.POST("/", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.LinkHandler.Create)
		links.PUT("/:id", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.LinkHandler.Update)
		links.POST("/:id/order/up", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.LinkHandler.MoveUp)
		links.POST("/:id/order/down", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.LinkHandler.MoveDown)
		links.DELETE("/:id", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.LinkHandler.Delete)
	}
}
