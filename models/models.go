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
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	SKU         string         `json:"sku" gorm:"unique;not null"`
	Description string         `json:"description"`
	Category    string         `json:"category"`
	Price       float64        `json:"price" gorm:"not null"`
	Cost        float64        `json:"cost" gorm:"not null"`
	HPP         float64        `json:"hpp" gorm:"default:0;comment:Harga Pokok Penjualan"`
	Quantity    int            `json:"quantity" gorm:"default:0"`
	MinStock    int            `json:"min_stock" gorm:"default:10"`
	MaxStock    int            `json:"max_stock" gorm:"default:1000"`
	Location    string         `json:"location"`
	Supplier    string         `json:"supplier"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
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
	ID              uint           `json:"id" gorm:"primaryKey"`
	SaleNumber      string         `json:"sale_number" gorm:"unique;not null"`
	UserID          uint           `json:"user_id" gorm:"not null"`
	User            User           `json:"user" gorm:"foreignKey:UserID"`
	CustomerName    string         `json:"customer_name"`
	Subtotal        float64        `json:"subtotal" gorm:"not null"`
	Tax             float64        `json:"tax" gorm:"default:0"`
	Discount        float64        `json:"discount" gorm:"default:0"`
	Total           float64        `json:"total" gorm:"not null"`
	PaymentMethod   string         `json:"payment_method" gorm:"not null"`     // cash, card, transfer, credit
	PaymentTerm     string         `json:"payment_term" gorm:"default:cash"`   // cash, net7, net15, net30, net60, net90
	PaymentStatus   string         `json:"payment_status" gorm:"default:paid"` // paid, pending, overdue
	DueDate         *time.Time     `json:"due_date"`
	PaidDate        *time.Time     `json:"paid_date"`
	AmountPaid      float64        `json:"amount_paid" gorm:"default:0"`
	AmountDue       float64        `json:"amount_due" gorm:"default:0"`
	Status          string         `json:"status" gorm:"default:completed"` // pending, completed, cancelled
	Items           []SaleItem     `json:"items" gorm:"foreignKey:SaleID"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// SaleItem represents items in a sale
type SaleItem struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	SaleID    uint    `json:"sale_id" gorm:"not null"`
	ProductID uint    `json:"product_id" gorm:"not null"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  int     `json:"quantity" gorm:"not null"`
	Price     float64 `json:"price" gorm:"not null"`
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
	IsActive      bool           `json:"is_active" gorm:"default:true"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

// PurchaseOrder represents purchase orders
type PurchaseOrder struct {
	ID            uint                `json:"id" gorm:"primaryKey"`
	PONumber      string              `json:"po_number" gorm:"unique;not null"`
	Supplier      string              `json:"supplier" gorm:"not null"`
	UserID        uint                `json:"user_id" gorm:"not null"`
	User          User                `json:"user" gorm:"foreignKey:UserID"`
	Status        string              `json:"status" gorm:"default:pending"`     // pending, approved, received, cancelled
	PaymentMethod string              `json:"payment_method" gorm:"default:net30"` // cash, net7, net15, net30, net60, net90, credit
	PaymentTerm   string              `json:"payment_term" gorm:"default:net30"`   // cash, net7, net15, net30, net60, net90
	PaymentStatus string              `json:"payment_status" gorm:"default:pending"` // pending, paid, overdue
	Total         float64             `json:"total" gorm:"not null"`
	DownPayment   float64             `json:"down_payment" gorm:"default:0"`         // downpayment amount for credit orders
	AmountPaid    float64             `json:"amount_paid" gorm:"default:0"`
	AmountDue     float64             `json:"amount_due" gorm:"default:0"`
	DueDate       *time.Time          `json:"due_date"`
	PaidDate      *time.Time          `json:"paid_date"`
	Notes         string              `json:"notes"`
	OrderDate     time.Time           `json:"order_date"`
	ExpectedDate  time.Time           `json:"expected_date"`
	ReceivedDate  *time.Time          `json:"received_date"`
	Items         []PurchaseOrderItem `json:"items" gorm:"foreignKey:PurchaseOrderID"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
	DeletedAt     gorm.DeletedAt      `json:"-" gorm:"index"`
}

// PurchaseOrderItem represents items in a purchase order
type PurchaseOrderItem struct {
	ID               uint    `json:"id" gorm:"primaryKey"`
	PurchaseOrderID  uint    `json:"purchase_order_id" gorm:"not null"`
	ProductID        uint    `json:"product_id" gorm:"not null"`
	Product          Product `json:"product" gorm:"foreignKey:ProductID"`
	QuantityOrdered  int     `json:"quantity_ordered" gorm:"not null"`
	QuantityReceived int     `json:"quantity_received" gorm:"default:0"`
	UnitCost         float64 `json:"unit_cost" gorm:"not null"`
	Total            float64 `json:"total" gorm:"not null"`
}

// PurchasePayment represents payment history for purchase orders
type PurchasePayment struct {
	ID                uint          `json:"id" gorm:"primaryKey"`
	PurchaseOrderID   uint          `json:"purchase_order_id" gorm:"not null"`
	PurchaseOrder     PurchaseOrder `json:"purchase_order" gorm:"foreignKey:PurchaseOrderID"`
	UserID            uint          `json:"user_id" gorm:"not null"`
	User              User          `json:"user" gorm:"foreignKey:UserID"`
	Amount            float64       `json:"amount" gorm:"not null"`
	PaymentMethod     string        `json:"payment_method" gorm:"not null"`
	PaymentType       string        `json:"payment_type" gorm:"not null"` // downpayment, payment, adjustment
	Notes             string        `json:"notes"`
	CreatedAt         time.Time     `json:"created_at"`
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
