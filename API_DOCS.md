# Inventory System API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication

All endpoints except `/auth/login` and `/auth/register` require authentication via JWT token in the Authorization header:

```bash
Authorization: Bearer <your_jwt_token>
```

---

## Authentication Endpoints

### 1. Login
**Endpoint:** `POST /auth/login`

**Description:** Authenticate user and get JWT token

**Request Body:**
```json
{
    "email": "admin@example.com",
    "password": "password123"
}
```

**Curl Example:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "password123"
  }'
```

**Response:**
```json
{
    "success": true,
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
        "id": 1,
        "email": "admin@example.com",
        "name": "Administrator",
        "role": "admin",
        "is_active": true
    }
}
```

### 2. Register
**Endpoint:** `POST /auth/register`

**Description:** Register new user

**Request Body:**
```json
{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123",
    "role": "employee"
}
```

**Curl Example:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123",
    "role": "employee"
  }'
```

---

## Product Management

### 1. Get All Products
**Endpoint:** `GET /products`

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 50)
- `category` (optional): Filter by category
- `search` (optional): Search in name, SKU, description
- `active` (optional): Filter by active status (true/false)

**Curl Example:**
```bash
curl -X GET "http://localhost:8080/api/v1/products?page=1&limit=10&search=laptop" \
  -H "Authorization: Bearer <your_token>"
```

**Response:**
```json
{
    "products": [
        {
            "id": 1,
            "name": "Laptop Dell",
            "sku": "DELL-001",
            "description": "Dell Laptop 14 inch",
            "category": "Electronics",
            "location": "A1-001",
            "is_active": true,
            "suppliers": [
                {
                    "id": 1,
                    "product_id": 1,
                    "supplier_id": 1,
                    "cost": 5000000,
                    "price": 6000000,
                    "stock": 10,
                    "min_stock": 5,
                    "is_active": true,
                    "supplier": {
                        "id": 1,
                        "name": "PT Supplier A",
                        "email": "supplier@example.com"
                    }
                }
            ]
        }
    ],
    "total": 1,
    "page": 1,
    "limit": 10
}
```

### 2. Get Single Product
**Endpoint:** `GET /products/{id}`

**Curl Example:**
```bash
curl -X GET http://localhost:8080/api/v1/products/1 \
  -H "Authorization: Bearer <your_token>"
```

### 3. Create Product
**Endpoint:** `POST /products`

**Required Role:** Manager or Admin

**Request Body:**
```json
{
    "name": "Laptop Asus",
    "sku": "ASUS-001",
    "description": "Asus Gaming Laptop",
    "category": "Electronics",
    "location": "A1-002"
}
```

**Curl Example:**
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop Asus",
    "sku": "ASUS-001",
    "description": "Asus Gaming Laptop",
    "category": "Electronics",
    "location": "A1-002"
  }'
```

### 4. Update Product
**Endpoint:** `PUT /products/{id}`

**Required Role:** Manager or Admin

**Request Body:**
```json
{
    "name": "Laptop Asus Updated",
    "sku": "ASUS-001-V2",
    "description": "Updated Asus Gaming Laptop",
    "category": "Electronics",
    "location": "A1-002",
    "is_active": true
}
```

**Curl Example:**
```bash
curl -X PUT http://localhost:8080/api/v1/products/1 \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop Asus Updated",
    "sku": "ASUS-001-V2",
    "description": "Updated Asus Gaming Laptop",
    "category": "Electronics",
    "location": "A1-002",
    "is_active": true
  }'
```

### 5. Delete Product
**Endpoint:** `DELETE /products/{id}`

**Required Role:** Admin only

**Curl Example:**
```bash
curl -X DELETE http://localhost:8080/api/v1/products/1 \
  -H "Authorization: Bearer <your_token>"
```

### 6. Add Supplier to Product
**Endpoint:** `POST /products/{id}/suppliers`

**Required Role:** Manager or Admin

**Request Body:**
```json
{
    "supplier_id": 1,
    "cost": 5000000,
    "price": 6000000,
    "stock": 20,
    "min_stock": 5
}
```

**Curl Example:**
```bash
curl -X POST http://localhost:8080/api/v1/products/1/suppliers \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "supplier_id": 1,
    "cost": 5000000,
    "price": 6000000,
    "stock": 20,
    "min_stock": 5
  }'
