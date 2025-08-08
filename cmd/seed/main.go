package main

import (
	"fmt"
	"log"
	"math/rand"
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
	fmt.Println("Seeding comprehensive sample data with 50 records for each model...")

	// Seed Users (50 records)
	seedUsers()

	// Seed Suppliers (50 records)
	seedSuppliers()

	// Seed Products (50 records)
	seedProducts()

	// Seed ProductSupplier relationships
	seedProductSuppliers()

	// Seed Company Profile
	seedCompanyProfile()

	// Seed Purchase Orders (50 records)
	seedPurchaseOrders()

	// Seed Sales (50 records)
	seedSales()

	// Seed Sale Items
	seedSaleItems()

	// Seed Purchase Order Items
	seedPurchaseOrderItems()

	// Seed Returns (30 records)
	seedReturns()

	// Seed Exchanges (20 records)
	seedExchanges()

	// Seed Stock Movements (100 records)
	seedStockMovements()

	// Seed Activity Logs (100 records)
	seedActivityLogs()

	// Seed Sale Payments (25 records)
	seedSalePayments()

	// Seed Purchase Payments (25 records)
	seedPurchasePayments()

	fmt.Println("Comprehensive data seeding completed!")
}

func seedUsers() {
	var count int64
	database.DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		fmt.Printf("Users already exist (%d users found). Skipping.\n", count)
		return
	}

	users := []models.User{
		// Admin users (5)
		{Email: "admin@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "System Administrator", Role: "admin", IsActive: true},
		{Email: "admin2@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Senior Admin", Role: "admin", IsActive: true},
		{Email: "admin3@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "IT Administrator", Role: "admin", IsActive: true},
		{Email: "admin4@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Security Admin", Role: "admin", IsActive: true},
		{Email: "admin5@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Database Admin", Role: "admin", IsActive: true},
		
		// Manager users (15)
		{Email: "manager@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Store Manager", Role: "manager", IsActive: true},
		{Email: "manager2@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Operations Manager", Role: "manager", IsActive: true},
		{Email: "manager3@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Sales Manager", Role: "manager", IsActive: true},
		{Email: "manager4@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Warehouse Manager", Role: "manager", IsActive: true},
		{Email: "manager5@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Inventory Manager", Role: "manager", IsActive: true},
		{Email: "manager6@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Floor Manager", Role: "manager", IsActive: true},
		{Email: "manager7@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Shift Manager", Role: "manager", IsActive: true},
		{Email: "manager8@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Regional Manager", Role: "manager", IsActive: true},
		{Email: "manager9@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Assistant Manager", Role: "manager", IsActive: true},
		{Email: "manager10@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Department Manager", Role: "manager", IsActive: true},
		{Email: "manager11@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Team Lead Manager", Role: "manager", IsActive: true},
		{Email: "manager12@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "General Manager", Role: "manager", IsActive: true},
		{Email: "manager13@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Branch Manager", Role: "manager", IsActive: true},
		{Email: "manager14@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Area Manager", Role: "manager", IsActive: true},
		{Email: "manager15@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Project Manager", Role: "manager", IsActive: true},
		
		// Employee users (30)
		{Email: "employee@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "John Smith", Role: "employee", IsActive: true},
		{Email: "employee2@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Sarah Johnson", Role: "employee", IsActive: true},
		{Email: "employee3@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Mike Wilson", Role: "employee", IsActive: true},
		{Email: "employee4@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Lisa Chen", Role: "employee", IsActive: true},
		{Email: "employee5@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "David Brown", Role: "employee", IsActive: true},
		{Email: "employee6@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Emma Davis", Role: "employee", IsActive: true},
		{Email: "employee7@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Carlos Rodriguez", Role: "employee", IsActive: true},
		{Email: "employee8@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Jennifer White", Role: "employee", IsActive: true},
		{Email: "employee9@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Michael Brown", Role: "employee", IsActive: true},
		{Email: "employee10@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Ryan Kim", Role: "employee", IsActive: true},
		{Email: "employee11@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Ashley Taylor", Role: "employee", IsActive: true},
		{Email: "employee12@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Kevin Martinez", Role: "employee", IsActive: true},
		{Email: "employee13@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Amanda Wilson", Role: "employee", IsActive: true},
		{Email: "employee14@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Chris Garcia", Role: "employee", IsActive: true},
		{Email: "employee15@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Laura Anderson", Role: "employee", IsActive: true},
		{Email: "employee16@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Daniel Lee", Role: "employee", IsActive: true},
		{Email: "employee17@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Jessica Thompson", Role: "employee", IsActive: true},
		{Email: "employee18@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Thomas Clark", Role: "employee", IsActive: true},
		{Email: "employee19@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Nicole Walker", Role: "employee", IsActive: true},
		{Email: "employee20@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "James Miller", Role: "employee", IsActive: true},
		{Email: "employee21@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Michelle Lopez", Role: "employee", IsActive: true},
		{Email: "employee22@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Robert Jackson", Role: "employee", IsActive: true},
		{Email: "employee23@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Stephanie Moore", Role: "employee", IsActive: true},
		{Email: "employee24@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Matthew Harris", Role: "employee", IsActive: true},
		{Email: "employee25@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Kimberly Young", Role: "employee", IsActive: true},
		{Email: "employee26@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Andrew Scott", Role: "employee", IsActive: true},
		{Email: "employee27@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Rachel Green", Role: "employee", IsActive: true},
		{Email: "employee28@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Brandon King", Role: "employee", IsActive: true},
		{Email: "employee29@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Amy Wright", Role: "employee", IsActive: true},
		{Email: "employee30@inventory.com", Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", Name: "Jason Hill", Role: "employee", IsActive: true},
	}

	for _, user := range users {
		result := database.DB.Create(&user)
		if result.Error != nil {
			log.Printf("Error creating user %s: %v", user.Email, result.Error)
		} else {
			fmt.Printf("Created user: %s (%s)\n", user.Name, user.Role)
		}
	}
}

func seedSuppliers() {
	var count int64
	database.DB.Model(&models.Supplier{}).Count(&count)
	if count > 0 {
		fmt.Printf("Suppliers already exist (%d suppliers found). Skipping.\n", count)
		return
	}

	suppliers := []models.Supplier{
		{Name: "Tech Supplies Inc", Category: "Electronics", Email: "orders@techsupplies.com", Phone: "+1-555-0101", Address: "123 Tech Street, Silicon Valley, CA 94000", ContactPerson: "John Smith", Website: "https://techsupplies.com", IsActive: true},
		{Name: "Office Furniture Co", Category: "Furniture", Email: "sales@officefurniture.com", Phone: "+1-555-0202", Address: "456 Business Ave, New York, NY 10001", ContactPerson: "Sarah Johnson", Website: "https://officefurniture.com", IsActive: true},
		{Name: "Paper Plus", Category: "Stationery", Email: "support@paperplus.com", Phone: "+1-555-0303", Address: "789 Paper Mill Rd, Portland, OR 97201", ContactPerson: "Mike Wilson", Website: "https://paperplus.com", IsActive: true},
		{Name: "Cable World", Category: "Electronics", Email: "info@cableworld.com", Phone: "+1-555-0404", Address: "321 Cable Ave, Austin, TX 73301", ContactPerson: "David Lee", Website: "https://cableworld.com", IsActive: true},
		{Name: "Audio Tech", Category: "Electronics", Email: "sales@audiotech.com", Phone: "+1-555-0505", Address: "654 Sound St, Nashville, TN 37201", ContactPerson: "Emma Davis", Website: "https://audiotech.com", IsActive: true},
		{Name: "Bright Lights Ltd", Category: "Lighting", Email: "orders@brightlights.com", Phone: "+1-555-0606", Address: "987 Light Blvd, Miami, FL 33101", ContactPerson: "Carlos Rodriguez", Website: "https://brightlights.com", IsActive: true},
		{Name: "Kitchen Supplies", Category: "Kitchen", Email: "support@kitchensupplies.com", Phone: "+1-555-0707", Address: "147 Kitchen Rd, Chicago, IL 60601", ContactPerson: "Lisa Chen", Website: "https://kitchensupplies.com", IsActive: true},
		{Name: "Phone Accessories Inc", Category: "Electronics", Email: "sales@phoneaccessories.com", Phone: "+1-555-0808", Address: "258 Mobile St, San Francisco, CA 94102", ContactPerson: "Ryan Kim", Website: "https://phoneaccessories.com", IsActive: true},
		{Name: "Office Organizers", Category: "Office Supplies", Email: "info@officeorganizers.com", Phone: "+1-555-0909", Address: "369 Organize Ave, Boston, MA 02101", ContactPerson: "Jennifer White", Website: "https://officeorganizers.com", IsActive: true},
		{Name: "Hydration Station", Category: "Lifestyle", Email: "orders@hydrationstation.com", Phone: "+1-555-1010", Address: "741 Water Way, Denver, CO 80201", ContactPerson: "Michael Brown", Website: "https://hydrationstation.com", IsActive: true},
		{Name: "Computer Components Corp", Category: "Electronics", Email: "sales@compcomponents.com", Phone: "+1-555-1111", Address: "852 CPU Blvd, Seattle, WA 98101", ContactPerson: "Ashley Taylor", Website: "https://compcomponents.com", IsActive: true},
		{Name: "Sports Equipment Pro", Category: "Sports", Email: "info@sportsequipment.com", Phone: "+1-555-1212", Address: "963 Sports Ave, Dallas, TX 75201", ContactPerson: "Kevin Martinez", Website: "https://sportsequipment.com", IsActive: true},
		{Name: "Medical Supply Direct", Category: "Medical", Email: "orders@medicalsupply.com", Phone: "+1-555-1313", Address: "147 Health St, Atlanta, GA 30301", ContactPerson: "Amanda Wilson", Website: "https://medicalsupply.com", IsActive: true},
		{Name: "Auto Parts Warehouse", Category: "Automotive", Email: "sales@autoparts.com", Phone: "+1-555-1414", Address: "258 Motor Dr, Detroit, MI 48201", ContactPerson: "Chris Garcia", Website: "https://autoparts.com", IsActive: true},
		{Name: "Garden Tools & More", Category: "Garden", Email: "support@gardentools.com", Phone: "+1-555-1515", Address: "369 Garden Way, Phoenix, AZ 85001", ContactPerson: "Laura Anderson", Website: "https://gardentools.com", IsActive: true},
		{Name: "Fashion Accessories Ltd", Category: "Fashion", Email: "info@fashionaccessories.com", Phone: "+1-555-1616", Address: "741 Fashion St, Los Angeles, CA 90001", ContactPerson: "Daniel Lee", Website: "https://fashionaccessories.com", IsActive: true},
		{Name: "Home Decor Central", Category: "Home Decor", Email: "orders@homedecor.com", Phone: "+1-555-1717", Address: "852 Decor Blvd, San Diego, CA 92101", ContactPerson: "Jessica Thompson", Website: "https://homedecor.com", IsActive: true},
		{Name: "Industrial Equipment Co", Category: "Industrial", Email: "sales@industrialequip.com", Phone: "+1-555-1818", Address: "963 Industry Ave, Houston, TX 77001", ContactPerson: "Thomas Clark", Website: "https://industrialequip.com", IsActive: true},
		{Name: "Pet Supplies Plus", Category: "Pets", Email: "support@petsupplies.com", Phone: "+1-555-1919", Address: "147 Pet St, Philadelphia, PA 19101", ContactPerson: "Nicole Walker", Website: "https://petsupplies.com", IsActive: true},
		{Name: "Craft Materials Inc", Category: "Crafts", Email: "info@craftmaterials.com", Phone: "+1-555-2020", Address: "258 Craft Ave, Minneapolis, MN 55401", ContactPerson: "James Miller", Website: "https://craftmaterials.com", IsActive: true},
		{Name: "Electronic Components", Category: "Electronics", Email: "orders@electroniccomp.com", Phone: "+1-555-2121", Address: "369 Circuit Dr, Portland, OR 97201", ContactPerson: "Michelle Lopez", Website: "https://electroniccomp.com", IsActive: true},
		{Name: "Packaging Solutions", Category: "Packaging", Email: "sales@packagingsol.com", Phone: "+1-555-2222", Address: "741 Package Way, Tampa, FL 33601", ContactPerson: "Robert Jackson", Website: "https://packagingsol.com", IsActive: true},
		{Name: "Safety Equipment Co", Category: "Safety", Email: "support@safetyequip.com", Phone: "+1-555-2323", Address: "852 Safety Blvd, Cleveland, OH 44101", ContactPerson: "Stephanie Moore", Website: "https://safetyequip.com", IsActive: true},
		{Name: "Cleaning Supplies Direct", Category: "Cleaning", Email: "info@cleaningsupplies.com", Phone: "+1-555-2424", Address: "963 Clean Ave, Las Vegas, NV 89101", ContactPerson: "Matthew Harris", Website: "https://cleaningsupplies.com", IsActive: true},
		{Name: "Musical Instruments Plus", Category: "Music", Email: "orders@musicalinstruments.com", Phone: "+1-555-2525", Address: "147 Music St, New Orleans, LA 70112", ContactPerson: "Kimberly Young", Website: "https://musicalinstruments.com", IsActive: true},
		{Name: "Art Supplies Unlimited", Category: "Art", Email: "sales@artsupplies.com", Phone: "+1-555-2626", Address: "258 Artist Ave, Santa Fe, NM 87501", ContactPerson: "Andrew Scott", Website: "https://artsupplies.com", IsActive: true},
		{Name: "Photography Equipment", Category: "Photography", Email: "support@photoequip.com", Phone: "+1-555-2727", Address: "369 Photo Dr, Nashville, TN 37201", ContactPerson: "Rachel Green", Website: "https://photoequip.com", IsActive: true},
		{Name: "Construction Materials", Category: "Construction", Email: "info@constructionmat.com", Phone: "+1-555-2828", Address: "741 Build Way, Phoenix, AZ 85001", ContactPerson: "Brandon King", Website: "https://constructionmat.com", IsActive: true},
		{Name: "Textile Supplies Co", Category: "Textiles", Email: "orders@textilesupplies.com", Phone: "+1-555-2929", Address: "852 Fabric Blvd, Charlotte, NC 28201", ContactPerson: "Amy Wright", Website: "https://textilesupplies.com", IsActive: true},
		{Name: "Laboratory Equipment", Category: "Laboratory", Email: "sales@labequipment.com", Phone: "+1-555-3030", Address: "963 Lab Ave, Raleigh, NC 27601", ContactPerson: "Jason Hill", Website: "https://labequipment.com", IsActive: true},
		{Name: "Bakery Supplies Inc", Category: "Food Service", Email: "support@bakerysupplies.com", Phone: "+1-555-3131", Address: "147 Bake St, Portland, ME 04101", ContactPerson: "Maria Gonzalez", Website: "https://bakerysupplies.com", IsActive: true},
		{Name: "Hardware Solutions", Category: "Hardware", Email: "info@hardwaresol.com", Phone: "+1-555-3232", Address: "258 Tool Ave, Milwaukee, WI 53201", ContactPerson: "Steven Davis", Website: "https://hardwaresol.com", IsActive: true},
		{Name: "Gaming Equipment Co", Category: "Gaming", Email: "orders@gamingequip.com", Phone: "+1-555-3333", Address: "369 Game Dr, Austin, TX 78701", ContactPerson: "Lisa Rodriguez", Website: "https://gamingequip.com", IsActive: true},
		{Name: "Travel Accessories", Category: "Travel", Email: "sales@travelaccessories.com", Phone: "+1-555-3434", Address: "741 Travel Way, Orlando, FL 32801", ContactPerson: "Mark Thompson", Website: "https://travelaccessories.com", IsActive: true},
		{Name: "Baby Products Plus", Category: "Baby", Email: "support@babyproducts.com", Phone: "+1-555-3535", Address: "852 Baby Blvd, San Antonio, TX 78201", ContactPerson: "Jennifer Lopez", Website: "https://babyproducts.com", IsActive: true},
		{Name: "Senior Care Supplies", Category: "Healthcare", Email: "info@seniorcare.com", Phone: "+1-555-3636", Address: "963 Care Ave, Jacksonville, FL 32201", ContactPerson: "David Wilson", Website: "https://seniorcare.com", IsActive: true},
		{Name: "Educational Materials", Category: "Education", Email: "orders@edumatrerials.com", Phone: "+1-555-3737", Address: "147 Learn St, Columbus, OH 43201", ContactPerson: "Sarah Mitchell", Website: "https://edumaterials.com", IsActive: true},
		{Name: "Fitness Equipment Pro", Category: "Fitness", Email: "sales@fitnessequip.com", Phone: "+1-555-3838", Address: "258 Fit Ave, Indianapolis, IN 46201", ContactPerson: "Michael Garcia", Website: "https://fitnessequip.com", IsActive: true},
		{Name: "Beauty Supply Central", Category: "Beauty", Email: "support@beautysupply.com", Phone: "+1-555-3939", Address: "369 Beauty Dr, Memphis, TN 38101", ContactPerson: "Amanda Martinez", Website: "https://beautysupply.com", IsActive: true},
		{Name: "Jewelry Components", Category: "Jewelry", Email: "info@jewelrycomp.com", Phone: "+1-555-4040", Address: "741 Gem Way, Baltimore, MD 21201", ContactPerson: "Christopher Lee", Website: "https://jewelrycomp.com", IsActive: true},
		{Name: "Event Planning Supplies", Category: "Events", Email: "orders@eventplanning.com", Phone: "+1-555-4141", Address: "852 Event Blvd, Washington, DC 20001", ContactPerson: "Nicole Johnson", Website: "https://eventplanning.com", IsActive: true},
		{Name: "Outdoor Gear Co", Category: "Outdoor", Email: "sales@outdoorgear.com", Phone: "+1-555-4242", Address: "963 Outdoor Ave, Denver, CO 80201", ContactPerson: "James Wilson", Website: "https://outdoorgear.com", IsActive: true},
		{Name: "Technology Solutions", Category: "Technology", Email: "support@techsolutions.com", Phone: "+1-555-4343", Address: "147 Tech St, Seattle, WA 98101", ContactPerson: "Michelle Brown", Website: "https://techsolutions.com", IsActive: true},
		{Name: "Restaurant Equipment", Category: "Food Service", Email: "info@restaurantequip.com", Phone: "+1-555-4444", Address: "258 Restaurant Ave, Las Vegas, NV 89101", ContactPerson: "Robert Miller", Website: "https://restaurantequip.com", IsActive: true},
		{Name: "Marine Supplies Inc", Category: "Marine", Email: "orders@marinesupplies.com", Phone: "+1-555-4545", Address: "369 Marine Dr, Miami, FL 33101", ContactPerson: "Stephanie Davis", Website: "https://marinesupplies.com", IsActive: true},
		{Name: "Agricultural Equipment", Category: "Agriculture", Email: "sales@agequipment.com", Phone: "+1-555-4646", Address: "741 Farm Way, Des Moines, IA 50301", ContactPerson: "Matthew Anderson", Website: "https://agequipment.com", IsActive: true},
		{Name: "Renewable Energy Co", Category: "Energy", Email: "support@renewableenergy.com", Phone: "+1-555-4747", Address: "852 Solar Blvd, San Diego, CA 92101", ContactPerson: "Kimberly Taylor", Website: "https://renewableenergy.com", IsActive: true},
		{Name: "Security Systems Plus", Category: "Security", Email: "info@securitysystems.com", Phone: "+1-555-4848", Address: "963 Security Ave, Atlanta, GA 30301", ContactPerson: "Andrew White", Website: "https://securitysystems.com", IsActive: true},
		{Name: "Printing Supplies Direct", Category: "Printing", Email: "orders@printingsupplies.com", Phone: "+1-555-4949", Address: "147 Print St, Boston, MA 02101", ContactPerson: "Rachel Martinez", Website: "https://printingsupplies.com", IsActive: true},
		{Name: "HVAC Solutions", Category: "HVAC", Email: "sales@hvacsolutions.com", Phone: "+1-555-5050", Address: "258 Climate Ave, Chicago, IL 60601", ContactPerson: "Brandon Thompson", Website: "https://hvacsolutions.com", IsActive: true},
	}

	for _, supplier := range suppliers {
		result := database.DB.Create(&supplier)
		if result.Error != nil {
			log.Printf("Error creating supplier %s: %v", supplier.Name, result.Error)
		} else {
			fmt.Printf("Created supplier: %s (%s)\n", supplier.Name, supplier.Category)
		}
	}
}

func seedProducts() {
	var count int64
	database.DB.Model(&models.Product{}).Count(&count)
	if count > 0 {
		fmt.Printf("Products already exist (%d products found). Skipping.\n", count)
		return
	}

	products := []models.Product{
		{Name: "Wireless Mouse", SKU: "WM001", Description: "Bluetooth wireless mouse with optical sensor", Category: "Electronics", Location: "A1-B2", IsActive: true},
		{Name: "USB Cable Type-C", SKU: "UC001", Description: "High-speed USB Type-C cable 3ft", Category: "Electronics", Location: "A2-C1", IsActive: true},
		{Name: "Office Chair", SKU: "OC001", Description: "Ergonomic office chair with lumbar support", Category: "Furniture", Location: "B1-A1", IsActive: true},
		{Name: "Notebook A4", SKU: "NB001", Description: "Ruled notebook 200 pages", Category: "Stationery", Location: "C1-D2", IsActive: true},
		{Name: "Bluetooth Headphones", SKU: "BH001", Description: "Noise-cancelling wireless headphones", Category: "Electronics", Location: "A1-C3", IsActive: true},
		{Name: "Desk Lamp LED", SKU: "DL001", Description: "Adjustable LED desk lamp with USB charging", Category: "Lighting", Location: "B2-A2", IsActive: true},
		{Name: "Coffee Mug", SKU: "CM001", Description: "Ceramic coffee mug 12oz", Category: "Kitchen", Location: "D1-B1", IsActive: true},
		{Name: "Smartphone Case", SKU: "SC001", Description: "Protective smartphone case with screen protector", Category: "Electronics", Location: "A3-B1", IsActive: true},
		{Name: "Desk Organizer", SKU: "DO001", Description: "Multi-compartment desk organizer", Category: "Office Supplies", Location: "C2-A3", IsActive: true},
		{Name: "Water Bottle", SKU: "WB001", Description: "Stainless steel water bottle 750ml", Category: "Lifestyle", Location: "D2-C1", IsActive: true},
		{Name: "Laptop Stand", SKU: "LS001", Description: "Adjustable aluminum laptop stand", Category: "Electronics", Location: "A1-D1", IsActive: true},
		{Name: "Wireless Keyboard", SKU: "WK001", Description: "Compact wireless keyboard with backlight", Category: "Electronics", Location: "A2-B3", IsActive: true},
		{Name: "Monitor 24inch", SKU: "MN001", Description: "24-inch Full HD IPS monitor", Category: "Electronics", Location: "B1-C2", IsActive: true},
		{Name: "Printer Paper A4", SKU: "PP001", Description: "Premium white printer paper 500 sheets", Category: "Stationery", Location: "C1-A2", IsActive: true},
		{Name: "Webcam HD", SKU: "WC001", Description: "1080p HD webcam with microphone", Category: "Electronics", Location: "A3-D2", IsActive: true},
		{Name: "Standing Desk", SKU: "SD001", Description: "Height adjustable standing desk", Category: "Furniture", Location: "B2-C1", IsActive: true},
		{Name: "Desk Fan", SKU: "DF001", Description: "USB powered desktop fan", Category: "Electronics", Location: "D1-A3", IsActive: true},
		{Name: "Pen Set", SKU: "PS001", Description: "Professional ballpoint pen set", Category: "Stationery", Location: "C2-B2", IsActive: true},
		{Name: "Phone Charger", SKU: "PC001", Description: "Universal smartphone charger cable", Category: "Electronics", Location: "A1-B3", IsActive: true},
		{Name: "Bookshelf", SKU: "BS001", Description: "5-tier wooden bookshelf", Category: "Furniture", Location: "B1-D3", IsActive: true},
		{Name: "Stapler", SKU: "ST001", Description: "Heavy-duty office stapler", Category: "Office Supplies", Location: "C1-C3", IsActive: true},
		{Name: "Scissors", SKU: "SS001", Description: "Professional office scissors", Category: "Office Supplies", Location: "C2-D1", IsActive: true},
		{Name: "Calculator", SKU: "CA001", Description: "Scientific calculator with display", Category: "Office Supplies", Location: "D1-C2", IsActive: true},
		{Name: "File Folder", SKU: "FF001", Description: "Manila file folders pack of 50", Category: "Stationery", Location: "C1-B1", IsActive: true},
		{Name: "Whiteboard", SKU: "WH001", Description: "Magnetic whiteboard with markers", Category: "Office Supplies", Location: "B2-D2", IsActive: true},
		{Name: "Power Bank", SKU: "PB001", Description: "10000mAh portable power bank", Category: "Electronics", Location: "A2-D3", IsActive: true},
		{Name: "Ethernet Cable", SKU: "EC001", Description: "Cat6 ethernet cable 10ft", Category: "Electronics", Location: "A3-C1", IsActive: true},
		{Name: "Desk Pad", SKU: "DP001", Description: "Large desk pad mouse mat", Category: "Office Supplies", Location: "D2-A1", IsActive: true},
		{Name: "Paper Clips", SKU: "PCL001", Description: "Standard paper clips box of 100", Category: "Stationery", Location: "C2-C2", IsActive: true},
		{Name: "Tape Dispenser", SKU: "TD001", Description: "Desktop tape dispenser with tape", Category: "Office Supplies", Location: "D1-D1", IsActive: true},
		{Name: "Wireless Speaker", SKU: "WS001", Description: "Portable Bluetooth speaker", Category: "Electronics", Location: "A1-A1", IsActive: true},
		{Name: "Desk Calendar", SKU: "DC001", Description: "Monthly desk calendar planner", Category: "Stationery", Location: "C1-D1", IsActive: true},
		{Name: "Storage Box", SKU: "SB001", Description: "Plastic storage container with lid", Category: "Office Supplies", Location: "B1-B2", IsActive: true},
		{Name: "Label Maker", SKU: "LM001", Description: "Electronic label maker with tape", Category: "Office Supplies", Location: "D2-B3", IsActive: true},
		{Name: "Hole Punch", SKU: "HP001", Description: "3-hole punch for documents", Category: "Office Supplies", Location: "C2-A1", IsActive: true},
		{Name: "Desk Clock", SKU: "DCL001", Description: "Digital desk clock with alarm", Category: "Office Supplies", Location: "D1-B2", IsActive: true},
		{Name: "USB Hub", SKU: "UH001", Description: "4-port USB 3.0 hub", Category: "Electronics", Location: "A2-A2", IsActive: true},
		{Name: "Mouse Pad", SKU: "MP001", Description: "Ergonomic mouse pad with wrist rest", Category: "Electronics", Location: "A3-B2", IsActive: true},
		{Name: "Document Tray", SKU: "DT001", Description: "Stackable document tray set", Category: "Office Supplies", Location: "B2-B1", IsActive: true},
		{Name: "Shredder", SKU: "SH001", Description: "Personal paper shredder", Category: "Office Supplies", Location: "D2-D2", IsActive: true},
		{Name: "Printer Ink", SKU: "PI001", Description: "Black printer ink cartridge", Category: "Stationery", Location: "C1-A1", IsActive: true},
		{Name: "Rubber Stamps", SKU: "RS001", Description: "Custom rubber stamp set", Category: "Office Supplies", Location: "C2-B1", IsActive: true},
		{Name: "Binder Clips", SKU: "BC001", Description: "Assorted binder clips pack", Category: "Stationery", Location: "D1-C1", IsActive: true},
		{Name: "Highlighters", SKU: "HL001", Description: "Fluorescent highlighter pen set", Category: "Stationery", Location: "C1-B2", IsActive: true},
		{Name: "Correction Fluid", SKU: "CF001", Description: "White correction fluid pen", Category: "Stationery", Location: "C2-C1", IsActive: true},
		{Name: "Envelope Set", SKU: "ES001", Description: "Business envelope pack of 100", Category: "Stationery", Location: "D1-A1", IsActive: true},
		{Name: "Rubber Bands", SKU: "RB001", Description: "Assorted rubber bands pack", Category: "Office Supplies", Location: "C1-C1", IsActive: true},
		{Name: "Index Cards", SKU: "IC001", Description: "Ruled index cards pack of 200", Category: "Stationery", Location: "C2-D2", IsActive: true},
		{Name: "Push Pins", SKU: "PP002", Description: "Colorful push pins pack", Category: "Office Supplies", Location: "D2-A2", IsActive: true},
		{Name: "Clipboard", SKU: "CB001", Description: "Legal size clipboard with storage", Category: "Office Supplies", Location: "B1-A2", IsActive: true},
	}

	for _, product := range products {
		result := database.DB.Create(&product)
		if result.Error != nil {
			log.Printf("Error creating product %s: %v", product.Name, result.Error)
		} else {
			fmt.Printf("Created product: %s (SKU: %s)\n", product.Name, product.SKU)
		}
	}
}

func seedProductSuppliers() {
	var count int64
	database.DB.Model(&models.ProductSupplier{}).Count(&count)
	if count > 0 {
		fmt.Printf("Product-Supplier relationships already exist (%d relationships found). Skipping.\n", count)
		return
	}

	// Get all products and suppliers
	var products []models.Product
	var suppliers []models.Supplier
	database.DB.Find(&products)
	database.DB.Find(&suppliers)

	if len(products) == 0 || len(suppliers) == 0 {
		log.Println("No products or suppliers found for creating relationships")
		return
	}

	// Create relationships - each product will have 1-3 suppliers
	rand.Seed(time.Now().UnixNano())
	
	for _, product := range products {
		// Random number of suppliers per product (1-3)
		numSuppliers := rand.Intn(3) + 1
		
		// Select random suppliers for this product
		selectedSuppliers := make(map[uint]bool)
		for i := 0; i < numSuppliers && i < len(suppliers); i++ {
			var supplier models.Supplier
			for {
				supplier = suppliers[rand.Intn(len(suppliers))]
				if !selectedSuppliers[supplier.ID] {
					selectedSuppliers[supplier.ID] = true
					break
				}
			}
			
			// Generate realistic pricing based on category
			baseCost := generateBaseCost(product.Category)
			cost := baseCost + float64(rand.Intn(1000))/100 // Add some variation
			price := cost * (1.2 + float64(rand.Intn(80))/100) // 20-100% markup
			stock := rand.Intn(100) + 10 // 10-110 stock
			minStock := rand.Intn(20) + 5 // 5-25 min stock
			
			productSupplier := models.ProductSupplier{
				ProductID:  product.ID,
				SupplierID: supplier.ID,
				Cost:       cost,
				Price:      price,
				Stock:      stock,
				MinStock:   minStock,
				IsActive:   true,
			}
			
			result := database.DB.Create(&productSupplier)
			if result.Error != nil {
				log.Printf("Error creating product-supplier relationship: %v", result.Error)
			} else {
				fmt.Printf("Created relationship: %s -> %s (Stock: %d)\n", product.SKU, supplier.Name, stock)
			}
		}
	}
}

func generateBaseCost(category string) float64 {
	costs := map[string]float64{
		"Electronics":     50.0,
		"Furniture":       200.0,
		"Stationery":      5.0,
		"Office Supplies": 15.0,
		"Lighting":        30.0,
		"Kitchen":         20.0,
		"Lifestyle":       25.0,
	}
	
	if cost, exists := costs[category]; exists {
		return cost
	}
	return 20.0 // default
}

func seedCompanyProfile() {
	var count int64
	database.DB.Model(&models.CompanyProfile{}).Count(&count)
	if count > 0 {
		fmt.Printf("Company profile already exists. Skipping.\n")
		return
	}

	profile := models.CompanyProfile{
		CompanyName:     "Inventory Management System Corp",
		CompanyAddress:  "123 Business Street, Suite 100, Business City, BC 12345",
		CompanyPhone:    "+1-555-INVENTORY",
		CompanyEmail:    "contact@inventorysystem.com",
		CompanyWebsite:  "https://www.inventorysystem.com",
		TaxNumber:       "TAX123456789",
		BusinessLicense: "BL987654321",
		InvoiceFooter:   "Thank you for your business! Payment is due within 30 days.",
		BankAccount:     "Business Bank - Account: 1234567890",
		Currency:        "USD",
		IsActive:        true,
	}

	result := database.DB.Create(&profile)
	if result.Error != nil {
		log.Printf("Error creating company profile: %v", result.Error)
	} else {
		fmt.Printf("Created company profile: %s\n", profile.CompanyName)
	}
}

func seedPurchaseOrders() {
	var count int64
	database.DB.Model(&models.PurchaseOrder{}).Count(&count)
	if count > 0 {
		fmt.Printf("Purchase orders already exist (%d orders found). Skipping.\n", count)
		return
	}

	// Get suppliers for foreign keys
	var suppliers []models.Supplier
	database.DB.Find(&suppliers)
	if len(suppliers) == 0 {
		log.Println("No suppliers found for creating purchase orders")
		return
	}

	rand.Seed(time.Now().UnixNano())
	paymentMethods := []string{"cash", "net7", "net15", "net30", "net60", "net90", "credit"}
	statuses := []string{"pending", "paid", "overdue"}

	for i := 1; i <= 50; i++ {
		supplier := suppliers[rand.Intn(len(suppliers))]
		paymentMethod := paymentMethods[rand.Intn(len(paymentMethods))]
		paymentDays := getPaymentDays(paymentMethod)
		status := statuses[rand.Intn(len(statuses))]
		
		total := float64(rand.Intn(5000)+500) + float64(rand.Intn(100))/100
		downPayment := 0.0
		amountPaid := 0.0
		amountDue := total
		
		if status == "paid" {
			amountPaid = total
			amountDue = 0.0
		} else if rand.Intn(3) == 0 { // 1/3 chance of down payment
			downPayment = total * 0.2 // 20% down payment
			amountPaid = downPayment
			amountDue = total - downPayment
		}
		
		orderDate := time.Now().AddDate(0, 0, -rand.Intn(90)) // Last 90 days
		
		po := models.PurchaseOrder{
			PONumber:      fmt.Sprintf("PO-2025-%03d", i),
			SupplierID:    supplier.ID,
			UserID:        1, // Admin user
			PaymentMethod: paymentMethod,
			PaymentDays:   paymentDays,
			PaymentStatus: status,
			Total:         total,
			DownPayment:   downPayment,
			AmountPaid:    amountPaid,
			AmountDue:     amountDue,
			Notes:         fmt.Sprintf("Purchase order for %s supplies", supplier.Category),
			OrderDate:     orderDate,
		}
		
		// Set due date
		if paymentDays > 0 {
			dueDate := orderDate.AddDate(0, 0, paymentDays)
			po.DueDate = &dueDate
			
			// Update status to overdue if past due
			if dueDate.Before(time.Now()) && po.PaymentStatus == "pending" {
				po.PaymentStatus = "overdue"
			}
		}
		
		// Set paid date for paid orders
		if po.PaymentStatus == "paid" {
			paidDate := orderDate.AddDate(0, 0, rand.Intn(paymentDays+1))
			po.PaidDate = &paidDate
		}

		result := database.DB.Create(&po)
		if result.Error != nil {
			log.Printf("Error creating purchase order %s: %v", po.PONumber, result.Error)
		} else {
			fmt.Printf("Created purchase order: %s (Supplier: %s, Total: $%.2f)\n", po.PONumber, supplier.Name, po.Total)
		}
	}
}

func getPaymentDays(method string) int {
	switch method {
	case "cash":
		return 0
	case "net7":
		return 7
	case "net15":
		return 15
	case "net30":
		return 30
	case "net60":
		return 60
	case "net90":
		return 90
	case "credit":
		return 30
	default:
		return 0
	}
}

func seedSales() {
	var count int64
	database.DB.Model(&models.Sale{}).Count(&count)
	if count > 0 {
		fmt.Printf("Sales already exist (%d sales found). Skipping.\n", count)
		return
	}

	customerNames := []string{
		"John Smith", "Sarah Johnson", "Mike Wilson", "Lisa Chen", "David Brown",
		"Emma Davis", "Carlos Rodriguez", "Jennifer White", "Michael Brown", "Ryan Kim",
		"Ashley Taylor", "Kevin Martinez", "Amanda Wilson", "Chris Garcia", "Laura Anderson",
		"Daniel Lee", "Jessica Thompson", "Thomas Clark", "Nicole Walker", "James Miller",
		"Michelle Lopez", "Robert Jackson", "Stephanie Moore", "Matthew Harris", "Kimberly Young",
		"Andrew Scott", "Rachel Green", "Brandon King", "Amy Wright", "Jason Hill",
		"Maria Gonzalez", "Steven Davis", "Lisa Rodriguez", "Mark Thompson", "Jennifer Lopez",
		"David Wilson", "Sarah Mitchell", "Michael Garcia", "Amanda Martinez", "Christopher Lee",
		"Nicole Johnson", "James Wilson", "Michelle Brown", "Robert Miller", "Stephanie Davis",
		"Matthew Anderson", "Kimberly Taylor", "Andrew White", "Rachel Martinez", "Brandon Thompson",
	}

	paymentMethods := []string{"cash", "card", "transfer", "credit"}
	statuses := []string{"paid", "pending", "overdue"}

	rand.Seed(time.Now().UnixNano())

	for i := 1; i <= 50; i++ {
		customer := customerNames[i-1]
		paymentMethod := paymentMethods[rand.Intn(len(paymentMethods))]
		paymentDays := 0
		if paymentMethod == "credit" {
			paymentDays = []int{15, 30, 60}[rand.Intn(3)]
		}
		
		status := statuses[rand.Intn(len(statuses))]
		if paymentMethod != "credit" {
			status = "paid" // Non-credit sales are typically paid immediately
		}
		
		subtotal := float64(rand.Intn(2000)+100) + float64(rand.Intn(100))/100
		tax := subtotal * 0.08 // 8% tax
		discount := 0.0
		if rand.Intn(4) == 0 { // 25% chance of discount
			discount = subtotal * 0.05 // 5% discount
		}
		total := subtotal + tax - discount
		
		downPayment := 0.0
		amountPaid := total
		amountDue := 0.0
		
		if status != "paid" {
			amountPaid = 0.0
			amountDue = total
			if rand.Intn(3) == 0 { // 1/3 chance of down payment for credit
				downPayment = total * 0.3
				amountPaid = downPayment
				amountDue = total - downPayment
			}
		}
		
		createdAt := time.Now().AddDate(0, 0, -rand.Intn(60)) // Last 60 days
		
		sale := models.Sale{
			SaleNumber:    fmt.Sprintf("SALE-2025-%03d", i),
			UserID:        uint(rand.Intn(30) + 1), // Random user 1-30
			CustomerName:  customer,
			Subtotal:      subtotal,
			Tax:           tax,
			Discount:      discount,
			Total:         total,
			PaymentMethod: paymentMethod,
			PaymentDays:   paymentDays,
			PaymentStatus: status,
			DownPayment:   downPayment,
			AmountPaid:    amountPaid,
			AmountDue:     amountDue,
			CreatedAt:     createdAt,
		}
		
		// Set due date for credit sales
		if paymentDays > 0 {
			dueDate := createdAt.AddDate(0, 0, paymentDays)
			sale.DueDate = &dueDate
			
			if dueDate.Before(time.Now()) && sale.PaymentStatus == "pending" {
				sale.PaymentStatus = "overdue"
			}
		}
		
		// Set paid date for paid sales
		if sale.PaymentStatus == "paid" {
			paidDate := createdAt.Add(time.Hour * time.Duration(rand.Intn(24)))
			sale.PaidDate = &paidDate
		}

		result := database.DB.Create(&sale)
		if result.Error != nil {
			log.Printf("Error creating sale %s: %v", sale.SaleNumber, result.Error)
		} else {
			fmt.Printf("Created sale: %s (Customer: %s, Total: $%.2f)\n", sale.SaleNumber, sale.CustomerName, sale.Total)
		}
	}
}

func seedSaleItems() {
	var count int64
	database.DB.Model(&models.SaleItem{}).Count(&count)
	if count > 0 {
		fmt.Printf("Sale items already exist (%d items found). Skipping.\n", count)
		return
	}

	// Get sales and products with suppliers
	var sales []models.Sale
	var productSuppliers []models.ProductSupplier
	database.DB.Limit(50).Find(&sales)
	database.DB.Preload("Product").Find(&productSuppliers)

	if len(sales) == 0 || len(productSuppliers) == 0 {
		log.Println("No sales or product suppliers found for creating sale items")
		return
	}

	rand.Seed(time.Now().UnixNano())

	for _, sale := range sales {
		// Each sale will have 1-5 items
		numItems := rand.Intn(5) + 1
		
		for j := 0; j < numItems; j++ {
			ps := productSuppliers[rand.Intn(len(productSuppliers))]
			quantity := rand.Intn(5) + 1
			
			saleItem := models.SaleItem{
				SaleID:    sale.ID,
				ProductID: ps.ProductID,
				Quantity:  quantity,
				Price:     ps.Price,
				Cost:      ps.Cost,
				Total:     ps.Price * float64(quantity),
			}
			
			result := database.DB.Create(&saleItem)
			if result.Error != nil {
				log.Printf("Error creating sale item: %v", result.Error)
			} else {
				fmt.Printf("Created sale item: Product %d x%d for Sale %d\n", ps.ProductID, quantity, sale.ID)
			}
		}
	}
}

func seedPurchaseOrderItems() {
	var count int64
	database.DB.Model(&models.PurchaseOrderItem{}).Count(&count)
	if count > 0 {
		fmt.Printf("Purchase order items already exist (%d items found). Skipping.\n", count)
		return
	}

	// Get purchase orders and product suppliers
	var purchaseOrders []models.PurchaseOrder
	var productSuppliers []models.ProductSupplier
	database.DB.Limit(50).Find(&purchaseOrders)
	database.DB.Find(&productSuppliers)

	if len(purchaseOrders) == 0 || len(productSuppliers) == 0 {
		log.Println("No purchase orders or product suppliers found for creating purchase order items")
		return
	}

	rand.Seed(time.Now().UnixNano())

	for _, po := range purchaseOrders {
		// Each PO will have 1-4 items
		numItems := rand.Intn(4) + 1
		
		for j := 0; j < numItems; j++ {
			ps := productSuppliers[rand.Intn(len(productSuppliers))]
			quantityOrdered := rand.Intn(50) + 10
			quantityReceived := quantityOrdered
			if rand.Intn(10) == 0 { // 10% chance of partial delivery
				quantityReceived = rand.Intn(quantityOrdered)
			}
			
			poItem := models.PurchaseOrderItem{
				PurchaseOrderID:    po.ID,
				ProductID:          ps.ProductID,
				ProductSupplierID:  &ps.ID,
				QuantityOrdered:    quantityOrdered,
				QuantityReceived:   quantityReceived,
				UnitCost:           ps.Cost,
				Total:              ps.Cost * float64(quantityOrdered),
			}
			
			result := database.DB.Create(&poItem)
			if result.Error != nil {
				log.Printf("Error creating purchase order item: %v", result.Error)
			} else {
				fmt.Printf("Created PO item: Product %d x%d for PO %d\n", ps.ProductID, quantityOrdered, po.ID)
			}
		}
	}
}

func seedReturns() {
	var count int64
	database.DB.Model(&models.Return{}).Count(&count)
	if count > 0 {
		fmt.Printf("Returns already exist (%d returns found). Skipping.\n", count)
		return
	}

	// Get some sales with items
	var sales []models.Sale
	database.DB.Preload("Items.Product").Limit(30).Find(&sales)

	if len(sales) == 0 {
		log.Println("No sales found for creating returns")
		return
	}

	rand.Seed(time.Now().UnixNano())
	reasons := []string{"Defective product", "Wrong item", "Customer changed mind", "Size issue", "Damaged in shipping"}
	methods := []string{"cash", "card", "transfer", "store_credit"}
	conditions := []string{"good", "damaged", "expired"}

	createdCount := 0
	for _, sale := range sales {
		if len(sale.Items) == 0 || rand.Intn(5) != 0 { // 20% chance of return
			continue
		}

		if createdCount >= 30 {
			break
		}

		item := sale.Items[rand.Intn(len(sale.Items))]
		returnQty := rand.Intn(item.Quantity) + 1

		returnRecord := models.Return{
			ReturnNumber: fmt.Sprintf("RET-2025-%03d", createdCount+1),
			SaleID:       sale.ID,
			UserID:       1,
			Subtotal:     item.Price * float64(returnQty),
			Tax:          item.Price * float64(returnQty) * 0.08,
			Discount:     0,
			Total:        item.Price * float64(returnQty) * 1.08,
			Reason:       reasons[rand.Intn(len(reasons))],
			RefundMethod: methods[rand.Intn(len(methods))],
			RefundAmount: item.Price * float64(returnQty) * 1.08,
			TotalCost:    item.Cost * float64(returnQty),
			ProfitLoss:   -(item.Price - item.Cost) * float64(returnQty),
			CreatedAt:    time.Now().AddDate(0, 0, -rand.Intn(30)),
		}

		result := database.DB.Create(&returnRecord)
		if result.Error != nil {
			log.Printf("Error creating return: %v", result.Error)
		} else {
			// Create return item
			returnItem := models.ReturnItem{
				ReturnID:   returnRecord.ID,
				SaleItemID: item.ID,
				ProductID:  item.ProductID,
				Quantity:   returnQty,
				Price:      item.Price,
				Cost:       item.Cost,
				Total:      item.Price * float64(returnQty),
				TotalCost:  item.Cost * float64(returnQty),
				Condition:  conditions[rand.Intn(len(conditions))],
			}

			database.DB.Create(&returnItem)
			fmt.Printf("Created return: %s (Amount: $%.2f)\n", returnRecord.ReturnNumber, returnRecord.RefundAmount)
			createdCount++
		}
	}
}

func seedExchanges() {
	var count int64
	database.DB.Model(&models.Exchange{}).Count(&count)
	if count > 0 {
		fmt.Printf("Exchanges already exist (%d exchanges found). Skipping.\n", count)
		return
	}

	// Get some sales with items and products
	var sales []models.Sale
	var productSuppliers []models.ProductSupplier
	database.DB.Preload("Items.Product").Limit(20).Find(&sales)
	database.DB.Preload("Product").Find(&productSuppliers)

	if len(sales) == 0 || len(productSuppliers) == 0 {
		log.Println("No sales or products found for creating exchanges")
		return
	}

	rand.Seed(time.Now().UnixNano())
	reasons := []string{"Size exchange", "Color preference", "Model upgrade", "Feature exchange"}

	createdCount := 0
	for _, sale := range sales {
		if len(sale.Items) == 0 || rand.Intn(10) != 0 { // 10% chance of exchange
			continue
		}

		if createdCount >= 20 {
			break
		}

		oldItem := sale.Items[rand.Intn(len(sale.Items))]
		newPS := productSuppliers[rand.Intn(len(productSuppliers))]
		
		oldQty := rand.Intn(oldItem.Quantity) + 1
		newQty := oldQty

		oldValue := oldItem.Price * float64(oldQty)
		newValue := newPS.Price * float64(newQty)
		difference := newValue - oldValue

		exchange := models.Exchange{
			ExchangeNumber: fmt.Sprintf("EXC-2025-%03d", createdCount+1),
			SaleID:         sale.ID,
			UserID:         1,
			Reason:         reasons[rand.Intn(len(reasons))],
			TotalOldValue:  oldValue,
			TotalNewValue:  newValue,
			TotalOldCost:   oldItem.Cost * float64(oldQty),
			TotalNewCost:   newPS.Cost * float64(newQty),
			Difference:     difference,
			ProfitImpact:   (newValue - newPS.Cost*float64(newQty)) - (oldValue - oldItem.Cost*float64(oldQty)),
			PaymentMethod:  "cash",
			CreatedAt:      time.Now().AddDate(0, 0, -rand.Intn(45)),
		}

		result := database.DB.Create(&exchange)
		if result.Error != nil {
			log.Printf("Error creating exchange: %v", result.Error)
		} else {
			// Create old and new items
			oldExItem := models.ExchangeOldItem{
				ExchangeID: exchange.ID,
				SaleItemID: oldItem.ID,
				ProductID:  oldItem.ProductID,
				Quantity:   oldQty,
				Price:      oldItem.Price,
				Cost:       oldItem.Cost,
				Total:      oldValue,
				TotalCost:  oldItem.Cost * float64(oldQty),
				Condition:  "good",
			}

			newExItem := models.ExchangeNewItem{
				ExchangeID: exchange.ID,
				ProductID:  newPS.ProductID,
				Quantity:   newQty,
				Price:      newPS.Price,
				Cost:       newPS.Cost,
				Total:      newValue,
				TotalCost:  newPS.Cost * float64(newQty),
			}

			database.DB.Create(&oldExItem)
			database.DB.Create(&newExItem)
			fmt.Printf("Created exchange: %s (Difference: $%.2f)\n", exchange.ExchangeNumber, exchange.Difference)
			createdCount++
		}
	}
}

func seedStockMovements() {
	var count int64
	database.DB.Model(&models.StockMovement{}).Count(&count)
	if count > 0 {
		fmt.Printf("Stock movements already exist (%d movements found). Skipping.\n", count)
		return
	}

	// Get products
	var products []models.Product
	database.DB.Limit(50).Find(&products)

	if len(products) == 0 {
		log.Println("No products found for creating stock movements")
		return
	}

	rand.Seed(time.Now().UnixNano())
	types := []string{"in", "out", "adjustment"}
	references := []string{"PO-2025-001", "SALE-2025-001", "ADJ-001", "RET-2025-001", "EXC-2025-001"}
	notes := []string{
		"Initial stock", "Sale transaction", "Stock adjustment", "Return processing",
		"Exchange processing", "Inventory count correction", "Damage write-off",
		"Supplier delivery", "Customer purchase", "Internal transfer",
	}

	// Get existing users
	var users []models.User
	database.DB.Find(&users)
	if len(users) == 0 {
		fmt.Println("No users found. Cannot create stock movements.")
		return
	}

	for i := 0; i < 100; i++ {
		product := products[rand.Intn(len(products))]
		user := users[rand.Intn(len(users))]
		movementType := types[rand.Intn(len(types))]
		quantity := rand.Intn(50) + 1

		movement := models.StockMovement{
			ProductID: product.ID,
			UserID:    user.ID,
			Type:      movementType,
			Quantity:  quantity,
			Reference: references[rand.Intn(len(references))],
			Notes:     notes[rand.Intn(len(notes))],
			CreatedAt: time.Now().AddDate(0, 0, -rand.Intn(90)),
		}

		result := database.DB.Create(&movement)
		if result.Error != nil {
			fmt.Printf("Error creating stock movement: %v\n", result.Error)
		} else {
			fmt.Printf("Created stock movement: Product %d %s %d units\n", product.ID, movementType, quantity)
		}
	}
}

func seedActivityLogs() {
	var count int64
	database.DB.Model(&models.ActivityLog{}).Count(&count)
	if count > 0 {
		fmt.Printf("Activity logs already exist (%d logs found). Skipping.\n", count)
		return
	}

	// Get existing users
	var users []models.User
	database.DB.Find(&users)
	if len(users) == 0 {
		fmt.Println("No users found. Cannot create activity logs.")
		return
	}

	actions := []string{
		"login", "logout", "create_product", "update_product", "delete_product",
		"create_sale", "void_sale", "create_purchase_order", "update_purchase_order",
		"create_supplier", "update_supplier", "create_user", "update_user",
		"create_return", "create_exchange", "adjust_stock", "backup_database",
	}

	resources := []string{"product", "sale", "purchase_order", "supplier", "user", "return", "exchange", "stock"}

	details := []string{
		"User logged in successfully", "Created new product", "Updated product information",
		"Processed sale transaction", "Voided sale transaction", "Created purchase order",
		"Updated supplier information", "Adjusted stock levels", "Created return request",
		"Processed exchange request", "Updated user permissions", "Performed database backup",
	}

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 100; i++ {
		// Use actual user ID from existing users
		user := users[rand.Intn(len(users))]
		
		activityLog := models.ActivityLog{
			UserID:     user.ID,
			Action:     actions[rand.Intn(len(actions))],
			Resource:   resources[rand.Intn(len(resources))],
			ResourceID: uint(rand.Intn(100) + 1),
			Details:    details[rand.Intn(len(details))],
			IPAddress:  fmt.Sprintf("192.168.1.%d", rand.Intn(254)+1),
			UserAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			CreatedAt:  time.Now().AddDate(0, 0, -rand.Intn(90)),
		}

		result := database.DB.Create(&activityLog)
		if result.Error != nil {
			fmt.Printf("Error creating activity log: %v\n", result.Error)
		} else {
			fmt.Printf("Created activity log: User %d %s %s\n", activityLog.UserID, activityLog.Action, activityLog.Resource)
		}
	}
}

func seedSalePayments() {
	var count int64
	database.DB.Model(&models.SalePayment{}).Count(&count)
	if count > 0 {
		fmt.Printf("Sale payments already exist (%d payments found). Skipping.\n", count)
		return
	}

	// Get some credit sales
	var sales []models.Sale
	database.DB.Where("payment_method = ? AND payment_status != ?", "credit", "paid").Limit(25).Find(&sales)

	if len(sales) == 0 {
		log.Println("No credit sales found for creating payments")
		return
	}

	rand.Seed(time.Now().UnixNano())
	paymentTypes := []string{"downpayment", "payment", "adjustment"}
	methods := []string{"cash", "card", "transfer"}

	for _, sale := range sales {
		// Create 1-3 payments per sale
		numPayments := rand.Intn(3) + 1
		
		for i := 0; i < numPayments; i++ {
			amount := float64(rand.Intn(int(sale.AmountDue))) + float64(rand.Intn(100))/100
			if amount > sale.AmountDue {
				amount = sale.AmountDue
			}

			payment := models.SalePayment{
				SaleID:        sale.ID,
				UserID:        uint(rand.Intn(10) + 1),
				Amount:        amount,
				PaymentMethod: methods[rand.Intn(len(methods))],
				PaymentType:   paymentTypes[rand.Intn(len(paymentTypes))],
				Notes:         "Payment processing",
				CreatedAt:     time.Now().AddDate(0, 0, -rand.Intn(30)),
			}

			result := database.DB.Create(&payment)
			if result.Error != nil {
				log.Printf("Error creating sale payment: %v", result.Error)
			} else {
				fmt.Printf("Created sale payment: Sale %d - $%.2f\n", sale.ID, amount)
			}
		}
	}
}

func seedPurchasePayments() {
	var count int64
	database.DB.Model(&models.PurchasePayment{}).Count(&count)
	if count > 0 {
		fmt.Printf("Purchase payments already exist (%d payments found). Skipping.\n", count)
		return
	}

	// Get some pending purchase orders
	var purchaseOrders []models.PurchaseOrder
	database.DB.Where("payment_status != ?", "paid").Limit(25).Find(&purchaseOrders)

	if len(purchaseOrders) == 0 {
		log.Println("No pending purchase orders found for creating payments")
		return
	}

	rand.Seed(time.Now().UnixNano())
	paymentTypes := []string{"downpayment", "payment", "adjustment"}
	methods := []string{"cash", "card", "transfer"}

	for _, po := range purchaseOrders {
		// Create 1-2 payments per PO
		numPayments := rand.Intn(2) + 1
		
		for i := 0; i < numPayments; i++ {
			amount := float64(rand.Intn(int(po.AmountDue))) + float64(rand.Intn(100))/100
			if amount > po.AmountDue {
				amount = po.AmountDue
			}

			payment := models.PurchasePayment{
				PurchaseOrderID: po.ID,
				UserID:          uint(rand.Intn(10) + 1),
				Amount:          amount,
				PaymentMethod:   methods[rand.Intn(len(methods))],
				PaymentType:     paymentTypes[rand.Intn(len(paymentTypes))],
				Notes:           "Payment processing",
				CreatedAt:       time.Now().AddDate(0, 0, -rand.Intn(30)),
			}

			result := database.DB.Create(&payment)
			if result.Error != nil {
				log.Printf("Error creating purchase payment: %v", result.Error)
			} else {
				fmt.Printf("Created purchase payment: PO %d - $%.2f\n", po.ID, amount)
			}
		}
	}
}