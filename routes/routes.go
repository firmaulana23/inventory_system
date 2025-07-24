package routes

import (
	"inventory_system/handlers"
	"inventory_system/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all API routes
func SetupRoutes() *gin.Engine {
	router := gin.Default()

	// Add CORS middleware
	router.Use(middleware.CORSMiddleware())

	// Serve static files
	router.Static("/static", "./templates")
	router.StaticFile("/login.html", "./templates/login.html")
	router.StaticFile("/admin-dashboard.html", "./templates/admin_dashboard.html")
	router.StaticFile("/dashboard.html", "./templates/employee_dashboard.html")
	// Add routes for new pages
	router.StaticFile("/manage_products", "./templates/manage_products.html")
	router.StaticFile("/sales_history", "./templates/sales_history.html")
	router.StaticFile("/stock_movements", "./templates/stock_movements.html")
	router.StaticFile("/suppliers", "./templates/suppliers.html")
	router.StaticFile("/purchase_orders", "./templates/purchase_orders.html")
	router.StaticFile("/users", "./templates/users.html")
	router.StaticFile("/activity_logs", "./templates/activity_logs.html")
	router.StaticFile("/settings", "./templates/settings.html")
	router.StaticFile("/contact_supplier", "./templates/contact_supplier.html")
	router.StaticFile("/404.html", "./templates/404.html")

	// Redirect root to login
	router.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/login.html")
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// API version 1
	v1 := router.Group("/api/v1")

	// Public routes (no authentication required)
	auth := v1.Group("/auth")
	{
		auth.POST("/login", handlers.Login)
		auth.POST("/register", handlers.Register)
	}

	// Protected routes (authentication required)
	protected := v1.Group("/")
	protected.Use(middleware.AuthMiddleware())
	protected.Use(middleware.ActivityLogger()) // Log user activities
	{
		// User profile
		protected.GET("/profile", handlers.GetProfile)

		// Products (all authenticated users can view)
		products := protected.Group("/products")
		{
			products.GET("", handlers.GetProducts)
			products.GET("/search", handlers.SearchProducts)
			products.GET("/categories", handlers.GetProductCategories)
			products.GET("/low-stock", handlers.GetLowStockProducts)
			products.GET("/:id", handlers.GetProduct)
			products.GET("/:id/suppliers", handlers.GetProductSuppliers)
		}

		// POS (all authenticated users can make sales)
		pos := protected.Group("/pos")
		{
			pos.POST("/sales", handlers.CreateSale)
			pos.GET("/sales", handlers.GetSales)
			pos.GET("/sales/:id", handlers.GetSale)
			pos.GET("/sales/:id/payments", handlers.GetSalePayments)
			pos.GET("/sales/summary", handlers.GetSalesSummary)
			pos.GET("/sales/export", handlers.ExportSales)
			pos.GET("/sales/export-excel", handlers.ExportSalesExcel)
			pos.GET("/reports", handlers.GetSalesReport)
			pos.GET("/sales/overdue", handlers.GetOverdueSales)
		}

		// Stock movements (view only for employees)
		protected.GET("/stock-movements", handlers.GetStockMovements)
		
		// Company profile for invoices (accessible to all authenticated users)
		protected.GET("/company-profile/invoice", handlers.GetCompanyProfileForInvoice)
		
		// Suppliers (view only for employees)
		suppliers := protected.Group("/suppliers")
		{
			suppliers.GET("", handlers.GetSuppliers)
			suppliers.GET("/:id", handlers.GetSupplier)
			suppliers.GET("/search", handlers.SearchSuppliers)
		}

		// Manager level routes
		manager := protected.Group("/")
		manager.Use(middleware.ManagerMiddleware())
		{
			// Product management
			manager.POST("/products", handlers.CreateProduct)
			manager.PUT("/products/:id", handlers.UpdateProduct)
			manager.POST("/products/:id/suppliers", handlers.AddProductSupplier)
			manager.PUT("/products/:id/suppliers/:supplier_id", handlers.UpdateProductSupplier)
			manager.DELETE("/products/:id/suppliers/:supplier_id", handlers.RemoveProductSupplier)
			manager.POST("/products/:id/suppliers/:supplier_id/adjust-stock", handlers.AdjustSupplierStock)
			
			// Supplier management
			manager.POST("/suppliers", handlers.CreateSupplier)
			manager.PUT("/suppliers/:id", handlers.UpdateSupplier)
			manager.DELETE("/suppliers/:id", handlers.DeleteSupplier)

			// Sales management
			manager.PUT("/pos/sales/:id/void", handlers.VoidSale)
			manager.DELETE("/pos/sales/:id", handlers.DeleteSale)
			manager.POST("/pos/sales/:id/payment", handlers.RecordSalePayment)

			// Purchase order management
			manager.POST("/purchase-orders", handlers.CreatePurchaseOrder)
			manager.GET("/purchase-orders", handlers.GetPurchaseOrders)
			manager.GET("/purchase-orders/:id", handlers.GetPurchaseOrder)
			manager.PUT("/purchase-orders/:id", handlers.UpdatePurchaseOrder)
			manager.GET("/purchase-orders/overdue", handlers.GetOverduePurchaseOrders)
			manager.GET("/purchase-orders/summary", handlers.GetPurchaseOrdersSummary)
			manager.POST("/purchase-orders/:id/payment", handlers.RecordPurchasePayment)
			manager.GET("/purchase-orders/:id/payments", handlers.GetPurchasePaymentHistory)
			manager.DELETE("/purchase-orders/:id", handlers.DeletePurchaseOrder)

		}

		// Admin only routes
		admin := protected.Group("/admin")
		admin.Use(middleware.AdminMiddleware())
		{
			// Dashboard
			admin.GET("/dashboard/stats", handlers.GetDashboardStats)
			admin.GET("/dashboard/sales-report", handlers.GetAdminSalesReport)
			admin.GET("/system/logs", handlers.GetSystemLogs)
			admin.GET("/users/:id/activity", handlers.GetUserActivity)

			// System management
			admin.PUT("/system/settings", handlers.UpdateSystemSettings)
			
			// Company profile management
			admin.GET("/company-profile", handlers.GetCompanyProfile)
			admin.POST("/company-profile", handlers.CreateOrUpdateCompanyProfile)
			admin.PUT("/company-profile", handlers.CreateOrUpdateCompanyProfile)
			admin.PUT("/company-profile/logo", handlers.UpdateCompanyLogo)
			admin.DELETE("/company-profile/:id", handlers.DeleteCompanyProfile)
			admin.POST("/system/backup", handlers.BackupDatabase)
			admin.POST("/system/restore", handlers.RestoreDatabase)
			admin.GET("/system/backups", handlers.ListBackups)
			admin.DELETE("/system/backups/:filename", handlers.DeleteBackup)

			// User management
			admin.GET("/users", handlers.GetUsers)
			admin.POST("/users", handlers.Register) // Admin can create users
			admin.GET("/users/:id", handlers.GetUser)
			admin.PUT("/users/:id", handlers.UpdateUser)
			admin.DELETE("/users/:id", handlers.DeleteUser)

			// Product management (admin can delete)
			admin.DELETE("/products/:id", handlers.DeleteProduct)
		}
	}

	return router
}
