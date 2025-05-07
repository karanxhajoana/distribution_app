package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"pack-sizer/internal/services"
)

type PackHandler struct {
	packSizesManager *services.PackSizesManager
}

func NewPackHandler(packSizesManager *services.PackSizesManager) *PackHandler {
	return &PackHandler{
		packSizesManager: packSizesManager,
	}
}

// GetPackSizes returns the current pack sizes
func (h *PackHandler) GetPackSizes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"sizes": h.packSizesManager.GetSizes()})
}

// AddPackSize adds a new pack size
func (h *PackHandler) AddPackSize(c *gin.Context) {
	var request struct {
		Size int `json:"size"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if request.Size <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pack size must be a positive integer"})
		return
	}

	h.packSizesManager.AddSize(request.Size)
	c.JSON(http.StatusOK, gin.H{"sizes": h.packSizesManager.GetSizes()})
}

// RemovePackSize removes a pack size
func (h *PackHandler) RemovePackSize(c *gin.Context) {
	size := c.Param("size")
	var sizeInt int
	if _, err := fmt.Sscanf(size, "%d", &sizeInt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pack size"})
		return
	}

	h.packSizesManager.RemoveSize(sizeInt)
	c.JSON(http.StatusOK, gin.H{"sizes": h.packSizesManager.GetSizes()})
}

// UpdatePackSize updates an existing pack size
func (h *PackHandler) UpdatePackSize(c *gin.Context) {
	var request struct {
		OldSize int `json:"oldSize"`
		NewSize int `json:"newSize"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if request.NewSize <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "New pack size must be a positive integer"})
		return
	}

	if err := h.packSizesManager.UpdateSize(request.OldSize, request.NewSize); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sizes": h.packSizesManager.GetSizes()})
}

// CalculatePacks calculates the optimal pack distribution
func (h *PackHandler) CalculatePacks(c *gin.Context) {
	quantity := c.Query("quantity")
	orderQuantity, err := strconv.Atoi(quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantity parameter"})
		return
	}

	if orderQuantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order quantity must be a positive integer"})
		return
	}

	packSizes := h.packSizesManager.GetSizes()
	result := services.CalculatePacks(orderQuantity, packSizes)
	c.JSON(http.StatusOK, result)
} 