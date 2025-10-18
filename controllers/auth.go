package controllers

import (
	"go-api/config"
	"go-api/models"
	"go-api/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Register user baru
func Register(c *gin.Context) {
	var input struct {
		Email    string  `json:"email"`
		Password string  `json:"password"`
		FullName string  `json:"fullName"`
		Phone    string  `json:"phone"`
		Avatar   *string `json:"avatar"`
		Fb       string  `json:"fb"`
		Ig       string  `json:"ig"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek apakah email sudah terdaftar
	var existing models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah terdaftar"})
		return
	}
	defaultAvatar := "https://avatar.iran.liara.run/public"

	if input.Avatar == nil {
		input.Avatar = &defaultAvatar
	}

	if input.Phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please fill your phone number",
		})
		return
	}
	if input.Fb == "" && input.Ig == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Please fill at least one: Instagram or Facebook.",
		})
		return
	}

	hashed, _ := utils.HashPassword(input.Password)

	user := models.User{
		ID:       uuid.New(),
		Email:    input.Email,
		Password: hashed,
		FullName: input.FullName,
		Phone:    input.Phone,
		Avatar:   input.Avatar,
		Ig:       input.Ig,
		Fb:       input.Fb,
		Role:     "USER",
	}

	config.DB.Create(&user)
	token := uuid.NewString()
	session := models.Session{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     token,
		Device:    "Device",
		IP:        c.ClientIP(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := config.DB.Create(&session).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Terjadi kesalahan saat mendaftar", "err": err})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Registrasi berhasil",
		"token":   token,
	})
}

// Login user
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Device   string `json:"device"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email tidak ditemukan"})
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password salah"})
		return
	}

	// Buat session baru
	token := uuid.NewString()
	session := models.Session{
		UserID:    user.ID,
		Token:     token,
		Device:    "Device",
		IP:        c.ClientIP(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	config.DB.Create(&session)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   token,
	})
}

func GetMe(c *gin.Context) {
	userData, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak ditemukan"})
		return
	}

	user := userData.(models.User)
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
