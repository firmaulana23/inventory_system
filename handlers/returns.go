package handlers

import (
	"fmt"
	"inventory_system/database"
	"inventory_system/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ReturnRequest represents the request structure for creating a return
type ReturnRequest struct {
	SaleID       uint                `json:"sale_id" binding:"required"`
	Reason       string              `json:"reason"`
	RefundMethod string              `json:"refund_method" binding:"required"`
	Items        []ReturnItemRequest `json:"items" binding:"required,min=1"`
}

// ReturnItemRequest represents the request structure for return items
type ReturnItemRequest struct {
	SaleItemID uint    `json:"sale_item_id" binding:"required"`
	ProductID  uint    `json:"product_id" binding:"required"`
	Quantity   int     `json:"quantity" binding:"required,min=1"`
	Price      float64 `json:"price" binding:"required,min=0"`
	Condition  string  `json:"condition" binding:"required"`
}

// ExchangeRequest represents the request structure for creating an exchange
type ExchangeRequest struct {
	SaleID        uint                     `json:"sale_id" binding:"required"`
	Reason        string                   `json:"reason"`
	PaymentMethod string                   `json:"payment_method"` // For difference payment/refund
	OldItems      []ExchangeOldItemRequest `json:"old_items" binding:"required,min=1"`
	NewItems      []ExchangeNewItemRequest `json:"new_items" binding:"required,min=1"`
}

// ExchangeOldItemRequest represents the request structure for old items in exchange
type ExchangeOldItemRequest struct {
	SaleItemID uint    `json:"sale_item_id" binding:"required"`
	ProductID  uint    `json:"product_id" binding:"required"`
	Quantity   int     `json:"quantity" binding:"required,min=1"`
	Price      float64 `json:"price" binding:"required,min=0"`
	Condition  string  `json:"condition" binding:"required"`
}

// ExchangeNewItemRequest represents the request structure for new items in exchange
type ExchangeNewItemRequest struct {
	ProductID  uint    `json:"product_id" binding:"required"`
	Quantity   int     `json:"quantity" binding:"required,min=1"`
	Price      float64 `json:"price" binding:"required,min=0"`
}

// CreateReturn handles creating a new return
func CreateReturn(c *gin.Context) {
	var request ReturnRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate refund method
	validRefundMethods := []string{"cash", "card", "transfer", "store_credit"}
	validMethod := false
	for _, method := range validRefundMethods {
		if request.RefundMethod == method {
			validMethod = true
			break
		}
	}
	if !validMethod {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid refund method"})
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
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Verify sale exists and get sale details
	var sale models.Sale
	if err := tx.Preload("Items.Product").First(&sale, request.SaleID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Sale not found"})
		return
	}

	// Generate return number
	returnNumber := fmt.Sprintf("RET-%d-%d", time.Now().Year(), time.Now().Unix()%100000)

	// Calculate totals
	var subtotal, total, totalCost float64
	returnItems := make([]models.ReturnItem, len(request.Items))

	for i, itemReq := range request.Items {
		// Validate sale item exists and belongs to this sale
		var saleItem models.SaleItem
		found := false
		for _, si := range sale.Items {
			if si.ID == itemReq.SaleItemID && si.ProductID == itemReq.ProductID {
				saleItem = si
				found = true
				break
			}
		}
		
		if !found {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Sale item %d not found in sale %d", itemReq.SaleItemID, request.SaleID)})
			return
		}

		// Validate quantity doesn't exceed original quantity
		if itemReq.Quantity > saleItem.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Return quantity (%d) cannot exceed original quantity (%d) for product %s", 
				itemReq.Quantity, saleItem.Quantity, saleItem.Product.Name)})
			return
		}

		// Validate condition
		validConditions := []string{"good", "damaged", "expired"}
		validCondition := false
		for _, condition := range validConditions {
			if itemReq.Condition == condition {
				validCondition = true
				break
			}
		}
		if !validCondition {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item condition. Must be: good, damaged, or expired"})
			return
		}

		itemTotal := float64(itemReq.Quantity) * itemReq.Price
		itemTotalCost := float64(itemReq.Quantity) * saleItem.Cost // Use original cost from sale
		subtotal += itemTotal
		totalCost += itemTotalCost

		returnItems[i] = models.ReturnItem{
			SaleItemID: itemReq.SaleItemID,
			ProductID:  itemReq.ProductID,
			Quantity:   itemReq.Quantity,
			Price:      itemReq.Price,
			Cost:       saleItem.Cost,
			Total:      itemTotal,
			TotalCost:  itemTotalCost,
			Condition:  itemReq.Condition,
		}
	}

	// Calculate proportional tax and discount
	originalRatio := subtotal / sale.Subtotal
	tax := sale.Tax * originalRatio
	discount := sale.Discount * originalRatio
	total = subtotal + tax - discount

	// Calculate profit loss (negative value since it's a loss)
	profitLoss := -(total - totalCost) // Original profit that is now lost

	// Create return record
	returnRecord := models.Return{
		ReturnNumber: returnNumber,
		SaleID:       request.SaleID,
		UserID:       userID.(uint),
		Subtotal:     subtotal,
		Tax:          tax,
		Discount:     discount,
		Total:        total,
		Reason:       request.Reason,
		RefundMethod: request.RefundMethod,
		RefundAmount: total,
		TotalCost:    totalCost,
		ProfitLoss:   profitLoss,
	}

	if err := tx.Create(&returnRecord).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create return"})
		return
	}

	// Set return ID for items
	for i := range returnItems {
		returnItems[i].ReturnID = returnRecord.ID
	}

	// Create return items
	if err := tx.Create(&returnItems).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create return items"})
		return
	}

	// Update stock for returned items (only for items in good condition)
	for _, item := range returnItems {
		if item.Condition == "good" {
			// Find the product and update stock
			var product models.Product
			if err := tx.Preload("Suppliers").First(&product, item.ProductID).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find product for stock update"})
				return
			}

			// Add stock back to the first active supplier (or create logic to determine which supplier)
			if len(product.Suppliers) > 0 {
				for i, supplier := range product.Suppliers {
					if supplier.IsActive {
						product.Suppliers[i].Stock += item.Quantity
						if err := tx.Save(&product.Suppliers[i]).Error; err != nil {
							tx.Rollback()
							c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update supplier stock"})
							return
						}
						break
					}
				}
			}

			// Record stock movement
			movement := models.StockMovement{
				ProductID: item.ProductID,
				UserID:    userID.(uint),
				Type:      "in",
				Quantity:  item.Quantity,
				Reference: returnRecord.ReturnNumber,
				Notes:     fmt.Sprintf("Return from sale %s - condition: %s", sale.SaleNumber, item.Condition),
			}

			if err := tx.Create(&movement).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record stock movement"})
				return
			}
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete return transaction"})
		return
	}

	// Get complete return with relations
	var completeReturn models.Return
	database.DB.Preload("Items.Product").Preload("Items.SaleItem").Preload("Sale").Preload("User").First(&completeReturn, returnRecord.ID)

	c.JSON(http.StatusCreated, completeReturn)
}

