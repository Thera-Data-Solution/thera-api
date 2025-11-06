package routes

import (
	initpkg "thera-api/init"

	"github.com/gin-gonic/gin"
)

func RegisterTranslationRoutes(router *gin.RouterGroup, c *initpkg.Container) {
	tr := router.Group("/translations")
	{
		tr.GET("", c.TranslationHandler.GetAll)
		tr.GET("/:id", c.TranslationHandler.GetByID)
		tr.POST("", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.TranslationHandler.Create)
		tr.PUT("/:id", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.TranslationHandler.Update)
		tr.DELETE("/:id", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.TranslationHandler.Delete)
	}
}
