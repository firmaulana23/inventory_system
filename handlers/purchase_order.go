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
	"gorm.io/gorm"
)

// CreatePurchaseOrderRequest represents the request body for creating a purchase order
type CreatePurchaseOrderRequest struct {
	SupplierID    uint                      `json:"supplier_id" binding:"required"`
	PaymentMethod string                    `json:"payment_method"`
	PaymentDays   int                       `json:"payment_days"` // Number of days for payment due
	DownPayment   float64                   `json:"down_payment"`
	Notes         string                    `json:"notes"`
	OrderDate     string                    `json:"order_date" binding:"required"`
	Items         []CreatePurchaseOrderItem `json:"items" binding:"required,min=1"`
}

// CreatePurchaseOrderItem represents an item in the purchase order request
type CreatePurchaseOrderItem struct {
	SKU               string  `json:"sku" binding:"required"`
	ProductName       string  `json:"product_name"`
	Category          string  `json:"category"`
	Description       string  `json:"description"`
	Quantity          int     `json:"quantity" binding:"required,min=1"`
	UnitCost          float64 `json:"unit_cost" binding:"required,min=0"`
	ProductSupplierID *uint   `json:"product_supplier_id"` // Link to specific supplier for existing products
}

// CreatePurchaseOrder creates a new purchase order with SKU-based product handling
func CreatePurchaseOrder(c *gin.Context) {
	var req CreatePurchaseOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Failed to bind JSON: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Received purchase order request: %+v\n", req)

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Set default payment method and days if not provided
	if req.PaymentMethod == "" {
		req.PaymentMethod = "cash"
	}
	if req.PaymentDays == 0 {
		if req.PaymentMethod == "credit" {
			req.PaymentDays = 30
		}
		// For non-credit payments, PaymentDays remains 0 (immediate payment)
	}

	// Validate payment method
	validPaymentMethods := []string{"cash", "transfer", "credit", "qris"}
	if !slices.Contains(validPaymentMethods, req.PaymentMethod) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment method"})
		return
	}

	// Validate payment days
	if req.PaymentDays < 0 || req.PaymentDays > 365 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment days must be between 0 and 365"})
		return
	}

	// Parse dates
	orderDate, err := time.Parse("2006-01-02", req.OrderDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order date format. Use YYYY-MM-DD"})
		return
	}

	// Generate PO number
	poNumber := generatePONumber()

	// Calculate due date based on payment days
	var dueDate *time.Time
	var paymentStatus string
	var paidDate *time.Time

	if req.PaymentMethod == "credit" && req.PaymentDays >= 0 {
		paymentStatus = "pending"
		calculatedDueDate := orderDate.AddDate(0, 0, req.PaymentDays)
		dueDate = &calculatedDueDate
	} else {
		paymentStatus = "paid" // Will be paid upon receipt
		dueDate = &orderDate
		paidDate = &orderDate // Set paid date to order date for cash/transfer
	}

	// Start transaction
	tx := database.DB.Begin()

	// Create purchase order
	po := models.PurchaseOrder{
		PONumber:      poNumber,
		SupplierID:    req.SupplierID,
		UserID:        userID.(uint),
		PaymentMethod: req.PaymentMethod,
		PaymentDays:   req.PaymentDays,
		PaymentStatus: paymentStatus,
		PaidDate:      paidDate,
		DueDate:       dueDate,
		Notes:         req.Notes,
		OrderDate:     orderDate,
	}

	if err := tx.Create(&po).Error; err != nil {
		tx.Rollback()
		fmt.Printf("Failed to create purchase order: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create purchase order: %v", err)})
		return
	}

	// Load supplier to get supplier name
	var supplier models.Supplier
	if err := tx.First(&supplier, po.SupplierID).Error; err != nil {
		tx.Rollback()
		fmt.Printf("Failed to find supplier with ID %d: %v\n", po.SupplierID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid supplier ID %d: %v", po.SupplierID, err)})
		return
	}

	var totalAmount float64

	// Process each item
	for i, item := range req.Items {
		fmt.Printf("Processing item %d: %+v\n", i+1, item)

		// Find or create product by SKU
		product, productSupplier, err := findOrCreateProductWithSupplier(tx, item, req.SupplierID)
		if err != nil {
			tx.Rollback()
			fmt.Printf("Failed to process product with SKU %s: %v\n", item.SKU, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to process product with SKU %s: %v", item.SKU, err)})
			return
		}

		// Update supplier-specific stock
		if productSupplier != nil {
			productSupplier.Stock += item.Quantity
			if err := tx.Save(productSupplier).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update supplier stock"})
				return
			}
		}

		// Create stock movement record
		stockMovement := models.StockMovement{
			ProductID: product.ID,
			UserID:    userID.(uint),
			Type:      "in",
			Quantity:  item.Quantity,
			Reference: po.PONumber,
			Notes:     fmt.Sprintf("Purchase Order %s - %s", po.PONumber, supplier.Name),
		}
		if err := tx.Create(&stockMovement).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create stock movement record"})
			return
		}

		// Create purchase order item
		total := float64(item.Quantity) * item.UnitCost
		poItem := models.PurchaseOrderItem{
			PurchaseOrderID:  po.ID,
			ProductID:        product.ID,
			QuantityOrdered:  item.Quantity,
			QuantityReceived: item.Quantity, // Mark as received since we're updating inventory
			UnitCost:         item.UnitCost,
			Total:            total,
		}

		// Link to specific product-supplier relationship if available
		if productSupplier != nil {
			poItem.ProductSupplierID = &productSupplier.ID
		}

		if err := tx.Create(&poItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create purchase order item"})
			return
		}

		totalAmount += total
	}

	// Update purchase order total and amount due
	po.Total = totalAmount
	po.DownPayment = req.DownPayment

	// Calculate amounts based on payment method
	if req.PaymentMethod == "credit" {
		// For credit payments, handle downpayment
		if req.DownPayment > 0 {
			// Validate downpayment doesn't exceed total
			if req.DownPayment > totalAmount {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "Downpayment cannot exceed total amount"})
				return
			}
			po.AmountPaid = req.DownPayment
			po.AmountDue = totalAmount - req.DownPayment

			// If downpayment covers the full amount, mark as paid
			if po.AmountDue <= 0.01 {
				po.PaymentStatus = "paid"
				now := time.Now()
				po.PaidDate = &now
			}
		} else {
			// Credit with no downpayment - nothing paid yet
			po.AmountPaid = 0
			po.AmountDue = totalAmount
		}
	} else {
		// For cash or transfer, mark as paid immediately
		po.AmountDue = 0
		po.AmountPaid = totalAmount
	}

	if err := tx.Save(&po).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update purchase order total"})
		return
	}

	// Record downpayment in payment history if applicable
	if req.PaymentMethod == "credit" && req.DownPayment > 0 {
		payment := models.PurchasePayment{
			PurchaseOrderID: po.ID,
			UserID:          userID.(uint),
			Amount:          req.DownPayment,
			PaymentMethod:   req.PaymentMethod,
			PaymentType:     "downpayment",
			Notes:           fmt.Sprintf("Down payment for PO %s", po.PONumber),
		}

		if err := tx.Create(&payment).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record downpayment history"})
			return
		}
	}

	// Commit transaction
	tx.Commit()

	// Load complete purchase order with relationships
	var completePO models.PurchaseOrder
	database.DB.Preload("User").Preload("Supplier").Preload("Items.Product").Preload("Items.ProductSupplier.Supplier").First(&completePO, po.ID)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    completePO,
	})
}

