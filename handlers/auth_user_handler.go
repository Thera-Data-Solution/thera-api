package handlers

import (
	"net/http"
	"thera-api/logger"
	"thera-api/models"
	"thera-api/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthUserHandler struct {
	Service            *services.AuthUserService
	PasswordResetService *services.PasswordResetService
	ProfileService     *services.ProfileService
}

func (h *AuthUserHandler) Register(c *gin.Context) {
	tenantId := c.GetHeader("x-tenant-id")
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		FullName string `json:"fullName"`
		Phone    string `json:"phone"`
		Ig       string `json:"ig"`
		Fb       string `json:"fb"`
		Address  string `jsong:"address"`
	}
	logger.Log.Info("Menerima request register user",
		zap.String("email", req.Email),
		zap.String("fullName", req.FullName),
		zap.String("phone", req.Phone),
		zap.String("tenantId", tenantId),
	)

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Log.Error("Gagal bind JSON", zap.Error(err))

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	if tenantId == "" {
		logger.Log.Error("Tenant ID kosong pada request register")

		c.JSON(http.StatusBadRequest, gin.H{"error": "error code 001"})
		return
	}

	session, err := h.Service.RegisterUser(req.Email, req.Password, req.FullName, req.Phone, req.Address, req.Ig, req.Fb, tenantId)
	if err != nil {
		logger.Log.Error("Gagal register user",
			zap.String("email", req.Email),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	logger.Log.Info("Berhasil register user",
		zap.String("email", req.Email),
		zap.String("tenantId", tenantId),
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "registered",
		"token":   session.Token,
	})
}

func (h *AuthUserHandler) Login(c *gin.Context) {
	tenantId := c.GetHeader("x-tenant-id")
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if tenantId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error code 001"})
		return
	}

	session, err := h.Service.LoginUser(req.Email, req.Password, tenantId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": session.Token})
}

func (h *AuthUserHandler) Me(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}

	auth := authData.(gin.H)
	userType := auth["userType"].(string)

	switch userType {
	case "user":
		user := auth["user"].(*models.User)
		c.JSON(http.StatusOK, gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"fullName": user.FullName,
			"avatar":   user.Avatar,
			"phone":    user.Phone,
			"address":  user.Address,
			"ig":       user.Ig,
			"fb":       user.Fb,
		})
	case "tenant":
		user := auth["user"].(*models.TenantUser)
		c.JSON(http.StatusOK, gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"avatar":   user.Avatar,
			"fullName": user.FullName,
			"role":     user.Role,
		})
	default:
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tipe pengguna tidak dikenali"})
	}
}

func (h *AuthUserHandler) ForgotPassword(c *gin.Context) {
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

	err := h.PasswordResetService.ForgotPasswordUser(req.Email, tenantId, ipAddress)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Always return success to prevent email enumeration
	c.JSON(http.StatusOK, gin.H{
		"message": "Jika email terdaftar, instruksi reset password telah dikirim ke email Anda",
	})
}

func (h *AuthUserHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.PasswordResetService.ResetPasswordUser(req.Token, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password berhasil direset",
	})
}

func (h *AuthUserHandler) UpdateProfile(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}

	auth := authData.(gin.H)
	userType := auth["userType"].(string)

	if userType != "user" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak valid"})
		return
	}

	var userIdentifier string
	if uid, ok := auth["userId"].(*string); ok && uid != nil {
		userIdentifier = *uid
	} else if user, ok := auth["user"].(*models.User); ok && user != nil {
		userIdentifier = user.ID
	}

	if userIdentifier == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID tidak ditemukan"})
		return
	}

	var req struct {
		FullName *string `json:"fullName"`
		Phone    *string `json:"phone"`
		Address  *string `json:"address"`
		Ig       *string `json:"ig"`
		Fb       *string `json:"fb"`
		Avatar   *string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.ProfileService.UpdateUserProfile(
		userIdentifier,
		req.FullName,
		req.Phone,
		req.Address,
		req.Ig,
		req.Fb,
		req.Avatar,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile berhasil diupdate",
		"data":    user,
	})
}
