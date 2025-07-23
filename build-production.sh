#!/bin/bash

# Production build script for PikaLink
# Usage: ./build-production.sh your-domain.com

if [ -z "$1" ]; then
    echo "Usage: $0 <domain>"
    echo "Example: $0 mysite.com"
    echo "Example: $0 localhost:8080 (for local production testing)"
    exit 1
fi

DOMAIN=$1

# Determine protocol (use https unless localhost)
if [[ "$DOMAIN" == localhost* ]]; then
    PROTOCOL="http"
else
    PROTOCOL="https"
fi

echo "Building PikaLink for domain: $DOMAIN"
echo "Using protocol: $PROTOCOL"

# Set environment variables for frontend build
export REACT_APP_API_URL="${PROTOCOL}://${DOMAIN}/api"
export REACT_APP_SHORT_LINK_URL="${PROTOCOL}://${DOMAIN}"

echo "Frontend environment:"
echo "  REACT_APP_API_URL=$REACT_APP_API_URL"
echo "  REACT_APP_SHORT_LINK_URL=$REACT_APP_SHORT_LINK_URL"

# Build frontend
echo "Building frontend..."
cd frontend
npm run build
cd ..

# Copy frontend build to backend
echo "Copying frontend build to backend..."
rm -rf backend/frontend
cp -r frontend/build backend/frontend

# Build backend
echo "Building backend..."
cd backend
go build -o pikalink .
cd ..

echo ""
echo "Build complete! To run in production:"
echo "1. Set environment variable: export CORS_ORIGINS=\"${PROTOCOL}://${DOMAIN}\""
echo "2. Run the backend: cd backend && ./pikalink"
echo ""
echo "Or use Docker with environment variables:"
echo "docker run -e CORS_ORIGINS=\"${PROTOCOL}://${DOMAIN}\" -p 8080:8080 your-pikalink-image"