// GetReturns handles getting returns with pagination and filters
func GetReturns(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	query := database.DB.Model(&models.Return{}).Preload("Items.Product").Preload("Sale").Preload("User")
	countQuery := database.DB.Model(&models.Return{})

	// Apply filters
	filters := make(map[string]interface{})

	// Filter by date range
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("created_at >= ?", startDate)
		countQuery = countQuery.Where("created_at >= ?", startDate)
		filters["start_date"] = startDate
	}

	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("created_at <= ?", endDate)
		countQuery = countQuery.Where("created_at <= ?", endDate)
		filters["end_date"] = endDate
	}

	// Filter by refund method
	if refundMethod := c.Query("refund_method"); refundMethod != "" {
		query = query.Where("refund_method = ?", refundMethod)
		countQuery = countQuery.Where("refund_method = ?", refundMethod)
		filters["refund_method"] = refundMethod
	}

	// Get total count
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count returns"})
		return
	}

	// Get returns
	var returns []models.Return
	result := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&returns)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch returns"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"returns": returns,
		"total":   total,
		"page":    page,
		"limit":   limit,
		"filters": filters,
	})
}

// GetReturn handles getting a single return by ID
func GetReturn(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid return ID"})
		return
	}

	var returnRecord models.Return
	result := database.DB.Preload("Items.Product").Preload("Items.SaleItem").Preload("Sale").Preload("User").First(&returnRecord, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Return not found"})
		return
	}

	c.JSON(http.StatusOK, returnRecord)
}