```

### 7. Update Product Supplier
**Endpoint:** `PUT /products/{id}/suppliers/{supplier_id}`

**Required Role:** Manager or Admin

**Request Body:**
```json
{
    "cost": 4800000,
    "price": 5800000,
    "stock": 15,
    "min_stock": 3,
    "is_active": true
}
```

**Curl Example:**
```bash
curl -X PUT http://localhost:8080/api/v1/products/1/suppliers/1 \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "cost": 4800000,
    "price": 5800000,
    "stock": 15,
    "min_stock": 3,
    "is_active": true
  }'
```

### 8. Remove Supplier from Product
**Endpoint:** `DELETE /products/{id}/suppliers/{supplier_id}`

**Required Role:** Manager or Admin

**Curl Example:**
```bash
curl -X DELETE http://localhost:8080/api/v1/products/1/suppliers/1 \
  -H "Authorization: Bearer <your_token>"
```

### 9. Adjust Stock
**Endpoint:** `POST /products/{id}/suppliers/{supplier_id}/adjust-stock`

**Required Role:** Manager or Admin

**Request Body:**
```json
{
    "quantity": 10,
    "type": "in",
    "notes": "Restocking from purchase order PO-001"
}
```

**Types:**
- `in`: Add stock
- `out`: Remove stock  
- `adjustment`: Set absolute stock value

**Curl Example:**
```bash
curl -X POST http://localhost:8080/api/v1/products/1/suppliers/1/adjust-stock \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 10,
    "type": "in",
    "notes": "Restocking from purchase order PO-001"
  }'
```

### 10. Get Product Categories
**Endpoint:** `GET /products/categories`

**Curl Example:**
```bash
curl -X GET http://localhost:8080/api/v1/products/categories \
  -H "Authorization: Bearer <your_token>"
```

**Response:**
```json
[
    "Electronics",
    "Furniture",
    "Clothing",
    "Books"
]
```

### 11. Get Low Stock Products
**Endpoint:** `GET /products/low-stock`

**Curl Example:**
```bash
curl -X GET http://localhost:8080/api/v1/products/low-stock \
  -H "Authorization: Bearer <your_token>"
```

### 12. Search Products
**Endpoint:** `GET /products/search`

**Query Parameters:**
- `q` (required): Search query

**Curl Example:**
```bash
curl -X GET "http://localhost:8080/api/v1/products/search?q=laptop" \
  -H "Authorization: Bearer <your_token>"
```

---

## Supplier Management

### 1. Get All Suppliers
**Endpoint:** `GET /suppliers`

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20)

**Curl Example:**
```bash
curl -X GET "http://localhost:8080/api/v1/suppliers?page=1&limit=10" \
  -H "Authorization: Bearer <your_token>"
```

**Response:**
```json
{
    "success": true,
    "data": [
        {
            "id": 1,
            "name": "PT Supplier A",
            "category": "Electronics",
            "email": "supplier@example.com",
            "phone": "081234567890",
            "address": "Jl. Supplier No. 1",
            "contact_person": "John Doe",
            "website": "https://supplier.com",
            "is_active": true,
            "created_at": "2024-01-01T00:00:00Z",
            "updated_at": "2024-01-01T00:00:00Z"
        }
    ],
    "total": 1,
    "page": 1,
    "limit": 10
}
```

### 2. Get Single Supplier
**Endpoint:** `GET /suppliers/{id}`

**Curl Example:**
```bash
curl -X GET http://localhost:8080/api/v1/suppliers/1 \
  -H "Authorization: Bearer <your_token>"
