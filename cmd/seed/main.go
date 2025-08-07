package main

import (
	"fmt"
	"log"
	"time"

	"inventory_system/database"
	"inventory_system/models"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Initialize database
	database.InitDatabase()

	// Seed sample data
	seedData()
}

func seedData() {
	fmt.Println("Seeding sample data...")

	// Sample suppliers first
	suppliers := []models.Supplier{
		{
			Name:          "Tech Supplies Inc",
			Email:         "orders@techsupplies.com",
			Phone:         "+1-555-0101",
			Address:       "123 Tech Street, Silicon Valley, CA 94000",
			ContactPerson: "John Smith",
			IsActive:      true,
		},
		{
			Name:          "Office Furniture Co",
			Email:         "sales@officefurniture.com",
			Phone:         "+1-555-0202",
			Address:       "456 Business Ave, New York, NY 10001",
			ContactPerson: "Sarah Johnson",
			IsActive:      true,
		},
		{
			Name:          "Paper Plus",
			Email:         "support@paperplus.com",
			Phone:         "+1-555-0303",
			Address:       "789 Paper Mill Rd, Portland, OR 97201",
			ContactPerson: "Mike Wilson",
			IsActive:      true,
		},
		{
			Name:          "Cable World",
			Email:         "info@cableworld.com",
			Phone:         "+1-555-0404",
			Address:       "321 Cable Ave, Austin, TX 73301",
			ContactPerson: "David Lee",
			IsActive:      true,
		},
		{
			Name:          "Audio Tech",
			Email:         "sales@audiotech.com",
			Phone:         "+1-555-0505",
			Address:       "654 Sound St, Nashville, TN 37201",
			ContactPerson: "Emma Davis",
			IsActive:      true,
		},
		{
			Name:          "Bright Lights Ltd",
			Email:         "orders@brightlights.com",
			Phone:         "+1-555-0606",
			Address:       "987 Light Blvd, Miami, FL 33101",
			ContactPerson: "Carlos Rodriguez",
			IsActive:      true,
		},
		{
			Name:          "Kitchen Supplies",
			Email:         "support@kitchensupplies.com",
			Phone:         "+1-555-0707",
			Address:       "147 Kitchen Rd, Chicago, IL 60601",
			ContactPerson: "Lisa Chen",
			IsActive:      true,
		},
		{
			Name:          "Phone Accessories Inc",
			Email:         "sales@phoneaccessories.com",
			Phone:         "+1-555-0808",
			Address:       "258 Mobile St, San Francisco, CA 94102",
			ContactPerson: "Ryan Kim",
			IsActive:      true,
		},
		{
			Name:          "Office Organizers",
			Email:         "info@officeorganizers.com",
			Phone:         "+1-555-0909",
			Address:       "369 Organize Ave, Boston, MA 02101",
			ContactPerson: "Jennifer White",
			IsActive:      true,
		},
		{
			Name:          "Hydration Station",
			Email:         "orders@hydrationstation.com",
			Phone:         "+1-555-1010",
			Address:       "741 Water Way, Denver, CO 80201",
			ContactPerson: "Michael Brown",
			IsActive:      true,
		},
	}

	// Check if suppliers already exist
	var supplierCount int64
	database.DB.Model(&models.Supplier{}).Count(&supplierCount)

	var supplierMap = make(map[string]uint)
	if supplierCount > 0 {
		fmt.Printf("Suppliers already exist (%d suppliers found). Using existing suppliers.\n", supplierCount)
		// Get existing suppliers
		var existingSuppliers []models.Supplier
		database.DB.Find(&existingSuppliers)
		for _, supplier := range existingSuppliers {
			supplierMap[supplier.Name] = supplier.ID
		}
	} else {
		for _, supplier := range suppliers {
			result := database.DB.Create(&supplier)
			if result.Error != nil {
				log.Printf("Error creating supplier %s: %v", supplier.Name, result.Error)
			} else {
				fmt.Printf("Created supplier: %s\n", supplier.Name)
				supplierMap[supplier.Name] = supplier.ID
			}
		}
	}

	// Sample products
	products := []models.Product{
		{
			Name:        "Wireless Mouse",
			SKU:         "WM001",
			Description: "Bluetooth wireless mouse with optical sensor",
			Category:    "Electronics",
			Location:    "A1-B2",
			IsActive:    true,
		},
		{
			Name:        "USB Cable Type-C",
			SKU:         "UC001",
			Description: "High-speed USB Type-C cable 3ft",
			Category:    "Electronics",
			Location:    "A2-C1",
			IsActive:    true,
		},
		{
			Name:        "Office Chair",
			SKU:         "OC001",
			Description: "Ergonomic office chair with lumbar support",
			Category:    "Furniture",
			Location:    "B1-A1",
			IsActive:    true,
		},
		{
			Name:        "Notebook A4",
			SKU:         "NB001",
			Description: "Ruled notebook 200 pages",
			Category:    "Stationery",
			Location:    "C1-D2",
			IsActive:    true,
		},
		{
			Name:        "Bluetooth Headphones",
			SKU:         "BH001",
			Description: "Noise-cancelling wireless headphones",
			Category:    "Electronics",
			Location:    "A1-C3",
			IsActive:    true,
		},
		{
			Name:        "Desk Lamp LED",
			SKU:         "DL001",
			Description: "Adjustable LED desk lamp with USB charging",
			Category:    "Lighting",
			Location:    "B2-A2",
			IsActive:    true,
		},
		{
			Name:        "Coffee Mug",
			SKU:         "CM001",
			Description: "Ceramic coffee mug 12oz",
			Category:    "Kitchen",
			Location:    "D1-B1",
			IsActive:    true,
		},
		{
			Name:        "Smartphone Case",
			SKU:         "SC001",
			Description: "Protective smartphone case with screen protector",
			Category:    "Electronics",
			Location:    "A3-B1",
			IsActive:    true,
		},
		{
			Name:        "Desk Organizer",
			SKU:         "DO001",
			Description: "Multi-compartment desk organizer",
			Category:    "Office Supplies",
			Location:    "C2-A3",
			IsActive:    true,
		},
		{
			Name:        "Water Bottle",
			SKU:         "WB001",
			Description: "Stainless steel water bottle 750ml",
			Category:    "Lifestyle",
			Location:    "D2-C1",
			IsActive:    true,
		},
	}

	// Product-supplier relationships with pricing and stock
	productSupplierData := []struct {
		ProductSKU   string
		SupplierName string
		Cost         float64
		Price        float64
		Stock        int
		MinStock     int
	}{
		{"WM001", "Tech Supplies Inc", 15.50, 29.99, 50, 10},
		{"UC001", "Cable World", 6.00, 12.99, 100, 20},
		{"OC001", "Office Furniture Co", 120.00, 199.99, 15, 5},
		{"NB001", "Paper Plus", 2.50, 4.99, 200, 50},
		{"BH001", "Audio Tech", 45.00, 89.99, 25, 8},
		{"DL001", "Bright Lights Ltd", 20.00, 39.99, 30, 10},
		{"CM001", "Kitchen Supplies", 4.00, 8.99, 75, 15},
		{"SC001", "Phone Accessories Inc", 8.50, 19.99, 60, 15},
		{"DO001", "Office Organizers", 12.00, 24.99, 40, 10},
		{"WB001", "Hydration Station", 8.00, 16.99, 80, 20},
		// Add some alternative suppliers for some products
		{"WM001", "Office Furniture Co", 16.00, 32.99, 25, 5},
		{"UC001", "Tech Supplies Inc", 6.50, 13.99, 75, 15},
		{"BH001", "Tech Supplies Inc", 47.00, 92.99, 15, 5},
	}

	// Check if products already exist
	var productCount int64
	database.DB.Model(&models.Product{}).Count(&productCount)

	var productMap = make(map[string]uint)
	if productCount > 0 {
		fmt.Printf("Products already exist (%d products found). Using existing products.\n", productCount)
		// Get existing products
		var existingProducts []models.Product
		database.DB.Find(&existingProducts)
		for _, product := range existingProducts {
			productMap[product.SKU] = product.ID
		}
	} else {
		for _, product := range products {
			result := database.DB.Create(&product)
			if result.Error != nil {
				log.Printf("Error creating product %s: %v", product.Name, result.Error)
			} else {
				fmt.Printf("Created product: %s (SKU: %s)\n", product.Name, product.SKU)
				productMap[product.SKU] = product.ID
			}
		}
	}

	// Create product-supplier relationships
	var psCount int64
	database.DB.Model(&models.ProductSupplier{}).Count(&psCount)

	if psCount > 0 {
		fmt.Printf("Product-supplier relationships already exist (%d relationships found). Skipping.\n", psCount)
	} else {
		for _, psData := range productSupplierData {
			productID, productExists := productMap[psData.ProductSKU]
			supplierID, supplierExists := supplierMap[psData.SupplierName]

			if !productExists || !supplierExists {
				log.Printf("Skipping product-supplier relationship: Product %s or Supplier %s not found", 
					psData.ProductSKU, psData.SupplierName)
				continue
			}

			productSupplier := models.ProductSupplier{
				ProductID:  productID,
				SupplierID: supplierID,
				Cost:       psData.Cost,
				Price:      psData.Price,
				Stock:      psData.Stock,
				MinStock:   psData.MinStock,
				IsActive:   true,
			}

			result := database.DB.Create(&productSupplier)
			if result.Error != nil {
				log.Printf("Error creating product-supplier relationship for %s-%s: %v", 
					psData.ProductSKU, psData.SupplierName, result.Error)
			} else {
				fmt.Printf("Created product-supplier relationship: %s -> %s (Stock: %d)\n", 
					psData.ProductSKU, psData.SupplierName, psData.Stock)
			}
		}
	}

	// Sample users
	users := []models.User{
		{
			Email:    "admin@inventory.com",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Name:     "Administrator",
			Role:     "admin",
			IsActive: true,
		},
		{
			Email:    "manager@inventory.com", 
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Name:     "Manager User",
			Role:     "manager",
			IsActive: true,
		},
		{
			Email:    "employee@inventory.com",
			Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
			Name:     "Employee User", 
			Role:     "employee",
			IsActive: true,
		},
	}

	// Check if users already exist
	var userCount int64
	database.DB.Model(&models.User{}).Count(&userCount)

	if userCount > 0 {
		fmt.Printf("Users already exist (%d users found). Skipping user seeding.\n", userCount)
	} else {
		for _, user := range users {
			result := database.DB.Create(&user)
			if result.Error != nil {
				log.Printf("Error creating user %s: %v", user.Email, result.Error)
			} else {
				fmt.Printf("Created user: %s (%s)\n", user.Name, user.Role)
			}
		}
	}

	// Sample purchase orders
	purchaseOrders := []models.PurchaseOrder{
		{
			PONumber:      "PO-2025-001",
			SupplierID:    supplierMap["Tech Supplies Inc"],
			UserID:        1, // Admin user
			PaymentMethod: "net30",
			PaymentDays:   30,
			PaymentStatus: "pending",
			Total:         1500.00,
			DownPayment:   0,
			AmountPaid:    0,
			AmountDue:     1500.00,
			Notes:         "Monthly electronics restock",
			OrderDate:     time.Now().AddDate(0, 0, -15), // 15 days ago
		},
		{
			PONumber:      "PO-2025-002", 
			SupplierID:    supplierMap["Office Furniture Co"],
			UserID:        1,
			PaymentMethod: "net15",
			PaymentDays:   15,
			PaymentStatus: "paid",
			Total:         2499.75,
			DownPayment:   500.00,
			AmountPaid:    2499.75,
			AmountDue:     0,
			Notes:         "Office chairs for new employees",
			OrderDate:     time.Now().AddDate(0, 0, -20),
		},
		{
			PONumber:      "PO-2025-003",
			SupplierID:    supplierMap["Paper Plus"],
			UserID:        1,
			PaymentMethod: "cash",
			PaymentDays:   0,
			PaymentStatus: "paid",
			Total:         299.50,
			DownPayment:   0,
			AmountPaid:    299.50,
			AmountDue:     0,
			Notes:         "Stationery supplies bulk order",
			OrderDate:     time.Now().AddDate(0, 0, -10),
		},
		{
			PONumber:      "PO-2025-004",
			SupplierID:    supplierMap["Cable World"],
			UserID:        1,
			PaymentMethod: "net60",
			PaymentDays:   60,
			PaymentStatus: "pending",
			Total:         850.25,
			DownPayment:   200.00,
			AmountPaid:    200.00,
			AmountDue:     650.25,
			Notes:         "USB cables and connectors",
			OrderDate:     time.Now().AddDate(0, 0, -5),
		},
		{
			PONumber:      "PO-2025-005",
			SupplierID:    supplierMap["Audio Tech"],
			UserID:        1,
			PaymentMethod: "net30",
			PaymentDays:   30,
			PaymentStatus: "overdue",
			Total:         1799.99,
			DownPayment:   0,
			AmountPaid:    0,
			AmountDue:     1799.99,
			Notes:         "Premium headphones for resale",
			OrderDate:     time.Now().AddDate(0, 0, -45), // 45 days ago (overdue)
		},
		{
			PONumber:      "PO-2025-006",
			SupplierID:    supplierMap["Bright Lights Ltd"],
			UserID:        1,
			PaymentMethod: "net15",
			PaymentDays:   15,
			PaymentStatus: "pending",
			Total:         679.80,
			DownPayment:   0,
			AmountPaid:    0,
			AmountDue:     679.80,
			Notes:         "LED lamps for office upgrade",
			OrderDate:     time.Now().AddDate(0, 0, -8),
		},
		{
			PONumber:      "PO-2025-007",
			SupplierID:    supplierMap["Kitchen Supplies"],
			UserID:        1,
			PaymentMethod: "credit",
			PaymentDays:   90,
			PaymentStatus: "pending",
			Total:         445.75,
			DownPayment:   100.00,
			AmountPaid:    100.00,
			AmountDue:     345.75,
			Notes:         "Break room supplies",
			OrderDate:     time.Now().AddDate(0, 0, -3),
		},
		{
			PONumber:      "PO-2025-008",
			SupplierID:    supplierMap["Phone Accessories Inc"],
			UserID:        1,
			PaymentMethod: "net30",
			PaymentDays:   30,
			PaymentStatus: "paid",
			Total:         1299.60,
			DownPayment:   0,
			AmountPaid:    1299.60,
			AmountDue:     0,
			Notes:         "Smartphone cases and accessories",
			OrderDate:     time.Now().AddDate(0, 0, -12),
		},
		{
			PONumber:      "PO-2025-009",
			SupplierID:    supplierMap["Office Organizers"],
			UserID:        1,
			PaymentMethod: "net7",
			PaymentDays:   7,
			PaymentStatus: "pending",
			Total:         324.95,
			DownPayment:   0,
			AmountPaid:    0,
			AmountDue:     324.95,
			Notes:         "Desk organizers and storage solutions",
			OrderDate:     time.Now().AddDate(0, 0, -2),
		},
		{
			PONumber:      "PO-2025-010",
			SupplierID:    supplierMap["Hydration Station"],
			UserID:        1,
			PaymentMethod: "net30",
			PaymentDays:   30,
			PaymentStatus: "pending",
			Total:         567.20,
			DownPayment:   0,
			AmountPaid:    0,
			AmountDue:     567.20,
			Notes:         "Water bottles for employee gifts",
			OrderDate:     time.Now().AddDate(0, 0, -6),
		},
		{
			PONumber:      "PO-2025-011",
			SupplierID:    supplierMap["Tech Supplies Inc"],
			UserID:        1,
			PaymentMethod: "net15",
			PaymentDays:   15,
			PaymentStatus: "overdue",
			Total:         999.99,
			DownPayment:   0,
			AmountPaid:    0,
			AmountDue:     999.99,
			Notes:         "Additional electronics inventory",
			OrderDate:     time.Now().AddDate(0, 0, -25), // 25 days ago (overdue)
		},
		{
			PONumber:      "PO-2025-012",
			SupplierID:    supplierMap["Cable World"],
			UserID:        1,
			PaymentMethod: "cash",
			PaymentDays:   0,
			PaymentStatus: "paid",
			Total:         199.95,
			DownPayment:   0,
			AmountPaid:    199.95,
			AmountDue:     0,
			Notes:         "Replacement cables",
			OrderDate:     time.Now().AddDate(0, 0, -1),
		},
		{
			PONumber:      "PO-2025-013",
			SupplierID:    supplierMap["Paper Plus"],
			UserID:        1,
			PaymentMethod: "net60",
			PaymentDays:   60,
			PaymentStatus: "pending",
			Total:         750.00,
			DownPayment:   150.00,
			AmountPaid:    150.00,
			AmountDue:     600.00,
			Notes:         "Bulk paper and notebooks",
			OrderDate:     time.Now().AddDate(0, 0, -4),
		},
		{
			PONumber:      "PO-2025-014",
			SupplierID:    supplierMap["Audio Tech"],
			UserID:        1,
			PaymentMethod: "net30",
			PaymentDays:   30,
			PaymentStatus: "pending",
			Total:         1250.50,
			DownPayment:   0,
			AmountPaid:    0,
			AmountDue:     1250.50,
			Notes:         "Audio equipment for conference room",
			OrderDate:     time.Now().AddDate(0, 0, -7),
		},
		{
			PONumber:      "PO-2025-015",
			SupplierID:    supplierMap["Office Furniture Co"],
			UserID:        1,
			PaymentMethod: "credit",
			PaymentDays:   90,
			PaymentStatus: "pending",
			Total:         3299.99,
			DownPayment:   1000.00,
			AmountPaid:    1000.00,
			AmountDue:     2299.99,
			Notes:         "Executive desk and chairs",
			OrderDate:     time.Now().AddDate(0, 0, -9),
		},
	}

	// Check if purchase orders already exist
	var poCount int64
	database.DB.Model(&models.PurchaseOrder{}).Count(&poCount)

	if poCount > 0 {
		fmt.Printf("Purchase orders already exist (%d orders found). Skipping purchase order seeding.\n", poCount)
	} else {
		for _, po := range purchaseOrders {
			// Set due date based on order date and payment days
			if po.PaymentDays > 0 {
				dueDate := po.OrderDate.AddDate(0, 0, po.PaymentDays)
				po.DueDate = &dueDate
			}
			
			// Set payment status to overdue if due date has passed
			if po.DueDate != nil && po.DueDate.Before(time.Now()) && po.PaymentStatus == "pending" {
				po.PaymentStatus = "overdue"
			}

			result := database.DB.Create(&po)
			if result.Error != nil {
				log.Printf("Error creating purchase order %s: %v", po.PONumber, result.Error)
			} else {
				fmt.Printf("Created purchase order: %s (Supplier: %s, Total: $%.2f)\n", 
					po.PONumber, getSupplierName(supplierMap, po.SupplierID), po.Total)
			}
		}
	}

	// Sample sales transactions
	sales := []models.Sale{
		{
			SaleNumber:    "SALE-2025-001",
			UserID:        1, // Admin user
			CustomerName:  "John Smith",
			Subtotal:      299.97,
			Tax:           23.98,
			Discount:      0,
			Total:         323.95,
			PaymentMethod: "cash",
			PaymentDays:   0,
			PaymentStatus: "paid",
			DownPayment:   0,
			AmountPaid:    323.95,
			AmountDue:     0,
			Status:        "completed",
		},
		{
			SaleNumber:    "SALE-2025-002",
			UserID:        1,
			CustomerName:  "Sarah Johnson",
			Subtotal:      1599.99,
			Tax:           127.99,
			Discount:      50.00,
			Total:         1677.98,
			PaymentMethod: "card",
			PaymentDays:   0,
			PaymentStatus: "paid",
			DownPayment:   0,
			AmountPaid:    1677.98,
			AmountDue:     0,
			Status:        "completed",
		},
		{
			SaleNumber:    "SALE-2025-003",
			UserID:        1,
			CustomerName:  "Mike Wilson",
			Subtotal:      89.99,
			Tax:           7.20,
			Discount:      0,
			Total:         97.19,
			PaymentMethod: "transfer",
			PaymentDays:   0,
			PaymentStatus: "paid",
			DownPayment:   0,
			AmountPaid:    97.19,
			AmountDue:     0,
			Status:        "completed",
		},
		{
			SaleNumber:    "SALE-2025-004",
			UserID:        1,
			CustomerName:  "Lisa Chen",
			Subtotal:      2499.75,
			Tax:           199.98,
			Discount:      100.00,
			Total:         2599.73,
			PaymentMethod: "credit",
			PaymentDays:   30,
			PaymentStatus: "pending",
			DownPayment:   800.00,
			AmountPaid:    800.00,
			AmountDue:     1799.73,
			Status:        "completed",
		},
		{
			SaleNumber:    "SALE-2025-005",
			UserID:        1,
			CustomerName:  "David Brown",
			Subtotal:      45.99,
			Tax:           3.68,
			Discount:      0,
			Total:         49.67,
			PaymentMethod: "cash",
			PaymentDays:   0,
			PaymentStatus: "paid",
			DownPayment:   0,
			AmountPaid:    49.67,
			AmountDue:     0,
			Status:        "completed",
		},
		{
			SaleNumber:    "SALE-2025-006",
			UserID:        1,
			CustomerName:  "Emma Davis",
			Subtotal:      179.98,
			Tax:           14.40,
			Discount:      20.00,
			Total:         174.38,
			PaymentMethod: "card",
			PaymentDays:   0,
			PaymentStatus: "paid",
			DownPayment:   0,
			AmountPaid:    174.38,
			AmountDue:     0,
			Status:        "completed",
		},
		{
			SaleNumber:    "SALE-2025-007",
			UserID:        1,
			CustomerName:  "Carlos Rodriguez",
			Subtotal:      679.80,
			Tax:           54.38,
			Discount:      0,
			Total:         734.18,
			PaymentMethod: "transfer",
			PaymentDays:   0,
			PaymentStatus: "paid",
			DownPayment:   0,
			AmountPaid:    734.18,
			AmountDue:     0,
			Status:        "completed",
		},
		{
			SaleNumber:    "SALE-2025-008",
			UserID:        1,
			CustomerName:  "Jennifer White",
			Subtotal:      124.95,
			Tax:           10.00,
			Discount:      15.00,
			Total:         119.95,
			PaymentMethod: "cash",
			PaymentDays:   0,
			PaymentStatus: "paid",
			DownPayment:   0,
			AmountPaid:    119.95,
			AmountDue:     0,
			Status:        "completed",
		},
		{
			SaleNumber:    "SALE-2025-009",
			UserID:        1,
			CustomerName:  "Michael Brown",
			Subtotal:      339.96,
			Tax:           27.20,
			Discount:      0,
			Total:         367.16,
			PaymentMethod: "credit",
			PaymentDays:   15,
			PaymentStatus: "overdue",
			DownPayment:   0,
			AmountPaid:    0,
			AmountDue:     367.16,
			Status:        "completed",
		},
		{
			SaleNumber:    "SALE-2025-010",
			UserID:        1,
			CustomerName:  "Ryan Kim",
			Subtotal:      59.97,
			Tax:           4.80,
			Discount:      0,
			Total:         64.77,
			PaymentMethod: "card",
			PaymentDays:   0,
			PaymentStatus: "paid",
			DownPayment:   0,
			AmountPaid:    64.77,
			AmountDue:     0,
			Status:        "completed",
		},
		{
			SaleNumber:    "SALE-2025-011",
			UserID:        1,
			CustomerName:  "Ashley Taylor",
			Subtotal:      1299.60,
			Tax:           103.97,
			Discount:      50.00,
			Total:         1353.57,
			PaymentMethod: "transfer",
			PaymentDays:   0,
			PaymentStatus: "paid",
			DownPayment:   0,
			AmountPaid:    1353.57,
			AmountDue:     0,
			Status:        "completed",
		},
		{
			SaleNumber:    "SALE-2025-012",
			UserID:        1,
			CustomerName:  "Kevin Martinez",
			Subtotal:      24.99,
			Tax:           2.00,
			Discount:      0,
			Total:         26.99,
			PaymentMethod: "cash",
			PaymentDays:   0,
			PaymentStatus: "paid",
			DownPayment:   0,
			AmountPaid:    26.99,
			AmountDue:     0,
			Status:        "completed",
		},
		{
			SaleNumber:    "SALE-2025-013",
			UserID:        1,
			CustomerName:  "Amanda Wilson",
			Subtotal:      449.95,
			Tax:           36.00,
			Discount:      25.00,
			Total:         460.95,
			PaymentMethod: "credit",
			PaymentDays:   30,
			PaymentStatus: "pending",
			DownPayment:   150.00,
			AmountPaid:    150.00,
			AmountDue:     310.95,
			Status:        "completed",
		},
		{
			SaleNumber:    "SALE-2025-014",
			UserID:        1,
			CustomerName:  "Chris Garcia",
			Subtotal:      799.99,
			Tax:           64.00,
			Discount:      0,
			Total:         863.99,
			PaymentMethod: "card",
			PaymentDays:   0,
			PaymentStatus: "paid",
			DownPayment:   0,
			AmountPaid:    863.99,
			AmountDue:     0,
			Status:        "completed",
		},
		{
			SaleNumber:    "SALE-2025-015",
			UserID:        1,
			CustomerName:  "Laura Anderson",
			Subtotal:      167.94,
			Tax:           13.44,
			Discount:      10.00,
			Total:         171.38,
			PaymentMethod: "cash",
			PaymentDays:   0,
			PaymentStatus: "paid",
			DownPayment:   0,
			AmountPaid:    171.38,
			AmountDue:     0,
			Status:        "completed",
		},
		{
			SaleNumber:    "SALE-2025-016",
			UserID:        1,
			CustomerName:  "Daniel Lee",
			Subtotal:      3299.99,
			Tax:           264.00,
			Discount:      200.00,
			Total:         3363.99,
			PaymentMethod: "credit",
			PaymentDays:   60,
			PaymentStatus: "pending",
			DownPayment:   1000.00,
			AmountPaid:    1000.00,
			AmountDue:     2363.99,
			Status:        "completed",
		},
		{
			SaleNumber:    "SALE-2025-017",
			UserID:        1,
			CustomerName:  "Jessica Thompson",
			Subtotal:      89.98,
			Tax:           7.20,
			Discount:      0,
			Total:         97.18,
			PaymentMethod: "transfer",
			PaymentDays:   0,
			PaymentStatus: "paid",
			DownPayment:   0,
			AmountPaid:    97.18,
			AmountDue:     0,
			Status:        "completed",
		},
		{
			SaleNumber:    "SALE-2025-018",
			UserID:        1,
			CustomerName:  "Thomas Clark",
			Subtotal:      1250.50,
			Tax:           100.04,
			Discount:      75.00,
			Total:         1275.54,
			PaymentMethod: "card",
			PaymentDays:   0,
			PaymentStatus: "paid",
			DownPayment:   0,
			AmountPaid:    1275.54,
			AmountDue:     0,
			Status:        "completed",
		},
	}

	// Check if sales already exist
	var salesCount int64
	database.DB.Model(&models.Sale{}).Count(&salesCount)

	if salesCount > 0 {
		fmt.Printf("Sales already exist (%d sales found). Skipping sales seeding.\n", salesCount)
	} else {
		for i, sale := range sales {
			// Set created date to spread across the last 30 days
			daysAgo := 30 - (i % 30)
			sale.CreatedAt = time.Now().AddDate(0, 0, -daysAgo)
			
			// Set due date for credit sales
			if sale.PaymentDays > 0 {
				dueDate := sale.CreatedAt.AddDate(0, 0, sale.PaymentDays)
				sale.DueDate = &dueDate
				
				// Set overdue status if due date has passed
				if dueDate.Before(time.Now()) && sale.PaymentStatus == "pending" {
					sale.PaymentStatus = "overdue"
				}
			}
			
			// Set paid date for completed payments
			if sale.PaymentStatus == "paid" {
				paidDate := sale.CreatedAt.Add(time.Hour * 2) // Paid 2 hours after creation
				sale.PaidDate = &paidDate
			}

			result := database.DB.Create(&sale)
			if result.Error != nil {
				log.Printf("Error creating sale %s: %v", sale.SaleNumber, result.Error)
			} else {
				fmt.Printf("Created sale: %s (Customer: %s, Total: $%.2f)\n", 
					sale.SaleNumber, sale.CustomerName, sale.Total)
			}
		}
	}

	// Create sample sale items for some sales
	var itemCount int64
	database.DB.Model(&models.SaleItem{}).Count(&itemCount)

	if itemCount == 0 {
		// Get created sales
		var createdSales []models.Sale
		database.DB.Limit(10).Find(&createdSales)
		
		// Get products for sale items
		var products []models.Product
		database.DB.Preload("Suppliers").Find(&products)
		
		if len(createdSales) > 0 && len(products) > 0 {
			saleItemsData := []struct {
				SaleIndex   int
				ProductSKU  string
				Quantity    int
			}{
				{0, "WM001", 3},   // SALE-2025-001: 3x Wireless Mouse
				{0, "UC001", 2},   // SALE-2025-001: 2x USB Cable
				{0, "NB001", 5},   // SALE-2025-001: 5x Notebook
				
				{1, "BH001", 1},   // SALE-2025-002: 1x Bluetooth Headphones
				{1, "SC001", 2},   // SALE-2025-002: 2x Smartphone Case
				
				{2, "BH001", 1},   // SALE-2025-003: 1x Bluetooth Headphones
				
				{3, "OC001", 1},   // SALE-2025-004: 1x Office Chair
				{3, "DL001", 2},   // SALE-2025-004: 2x Desk Lamp
				
				{4, "NB001", 9},   // SALE-2025-005: 9x Notebook
				{4, "CM001", 1},   // SALE-2025-005: 1x Coffee Mug
			}

			for _, itemData := range saleItemsData {
				if itemData.SaleIndex >= len(createdSales) {
					continue
				}
				
				sale := createdSales[itemData.SaleIndex]
				
				// Find product by SKU
				var product *models.Product
				for _, p := range products {
					if p.SKU == itemData.ProductSKU {
						product = &p
						break
					}
				}
				
				if product != nil && len(product.Suppliers) > 0 {
					supplier := product.Suppliers[0] // Use first supplier
					
					saleItem := models.SaleItem{
						SaleID:    sale.ID,
						ProductID: product.ID,
						Quantity:  itemData.Quantity,
						Price:     supplier.Price,
						Cost:      supplier.Cost,
						Total:     supplier.Price * float64(itemData.Quantity),
					}
					
					result := database.DB.Create(&saleItem)
					if result.Error != nil {
						log.Printf("Error creating sale item: %v", result.Error)
					} else {
						fmt.Printf("Created sale item: %s x%d for %s\n", 
							product.Name, itemData.Quantity, sale.SaleNumber)
					}
				}
			}
		}
	} else {
		fmt.Printf("Sale items already exist (%d items found). Skipping sale items seeding.\n", itemCount)
	}

	fmt.Println("Sample data seeding completed!")
}

// Helper function to get supplier name by ID
func getSupplierName(supplierMap map[string]uint, supplierID uint) string {
	for name, id := range supplierMap {
		if id == supplierID {
			return name
		}
	}
	return "Unknown"
}
