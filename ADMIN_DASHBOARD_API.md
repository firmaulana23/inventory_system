# Inventory System API Documentation

## Admin Dashboard API Endpoints

### Authentication Required
All admin dashboard endpoints require authentication with a valid JWT token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

### Admin Role Required
All admin dashboard endpoints require the user to have `admin` role.

---

## Dashboard Statistics

### GET `/api/v1/admin/dashboard/stats`
Returns comprehensive dashboard statistics.

**Response:**
```json
{
  "success": true,
  "data": {
    "total_users": 5,
    "total_products": 10,
    "total_sales": 25,
    "today_sales": 3,
    "total_revenue": 1250.50,
    "today_revenue": 150.75,
    "low_stock_products": 2,
    "recent_sales": [
      {
        "id": 1,
        "sale_number": "SALE-20250710-1234567890",
        "customer_name": "John Doe",
        "total": 99.99,
        "created_at": "2025-07-10T10:30:00Z",
        "user": {
          "id": 1,
          "name": "Admin User"
        }
      }
    ],
    "top_products": [
      {
        "product_id": 1,
        "product_name": "Wireless Mouse",
        "total_sold": 15,
        "revenue": 449.85
      }
    ],
    "sales_chart": [
      {
        "date": "2025-07-04",
        "sales": 5,
        "revenue": 250.00
      }
    ]
  }
}
```

---

## Sales Reports

### GET `/api/v1/admin/dashboard/sales-report`
Returns detailed sales report for a specified date range.

**Query Parameters:**
- `start_date` (optional): Start date in YYYY-MM-DD format (default: 30 days ago)
- `end_date` (optional): End date in YYYY-MM-DD format (default: today)

**Response:**
```json
{
  "success": true,
  "data": {
    "report": [
      {
        "date": "2025-07-10",
        "total_sales": 3,
        "total_revenue": 150.75,
        "total_items": 8
      }
    ],
    "start_date": "2025-06-10",
    "end_date": "2025-07-10"
  }
}
```

---

## System Logs

### GET `/api/v1/admin/system/logs`
Returns system activity logs.

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 50)

**Response:**
```json
{
  "success": true,
  "data": {
    "logs": [
      {
        "id": 1,
        "action": "login",
        "resource": "/api/v1/auth/login",
        "resource_id": 0,
        "details": "POST /api/v1/auth/login",
        "ip_address": "192.168.1.100",
        "user_agent": "Mozilla/5.0...",
        "created_at": "2025-07-10T10:30:00Z",
        "user": {
          "id": 1,
          "name": "Admin User",
          "email": "admin@inventory.com"
        }
      }
    ],
    "total": 25,
    "page": 1,
    "limit": 50
  }
}
```

---

## User Activity

### GET `/api/v1/admin/users/:id/activity`
Returns activity logs for a specific user.

**Parameters:**
- `id`: User ID

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20)

**Response:**
```json
{
  "success": true,
  "data": {
    "logs": [
      {
        "id": 5,
        "action": "create_sale",
        "resource": "/api/v1/pos/sales",
        "resource_id": 10,
        "details": "Created new sale transaction",
        "ip_address": "192.168.1.101",
        "created_at": "2025-07-10T14:45:00Z"
      }
    ],
    "total": 12,
    "page": 1,
    "limit": 20
  }
}
```

---

## System Management

### PUT `/api/v1/admin/system/settings`
Updates system configuration settings.

**Request Body:**
```json
{
  "store_name": "My Inventory Store",
  "store_address": "123 Main St, City, State 12345",
  "store_phone": "+1-555-0123",
  "low_stock_threshold": 10,
  "tax_rate": 0.08
}
```

**Response:**
```json
{
  "success": true,
  "message": "Settings updated successfully",
  "data": {
    "store_name": "My Inventory Store",
    "store_address": "123 Main St, City, State 12345",
    "store_phone": "+1-555-0123",
    "low_stock_threshold": 10,
    "tax_rate": 0.08
  }
}
```

### POST `/api/v1/admin/system/backup`
Creates a database backup.

**Response:**
```json
{
  "success": true,
  "message": "Database backup created successfully",
  "data": {
    "filename": "backup_20250710_143000.sql",
    "created_at": "2025-07-10T14:30:00Z",
    "size": "2.5 MB"
  }
}
```

### POST `/api/v1/admin/system/restore`
Restores database from backup.

**Form Data:**
- `filename`: Backup filename to restore from

**Response:**
```json
{
  "success": true,
  "message": "Database restored successfully from backup_20250710_143000.sql"
}
```

---

## User Management

### GET `/api/v1/admin/users`
Returns list of all users.

**Response:**
```json
[
  {
    "id": 1,
    "email": "admin@inventory.com",
    "name": "Admin User",
    "role": "admin",
    "is_active": true,
    "created_at": "2025-07-10T08:00:00Z",
    "updated_at": "2025-07-10T08:00:00Z"
  }
]
```

### GET `/api/v1/admin/users/:id`
Returns details of a specific user.

**Parameters:**
- `id`: User ID

### PUT `/api/v1/admin/users/:id`
Updates user information.

**Parameters:**
- `id`: User ID

**Request Body:**
```json
{
  "name": "Updated Name",
  "role": "manager",
  "is_active": true
}
```

### DELETE `/api/v1/admin/users/:id`
Soft deletes a user.

**Parameters:**
- `id`: User ID

**Response:**
```json
{
  "message": "User deleted successfully"
}
```

---

## Product Management (Admin)

### DELETE `/api/v1/admin/products/:id`
Permanently deletes a product (admin only).

**Parameters:**
- `id`: Product ID

**Response:**
```json
{
  "message": "Product deleted successfully"
}
```

---

## Error Responses

All endpoints may return error responses in the following format:

```json
{
  "success": false,
  "error": "Error message describing what went wrong"
}
```

Common error status codes:
- `400 Bad Request`: Invalid input data
- `401 Unauthorized`: Missing or invalid authentication token
- `403 Forbidden`: Insufficient permissions (not admin)
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server-side error

---

## Dashboard Frontend

The admin dashboard is accessible at:
- **URL:** `/admin-dashboard.html`
- **Login:** `/login.html`

### Features:
1. **Real-time Statistics** - Overview cards showing key metrics
2. **Sales Chart** - Visual representation of sales over time
3. **Top Products** - Best-selling products list
4. **Recent Activity** - Latest sales transactions
5. **Quick Actions** - One-click access to management functions
6. **Data Export** - Download sales reports as CSV
7. **Database Backup** - Create system backups

### Dashboard Authentication Flow:
1. User logs in at `/login.html`
2. System validates credentials via `/api/v1/auth/login`
3. If user has admin role, redirect to `/admin-dashboard.html`
4. Dashboard loads data from various admin endpoints
5. Real-time updates every 5 minutes

### Sample Admin Credentials:
- **Email:** admin@inventory.com
- **Password:** admin123
