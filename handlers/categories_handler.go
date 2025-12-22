package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"thera-api/models"
	"thera-api/services"
	"thera-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"gorm.io/datatypes"
)

type CategoriesHandler struct {
	Service *services.CategoriesService
}

func (h *CategoriesHandler) Create(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}
	auth := authData.(gin.H)
	name := c.PostForm("name")
	description := c.PostForm("description")
	descriptionEn := c.PostForm("descriptionEn")
	slug := c.PostForm("slug")
	start := utils.ParseInt(c.PostForm("start"))
	end := utils.ParseInt(c.PostForm("end"))
	location := c.PostForm("location")
	price := utils.ParseFloat64(c.PostForm("price"))
	isGroup := c.PostForm("isGroup") == "true"
	isFree := c.PostForm("isFree") == "true"
	isPayAsYouWish := c.PostForm("isPayAsYouWish") == "true"
	isManual := c.PostForm("isManual") == "true"
	disable := c.PostForm("disable") == "true"
	customFieldsRaw := c.PostForm("customFields")

	tenantId := auth["tenantId"].(string)

	var customFields datatypes.JSON

	if customFieldsRaw != "" {
		var temp []models.CategoryCustomField

		if err := json.Unmarshal([]byte(customFieldsRaw), &temp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "format customFields tidak valid",
			})
			return
		}

		customFields = datatypes.JSON(customFieldsRaw)
	}

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

	category, err := h.Service.CreateCategory(
		name, &description, &descriptionEn, slug, imageURL, start, end, &location, &price,
		isGroup, isFree, isPayAsYouWish, isManual, disable, &tenantId, customFields,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, category)
}

func (h *CategoriesHandler) GetAll(c *gin.Context) {
	tenantId := c.GetHeader("x-tenant-id")
	if tenantId == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error 01"})
		return
	}
	categories, err := h.Service.GetAllCategories(tenantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, categories)
}

func (h *CategoriesHandler) GetByID(c *gin.Context) {
	tenantId := c.GetHeader("x-tenant-id")
	id := c.Param("id")
	category, err := h.Service.GetCategoryByID(id, tenantId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "kategori tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, category)
}

// PUT /categories/:id
func (h *CategoriesHandler) Update(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}
	auth := authData.(gin.H)
	id := c.Param("id")
	name := c.PostForm("name")
	description := c.PostForm("description")
	descriptionEn := c.PostForm("descriptionEn")
	slug := c.PostForm("slug")
	start := utils.ParseInt(c.PostForm("start"))
	end := utils.ParseInt(c.PostForm("end"))
	location := c.PostForm("location")
	price := utils.ParseFloat64(c.PostForm("price"))
	isGroup := c.PostForm("isGroup") == "true"
	isFree := c.PostForm("isFree") == "true"
	isPayAsYouWish := c.PostForm("isPayAsYouWish") == "true"
	isManual := c.PostForm("isManual") == "true"
	disable := c.PostForm("disable") == "true"
	tenantId := auth["tenantId"].(string)
	customFieldsRaw := c.PostForm("customFields")

	var imageURL *string
	var customFields datatypes.JSON

	if customFieldsRaw != "" {
		var temp []models.CategoryCustomField

		if err := json.Unmarshal([]byte(customFieldsRaw), &temp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "format customFields tidak valid",
			})
			return
		}

		customFields = datatypes.JSON(customFieldsRaw)
	}

	file, fileHeader, err := c.Request.FormFile("image")
	if err == nil {
		uploader, _ := utils.NewMinIOUploader()
		url, err := uploader.UploadFile(c, file, fileHeader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal upload gambar", "detail": err.Error()})
			return
		}
		imageURL = &url

		oldCategory, _ := h.Service.GetCategoryByIDAndTenant(id, tenantId)
		if oldCategory != nil && oldCategory.Image != nil && *oldCategory.Image != "" {
			oldObject := strings.TrimPrefix(*oldCategory.Image, fmt.Sprintf("%s/%s/", strings.TrimRight(uploader.Endpoint, "/"), uploader.BucketName))
			_ = uploader.Client.RemoveObject(c, uploader.BucketName, oldObject, minio.RemoveObjectOptions{})
		}
	}

	category, err := h.Service.UpdateCategory(
		id, &name, &description, &descriptionEn, &slug, imageURL, &start, &end,
		&location, &price, &isGroup, &isFree, &isPayAsYouWish, &isManual, &disable, tenantId, customFields,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, category)
}

// DELETE /categories/:id
func (h *CategoriesHandler) Delete(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}
	auth := authData.(gin.H)
	id := c.Param("id")
	if err := h.Service.DeleteCategory(id, auth["tenantId"].(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "kategori berhasil dihapus"})
}
