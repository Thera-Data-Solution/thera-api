package routes

import (
	"thera-api/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(authController *controllers.AuthController) *gin.Engine {
	r := gin.Default()
	AuthRoutes(r, authController)

	return r
}
