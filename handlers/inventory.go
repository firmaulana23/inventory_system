package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"inventory_system/database"
	"inventory_system/models"

	"github.com/gin-gonic/gin"
)

// CreateProduct creates a new product
func CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if SKU already exists
	var existingProduct models.Product
	result := database.DB.Where("sku = ?", product.SKU).First(&existingProduct)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Product with this SKU already exists"})
		return
	}

	result = database.DB.Create(&product)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, product)
}

// GetProducts returns all products with optional filtering
func GetProducts(c *gin.Context) {
	var products []models.Product
	query := database.DB.Model(&models.Product{})

	// Add filters
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}
	if search := c.Query("search"); search != "" {
		query = query.Where("name LIKE ? OR sku LIKE ? OR description LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if lowStock := c.Query("low_stock"); lowStock == "true" {
		query = query.Where("quantity <= min_stock")
	}
	if active := c.Query("active"); active != "" {
		isActive := active == "true"
		query = query.Where("is_active = ?", isActive)
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset := (page - 1) * limit

	var total int64
	query.Count(&total)

	result := query.Offset(offset).Limit(limit).Find(&products)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"total":    total,
		"page":     page,
		"limit":    limit,
	})
}

// GetProduct returns a specific product
func GetProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	result := database.DB.First(&product, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// UpdateProduct updates product information
func UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var product models.Product
	result := database.DB.First(&product, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result = database.DB.Save(&product)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// DeleteProduct soft deletes a product
func DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	result := database.DB.Delete(&models.Product{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// AdjustStock adjusts product stock
func AdjustStock(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var request struct {
		Quantity int    `json:"quantity" binding:"required"`
		Type     string `json:"type" binding:"required"` // in, out, adjustment
		Notes    string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate adjustment type
	if request.Type != "in" && request.Type != "out" && request.Type != "adjustment" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid adjustment type"})
		return
	}

	var product models.Product
	result := database.DB.First(&product, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Get user ID from context
	userID, _ := c.Get("user_id")

	// Start transaction
	tx := database.DB.Begin()

	// Calculate new quantity
	var newQuantity int
	switch request.Type {
	case "in":
		newQuantity = product.Quantity + request.Quantity
	case "out":
		newQuantity = product.Quantity - request.Quantity
		if newQuantity < 0 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
			return
		}
	case "adjustment":
		newQuantity = request.Quantity
	}

	// Update product quantity
	product.Quantity = newQuantity
	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stock"})
		return
	}

	// Create stock movement record
	movement := models.StockMovement{
		ProductID: product.ID,
		UserID:    userID.(uint),
		Type:      request.Type,
		Quantity:  request.Quantity,
		Notes:     request.Notes,
	}

	if err := tx.Create(&movement).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record stock movement"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Stock adjusted successfully",
		"product": product,
	})
}

// GetStockMovements returns stock movement history
func GetStockMovements(c *gin.Context) {
	var movements []models.StockMovement
	query := database.DB.Preload("Product").Preload("User")

	// Filter by product if specified
	if productID := c.Query("product_id"); productID != "" {
		query = query.Where("product_id = ?", productID)
	}

	// Filter by type if specified
	if movementType := c.Query("type"); movementType != "" {
		query = query.Where("type = ?", movementType)
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset := (page - 1) * limit

	var total int64
	query.Count(&total)

	result := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&movements)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stock movements"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movements": movements,
		"total":     total,
		"page":      page,
		"limit":     limit,
	})
}

// GetLowStockProducts returns products with low stock
func GetLowStockProducts(c *gin.Context) {
	var products []models.Product
	result := database.DB.Where("quantity <= min_stock AND is_active = ?", true).Find(&products)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch low stock products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetProductCategories returns unique product categories
func GetProductCategories(c *gin.Context) {
	var categories []string
	result := database.DB.Model(&models.Product{}).
		Distinct("category").
		Where("category != '' AND is_active = ?", true).
		Pluck("category", &categories)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// SearchProducts searches products by name, SKU, or barcode
func SearchProducts(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	var products []models.Product
	searchPattern := "%" + strings.ToLower(query) + "%"

	result := database.DB.Where(
		"(LOWER(name) LIKE ? OR LOWER(sku) LIKE ? OR LOWER(description) LIKE ?) AND is_active = ?",
		searchPattern, searchPattern, searchPattern, true,
	).Limit(20).Find(&products)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search products"})
		return
	}

	c.JSON(http.StatusOK, products)
}
