# POS Cafe Server

A Point of Sale (POS) system for cafes, built with Go and PostgreSQL.

## Overview

This is the server component for a POS system designed for cafes. It handles inventory management, order processing, payment flows, and financial reporting with proper transaction tracking and data consistency.

## Features

- Inventory management
- Order processing
- Payment flows
- Financial reporting
- Role-based access control
- Database transaction management
- RESTful API endpoints

## Quick Start

To run the POS Cafe server:

```bash
docker run -d \
  --name pos-cafee-server \
  -p 8080:8080 \
  -e DATABASE_URL=postgresql://user:password@host:port/database?sslmode=require \
  -e JWT_SECRET=your-super-secret-jwt-key \
  -e PORT=8080 \
  poscafee/server:latest
```

## Environment Variables

- `DATABASE_URL`: PostgreSQL connection string (required)
- `JWT_SECRET`: Secret key for JWT token signing (required)
- `JWT_EXPIRY`: JWT token expiry duration (default: 24h)
- `PORT`: Server port (default: 8080)
- `ENVIRONMENT`: Environment mode (default: development)

## Docker Compose Example

```yaml
version: '3.8'

services:
  app:
    image: poscafee/server:latest
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgresql://user:password@host:port/database?sslmode=require
      - JWT_SECRET=your-super-secret-jwt-key
      - PORT=8080
    restart: unless-stopped
```

## Multi-Platform Support

This image is built for multiple architectures. If needed, explicitly pull for your platform:

```bash
docker run --platform linux/amd64 poscafee/server:latest
```

## Versioning

- `latest`: Latest stable release
- `vX.Y.Z`: Specific version (e.g., v1.0.0)

## Architecture

Built for Linux AMD64 architecture. The image uses Alpine Linux for a minimal footprint and runs as a non-root user for security.

## Health Check

The container includes a built-in health check that verifies the server is responding on the `/health` endpoint.

## License

This project is licensed under the terms specified in the original repository.