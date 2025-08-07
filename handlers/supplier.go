package handlers

import (
	"inventory_system/database"
	"inventory_system/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetSuppliers returns all suppliers with pagination support
func GetSuppliers(c *gin.Context) {
	var suppliers []models.Supplier
	var total int64

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	// Build base query
	query := database.DB.Where("is_active = ?", true)

	// Count total records
	query.Model(&models.Supplier{}).Count(&total)

	// Get paginated results
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&suppliers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch suppliers",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    suppliers,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

// GetSupplier returns a specific supplier by ID
func GetSupplier(c *gin.Context) {
	id := c.Param("id")
	var supplier models.Supplier
	
	if err := database.DB.Where("id = ? AND is_active = ?", id, true).First(&supplier).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Supplier not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    supplier,
	})
}

// CreateSupplier creates a new supplier
func CreateSupplier(c *gin.Context) {
	var supplier models.Supplier
	
	if err := c.ShouldBindJSON(&supplier); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid input data: " + err.Error(),
		})
		return
	}

	// Validate required fields
	if supplier.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Supplier name is required",
		})
		return
	}

	// Set default values
	supplier.IsActive = true

	if err := database.DB.Create(&supplier).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create supplier: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Supplier created successfully",
		"data":    supplier,
	})
}

// UpdateSupplier updates an existing supplier
func UpdateSupplier(c *gin.Context) {
	id := c.Param("id")
	var supplier models.Supplier
	
	// Check if supplier exists
	if err := database.DB.Where("id = ? AND is_active = ?", id, true).First(&supplier).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Supplier not found",
		})
		return
	}

	// Bind updated data
	var updateData models.Supplier
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid input data: " + err.Error(),
		})
		return
	}

	// Validate required fields
	if updateData.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Supplier name is required",
		})
		return
	}

	// Update fields
	supplier.Name = updateData.Name
	supplier.Email = updateData.Email
	supplier.Phone = updateData.Phone
	supplier.Address = updateData.Address
	supplier.ContactPerson = updateData.ContactPerson
	supplier.Website = updateData.Website

	if err := database.DB.Save(&supplier).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to update supplier: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Supplier updated successfully",
		"data":    supplier,
	})
}

// DeleteSupplier soft deletes a supplier
func DeleteSupplier(c *gin.Context) {
	id := c.Param("id")
	var supplier models.Supplier
	
	// Check if supplier exists
	if err := database.DB.Where("id = ? AND is_active = ?", id, true).First(&supplier).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Supplier not found",
		})
		return
	}

	// Soft delete by setting is_active to false
	supplier.IsActive = false
	if err := database.DB.Save(&supplier).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to delete supplier: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Supplier deleted successfully",
	})
}

// SearchSuppliers searches suppliers by name, email, or contact person
func SearchSuppliers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		GetSuppliers(c)
		return
	}

	var suppliers []models.Supplier
	searchPattern := "%" + query + "%"
	
	if err := database.DB.Where("is_active = ? AND (name ILIKE ? OR email ILIKE ? OR contact_person ILIKE ? OR phone ILIKE ?)", 
		true, searchPattern, searchPattern, searchPattern, searchPattern).Find(&suppliers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to search suppliers",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    suppliers,
	})
}