# Inventory Management System

A comprehensive warehouse/inventory management system with Point of Sale (POS) functionality built with Go.

## Features

### User Authentication & Authorization
- JWT-based authentication
- Role-based access control (Admin, Manager, Employee)
- User registration and profile management

### Inventory Management
- Product catalog with SKU, pricing, and stock tracking
- Stock movements tracking (in, out, adjustments)
- Low stock alerts
- Category management
- Search functionality

### Point of Sale (POS)
- Complete sales transaction processing
- Multiple payment methods (cash, card, transfer)
- Real-time inventory updates
- Sales reporting and analytics
- Transaction void capability (Manager/Admin)

### Warehouse Management
- Stock adjustment capabilities
- Purchase order management
- Supplier management
- Movement history tracking

### Reporting & Analytics
- Sales reports by date range
- Top-selling products
- Payment method analytics
- Low stock reports

## API Endpoints

### Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration
- `GET /api/v1/profile` - Get user profile

### Products
- `GET /api/v1/products` - List products (with filters)
- `GET /api/v1/products/:id` - Get product details
- `POST /api/v1/products` - Create product (Manager+)
- `PUT /api/v1/products/:id` - Update product (Manager+)
- `DELETE /api/v1/products/:id` - Delete product (Admin)
- `POST /api/v1/products/:id/adjust-stock` - Adjust stock (Manager+)
- `GET /api/v1/products/search` - Search products
- `GET /api/v1/products/categories` - Get categories
- `GET /api/v1/products/low-stock` - Get low stock products

### Point of Sale
- `POST /api/v1/pos/sales` - Create sale
- `GET /api/v1/pos/sales` - List sales
- `GET /api/v1/pos/sales/:id` - Get sale details
- `PUT /api/v1/pos/sales/:id/void` - Void sale (Manager+)
- `GET /api/v1/pos/reports` - Sales reports

### Stock Management
- `GET /api/v1/stock-movements` - Get stock movement history

### User Management (Admin only)
- `GET /api/v1/admin/users` - List users
- `GET /api/v1/admin/users/:id` - Get user details
- `PUT /api/v1/admin/users/:id` - Update user
- `DELETE /api/v1/admin/users/:id` - Delete user

## Getting Started

### Prerequisites
- Go 1.23 or later
- PostgreSQL 16+ (or Docker for easy setup)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd inventory_system
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up the database using Docker:
```bash
docker-compose up -d postgres
```

4. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

5. Run database migrations:
```bash
go run migrations/migrate.go
```

6. Seed the database (optional):
```bash
go run cmd/seed/main.go
```

7. Run the application:
```bash
go run main.go
```

The server will start on `http://localhost:8080`

### Alternative Setup (without Docker)

If you prefer to use a local PostgreSQL installation:

1. Install PostgreSQL 16+
2. Create a database:
```sql
CREATE DATABASE inventory_system;
```
3. Update your `.env` file with the appropriate database credentials
4. Follow steps 5-7 above

### Default Admin Account
- Email: `admin@inventory.com`
- Password: `admin123`

### Docker Services

The `docker-compose.yml` includes:
- **PostgreSQL**: Database server on port 5432
- **Adminer**: Database management UI on port 8081 (http://localhost:8081)

## Configuration

Environment variables in `.env`:
- `PORT` - Server port (default: 8080)
- `DB_HOST` - PostgreSQL host (default: localhost)
- `DB_PORT` - PostgreSQL port (default: 5432)
- `DB_USER` - PostgreSQL username
- `DB_PASSWORD` - PostgreSQL password
- `DB_NAME` - PostgreSQL database name
- `JWT_SECRET` - JWT signing secret
- `ADMIN_EMAIL` - Default admin email
- `ADMIN_PASSWORD` - Default admin password

## Database Schema

The system uses PostgreSQL with the following main entities:
- Users (authentication and authorization)
- Products (inventory items)
- Sales & SaleItems (POS transactions)
- StockMovements (inventory tracking)
- Suppliers (vendor management)
- PurchaseOrders (procurement)

## Security Features

- Password hashing with bcrypt
- JWT token-based authentication
- Role-based access control
- CORS support
- Input validation and sanitization

## API Response Format

Success responses:
```json
{
  "data": {...},
  "message": "Success message"
}
```

Error responses:
```json
{
  "error": "Error message"
}
```

## Development

### Project Structure
```
inventory_system/
├── cmd/                    # Command-line tools
│   └── seed/              # Database seeding
├── database/              # Database connection
├── handlers/              # HTTP request handlers
├── middleware/            # HTTP middleware
├── migrations/            # Database migrations
├── models/                # Data models
├── routes/                # Route definitions
├── scripts/               # Setup scripts
├── templates/             # HTML templates
└── main.go               # Application entry point
```

### Tech Stack
- **Backend**: Go (Gin framework)
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT
- **Frontend**: HTML templates with vanilla JavaScript
- **Containerization**: Docker & Docker Compose

## Testing

Use tools like Postman or curl to test the API endpoints. Start with authentication:

```bash
# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@inventory.com", "password": "admin123"}'

# Use the returned token for authenticated requests
curl -X GET http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License.

## API Documentation

### Authentication

#### Login
- **POST** `/api/v1/auth/login`
- **Request Body:**
```json
{
  "email": "user@example.com",
  "password": "yourpassword"
}
```
- **Response:**
```json
{
  "token": "<JWT_TOKEN>",
  "user": { "id": 1, "email": "user@example.com", ... }
}
```
- **Example curl:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@inventory.com", "password": "admin123"}'
```

#### Register
- **POST** `/api/v1/auth/register`
- **Request Body:**
```json
{
  "email": "user@example.com",
  "password": "yourpassword",
  "name": "User Name",
  "role": "employee"
}
```
- **Response:**
```json
{
  "user": { "id": 2, "email": "user@example.com", ... }
}
```
- **Example curl:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "user2@example.com", "password": "pass123", "name": "User 2"}'
```

### Products

#### List Products
- **GET** `/api/v1/products`
- **Response:**
```json
[
  { "id": 1, "name": "Product A", "sku": "SKU001", ... },
  { "id": 2, "name": "Product B", ... }
]
```
- **Example curl:**
```bash
curl -X GET http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

