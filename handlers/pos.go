package handlers

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"time"

	"inventory_system/database"
	"inventory_system/models"

	"github.com/gin-gonic/gin"
)

// SaleRequest represents a POS sale request
type SaleRequest struct {
	CustomerName  string            `json:"customer_name"`
	PaymentMethod string            `json:"payment_method" binding:"required"`
	PaymentTerm   string            `json:"payment_term"`
	DownPayment   float64           `json:"down_payment"`
	Items         []SaleItemRequest `json:"items" binding:"required,min=1"`
	Discount      float64           `json:"discount"`
	Tax           float64           `json:"tax"`
}

// SaleItemRequest represents an item in a sale
type SaleItemRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

// CreateSale processes a new sale transaction
func CreateSale(c *gin.Context) {
	var request SaleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate payment method
	validPaymentMethods := []string{"cash", "card", "transfer", "credit"}
	if !slices.Contains(validPaymentMethods, request.PaymentMethod) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment method"})
		return
	}

	// Set default payment term if not provided
	if request.PaymentTerm == "" {
		if request.PaymentMethod == "credit" {
			request.PaymentTerm = "net30"
		} else {
			request.PaymentTerm = "cash"
		}
	}

	// Validate payment term
	validPaymentTerms := []string{"cash", "net7", "net15", "net30", "net60", "net90"}
	if !slices.Contains(validPaymentTerms, request.PaymentTerm) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment term"})
		return
	}

	// Get user ID from context
	userID, _ := c.Get("user_id")

	// Start transaction
	tx := database.DB.Begin()

	// Generate sale number
	saleNumber := generateSaleNumber()

	// Calculate due date based on payment term
	var dueDate *time.Time
	var paymentStatus string

	if request.PaymentTerm == "cash" {
		paymentStatus = "paid"
		now := time.Now()
		dueDate = &now
	} else {
		paymentStatus = "pending"
		var days int
		switch request.PaymentTerm {
		case "net7":
			days = 7
		case "net15":
			days = 15
		case "net30":
			days = 30
		case "net60":
			days = 60
		case "net90":
			days = 90
		}
		calculatedDueDate := time.Now().AddDate(0, 0, days)
		dueDate = &calculatedDueDate
	}

	// Create sale record
	sale := models.Sale{
		SaleNumber:    saleNumber,
		UserID:        userID.(uint),
		CustomerName:  request.CustomerName,
		PaymentMethod: request.PaymentMethod,
		PaymentTerm:   request.PaymentTerm,
		PaymentStatus: paymentStatus,
		DownPayment:   request.DownPayment,
		DueDate:       dueDate,
		Status:        "completed",
		Discount:      request.Discount,
		Tax:           request.Tax,
	}

	var subtotal float64
	var saleItems []models.SaleItem

	// Process each item
	for _, itemReq := range request.Items {
		var product models.Product
		if err := tx.First(&product, itemReq.ProductID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Product with ID %d not found", itemReq.ProductID)})
			return
		}

		// Check stock availability
		if product.Quantity < itemReq.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Insufficient stock for product %s. Available: %d, Requested: %d",
					product.Name, product.Quantity, itemReq.Quantity),
			})
			return
		}

		// Update product quantity
		product.Quantity -= itemReq.Quantity
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product stock"})
			return
		}

		// Create sale item
		itemTotal := float64(itemReq.Quantity) * product.Price
		saleItem := models.SaleItem{
			ProductID: product.ID,
			Quantity:  itemReq.Quantity,
			Price:     product.Price,
			Cost:      product.Cost,
			Total:     itemTotal,
		}

		saleItems = append(saleItems, saleItem)
		subtotal += itemTotal

		// Create stock movement record
		movement := models.StockMovement{
			ProductID: product.ID,
			UserID:    userID.(uint),
			Type:      "out",
			Quantity:  itemReq.Quantity,
			Reference: saleNumber,
			Notes:     "Sale transaction",
		}

		if err := tx.Create(&movement).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record stock movement"})
			return
		}
	}

	// Calculate total
	sale.Subtotal = subtotal
	sale.Total = subtotal + request.Tax - request.Discount

	// Set payment amounts based on payment term
	if request.PaymentTerm == "cash" {
		sale.AmountPaid = sale.Total
		sale.AmountDue = 0
		now := time.Now()
		sale.PaidDate = &now
	} else {
		// For credit sales, set downpayment as amount paid
		sale.AmountPaid = request.DownPayment
		sale.AmountDue = sale.Total - request.DownPayment
		
		// If downpayment covers the full amount, mark as paid
		if sale.AmountDue <= 0 {
			sale.PaymentStatus = "paid"
			now := time.Now()
			sale.PaidDate = &now
			sale.AmountDue = 0
		}
	}

	// Save sale
	if err := tx.Create(&sale).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create sale"})
		return
	}

	// Save sale items
	for i := range saleItems {
		saleItems[i].SaleID = sale.ID
	}

	if err := tx.Create(&saleItems).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create sale items"})
		return
	}

	// Record downpayment if it's a credit sale with downpayment
	if request.PaymentMethod == "credit" && request.DownPayment > 0 {
		salePayment := models.SalePayment{
			SaleID:        sale.ID,
			UserID:        userID.(uint),
			Amount:        request.DownPayment,
			PaymentMethod: request.PaymentMethod,
			PaymentType:   "downpayment",
			Notes:         "Initial downpayment for credit sale",
		}
		
		if err := tx.Create(&salePayment).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record downpayment"})
			return
		}
	}

	tx.Commit()

	// Load complete sale data with items and products
	var completeSale models.Sale
	database.DB.Preload("Items.Product").Preload("User").First(&completeSale, sale.ID)

	c.JSON(http.StatusCreated, completeSale)
}

