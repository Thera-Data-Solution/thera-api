package routes

import (
	initpkg "thera-api/init"

	"github.com/gin-gonic/gin"
)

func RegisterPartnerRoutes(router *gin.RouterGroup, c *initpkg.Container) {
	partner := router.Group("/partner")
	{
		partner.GET("", c.PartnerHandler.GetAll)
		partner.POST("", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.PartnerHandler.Create)
		partner.PUT("/:id", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.PartnerHandler.Update)
		partner.DELETE("/:id", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.PartnerHandler.Delete)
		// partner.GET("/:id", c.PartnerHandler.GetByID)
	}
}
