package handlers

import (
	"net/http"
	"thera-api/services"

	"github.com/gin-gonic/gin"
)

type TranslationHandler struct {
	Service *services.TranslationService
}

type TranslationRequest struct {
	Locale    string `json:"locale"`
	Namespace string `json:"namespace"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

// POST /translations
func (h *TranslationHandler) Create(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}
	auth := authData.(gin.H)

	var req TranslationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format request tidak valid"})
		return
	}

	t, err := h.Service.CreateTranslation(req.Locale, req.Namespace, req.Key, req.Value, auth["tenantId"].(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, t)
}

// GET /translations
func (h *TranslationHandler) GetAll(c *gin.Context) {
	tenantId := c.GetHeader("x-tenant-id")
	if tenantId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "header x-tenant-id wajib ada"})
		return
	}

	translations, err := h.Service.GetAllTranslations(tenantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, translations)
}

// GET /translations/:id
func (h *TranslationHandler) GetByID(c *gin.Context) {
	tenantId := c.GetHeader("x-tenant-id")
	if tenantId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "header x-tenant-id wajib ada"})
		return
	}

	id := c.Param("id")
	t, err := h.Service.GetTranslationByID(id, tenantId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "data tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, t)
}

type UpdateTranslationRequest struct {
	Locale    *string `json:"locale"`
	Namespace *string `json:"namespace"`
	Key       *string `json:"key"`
	Value     *string `json:"value"`
}

// PUT /translations/:id
func (h *TranslationHandler) Update(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}
	auth := authData.(gin.H)
	id := c.Param("id")

	var req UpdateTranslationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format request tidak valid"})
		return
	}

	t, err := h.Service.UpdateTranslation(id, req.Locale, req.Namespace, req.Key, req.Value, auth["tenantId"].(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, t)
}

// DELETE /translations/:id
func (h *TranslationHandler) Delete(c *gin.Context) {
	authData, exists := c.Get("auth")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth tidak ditemukan"})
		return
	}
	auth := authData.(gin.H)
	id := c.Param("id")

	if err := h.Service.DeleteTranslation(id, auth["tenantId"].(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "translation berhasil dihapus"})
}