// GetSales returns all sales with optional filtering
func GetSales(c *gin.Context) {
	var sales []models.Sale
	
	// Build base query for counting (without preloads to avoid issues)
	countQuery := database.DB.Model(&models.Sale{})
	
	// Build main query with preloads
	query := database.DB.Preload("Items.Product").Preload("User")

	// Apply filters to both queries
	filters := make(map[string]any)
	
	// Filter by date range
	if startDate := c.Query("start_date"); startDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", startDate); err == nil {
			filters["start_date"] = parsedDate
			query = query.Where("created_at >= ?", parsedDate)
			countQuery = countQuery.Where("created_at >= ?", parsedDate)
		}
	}
	if endDate := c.Query("end_date"); endDate != "" {
		if parsedDate, err := time.Parse("2006-01-02", endDate); err == nil {
			endDateTime := parsedDate.Add(24 * time.Hour)
			filters["end_date"] = endDateTime
			query = query.Where("created_at <= ?", endDateTime)
			countQuery = countQuery.Where("created_at <= ?", endDateTime)
		}
	}

	// Filter by user
	if userID := c.Query("user_id"); userID != "" {
		filters["user_id"] = userID
		query = query.Where("user_id = ?", userID)
		countQuery = countQuery.Where("user_id = ?", userID)
	}

	// Filter by payment method
	if paymentMethod := c.Query("payment_method"); paymentMethod != "" {
		filters["payment_method"] = paymentMethod
		query = query.Where("payment_method = ?", paymentMethod)
		countQuery = countQuery.Where("payment_method = ?", paymentMethod)
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset := (page - 1) * limit

	// Get total count
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count sales"})
		return
	}

	// Get sales with preloads
	result := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&sales)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sales"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sales": sales,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetSale returns a specific sale
func GetSale(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sale ID"})
		return
	}

	var sale models.Sale
	result := database.DB.Preload("Items.Product").Preload("User").First(&sale, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sale not found"})
		return
	}

	c.JSON(http.StatusOK, sale)
}

