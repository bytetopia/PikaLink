# Docker Publishing Guide for PikaLink

This guide walks you through the process of publishing your PikaLink Docker image to various container registries.

## Prerequisites

- Docker installed and running
- Docker Hub account (for Docker Hub publishing)
- GitHub account with Personal Access Token (for GitHub Container Registry)
- PowerShell (Windows) or Bash (Linux/macOS)

## Quick Start

```powershell
# Build and test locally
.\docker-build.ps1

# Test the application
docker run -p 8080:8080 pikalink:latest
# Visit http://localhost:8080/admin to verify it works

# Stop the container
docker stop $(docker ps -q --filter ancestor=pikalink:latest)
```

## Publishing to Docker Hub

### Step 1: Create Docker Hub Repository

1. Go to [Docker Hub](https://hub.docker.com)
2. Sign in to your account
3. Click "Create Repository"
4. Set repository name: `pikalink`
5. Choose visibility (Public or Private)
6. Click "Create"

### Step 2: Build Production Image

```powershell
# Build optimized production image
.\docker-build.ps1 -Tag "pikalink:v1.0.0" -NoBuildCache

# Verify the image was created
docker images | findstr pikalink
```

### Step 3: Login to Docker Hub

```powershell
# Login to Docker Hub
docker login

# Enter your Docker Hub username and password when prompted
```

### Step 4: Tag and Push

```powershell
# Replace 'yourusername' with your actual Docker Hub username
$DOCKER_USERNAME = "yourusername"

# Tag for Docker Hub
docker tag pikalink:v1.0.0 ${DOCKER_USERNAME}/pikalink:v1.0.0
docker tag pikalink:v1.0.0 ${DOCKER_USERNAME}/pikalink:latest

# Push to Docker Hub
docker push ${DOCKER_USERNAME}/pikalink:v1.0.0
docker push ${DOCKER_USERNAME}/pikalink:latest
```

### Step 5: Verify Publication

1. Go to your Docker Hub repository page
2. Verify both tags (`v1.0.0` and `latest`) are visible
3. Test pulling the image:

```powershell
# Remove local image to test pull
docker rmi ${DOCKER_USERNAME}/pikalink:latest

# Pull from Docker Hub
docker pull ${DOCKER_USERNAME}/pikalink:latest

# Test run
docker run -d -p 8080:8080 --name pikalink-test ${DOCKER_USERNAME}/pikalink:latest
```

## Publishing to GitHub Container Registry (GHCR)

### Step 1: Create Personal Access Token

1. Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Click "Generate new token (classic)"
3. Select scopes: `write:packages`, `read:packages`, `delete:packages`
4. Generate and copy the token

### Step 2: Login to GHCR

```powershell
# Set your GitHub username and token
$GITHUB_USERNAME = "yourusername"
$GITHUB_TOKEN = "your_personal_access_token"

# Login to GitHub Container Registry
echo $GITHUB_TOKEN | docker login ghcr.io -u $GITHUB_USERNAME --password-stdin
```

### Step 3: Tag and Push to GHCR

```powershell
# Tag for GitHub Container Registry
docker tag pikalink:v1.0.0 ghcr.io/${GITHUB_USERNAME}/pikalink:v1.0.0
docker tag pikalink:v1.0.0 ghcr.io/${GITHUB_USERNAME}/pikalink:latest

# Push to GHCR
docker push ghcr.io/${GITHUB_USERNAME}/pikalink:v1.0.0
docker push ghcr.io/${GITHUB_USERNAME}/pikalink:latest
```

### Step 4: Make Package Public (Optional)

1. Go to your GitHub profile → Packages
2. Find the `pikalink` package
3. Click on it → Package settings
4. Change visibility to Public if desired

## Publishing to AWS ECR

### Step 1: Setup AWS CLI

```powershell
# Install AWS CLI if not already installed
# Download from: https://aws.amazon.com/cli/

# Configure AWS credentials
aws configure
```

### Step 2: Create ECR Repository

```powershell
# Create ECR repository
aws ecr create-repository --repository-name pikalink --region us-east-1

# Get login token and login
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 123456789012.dkr.ecr.us-east-1.amazonaws.com
```

### Step 3: Tag and Push to ECR

```powershell
# Replace with your actual AWS account ID and region
$AWS_ACCOUNT_ID = "123456789012"
$AWS_REGION = "us-east-1"

# Tag for ECR
docker tag pikalink:v1.0.0 ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/pikalink:v1.0.0
docker tag pikalink:v1.0.0 ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/pikalink:latest

# Push to ECR
docker push ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/pikalink:v1.0.0
docker push ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/pikalink:latest
```

## Automated Publishing Script

Create an automated publishing script for easier deployment:

```powershell
# Save as: publish-docker.ps1
param(
    [Parameter(Mandatory=$true)]
    [string]$Version,
    
    [Parameter(Mandatory=$true)]
    [string]$Username,
    
    [ValidateSet("dockerhub", "ghcr", "ecr")]
    [string]$Registry = "dockerhub",
    
    [switch]$Latest = $true
)

Write-Host "Publishing PikaLink v$Version to $Registry..." -ForegroundColor Cyan

# Build image
.\docker-build.ps1 -Tag "pikalink:$Version" -NoBuildCache

if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Build failed" -ForegroundColor Red
    exit 1
}

# Tag and push based on registry
switch ($Registry) {
    "dockerhub" {
        docker tag "pikalink:$Version" "$Username/pikalink:$Version"
        docker push "$Username/pikalink:$Version"
        
        if ($Latest) {
            docker tag "pikalink:$Version" "$Username/pikalink:latest"
            docker push "$Username/pikalink:latest"
        }
    }
    "ghcr" {
        docker tag "pikalink:$Version" "ghcr.io/$Username/pikalink:$Version"
        docker push "ghcr.io/$Username/pikalink:$Version"
        
        if ($Latest) {
            docker tag "pikalink:$Version" "ghcr.io/$Username/pikalink:latest"
            docker push "ghcr.io/$Username/pikalink:latest"
        }
    }
    "ecr" {
        # ECR implementation would go here
        Write-Host "ECR publishing not implemented in this script" -ForegroundColor Yellow
    }
}

Write-Host "SUCCESS: Published pikalink:$Version to $Registry" -ForegroundColor Green
```

Usage:
```powershell
# Publish to Docker Hub
.\publish-docker.ps1 -Version "1.0.0" -Username "yourusername" -Registry "dockerhub"

# Publish to GitHub Container Registry
.\publish-docker.ps1 -Version "1.0.0" -Username "yourusername" -Registry "ghcr"
```

## User Deployment Instructions

Once published, users can deploy PikaLink using:

### Docker Run

```bash
# Using Docker Hub
docker run -d \
  --name pikalink \
  -p 8080:8080 \
  -v pikalink_data:/app/data \
  --restart unless-stopped \
  yourusername/pikalink:latest

# Using GitHub Container Registry
docker run -d \
  --name pikalink \
  -p 8080:8080 \
  -v pikalink_data:/app/data \
  --restart unless-stopped \
  ghcr.io/yourusername/pikalink:latest
```

### Docker Compose

Create a `docker-compose.yml` for users:

```yaml
version: '3.8'

services:
  pikalink:
    image: yourusername/pikalink:latest
    container_name: pikalink
    ports:
      - "8080:8080"
    volumes:
      - pikalink_data:/app/data
    restart: unless-stopped
    environment:
      - DB_PATH=/app/data/pikalink.db
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/admin"]
      interval: 30s
      timeout: 10s
      retries: 3

volumes:
  pikalink_data:
    driver: local
```

Then run:
```bash
docker-compose up -d
```

## Best Practices

1. **Versioning**: Always tag with semantic versions (e.g., `v1.0.0`)
2. **Security**: Regularly update base images and dependencies
3. **Size Optimization**: Use multi-stage builds and Alpine Linux
4. **Health Checks**: Include health checks for better monitoring
5. **Documentation**: Update README with deployment instructions
6. **Testing**: Test pulled images before marking as latest
7. **Automation**: Consider GitHub Actions for automatic publishing

## Troubleshooting

### Common Issues

**Build Fails:**
```powershell
# Clear Docker cache
docker system prune -a

# Rebuild without cache
.\docker-build.ps1 -NoBuildCache
```

**Push Fails:**
```powershell
# Check if logged in
docker info | findstr Username

# Re-login if needed
docker login
```

**Image Too Large:**
- Check `.dockerignore` is properly configured
- Use multi-stage builds
- Remove unnecessary dependencies

**Health Check Fails:**
```powershell
# Test health check manually
docker exec -it container_name wget --quiet --tries=1 --spider http://localhost:8080/admin
```

## Security Considerations

1. **Use specific base image tags** instead of `latest`
2. **Scan images for vulnerabilities**:
   ```powershell
   docker scan yourusername/pikalink:latest
   ```
3. **Use secrets for sensitive data** instead of environment variables
4. **Run as non-root user** (already implemented in Dockerfile)
5. **Keep images updated** with latest security patches

## Next Steps

After publishing:

1. Update your project README with deployment instructions
2. Create GitHub releases with corresponding Docker tags
3. Set up automated builds with GitHub Actions
4. Consider setting up monitoring and logging
5. Plan for database backups and migrations

For more advanced deployment scenarios, consider using Kubernetes, Docker Swarm, or cloud-specific container services.