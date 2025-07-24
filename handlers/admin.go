package handlers

import (
	"fmt"
	"inventory_system/database"
	"inventory_system/models"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// DashboardStats represents admin dashboard statistics
type DashboardStats struct {
	TotalUsers         int64          `json:"total_users"`
	TotalProducts      int64          `json:"total_products"`
	TotalSales         int64          `json:"total_sales"`
	TodaySales         int64          `json:"today_sales"`
	TotalRevenue       float64        `json:"total_revenue"`
	TodayRevenue       float64        `json:"today_revenue"`
	TotalProfit        float64        `json:"total_profit"`
	TodayProfit        float64        `json:"today_profit"`
	TotalPurchasing    float64        `json:"total_purchasing"`
	TotalPurchasingPaid float64       `json:"total_purchasing_paid"`
	TotalPurchasingDue  float64       `json:"total_purchasing_due"`
	LowStockProducts   int64          `json:"low_stock_products"`
	RecentSales        []models.Sale  `json:"recent_sales"`
	TopProducts        []ProductStats `json:"top_products"`
	SalesChart         []SalesChart   `json:"sales_chart"`
}

type ProductStats struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	TotalSold   int     `json:"total_sold"`
	Revenue     float64 `json:"revenue"`
}

type SalesChart struct {
	Date    string  `json:"date"`
	Sales   int64   `json:"sales"`
	Revenue float64 `json:"revenue"`
}

