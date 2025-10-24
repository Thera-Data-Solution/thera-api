package routes

import (
	initpkg "thera-api/init"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, c *initpkg.Container) {
	api := r.Group("/api")

	// panggil file route lain
	RegisterAuthRoutes(api, c)

	// nanti kalau ada modul lain:
	// RegisterEventRoutes(api, c)
	// RegisterGalleryRoutes(api, c)
}