```

### 3. Create Supplier
**Endpoint:** `POST /suppliers`

**Required Role:** Manager or Admin

**Request Body:**
```json
{
    "name": "PT New Supplier",
    "category": "Electronics",
    "email": "newsupplier@example.com",
    "phone": "081234567891",
    "address": "Jl. New Supplier No. 1",
    "contact_person": "Jane Doe",
    "website": "https://newsupplier.com"
}
```

**Curl Example:**
```bash
curl -X POST http://localhost:8080/api/v1/suppliers \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "PT New Supplier",
    "category": "Electronics",
    "email": "newsupplier@example.com",
    "phone": "081234567891",
    "address": "Jl. New Supplier No. 1",
    "contact_person": "Jane Doe",
    "website": "https://newsupplier.com"
  }'
```

### 4. Update Supplier
**Endpoint:** `PUT /suppliers/{id}`

**Required Role:** Manager or Admin

**Request Body:**
```json
{
    "name": "PT Updated Supplier",
    "category": "Electronics",
    "email": "updated@example.com",
    "phone": "081234567892",
    "address": "Jl. Updated No. 1",
    "contact_person": "Updated Person",
    "website": "https://updated.com",
    "is_active": true
}
```

**Curl Example:**
```bash
curl -X PUT http://localhost:8080/api/v1/suppliers/1 \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "PT Updated Supplier",
    "category": "Electronics",
    "email": "updated@example.com",
    "phone": "081234567892",
    "address": "Jl. Updated No. 1",
    "contact_person": "Updated Person",
    "website": "https://updated.com",
    "is_active": true
  }'
```

### 5. Delete Supplier
**Endpoint:** `DELETE /suppliers/{id}`

**Required Role:** Manager or Admin

**Curl Example:**
```bash
curl -X DELETE http://localhost:8080/api/v1/suppliers/1 \
  -H "Authorization: Bearer <your_token>"
```

### 6. Search Suppliers
**Endpoint:** `GET /suppliers/search`

**Query Parameters:**
- `q` (required): Search query

**Curl Example:**
```bash
curl -X GET "http://localhost:8080/api/v1/suppliers/search?q=supplier" \
  -H "Authorization: Bearer <your_token>"
```

---

## POS/Sales Management

### 1. Create Sale
**Endpoint:** `POST /pos/sales`

**Request Body:**
```json
{
    "customer_name": "John Customer",
    "payment_method": "cash",
    "payment_days": 0,
    "down_payment": 0,
    "items": [
        {
            "product_id": 1,
            "quantity": 2,
            "price": 6000000
        }
    ],
    "tax": 600000,
    "discount": 100000
}
```

**Payment Methods:**
- `cash`: Immediate cash payment
- `card`: Credit/debit card
- `transfer`: Bank transfer
- `credit`: Credit sale (requires payment_days)

**Curl Example:**
```bash
curl -X POST http://localhost:8080/api/v1/pos/sales \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "customer_name": "John Customer",
    "payment_method": "cash",
    "items": [
        {
            "product_id": 1,
            "quantity": 2,
            "price": 6000000
        }
    ],
    "tax": 600000,
    "discount": 100000
  }'
```

### 2. Get All Sales
**Endpoint:** `GET /pos/sales`

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20)
- `start_date` (optional): Filter from date (YYYY-MM-DD)
- `end_date` (optional): Filter to date (YYYY-MM-DD)
- `payment_method` (optional): Filter by payment method
- `payment_status` (optional): Filter by payment status

**Curl Example:**
```bash
curl -X GET "http://localhost:8080/api/v1/pos/sales?page=1&limit=10&start_date=2024-01-01&end_date=2024-12-31" \
  -H "Authorization: Bearer <your_token>"
```

### 3. Get Single Sale
**Endpoint:** `GET /pos/sales/{id}`

**Curl Example:**
```bash
curl -X GET http://localhost:8080/api/v1/pos/sales/1 \
  -H "Authorization: Bearer <your_token>"
```

### 4. Void Sale
**Endpoint:** `PUT /pos/sales/{id}/void`

**Required Role:** Manager or Admin

**Curl Example:**
```bash
curl -X PUT http://localhost:8080/api/v1/pos/sales/1/void \
  -H "Authorization: Bearer <your_token>"