// GetDashboardStats returns admin dashboard statistics
func GetDashboardStats(c *gin.Context) {
	var stats DashboardStats

	// Get current date for today's stats
	today := time.Now().Format("2006-01-02")

	// Count total users
	database.DB.Model(&models.User{}).Count(&stats.TotalUsers)

	// Count total products
	database.DB.Model(&models.Product{}).Count(&stats.TotalProducts)

	// Count total sales
	database.DB.Model(&models.Sale{}).Count(&stats.TotalSales)

	// Count today's sales
	database.DB.Model(&models.Sale{}).Where("created_at::date = ?", today).Count(&stats.TodaySales)

	// Calculate total revenue
	database.DB.Model(&models.Sale{}).Select("COALESCE(SUM(total), 0)").Scan(&stats.TotalRevenue)

	// Calculate today's revenue
	database.DB.Model(&models.Sale{}).Where("created_at::date = ?", today).Select("COALESCE(SUM(total), 0)").Scan(&stats.TodayRevenue)

	// Calculate total profit (revenue - cost)
	database.DB.Raw(`
		SELECT COALESCE(SUM(si.quantity * (si.price - si.cost)), 0) 
		FROM sale_items si
	`).Scan(&stats.TotalProfit)

	// Calculate today's profit
	database.DB.Raw(`
		SELECT COALESCE(SUM(si.quantity * (si.price - si.cost)), 0) 
		FROM sale_items si
		JOIN sales s ON si.sale_id = s.id
		WHERE s.created_at::date = ?
	`, today).Scan(&stats.TodayProfit)

	// Calculate total purchasing price
	database.DB.Model(&models.PurchaseOrder{}).Select("COALESCE(SUM(total), 0)").Scan(&stats.TotalPurchasing)

	// Calculate total purchasing amount paid
	database.DB.Model(&models.PurchaseOrder{}).Select("COALESCE(SUM(amount_paid), 0)").Scan(&stats.TotalPurchasingPaid)

	// Calculate total purchasing amount due
	database.DB.Model(&models.PurchaseOrder{}).Select("COALESCE(SUM(amount_due), 0)").Scan(&stats.TotalPurchasingDue)

	// Count low stock products (products where any supplier has stock <= min_stock)
	database.DB.Raw(`
		SELECT COUNT(DISTINCT p.id) 
		FROM products p 
		JOIN product_suppliers ps ON p.id = ps.product_id 
		WHERE ps.is_active = true AND ps.stock <= ps.min_stock
	`).Scan(&stats.LowStockProducts)

	// Get recent sales (last 10)
	database.DB.Preload("User").Preload("Items.Product").Order("created_at desc").Limit(10).Find(&stats.RecentSales)

	// Get top products by sales
	database.DB.Raw(`
		SELECT si.product_id, p.name as product_name, 
		       SUM(si.quantity) as total_sold, 
		       SUM(si.quantity * si.price) as revenue
		FROM sale_items si
		JOIN products p ON si.product_id = p.id
		GROUP BY si.product_id, p.name
		ORDER BY total_sold DESC
		LIMIT 5
	`).Scan(&stats.TopProducts)

	// Get sales chart data for last 7 days
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		var chartData SalesChart
		chartData.Date = date

		database.DB.Model(&models.Sale{}).Where("created_at::date = ?", date).Count(&chartData.Sales)
		database.DB.Model(&models.Sale{}).Where("created_at::date = ?", date).Select("COALESCE(SUM(total), 0)").Scan(&chartData.Revenue)

		stats.SalesChart = append(stats.SalesChart, chartData)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// GetSystemLogs returns system activity logs
func GetSystemLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	var logs []models.ActivityLog
	var total int64

	offset := (page - 1) * limit

	query := database.DB.Model(&models.ActivityLog{}).Preload("User")
	query.Count(&total)
	query.Order("created_at desc").Offset(offset).Limit(limit).Find(&logs)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"logs":  logs,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// GetUserActivity returns specific user activity
func GetUserActivity(c *gin.Context) {
	userID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	var logs []models.ActivityLog
	var total int64

	offset := (page - 1) * limit

	query := database.DB.Model(&models.ActivityLog{}).Where("user_id = ?", userID).Preload("User")
	query.Count(&total)
	query.Order("created_at desc").Offset(offset).Limit(limit).Find(&logs)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"logs":  logs,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// GetAdminSalesReport returns detailed sales report for admin dashboard
func GetAdminSalesReport(c *gin.Context) {
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	type SalesReportData struct {
		Date         string  `json:"date"`
		TotalSales   int64   `json:"total_sales"`
		TotalRevenue float64 `json:"total_revenue"`
		TotalItems   int64   `json:"total_items"`
	}

	var reportData []SalesReportData

	database.DB.Raw(`
		SELECT s.created_at::date as date,
		       COUNT(s.id) as total_sales,
		       COALESCE(SUM(s.total), 0) as total_revenue,
		       COALESCE(SUM(si.quantity), 0) as total_items
		FROM sales s
		LEFT JOIN sale_items si ON s.id = si.sale_id
		WHERE s.created_at::date BETWEEN ? AND ?
		GROUP BY s.created_at::date
		ORDER BY date ASC
	`, startDate, endDate).Scan(&reportData)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"report":     reportData,
			"start_date": startDate,
			"end_date":   endDate,
		},
	})
}

// UpdateSystemSettings updates system configuration
func UpdateSystemSettings(c *gin.Context) {
	var settings struct {
		StoreName         string  `json:"store_name"`
		StoreAddress      string  `json:"store_address"`
		StorePhone        string  `json:"store_phone"`
		LowStockThreshold int     `json:"low_stock_threshold"`
		TaxRate           float64 `json:"tax_rate"`
	}

	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid input data",
		})
		return
	}

	// In a real application, you would save these to a settings table
	// For now, we'll just return success
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Settings updated successfully",
		"data":    settings,
	})
}

