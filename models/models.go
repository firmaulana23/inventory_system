package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a system user
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"unique;not null"`
	Password  string         `json:"-" gorm:"not null"`
	Name      string         `json:"name" gorm:"not null"`
	Role      string         `json:"role" gorm:"default:employee"` // admin, manager, employee
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Product represents an inventory item
type Product struct {
	ID          uint              `json:"id" gorm:"primaryKey"`
	Name        string            `json:"name" gorm:"not null"`
	SKU         string            `json:"sku" gorm:"unique;not null"`
	Description string            `json:"description"`
	Category    string            `json:"category"`
	Location    string            `json:"location"`
	Suppliers   []ProductSupplier `json:"suppliers" gorm:"foreignKey:ProductID"` // Multiple suppliers relationship
	IsActive    bool              `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	DeletedAt   gorm.DeletedAt    `json:"-" gorm:"index"`
}

// GetTotalStock returns total stock across all suppliers
func (p *Product) GetTotalStock() int {
	total := 0
	for _, supplier := range p.Suppliers {
		if supplier.IsActive {
			total += supplier.Stock
		}
	}
	return total
}

// GetLowestPrice returns the lowest selling price among suppliers
func (p *Product) GetLowestPrice() float64 {
	if len(p.Suppliers) == 0 {
		return 0
	}
	
	lowest := float64(0)
	first := true
	
	for _, supplier := range p.Suppliers {
		if supplier.IsActive {
			if first || supplier.Price < lowest {
				lowest = supplier.Price
				first = false
			}
		}
	}
	
	return lowest
}

// GetLowestCost returns the lowest cost among suppliers
func (p *Product) GetLowestCost() float64 {
	if len(p.Suppliers) == 0 {
		return 0
	}
	
	lowest := float64(0)
	first := true
	
	for _, supplier := range p.Suppliers {
		if supplier.IsActive {
			if first || supplier.Cost < lowest {
				lowest = supplier.Cost
				first = false
			}
		}
	}
	
	return lowest
}

// IsLowStock checks if any supplier has low stock
func (p *Product) IsLowStock() bool {
	for _, supplier := range p.Suppliers {
		if supplier.IsActive && supplier.Stock <= supplier.MinStock {
			return true
		}
	}
	return false
}

// StockMovement represents inventory movements
type StockMovement struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ProductID uint      `json:"product_id" gorm:"not null"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	Type      string    `json:"type" gorm:"not null"` // in, out, adjustment
	Quantity  int       `json:"quantity" gorm:"not null"`
	Reference string    `json:"reference"` // PO number, sale ID, etc.
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
}

// Sale represents a POS transaction
type Sale struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	SaleNumber    string         `json:"sale_number" gorm:"unique;not null"`
	UserID        uint           `json:"user_id" gorm:"not null"`
	User          User           `json:"user" gorm:"foreignKey:UserID"`
	CustomerName  string         `json:"customer_name"`
	Subtotal      float64        `json:"subtotal" gorm:"not null"`
	Tax           float64        `json:"tax" gorm:"default:0"`
	Discount      float64        `json:"discount" gorm:"default:0"`
	Total         float64        `json:"total" gorm:"not null"`
	PaymentMethod string         `json:"payment_method" gorm:"not null"`     // cash, card, transfer, credit
	PaymentDays   int            `json:"payment_days" gorm:"default:0"`      // Number of days for payment due (0 = immediate)
	PaymentStatus string         `json:"payment_status" gorm:"default:paid"` // paid, pending, overdue
	DownPayment   float64        `json:"down_payment" gorm:"default:0"`      // downpayment amount for credit sales
	DueDate       *time.Time     `json:"due_date"`
	PaidDate      *time.Time     `json:"paid_date"`
	AmountPaid    float64        `json:"amount_paid" gorm:"default:0"`
	AmountDue     float64        `json:"amount_due" gorm:"default:0"`
	Status        string         `json:"status" gorm:"default:completed"` // pending, completed, cancelled
	Items         []SaleItem     `json:"items" gorm:"foreignKey:SaleID"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

// SaleItem represents items in a sale
type SaleItem struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	SaleID    uint    `json:"sale_id" gorm:"not null"`
	ProductID uint    `json:"product_id" gorm:"not null"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  int     `json:"quantity" gorm:"not null"`
	Price     float64 `json:"price" gorm:"not null"`
	Cost      float64 `json:"cost" gorm:"not null"`
	Total     float64 `json:"total" gorm:"not null"`
}

// Supplier represents product suppliers
type Supplier struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Name          string         `json:"name" gorm:"not null"`
	Email         string         `json:"email"`
	Phone         string         `json:"phone"`
	Address       string         `json:"address"`
	ContactPerson string         `json:"contact_person"`
	Website       string         `json:"website"`
	IsActive      bool           `json:"is_active" gorm:"default:true"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

