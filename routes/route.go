package routes

import (
	initpkg "thera-api/init"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, c *initpkg.Container) {
	api := r.Group("/api")

	// panggil file route lain
	RegisterAuthRoutes(api, c)
	RegisterTenantRoutes(api, c)
	RegisterCategoriesRoutes(api, c)
	RegisterSchedulesRoutes(api, c)
	RegisterBookingRoutes(api, c)
	RegisterHeroRoutes(api, c)

	// nanti kalau ada modul lain:
	// RegisterEventRoutes(api, c)
	// RegisterGalleryRoutes(api, c)
}
