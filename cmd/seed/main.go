package main

import (
	"fmt"
	"log"

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

	// Sample products
	products := []models.Product{
		{
			Name:        "Wireless Mouse",
			SKU:         "WM001",
			Description: "Bluetooth wireless mouse with optical sensor",
			Category:    "Electronics",
			Price:       29.99,
			Cost:        15.50,
			Quantity:    50,
			MinStock:    10,
			MaxStock:    200,
			Location:    "A1-B2",
			Supplier:    "Tech Supplies Inc",
			IsActive:    true,
		},
		{
			Name:        "USB Cable Type-C",
			SKU:         "UC001",
			Description: "High-speed USB Type-C cable 3ft",
			Category:    "Electronics",
			Price:       12.99,
			Cost:        6.00,
			Quantity:    100,
			MinStock:    20,
			MaxStock:    500,
			Location:    "A2-C1",
			Supplier:    "Cable World",
			IsActive:    true,
		},
		{
			Name:        "Office Chair",
			SKU:         "OC001",
			Description: "Ergonomic office chair with lumbar support",
			Category:    "Furniture",
			Price:       199.99,
			Cost:        120.00,
			Quantity:    15,
			MinStock:    5,
			MaxStock:    50,
			Location:    "B1-A1",
			Supplier:    "Office Furniture Co",
			IsActive:    true,
		},
		{
			Name:        "Notebook A4",
			SKU:         "NB001",
			Description: "Ruled notebook 200 pages",
			Category:    "Stationery",
			Price:       4.99,
			Cost:        2.50,
			Quantity:    200,
			MinStock:    50,
			MaxStock:    1000,
			Location:    "C1-D2",
			Supplier:    "Paper Plus",
			IsActive:    true,
		},
		{
			Name:        "Bluetooth Headphones",
			SKU:         "BH001",
			Description: "Noise-cancelling wireless headphones",
			Category:    "Electronics",
			Price:       89.99,
			Cost:        45.00,
			Quantity:    25,
			MinStock:    8,
			MaxStock:    100,
			Location:    "A1-C3",
			Supplier:    "Audio Tech",
			IsActive:    true,
		},
		{
			Name:        "Desk Lamp LED",
			SKU:         "DL001",
			Description: "Adjustable LED desk lamp with USB charging",
			Category:    "Lighting",
			Price:       39.99,
			Cost:        20.00,
			Quantity:    30,
			MinStock:    10,
			MaxStock:    150,
			Location:    "B2-A2",
			Supplier:    "Bright Lights Ltd",
			IsActive:    true,
		},
		{
			Name:        "Coffee Mug",
			SKU:         "CM001",
			Description: "Ceramic coffee mug 12oz",
			Category:    "Kitchen",
			Price:       8.99,
			Cost:        4.00,
			Quantity:    75,
			MinStock:    15,
			MaxStock:    300,
			Location:    "D1-B1",
			Supplier:    "Kitchen Supplies",
			IsActive:    true,
		},
		{
			Name:        "Smartphone Case",
			SKU:         "SC001",
			Description: "Protective smartphone case with screen protector",
			Category:    "Electronics",
			Price:       19.99,
			Cost:        8.50,
			Quantity:    60,
			MinStock:    15,
			MaxStock:    250,
			Location:    "A3-B1",
			Supplier:    "Phone Accessories Inc",
			IsActive:    true,
		},
		{
			Name:        "Desk Organizer",
			SKU:         "DO001",
			Description: "Multi-compartment desk organizer",
			Category:    "Office Supplies",
			Price:       24.99,
			Cost:        12.00,
			Quantity:    40,
			MinStock:    10,
			MaxStock:    120,
			Location:    "C2-A3",
			Supplier:    "Office Organizers",
			IsActive:    true,
		},
		{
			Name:        "Water Bottle",
			SKU:         "WB001",
			Description: "Stainless steel water bottle 750ml",
			Category:    "Lifestyle",
			Price:       16.99,
			Cost:        8.00,
			Quantity:    80,
			MinStock:    20,
			MaxStock:    400,
			Location:    "D2-C1",
			Supplier:    "Hydration Station",
			IsActive:    true,
		},
	}

	// Check if products already exist
	var count int64
	database.DB.Model(&models.Product{}).Count(&count)

	if count > 0 {
		fmt.Printf("Products already exist (%d products found). Skipping product seeding.\n", count)
	} else {
		for _, product := range products {
			result := database.DB.Create(&product)
			if result.Error != nil {
				log.Printf("Error creating product %s: %v", product.Name, result.Error)
			} else {
				fmt.Printf("Created product: %s (SKU: %s)\n", product.Name, product.SKU)
			}
		}
	}

	// Sample suppliers
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
	}

	// Check if suppliers already exist
	database.DB.Model(&models.Supplier{}).Count(&count)

	if count > 0 {
		fmt.Printf("Suppliers already exist (%d suppliers found). Skipping supplier seeding.\n", count)
	} else {
		for _, supplier := range suppliers {
			result := database.DB.Create(&supplier)
			if result.Error != nil {
				log.Printf("Error creating supplier %s: %v", supplier.Name, result.Error)
			} else {
				fmt.Printf("Created supplier: %s\n", supplier.Name)
			}
		}
	}

	fmt.Println("Sample data seeding completed!")
}
