package routes

import (
	initpkg "thera-api/init"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, c *initpkg.Container) {
	api := r.Group("/api")
	RegisterAuthRoutes(api, c)
	RegisterTenantRoutes(api, c)
	RegisterCategoriesRoutes(api, c)
	RegisterSchedulesRoutes(api, c)
	RegisterBookingRoutes(api, c)
	RegisterHeroRoutes(api, c)
	RegisterLinkRoutes(api, c)
}