// GetSaleForReturn handles getting sale details for return processing
func GetSaleForReturn(c *gin.Context) {
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

	// Get existing returns for this sale to calculate available quantities
	var existingReturns []models.Return
	database.DB.Preload("Items").Where("sale_id = ?", id).Find(&existingReturns)

	// Calculate returned quantities per sale item
	returnedQuantities := make(map[uint]int)
	for _, ret := range existingReturns {
		for _, item := range ret.Items {
			returnedQuantities[item.SaleItemID] += item.Quantity
		}
	}

	// Add available quantity to each sale item
	type SaleItemWithAvailable struct {
		models.SaleItem
		AvailableQuantity int `json:"available_quantity"`
	}

	type SaleWithAvailable struct {
		models.Sale
		Items []SaleItemWithAvailable `json:"items"`
	}

	saleWithAvailable := SaleWithAvailable{
		Sale:  sale,
		Items: make([]SaleItemWithAvailable, len(sale.Items)),
	}

	for i, item := range sale.Items {
		availableQty := item.Quantity - returnedQuantities[item.ID]
		saleWithAvailable.Items[i] = SaleItemWithAvailable{
			SaleItem:          item,
			AvailableQuantity: availableQty,
		}
	}

	c.JSON(http.StatusOK, saleWithAvailable)
}

