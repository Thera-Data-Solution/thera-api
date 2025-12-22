package routes

import (
	initpkg "thera-api/init"

	"github.com/gin-gonic/gin"
)

func RegisterCategoriesRoutes(router *gin.RouterGroup, c *initpkg.Container) {
	categories := router.Group("/categories")
	{
		categories.GET("", c.CategoryHandler.GetAll)
		categories.GET("/", c.CategoryHandler.GetAll)
		categories.GET("/:id", c.CategoryHandler.GetByID)
		categories.POST("", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.CategoryHandler.Create)
		categories.PUT("/:id", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.CategoryHandler.Update)
		categories.DELETE("/:id", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.CategoryHandler.Delete)
	}
}
