package routes

import (
	initpkg "thera-api/init"

	"github.com/gin-gonic/gin"
)

func RegisterGalleryRoutes(router *gin.RouterGroup, c *initpkg.Container) {
	gallery := router.Group("/gallery")
	{
		gallery.GET("/", c.GalleryHandler.GetAll)
		// gallery.GET("/:id", c.GalleryHandler.GetByID)
		// gallery.POST("/", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.GalleryHandler.Create)
		// gallery.PUT("/:id", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.GalleryHandler.Update)
		// gallery.DELETE("/:id", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.GalleryHandler.Delete)
	}
}