// CreateExchange handles creating a new product exchange
func CreateExchange(c *gin.Context) {
	var request ExchangeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
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
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Verify sale exists and get sale details
	var sale models.Sale
	if err := tx.Preload("Items.Product").First(&sale, request.SaleID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Sale not found"})
		return
	}

	// Generate exchange number
	exchangeNumber := fmt.Sprintf("EXC-%d-%d", time.Now().Year(), time.Now().Unix()%100000)

	// Validate and calculate old items total
	var totalOldValue, totalOldCost float64
	exchangeOldItems := make([]models.ExchangeOldItem, len(request.OldItems))

	for i, oldItemReq := range request.OldItems {
		// Validate sale item exists and belongs to this sale
		var saleItem models.SaleItem
		found := false
		for _, si := range sale.Items {
			if si.ID == oldItemReq.SaleItemID && si.ProductID == oldItemReq.ProductID {
				saleItem = si
				found = true
				break
			}
		}
		
		if !found {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Sale item %d not found in sale %d", oldItemReq.SaleItemID, request.SaleID)})
			return
		}

		// Validate quantity doesn't exceed original quantity
		if oldItemReq.Quantity > saleItem.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Exchange quantity (%d) cannot exceed original quantity (%d) for product %s", 
				oldItemReq.Quantity, saleItem.Quantity, saleItem.Product.Name)})
			return
		}

		// Validate condition
		validConditions := []string{"good", "damaged", "expired"}
		validCondition := false
		for _, condition := range validConditions {
			if oldItemReq.Condition == condition {
				validCondition = true
				break
			}
		}
		if !validCondition {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item condition. Must be: good, damaged, or expired"})
			return
		}

		itemTotal := float64(oldItemReq.Quantity) * oldItemReq.Price
		itemTotalCost := float64(oldItemReq.Quantity) * saleItem.Cost
		totalOldValue += itemTotal
		totalOldCost += itemTotalCost

		exchangeOldItems[i] = models.ExchangeOldItem{
			SaleItemID: oldItemReq.SaleItemID,
			ProductID:  oldItemReq.ProductID,
			Quantity:   oldItemReq.Quantity,
			Price:      oldItemReq.Price,
			Cost:       saleItem.Cost,
			Total:      itemTotal,
			TotalCost:  itemTotalCost,
			Condition:  oldItemReq.Condition,
		}
	}

	// Validate and calculate new items total
	var totalNewValue, totalNewCost float64
	exchangeNewItems := make([]models.ExchangeNewItem, len(request.NewItems))

	for i, newItemReq := range request.NewItems {
		// Validate product exists and has sufficient stock
		var product models.Product
		if err := tx.Preload("Suppliers").First(&product, newItemReq.ProductID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Product with ID %d not found", newItemReq.ProductID)})
			return
		}

		// Check stock availability and get cost from active supplier
		totalStock := 0
		var lowestCost float64
		firstCost := true
		for _, supplier := range product.Suppliers {
			if supplier.IsActive {
				totalStock += supplier.Stock
				if firstCost || supplier.Cost < lowestCost {
					lowestCost = supplier.Cost
					firstCost = false
				}
			}
		}

		if totalStock < newItemReq.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Insufficient stock for product %s. Available: %d, Requested: %d",
					product.Name, totalStock, newItemReq.Quantity),
			})
			return
		}

		itemTotal := float64(newItemReq.Quantity) * newItemReq.Price
		itemTotalCost := float64(newItemReq.Quantity) * lowestCost
		totalNewValue += itemTotal
		totalNewCost += itemTotalCost

		exchangeNewItems[i] = models.ExchangeNewItem{
			ProductID: newItemReq.ProductID,
			Quantity:  newItemReq.Quantity,
			Price:     newItemReq.Price,
			Cost:      lowestCost,
			Total:     itemTotal,
			TotalCost: itemTotalCost,
		}
	}

	// Calculate price difference
	difference := totalNewValue - totalOldValue

	// Calculate profit impact
	// Old profit (lost): (totalOldValue - totalOldCost) * -1 (negative because it's lost)
	// New profit (gained): totalNewValue - totalNewCost
	// Net profit impact = New profit - Lost profit = (totalNewValue - totalNewCost) - (totalOldValue - totalOldCost)
	profitImpact := (totalNewValue - totalNewCost) - (totalOldValue - totalOldCost)

	// Create exchange record
	exchange := models.Exchange{
		ExchangeNumber: exchangeNumber,
		SaleID:         request.SaleID,
		UserID:         userID.(uint),
		Reason:         request.Reason,
		TotalOldValue:  totalOldValue,
		TotalNewValue:  totalNewValue,
		TotalOldCost:   totalOldCost,
		TotalNewCost:   totalNewCost,
		Difference:     difference,
		ProfitImpact:   profitImpact,
		PaymentMethod:  request.PaymentMethod,
	}

	if err := tx.Create(&exchange).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create exchange"})
		return
	}

	// Set exchange ID for items
	for i := range exchangeOldItems {
		exchangeOldItems[i].ExchangeID = exchange.ID
	}
	for i := range exchangeNewItems {
		exchangeNewItems[i].ExchangeID = exchange.ID
	}

	// Create exchange items
	if err := tx.Create(&exchangeOldItems).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create exchange old items"})
		return
	}

	if err := tx.Create(&exchangeNewItems).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create exchange new items"})
		return
	}

	// Process stock changes
	// 1. Return old items to stock (only good condition)
	for _, oldItem := range exchangeOldItems {
		if oldItem.Condition == "good" {
			var product models.Product
			if err := tx.Preload("Suppliers").First(&product, oldItem.ProductID).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find product for stock update"})
				return
			}

			// Add stock back to the first active supplier
			if len(product.Suppliers) > 0 {
				for j, supplier := range product.Suppliers {
					if supplier.IsActive {
						product.Suppliers[j].Stock += oldItem.Quantity
						if err := tx.Save(&product.Suppliers[j]).Error; err != nil {
							tx.Rollback()
							c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update supplier stock"})
							return
						}
						break
					}
				}
			}

			// Record stock movement for returned item
			movement := models.StockMovement{
				ProductID: oldItem.ProductID,
				UserID:    userID.(uint),
				Type:      "in",
				Quantity:  oldItem.Quantity,
				Reference: exchange.ExchangeNumber,
				Notes:     fmt.Sprintf("Exchange return - condition: %s", oldItem.Condition),
			}

			if err := tx.Create(&movement).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record stock movement"})
				return
			}
		}
	}

	// 2. Reduce stock for new items
	for _, newItem := range exchangeNewItems {
		var product models.Product
		if err := tx.Preload("Suppliers").First(&product, newItem.ProductID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find product for stock update"})
			return
		}

		// Reduce stock from suppliers (FIFO basis)
		remainingQty := newItem.Quantity
		for j, supplier := range product.Suppliers {
			if supplier.IsActive && remainingQty > 0 {
				deductQty := remainingQty
				if supplier.Stock < remainingQty {
					deductQty = supplier.Stock
				}
				
				product.Suppliers[j].Stock -= deductQty
				remainingQty -= deductQty
				
				if err := tx.Save(&product.Suppliers[j]).Error; err != nil {
					tx.Rollback()
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update supplier stock"})
					return
				}
				
				if remainingQty == 0 {
					break
				}
			}
		}

		// Record stock movement for new item
		movement := models.StockMovement{
			ProductID: newItem.ProductID,
			UserID:    userID.(uint),
			Type:      "out",
			Quantity:  newItem.Quantity,
			Reference: exchange.ExchangeNumber,
			Notes:     "Exchange new item",
		}

		if err := tx.Create(&movement).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record stock movement"})
			return
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete exchange transaction"})
		return
	}

	// Get complete exchange with relations
	var completeExchange models.Exchange
	database.DB.Preload("OldItems.Product").Preload("OldItems.SaleItem").
		Preload("NewItems.Product").Preload("Sale").Preload("User").
		First(&completeExchange, exchange.ID)

	c.JSON(http.StatusCreated, completeExchange)
}

