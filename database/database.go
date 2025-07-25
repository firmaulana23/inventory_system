package database

import (
	"fmt"
	"inventory_system/models"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDatabase initializes the database connection
func InitDatabase() {
	var err error

	// Build PostgreSQL connection string
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres"
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "inventory_system"
	}

	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	timezone := os.Getenv("DB_TIMEZONE")
	if timezone == "" {
		timezone = "UTC"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		host, user, password, dbname, port, sslmode, timezone)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")

	// Auto migrate the schema first
	err = DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Supplier{},
		&models.ProductSupplier{},
		&models.StockMovement{},
		&models.Sale{},
		&models.SaleItem{},
		&models.SalePayment{},
		&models.PurchaseOrder{},
		&models.PurchaseOrderItem{},
		&models.PurchasePayment{},
		&models.ActivityLog{},
		&models.CompanyProfile{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database schema migration completed")

	// Create default admin user
	createDefaultAdmin()
}

// createDefaultAdmin creates a default admin user if none exists
func createDefaultAdmin() {
	var count int64
	DB.Model(&models.User{}).Where("role = ?", "admin").Count(&count)

	if count == 0 {
		adminEmail := os.Getenv("ADMIN_EMAIL")
		adminPassword := os.Getenv("ADMIN_PASSWORD")

		if adminEmail == "" {
			adminEmail = "admin@inventory.com"
		}
		if adminPassword == "" {
			adminPassword = "admin123"
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing admin password: %v", err)
			return
		}

		admin := models.User{
			Email:    adminEmail,
			Password: string(hashedPassword),
			Name:     "System Administrator",
			Role:     "admin",
			IsActive: true,
		}

		result := DB.Create(&admin)
		if result.Error != nil {
			log.Printf("Error creating admin user: %v", result.Error)
		} else {
			fmt.Printf("Default admin user created: %s\n", adminEmail)
		}
	}
}


// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
