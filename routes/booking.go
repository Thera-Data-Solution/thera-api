package routes

import (
	initpkg "thera-api/init"

	"github.com/gin-gonic/gin"
)

func RegisterBookingRoutes(router *gin.RouterGroup, c *initpkg.Container) {
	book := router.Group("/booking")
	{
		book.GET("", c.Middlewares.Handle(), c.BookHandler.GetByUserId)
		book.POST("", c.Middlewares.Handle(), c.BookHandler.Create)
		book.GET("/all", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.BookHandler.GetAll)
		book.GET("/one/:id", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.BookHandler.GetById)
		book.DELETE("/one/:id", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.BookHandler.Cancel)
	}
}
