# 🏪 Warehouse/Inventory Management System with Admin Dashboard

A comprehensive warehouse and inventory management system built with Go, featuring a complete Point of Sale (POS) system, user authentication, and a powerful admin dashboard.

## ✨ Features

### 🔐 Authentication & Authorization
- JWT-based authentication
- Role-based access control (Admin, Manager, Employee)
- User session management
- Secure password hashing with bcrypt

### 📊 Admin Dashboard
- **Real-time Statistics**: Total users, products, sales, revenue
- **Interactive Charts**: Sales trends over time with Chart.js
- **Top Products**: Best-selling products with revenue tracking
- **Recent Activity**: Latest sales transactions
- **Low Stock Alerts**: Products below minimum threshold
- **Quick Actions**: One-click access to management functions
- **Data Export**: Download sales reports as CSV
- **System Management**: Database backup and restore

### 📦 Inventory Management
- Product catalog with SKU, pricing, and stock tracking
- Stock movements logging (in, out, adjustments)
- Category management and search functionality
- Supplier management
- Purchase order system
- Automatic low stock alerts

### 💰 Point of Sale (POS)
- **Employee Dashboard**: Streamlined POS interface
- **Real-time Inventory**: Live stock updates during sales
- **Multiple Payment Methods**: Cash, Card, Transfer
- **Customer Management**: Optional customer information
- **Receipt Generation**: Automated sale number generation
- **Transaction History**: Complete sales audit trail
- **Sale Void**: Manager/Admin can void transactions

### 📈 Reporting & Analytics
- Sales reports by date range
- Payment method analytics
- Product performance metrics
- User activity logging
- System audit trails

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- SQLite (embedded)

### Installation

1. **Clone and setup the project:**
```bash
cd /Users/fmaulana/Documents/zona/inventory_system
go mod tidy
```

2. **Configure environment:**
```bash
cp .env.example .env
# Edit .env with your preferred settings
```

3. **Add sample data:**
```bash
go run cmd/seed/main.go
```

4. **Start the server:**
```bash
go run main.go
```

5. **Access the system:**
- **Login Page**: http://localhost:8080
- **Admin Dashboard**: http://localhost:8080/admin-dashboard.html
- **Employee Dashboard**: http://localhost:8080/dashboard.html

## 👤 Default Accounts

### Admin Account
- **Email**: admin@inventory.com
- **Password**: admin123
- **Access**: Full system access, admin dashboard

### Test Employee
You can create employee accounts through the admin dashboard or register new users via the API.

## 🌐 Web Interfaces

### 🔑 Login Page (`/login.html`)
- Clean, responsive login form
- Auto-redirect based on user role
- Remember me functionality
- Error handling with user feedback

### 👑 Admin Dashboard (`/admin-dashboard.html`)
- **Overview Cards**: Key metrics at a glance
- **Sales Chart**: Visual trends with dual-axis (count & revenue)
- **Top Products**: Best performers by sales volume
- **Recent Sales**: Latest transactions with details
- **Quick Actions**: 
  - User Management
  - Product Management
  - Reports & Analytics
  - System Settings
  - Database Backup
  - Data Export

### 🛒 Employee Dashboard (`/dashboard.html`)
- **POS Interface**: Full point-of-sale system
- **Product Search**: Real-time product lookup
- **Shopping Cart**: Add/remove items with quantity controls
- **Payment Processing**: Multiple payment methods
- **Quick Stats**: Personal sales metrics
- **Low Stock Alerts**: Inventory warnings

## 🔌 API Endpoints

### Authentication
```
POST /api/v1/auth/login      # User login
POST /api/v1/auth/register   # User registration
GET  /api/v1/profile         # User profile
```