// PurchaseOrder represents purchase orders
type PurchaseOrder struct {
	ID            uint                `json:"id" gorm:"primaryKey"`
	PONumber      string              `json:"po_number" gorm:"unique;not null"`
	SupplierID    uint                `json:"supplier_id" gorm:"not null"`
	Supplier      Supplier            `json:"supplier" gorm:"foreignKey:SupplierID"`
	UserID        uint                `json:"user_id" gorm:"not null"`
	User          User                `json:"user" gorm:"foreignKey:UserID"`
	PaymentMethod string              `json:"payment_method" gorm:"default:net30"`   // cash, net7, net15, net30, net60, net90, credit
	PaymentDays   int                 `json:"payment_days" gorm:"default:30"`       // Number of days for payment due
	PaymentStatus string              `json:"payment_status" gorm:"default:pending"` // pending, paid, overdue
	Total         float64             `json:"total" gorm:"not null"`
	DownPayment   float64             `json:"down_payment" gorm:"default:0"` // downpayment amount for credit orders
	AmountPaid    float64             `json:"amount_paid" gorm:"default:0"`
	AmountDue     float64             `json:"amount_due" gorm:"default:0"`
	DueDate       *time.Time          `json:"due_date"`
	PaidDate      *time.Time          `json:"paid_date"`
	Notes         string              `json:"notes"`
	OrderDate     time.Time           `json:"order_date"`
	ReceivedDate  *time.Time          `json:"received_date"`
	Items         []PurchaseOrderItem `json:"items" gorm:"foreignKey:PurchaseOrderID"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
	DeletedAt     gorm.DeletedAt      `json:"-" gorm:"index"`
}

// PurchaseOrderItem represents items in a purchase order
type PurchaseOrderItem struct {
	ID                 uint            `json:"id" gorm:"primaryKey"`
	PurchaseOrderID    uint            `json:"purchase_order_id" gorm:"not null"`
	ProductID          uint            `json:"product_id" gorm:"not null"`
	Product            Product         `json:"product" gorm:"foreignKey:ProductID"`
	ProductSupplierID  *uint           `json:"product_supplier_id"` // Link to specific supplier for this product
	ProductSupplier    *ProductSupplier `json:"product_supplier" gorm:"foreignKey:ProductSupplierID"`
	QuantityOrdered    int             `json:"quantity_ordered" gorm:"not null"`
	QuantityReceived   int             `json:"quantity_received" gorm:"default:0"`
	UnitCost           float64         `json:"unit_cost" gorm:"not null"`
	Total              float64         `json:"total" gorm:"not null"`
}

// PurchasePayment represents payment history for purchase orders
type PurchasePayment struct {
	ID              uint          `json:"id" gorm:"primaryKey"`
	PurchaseOrderID uint          `json:"purchase_order_id" gorm:"not null"`
	PurchaseOrder   PurchaseOrder `json:"purchase_order" gorm:"foreignKey:PurchaseOrderID"`
	UserID          uint          `json:"user_id" gorm:"not null"`
	User            User          `json:"user" gorm:"foreignKey:UserID"`
	Amount          float64       `json:"amount" gorm:"not null"`
	PaymentMethod   string        `json:"payment_method" gorm:"not null"`
	PaymentType     string        `json:"payment_type" gorm:"not null"` // downpayment, payment, adjustment
	Notes           string        `json:"notes"`
	CreatedAt       time.Time     `json:"created_at"`
}

// SalePayment represents payment history for sales
type SalePayment struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	SaleID        uint      `json:"sale_id" gorm:"not null"`
	Sale          Sale      `json:"sale" gorm:"foreignKey:SaleID"`
	UserID        uint      `json:"user_id" gorm:"not null"`
	User          User      `json:"user" gorm:"foreignKey:UserID"`
	Amount        float64   `json:"amount" gorm:"not null"`
	PaymentMethod string    `json:"payment_method" gorm:"not null"`
	PaymentType   string    `json:"payment_type" gorm:"not null"` // downpayment, payment, adjustment
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"created_at"`
}

// ProductSupplier represents the relationship between products and suppliers with pricing
type ProductSupplier struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	ProductID  uint           `json:"product_id" gorm:"not null"`
	Product    Product        `json:"product" gorm:"foreignKey:ProductID"`
	SupplierID uint           `json:"supplier_id" gorm:"not null"`
	Supplier   Supplier       `json:"supplier" gorm:"foreignKey:SupplierID"`
	Cost       float64        `json:"cost" gorm:"not null"`        // Cost from this supplier
	Price      float64        `json:"price" gorm:"not null"`       // Selling price for this supplier's stock
	Stock      int            `json:"stock" gorm:"default:0"`      // Current stock from this supplier
	MinStock   int            `json:"min_stock" gorm:"default:10"` // Minimum stock for this supplier
	IsActive   bool           `json:"is_active" gorm:"default:true"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

// ActivityLog represents system activity logs
type ActivityLog struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null"`
	User       User      `json:"user" gorm:"foreignKey:UserID"`
	Action     string    `json:"action" gorm:"not null"`
	Resource   string    `json:"resource"`
	ResourceID uint      `json:"resource_id"`
	Details    string    `json:"details"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	CreatedAt  time.Time `json:"created_at"`
}

// CompanyProfile represents company information for invoices and system branding
type CompanyProfile struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	CompanyName     string         `json:"company_name" gorm:"not null"`
	CompanyAddress  string         `json:"company_address"`
	CompanyPhone    string         `json:"company_phone"`
	CompanyEmail    string         `json:"company_email"`
	CompanyWebsite  string         `json:"company_website"`
	TaxNumber       string         `json:"tax_number"`       // NPWP or tax identification number
	BusinessLicense string         `json:"business_license"` // Business registration number
	LogoBase64      string         `json:"logo_base64" gorm:"type:text"` // Company logo in base64 format
	InvoiceFooter   string         `json:"invoice_footer"`   // Custom footer text for invoices
	BankAccount     string         `json:"bank_account"`     // Bank account information
	Currency        string         `json:"currency" gorm:"default:IDR"` // Default currency
	IsActive        bool           `json:"is_active" gorm:"default:true"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

