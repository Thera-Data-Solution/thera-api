package handlers

import (
	"net/http"
	"thera-api/models"
	"thera-api/services"

	"github.com/gin-gonic/gin"
)

type AuthUserHandler struct {
	Service *services.AuthUserService
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

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return

	}

	if tenantId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error code 001"})
		return
	}

	session, err := h.Service.RegisterUser(req.Email, req.Password, req.FullName, req.Phone, req.Address, req.Ig, req.Fb, tenantId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
