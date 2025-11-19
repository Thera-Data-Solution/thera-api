package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"thera-api/logger"
	"thera-api/services"
	"thera-api/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

type GalleryHandler struct {
	Service *services.GalleryService
}

func (h *GalleryHandler) Create(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}
	auth := authData.(gin.H)
	title := c.PostForm("title")
	description := c.PostForm("description")
	createdAt := time.Now()
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

	gallery, err := h.Service.CreateGallery(
		&title, &description, imageURL, createdAt, &tenantId,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gallery)
}

func (h *GalleryHandler) Update(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}
	id := c.Param("id")
	auth := authData.(gin.H)
	title := c.PostForm("title")
	description := c.PostForm("description")
	createdAt := time.Now()
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

		oldGallery, _ := h.Service.GetGalleryByIDAndTenant(id, tenantId)
		if oldGallery != nil && oldGallery.ImageUrl != "" {
			oldObject := strings.TrimPrefix(oldGallery.ImageUrl, fmt.Sprintf("%s/%s/", strings.TrimRight(uploader.Endpoint, "/"), uploader.BucketName))
			_ = uploader.Client.RemoveObject(c, uploader.BucketName, oldObject, minio.RemoveObjectOptions{})
		}
	}

	gallery, err := h.Service.UpdateGallery(
		id, &title, &description, imageURL, createdAt, tenantId,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gallery)
}

func (h *GalleryHandler) Delete(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}
	auth := authData.(gin.H)
	id := c.Param("id")
	if err := h.Service.DeleteGallery(id, auth["tenantId"].(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "gallery berhasil dihapus"})
}

func (h *GalleryHandler) GetAll(c *gin.Context) {
	tenantID := c.GetHeader("x-tenant-id")
	if tenantID == "" {
		logger.Log.Warn("GetAll Gallery: tenantId header is missing")
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenantId tidak ditemukan di header"})
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	gallery, err := h.Service.GetAllGallery(tenantID, page, pageSize)
	if err != nil {
		logger.Log.Error("Failed to get all gallery with pagination via service", zap.Error(err), zap.String("tenantId", tenantID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gallery)
}

func (h *GalleryHandler) GetByID(c *gin.Context) {
	tenantId := c.GetHeader("x-tenant-id")
	id := c.Param("id")
	gallery, err := h.Service.GetGalleryByID(id, tenantId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "kategori tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gallery)
}