// BackupDatabase creates a database backup using pg_dump
func BackupDatabase(c *gin.Context) {
	// Create backups directory if it doesn't exist
	backupDir := "backups"
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create backup directory: " + err.Error(),
		})
		return
	}

	// Generate backup filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	backupFilename := fmt.Sprintf("inventory_backup_%s.sql", timestamp)
	backupPath := filepath.Join(backupDir, backupFilename)

	// Get database connection info from environment
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "postgres"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "inventory_system"
	}

	// Prepare pg_dump command
	cmd := exec.Command("pg_dump",
		"-h", dbHost,
		"-p", dbPort,
		"-U", dbUser,
		"-d", dbName,
		"-f", backupPath,
		"--no-password",
		"--verbose",
		"--clean",
		"--if-exists",
		"--create",
	)

	// Set PGPASSWORD environment variable for authentication
	cmd.Env = append(os.Environ(), "PGPASSWORD="+dbPassword)

	// Execute backup command
	if err := cmd.Run(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to create database backup: " + err.Error(),
		})
		return
	}

	// Get file size
	fileInfo, err := os.Stat(backupPath)
	var fileSize string
	if err != nil {
		fileSize = "Unknown"
	} else {
		size := fileInfo.Size()
		if size < 1024 {
			fileSize = fmt.Sprintf("%d B", size)
		} else if size < 1024*1024 {
			fileSize = fmt.Sprintf("%.1f KB", float64(size)/1024)
		} else {
			fileSize = fmt.Sprintf("%.1f MB", float64(size)/(1024*1024))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Database backup created successfully",
		"data": gin.H{
			"filename":   backupFilename,
			"path":       backupPath,
			"created_at": time.Now(),
			"size":       fileSize,
		},
	})
}

// RestoreDatabase restores database from backup using psql
func RestoreDatabase(c *gin.Context) {
	filename := c.PostForm("filename")

	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Backup filename is required",
		})
		return
	}

	// Check if backup file exists
	backupDir := "backups"
	backupPath := filepath.Join(backupDir, filename)
	
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Backup file not found: " + filename,
		})
		return
	}

	// Get database connection info from environment
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "postgres"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "inventory_system"
	}

	// Prepare psql command to restore database
	cmd := exec.Command("psql",
		"-h", dbHost,
		"-p", dbPort,
		"-U", dbUser,
		"-d", "postgres", // Connect to postgres database first
		"-f", backupPath,
		"--no-password",
		"-v", "ON_ERROR_STOP=1",
	)

	// Set PGPASSWORD environment variable for authentication
	cmd.Env = append(os.Environ(), "PGPASSWORD="+dbPassword)

	// Execute restore command
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to restore database: " + err.Error(),
			"details": string(output),
		})
		return
	}

	// After successful restore, we need to reconnect the application to the database
	// Close existing connections
	sqlDB, err := database.DB.DB()
	if err == nil {
		sqlDB.Close()
	}

	// Reinitialize database connection
	database.InitDatabase()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Database restored successfully from " + filename,
		"data": gin.H{
			"filename":    filename,
			"restored_at": time.Now(),
		},
	})
}

// ListBackups returns list of available backup files
func ListBackups(c *gin.Context) {
	backupDir := "backups"
	
	// Check if backup directory exists
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    []any{},
			"message": "No backups directory found",
		})
		return
	}

	// Read backup directory
	files, err := os.ReadDir(backupDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to read backup directory: " + err.Error(),
		})
		return
	}

	type BackupFile struct {
		Name       string    `json:"name"`
		Size       string    `json:"size"`
		ModifiedAt time.Time `json:"modified_at"`
		Path       string    `json:"path"`
	}

	var backups []BackupFile

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".sql" {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		var fileSize string
		size := info.Size()
		if size < 1024 {
			fileSize = fmt.Sprintf("%d B", size)
		} else if size < 1024*1024 {
			fileSize = fmt.Sprintf("%.1f KB", float64(size)/1024)
		} else {
			fileSize = fmt.Sprintf("%.1f MB", float64(size)/(1024*1024))
		}

		backups = append(backups, BackupFile{
			Name:       file.Name(),
			Size:       fileSize,
			ModifiedAt: info.ModTime(),
			Path:       filepath.Join(backupDir, file.Name()),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    backups,
	})
}

// DeleteBackup deletes a backup file
func DeleteBackup(c *gin.Context) {
	filename := c.Param("filename")

	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Backup filename is required",
		})
		return
	}

	backupDir := "backups"
	backupPath := filepath.Join(backupDir, filename)

	// Check if backup file exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Backup file not found: " + filename,
		})
		return
	}

	// Delete the backup file
	if err := os.Remove(backupPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to delete backup file: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Backup file deleted successfully: " + filename,
	})
}
