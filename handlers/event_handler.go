package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"thera-api/services"
	"thera-api/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"gorm.io/datatypes"
)

type EventsHandler struct {
	Service *services.EventsService
}

func (h *EventsHandler) Create(c *gin.Context) {
	authData, _ := c.Get("auth")
	auth := authData.(gin.H)
	tenantId := auth["tenantId"].(string)

	name := c.PostForm("name")
	nameEn := c.PostForm("nameEn")
	description := c.PostForm("description")
	descriptionEn := c.PostForm("descriptionEn")
	slug := c.PostForm("slug")
	price := utils.ParseFloat64(c.PostForm("price"))
	capacity := utils.ParseInt(c.PostForm("capacity"))
	status := c.PostForm("status")
	if status == "" {
		status = "available"
	}

	startAt, _ := time.Parse(time.RFC3339, c.PostForm("startAt"))
	endAt, _ := time.Parse(time.RFC3339, c.PostForm("endAt"))

	customFieldsRaw := c.PostForm("customFields")
	var customFields datatypes.JSON
	if customFieldsRaw != "" {
		customFields = datatypes.JSON(customFieldsRaw)
	}

	var imageURL *string
	file, fileHeader, err := c.Request.FormFile("image")
	if err == nil {
		uploader, _ := utils.NewMinIOUploader()
		url, err := uploader.UploadFile(c, file, fileHeader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal upload gambar"})
			return
		}
		imageURL = &url
	}

	event, err := h.Service.CreateEvent(
		name, nameEn, &description, &descriptionEn, slug, imageURL,
		price, startAt, endAt, capacity, status, &tenantId, customFields,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, event)
}

func (h *EventsHandler) GetAll(c *gin.Context) {
	tenantId := c.GetHeader("x-tenant-id")
	if tenantId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenant_id required"})
		return
	}

	events, err := h.Service.GetAllEvents(tenantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, events)
}

func (h *EventsHandler) GetByID(c *gin.Context) {
	tenantId := c.GetHeader("x-tenant-id")
	id := c.Param("id")
	event, err := h.Service.GetEventByID(id, tenantId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, event)
}

func (h *EventsHandler) Update(c *gin.Context) {
	authData, _ := c.Get("auth")
	auth := authData.(gin.H)
	tenantId := auth["tenantId"].(string)
	id := c.Param("id")

	updates := make(map[string]interface{})

	// Handle Image Upload & Old Image Deletion
	file, fileHeader, err := c.Request.FormFile("image")
	if err == nil {
		uploader, _ := utils.NewMinIOUploader()
		url, err := uploader.UploadFile(c, file, fileHeader)
		if err == nil {
			updates["image"] = &url

			// Hapus gambar lama
			oldEvent, _ := h.Service.GetEventByID(id, tenantId)
			if oldEvent != nil && oldEvent.Image != nil && *oldEvent.Image != "" {
				oldObject := strings.TrimPrefix(*oldEvent.Image, fmt.Sprintf("%s/%s/", strings.TrimRight(uploader.Endpoint, "/"), uploader.BucketName))
				_ = uploader.Client.RemoveObject(c, uploader.BucketName, oldObject, minio.RemoveObjectOptions{})
			}
		}
	}

	if name := c.PostForm("name"); name != "" {
		updates["name"] = name
	}
	if nameEn := c.PostForm("nameEn"); nameEn != "" {
		updates["nameEn"] = nameEn
	}
	if desc := c.PostForm("description"); desc != "" {
		updates["description"] = &desc
	}
	if descEn := c.PostForm("descriptionEn"); descEn != "" {
		updates["descriptionEn"] = &descEn
	}
	if price := c.PostForm("price"); price != "" {
		updates["price"] = utils.ParseFloat64(price)
	}
	if cap := c.PostForm("capacity"); cap != "" {
		updates["capacity"] = utils.ParseInt(cap)
	}
	if stat := c.PostForm("status"); stat != "" {
		updates["status"] = stat
	}

	if start := c.PostForm("startAt"); start != "" {
		t, _ := time.Parse(time.RFC3339, start)
		updates["startAt"] = t
	}
	if end := c.PostForm("endAt"); end != "" {
		t, _ := time.Parse(time.RFC3339, end)
		updates["endAt"] = t
	}

	event, err := h.Service.UpdateEvent(id, tenantId, updates)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, event)
}

func (h *EventsHandler) Delete(c *gin.Context) {
	authData, _ := c.Get("auth")
	auth := authData.(gin.H)
	id := c.Param("id")
	if err := h.Service.DeleteEvent(id, auth["tenantId"].(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "event berhasil dihapus"})
}