// findOrCreateProductWithSupplier finds an existing product by SKU or creates a new one, and handles supplier relationship
func findOrCreateProductWithSupplier(tx *gorm.DB, item CreatePurchaseOrderItem, supplierID uint) (*models.Product, *models.ProductSupplier, error) {
	var product models.Product

	// Try to find existing product by SKU
	result := tx.Where("sku = ?", item.SKU).First(&product)
	if result.Error != nil {
		// Product doesn't exist, create new one
		product = models.Product{
			Name:        getProductName(item),
			SKU:         item.SKU,
			Description: item.Description,
			Category:    item.Category,
			IsActive:    true,
		}

		if err := tx.Create(&product).Error; err != nil {
			return nil, nil, err
		}
	}

	// Handle product-supplier relationship
	var productSupplier models.ProductSupplier

	// If ProductSupplierID is provided, use that specific relationship
	if item.ProductSupplierID != nil {
		if err := tx.Where("id = ? AND product_id = ? AND supplier_id = ?",
			*item.ProductSupplierID, product.ID, supplierID).First(&productSupplier).Error; err != nil {
			return nil, nil, fmt.Errorf("invalid product_supplier_id: %v", err)
		}
		return &product, &productSupplier, nil
	}

	// Check if product-supplier relationship exists
	result = tx.Where("product_id = ? AND supplier_id = ?", product.ID, supplierID).First(&productSupplier)
	if result.Error != nil {
		// Create new product-supplier relationship
		productSupplier = models.ProductSupplier{
			ProductID:  product.ID,
			SupplierID: supplierID,
			Cost:       item.UnitCost,
			Price:      item.UnitCost * 1.2, // Default markup of 20%
			Stock:      0,                   // Will be updated by caller
			MinStock:   10,                  // Default minimum stock
			IsActive:   true,
		}

		if err := tx.Create(&productSupplier).Error; err != nil {
			return nil, nil, err
		}
	} else {
		// Update cost if different (in case prices changed)
		if productSupplier.Cost != item.UnitCost {
			productSupplier.Cost = item.UnitCost
		}
	}

	return &product, &productSupplier, nil
}