#### Create Product
- **POST** `/api/v1/products`
- **Request Body:**
```json
{
  "name": "Product A",
  "sku": "SKU001",
  "category": "Category 1",
  "price": 10000,
  "quantity": 10
}
```
- **Response:**
```json
{
  "id": 1,
  "name": "Product A",
  "sku": "SKU001",
  ...
}
```
- **Example curl:**
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"name": "Product A", "sku": "SKU001", "category": "Category 1", "price": 10000, "quantity": 10}'
```

#### Update Product
- **PUT** `/api/v1/products/:id`
- **Request Body:**
```json
{
  "name": "Product A Updated",
  "category": "Category 2",
  "price": 12000,
  "quantity": 15
}
```
- **Response:**
```json
{
  "id": 1,
  "name": "Product A Updated",
  ...
}
```
- **Example curl:**
```bash
curl -X PUT http://localhost:8080/api/v1/products/1 \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"name": "Product A Updated", "category": "Category 2", "price": 12000, "quantity": 15}'
```

#### Delete Product
- **DELETE** `/api/v1/products/:id`
- **Response:**
```json
{
  "message": "Product deleted"
}
```
- **Example curl:**
```bash
curl -X DELETE http://localhost:8080/api/v1/products/1 \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

### Sales (POS)

#### Create Sale
- **POST** `/api/v1/pos/sales`
- **Request Body:**
```json
{
  "customer_name": "John Doe",
  "payment_method": "cash",
  "items": [
    { "product_id": 1, "quantity": 2 },
    { "product_id": 2, "quantity": 1 }
  ],
  "discount": 0,
  "tax": 0
}
```
- **Response:**
```json
{
  "id": 1,
  "sale_number": "S0001",
  ...
}
```
- **Example curl:**
```bash
curl -X POST http://localhost:8080/api/v1/pos/sales \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"customer_name": "John Doe", "payment_method": "cash", "items": [{"product_id": 1, "quantity": 2}]}'
```

### Users (Admin)

#### List Users
- **GET** `/api/v1/admin/users`
- **Response:**
```json
[
  { "id": 1, "email": "admin@inventory.com", "role": "admin", ... },
  { "id": 2, "email": "user2@example.com", ... }
]
```
- **Example curl:**
```bash
curl -X GET http://localhost:8080/api/v1/admin/users \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

#### Update User
- **PUT** `/api/v1/admin/users/:id`
- **Request Body:**
```json
{
  "name": "User Updated",
  "role": "manager"
}
```
- **Response:**
```json
{
  "id": 2,
  "name": "User Updated",
  ...
}
```
- **Example curl:**
```bash
curl -X PUT http://localhost:8080/api/v1/admin/users/2 \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"name": "User Updated", "role": "manager"}'
```

#### Delete User
- **DELETE** `/api/v1/admin/users/:id`
- **Response:**
```json
{
  "message": "User deleted"
}
```
- **Example curl:**
```bash
curl -X DELETE http://localhost:8080/api/v1/admin/users/2 \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

### Stock Movements

#### Get Stock Movement History
- **GET** `/api/v1/stock-movements`
- **Response:**
```json
[
  { "id": 1, "product_id": 1, "quantity": 10, "type": "in", ... },
  { "id": 2, "product_id": 2, "quantity": 5, "type": "out", ... }
]
```
- **Example curl:**
```bash
curl -X GET http://localhost:8080/api/v1/stock-movements \
  -H "Authorization: Bearer <JWT_TOKEN>"
```

### Reports

#### Sales Reports
- **GET** `/api/v1/pos/reports`
- **Response:**
```json
{
  "total_sales": 100000,
  "total_orders": 50,
  "sales_by_product": [
    { "product_id": 1, "quantity_sold": 20, "total_earned": 200000 },
    { "product_id": 2, "quantity_sold": 10, "total_earned": 100000 }
  ]
}
```
- **Example curl:**
```bash
curl -X GET http://localhost:8080/api/v1/pos/reports \
  -H "Authorization: Bearer <JWT_TOKEN>"
```