// GetSalesReport returns sales analytics
func GetSalesReport(c *gin.Context) {
	// Default to today's date
	startDate := c.DefaultQuery("start_date", time.Now().Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	parsedStartDate, _ := time.Parse("2006-01-02", startDate)
	parsedEndDate, _ := time.Parse("2006-01-02", endDate)
	parsedEndDate = parsedEndDate.Add(24 * time.Hour) // Include the end date

	// Total sales count
	var totalSales int64
	database.DB.Model(&models.Sale{}).
		Where("created_at >= ? AND created_at < ?", parsedStartDate, parsedEndDate).
		Count(&totalSales)

	// Total revenue
	var totalRevenue float64
	database.DB.Model(&models.Sale{}).
		Where("created_at >= ? AND created_at < ?", parsedStartDate, parsedEndDate).
		Select("COALESCE(SUM(total), 0)").
		Scan(&totalRevenue)

	// Sales by payment method
	var paymentMethodStats []struct {
		PaymentMethod string  `json:"payment_method"`
		Count         int64   `json:"count"`
		Total         float64 `json:"total"`
	}
	database.DB.Model(&models.Sale{}).
		Select("payment_method, COUNT(*) as count, COALESCE(SUM(total), 0) as total").
		Where("created_at >= ? AND created_at < ?", parsedStartDate, parsedEndDate).
		Group("payment_method").
		Scan(&paymentMethodStats)

	// Top selling products
	var topProducts []struct {
		ProductID   uint    `json:"product_id"`
		ProductName string  `json:"product_name"`
		TotalSold   int64   `json:"total_sold"`
		Revenue     float64 `json:"revenue"`
	}
	database.DB.Table("sale_items si").
		Select("si.product_id, p.name as product_name, SUM(si.quantity) as total_sold, SUM(si.total) as revenue").
		Joins("JOIN products p ON si.product_id = p.id").
		Joins("JOIN sales s ON si.sale_id = s.id").
		Where("s.created_at >= ? AND s.created_at < ?", parsedStartDate, parsedEndDate).
		Group("si.product_id, p.name").
		Order("total_sold DESC").
		Limit(10).
		Scan(&topProducts)

	report := gin.H{
		"period": gin.H{
			"start_date": startDate,
			"end_date":   endDate,
		},
		"summary": gin.H{
			"total_sales":   totalSales,
			"total_revenue": totalRevenue,
		},
		"payment_methods": paymentMethodStats,
		"top_products":    topProducts,
	}

	c.JSON(http.StatusOK, report)
}

// generateSaleNumber generates a unique sale number
func generateSaleNumber() string {
	now := time.Now()
	return fmt.Sprintf("SALE-%s-%d", now.Format("20060102"), now.Unix())
}

// VoidSale cancels a sale (manager/admin only)
func VoidSale(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sale ID"})
		return
	}

	// Get user ID from context
	userID, _ := c.Get("user_id")

	// Start transaction
	tx := database.DB.Begin()

	var sale models.Sale
	if err := tx.Preload("Items.Product").First(&sale, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Sale not found"})
		return
	}

	if sale.Status == "cancelled" {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sale is already cancelled"})
		return
	}

	// Restore stock for each item (only if product still exists)
	for _, item := range sale.Items {
		var product models.Product
		if err := tx.First(&product, item.ProductID).Error; err != nil {
			// If product doesn't exist anymore, skip stock restoration but continue with void
			// This can happen if the product was deleted after the sale was made
			continue
		}

		// Restore quantity
		product.Quantity += item.Quantity
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore stock"})
			return
		}

		// Create stock movement record
		movement := models.StockMovement{
			ProductID: product.ID,
			UserID:    userID.(uint),
			Type:      "in",
			Quantity:  item.Quantity,
			Reference: sale.SaleNumber,
			Notes:     "Sale void - stock restored",
		}

		if err := tx.Create(&movement).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record stock movement"})
			return
		}
	}

	// Update sale status
	sale.Status = "cancelled"
	if err := tx.Save(&sale).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to void sale"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Sale voided successfully", "sale": sale})
}

