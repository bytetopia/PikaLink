param(
    [Parameter(Mandatory=$true)]
    [string]$Version,
    
    [Parameter(Mandatory=$true)]
    [string]$Username,
    
    [ValidateSet("dockerhub")]
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
}

Write-Host "SUCCESS: Published pikalink:$Version to $Registry" -ForegroundColor Green