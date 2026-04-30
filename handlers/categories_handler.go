package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"thera-api/dto"
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
	catType := utils.ParseInt(c.PostForm("catType"))
	showBanner := c.PostForm("showBanner") == "false"

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
		catType, showBanner,
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

	data, err := h.Service.GetAllCategories(tenantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var response []dto.CategoriesResponse
	for _, cat := range data {
		res := dto.CategoriesResponse{
			ID:             cat.ID,
			Name:           cat.Name,
			Slug:           cat.Slug,
			Start:          cat.Start,
			End:            cat.End,
			IsGroup:        cat.IsGroup,
			IsFree:         cat.IsFree,
			IsPayAsYouWish: cat.IsPayAsYouWish,
			IsManual:       cat.IsManual,
		}
		if cat.NameEn == "" {
			res.NameEn = cat.Name
		}
		if cat.Description != nil {
			res.Description = *cat.Description
		}
		if cat.DescriptionEn != nil {
			res.DescriptionEn = *cat.DescriptionEn
		}
		if cat.Location != nil {
			res.Location = *cat.Location
		}
		if cat.Image != nil {
			res.Image = *cat.Image
		}
		if cat.Price != nil {
			res.Price = int(*cat.Price)
		}

		response = append(response, res)
	}

	c.JSON(http.StatusOK, response)
}

func (h *CategoriesHandler) GetCategoriesWithType(c *gin.Context) {
	tenantId := c.GetHeader("x-tenant-id")
	id := c.Param("id")

	if tenantId == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error 01"})
		return
	}

	data, err := h.Service.GetAllCategoriesWithType(tenantId, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Inisialisasi dengan make agar jika data kosong, JSON yang keluar tetap [] bukan null
	response := make([]dto.CategoriesResponse, 0)

	for _, cat := range data {
		res := dto.CategoriesResponse{
			ID:             cat.ID,
			Name:           cat.Name,
			Slug:           cat.Slug,
			Start:          cat.Start,
			End:            cat.End,
			IsGroup:        cat.IsGroup,
			IsFree:         cat.IsFree,
			IsPayAsYouWish: cat.IsPayAsYouWish,
			IsManual:       cat.IsManual,
		}

		if cat.NameEn == "" {
			res.NameEn = cat.Name
		} else {
			res.NameEn = cat.NameEn
		}

		if cat.Description != nil {
			res.Description = *cat.Description
		}
		if cat.DescriptionEn != nil {
			res.DescriptionEn = *cat.DescriptionEn
		}
		if cat.Location != nil {
			res.Location = *cat.Location
		}
		if cat.Image != nil {
			res.Image = *cat.Image
		}
		if cat.Price != nil {
			res.Price = int(*cat.Price)
		}

		response = append(response, res)
	}

	// Mengembalikan response (otomatis [] jika len 0 karena sudah di-make)
	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
}

func (h *CategoriesHandler) GetAllAsAdmin(c *gin.Context) {
	tenantId := c.GetHeader("x-tenant-id")
	id := c.Param("id")

	if tenantId == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error 01"})
		return
	}

	data, err := h.Service.GetAllCategoriesAsAdmin(tenantId, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Inisialisasi dengan make agar jika data kosong, JSON yang keluar tetap [] bukan null
	response := make([]dto.AdminCategoriesResponse, 0)

	for _, cat := range data {
		res := dto.AdminCategoriesResponse{
			ID:             cat.ID,
			Name:           cat.Name,
			Slug:           cat.Slug,
			Start:          cat.Start,
			End:            cat.End,
			IsGroup:        cat.IsGroup,
			IsFree:         cat.IsFree,
			IsPayAsYouWish: cat.IsPayAsYouWish,
			IsManual:       cat.IsManual,
			Disable:        cat.Disable,
		}

		if cat.NameEn == "" {
			res.NameEn = cat.Name
		} else {
			res.NameEn = cat.NameEn
		}

		if cat.Description != nil {
			res.Description = *cat.Description
		}
		if cat.DescriptionEn != nil {
			res.DescriptionEn = *cat.DescriptionEn
		}
		if cat.Location != nil {
			res.Location = *cat.Location
		}
		if cat.Image != nil {
			res.Image = *cat.Image
		}
		if cat.Price != nil {
			res.Price = int(*cat.Price)
		}

		response = append(response, res)
	}

	// Mengembalikan response (otomatis [] jika len 0 karena sudah di-make)
	c.JSON(http.StatusOK, gin.H{
		"data": response,
	})
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
	nameRaw := strings.TrimSpace(c.PostForm("name"))
	nameEn := strings.TrimSpace(c.PostForm("nameEn"))
	description := c.PostForm("description")
	descriptionEn := c.PostForm("descriptionEn")
	slugInput := strings.TrimSpace(c.PostForm("slug"))
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

	var slugPtr *string
	if slugInput != "" {
		slugPtr = &slugInput
	}

	var namePtr *string
	var nameEnPtr *string
	if nameRaw != "" {
		namePtr = &nameRaw
	}

	if nameEn != "" {
		nameEnPtr = &nameEn
	}

	category, err := h.Service.UpdateCategory(
		id, namePtr, nameEnPtr, &description, &descriptionEn, slugPtr, imageURL, &start, &end,
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
