package handlers

import (
	"net/http"
	"thera-api/services"

	"github.com/gin-gonic/gin"
)

type TenantHandler struct {
	Service *services.TenantService
}

// POST /tenants
func (h *TenantHandler) Create(c *gin.Context) {
	var req struct {
		Name string  `json:"name"`
		Logo *string `json:"logo"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenant, err := h.Service.CreateTenant(req.Name, req.Logo)
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
	var req struct {
		Name     *string `json:"name"`
		Logo     *string `json:"logo"`
		IsActive *bool   `json:"isActive"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenant, err := h.Service.UpdateTenant(id, req.Name, req.Logo, req.IsActive)
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
