package routes

import (
	initpkg "thera-api/init"

	"github.com/gin-gonic/gin"
)

func RegisterHeroRoutes(router *gin.RouterGroup, c *initpkg.Container) {
	hero := router.Group("/hero")
	{
		hero.GET("/", c.HeroHandler.GetAll)
		hero.GET("/:id", c.HeroHandler.GetByID)
		hero.POST("/", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.HeroHandler.Create)
		hero.PUT("/:id", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.HeroHandler.Update)
		hero.DELETE("/:id", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.HeroHandler.Delete)
	}
}
