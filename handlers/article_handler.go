// handlers/article_handler.go
package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"thera-api/dto"
	"thera-api/logger"
	"thera-api/services"
	"thera-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ArticleHandler struct {
	Service *services.ArticleService
}

// NewArticleHandler creates a new instance of ArticleHandler
func NewArticleHandler(service *services.ArticleService) *ArticleHandler {
	return &ArticleHandler{Service: service}
}

// Create handles the creation of a new Article using form-data (including file upload)
func (h *ArticleHandler) Create(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}
	auth := authData.(gin.H)
	tenantID := auth["tenantId"].(string)

	// Ambil data dari form-data
	title := c.PostForm("title")
	slug := c.PostForm("slug")
	bodyContent := c.PostForm("body")
	excerptContent := c.PostForm("excerpt")
	isPublishedStr := c.PostForm("isPublished")
	articleTypeStr := c.PostForm("articleType")

	// Konversi string ke tipe data yang sesuai
	isPublished := isPublishedStr == "true"
	var body *string
	if bodyContent != "" {
		body = &bodyContent
	}
	var excerpt *string
	if excerptContent != "" {
		excerpt = &excerptContent
	}
	var articleType *string
	if articleTypeStr != "" {
		articleType = &articleTypeStr
	} else {
		// articleType is required, so if empty, return error
		logger.Log.Error("Failed to create article: articleType is missing")
		c.JSON(http.StatusBadRequest, gin.H{"error": "articleType tidak boleh kosong"})
		return
	}

	// Validasi input dasar (sesuai binding:required di DTO, kita lakukan manual di sini)
	if title == "" || slug == "" {
		logger.Log.Error("Failed to create article: title or slug is missing")
		c.JSON(http.StatusBadRequest, gin.H{"error": "title dan slug tidak boleh kosong"})
		return
	}
	// Tambahan: Validasi ArticleType agar isinya PAGE atau BLOG
	if *articleType != "PAGE" && *articleType != "BLOG" {
		logger.Log.Error("Failed to create article: invalid articleType", zap.String("articleType", *articleType))
		c.JSON(http.StatusBadRequest, gin.H{"error": "articleType harus 'PAGE' atau 'BLOG'"})
		return
	}

	var coverImageURL *string
	file, fileHeader, err := c.Request.FormFile("coverImage")
	if err == nil { // File coverImage ditemukan
		uploader, minioErr := utils.NewMinIOUploader() // Pastikan NewMinIOUploader mengembalikan uploader dan error
		if minioErr != nil {
			logger.Log.Error("Failed to initialize MinIO uploader", zap.Error(minioErr))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal inisialisasi uploader MinIO"})
			return
		}
		url, uploadErr := uploader.UploadFile(c, file, fileHeader)
		if uploadErr != nil {
			logger.Log.Error("Failed to upload cover image", zap.Error(uploadErr))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal upload gambar cover", "detail": uploadErr.Error()})
			return
		}
		coverImageURL = &url
	} else if err != http.ErrMissingFile {
		// Ada error lain selain file tidak ada (misal corrupted upload)
		logger.Log.Error("Error processing cover image file", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "gagal memproses file gambar cover"})
		return
	}
	// Jika http.ErrMissingFile, coverImageURL akan tetap nil, yang berarti tidak ada cover image

	article, err := h.Service.CreateArticle(
		title, slug, coverImageURL, body,
		excerpt, articleType, isPublished, &tenantID,
	)
	if err != nil {
		logger.Log.Error("Failed to create article via service", zap.Error(err), zap.String("tenantId", tenantID))
		// Cek error untuk slug duplikat
		if strings.Contains(err.Error(), "slug sudah ada") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()}) // 409 Conflict
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, article)
}

