package handlers

import (
	"encoding/json"
	"net/http"
	"thera-api/models"
	"thera-api/services"
	"thera-api/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

type BookedHandler struct {
	Service *services.BookedService
}

func NewBookedHandler(service *services.BookedService) *BookedHandler {
	return &BookedHandler{Service: service}
}

func (h *BookedHandler) Create(c *gin.Context) {
	var req struct {
		ScheduleId   string `json:"scheduleId" binding:"required"`
		Type         int    `json:"type"` // 1 untuk schedule biasa, 2 untuk event
		CustomAnswer string `json:"customAnswer"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authData, _ := c.Get("auth")
	auth := authData.(gin.H)

	tenantId := auth["tenantId"].(string)
	userType := auth["userType"].(string)

	if userType != "user" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak valid"})
		return
	}

	var userIdentifier string
	if uid, ok := auth["userId"].(*string); ok && uid != nil {
		userIdentifier = *uid
	}

	var customAnswer datatypes.JSON

	if req.CustomAnswer != "" {
		var temp []models.BookedCustomField

		if err := json.Unmarshal([]byte(req.CustomAnswer), &temp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "format customAnswer tidak valid",
			})
			return
		}

		customAnswer = datatypes.JSON(req.CustomAnswer)
	}

	if err := h.Service.Create(
		userIdentifier,
		req.ScheduleId,
		tenantId,
		customAnswer,
		req.Type,
	); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "booking berhasil dibuat",
	})
}

func (h *BookedHandler) Cancel(c *gin.Context) {
	var req struct {
		ScheduleId string `json:"scheduleId" binding:"required"`
		Type       int    `json:"type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	authData, _ := c.Get("auth")
	auth := authData.(gin.H)
	tenantId := auth["tenantId"].(string)

	err := h.Service.Cancel(req.ScheduleId, req.Type, auth["userId"].(*string), tenantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "booking berhasil dibatalkan"})
}

func (h *BookedHandler) ChangeStatus(c *gin.Context) {
	id := c.Query("id")
	status := c.Query("status")

	auth := c.MustGet("auth").(gin.H)
	tenantId := auth["tenantId"].(string)

	if id == "" || status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id dan status diperlukan"})
		return
	}

	if err := h.Service.AdminChangeStatus(id, status, tenantId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status booking berhasil diperbarui"})
}

func (h *BookedHandler) CloseSchedule(c *gin.Context) {
	scheduleId := c.Query("scheduleId")
	auth := c.MustGet("auth").(gin.H)
	tenantId := auth["tenantId"].(string)

	if scheduleId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "scheduleId diperlukan"})
		return
	}

	if err := h.Service.CloseSchedule(scheduleId, tenantId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "jadwal telah diselesaikan (CLOSED)"})
}

func (h *BookedHandler) GetAllAdmin(c *gin.Context) {
	auth := c.MustGet("auth").(gin.H)
	tenantId := auth["tenantId"].(string)

	limit := utils.ParseInt(c.DefaultQuery("limit", "10"))
	offset := utils.ParseInt(c.DefaultQuery("offset", "0"))

	data, total, err := h.Service.GetAllForAdmin(tenantId, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":  total,
		"limit":  limit,
		"offset": offset,
		"data":   data,
	})
}
