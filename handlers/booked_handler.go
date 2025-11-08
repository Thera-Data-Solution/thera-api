package handlers

import (
	"net/http"
	"strconv"
	"thera-api/services"

	"github.com/gin-gonic/gin"
)

type BookedHandler struct {
	Service *services.BookedService
}

func NewBookedHandler(service *services.BookedService) *BookedHandler {
	return &BookedHandler{Service: service}
}

func (h *BookedHandler) Create(c *gin.Context) {
	var req struct {
		ScheduleId string `json:"scheduleId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authData, _ := c.Get("auth")
	auth := authData.(gin.H)
	tenantId := auth["tenantId"].(string)
	userType := auth["userType"].(string)
	var userIdentifier string

	if userType != "user" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan"})
		return
	}
	if uid, ok := auth["userId"].(*string); ok && uid != nil {
		userIdentifier = *uid
	}

	err := h.Service.Create(userIdentifier, req.ScheduleId, tenantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "booking berhasil dibuat"})
}

func (h *BookedHandler) GetByUserId(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists || authData == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	auth, ok := authData.(gin.H)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid auth data type"})
		return
	}

	tenantId := auth["tenantId"].(string)
	userType := auth["userType"].(string)

	var userIdentifier string

	if userType != "user" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan"})
		return
	}
	if uid, ok := auth["userId"].(*string); ok && uid != nil {
		userIdentifier = *uid
	}

	booked, err := h.Service.GetByUser(tenantId, userIdentifier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, booked)
}

func (h *BookedHandler) GetById(c *gin.Context) {
	id := c.Param("id")
	authData, _ := c.Get("auth")
	auth := authData.(gin.H)
	tenantId := auth["tenantId"].(string)

	booked, err := h.Service.GetById(id, tenantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, booked)
}
func (h *BookedHandler) GetAll(c *gin.Context) {
	authData, _ := c.Get("auth")
	auth := authData.(gin.H)
	tenantId := auth["tenantId"].(string)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	booked, total, err := h.Service.GetAll(tenantId, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  booked,
		"page":  page,
		"limit": limit,
		"total": total,
	})
}

func (h *BookedHandler) Cancel(c *gin.Context) {
	id := c.Param("id")
	authData, _ := c.Get("auth")
	auth := authData.(gin.H)
	tenantId := auth["tenantId"].(string)

	err := h.Service.Cancel(id, tenantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "kategori berhasil dihapus"})
}
