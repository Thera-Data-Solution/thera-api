package handlers

import (
	"fmt"
	"net/http"
	"thera-api/services"

	"github.com/gin-gonic/gin"
)

type GalleryHandler struct {
	Service *services.GalleryService
}

func (h *GalleryHandler) GetAll(c *gin.Context) {
	fmt.Println("testing disini 0")
	tenantId := c.GetHeader("x-tenant-id")
	if tenantId == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error 01"})
		return
	}
	fmt.Println("testing disini 0.1")
	gallery, err := h.Service.GetAllGallery(tenantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("testing disini 0.2")
	c.JSON(http.StatusOK, gallery)
}