// GetExchanges handles getting exchanges with pagination and filters
func GetExchanges(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	query := database.DB.Model(&models.Exchange{}).Preload("OldItems.Product").Preload("NewItems.Product").Preload("Sale").Preload("User")
	countQuery := database.DB.Model(&models.Exchange{})

	// Apply filters
	filters := make(map[string]interface{})

	// Filter by date range
	if startDate := c.Query("start_date"); startDate != "" {
		query = query.Where("created_at >= ?", startDate)
		countQuery = countQuery.Where("created_at >= ?", startDate)
		filters["start_date"] = startDate
	}

	if endDate := c.Query("end_date"); endDate != "" {
		query = query.Where("created_at <= ?", endDate)
		countQuery = countQuery.Where("created_at <= ?", endDate)
		filters["end_date"] = endDate
	}

	// Get total count
	var total int64
	if err := countQuery.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count exchanges"})
		return
	}

	// Get exchanges
	var exchanges []models.Exchange
	result := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&exchanges)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch exchanges"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"exchanges": exchanges,
		"total":     total,
		"page":      page,
		"limit":     limit,
		"filters":   filters,
	})
}

// GetExchange handles getting a single exchange by ID
func GetExchange(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exchange ID"})
		return
	}

	var exchange models.Exchange
	result := database.DB.Preload("OldItems.Product").Preload("OldItems.SaleItem").
		Preload("NewItems.Product").Preload("Sale").Preload("User").
		First(&exchange, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exchange not found"})
		return
	}

	c.JSON(http.StatusOK, exchange)
}