// getProductName returns the product name, using SKU as fallback
func getProductName(item CreatePurchaseOrderItem) string {
	if item.ProductName != "" {
		return item.ProductName
	}
	return item.SKU // Use SKU as name if no name provided
}

// generatePONumber generates a unique purchase order number
func generatePONumber() string {
	now := time.Now()
	return fmt.Sprintf("PO-%d%02d%02d-%d", now.Year(), now.Month(), now.Day(), now.Unix()%10000)
}

// GetPurchaseOrders returns all purchase orders
func GetPurchaseOrders(c *gin.Context) {
	var purchaseOrders []models.PurchaseOrder
	query := database.DB.Model(&models.PurchaseOrder{}).Preload("User").Preload("Supplier")

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset := (page - 1) * limit

	var total int64
	query.Count(&total)

	result := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&purchaseOrders)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch purchase orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"purchase_orders": purchaseOrders,
		"total":           total,
		"page":            page,
		"limit":           limit,
	})
}

// GetPurchaseOrder returns a specific purchase order
func GetPurchaseOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid purchase order ID"})
		return
	}

	var po models.PurchaseOrder
	result := database.DB.Preload("User").Preload("Supplier").Preload("Items.Product").Preload("Items.ProductSupplier.Supplier").First(&po, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Purchase order not found"})
		return
	}

	c.JSON(http.StatusOK, po)
}

// RecordPurchasePayment records a payment for a purchase order
func RecordPurchasePayment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid purchase order ID"})
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
	validPaymentMethods := []string{"cash", "card", "transfer", "check", "wire"}
	if !slices.Contains(validPaymentMethods, request.PaymentMethod) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment method"})
		return
	}

	// Start transaction
	tx := database.DB.Begin()

	var po models.PurchaseOrder
	if err := tx.First(&po, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Purchase order not found"})
		return
	}

	// Check if payment amount is valid
	if request.Amount > po.AmountDue {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment amount exceeds amount due"})
		return
	}

	// Update payment amounts
	po.AmountPaid += request.Amount
	po.AmountDue -= request.Amount

	// Update payment status
	if po.AmountDue <= 0.01 { // Allow for small rounding differences
		po.PaymentStatus = "paid"
		now := time.Now()
		po.PaidDate = &now
		po.AmountDue = 0 // Ensure it's exactly 0
	} else {
		// Check if overdue
		if po.DueDate != nil && time.Now().After(*po.DueDate) {
			po.PaymentStatus = "overdue"
		} else {
			po.PaymentStatus = "pending"
		}
	}

	if err := tx.Save(&po).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update purchase order payment"})
		return
	}

	// Get user ID from context
	userID, _ := c.Get("user_id")

	// Record payment in payment history
	payment := models.PurchasePayment{
		PurchaseOrderID: po.ID,
		UserID:          userID.(uint),
		Amount:          request.Amount,
		PaymentMethod:   request.PaymentMethod,
		PaymentType:     "payment",
		Notes:           request.Notes,
	}

	if err := tx.Create(&payment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record payment history"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"message":        "Payment recorded successfully",
		"purchase_order": po,
	})
}

// GetOverduePurchaseOrders returns purchase orders that are overdue for payment
func GetOverduePurchaseOrders(c *gin.Context) {
	var purchaseOrders []models.PurchaseOrder

	query := database.DB.Preload("User").Preload("Supplier").Preload("Items.Product").Preload("Items.ProductSupplier.Supplier").
		Where("payment_status = ? OR (payment_status = ? AND due_date < ?)",
			"overdue", "pending", time.Now())

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset := (page - 1) * limit

	var total int64
	database.DB.Model(&models.PurchaseOrder{}).
		Where("payment_status = ? OR (payment_status = ? AND due_date < ?)",
			"overdue", "pending", time.Now()).
		Count(&total)

	result := query.Order("due_date ASC").Offset(offset).Limit(limit).Find(&purchaseOrders)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch overdue purchase orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"purchase_orders": purchaseOrders,
		"total":           total,
		"page":            page,
		"limit":           limit,
	})
}

