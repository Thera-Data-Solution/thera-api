package routes

import (
	initpkg "thera-api/init" // Sesuaikan path jika perlu

	"github.com/gin-gonic/gin"
)

func RegisterArticleRoutes(router *gin.RouterGroup, c *initpkg.Container) {
	articles := router.Group("/articles")
	{
		articles.GET("", c.ArticleHandler.GetAll)
		articles.GET("/:id", c.ArticleHandler.GetByID)

		articles.POST("", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.ArticleHandler.Create)
		articles.PUT("/:id", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.ArticleHandler.Update)
		articles.DELETE("/:id", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.ArticleHandler.Delete)
	}
}
