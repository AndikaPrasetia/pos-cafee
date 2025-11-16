# Build stage
FROM golang:1.24.3-alpine AS builder

# Set environment variables
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# Install git (needed for go mod download in some cases)
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary with optimizations
RUN go build -a -installsuffix cgo -ldflags="-w -s" -o main cmd/server/main.go

# Create non-root user in builder stage
RUN adduser -D -s /bin/sh -u 1001 appuser

# Final multi-stage build
FROM alpine:3.19

# Install only necessary runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Copy the non-root user from builder
COPY --from=builder /etc/passwd /etc/passwd

# Set working directory
WORKDIR /home/appuser

# Copy the binary from builder stage
COPY --from=builder --chown=1001:1001 /app/main .
COPY --from=builder --chown=1001:1001 /app/db/migrations ./db/migrations

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Command to run the application
ENTRYPOINT ["./main"]