// DeletePurchaseOrder deletes a purchase order and reverses stock movements
func DeletePurchaseOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid purchase order ID"})
		return
	}

	// Get user ID from context
	userID, _ := c.Get("user_id")

	// Start transaction
	tx := database.DB.Begin()

	var po models.PurchaseOrder
	if err := tx.Preload("Items.Product").Preload("Items.ProductSupplier").First(&po, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Purchase order not found"})
		return
	}

	// Reverse stock movements for each item (only if status is not cancelled)
	if po.PaymentStatus != "cancelled" {
		for _, item := range po.Items {
			var product models.Product
			if err := tx.First(&product, item.ProductID).Error; err != nil {
				// If product doesn't exist anymore, skip stock restoration but continue with deletion
				continue
			}

			// Reverse supplier-specific stock if ProductSupplierID is linked
			if item.ProductSupplierID != nil {
				var productSupplier models.ProductSupplier
				if err := tx.First(&productSupplier, *item.ProductSupplierID).Error; err == nil {
					// Reduce supplier stock by the received quantity
					newStock := productSupplier.Stock - item.QuantityReceived
					if newStock < 0 {
						newStock = 0 // Don't allow negative stock
					}
					productSupplier.Stock = newStock

					if err := tx.Save(&productSupplier).Error; err != nil {
						tx.Rollback()
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reverse supplier stock"})
						return
					}
				}
			}

			// Create a reversing stock movement record
			movement := models.StockMovement{
				ProductID: product.ID,
				UserID:    userID.(uint),
				Type:      "out",
				Quantity:  item.QuantityReceived,
				Reference: po.PONumber,
				Notes:     "Purchase order deleted - stock reversed",
			}

			if err := tx.Create(&movement).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record stock movement"})
				return
			}
		}
	}

	// Delete purchase order items first (foreign key constraint)
	if err := tx.Where("purchase_order_id = ?", po.ID).Delete(&models.PurchaseOrderItem{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete purchase order items"})
		return
	}

	// Delete the purchase order
	if err := tx.Delete(&po).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete purchase order"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Purchase order deleted successfully"})
}

// GetPurchasePaymentHistory returns payment history for a specific purchase order
func GetPurchasePaymentHistory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid purchase order ID"})
		return
	}

	// Get purchase order basic info
	var po models.PurchaseOrder
	if err := database.DB.First(&po, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Purchase order not found"})
		return
	}

	// Get payment history
	var payments []models.PurchasePayment
	result := database.DB.Preload("User").
		Where("purchase_order_id = ?", id).
		Order("created_at DESC").
		Find(&payments)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payment history"})
		return
	}

	// Calculate summary
	var totalPaid float64
	for _, payment := range payments {
		totalPaid += payment.Amount
	}

	c.JSON(http.StatusOK, gin.H{
		"purchase_order": gin.H{
			"id":             po.ID,
			"po_number":      po.PONumber,
			"supplier":       po.Supplier,
			"total":          po.Total,
			"amount_paid":    po.AmountPaid,
			"amount_due":     po.AmountDue,
			"payment_status": po.PaymentStatus,
			"due_date":       po.DueDate,
			"paid_date":      po.PaidDate,
		},
		"payment_history": payments,
		"summary": gin.H{
			"total_payments": len(payments),
			"total_paid":     totalPaid,
			"remaining_due":  po.AmountDue,
		},
	})
}

// UpdatePurchaseOrderRequest represents the request body for updating a purchase order
type UpdatePurchaseOrderRequest struct {
	SupplierID    uint    `json:"supplier_id"`
	PaymentMethod string  `json:"payment_method"`
	PaymentDays   int     `json:"payment_days"` // Number of days for payment due
	Notes         string  `json:"notes"`
	Status        string  `json:"status"`
	DownPayment   float64 `json:"down_payment"`
}

