# Multi-stage build for PikaLink
FROM node:18-alpine AS frontend-builder

# Build frontend
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci --only=production
COPY frontend/ ./
RUN npm run build

# Build backend
FROM golang:1.24-alpine AS backend-builder

# Accept version as build argument
ARG VERSION=unknown

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
# Copy frontend build from previous stage
COPY --from=frontend-builder /app/frontend/build ./frontend

# Build the Go binary
RUN CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static"' -o pikalink main.go

# Create version file
RUN echo -n "$VERSION" > version.txt

# Final stage
FROM alpine:latest

# Accept version as build argument in final stage
ARG VERSION=unknown

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy the binary and frontend files
COPY --from=backend-builder /app/backend/pikalink .
COPY --from=backend-builder /app/backend/frontend ./frontend
COPY --from=backend-builder /app/backend/version.txt .

# Create data directory for SQLite database
RUN mkdir -p /app/data

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8080/admin/ || exit 1

# Run the binary
CMD ["./pikalink"]