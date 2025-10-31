package handlers

import (
	"fmt"
	"net/http"
	"thera-api/services"
	"time"

	"github.com/gin-gonic/gin"
)

type SchedulesHandler struct {
	Service *services.SchedulesService
}

// POST /schedules
type ScheduleRequest struct {
	DateTime   string `json:"dateTime"`
	CategoryId string `json:"categoryId"`
	Status     string `json:"status"`
}

func (h *SchedulesHandler) Create(c *gin.Context) {
	var req ScheduleRequest
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}
	auth := authData.(gin.H)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format request tidak valid"})
		return
	}

	dateTime, err := time.Parse(time.RFC3339, req.DateTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format tanggal tidak valid, gunakan RFC3339"})
		return
	}

	fmt.Println(auth["tenantId"].(string))

	schedule, err := h.Service.CreateSchedule(dateTime, req.CategoryId, req.Status, auth["tenantId"].(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, schedule)
}

// GET /schedules
func (h *SchedulesHandler) GetAll(c *gin.Context) {
	tenantId := c.GetHeader("x-tenant-id")
	if tenantId == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error 01"})
		return
	}
	schedules, err := h.Service.GetAllSchedules(tenantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, schedules)
}

// GET /schedules/:id
func (h *SchedulesHandler) GetByID(c *gin.Context) {
	tenantId := c.GetHeader("x-tenant-id")
	if tenantId == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error 01"})
		return
	}
	id := c.Param("id")
	schedule, err := h.Service.GetScheduleByID(id, tenantId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "jadwal tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, schedule)
}

type UpdateScheduleRequest struct {
	DateTime   string  `json:"dateTime"`
	CategoryId *string `json:"categoryId"`
	Status     *string `json:"status"`
}

func (h *SchedulesHandler) Update(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}
	auth := authData.(gin.H)
	id := c.Param("id")

	var req UpdateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format request tidak valid"})
		return
	}

	var dateTime *time.Time
	if req.DateTime != "" {
		parsed, err := time.Parse(time.RFC3339, req.DateTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "format tanggal tidak valid, gunakan RFC3339"})
			return
		}
		dateTime = &parsed
	}

	// Jika categoryId dan status tidak disediakan, biarkan nil
	var categoryIdPtr *string
	if req.CategoryId != nil {
		categoryIdPtr = req.CategoryId
	}

	var statusPtr *string
	if req.Status != nil {
		statusPtr = req.Status
	}

	schedule, err := h.Service.UpdateSchedule(id, dateTime, categoryIdPtr, statusPtr, auth["tenantId"].(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

// DELETE /schedules/:id
func (h *SchedulesHandler) Delete(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}
	auth := authData.(gin.H)
	id := c.Param("id")
	if err := h.Service.DeleteSchedule(id, auth["tenantId"].(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "jadwal berhasil dihapus"})
}