// GetPurchaseOrdersSummary returns summary statistics for purchase orders
func GetPurchaseOrdersSummary(c *gin.Context) {
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

	db := database.DB

	// Calculate total purchase orders count
	var totalOrders int64
	db.Model(&models.PurchaseOrder{}).
		Where("created_at BETWEEN ? AND ?", start, end).
		Count(&totalOrders)

	// Calculate total amount
	var totalAmount float64
	db.Model(&models.PurchaseOrder{}).
		Where("created_at BETWEEN ? AND ?", start, end).
		Select("COALESCE(SUM(total), 0)").
		Scan(&totalAmount)

	// Calculate pending payments (credit orders with amount due > 0)
	var pendingAmount float64
	db.Model(&models.PurchaseOrder{}).
		Where("created_at BETWEEN ? AND ? AND payment_method = ? AND amount_due > ?", start, end, "credit", 0).
		Select("COALESCE(SUM(amount_due), 0)").
		Scan(&pendingAmount)

	// Calculate overdue payments (credit orders past due date with amount due > 0)
	var overdueAmount float64
	db.Model(&models.PurchaseOrder{}).
		Where("created_at BETWEEN ? AND ? AND payment_method = ? AND amount_due > ? AND due_date < ?",
			start, end, "credit", 0, time.Now()).
		Select("COALESCE(SUM(amount_due), 0)").
		Scan(&overdueAmount)

	c.JSON(http.StatusOK, gin.H{
		"total_orders":   totalOrders,
		"total_amount":   totalAmount,
		"pending_amount": pendingAmount,
		"overdue_amount": overdueAmount,
		"period": gin.H{
			"start_date": startDate,
			"end_date":   endDate,
		},
	})
}

// UpdatePurchaseOrder updates an existing purchase order
func UpdatePurchaseOrder(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid purchase order ID"})
		return
	}

	var req UpdatePurchaseOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Start transaction
	tx := database.DB.Begin()

	// Get existing purchase order
	var po models.PurchaseOrder
	if err := tx.First(&po, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Purchase order not found"})
		return
	}

	// Store original values for comparison
	originalDownPayment := po.DownPayment
	originalPaymentDays := po.PaymentDays

	// Validate payment method if provided
	if req.PaymentMethod != "" {
		validPaymentMethods := []string{"cash", "transfer", "credit", "qris"}
		if !slices.Contains(validPaymentMethods, req.PaymentMethod) {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment method"})
			return
		}
		po.PaymentMethod = req.PaymentMethod
	}

	// Validate payment days if provided
	if req.PaymentDays != 0 {
		if req.PaymentDays < 0 || req.PaymentDays > 365 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Payment days must be between 0 and 365"})
			return
		}
		po.PaymentDays = req.PaymentDays
	}

	// Update basic fields
	if req.SupplierID != 0 {
		po.SupplierID = req.SupplierID
	}
	if req.Notes != "" {
		po.Notes = req.Notes
	}

	// Handle downpayment changes for credit orders
	if req.PaymentMethod == "credit" && req.DownPayment != originalDownPayment {
		// Validate downpayment doesn't exceed total
		if req.DownPayment > po.Total {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Downpayment cannot exceed total amount"})
			return
		}

		// Calculate the difference
		paymentDifference := req.DownPayment - originalDownPayment

		// Update amounts
		po.DownPayment = req.DownPayment
		po.AmountPaid += paymentDifference
		po.AmountDue -= paymentDifference

		// Record payment history for the adjustment
		if paymentDifference != 0 {
			paymentType := "adjustment"
			notes := fmt.Sprintf("Downpayment adjusted from Rp%.2f to Rp%.2f", originalDownPayment, req.DownPayment)

			payment := models.PurchasePayment{
				PurchaseOrderID: po.ID,
				UserID:          userID.(uint),
				Amount:          paymentDifference,
				PaymentMethod:   req.PaymentMethod,
				PaymentType:     paymentType,
				Notes:           notes,
			}

			if err := tx.Create(&payment).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record payment adjustment"})
				return
			}
		}

		// Update payment status if needed
		if po.AmountDue <= 0.01 {
			po.PaymentStatus = "paid"
			now := time.Now()
			po.PaidDate = &now
			po.AmountDue = 0
		} else {
			po.PaymentStatus = "pending"
		}
	}

	// Update due date if payment days changed
	if req.PaymentDays != 0 && req.PaymentDays != originalPaymentDays {
		if req.PaymentMethod == "credit" && req.PaymentDays > 0 {
			calculatedDueDate := po.OrderDate.AddDate(0, 0, req.PaymentDays)
			po.DueDate = &calculatedDueDate
		}
	}

	// Save the updated purchase order
	if err := tx.Save(&po).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update purchase order"})
		return
	}

	// Commit transaction
	tx.Commit()

	// Load complete purchase order with relationships
	var updatedPO models.PurchaseOrder
	database.DB.Preload("User").Preload("Supplier").Preload("Items.Product").Preload("Items.ProductSupplier.Supplier").First(&updatedPO, po.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Purchase order updated successfully",
		"data":    updatedPO,
	})
}
