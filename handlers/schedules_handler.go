package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"thera-api/dto"
	"thera-api/services"
	"time"

	"github.com/gin-gonic/gin"
)

type SchedulesHandler struct {
	Service *services.SchedulesService
}

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

	fmt.Printf("Schedule created: %+v\n", schedule)

	c.JSON(http.StatusCreated, gin.H{"message": "jadwal berhasil dibuat"})
}

// GET /schedules
func (h *SchedulesHandler) GetAll(c *gin.Context) {
	tenantId := c.GetHeader("x-tenant-id")
	if tenantId == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error 01"})
		return
	}
	schedules, err := h.Service.GetAllSchedules(tenantId)

	var response []dto.ScheduleResponse
	for _, schedule := range schedules {
		response = append(response, dto.ScheduleResponse{
			ID:           schedule.ID,
			DateTime:     schedule.DateTime,
			CategoryId:   schedule.CategoryId,
			CategoryName: schedule.Categories.Name,
			Status:       schedule.Status,
		})
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *SchedulesHandler) GetByCatID(c *gin.Context) {
	tenantId := c.GetHeader("x-tenant-id")
	if tenantId == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error 01"})
		return
	}
	slug := c.Param("slug")
	date := c.Query("date")
	schedule, err := h.Service.GetScheduleByCatID(slug, tenantId, date)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "jadwal tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, schedule)
}

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

func (h *SchedulesHandler) GetScheduleByID(c *gin.Context) {
	tenantId := c.GetHeader("x-tenant-id")
	if tenantId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "x-tenant-id header is required"})
		return
	}

	id := c.Param("id")
	schedule, err := h.Service.GetScheduleByID(id, tenantId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "jadwal tidak ditemukan"})
		return
	}
	if schedule == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "jadwal tidak ditemukan"})
		return
	}
	var customFieldsData []dto.CustomField
	if len(schedule.Categories.CustomFields) > 0 {
		err := json.Unmarshal(schedule.Categories.CustomFields, &customFieldsData)
		if err != nil {
			fmt.Printf("Error unmarshal custom fields: %v\n", err)
		}
	}
	response := dto.ScheduleIdResponse{
		ID:             schedule.ID,
		DateTime:       schedule.DateTime,
		Name:           schedule.Categories.Name,
		NameEn:         schedule.Categories.NameEn,
		Status:         schedule.Status,
		Start:          schedule.Categories.Start,
		End:            schedule.Categories.End,
		IsGroup:        schedule.Categories.IsGroup,
		IsFree:         schedule.Categories.IsFree,
		IsPayAsYouWish: schedule.Categories.IsPayAsYouWish,
		IsManual:       schedule.Categories.IsManual,
	}
	if schedule.Categories.Description != nil {
		response.Description = *schedule.Categories.Description
	}
	if schedule.Categories.DescriptionEn != nil {
		response.DescriptionEn = *schedule.Categories.DescriptionEn
	}
	if schedule.Categories.Image != nil {
		response.Image = *schedule.Categories.Image
	}
	if schedule.Categories.Location != nil {
		response.Location = *schedule.Categories.Location
	}
	if schedule.Categories.Price != nil {
		response.Price = int(*schedule.Categories.Price)
	}
	response.CustomFields = customFieldsData

	c.JSON(http.StatusOK, response)
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
