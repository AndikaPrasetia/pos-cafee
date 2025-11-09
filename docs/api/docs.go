// Package docs provides API documentation for the POS Cafe system
package docs

/*
# POS Cafe API Documentation

## Authentication

### POST /api/auth/login
Authenticate user and return JWT token

Request:
```json
{
  "username": "string (required)",
  "password": "string (required, min 8 chars)"
}
```

Response (200 OK):
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "username": "string",
      "email": "string",
      "role": "string (cashier|manager|admin)",
      "first_name": "string",
      "last_name": "string"
    },
    "token": "string (JWT token)",
    "expires_in": "integer (seconds until expiry)"
  }
}
```

### POST /api/auth/register
Create a new user account

Request:
```json
{
  "username": "string (required, 3-50 chars, alphanumeric and underscores only)",
  "email": "string (required, valid email format)",
  "password": "string (required, min 8 chars with uppercase, lowercase, number, special char)",
  "role": "string (required, one of: cashier, manager, admin)",
  "first_name": "string (required, 1-100 chars)",
  "last_name": "string (required, 1-100 chars)"
}
```

### GET /api/auth/profile
Get current user profile information (requires authentication)

Response (200 OK):
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "username": "string",
    "email": "string",
    "role": "string (cashier|manager|admin)",
    "first_name": "string",
    "last_name": "string",
    "is_active": "boolean",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}
```

## Menu Management

### GET /api/menu/categories
List all active menu categories (requires manager role)
Query parameters: active (boolean), limit (int), offset (int)

### POST /api/menu/categories
Create a new menu category (requires manager role)

Request:
```json
{
  "name": "string (required, 1-100 chars)",
  "description": "string (optional, max 500 chars)"
}
```

### GET /api/menu/items
List all menu items (requires manager role)
Query parameters: category_id (uuid), is_available (boolean), limit (int), offset (int)

### POST /api/menu/items
Create a new menu item (requires manager role)

Request:
```json
{
  "name": "string (required, 1-255 chars)",
  "category_id": "uuid (required)",
  "description": "string (optional)",
  "price": "decimal string (required, positive)",
  "cost": "decimal string (required, positive, <= price)",
  "is_available": "boolean (default true)"
}
```

## Order Processing

### POST /api/orders
Create a new draft order (requires cashier role)

Request:
```json
{
  "items": [
    {
      "menu_item_id": "uuid (required)",
      "quantity": "integer (required, positive)"
    }
  ]
}
```

### PUT /api/orders/{id}/complete
Complete an order and process payment (requires cashier role)

Request:
```json
{
  "payment_method": "string (cash|card|qris|transfer)",
  "discount_amount": "decimal string (optional)",
  "tax_amount": "decimal string (optional)"
}
```

### PUT /api/orders/{id}/cancel
Cancel an order (requires cashier role for draft orders, manager/admin for completed orders)

Request:
```json
{
  "reason": "string (required for completed orders)"
}
```

## Inventory Management

### GET /api/inventory
List inventory levels (requires manager role)
Query parameters: low_stock_only (boolean), limit (int), offset (int)

### POST /api/inventory/adjust
Make manual stock adjustment (requires manager role)

Request:
```json
{
  "menu_item_id": "uuid (required)",
  "quantity": "integer (required, positive for addition, negative for subtraction)",
  "reason": "string (required, explanation for adjustment)"
}
```

## Reporting

### GET /api/reports/daily-sales
Get daily sales report (requires manager role)
Query parameters: date (string, required in YYYY-MM-DD format)

### GET /api/reports/financial-summary
Get financial summary report (requires manager role)
Query parameters: start_date (string, required in YYYY-MM-DD), end_date (string, required in YYYY-MM-DD)
*/
func APIDocs() {}