// DeleteSale permanently deletes a sale (manager/admin only)
func DeleteSale(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sale ID"})
		return
	}

	// Get user ID from context
	userID, _ := c.Get("user_id")

	// Start transaction
	tx := database.DB.Begin()

	var sale models.Sale
	if err := tx.Preload("Items.Product").First(&sale, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Sale not found"})
		return
	}

	// Restore stock for each item before deletion (only if product still exists)
	for _, item := range sale.Items {
		var product models.Product
		if err := tx.First(&product, item.ProductID).Error; err != nil {
			// If product doesn't exist anymore, skip stock restoration but continue with deletion
			// This can happen if the product was deleted after the sale was made
			continue
		}

		// Restore quantity
		product.Quantity += item.Quantity
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore stock"})
			return
		}

		// Create stock movement record
		movement := models.StockMovement{
			ProductID: product.ID,
			UserID:    userID.(uint),
			Type:      "in",
			Quantity:  item.Quantity,
			Reference: sale.SaleNumber,
			Notes:     "Sale deleted - stock restored",
		}

		if err := tx.Create(&movement).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record stock movement"})
			return
		}
	}

	// Delete sale items first (foreign key constraint)
	if err := tx.Where("sale_id = ?", sale.ID).Delete(&models.SaleItem{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete sale items"})
		return
	}

	// Delete the sale
	if err := tx.Delete(&sale).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete sale"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Sale deleted successfully"})
}

// RecordSalePayment records a payment for a sale
func RecordSalePayment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sale ID"})
		return
	}

	var request struct {
		Amount        float64 `json:"amount" binding:"required,min=0"`
		PaymentMethod string  `json:"payment_method" binding:"required"`
		Notes         string  `json:"notes"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate payment method
	validPaymentMethods := []string{"cash", "card", "transfer", "credit"}
	if !slices.Contains(validPaymentMethods, request.PaymentMethod) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment method"})
		return
	}

	// Start transaction
	tx := database.DB.Begin()

	var sale models.Sale
	if err := tx.First(&sale, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Sale not found"})
		return
	}

	// Check if payment amount is valid
	if request.Amount > sale.AmountDue {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment amount exceeds amount due"})
		return
	}

	// Update payment amounts
	sale.AmountPaid += request.Amount
	sale.AmountDue -= request.Amount

	// Update payment status
	if sale.AmountDue <= 0.01 { // Allow for small rounding differences
		sale.PaymentStatus = "paid"
		now := time.Now()
		sale.PaidDate = &now
		sale.AmountDue = 0 // Ensure it's exactly 0
	} else {
		// Check if overdue
		if sale.DueDate != nil && time.Now().After(*sale.DueDate) {
			sale.PaymentStatus = "overdue"
		} else {
			sale.PaymentStatus = "pending"
		}
	}

	if err := tx.Save(&sale).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update sale payment"})
		return
	}

	// Record payment in SalePayment table
	userID, _ := c.Get("user_id")
	salePayment := models.SalePayment{
		SaleID:        sale.ID,
		UserID:        userID.(uint),
		Amount:        request.Amount,
		PaymentMethod: request.PaymentMethod,
		PaymentType:   "payment",
		Notes:         request.Notes,
	}

	if err := tx.Create(&salePayment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record payment"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment recorded successfully",
		"sale":    sale,
	})
}

// GetOverdueSales returns sales that are overdue for payment
func GetOverdueSales(c *gin.Context) {
	var sales []models.Sale

	query := database.DB.Preload("Items.Product").Preload("User").
		Where("payment_status = ? OR (payment_status = ? AND due_date < ?)", 
			"overdue", "pending", time.Now())

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset := (page - 1) * limit

	var total int64
	database.DB.Model(&models.Sale{}).
		Where("payment_status = ? OR (payment_status = ? AND due_date < ?)", 
			"overdue", "pending", time.Now()).
		Count(&total)

	result := query.Order("due_date ASC").Offset(offset).Limit(limit).Find(&sales)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch overdue sales"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sales": sales,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// GetSalePayments returns payment history for a sale
func GetSalePayments(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sale ID"})
		return
	}

	var payments []models.SalePayment
	result := database.DB.Where("sale_id = ?", id).
		Preload("User").
		Order("created_at DESC").
		Find(&payments)
	
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payment history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payments": payments,
		"total":    len(payments),
	})
}
