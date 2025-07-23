#!/bin/bash

# Automated publishing script for PikaLink

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
REGISTRY="dockerhub"
LATEST=true

# Function to display help
show_help() {
    echo "Automated publishing script for PikaLink"
    echo
    echo "Usage: $0 -v VERSION -u USERNAME [OPTIONS]"
    echo
    echo "Required options:"
    echo "  -v, --version VERSION    Version to publish"
    echo "  -u, --username USERNAME  Username for registry"
    echo
    echo "Optional:"
    echo "  -r, --registry REGISTRY  Registry (dockerhub|ghcr|ecr, default: dockerhub)"
    echo "  --no-latest             Don't tag as latest"
    echo "  -h, --help              Show this help message"
    echo
    echo "Examples:"
    echo "  $0 -v 1.0.0 -u yourusername"
    echo "  $0 -v 1.0.0 -u yourusername -r ghcr"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        -u|--username)
            USERNAME="$2"
            shift 2
            ;;
        -r|--registry)
            REGISTRY="$2"
            shift 2
            ;;
        --no-latest)
            LATEST=false
            shift
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            echo -e "${RED}ERROR: Unknown option $1${NC}" >&2
            show_help
            exit 1
            ;;
    esac
done

# Check required parameters
if [[ -z "$VERSION" || -z "$USERNAME" ]]; then
    echo -e "${RED}ERROR: Version and username are required${NC}" >&2
    show_help
    exit 1
fi

# Validate registry
if [[ ! "$REGISTRY" =~ ^(dockerhub|ghcr|ecr)$ ]]; then
    echo -e "${RED}ERROR: Invalid registry. Must be dockerhub, ghcr, or ecr${NC}" >&2
    exit 1
fi

echo -e "${CYAN}Publishing PikaLink v$VERSION to $REGISTRY...${NC}"

# Build image
if ! ./docker-build.sh -t "pikalink:$VERSION" -n; then
    echo -e "${RED}ERROR: Build failed${NC}" >&2
    exit 1
fi

# Tag and push based on registry
case $REGISTRY in
    "dockerhub")
        echo -e "${CYAN}Tagging and pushing to Docker Hub...${NC}"
        docker tag "pikalink:$VERSION" "$USERNAME/pikalink:$VERSION"
        if ! docker push "$USERNAME/pikalink:$VERSION"; then
            echo -e "${RED}ERROR: Failed to push version tag${NC}" >&2
            exit 1
        fi
        
        if [[ "$LATEST" == true ]]; then
            docker tag "pikalink:$VERSION" "$USERNAME/pikalink:latest"
            if ! docker push "$USERNAME/pikalink:latest"; then
                echo -e "${RED}ERROR: Failed to push latest tag${NC}" >&2
                exit 1
            fi
        fi
        ;;
    "ghcr")
        echo -e "${CYAN}Tagging and pushing to GitHub Container Registry...${NC}"
        docker tag "pikalink:$VERSION" "ghcr.io/$USERNAME/pikalink:$VERSION"
        if ! docker push "ghcr.io/$USERNAME/pikalink:$VERSION"; then
            echo -e "${RED}ERROR: Failed to push version tag${NC}" >&2
            exit 1
        fi
        
        if [[ "$LATEST" == true ]]; then
            docker tag "pikalink:$VERSION" "ghcr.io/$USERNAME/pikalink:latest"
            if ! docker push "ghcr.io/$USERNAME/pikalink:latest"; then
                echo -e "${RED}ERROR: Failed to push latest tag${NC}" >&2
                exit 1
            fi
        fi
        ;;
    "ecr")
        echo -e "${YELLOW}ECR publishing not implemented in this script${NC}"
        echo -e "${YELLOW}Please use the manual ECR steps in the documentation${NC}"
        exit 1
        ;;
esac

echo -e "${GREEN}SUCCESS: Published pikalink:$VERSION to $REGISTRY${NC}"

# Display next steps
echo
echo -e "${GREEN}Next steps:${NC}"
case $REGISTRY in
    "dockerhub")
        echo "1. Visit https://hub.docker.com/u/$USERNAME to verify the publication"
        echo "2. Test: docker pull $USERNAME/pikalink:$VERSION"
        ;;
    "ghcr")
        echo "1. Visit https://github.com/$USERNAME?tab=packages to verify the publication"
        echo "2. Test: docker pull ghcr.io/$USERNAME/pikalink:$VERSION"
        ;;
esac