### Products
```
GET    /api/v1/products                    # List products
GET    /api/v1/products/:id               # Get product
POST   /api/v1/products                   # Create product (Manager+)
PUT    /api/v1/products/:id               # Update product (Manager+)
DELETE /api/v1/admin/products/:id         # Delete product (Admin)
POST   /api/v1/products/:id/adjust-stock  # Adjust stock (Manager+)
GET    /api/v1/products/search            # Search products
GET    /api/v1/products/low-stock         # Low stock alerts
```

### Point of Sale
```
POST /api/v1/pos/sales        # Create sale
GET  /api/v1/pos/sales        # List sales
GET  /api/v1/pos/sales/:id    # Get sale details
PUT  /api/v1/pos/sales/:id/void # Void sale (Manager+)
```

### Admin Dashboard
```
GET  /api/v1/admin/dashboard/stats        # Dashboard statistics
GET  /api/v1/admin/dashboard/sales-report # Detailed sales report
GET  /api/v1/admin/system/logs            # System activity logs
GET  /api/v1/admin/users/:id/activity     # User activity logs
PUT  /api/v1/admin/system/settings        # Update system settings
POST /api/v1/admin/system/backup          # Create database backup
```

### User Management (Admin)
```
GET    /api/v1/admin/users      # List users
GET    /api/v1/admin/users/:id  # Get user
PUT    /api/v1/admin/users/:id  # Update user
DELETE /api/v1/admin/users/:id  # Delete user
```

## 🛠️ Technical Stack

### Backend
- **Framework**: Gin (Go web framework)
- **Database**: SQLite with GORM ORM
- **Authentication**: JWT tokens
- **Security**: bcrypt password hashing
- **Configuration**: Environment variables

### Frontend
- **Styling**: Tailwind CSS
- **Icons**: Font Awesome
- **Charts**: Chart.js
- **Interactions**: Vanilla JavaScript
- **Responsive**: Mobile-first design

### Database Schema
- **Users**: Authentication and authorization
- **Products**: Inventory items with full details
- **Sales & SaleItems**: POS transaction records
- **StockMovements**: Inventory tracking
- **ActivityLogs**: System audit trail
- **Suppliers**: Vendor management
- **PurchaseOrders**: Procurement tracking

## 📊 Dashboard Features

### Real-time Updates
- Dashboard refreshes every 5 minutes
- Live inventory updates during sales
- Instant low stock notifications

### Data Visualization
- Sales trend charts with dual-axis display
- Revenue and volume tracking
- Weekly sales overview

### Export Capabilities
- CSV export for sales reports
- Date range selection
- Comprehensive transaction data

### System Management
- Database backup with timestamp
- User activity monitoring
- System settings configuration

## 🔒 Security Features

- **JWT Authentication**: Secure token-based auth
- **Role-based Access**: Admin, Manager, Employee roles
- **Password Security**: bcrypt hashing
- **Activity Logging**: Complete audit trail
- **CORS Protection**: Cross-origin request handling
- **Input Validation**: Comprehensive data validation

## 📱 Mobile Responsive

Both admin and employee dashboards are fully responsive and work seamlessly on:
- Desktop computers
- Tablets
- Mobile phones
- Various screen sizes

## 🎯 Use Cases

### For Retail Stores
- Point of sale operations
- Inventory tracking
- Sales reporting
- Staff management

### For Warehouses
- Stock level monitoring
- Product organization
- Supplier management
- Purchase order tracking

### For Small Businesses
- Complete inventory solution
- User role management
- Financial reporting
- Growth analytics

## 📈 Sample Data

The system includes sample data seeder with:
- 10 diverse products across categories
- 3 supplier records
- Complete product details (SKU, pricing, stock levels)
- Realistic inventory scenarios

## 🔮 Future Enhancements

- **Barcode Scanning**: Mobile barcode integration
- **Email Notifications**: Low stock and sales alerts
- **Advanced Reporting**: More detailed analytics
- **Multi-location**: Support for multiple warehouses
- **API Integration**: Third-party system connections
- **Mobile App**: Native mobile applications

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

**Built with ❤️ using Go and modern web technologies**
