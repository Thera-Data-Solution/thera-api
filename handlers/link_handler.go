package handlers

import (
	"net/http"
	"thera-api/dto"
	"thera-api/logger"
	"thera-api/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LinkHandler struct {
	Service *services.LinkService
}

// POST /links
func (h *LinkHandler) Create(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}
	auth := authData.(gin.H)

	var body dto.CreateLinkDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenantID := auth["tenantId"].(string)
	link, err := h.Service.CreateLink(body.Name, body.Value, body.Type, body.Icon, &tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, link)
}

// GET /links
func (h *LinkHandler) GetAll(c *gin.Context) {
	tenantId := c.GetHeader("x-tenant-id")
	if tenantId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenantId tidak ditemukan"})
		return
	}

	links, err := h.Service.GetAllLinks(tenantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, links)
}

// GET /links/:id
func (h *LinkHandler) GetByID(c *gin.Context) {
	tenantId := c.GetHeader("x-tenant-id")
	id := c.Param("id")

	link, err := h.Service.GetLinkByID(id, tenantId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "link tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, link)
}

// PUT /links/:id
func (h *LinkHandler) Update(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}
	auth := authData.(gin.H)
	id := c.Param("id")

	var body dto.UpdateLinkDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dtoMap := make(map[string]interface{})
	if body.Name != nil {
		dtoMap["name"] = *body.Name
	}
	if body.Value != nil {
		dtoMap["value"] = *body.Value
	}
	if body.Type != nil {
		dtoMap["type"] = *body.Type
	}
	if body.Icon != nil {
		dtoMap["icon"] = *body.Icon
	}
	if body.Order != nil {
		dtoMap["order"] = *body.Order
	}

	link, err := h.Service.UpdateLink(id, dtoMap, auth["tenantId"].(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, link)
}

// DELETE /links/:id
func (h *LinkHandler) Delete(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}
	auth := authData.(gin.H)
	id := c.Param("id")

	if err := h.Service.DeleteLink(id, auth["tenantId"].(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "link berhasil dihapus"})
}

func (h *LinkHandler) MoveUp(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}
	auth := authData.(gin.H)
	id := c.Param("id")
	tenantId := auth["tenantId"].(string)

	link, err := h.Service.MoveUp(id, tenantId)
	if err != nil {
		logger.Log.Error("Failed to move link up via service", zap.Error(err), zap.String("id", id), zap.String("tenantId", tenantId))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, link)
}

func (h *LinkHandler) MoveDown(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}
	auth := authData.(gin.H)
	id := c.Param("id")
	tenantId := auth["tenantId"].(string)

	link, err := h.Service.MoveDown(id, tenantId)
	if err != nil {
		logger.Log.Error("Failed to move link down via service", zap.Error(err), zap.String("id", id), zap.String("tenantId", tenantId))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, link)
}
