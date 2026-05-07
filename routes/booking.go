package routes

import (
	initpkg "thera-api/init"

	"github.com/gin-gonic/gin"
)

func RegisterBookingRoutes(router *gin.RouterGroup, c *initpkg.Container) {
	book := router.Group("/booking")
	{
		book.POST("", c.Middlewares.Handle(), c.BookHandler.Create)
		book.POST("/cancel", c.Middlewares.Handle(), c.BookHandler.Cancel)
	}

	adminBook := router.Group("/abook")
	{
		adminBook.GET("/", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.BookHandler.GetAllAdmin)
		adminBook.GET("/change-status", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.BookHandler.ChangeStatus)
		adminBook.GET("/close-schedule", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.BookHandler.CloseSchedule)

	}
}
