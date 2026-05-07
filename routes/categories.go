package routes

import (
	initpkg "thera-api/init"

	"github.com/gin-gonic/gin"
)

func RegisterCategoriesRoutes(router *gin.RouterGroup, c *initpkg.Container) {
	categories := router.Group("/categories")
	{
		categories.GET("", c.CategoryHandler.GetAll)
		categories.GET("/:id", c.CategoryHandler.GetByID)
	}

	adminCategories := router.Group("/cadm")
	{
		adminCategories.GET("", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.CategoryHandler.GetAllAsAdmin)
		adminCategories.GET("/:id", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.CategoryHandler.GetByIDAsAdmin)
		adminCategories.POST("", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.CategoryHandler.Create)
		adminCategories.PUT("/:id", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.CategoryHandler.Update)
		adminCategories.DELETE("/:id", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.CategoryHandler.Delete)
	}
}
