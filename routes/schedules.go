package routes

import (
	initpkg "thera-api/init"

	"github.com/gin-gonic/gin"
)

func RegisterSchedulesRoutes(router *gin.RouterGroup, c *initpkg.Container) {
	schedule := router.Group("/schedules")
	{
		schedule.GET("", c.ScheduleHandler.GetAll)
		schedule.GET("/:id", c.ScheduleHandler.GetByID)
		schedule.POST("", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.ScheduleHandler.Create)
		schedule.PUT("/:id", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.ScheduleHandler.Update)
		schedule.DELETE("/:id", c.Middlewares.Handle(), c.AtLeastAdmin.Handle(), c.ScheduleHandler.Delete)
	}
}
