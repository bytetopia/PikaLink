# Docker Publishing Guide for PikaLink

> ⚠️ Internal use only. You don't need to publish this to bytetopia/PikaLink docker hub repo.

## Automated Publishing Script

Automated publishing script for easier deployment:

Usage:
```bash

chmod +x publish-docker.sh

docker login -u "yourusername"

./publish-docker.sh -v "1.0.0" -u "yourusername"

```

---

## Manual Docker Build & Publish details

### Build image

```bash
# Build and test locally
./docker-build.sh

# Test the application
docker run -p 8080:8080 pikalink:latest
# Visit http://localhost:8080/admin/ to verify it works

# Stop the container
docker stop $(docker ps -q --filter ancestor=pikalink:latest)
```

### Publishing to Docker Hub

#### Step 1: Create Docker Hub Repository

1. Go to [Docker Hub](https://hub.docker.com)
2. Sign in to your account
3. Click "Create Repository"
4. Set repository name: `pikalink`
5. Choose visibility (Public or Private)
6. Click "Create"

#### Step 2: Build Production Image

```bash
# Build optimized production image
./docker-build.sh -t "pikalink:v1.0.0" -n

# Verify the image was created
docker images | grep pikalink
```

#### Step 3: Login to Docker Hub

```bash
# Login to Docker Hub
docker login -u "yourusername"

# Enter your Docker Hub username and password when prompted
```

#### Step 4: Tag and Push

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

#### Step 5: Verify Publication

1. Go to your Docker Hub repository page
2. Verify both tags (`v1.0.0` and `latest`) are visible
3. Test pulling the image:

```bash
# Remove local image to test pull
docker rmi ${DOCKER_USERNAME}/pikalink:latest

# Pull from Docker Hub
docker pull ${DOCKER_USERNAME}/pikalink:latest

# Test run
docker run -d -p 8080:8080 --name pikalink-test ${DOCKER_USERNAME}/pikalink:latest
```

## User Deployment Instructions

Once published, users can deploy PikaLink using:


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

```bash
# Clear Docker cache
docker system prune -a

# Rebuild without cache
./docker-build.sh -n
```

**Push Fails:**

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