```

### 5. Delete Sale
**Endpoint:** `DELETE /pos/sales/{id}`

**Required Role:** Manager or Admin

**Curl Example:**
```bash
curl -X DELETE http://localhost:8080/api/v1/pos/sales/1 \
  -H "Authorization: Bearer <your_token>"
```

### 6. Record Sale Payment
**Endpoint:** `POST /pos/sales/{id}/payment`

**Required Role:** Manager or Admin

**Request Body:**
```json
{
    "amount": 1000000,
    "payment_method": "transfer",
    "payment_type": "payment",
    "notes": "Partial payment via bank transfer"
}
```

**Payment Types:**
- `downpayment`: Initial down payment
- `payment`: Regular payment
- `adjustment`: Payment adjustment

**Curl Example:**
```bash
curl -X POST http://localhost:8080/api/v1/pos/sales/1/payment \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 1000000,
    "payment_method": "transfer",
    "payment_type": "payment",
    "notes": "Partial payment via bank transfer"
  }'
```

---

## Returns & Exchanges

### 1. Create Return
**Endpoint:** `POST /pos/returns`

**Request Body:**
```json
{
    "sale_id": 1,
    "reason": "Defective product",
    "refund_method": "cash",
    "items": [
        {
            "sale_item_id": 1,
            "product_id": 1,
            "quantity": 1,
            "price": 6000000,
            "condition": "damaged"
        }
    ]
}
```

**Refund Methods:** `cash`, `card`, `transfer`, `store_credit`
**Item Conditions:** `good`, `damaged`, `expired`

**Curl Example:**
```bash
curl -X POST http://localhost:8080/api/v1/pos/returns \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "sale_id": 1,
    "reason": "Defective product",
    "refund_method": "cash",
    "items": [
        {
            "sale_item_id": 1,
            "product_id": 1,
            "quantity": 1,
            "price": 6000000,
            "condition": "damaged"
        }
    ]
  }'
```

### 2. Get All Returns
**Endpoint:** `GET /pos/returns`

**Query Parameters:**
- `page`, `limit`, `start_date`, `end_date`, `refund_method`

**Curl Example:**
```bash
curl -X GET "http://localhost:8080/api/v1/pos/returns?page=1&limit=10" \
  -H "Authorization: Bearer <your_token>"
```

### 3. Create Exchange
**Endpoint:** `POST /pos/exchanges`

**Request Body:**
```json
{
    "sale_id": 1,
    "reason": "Size exchange",
    "payment_method": "cash",
    "old_items": [
        {
            "sale_item_id": 1,
            "product_id": 1,
            "quantity": 1,
            "price": 6000000,
            "condition": "good"
        }
    ],
    "new_items": [
        {
            "product_id": 2,
            "quantity": 1,
            "price": 6500000
        }
    ]
}
```

**Curl Example:**
```bash
curl -X POST http://localhost:8080/api/v1/pos/exchanges \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "sale_id": 1,
    "reason": "Size exchange",
    "payment_method": "cash",
    "old_items": [
        {
            "sale_item_id": 1,
            "product_id": 1,
            "quantity": 1,
            "price": 6000000,
            "condition": "good"
        }
    ],
    "new_items": [
        {
            "product_id": 2,
            "quantity": 1,
            "price": 6500000
        }
    ]
  }'
```

---

## Admin Dashboard

### 1. Get Dashboard Stats
**Endpoint:** `GET /admin/dashboard/stats`

**Required Role:** Admin only

**Curl Example:**
```bash
curl -X GET http://localhost:8080/api/v1/admin/dashboard/stats \
  -H "Authorization: Bearer <your_token>"
```

**Response:**
```json
{
    "success": true,
    "data": {
        "total_users": 10,
        "total_products": 50,
        "total_sales": 100,
        "today_sales": 5,
        "total_revenue": 50000000,
        "today_revenue": 2500000,
        "total_profit": 10000000,
        "today_profit": 500000,
        "low_stock_products": 3,
        "recent_sales": [...],
        "top_products": [...],
        "sales_chart": [...]
    }
}
```

### 2. Get System Logs
**Endpoint:** `GET /admin/system/logs`

**Required Role:** Admin only

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 50)

**Curl Example:**
```bash
curl -X GET "http://localhost:8080/api/v1/admin/system/logs?page=1&limit=20" \
  -H "Authorization: Bearer <your_token>"
