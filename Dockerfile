FROM golang:1.21-alpine AS builder

# Set environment variables
ENV CGO_ENABLED=0
ENV GOOS=linux

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN go build -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN adduser -D -s /bin/sh appuser

WORKDIR /home/appuser

# Copy the binary from builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/db/migrations ./db/migrations

# Change ownership to appuser
RUN chown appuser:appuser main
RUN chown -R appuser:appuser db/migrations

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Command to run the application
CMD ["./main"]