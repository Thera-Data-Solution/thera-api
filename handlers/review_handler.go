package handlers

import (
	"net/http"
	"thera-api/dto"
	"thera-api/services"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	Service *services.ReviewService
}

func (h *ReviewHandler) AdminGetAll(c *gin.Context) {
	auth := c.MustGet("auth").(gin.H)
	tenantId := auth["tenantId"].(string)

	reviews, err := h.Service.AdminGetAllReviews(tenantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reviews)
}

func (h *ReviewHandler) AdminUpdate(c *gin.Context) {
	id := c.Param("id")
	var req dto.AdminUpdateReviewRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.AdminUpdateReview(id, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ulasan berhasil diperbarui"})
}

func (h *ReviewHandler) GetMyHistory(c *gin.Context) {
	auth := c.MustGet("auth").(gin.H)
	userId := *auth["userId"].(*string)
	tenantId := auth["tenantId"].(string)
	userType := auth["userType"].(string)

	if userType != "user" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User tidak valid"})
		return
	}

	history, err := h.Service.GetHistory(userId, tenantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
}

func (h *ReviewHandler) SubmitReview(c *gin.Context) {
	auth := c.MustGet("auth").(gin.H)
	userId := *auth["userId"].(*string)
	tenantId := auth["tenantId"].(string)

	var req dto.ReviewUpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.SubmitReview(userId, req, tenantId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "review berhasil disimpan"})
}

func (h *ReviewHandler) GetLandingPageReviews(c *gin.Context) {
	tenantId := c.GetHeader("x-tenant-id")
	reviews, err := h.Service.GetPublicTestimonials(tenantId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, reviews)
}