```

### 3. User Management
**Endpoint:** `GET /admin/users`

**Required Role:** Admin only

**Curl Example:**
```bash
curl -X GET http://localhost:8080/api/v1/admin/users \
  -H "Authorization: Bearer <your_token>"
```

**Create User:**
```bash
curl -X POST http://localhost:8080/api/v1/admin/users \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New Employee",
    "email": "employee@example.com",
    "password": "password123",
    "role": "employee"
  }'
```

**Update User:**
```bash
curl -X PUT http://localhost:8080/api/v1/admin/users/1 \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Employee",
    "email": "updated@example.com",
    "role": "manager",
    "is_active": true
  }'
```

**Delete User:**
```bash
curl -X DELETE http://localhost:8080/api/v1/admin/users/1 \
  -H "Authorization: Bearer <your_token>"
```

---

## Stock Movements

### 1. Get Stock Movements
**Endpoint:** `GET /stock-movements`

**Query Parameters:**
- `page`, `limit`: Pagination
- `product_id`: Filter by product
- `type`: Filter by movement type (`in`, `out`, `adjustment`)

**Curl Example:**
```bash
curl -X GET "http://localhost:8080/api/v1/stock-movements?product_id=1&type=in" \
  -H "Authorization: Bearer <your_token>"
```

---

## Purchase Orders

### 1. Create Purchase Order
**Endpoint:** `POST /purchase-orders`

**Required Role:** Manager or Admin

**Request Body:**
```json
{
    "supplier_id": 1,
    "payment_method": "net30",
    "payment_days": 30,
    "down_payment": 1000000,
    "notes": "Urgent order",
    "items": [
        {
            "product_id": 1,
            "product_supplier_id": 1,
            "quantity_ordered": 10,
            "unit_cost": 5000000
        }
    ]
}
```

**Curl Example:**
```bash
curl -X POST http://localhost:8080/api/v1/purchase-orders \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "supplier_id": 1,
    "payment_method": "net30",
    "payment_days": 30,
    "down_payment": 1000000,
    "notes": "Urgent order",
    "items": [
        {
            "product_id": 1,
            "product_supplier_id": 1,
            "quantity_ordered": 10,
            "unit_cost": 5000000
        }
    ]
  }'
```

### 2. Get Purchase Orders
**Endpoint:** `GET /purchase-orders`

**Required Role:** Manager or Admin

**Curl Example:**
```bash
curl -X GET "http://localhost:8080/api/v1/purchase-orders?page=1&limit=10" \
  -H "Authorization: Bearer <your_token>"
```

---

## Error Responses

All endpoints may return these error responses:

### 400 Bad Request
```json
{
    "error": "Invalid input data"
}
```

### 401 Unauthorized
```json
{
    "error": "Invalid or missing token"
}
```

### 403 Forbidden
```json
{
    "error": "Insufficient permissions"
}
```

### 404 Not Found
```json
{
    "error": "Resource not found"
}
```

### 500 Internal Server Error
```json
{
    "error": "Internal server error"
}
```

---

## Role Permissions

### Employee
- View products, suppliers
- Create sales, returns, exchanges
- View stock movements
- View own profile

### Manager
- All employee permissions
- Create/update/delete products and suppliers
- Manage purchase orders
- Record payments
- Void/delete sales

### Admin
- All manager permissions
- User management
- System administration
- Dashboard access
- Delete products (permanent)
- System backups

---

## Rate Limiting

- **Authentication endpoints:** 5 requests per minute per IP
- **Other endpoints:** 100 requests per minute per user

---

## Pagination

Most list endpoints support pagination:

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Items per page (default: varies by endpoint)

**Response Format:**
```json
{
    "data": [...],
    "total": 100,
    "page": 1,
    "limit": 20
}
```