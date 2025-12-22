package handlers

import (
	"net/http"
	"thera-api/models"
	"thera-api/services"

	"github.com/gin-gonic/gin"
)

type AuthAdminHandler struct {
	Service            *services.AuthAdminService
	TenantService      *services.TenantService
	PasswordResetService *services.PasswordResetService
	ProfileService     *services.ProfileService
}

func (h *AuthAdminHandler) Register(c *gin.Context) {
	tenantId := c.GetHeader("x-tenant-id")
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		FullName string `json:"fullName"`
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

func (h *AuthAdminHandler) ForgotPassword(c *gin.Context) {
	tenantId := c.GetHeader("x-tenant-id")
	var req struct {
		Email string `json:"email" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if tenantId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error code 001"})
		return
	}

	// Get IP address from request
	ipAddress := c.ClientIP()
	if ipAddress == "" {
		ipAddress = c.GetHeader("X-Forwarded-For")
		if ipAddress == "" {
			ipAddress = c.GetHeader("X-Real-Ip")
		}
	}

	err := h.PasswordResetService.ForgotPasswordAdmin(req.Email, tenantId, ipAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Always return success to prevent email enumeration
	c.JSON(http.StatusOK, gin.H{
		"message": "Jika email terdaftar, instruksi reset password telah dikirim ke email Anda",
	})
}

func (h *AuthAdminHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.PasswordResetService.ResetPasswordAdmin(req.Token, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password berhasil direset",
	})
}

func (h *AuthAdminHandler) UpdateProfile(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}

	auth := authData.(gin.H)
	userType := auth["userType"].(string)

	if userType != "tenant" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Admin tidak valid"})
		return
	}

	var adminIdentifier string
	if uid, ok := auth["tenantUserId"].(*string); ok && uid != nil {
		adminIdentifier = *uid
	} else if admin, ok := auth["user"].(*models.TenantUser); ok && admin != nil {
		adminIdentifier = admin.ID
	}

	if adminIdentifier == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Admin ID tidak ditemukan"})
		return
	}

	var req struct {
		FullName *string `json:"fullName"`
		Avatar   *string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	admin, err := h.ProfileService.UpdateAdminProfile(
		adminIdentifier,
		req.FullName,
		req.Avatar,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile berhasil diupdate",
		"data":    admin,
	})
}
