package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"thera-api/services"
	"thera-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
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
	title := c.PostForm("title")
	subtitle := c.PostForm("subtitle")
	description := c.PostForm("description")
	buttonText := c.PostForm("buttonText")
	buttonLink := c.PostForm("buttonLink")
	themeType := c.PostForm("themeType")
	isActive := c.PostForm("isActive") == "true"
	tenantId := auth["tenantId"].(string)

	var imageURL *string
	file, fileHeader, err := c.Request.FormFile("image")
	if err == nil {
		uploader, _ := utils.NewMinIOUploader()
		url, err := uploader.UploadFile(c, file, fileHeader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal upload gambar", "detail": err.Error()})
			return
		}
		imageURL = &url
	}

	hero, err := h.Service.CreateHero(
		title, &subtitle, &description, imageURL, &buttonText, &buttonLink, &themeType, isActive, &tenantId,
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

func (h *HeroHandler) GetByID(c *gin.Context) {
	tenantId := c.GetHeader("x-tenant-id")
	id := c.Param("id")
	hero, err := h.Service.GetHeroByID(id, tenantId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "hero tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, hero)
}

func (h *HeroHandler) Update(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}
	auth := authData.(gin.H)
	id := c.Param("id")
	title := c.PostForm("title")
	subtitle := c.PostForm("subtitle")
	description := c.PostForm("description")
	buttonText := c.PostForm("buttonText")
	buttonLink := c.PostForm("buttonLink")
	themeType := c.PostForm("themeType")
	isActive := c.PostForm("isActive") == "true"
	tenantId := auth["tenantId"].(string)

	var imageURL *string
	file, fileHeader, err := c.Request.FormFile("image")
	if err == nil {
		uploader, _ := utils.NewMinIOUploader()
		url, err := uploader.UploadFile(c, file, fileHeader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal upload gambar", "detail": err.Error()})
			return
		}
		imageURL = &url

		oldHero, _ := h.Service.GetHeroByID(id, tenantId)
		if oldHero != nil && oldHero.Image != nil && *oldHero.Image != "" {
			oldObject := strings.TrimPrefix(*oldHero.Image, fmt.Sprintf("%s/%s/", strings.TrimRight(uploader.Endpoint, "/"), uploader.BucketName))
			_ = uploader.Client.RemoveObject(c, uploader.BucketName, oldObject, minio.RemoveObjectOptions{})
		}
	}

	hero, err := h.Service.UpdateHero(
		id, &title, &subtitle, &description, imageURL, &buttonText, &buttonLink, &themeType, &isActive, tenantId,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, hero)
}

func (h *HeroHandler) Delete(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}
	auth := authData.(gin.H)
	id := c.Param("id")
	if err := h.Service.DeleteHero(id, auth["tenantId"].(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "hero berhasil dihapus"})
}
