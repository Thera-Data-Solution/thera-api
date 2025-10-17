package routes

import "github.com/gin-gonic/gin"

func SetupRoutes() *gin.Engine {
	r := gin.Default()

	// daftar semua group route di sini
	UserRoutes(r)

	return r
}
