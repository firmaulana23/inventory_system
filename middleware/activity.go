package middleware

import (
	"inventory_system/database"
	"inventory_system/models"
	"time"

	"github.com/gin-gonic/gin"
)

// ActivityLogger middleware logs user activities
func ActivityLogger() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Process request first
		c.Next()

		// Only log successful requests and specific actions
		if c.Writer.Status() >= 400 {
			return
		}

		// Get user from context
		userInterface, exists := c.Get("user")
		if !exists {
			return
		}

		user, ok := userInterface.(*models.User)
		if !ok {
			return
		}

		// Determine action based on method and path
		action := getActionFromRequest(c)
		if action == "" {
			return // Skip logging for read-only operations
		}

		// Create activity log
		log := models.ActivityLog{
			UserID:     user.ID,
			Action:     action,
			Resource:   c.Request.URL.Path,
			ResourceID: getResourceID(c),
			Details:    getActionDetails(c),
			IPAddress:  c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
			CreatedAt:  time.Now(),
		}

		// Save to database (non-blocking)
		go func() {
			database.DB.Create(&log)
		}()
	})
}

// getActionFromRequest determines the action based on HTTP method and path
func getActionFromRequest(c *gin.Context) string {
	method := c.Request.Method
	path := c.Request.URL.Path

	switch method {
	case "POST":
		if contains(path, "/auth/login") {
			return "login"
		} else if contains(path, "/auth/register") {
			return "register"
		} else if contains(path, "/products") {
			return "create_product"
		} else if contains(path, "/pos/sales") {
			return "create_sale"
		} else if contains(path, "/users") {
			return "create_user"
		} else if contains(path, "/adjust-stock") {
			return "adjust_stock"
		} else if contains(path, "/backup") {
			return "backup_database"
		}
	case "PUT":
		if contains(path, "/products") {
			return "update_product"
		} else if contains(path, "/users") {
			return "update_user"
		} else if contains(path, "/settings") {
			return "update_settings"
		} else if contains(path, "/void") {
			return "void_sale"
		}
	case "DELETE":
		if contains(path, "/products") {
			return "delete_product"
		} else if contains(path, "/users") {
			return "delete_user"
		}
	}

	return ""
}

// getResourceID extracts resource ID from URL parameters
func getResourceID(c *gin.Context) uint {
	if id := c.Param("id"); id != "" {
		// Convert string to uint (simple conversion, should add proper error handling)
		var resourceID uint
		if parsed := parseUint(id); parsed > 0 {
			resourceID = parsed
		}
		return resourceID
	}
	return 0
}

// getActionDetails creates a description of the action
func getActionDetails(c *gin.Context) string {
	method := c.Request.Method
	path := c.Request.URL.Path

	if method == "POST" && contains(path, "/pos/sales") {
		return "Created new sale transaction"
	} else if method == "PUT" && contains(path, "/products") {
		return "Updated product information"
	} else if method == "DELETE" && contains(path, "/users") {
		return "Deleted user account"
	} else if contains(path, "/adjust-stock") {
		return "Adjusted product stock levels"
	}

	return method + " " + path
}

// Helper functions
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr)))
}

func parseUint(s string) uint {
	var result uint
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			result = result*10 + uint(ch-'0')
		} else {
			return 0
		}
	}
	return result
}
