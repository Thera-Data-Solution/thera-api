package routes

import (
	initpkg "thera-api/init"

	"github.com/gin-gonic/gin"
)

func RegisterReviewRoutes(router *gin.RouterGroup, c *initpkg.Container) {
	router.GET("/testimonials", c.ReviewHandler.GetLandingPageReviews)

	user := router.Group("/history")
	user.Use(c.Middlewares.Handle())
	{
		user.GET("", c.ReviewHandler.GetMyHistory)
		user.POST("/review", c.ReviewHandler.SubmitReview)
	}

	admin := router.Group("/radm")
	admin.Use(c.Middlewares.Handle(), c.AtLeastAdmin.Handle())
	{
		admin.GET("", c.ReviewHandler.AdminGetAll)
		admin.PUT("/:id", c.ReviewHandler.AdminUpdate)
	}
}
