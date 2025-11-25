package handlers

import (
	"net/http"
	"strconv"
	"thera-api/logger"
	"thera-api/services"
	"thera-api/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PartnerHandler struct {
	Service *services.PartnerService
}

func (h *PartnerHandler) Create(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}
	auth := authData.(gin.H)
	tenantId := auth["tenantId"].(string)

	var logoURL *string
	file, fileHeader, err := c.Request.FormFile("logo")
	if err == nil {
		uploader, _ := utils.NewMinIOUploader()
		url, err := uploader.UploadFile(c, file, fileHeader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal upload logo", "detail": err.Error()})
			return
		}
		logoURL = &url
	}

	partner, err := h.Service.CreatePartner(
		logoURL, &tenantId,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, partner)
}

func (h *PartnerHandler) Update(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}
	auth := authData.(gin.H)
	id := c.Param("id")
	tenantId := auth["tenantId"].(string)

	var logoURL *string
	file, fileHeader, err := c.Request.FormFile("logo")
	if err == nil {
		uploader, _ := utils.NewMinIOUploader()
		url, err := uploader.UploadFile(c, file, fileHeader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal upload logo", "detail": err.Error()})
			return
		}
		logoURL = &url
	}

	partner, err := h.Service.UpdatePartner(
		id, logoURL, &tenantId,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, partner)
}

func (h *PartnerHandler) Delete(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}
	auth := authData.(gin.H)
	id := c.Param("id")
	if err := h.Service.Delete(id, auth["tenantId"].(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "partner berhasil dihapus"})
}

func (h *PartnerHandler) GetAll(c *gin.Context) {
	tenantID := c.GetHeader("x-tenant-id")
	if tenantID == "" {
		logger.Log.Warn("GetAll Partner: tenantId header is missing")
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

	partner, err := h.Service.GetAllPartners(tenantID, page, pageSize)
	if err != nil {
		logger.Log.Error("Failed to get all partner with pagination via service", zap.Error(err), zap.String("tenantId", tenantID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, partner)
}
