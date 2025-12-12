# Deployment Documentation

## Environment Setup

### Prerequisites
- Go 1.21 or higher
- PostgreSQL database (12 or higher)
- Git for version control

### Environment Variables
Create a `.env` file with the following variables:

```env
# Database Configuration
DATABASE_URL=postgresql://username:password@localhost:5432/pos_cafe?sslmode=disable
DB_HOST=localhost
DB_PORT=5432
DB_USER=username
DB_PASSWORD=password
DB_NAME=pos_cafe

# Application Configuration
APP_ENV=production
APP_PORT=8080
JWT_SECRET=a-very-long-and-secure-jwt-secret-key-here-make-it-random
JWT_EXPIRY=24h

# Logging Configuration
LOG_LEVEL=info
```

## Logging Configuration

The application provides structured logging with different formats based on the environment:

- **Development (`APP_ENV=development`)**: Uses text formatting with colors for improved readability
- **Production (`APP_ENV=production`)**: Uses JSON formatting for better integration with log aggregation systems

Supported log levels (controlled by `LOG_LEVEL` environment variable):
- `debug`: Detailed information for debugging purposes
- `info`: General information about application flow
- `warn`: Warning about potential issues
- `error`: Error events that don't prevent application flow
- `fatal`: Critical errors that cause application termination
- `panic`: Errors that cause panic

The system also includes audit logging for sensitive operations like user logins, data modifications, and financial operations.

## Database Setup

### Running Migrations
```bash
# Install migrate CLI tool if not already installed
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run database migrations
migrate -path db/migrations -database $DATABASE_URL up
```

### Creating New Migrations
```bash
migrate create -ext sql -dir db/migrations -seq add_new_feature
```

## Building and Running

### Building the Binary
```bash
go build -o pos-server cmd/server/main.go
```

### Running the Server
```bash
# Set environment variables
export APP_ENV=production
export DATABASE_URL=postgresql://prod-credentials...
export JWT_SECRET=production-secret-key

# Run the server
./pos-server
```

## Docker Deployment (Optional)

### Creating Dockerfile
```Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/db/migrations ./db/migrations

EXPOSE 8080
CMD ["./main"]
```

### Building and Running Docker Container
```bash
# Build the image
docker build -t pos-cafe .

# Run the container
docker run -p 8080:8080 -e DATABASE_URL=... -e JWT_SECRET=... pos-cafe
```

## Health Checks

The application provides a health check endpoint at `/health` which returns:
- Status 200 OK when the application is running
- Status 500 when there are issues

## Configuration for Production

### Performance Settings
- Use connection pooling with appropriate limits
- Set appropriate timeout values for requests
- Configure proper logging levels

### Security Settings
- Use strong, unique JWT secrets
- Set secure HTTP headers
- Implement proper CORS policies for production domains
- Enable SSL/TLS in production