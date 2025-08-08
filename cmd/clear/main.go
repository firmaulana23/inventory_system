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

	fmt.Println("Clearing all data except admin user...")

	// Clear data in reverse dependency order
	database.DB.Where("1=1").Delete(&models.ActivityLog{})
	database.DB.Where("1=1").Delete(&models.StockMovement{})
	database.DB.Where("1=1").Delete(&models.SalePayment{})
	database.DB.Where("1=1").Delete(&models.PurchasePayment{})
	database.DB.Where("1=1").Delete(&models.PurchaseOrderItem{})
	database.DB.Where("1=1").Delete(&models.ReturnItem{})
	database.DB.Where("1=1").Delete(&models.ExchangeOldItem{})
	database.DB.Where("1=1").Delete(&models.ExchangeNewItem{})
	database.DB.Where("1=1").Delete(&models.Exchange{})
	database.DB.Where("1=1").Delete(&models.Return{})
	database.DB.Where("1=1").Delete(&models.SaleItem{})
	database.DB.Where("1=1").Delete(&models.Sale{})
	database.DB.Where("1=1").Delete(&models.PurchaseOrder{})
	database.DB.Where("1=1").Delete(&models.ProductSupplier{})
	database.DB.Where("1=1").Delete(&models.Product{})
	database.DB.Where("1=1").Delete(&models.Supplier{})
	database.DB.Where("role != ?", "admin").Delete(&models.User{})
	database.DB.Where("1=1").Delete(&models.CompanyProfile{})

	fmt.Println("Data cleared successfully! Now you can run the seed program.")
}