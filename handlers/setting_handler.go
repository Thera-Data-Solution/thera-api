package handlers

import (
	"net/http"
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
	tenantId := c.GetHeader("x-tenant-id")
	if tenantId == "" {
		logger.Log.Warn("GetAll Settings: tenantId header is missing")
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenantId tidak ditemukan di header"})
		return
	}

	settings, err := h.service.FindByTenantId(tenantId)
	response := dto.SettingGetResponse{
		AppName:        settings.AppName,
		AppLogo:        settings.AppLogo,
		AppTitle:       settings.AppTitle,
		AppDescription: settings.AppDescription,
		AppTheme:       settings.AppTheme,
		Timezone:       settings.Timezone,
	}
	if err != nil {
		logger.Log.Error("failed to get settings", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get settings"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *SettingHandler) GetAllAdmin(c *gin.Context) {
	tenantId := c.GetHeader("x-tenant-id")
	_, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}
	if tenantId == "" {
		logger.Log.Warn("GetAll Settings: tenantId header is missing")
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenantId tidak ditemukan di header"})
		return
	}

	settings, err := h.service.FindByTenantId(tenantId)

	if err != nil {
		logger.Log.Error("failed to get settings", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get settings"})
		return
	}

	c.JSON(http.StatusOK, settings)
}

func (h *SettingHandler) GetById(c *gin.Context) {
	id := c.Param("id")

	setting, err := h.service.FindById(id)
	response := dto.SettingGetResponse{
		AppName:        setting.AppName,
		AppLogo:        setting.AppLogo,
		AppTitle:       setting.AppTitle,
		AppDescription: setting.AppDescription,
		AppTheme:       setting.AppTheme,
		Timezone:       setting.Timezone,
	}
	if err != nil {
		logger.Log.Error("failed to get setting by id", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Setting not found"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *SettingHandler) GetByTenantId(c *gin.Context) {
	tenantId := c.Param("tenantId")

	setting, err := h.service.FindByTenantId(tenantId)
	response := dto.SettingGetResponse{
		AppName:        setting.AppName,
		AppLogo:        setting.AppLogo,
		AppTitle:       setting.AppTitle,
		AppDescription: setting.AppDescription,
		AppTheme:       setting.AppTheme,
		Timezone:       setting.Timezone,
	}
	if err != nil {
		logger.Log.Error("failed to get setting by tenant id", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Setting not found"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *SettingHandler) Upsert(c *gin.Context) {
	var dto dto.SettingRequestBody
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}
	auth := authData.(gin.H)
	tenantId := auth["tenantId"].(string)
	dto.TenantId = &tenantId

	dto.AppName = c.PostForm("appName")
	dto.AppTitle = c.PostForm("appTitle")
	appDescription := c.PostForm("appDescription")
	dto.AppDescription = &appDescription
	appTheme := c.PostForm("appTheme")
	dto.AppTheme = &appTheme
	dto.EnableChatBot = c.PostForm("enableChatBot") == "true"
	dto.EnableFacilitator = c.PostForm("enableFacilitator") == "true"
	dto.EnablePaymentGateway = c.PostForm("enablePaymentGateway") == "true"
	timezone := c.PostForm("timezone")
	dto.Timezone = &timezone
	telegramBotToken := c.PostForm("telegramBotToken")
	if telegramBotToken != "" {
		dto.TelegramBotToken = &telegramBotToken
	}

	telegramChatId := c.PostForm("telegramChatId")
	if telegramChatId != "" {
		dto.TelegramChatId = &telegramChatId
	}

	mailKey := c.PostForm("mailKey")
	if mailKey != "" {
		dto.MailKey = &mailKey
	}

	mailSecret := c.PostForm("mailSecret")
	if mailSecret != "" {
		dto.MailSecret = &mailSecret
	}

	discordReportId := c.PostForm("discordReportId")
	if discordReportId != "" {
		dto.DiscordReportId = &discordReportId
	}

	// Handle file upload
	file, fileHeader, err := c.Request.FormFile("appLogo")
	if err == nil {
		uploader, minioErr := utils.NewMinIOUploader()
		if minioErr != nil {
			logger.Log.Error("Failed to initialize MinIO uploader", zap.Error(minioErr))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal inisialisasi uploader MinIO"})
			return
		}
		url, minioErr := uploader.UploadWithoutDecode(c, file, fileHeader)
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

	c.JSON(http.StatusOK, gin.H{"message": "Setting upserted successfully", "data": setting})
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