// GetAll handles retrieving all articles with pagination
// (Tidak ada perubahan, karena ini GET request dan menggunakan query params)
func (h *ArticleHandler) GetAll(c *gin.Context) {
	tenantID := c.GetHeader("x-tenant-id")
	if tenantID == "" {
		logger.Log.Warn("GetAll Articles: tenantId header is missing")
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

	articleTypeQuery := c.Query("type") // Filter by article type (PAGE or BLOG)
	var articleType *string
	if articleTypeQuery != "" {
		articleType = &articleTypeQuery
	}

	articles, err := h.Service.GetAllArticles(tenantID, page, pageSize, articleType)
	if err != nil {
		logger.Log.Error("Failed to get all articles with pagination via service", zap.Error(err), zap.String("tenantId", tenantID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, articles)
}

// GetByID handles retrieving a single article by ID
// (Tidak ada perubahan)
func (h *ArticleHandler) GetByID(c *gin.Context) {
	tenantID := c.GetHeader("x-tenant-id")
	if tenantID == "" {
		logger.Log.Warn("GetArticleByID: tenantId header is missing")
		c.JSON(http.StatusBadRequest, gin.H{"error": "tenantId tidak ditemukan di header"})
		return
	}
	id := c.Param("id")

	article, err := h.Service.GetArticleByID(id, tenantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || err.Error() == "artikel tidak ditemukan" {
			c.JSON(http.StatusNotFound, gin.H{"error": "artikel tidak ditemukan"})
		} else {
			logger.Log.Error("Failed to get article by ID via service", zap.Error(err), zap.String("id", id), zap.String("tenantId", tenantID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, article)
}

// Update handles updating an existing Article using form-data (including file upload)
func (h *ArticleHandler) Update(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}
	auth := authData.(gin.H)
	tenantID := auth["tenantId"].(string)
	id := c.Param("id")

	// Buat DTO untuk update dari form-data
	updateDTO := dto.UpdateArticleDTO{}

	// Ambil nilai dari form-data. Gunakan pointer untuk field opsional.
	if val := c.PostForm("title"); val != "" {
		updateDTO.Title = &val
	}
	if val := c.PostForm("slug"); val != "" {
		updateDTO.Slug = &val
	}
	if val := c.PostForm("body"); val != "" {
		updateDTO.Body = &val
	}
	if val := c.PostForm("excerpt"); val != "" {
		updateDTO.Excerpt = &val
	}
	if val := c.PostForm("isPublished"); val != "" {
		isPublished := val == "true"
		updateDTO.IsPublished = &isPublished
	}
	if val := c.PostForm("articleType"); val != "" {
		// Tambahan: Validasi ArticleType agar isinya PAGE atau BLOG
		if val != "PAGE" && val != "BLOG" {
			logger.Log.Error("Failed to update article: invalid articleType", zap.String("articleType", val))
			c.JSON(http.StatusBadRequest, gin.H{"error": "articleType harus 'PAGE' atau 'BLOG'"})
			return
		}
		updateDTO.ArticleType = &val
	}

	// Cek apakah ada request untuk menghapus coverImage (misalnya dengan mengirim "coverImage": "null" atau string kosong)
	// Jika ada field 'removeCoverImage' = "true" di form data, kita bisa menghapus gambar.
	if c.PostForm("removeCoverImage") == "true" {
		nilString := ""                   // Nilai kosong untuk pointer string
		updateDTO.CoverImage = &nilString // Set coverImage menjadi nil (string kosong)
	}

	var newCoverImageURL *string
	file, fileHeader, err := c.Request.FormFile("coverImage")
	if err == nil { // File coverImage ditemukan
		uploader, minioErr := utils.NewMinIOUploader()
		if minioErr != nil {
			logger.Log.Error("Failed to initialize MinIO uploader", zap.Error(minioErr))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal inisialisasi uploader MinIO"})
			return
		}
		url, uploadErr := uploader.UploadFile(c, file, fileHeader)
		if uploadErr != nil {
			logger.Log.Error("Failed to upload cover image", zap.Error(uploadErr))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal upload gambar cover", "detail": uploadErr.Error()})
			return
		}
		newCoverImageURL = &url
		updateDTO.CoverImage = newCoverImageURL // Set URL gambar baru di DTO

		// Optional: Delete old image if a new one is uploaded
		oldArticle, _ := h.Service.GetArticleByID(id, tenantID)
		if oldArticle != nil && oldArticle.CoverImage != nil && *oldArticle.CoverImage != "" {
			oldObject := strings.TrimPrefix(*oldArticle.CoverImage, fmt.Sprintf("%s/%s/", strings.TrimRight(uploader.Endpoint, "/"), uploader.BucketName))
			_ = uploader.Client.RemoveObject(c, uploader.BucketName, oldObject, minio.RemoveObjectOptions{})
			logger.Log.Info("Old cover image deleted during update", zap.String("oldObject", oldObject))
		}
	} else if err != http.ErrMissingFile {
		logger.Log.Error("Error processing cover image file during update", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "gagal memproses file gambar cover"})
		return
	}
	// Jika http.ErrMissingFile, maka tidak ada file baru yang diunggah.
	// Jika removeCoverImage=true, CoverImage di DTO sudah diatur ke nil.
	// Jika tidak ada file dan tidak ada removeCoverImage, maka CoverImage di DTO tetap nil,
	// yang berarti tidak ada perubahan pada CoverImage yang ada di DB.

	article, err := h.Service.UpdateArticle(id, updateDTO, tenantID)
	if err != nil {
		logger.Log.Error("Failed to update article via service", zap.Error(err), zap.String("id", id), zap.String("tenantId", tenantID))
		if strings.Contains(err.Error(), "slug sudah ada") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()}) // 409 Conflict
		} else if strings.Contains(err.Error(), "artikel tidak ditemukan") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()}) // 404 Not Found
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, article)
}

// Delete handles deleting an Article
// (Tidak ada perubahan)
func (h *ArticleHandler) Delete(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}
	auth := authData.(gin.H)
	tenantID := auth["tenantId"].(string)
	id := c.Param("id")

	// Optional: Delete cover image from MinIO before deleting article record
	articleToDelete, err := h.Service.GetArticleByID(id, tenantID)
	if err == nil && articleToDelete != nil && articleToDelete.CoverImage != nil && *articleToDelete.CoverImage != "" {
		uploader, minioErr := utils.NewMinIOUploader()
		if minioErr != nil {
			logger.Log.Error("Failed to initialize MinIO uploader for deletion", zap.Error(minioErr))
			// Lanjutkan saja, jangan blokir penghapusan artikel karena uploader error
		} else {
			oldObject := strings.TrimPrefix(*articleToDelete.CoverImage, fmt.Sprintf("%s/%s/", strings.TrimRight(uploader.Endpoint, "/"), uploader.BucketName))
			_ = uploader.Client.RemoveObject(c, uploader.BucketName, oldObject, minio.RemoveObjectOptions{})
			logger.Log.Info("Cover image deleted from MinIO during article deletion", zap.String("oldObject", oldObject))
		}
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Log.Error("Failed to fetch article for image deletion during article delete", zap.Error(err), zap.String("id", id))
	}

	if err := h.Service.DeleteArticle(id, tenantID); err != nil {
		logger.Log.Error("Failed to delete article via service", zap.Error(err), zap.String("id", id), zap.String("tenantId", tenantID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "artikel berhasil dihapus"})
}
