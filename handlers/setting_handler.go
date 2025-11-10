package handlers

import (
	"net/http"
	"strconv"
	"thera-api/dto"
	"thera-api/logger"
	"thera-api/services"
	"thera-api/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SettingHandler struct {
	service services.SettingService
}

func NewSettingHandler(service services.SettingService) *SettingHandler {
	return &SettingHandler{service: service}
}

func (h *SettingHandler) GetAll(c *gin.Context) {
	settings, err := h.service.FindAll()
	if err != nil {
		logger.Log.Error("failed to get all settings", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get all settings"})
		return
	}

	c.JSON(http.StatusOK, settings)
}

func (h *SettingHandler) GetById(c *gin.Context) {
	id := c.Param("id")

	setting, err := h.service.FindById(id)
	if err != nil {
		logger.Log.Error("failed to get setting by id", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Setting not found"})
		return
	}

	c.JSON(http.StatusOK, setting)
}

func (h *SettingHandler) GetByTenantId(c *gin.Context) {
	tenantId := c.Param("tenantId")

	setting, err := h.service.FindByTenantId(tenantId)
	if err != nil {
		logger.Log.Error("failed to get setting by tenant id", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Setting not found"})
		return
	}

	c.JSON(http.StatusOK, setting)
}

func (h *SettingHandler) Upsert(c *gin.Context) {
	var dto dto.SettingRequestBody

	// Parse form data manually
	dto.AppName = c.PostForm("appName")
	dto.AppTitle = c.PostForm("appTitle")
	appDescription := c.PostForm("appDescription")
	dto.AppDescription = &appDescription
	appTheme := c.PostForm("appTheme")
	dto.AppTheme = &appTheme
	dto.AppMainColor = c.PostForm("appMainColor")
	dto.AppHeaderColor = c.PostForm("appHeaderColor")
	dto.AppFooterColor = c.PostForm("appFooterColor")
	fontSize, _ := strconv.Atoi(c.PostForm("fontSize"))
	dto.FontSize = fontSize
	appDecoration := c.PostForm("appDecoration")
	dto.AppDecoration = &appDecoration
	dto.EnableChatBot = c.PostForm("enableChatBot") == "true"
	dto.EnableFacilitator = c.PostForm("enableFacilitator") == "true"
	dto.EnablePaymentGateway = c.PostForm("enablePaymentGateway") == "true"
	metaOg := c.PostForm("metaOg")
	dto.MetaOg = &metaOg
	timezone := c.PostForm("timezone")
	dto.Timezone = &timezone
	tenantId := c.PostForm("tenantId")
	dto.TenantId = &tenantId

	// Handle file upload
	file, fileHeader, err := c.Request.FormFile("appLogo")
	if err == nil {
		uploader, minioErr := utils.NewMinIOUploader()
		if minioErr != nil {
			logger.Log.Error("Failed to initialize MinIO uploader", zap.Error(minioErr))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal inisialisasi uploader MinIO"})
			return
		}
		url, minioErr := uploader.UploadFile(c, file, fileHeader)
		if minioErr != nil {
			logger.Log.Error("failed to upload file", zap.Error(minioErr))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
			return
		}
		dto.AppLogo = url
	}

	setting, err := h.service.Upsert(dto)
	if err != nil {
		logger.Log.Error("failed to upsert setting", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upsert setting"})
		return
	}

	c.JSON(http.StatusOK, setting)
}

func (h *SettingHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(id); err != nil {
		logger.Log.Error("failed to delete setting", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete setting"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Setting deleted successfully"})
}
