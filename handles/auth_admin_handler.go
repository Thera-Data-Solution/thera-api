package handlers

import (
	"net/http"
	"thera-api/models"
	"thera-api/services"

	"github.com/gin-gonic/gin"
)

type AuthAdminHandler struct {
	Service *services.AuthAdminService
}

func (h *AuthAdminHandler) Register(c *gin.Context) {
	tenantId := c.GetHeader("x-tenant-id")
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		FullName string `json:"fullName" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if tenantId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error code 001"})
		return
	}

	session, err := h.Service.RegisterAdmin(req.Email, req.Password, req.FullName, tenantId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "registered",
		"token":   session.Token,
	})
}

func (h *AuthAdminHandler) Login(c *gin.Context) {
	tenantId := c.GetHeader("x-tenant-id")
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		TenantId string `json:"tenantId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if tenantId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error code 001"})
		return
	}

	session, err := h.Service.LoginAdmin(req.Email, req.Password, tenantId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": session.Token})
}

func (h *AuthAdminHandler) Me(c *gin.Context) {
	sessionVal, exists := c.Get("session")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "session tidak ditemukan"})
		return
	}
	session := sessionVal.(*models.Session)

	admin, err := h.Service.AdminRepo.FindByID(*session.TenantUserId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "pengguna tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       admin.ID,
		"email":    admin.Email,
		"fullName": admin.FullName,
		"avatar":   admin.Avatar,
		"role":     admin.Role,
		"tenantId": admin.TenantId,
	})
}
