package handlers

import (
	"net/http"
	"thera-api/services"

	"github.com/gin-gonic/gin"
)

type AuthAdminHandler struct {
	Service       *services.AuthAdminService
	TenantService *services.TenantService
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
