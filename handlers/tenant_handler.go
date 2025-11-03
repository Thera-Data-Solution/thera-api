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

type TenantHandler struct {
	Service *services.TenantService
}

// POST /tenants
func (h *TenantHandler) Create(c *gin.Context) {
	name := c.PostForm("name")
	file, fileHeader, err := c.Request.FormFile("logo")

	var logoURL *string

	if err == nil {
		uploader, _ := utils.NewMinIOUploader()
		url, err := uploader.UploadFile(c, file, fileHeader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":  "gagal upload logo",
				"detail": err.Error(), // tambahkan baris ini
			})
			return
		}

		logoURL = &url
	}

	tenant, err := h.Service.CreateTenant(name, logoURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tenant)
}

// GET /tenants
func (h *TenantHandler) GetAll(c *gin.Context) {
	tenants, err := h.Service.GetAllTenants()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tenants)
}

// GET /tenants/:id
func (h *TenantHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	tenant, err := h.Service.GetTenantByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tenant tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, tenant)
}

func (h *TenantHandler) Update(c *gin.Context) {
	id := c.Param("id")

	// Bind form-data karena kita bisa dapat file di multipart
	name := c.PostForm("name")
	isActiveStr := c.PostForm("isActive")

	var isActive *bool
	if isActiveStr != "" {
		val := isActiveStr == "true"
		isActive = &val
	}

	var logoURL *string

	// Cek apakah ada file logo yang diupload
	file, fileHeader, err := c.Request.FormFile("logo")
	if err == nil {
		// Upload file baru
		uploader, _ := utils.NewMinIOUploader()
		url, err := uploader.UploadFile(c, file, fileHeader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":  "gagal upload logo",
				"detail": err.Error(),
			})
			return
		}
		logoURL = &url

		// Hapus logo lama jika ada
		oldTenant, _ := h.Service.GetTenantByID(id)
		if oldTenant != nil && oldTenant.Logo != nil && *oldTenant.Logo != "" {
			// Extract object name dari URL
			oldObject := strings.TrimPrefix(*oldTenant.Logo, fmt.Sprintf("%s/%s/", strings.TrimRight(uploader.Endpoint, "/"), uploader.BucketName))
			_ = uploader.Client.RemoveObject(c, uploader.BucketName, oldObject, minio.RemoveObjectOptions{})
		}
	}

	// Update tenant
	tenant, err := h.Service.UpdateTenant(id, &name, logoURL, isActive)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tenant)
}

func (h *TenantHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.Service.DeleteTenant(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "tenant berhasil dihapus"})
}
