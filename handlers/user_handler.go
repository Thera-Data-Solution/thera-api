package handlers

import (
	"math"
	"net/http"
	"strconv"
	"thera-api/dto"
	"thera-api/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service *services.UserService
}

func (h *UserHandler) GetAllByTenantId(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}
	auth := authData.(gin.H)
	tenantId := auth["tenantId"].(string)

	// Ambil query param pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// Ambil data + total count dari service
	users, total, err := h.Service.GetAllByTenantPaginated(tenantId, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Mapping ke DTO
	var data []dto.UserNih
	for _, u := range users {
		data = append(data, dto.UserNih{
			ID:       u.ID,
			FullName: u.FullName,
			Email:    u.Email,
			Phone:    u.Phone,
			Address:  u.Address,
			Ig:       u.Ig,
			Fb:       u.Fb,
			Avatar:   u.Avatar,
			Disable:  u.Disable,
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	// Response final
	resp := dto.UserResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	users, tenantNames, total, err := h.Service.GetAllPaginated(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var data []dto.UserNih
	for i, u := range users {
		data = append(data, dto.UserNih{
			ID:         u.ID,
			FullName:   u.FullName,
			Email:      u.Email,
			Phone:      u.Phone,
			Address:    u.Address,
			Ig:         u.Ig,
			Fb:         u.Fb,
			Avatar:     u.Avatar,
			Disable:    u.Disable,
			TenantName: tenantNames[i],
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	// Response final
	resp := dto.UserResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) DisableUser(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}
	auth := authData.(gin.H)
	tenantId := auth["tenantId"].(string)
	id := c.Param("id")

	err := h.Service.DisableUser(id, tenantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User berhasil dinonaktifkan"})
}
