# POS Cafe System

A comprehensive Point of Sale (POS) system designed specifically for cafe businesses, built with Go and PostgreSQL. This system handles inventory management, order processing, payment flows, and financial reporting with proper transaction tracking and data consistency.

## 🚀 Features

- **User Authentication & Authorization**: Role-based access control (admin, manager, cashier)
- **Menu Management**: Create, update, and manage categories and menu items
- **Order Processing**: Complete order lifecycle from creation to completion
- **Inventory Management**: Real-time stock tracking with low-stock alerts
- **Financial Reporting**: Daily sales, financial summaries, and top-selling items reports
- **Expense Tracking**: Record and manage business expenses
- **Docker Support**: Containerized deployment with production-ready configurations
- **API Documentation**: Comprehensive RESTful API with proper error handling

## 🛠️ Tech Stack

- **Language**: Go 1.24.3
- **Web Framework**: Gin-Gonic
- **Database**: PostgreSQL
- **Database Tools**: 
  - sqlc for type-safe SQL queries
  - Golang Migrate for database migrations
- **Authentication**: JWT-based authentication
- **Logging**: Logrus
- **Configuration**: Viper
- **Validation**: Go Playground Validator
- **Monetary Calculations**: Shopspring Decimal
- **Containerization**: Docker

## 📋 Prerequisites

- Go 1.24.3 or higher
- PostgreSQL 12+ (or Neon PostgreSQL)
- Docker and Docker Compose (optional, for containerized deployment)

## 🚀 Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/AndikaPrasetia/pos-cafee.git
cd pos-cafee
```

### 2. Set up environment variables

Copy the example environment file and configure your settings:

```bash
cp .env.example .env
```

Edit the `.env` file with your database and application settings:

```env
# Database Configuration
DATABASE_URL=postgresql://username:password@localhost:5432/pos_cafe

# Server Configuration
PORT=8080
ENVIRONMENT=development

# JWT Configuration
JWT_SECRET=your_very_secure_random_string_here
JWT_EXPIRY=24h
```

### 3. Run database migrations

```bash
# If using migrate CLI tool
migrate -path db/migrations -database $DATABASE_URL up
```

### 4. Install dependencies

```bash
go mod tidy
```

### 5. Run the application

```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

## 🐳 Docker Deployment

### Using Docker Compose (Recommended for Development)

```bash
docker-compose up -d
```

### Building and Running the Docker Image

```bash
# Build the image
docker build -t pos-cafee .

# Run the container
docker run -p 8080:8080 --env-file .env pos-cafee
```

### Production Deployment

The application is optimized for deployment on platforms like Koyeb, Heroku, or AWS with proper environment configuration.

## 📚 API Documentation

The API follows RESTful principles and is secured with JWT authentication. Different endpoints require different user roles:

- **Public endpoints**: No authentication required (login, register)
- **Cashier**: Can handle orders
- **Manager**: Can manage menu, inventory, and reports
- **Admin**: Full access, including maintenance features

For detailed API documentation, refer to `api_docs.md` in the project root.

### Key Endpoints

- `POST /api/auth/login` - User authentication
- `GET /api/orders` - List orders (requires authentication)
- `POST /api/orders` - Create new order
- `PUT /api/orders/:id/complete` - Complete an order
- `GET /api/inventory` - List inventory items
- `GET /api/reports/daily-sales` - Daily sales report
- `GET /api/health` - Health check endpoint

## 🧪 Testing

To run all tests:

```bash
go test ./...
```

The project aims for 80%+ test coverage across all business logic and API handlers.

## 🔧 Configuration

The application is configured using environment variables. Key settings include:

- `DATABASE_URL`: PostgreSQL connection string
- `PORT`: Server port (default 8080)
- `ENVIRONMENT`: Environment mode ("development", "production")
- `JWT_SECRET`: Secret key for JWT token signing
- `JWT_EXPIRY`: JWT token expiry duration
- `REDIS_URL`: Redis connection URL for caching (format: redis://host:port or rediss:// for SSL)

## 🗄️ Redis Configuration

The application uses Redis for caching frequently accessed data like menu items, categories, and reports.

### Local Development
For local development, you can set up Redis one of these ways:

1. **Using a local Redis server**:
   - Install Redis locally and start the service
   - Set REDIS_URL to your local Redis instance, e.g., `redis://localhost:6379`

2. **Using Redis provided by Docker Compose** (if you add it back):
   - Run Redis in a Docker container alongside the app
   - Set REDIS_URL to point to the Redis container, e.g., `redis://redis:6379`

### Production Deployment
For production deployment on Render or similar platforms:
- Obtain your Redis URL from your Redis service provider
- Set the REDIS_URL environment variable to your Redis URL in your deployment settings

## 🔐 Security Features

- Input validation on all endpoints
- Parameterized queries to prevent SQL injection
- JWT-based authentication with role-based authorization
- Non-root containers for Docker deployments
- SSL mode for database connections (especially with Neon)
- Rate limiting and request size limits

## 📈 Performance Optimization

- Indexes on frequently queried database fields
- Optimized SQL queries for large datasets
- Efficient memory usage with Go's garbage collector
- Connection pooling for database operations
- Proper timeout settings for HTTP requests

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🐛 Issues & Support

If you encounter any issues or have questions, please file an issue in the GitHub repository.

## 🙏 Acknowledgments

- Built with the Gin-Gonic web framework
- Database powered by PostgreSQL
- Containerized with Docker
- Inspired by modern POS systems designed for cafe businesses