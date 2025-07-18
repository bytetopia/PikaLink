# Docker build script for PikaLink

param(
    [string]$Tag = "pikalink:latest",
    [switch]$Push = $false,
    [string]$Registry = "",
    [switch]$NoBuildCache = $false
)

Write-Host "Building PikaLink Docker image..." -ForegroundColor Cyan

# Check if we're in the right directory
if (-not (Test-Path "Dockerfile")) {
    Write-Host "ERROR: Dockerfile not found. Please run this script from the project root directory" -ForegroundColor Red
    exit 1
}

# Build arguments
$BuildArgs = @()
if ($NoBuildCache) {
    $BuildArgs += "--no-cache"
}

# Build the image
Write-Host "Building Docker image with tag: $Tag" -ForegroundColor Blue
$BuildCommand = "docker build $($BuildArgs -join ' ') -t $Tag ."
Write-Host "Executing: $BuildCommand" -ForegroundColor Gray

Invoke-Expression $BuildCommand

if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Docker build failed" -ForegroundColor Red
    exit 1
}

Write-Host "SUCCESS: Docker image built successfully!" -ForegroundColor Green

# Push to registry if requested
if ($Push -and $Registry) {
    $FullTag = "$Registry/$Tag"
    Write-Host "Tagging image for registry: $FullTag" -ForegroundColor Blue
    docker tag $Tag $FullTag
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Pushing to registry: $Registry" -ForegroundColor Blue
        docker push $FullTag
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "SUCCESS: Image pushed to registry!" -ForegroundColor Green
        } else {
            Write-Host "ERROR: Failed to push image to registry" -ForegroundColor Red
            exit 1
        }
    } else {
        Write-Host "ERROR: Failed to tag image for registry" -ForegroundColor Red
        exit 1
    }
} elseif ($Push) {
    Write-Host "WARNING: Push requested but no registry specified" -ForegroundColor Yellow
}

Write-Host "Image ready: $Tag" -ForegroundColor Green
Write-Host "To run: docker run -p 8080:8080 $Tag" -ForegroundColor Green