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

	// Run custom migrations first
	err = runCustomMigrations()
	if err != nil {
		log.Fatal("Failed to run custom migrations:", err)
	}

	// Auto migrate the schema
	err = DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.StockMovement{},
		&models.Sale{},
		&models.SaleItem{},
		&models.Supplier{},
		&models.PurchaseOrder{},
		&models.PurchaseOrderItem{},
		&models.PurchasePayment{},
		&models.ActivityLog{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database migration completed")

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

// runCustomMigrations handles custom database migrations
func runCustomMigrations() error {
	// Check if HPP column exists
	var count int64
	err := DB.Raw("SELECT COUNT(*) FROM information_schema.columns WHERE table_name = 'products' AND column_name = 'hpp'").Scan(&count).Error
	if err != nil {
		return err
	}

	// If HPP column doesn't exist, add it with existing data handling
	if count == 0 {
		log.Println("Adding HPP column to products table...")
		
		// Add column as nullable first
		err = DB.Exec("ALTER TABLE products ADD COLUMN hpp DECIMAL(10,2)").Error
		if err != nil {
			return err
		}

		// Update existing products to set HPP equal to cost
		err = DB.Exec("UPDATE products SET hpp = COALESCE(cost, 0)").Error
		if err != nil {
			return err
		}

		// Make column NOT NULL with default
		err = DB.Exec("ALTER TABLE products ALTER COLUMN hpp SET NOT NULL").Error
		if err != nil {
			return err
		}

		err = DB.Exec("ALTER TABLE products ALTER COLUMN hpp SET DEFAULT 0.00").Error
		if err != nil {
			return err
		}

		log.Println("HPP column added successfully")
	}

	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
