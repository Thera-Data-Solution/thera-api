package routes

import (
	initpkg "thera-api/init"

	"github.com/gin-gonic/gin"
)

func RegisterReviewRoutes(router *gin.RouterGroup, c *initpkg.Container) {
	// Publik (Landing Page)
	router.GET("/testimonials", c.ReviewHandler.GetLandingPageReviews)

	// User (Authenticated)
	user := router.Group("/history")
	user.Use(c.Middlewares.Handle())
	{
		user.GET("", c.ReviewHandler.GetMyHistory)
		user.POST("/review", c.ReviewHandler.SubmitReview)
	}

	// Admin (Khusus moderasi)
	// admin := router.Group("/radm") ... dst
}
