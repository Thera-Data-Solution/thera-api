package controllers

import (
	"fmt"
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
		fmt.Println("Terjadi kesalahan", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}
	if r.TenantId == nil {
		fmt.Println("Terjadi kesalahan", r)
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid input"})
		return
	}

	_, errnotenant := ac.TenantRepo.FindTenantById(*r.TenantId)
	if errnotenant != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Terjadi kesalahan sistem"})
		return
	}

	existing, _ := ac.TenantUserRepo.FindTenantUserByEmailAndTenant(r.Email, r.TenantId)
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

	ec := ac.TenantUserRepo.CreateTenantUser(tenantUser)
	if ec != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": ec.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	byemail, err := ac.TenantUserRepo.FindTenantUserByEmail(r.Email)
	if err != nil || byemail == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	if byemail.Role == nil || *byemail.Role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid role"})
		return
	}

	currentRole := *byemail.Role
	hasReqTenantId := r.TenantId != nil && *r.TenantId != ""

	var user *models.TenantUser

	switch currentRole {
	case "SU":
		user = byemail

	case "ADMIN":
		if !hasReqTenantId {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tenant ID wajib diisi untuk ADMIN"})
			return
		}

		tu, err := ac.TenantUserRepo.FindTenantUserByEmailAndTenant(r.Email, r.TenantId)
		if err != nil || tu == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
			return
		}
		user = tu

	case "USER":
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: User role tidak diizinkan login di sini"})
		return

	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid role access"})
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	expHours := 24
	if v := os.Getenv("SESSION_EXPIRE_HOURS"); v != "" {
		if h, err := strconv.Atoi(v); err == nil {
			expHours = h
		}
	}
	exp := time.Now().Add(time.Hour * time.Duration(expHours)).UTC()

	sessionTenantId := r.TenantId
	if sessionTenantId == nil || *sessionTenantId == "" {
		sessionTenantId = user.TenantId
	}

	if err := ac.SessionRepo.DeleteByTenantUserId(user.ID); err != nil {
		fmt.Println("⚠️  Warning: gagal menghapus session lama:", err)
	}

	s := &models.Session{
		ID:           uuid.New().String(),
		TenantUserId: &user.ID,
		Token:        uuid.New().String(),
		Device:       r.Device,
		IP:           r.IP,
		ExpiresAt:    exp.Format(time.RFC3339),
		TenantId:     sessionTenantId,
	}

	if err := ac.SessionRepo.CreateSession(s); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat session", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":     s.Token,
		"expiresAt": s.ExpiresAt,
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"fullName": user.FullName,
			"role":     user.Role,
			"tenantId": user.TenantId,
		},
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
