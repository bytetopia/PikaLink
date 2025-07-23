# Docker Publishing Guide for PikaLink

This guide walks you through the process of publishing your PikaLink Docker image to various container registries.

## Prerequisites

- Docker installed and running
- Docker Hub account (for Docker Hub publishing)
- GitHub account with Personal Access Token (for GitHub Container Registry)
- PowerShell (Windows) or Bash (Linux/macOS)

## Build Scripts

This project includes build scripts for both Windows and Linux environments:
- `docker-build.ps1` - PowerShell script for Windows
- `docker-build.sh` - Bash script for Linux/macOS

## Quick Start

### Windows (PowerShell)
```powershell
# Build and test locally
.\docker-build.ps1

# Test the application
docker run -p 8080:8080 pikalink:latest
# Visit http://localhost:8080/admin to verify it works

# Stop the container
docker stop $(docker ps -q --filter ancestor=pikalink:latest)
```

### Linux/macOS (Bash)
```bash
# Build and test locally
./docker-build.sh

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

#### Windows (PowerShell)
```powershell
# Build optimized production image
.\docker-build.ps1 -Tag "pikalink:v1.0.0" -NoBuildCache

# Verify the image was created
docker images | findstr pikalink
```

#### Linux/macOS (Bash)
```bash
# Build optimized production image
./docker-build.sh -t "pikalink:v1.0.0" -n

# Verify the image was created
docker images | grep pikalink
```

### Step 3: Login to Docker Hub

#### Windows (PowerShell)
```powershell
# Login to Docker Hub
docker login

# Enter your Docker Hub username and password when prompted
```

#### Linux/macOS (Bash)
```bash
# Login to Docker Hub
docker login

# Enter your Docker Hub username and password when prompted
```

### Step 4: Tag and Push

#### Windows (PowerShell)
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

#### Linux/macOS (Bash)
```bash
# Replace 'yourusername' with your actual Docker Hub username
DOCKER_USERNAME="yourusername"

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

#### Windows (PowerShell)
```powershell
# Remove local image to test pull
docker rmi ${DOCKER_USERNAME}/pikalink:latest

# Pull from Docker Hub
docker pull ${DOCKER_USERNAME}/pikalink:latest

# Test run
docker run -d -p 8080:8080 --name pikalink-test ${DOCKER_USERNAME}/pikalink:latest
```

#### Linux/macOS (Bash)
```bash
# Remove local image to test pull
docker rmi ${DOCKER_USERNAME}/pikalink:latest

# Pull from Docker Hub
docker pull ${DOCKER_USERNAME}/pikalink:latest

# Test run
docker run -d -p 8080:8080 --name pikalink-test ${DOCKER_USERNAME}/pikalink:latest
```


## Automated Publishing Script

Automated publishing script for easier deployment:

Usage:
```bash
# Make script executable
chmod +x publish-docker.sh

# Publish to Docker Hub
./publish-docker.sh -v "1.0.0" -u "yourusername"

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
      - DATA_PATH=/app/data
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

Windows (PowerShell):
```powershell
# Clear Docker cache
docker system prune -a

# Rebuild without cache
.\docker-build.ps1 -NoBuildCache
```

Linux/macOS (Bash):
```bash
# Clear Docker cache
docker system prune -a

# Rebuild without cache
./docker-build.sh -n
```

**Push Fails:**

Windows (PowerShell):
```powershell
# Check if logged in
docker info | findstr Username

# Re-login if needed
docker login
```

Linux/macOS (Bash):
```bash
# Check if logged in
docker info | grep Username

# Re-login if needed
docker login
```

**Image Too Large:**
- Check `.dockerignore` is properly configured
- Use multi-stage builds
- Remove unnecessary dependencies

**Health Check Fails:**

Windows (PowerShell):
```powershell
# Test health check manually
docker exec -it container_name wget --quiet --tries=1 --spider http://localhost:8080/admin/
```

Linux/macOS (Bash):
```bash
# Test health check manually
docker exec -it container_name wget --quiet --tries=1 --spider http://localhost:8080/admin/
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