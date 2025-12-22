package handlers

import (
	"fmt"
	"net/http"
	"thera-api/services"
	"thera-api/utils"

	"github.com/gin-gonic/gin"
)

type HeroHandler struct {
	Service *services.HeroService
}

func (h *HeroHandler) Create(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}
	auth := authData.(gin.H)
	tenantId := auth["tenantId"].(string)

	title := c.PostForm("title")
	subtitle := c.PostForm("subtitle")
	description := c.PostForm("description")
	buttonText := c.PostForm("buttonText")
	buttonLink := c.PostForm("buttonLink")
	themeType := c.PostForm("themeType")
	isActive := c.PostForm("isActive") == "true"

	// === Sama persis seperti Article ===
	var imageURL *string
	fmt.Println("CONTENT TYPE:", c.Request.Header.Get("Content-Type"))
	c.Request.ParseMultipartForm(50 << 20)
	fmt.Println("FILES:", c.Request.MultipartForm.File)

	file, fileHeader, err := c.Request.FormFile("image")
	if err == nil {
		fmt.Println("ERR FORMFILE:", err)
		uploader, minioErr := utils.NewMinIOUploader()
		if minioErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal inisialisasi uploader MinIO"})
			return
		}

		url, uploadErr := uploader.UploadFile(c, file, fileHeader)
		if uploadErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal upload gambar", "detail": uploadErr.Error()})
			return
		}

		imageURL = &url
	} else if err != http.ErrMissingFile {
		// Error lain selain tidak ada file
		c.JSON(http.StatusBadRequest, gin.H{"error": "gagal memproses file image", "detail": err.Error()})
		return
	}

	hero, err := h.Service.CreateHero(
		title, &subtitle, &description, imageURL,
		&buttonText, &buttonLink, &themeType,
		isActive, &tenantId,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, hero)
}

func (h *HeroHandler) GetAll(c *gin.Context) {
	tenantId := c.GetHeader("x-tenant-id")
	if tenantId == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "tenantId kosong"})
		return
	}
	heroes, err := h.Service.GetAllHeroes(tenantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, heroes)
}
