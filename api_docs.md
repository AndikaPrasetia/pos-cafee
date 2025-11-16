# POS Cafe API Documentation

## Table of Contents
1. [Authentication Endpoints](#authentication-endpoints)
2. [Menu Management Endpoints](#menu-management-endpoints)
3. [Order Processing Endpoints](#order-processing-endpoints)
4. [Inventory Management Endpoints](#inventory-management-endpoints)
5. [Expense Management Endpoints](#expense-management-endpoints)
6. [Reporting Endpoints](#reporting-endpoints)
7. [Maintenance Endpoints](#maintenance-endpoints)

---

## Authentication Endpoints

### POST /api/auth/register
Create a new user account

**Request:**
```json
{
  "username": "kasir",
  "email": "kasir@test.com",
  "password": "@Dm!n123",
  "role": "cashier",
  "first_name": "pak",
  "last_name": "kasir"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": "uuid",
    "username": "string",
    "email": "string",
    "role": "string (cashier|manager|admin)",
    "first_name": "string",
    "last_name": "string",
    "is_active": true,
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}
```

**Response (400 Bad Request):**
```json
{
  "success": false,
  "message": "Validation error message",
  "errors": {
    "field": "error message"
  }
}
```

### POST /api/auth/login
Authenticate user and return JWT token

**Request:**
```json
{
  "username": "string (required)",
  "password": "string (required, min 8 chars)"
}
```

**Response (200 OK):**
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

### GET /api/auth/profile
Get current user profile information (requires authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Response (200 OK):**
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

### PUT /api/auth/change-password
Change user password (requires authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Request:**
```json
{
  "current_password": "string (required)",
  "new_password": "string (required, min 8 chars with uppercase, lowercase, number, special char)"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Password changed successfully"
}
```

### POST /api/auth/logout
Logout user and invalidate session

**Headers:**
```
Authorization: Bearer {token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

---

## Menu Management Endpoints

### GET /api/menu/categories
List all active menu categories (requires authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Query Parameters:**
- active: boolean (optional, default true)
- limit: integer (optional, default 50)
- offset: integer (optional, default 0)

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "categories": [
      {
        "id": "uuid",
        "name": "string",
        "description": "string",
        "is_active": "boolean",
        "created_at": "timestamp",
        "updated_at": "timestamp"
      }
    ],
    "pagination": {
      "limit": "integer",
      "offset": "integer",
      "total": "integer"
    }
  }
}
```

### POST /api/menu/categories
Create a new menu category (requires manager role)

**Headers:**
```
Authorization: Bearer {token}
```

**Request:**
```json
{
  "name": "string (required, 1-100 chars)",
  "description": "string (optional, max 500 chars)"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "string",
    "description": "string",
    "is_active": true,
    "created_at": "timestamp",
    "updated_at": "timestamp"
  },
  "message": "Category created successfully"
}
```

### GET /api/menu/categories/{id}
Get a specific menu category by ID

**Headers:**
```
Authorization: Bearer {token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "string",
    "description": "string",
    "is_active": "boolean",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}
```

### PUT /api/menu/categories/{id}
Update a menu category (requires manager role)

**Headers:**
```
Authorization: Bearer {token}
```

**Request:**
```json
{
  "name": "string (1-100 chars)",
  "description": "string (max 500 chars)",
  "is_active": "boolean"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "string",
    "description": "string",
    "is_active": "boolean",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  },
  "message": "Category updated successfully"
}
```

### DELETE /api/menu/categories/{id}
Delete a menu category (requires manager role)

**Headers:**
```
Authorization: Bearer {token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Category deleted successfully"
}
```

### GET /api/menu/items
List all menu items (requires authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Query Parameters:**
- category_id: uuid (optional)
- is_available: boolean (optional)
- limit: integer (optional, default 50)
- offset: integer (optional, default 0)

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "menu_items": [
      {
        "id": "uuid",
        "name": "string",
        "category_id": "uuid",
        "category_name": "string",
        "description": "string",
        "price": "decimal string",
        "cost": "decimal string",
        "is_available": "boolean",
        "created_at": "timestamp",
        "updated_at": "timestamp"
      }
    ],
    "pagination": {
      "limit": "integer",
      "offset": "integer",
      "total": "integer"
    }
  }
}
```

### POST /api/menu/items
Create a new menu item (requires manager role)

**Headers:**
```
Authorization: Bearer {token}
```

**Request:**
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

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "string",
    "category_id": "uuid",
    "category_name": "string",
    "description": "string",
    "price": "decimal string",
    "cost": "decimal string",
    "is_available": "boolean",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  },
  "message": "Menu item created successfully"
}
```

### GET /api/menu/items/{id}
Get a specific menu item by ID

**Headers:**
```
Authorization: Bearer {token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "string",
    "category_id": "uuid",
    "category_name": "string",
    "description": "string",
    "price": "decimal string",
    "cost": "decimal string",
    "is_available": "boolean",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  }
}
```

### PUT /api/menu/items/{id}
Update a menu item (requires manager role)

**Headers:**
```
Authorization: Bearer {token}
```

**Request:**
```json
{
  "name": "string (1-255 chars)",
  "category_id": "uuid",
  "description": "string",
  "price": "decimal string (positive)",
  "cost": "decimal string (positive, <= price)",
  "is_available": "boolean"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "name": "string",
    "category_id": "uuid",
    "category_name": "string",
    "description": "string",
    "price": "decimal string",
    "cost": "decimal string",
    "is_available": "boolean",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  },
  "message": "Menu item updated successfully"
}
```

### DELETE /api/menu/items/{id}
Delete a menu item (requires manager role)

**Headers:**
```
Authorization: Bearer {token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Menu item deleted successfully"
}
```

---

## Order Processing Endpoints

### GET /api/orders
List orders with optional filtering and pagination

**Headers:**
```
Authorization: Bearer {token}
```

**Query Parameters:**
- status: string (draft|pending|completed|cancelled)
- user_id: uuid
- start_date: string (YYYY-MM-DD)
- end_date: string (YYYY-MM-DD)
- limit: integer (default 50)
- offset: integer (default 0)

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "orders": [
      {
        "id": "uuid",
        "order_number": "string",
        "user_id": "uuid",
        "username": "string",
        "first_name": "string",
        "last_name": "string",
        "status": "string (draft|pending|completed|cancelled)",
        "total_amount": "decimal string",
        "discount_amount": "decimal string",
        "tax_amount": "decimal string",
        "payment_method": "string (cash|card|qris|transfer)",
        "payment_status": "string (pending|paid|failed)",
        "completed_at": "timestamp or null",
        "created_at": "timestamp",
        "updated_at": "timestamp"
      }
    ],
    "pagination": {
      "limit": "integer",
      "offset": "integer",
      "total": "integer"
    }
  }
}
```

### POST /api/orders
Create a new draft order (requires cashier role)

**Headers:**
```
Authorization: Bearer {token}
```

**Request:**
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

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "order_number": "string",
    "user_id": "uuid",
    "status": "draft",
    "total_amount": "decimal string",
    "items": [
      {
        "id": "uuid",
        "order_id": "uuid",
        "menu_item_id": "uuid",
        "menu_item_name": "string",
        "quantity": "integer",
        "unit_price": "decimal string",
        "total_price": "decimal string"
      }
    ]
  },
  "message": "Order created successfully"
}
```

### GET /api/orders/{id}
Get a specific order by ID with its items

**Headers:**
```
Authorization: Bearer {token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "order_number": "string",
    "user_id": "uuid",
    "username": "string",
    "first_name": "string",
    "last_name": "string",
    "status": "string (draft|pending|completed|cancelled)",
    "total_amount": "decimal string",
    "discount_amount": "decimal string",
    "tax_amount": "decimal string",
    "payment_method": "string (cash|card|qris|transfer)",
    "payment_status": "string (pending|paid|failed)",
    "completed_at": "timestamp or null",
    "created_at": "timestamp",
    "updated_at": "timestamp",
    "items": [
      {
        "id": "uuid",
        "order_id": "uuid",
        "menu_item_id": "uuid",
        "menu_item_name": "string",
        "quantity": "integer",
        "unit_price": "decimal string",
        "total_price": "decimal string"
      }
    ]
  }
}
```

### POST /api/orders/{id}/items
Add an item to an existing order (requires cashier role)

**Headers:**
```
Authorization: Bearer {token}
```

**Request:**
```json
{
  "menu_item_id": "uuid (required)",
  "quantity": "integer (required, positive)"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "order_id": "uuid",
    "menu_item_id": "uuid",
    "menu_item_name": "string",
    "quantity": "integer",
    "unit_price": "decimal string",
    "total_price": "decimal string"
  },
  "message": "Item added to order successfully"
}
```

### PUT /api/orders/{id}/complete
Complete an order and process payment (requires cashier role)

**Headers:**
```
Authorization: Bearer {token}
```

**Request:**
```json
{
  "payment_method": "string (cash|card|qris|transfer)",
  "payment_status": "string (pending|paid|failed)",
  "discount_amount": "decimal string (optional)",
  "tax_amount": "decimal string (optional)"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "order_number": "string",
    "user_id": "uuid",
    "status": "completed",
    "total_amount": "decimal string",
    "discount_amount": "decimal string",
    "tax_amount": "decimal string",
    "payment_method": "string",
    "payment_status": "string",
    "completed_at": "timestamp",
    "created_at": "timestamp",
    "updated_at": "timestamp"
  },
  "message": "Order completed successfully"
}
```

### PUT /api/orders/{id}/cancel
Cancel an order (requires cashier role for draft orders, manager/admin for completed orders)

**Headers:**
```
Authorization: Bearer {token}
```

**Request:**
```json
{
  "reason": "string (required for completed orders)"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "status": "cancelled"
  },
  "message": "Order cancelled successfully"
}
```

---

## Inventory Management Endpoints

### GET /api/inventory
List inventory levels (requires authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Query Parameters:**
- low_stock_only: boolean (optional)
- limit: integer (optional, default 50)
- offset: integer (optional, default 0)

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "inventory": [
      {
        "id": "uuid",
        "menu_item_id": "uuid",
        "menu_item_name": "string",
        "current_stock": "integer",
        "minimum_stock": "integer",
        "unit": "string",
        "last_updated_at": "timestamp",
        "last_updated_by_username": "string",
        "stock_status": "string (OK|LOW)"
      }
    ],
    "pagination": {
      "limit": "integer",
      "offset": "integer",
      "total": "integer"
    }
  }
}
```

### GET /api/inventory/low-stock
Get items with stock below minimum threshold (requires authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Query Parameters:**
- limit: integer (optional, default 50)
- offset: integer (optional, default 0)

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "inventory": [
      {
        "id": "uuid",
        "menu_item_id": "uuid",
        "menu_item_name": "string",
        "current_stock": "integer",
        "minimum_stock": "integer",
        "unit": "string",
        "last_updated_at": "timestamp",
        "last_updated_by_username": "string",
        "stock_status": "string (LOW)"
      }
    ],
    "pagination": {
      "limit": "integer",
      "offset": "integer",
      "total": "integer"
    }
  }
}
```

### POST /api/inventory/adjust
Make manual stock adjustment (requires manager role)

**Headers:**
```
Authorization: Bearer {token}
```

**Request:**
```json
{
  "menu_item_id": "uuid (required)",
  "quantity": "integer (required, positive for addition, negative for subtraction)",
  "reason": "string (required, explanation for adjustment)"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "menu_item_id": "uuid",
    "transaction_type": "adjustment",
    "quantity": "integer",
    "previous_stock": "integer",
    "current_stock": "integer",
    "reason": "string",
    "user_id": "uuid",
    "created_at": "timestamp"
  },
  "message": "Inventory adjusted successfully"
}
```

### GET /api/inventory/transactions
List stock transactions with optional filtering (requires authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Query Parameters:**
- menu_item_id: uuid (optional)
- transaction_type: string (in|out|adjustment) (optional)
- start_date: string (YYYY-MM-DD) (optional)
- end_date: string (YYYY-MM-DD) (optional)
- limit: integer (default 50)
- offset: integer (default 0)

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "transactions": [
      {
        "id": "uuid",
        "menu_item_id": "uuid",
        "menu_item_name": "string",
        "transaction_type": "string (in|out|adjustment)",
        "quantity": "integer",
        "previous_stock": "integer",
        "current_stock": "integer",
        "reason": "string",
        "reference_type": "string or null",
        "reference_id": "uuid or null",
        "user_id": "uuid or null",
        "username": "string or null",
        "created_at": "timestamp"
      }
    ],
    "pagination": {
      "limit": "integer",
      "offset": "integer",
      "total": "integer"
    }
  }
}
```

---

## Expense Management Endpoints

### GET /api/expenses
List all expenses with optional filtering and pagination (requires authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Query Parameters:**
- start_date: string (YYYY-MM-DD) (optional)
- end_date: string (YYYY-MM-DD) (optional)
- category: string (optional)
- limit: integer (default 50)
- offset: integer (default 0)

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "expenses": [
      {
        "id": "uuid",
        "category": "string",
        "description": "string",
        "amount": "decimal string",
        "date": "date string (YYYY-MM-DD)",
        "user_id": "uuid",
        "username": "string",
        "created_at": "timestamp"
      }
    ],
    "pagination": {
      "limit": "integer",
      "offset": "integer",
      "total": "integer"
    }
  }
}
```

### POST /api/expenses
Create a new expense (requires manager role)

**Headers:**
```
Authorization: Bearer {token}
```

**Request:**
```json
{
  "category": "string (required, 1-100 chars)",
  "description": "string (optional, max 500 chars)",
  "amount": "decimal string (required, positive)",
  "date": "string (required, YYYY-MM-DD)"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "category": "string",
    "description": "string",
    "amount": "decimal string",
    "date": "date string (YYYY-MM-DD)",
    "user_id": "uuid",
    "username": "string",
    "created_at": "timestamp"
  },
  "message": "Expense created successfully"
}
```

### GET /api/expenses/{id}
Get a specific expense by ID (requires authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "category": "string",
    "description": "string",
    "amount": "decimal string",
    "date": "date string (YYYY-MM-DD)",
    "user_id": "uuid",
    "username": "string",
    "created_at": "timestamp"
  }
}
```

### PUT /api/expenses/{id}
Update an existing expense (requires manager role)

**Headers:**
```
Authorization: Bearer {token}
```

**Request:**
```json
{
  "category": "string (1-100 chars)",
  "description": "string (max 500 chars)",
  "amount": "decimal string (positive)",
  "date": "string (YYYY-MM-DD)"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "category": "string",
    "description": "string",
    "amount": "decimal string",
    "date": "date string (YYYY-MM-DD)",
    "user_id": "uuid",
    "username": "string",
    "created_at": "timestamp"
  },
  "message": "Expense updated successfully"
}
```

### DELETE /api/expenses/{id}
Delete an expense (requires manager role)

**Headers:**
```
Authorization: Bearer {token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Expense deleted successfully"
}
```

### GET /api/expenses/summary
Get expense summary for a date range (requires authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Query Parameters:**
- start_date: string (YYYY-MM-DD) (required)
- end_date: string (YYYY-MM-DD) (required)

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "summary": {
      "total_expenses": "decimal string",
      "categories": [
        {
          "category": "string",
          "total_amount": "decimal string"
        }
      ],
      "date_range": {
        "start_date": "string (YYYY-MM-DD)",
        "end_date": "string (YYYY-MM-DD)"
      }
    }
  }
}
```

---

## Reporting Endpoints

### GET /api/reports/daily-sales
Get daily sales report (requires authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Query Parameters:**
- date: string (YYYY-MM-DD) (default today)

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "report": {
      "sale_date": "string (YYYY-MM-DD)",
      "total_orders": "integer",
      "total_sales": "decimal string",
      "total_discount": "decimal string",
      "total_tax": "decimal string"
    }
  }
}
```

### GET /api/reports/financial-summary
Get financial summary report (requires authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Query Parameters:**
- start_date: string (YYYY-MM-DD) (required)
- end_date: string (YYYY-MM-DD) (required)

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "report": {
      "date_range": {
        "start_date": "string (YYYY-MM-DD)",
        "end_date": "string (YYYY-MM-DD)"
      },
      "sales": {
        "total_orders": "integer",
        "total_sales": "decimal string",
        "total_discount": "decimal string",
        "total_tax": "decimal string"
      },
      "expenses": {
        "total_expenses": "decimal string",
        "by_category": [
          {
            "category": "string",
            "amount": "decimal string"
          }
        ]
      },
      "net_profit": "decimal string"
    }
  }
}
```

### GET /api/reports/sales-by-category
Get sales report by category (requires authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Query Parameters:**
- start_date: string (YYYY-MM-DD) (required)
- end_date: string (YYYY-MM-DD) (required)

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "report": [
      {
        "category_name": "string",
        "total_quantity_sold": "integer",
        "total_revenue": "decimal string"
      }
    ]
  }
}
```

### GET /api/reports/top-selling-items
Get top selling items report (requires authentication)

**Headers:**
```
Authorization: Bearer {token}
```

**Query Parameters:**
- start_date: string (YYYY-MM-DD) (required)
- end_date: string (YYYY-MM-DD) (required)
- limit: integer (default 10)

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "report": [
      {
        "menu_item_id": "uuid",
        "menu_item_name": "string",
        "description": "string",
        "category_name": "string",
        "total_quantity_sold": "integer",
        "total_revenue": "decimal string",
        "times_ordered": "integer"
      }
    ]
  }
}
```

---

## Maintenance Endpoints

### GET /api/health
Check system health

**Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "status": "healthy",
    "timestamp": "timestamp"
  }
}
```

### POST /api/backup
Create database backup (requires admin role)

**Headers:**
```
Authorization: Bearer {token}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Database backup created successfully"
}
```

---

## Error Response Format

All error responses follow this format:

```json
{
  "success": false,
  "message": "Error message",
  "errors": {
    "field": "field specific error message"
  }
}
```

## Authentication

Most endpoints require authentication with a Bearer token in the Authorization header. Use the login endpoint to obtain a token.

## Data Types

- **uuid**: A standard UUID string in the format `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`
- **timestamp**: ISO 8601 format string (e.g., "2025-10-24T14:30:00Z")
- **date string**: YYYY-MM-DD format (e.g., "2025-10-24")
- **decimal string**: String representation of decimal value (e.g., "123.45")
- **boolean**: JSON boolean (true or false)
- **integer**: JSON integer