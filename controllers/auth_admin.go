package controllers

import (
	"net/http"
	"net/url"
	"os"
	"strconv"
	"thera-api/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (ac *AuthController) AdminRegister(c *gin.Context) {
	type input struct {
		Email    string  `json:"email"`
		Password string  `json:"password"`
		FullName string  `json:"fullName"`
		TenantId *string `json:"tenantId"`
	}
	var r input
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	existing, _ := ac.TenantUserRepo.FindByEmailAndTenant(r.Email, r.TenantId)
	if existing != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email sudah terdaftar untuk tenant ini"})
		return
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	avatar := "https://api.dicebear.com/9.x/fun-emoji/svg?seed=" + url.QueryEscape(r.FullName)
	tenantUser := &models.TenantUser{
		ID:       uuid.New().String(),
		Email:    r.Email,
		Password: string(hashed),
		FullName: r.FullName,
		TenantId: r.TenantId,
		Avatar:   &avatar,
	}

	err := ac.TenantUserRepo.CreateTenantUser(tenantUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": tenantUser})
}

func (ac *AuthController) AdminLogin(c *gin.Context) {
	type req struct {
		Email    string  `json:"email"`
		Password string  `json:"password"`
		TenantId *string `json:"tenantId"`
		Device   *string `json:"device"`
		IP       *string `json:"ip"`
	}
	var r req
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	tu, err := ac.TenantUserRepo.FindByEmailAndTenant(r.Email, r.TenantId)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "email/password wrong"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(tu.Password), []byte(r.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "email/password wrong"})
		return
	}

	token := uuid.New().String()
	expHours := 24
	if v := os.Getenv("SESSION_EXPIRE_HOURS"); v != "" {
		if h, err := strconv.Atoi(v); err == nil {
			expHours = h
		}
	}
	exp := time.Now().Add(time.Hour * time.Duration(expHours)).UTC()

	s := &models.Session{
		ID:           uuid.New().String(),
		TenantUserId: &tu.ID,
		Token:        token,
		Device:       r.Device,
		IP:           r.IP,
		ExpiresAt:    exp.Format(time.RFC3339),
		TenantId:     r.TenantId,
	}
	ac.SessionRepo.DeleteByTenantUserId(tu.ID)
	if err := ac.SessionRepo.CreateSession(s); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat session", "err": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":     token,
		"expiresAt": s.ExpiresAt,
	})
}

func (ac *AuthController) AdminMe(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, user)
}
