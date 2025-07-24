package handlers

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"inventory_system/database"
	"inventory_system/models"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
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
	ProductID  uint    `json:"product_id" binding:"required"`
	SupplierID *uint   `json:"supplier_id"` // Optional supplier selection
	Quantity   int     `json:"quantity" binding:"required,min=1"`
	Price      *float64 `json:"price"`      // Optional price override
	Cost       *float64 `json:"cost"`       // Optional cost override
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
		if err := tx.Preload("Suppliers.Supplier").First(&product, itemReq.ProductID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Product with ID %d not found", itemReq.ProductID)})
			return
		}

		var usePrice, useCost float64
		var supplierName string

		if itemReq.SupplierID != nil {
			// User selected a specific supplier
			fmt.Printf("Looking for supplier ID: %d\n", *itemReq.SupplierID)
			var selectedSupplier *models.ProductSupplier
			for _, supplier := range product.Suppliers {
				fmt.Printf("Checking supplier - ProductSupplier ID: %d, Company SupplierID: %d, IsActive: %t\n", supplier.ID, supplier.SupplierID, supplier.IsActive)
				// Check both ProductSupplier ID and actual supplier company ID
				if (supplier.SupplierID == *itemReq.SupplierID || supplier.ID == *itemReq.SupplierID) && supplier.IsActive {
					selectedSupplier = &supplier
					break
				}
			}

			if selectedSupplier == nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Selected supplier ID %d not found or inactive for product %s", *itemReq.SupplierID, product.Name)})
				return
			}

			// Check stock availability from selected supplier
			if selectedSupplier.Stock < itemReq.Quantity {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("Insufficient stock for product %s from selected supplier. Available: %d, Requested: %d",
						product.Name, selectedSupplier.Stock, itemReq.Quantity),
				})
				return
			}

			// Use supplier's price and cost (or override if provided)
			if itemReq.Price != nil {
				usePrice = *itemReq.Price
			} else {
				usePrice = selectedSupplier.Price
			}

			if itemReq.Cost != nil {
				useCost = *itemReq.Cost
			} else {
				useCost = selectedSupplier.Cost
			}

			supplierName = selectedSupplier.Supplier.Name

			// Update stock from selected supplier
			for i, supplier := range product.Suppliers {
				if supplier.ID == selectedSupplier.ID {
					product.Suppliers[i].Stock -= itemReq.Quantity
					if err := tx.Save(&product.Suppliers[i]).Error; err != nil {
						tx.Rollback()
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update supplier stock"})
						return
					}
					break
				}
			}
		} else {
			// Legacy mode - use lowest price supplier (FIFO approach)
			totalStock := product.GetTotalStock()
			usePrice = product.GetLowestPrice()
			useCost = product.GetLowestCost()

			// Check stock availability
			if totalStock < itemReq.Quantity {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("Insufficient stock for product %s. Available: %d, Requested: %d",
						product.Name, totalStock, itemReq.Quantity),
				})
				return
			}

			// Update stock from suppliers (FIFO approach - use cheapest supplier first)
			remainingQty := itemReq.Quantity
			for i, supplier := range product.Suppliers {
				if remainingQty <= 0 || !supplier.IsActive {
					continue
				}

				if supplier.Stock > 0 {
					deductQty := remainingQty
					if supplier.Stock < remainingQty {
						deductQty = supplier.Stock
					}

					product.Suppliers[i].Stock -= deductQty
					remainingQty -= deductQty

					if err := tx.Save(&product.Suppliers[i]).Error; err != nil {
						tx.Rollback()
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update supplier stock"})
						return
					}
				}
			}
		}

		// Create sale item
		itemTotal := float64(itemReq.Quantity) * usePrice
		saleItem := models.SaleItem{
			ProductID: product.ID,
			Quantity:  itemReq.Quantity,
			Price:     usePrice,
			Cost:      useCost,
			Total:     itemTotal,
		}

		saleItems = append(saleItems, saleItem)
		subtotal += itemTotal

		// Create stock movement record
		notes := "Sale transaction"
		if supplierName != "" {
			notes = fmt.Sprintf("Sale transaction - Supplier: %s", supplierName)
		}
		
		movement := models.StockMovement{
			ProductID: product.ID,
			UserID:    userID.(uint),
			Type:      "out",
			Quantity:  itemReq.Quantity,
			Reference: saleNumber,
			Notes:     notes,
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

	// Filter by payment status
	if paymentStatus := c.Query("payment_status"); paymentStatus != "" {
		filters["payment_status"] = paymentStatus
		query = query.Where("payment_status = ?", paymentStatus)
		countQuery = countQuery.Where("payment_status = ?", paymentStatus)
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

// generateSaleNumber generates a unique sale number with format A-0001, B-0001, etc.
func generateSaleNumber() string {
	// Get the last sale number from database
	var lastSale models.Sale
	result := database.DB.Order("id DESC").First(&lastSale)
	
	if result.Error != nil {
		// If no previous sales, start with A-0001
		return "A-0001"
	}
	
	// Parse the last sale number and increment it
	return incrementSaleNumber(lastSale.SaleNumber)
}

// incrementSaleNumber increments the sale number alphabetically
func incrementSaleNumber(lastNumber string) string {
	if lastNumber == "" {
		return "A-0001"
	}
	
	// Split the number (e.g., "A-0001" -> "A", "0001")
	parts := strings.Split(lastNumber, "-")
	if len(parts) != 2 {
		return "A-0001" // fallback if format is unexpected
	}
	
	prefix := parts[0]
	numStr := parts[1]
	
	// Parse the number part
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return "A-0001" // fallback if number parsing fails
	}
	
	// Increment logic
	if num < 9999 {
		// Increment number within same prefix
		num++
		return fmt.Sprintf("%s-%04d", prefix, num)
	}
	
	// Number reached 9999, need to increment prefix
	newPrefix := incrementPrefix(prefix)
	return fmt.Sprintf("%s-0001", newPrefix)
}

// incrementPrefix increments the alphabetic prefix
func incrementPrefix(prefix string) string {
	if prefix == "" {
		return "A"
	}
	
	runes := []rune(prefix)
	
	// Start from the last character and work backwards
	for i := len(runes) - 1; i >= 0; i-- {
		if runes[i] < 'Z' {
			// Can increment this character
			runes[i]++
			return string(runes)
		}
		// Character is Z, set to A and continue to next position
		runes[i] = 'A'
	}
	
	// All characters were Z, need to add a new character
	// Z -> AA, ZZ -> AAA, etc.
	return "A" + string(runes)
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

		// TODO: Restore quantity to specific suppliers (complex logic needed)
		// For now, stock restoration is disabled - need to implement supplier-specific restoration

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

		// TODO: Restore quantity to specific suppliers (complex logic needed)
		// For now, stock restoration is disabled - need to implement supplier-specific restoration

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

// GetSalesSummary returns summary statistics for sales
func GetSalesSummary(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// Set default dates if not provided (last 30 days)
	if startDate == "" {
		startDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	// Parse dates
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
		return
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
		return
	}
	// Add 23:59:59 to end date to include the entire day
	end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	db := database.GetDB()

	// Calculate total sales amount
	var totalSales float64
	db.Model(&models.Sale{}).
		Where("created_at BETWEEN ? AND ?", start, end).
		Select("COALESCE(SUM(total), 0)").
		Scan(&totalSales)

	// Calculate total number of transactions
	var totalTransactions int64
	db.Model(&models.Sale{}).
		Where("created_at BETWEEN ? AND ?", start, end).
		Count(&totalTransactions)

	// Calculate pending payments (credit sales with amount due > 0)
	var pendingPayments float64
	db.Model(&models.Sale{}).
		Where("created_at BETWEEN ? AND ? AND payment_method = ? AND amount_due > ?", start, end, "credit", 0).
		Select("COALESCE(SUM(amount_due), 0)").
		Scan(&pendingPayments)

	// Calculate overdue payments (credit sales past due date with amount due > 0)
	var overduePayments float64
	db.Model(&models.Sale{}).
		Where("created_at BETWEEN ? AND ? AND payment_method = ? AND amount_due > ? AND due_date < ?", 
			start, end, "credit", 0, time.Now()).
		Select("COALESCE(SUM(amount_due), 0)").
		Scan(&overduePayments)

	c.JSON(http.StatusOK, gin.H{
		"total_sales":       totalSales,
		"total_transactions": totalTransactions,
		"pending_payments":   pendingPayments,
		"overdue_payments":   overduePayments,
		"period": gin.H{
			"start_date": startDate,
			"end_date":   endDate,
		},
	})
}

// ExportSales exports sales data as CSV
func ExportSales(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	paymentMethod := c.Query("payment_method")
	paymentStatus := c.Query("payment_status")

	// Set default dates if not provided (last 30 days)
	if startDate == "" {
		startDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	// Parse dates
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
		return
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
		return
	}
	// Add 23:59:59 to end date to include the entire day
	end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	db := database.GetDB()

	// Build query with filters
	query := db.Preload("User").Preload("Items.Product").Where("created_at BETWEEN ? AND ?", start, end)

	if paymentMethod != "" {
		query = query.Where("payment_method = ?", paymentMethod)
	}

	if paymentStatus != "" {
		switch paymentStatus {
		case "paid":
			query = query.Where("payment_method != ? OR (payment_method = ? AND amount_due = ?)", "credit", "credit", 0)
		case "pending":
			query = query.Where("payment_method = ? AND amount_due > ? AND (due_date IS NULL OR due_date >= ?)", "credit", 0, time.Now())
		case "overdue":
			query = query.Where("payment_method = ? AND amount_due > ? AND due_date < ?", "credit", 0, time.Now())
		}
	}

	var sales []models.Sale
	if err := query.Find(&sales).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sales data"})
		return
	}

	// Create CSV content
	csvContent := "Sale Number,Customer Name,Date,Payment Method,Payment Status,Total,Amount Due,Cashier,Items\n"

	for _, sale := range sales {
		// Determine payment status
		paymentStatusStr := "paid"
		if sale.PaymentMethod == "credit" {
			if sale.AmountDue > 0 {
				if sale.DueDate != nil && sale.DueDate.Before(time.Now()) {
					paymentStatusStr = "overdue"
				} else {
					paymentStatusStr = "pending"
				}
			}
		}

		// Get items summary
		itemsSummary := ""
		if len(sale.Items) > 0 {
			for i, item := range sale.Items {
				if i > 0 {
					itemsSummary += "; "
				}
				productName := "Unknown Product"
				if item.Product.ID != 0 {
					productName = item.Product.Name
				}
				itemsSummary += fmt.Sprintf("%s (Qty: %d, Price: Rp%.2f)", productName, item.Quantity, item.Price)
			}
		}

		// Escape CSV fields that contain commas or quotes
		customerName := sale.CustomerName
		if customerName == "" {
			customerName = "Walk-in"
		}
		customerName = escapeCSV(customerName)
		
		cashierName := "Unknown"
		if sale.User.ID != 0 {
			cashierName = escapeCSV(sale.User.Name)
		}

		itemsSummary = escapeCSV(itemsSummary)

		csvContent += fmt.Sprintf("%s,%s,%s,%s,%s,%.2f,%.2f,%s,%s\n",
			escapeCSV(sale.SaleNumber),
			customerName,
			sale.CreatedAt.Format("2006-01-02 15:04:05"),
			escapeCSV(sale.PaymentMethod),
			paymentStatusStr,
			sale.Total,
			sale.AmountDue,
			cashierName,
			itemsSummary,
		)
	}

	// Set headers for CSV download
	filename := fmt.Sprintf("sales_report_%s.csv", time.Now().Format("2006-01-02"))
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Length", fmt.Sprintf("%d", len(csvContent)))

	c.String(http.StatusOK, csvContent)
}

// Helper function to escape CSV fields
func escapeCSV(field string) string {
	// If field contains comma, quote, or newline, wrap in quotes and escape internal quotes
	if fmt.Sprintf("%q", field) != "\""+field+"\"" {
		field = fmt.Sprintf("%q", field)
	}
	return field
}

// ExportSalesExcel exports sales data as Excel file
func ExportSalesExcel(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	paymentMethod := c.Query("payment_method")
	paymentStatus := c.Query("payment_status")

	// Set default dates if not provided (last 30 days)
	if startDate == "" {
		startDate = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	// Parse dates
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
		return
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
		return
	}
	// Add 23:59:59 to end date to include the entire day
	end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	db := database.GetDB()

	// Build query with filters
	query := db.Preload("User").Preload("Items.Product").Where("created_at BETWEEN ? AND ?", start, end)

	if paymentMethod != "" {
		query = query.Where("payment_method = ?", paymentMethod)
	}

	if paymentStatus != "" {
		switch paymentStatus {
		case "paid":
			query = query.Where("payment_method != ? OR (payment_method = ? AND amount_due = ?)", "credit", "credit", 0)
		case "pending":
			query = query.Where("payment_method = ? AND amount_due > ? AND (due_date IS NULL OR due_date >= ?)", "credit", 0, time.Now())
		case "overdue":
			query = query.Where("payment_method = ? AND amount_due > ? AND due_date < ?", "credit", 0, time.Now())
		}
	}

	var sales []models.Sale
	if err := query.Find(&sales).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sales data"})
		return
	}

	// Create new Excel file
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Create a new worksheet
	sheetName := "Sales Report"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Excel sheet"})
		return
	}

	// Set the worksheet as active
	f.SetActiveSheet(index)

	// Define headers
	headers := []string{
		"Sale Number", "Customer Name", "Date", "Payment Method", 
		"Payment Status", "Total", "Amount Due", "Cashier", "Items Count", "Items Detail",
	}

	// Create header style
	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 12,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#366092"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Excel style"})
		return
	}

	// Create data style
	dataStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{
			Vertical: "center",
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Excel data style"})
		return
	}

	// Set headers
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Add data rows
	for i, sale := range sales {
		row := i + 2 // Start from row 2 (after headers)
		
		// Determine payment status
		paymentStatusStr := "Paid"
		if sale.PaymentMethod == "credit" {
			if sale.AmountDue > 0 {
				if sale.DueDate != nil && sale.DueDate.Before(time.Now()) {
					paymentStatusStr = "Overdue"
				} else {
					paymentStatusStr = "Pending"
				}
			}
		}

		// Get items summary and count
		itemsCount := len(sale.Items)
		itemsSummary := ""
		if len(sale.Items) > 0 {
			for j, item := range sale.Items {
				if j > 0 {
					itemsSummary += "; "
				}
				productName := "Unknown Product"
				if item.Product.ID != 0 {
					productName = item.Product.Name
				}
				itemsSummary += fmt.Sprintf("%s (Qty: %d, Price: Rp%.2f)", productName, item.Quantity, item.Price)
			}
		}

		customerName := sale.CustomerName
		if customerName == "" {
			customerName = "Walk-in"
		}
		
		cashierName := "Unknown"
		if sale.User.ID != 0 {
			cashierName = sale.User.Name
		}

		// Set cell values
		data := []interface{}{
			sale.SaleNumber,
			customerName,
			sale.CreatedAt.Format("2006-01-02 15:04:05"),
			sale.PaymentMethod,
			paymentStatusStr,
			sale.Total,
			sale.AmountDue,
			cashierName,
			itemsCount,
			itemsSummary,
		}

		for j, value := range data {
			cell := fmt.Sprintf("%c%d", 'A'+j, row)
			f.SetCellValue(sheetName, cell, value)
			f.SetCellStyle(sheetName, cell, cell, dataStyle)
		}
	}

	// Auto-fit columns
	cols := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}
	for _, col := range cols {
		f.SetColWidth(sheetName, col, col, 15)
	}
	
	// Make the Items Detail column wider
	f.SetColWidth(sheetName, "J", "J", 50)

	// Add summary information at the top
	f.InsertRows(sheetName, 1, 3)
	
	// Add title
	f.SetCellValue(sheetName, "A1", "Sales Report")
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 16,
		},
	})
	f.SetCellStyle(sheetName, "A1", "A1", titleStyle)

	// Add date range
	f.SetCellValue(sheetName, "A2", fmt.Sprintf("Period: %s to %s", startDate, endDate))
	
	// Add total count
	f.SetCellValue(sheetName, "A3", fmt.Sprintf("Total Records: %d", len(sales)))

	// Set the header row (now row 4)
	for i, header := range headers {
		cell := fmt.Sprintf("%c4", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, headerStyle)
	}

	// Generate filename
	filename := fmt.Sprintf("sales_report_%s.xlsx", time.Now().Format("2006-01-02"))
	
	// Set headers for Excel download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	// Write the Excel file to the response
	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate Excel file"})
		return
	}
}
