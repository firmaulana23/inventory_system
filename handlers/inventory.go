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
	var request struct {
		Name        string `json:"name" binding:"required"`
		SKU         string `json:"sku" binding:"required"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Location    string `json:"location"`
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Start transaction
	tx := database.DB.Begin()

	// Check if SKU already exists
	var existingProduct models.Product
	result := tx.Where("sku = ?", request.SKU).First(&existingProduct)
	if result.Error == nil {
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{"error": "Product with this SKU already exists"})
		return
	}

	product := models.Product{
		Name:        request.Name,
		SKU:         request.SKU,
		Description: request.Description,
		Category:    request.Category,
		Location:    request.Location,
		IsActive:    true,
	}

	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		// Log the detailed error for debugging
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create product",
			"details": err.Error(),
		})
		return
	}

	// Commit transaction
	tx.Commit()

	c.JSON(http.StatusCreated, product)
}

// GetProducts returns all products with optional filtering
func GetProducts(c *gin.Context) {
	var products []models.Product
	query := database.DB.Model(&models.Product{}).Preload("Suppliers.Supplier")

	// Add filters
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}
	if search := c.Query("search"); search != "" {
		query = query.Where("name LIKE ? OR sku LIKE ? OR description LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	// Note: Low stock filtering now requires loading suppliers relationship
	// This will be handled after loading the products
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

	var request struct {
		Name        string `json:"name" binding:"required"`
		SKU         string `json:"sku" binding:"required"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Location    string `json:"location"`
		IsActive    *bool  `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	product.Name = request.Name
	product.SKU = request.SKU
	product.Description = request.Description
	product.Category = request.Category
	product.Location = request.Location
	
	if request.IsActive != nil {
		product.IsActive = *request.IsActive
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

// AdjustSupplierStock adjusts stock for a specific supplier
func AdjustSupplierStock(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	supplierID, err := strconv.ParseUint(c.Param("supplier_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid supplier ID"})
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

	var productSupplier models.ProductSupplier
	result := database.DB.Where("product_id = ? AND supplier_id = ?", productID, supplierID).First(&productSupplier)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product-supplier relationship not found"})
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
		newQuantity = productSupplier.Stock + request.Quantity
	case "out":
		newQuantity = productSupplier.Stock - request.Quantity
		if newQuantity < 0 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock for this supplier"})
			return
		}
	case "adjustment":
		newQuantity = request.Quantity
	}

	// Update supplier stock
	productSupplier.Stock = newQuantity
	if err := tx.Save(&productSupplier).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stock"})
		return
	}

	// Create stock movement record
	movement := models.StockMovement{
		ProductID: uint(productID),
		UserID:    userID.(uint),
		Type:      request.Type,
		Quantity:  request.Quantity,
		Reference: "Supplier ID: " + strconv.FormatUint(supplierID, 10),
		Notes:     request.Notes,
	}

	if err := tx.Create(&movement).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record stock movement"})
		return
	}

	tx.Commit()

	// Load with supplier info
	database.DB.Preload("Supplier").First(&productSupplier, productSupplier.ID)

	c.JSON(http.StatusOK, gin.H{
		"message":         "Stock adjusted successfully",
		"product_supplier": productSupplier,
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

// GetLowStockProducts returns products with low stock from any supplier
func GetLowStockProducts(c *gin.Context) {
	var products []models.Product
	result := database.DB.Preload("Suppliers.Supplier").Where("is_active = ?", true).Find(&products)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	// Filter products that have low stock from any supplier
	var lowStockProducts []models.Product
	for _, product := range products {
		if product.IsLowStock() {
			lowStockProducts = append(lowStockProducts, product)
		}
	}

	c.JSON(http.StatusOK, lowStockProducts)
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

	result := database.DB.Preload("Suppliers.Supplier").Where(
		"(LOWER(name) LIKE ? OR LOWER(sku) LIKE ? OR LOWER(description) LIKE ?) AND is_active = ?",
		searchPattern, searchPattern, searchPattern, true,
	).Limit(20).Find(&products)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// AddProductSupplier adds a supplier to a product with pricing
func AddProductSupplier(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var request struct {
		SupplierID uint    `json:"supplier_id" binding:"required"`
		Cost       float64 `json:"cost" binding:"required"`
		Price      float64 `json:"price" binding:"required"`
		Stock      int     `json:"stock"`
		MinStock   int     `json:"min_stock"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if product exists
	var product models.Product
	if err := database.DB.First(&product, productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Check if supplier exists
	var supplier models.Supplier
	if err := database.DB.First(&supplier, request.SupplierID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Supplier not found"})
		return
	}

	// Check if product-supplier relationship already exists
	var existing models.ProductSupplier
	result := database.DB.Where("product_id = ? AND supplier_id = ?", productID, request.SupplierID).First(&existing)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Supplier already linked to this product"})
		return
	}

	// Create new product-supplier relationship
	productSupplier := models.ProductSupplier{
		ProductID:  uint(productID),
		SupplierID: request.SupplierID,
		Cost:       request.Cost,
		Price:      request.Price,
		Stock:      request.Stock,
		MinStock:   request.MinStock,
	}

	if err := database.DB.Create(&productSupplier).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add supplier to product"})
		return
	}

	// Load the relationship with supplier info
	database.DB.Preload("Supplier").First(&productSupplier, productSupplier.ID)

	c.JSON(http.StatusCreated, productSupplier)
}

// UpdateProductSupplier updates product-supplier relationship
func UpdateProductSupplier(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	supplierID, err := strconv.ParseUint(c.Param("supplier_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid supplier ID"})
		return
	}

	var request struct {
		Cost     float64 `json:"cost" binding:"required"`
		Price    float64 `json:"price" binding:"required"`
		Stock    int     `json:"stock"`
		MinStock int     `json:"min_stock"`
		IsActive bool    `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var productSupplier models.ProductSupplier
	if err := database.DB.Where("product_id = ? AND supplier_id = ?", productID, supplierID).First(&productSupplier).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product-supplier relationship not found"})
		return
	}

	// Update fields
	productSupplier.Cost = request.Cost
	productSupplier.Price = request.Price
	productSupplier.Stock = request.Stock
	productSupplier.MinStock = request.MinStock
	productSupplier.IsActive = request.IsActive

	if err := database.DB.Save(&productSupplier).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product-supplier relationship"})
		return
	}

	// Load with supplier info
	database.DB.Preload("Supplier").First(&productSupplier, productSupplier.ID)

	c.JSON(http.StatusOK, productSupplier)
}

// RemoveProductSupplier removes a supplier from a product
func RemoveProductSupplier(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	supplierID, err := strconv.ParseUint(c.Param("supplier_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid supplier ID"})
		return
	}

	result := database.DB.Where("product_id = ? AND supplier_id = ?", productID, supplierID).Delete(&models.ProductSupplier{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove supplier from product"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product-supplier relationship not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Supplier removed from product successfully"})
}

// GetProductSuppliers returns all suppliers for a product
func GetProductSuppliers(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var productSuppliers []models.ProductSupplier
	result := database.DB.Preload("Supplier").Where("product_id = ? AND is_active = ?", productID, true).Find(&productSuppliers)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product suppliers"})
		return
	}

	c.JSON(http.StatusOK, productSuppliers)
}
