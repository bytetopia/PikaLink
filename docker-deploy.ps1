# Docker deployment script for PikaLink

param(
    [ValidateSet("up", "down", "restart", "logs", "build")]
    [string]$Action = "up",
    [switch]$Detached = $true,
    [switch]$Build = $false,
    [string]$Service = ""
)

Write-Host "PikaLink Docker Deployment" -ForegroundColor Cyan

# Check if docker-compose.yml exists
if (-not (Test-Path "docker-compose.yml")) {
    Write-Host "ERROR: docker-compose.yml not found" -ForegroundColor Red
    exit 1
}

# Build compose command
$ComposeArgs = @()
if ($Build) {
    $ComposeArgs += "--build"
}
if ($Detached -and $Action -eq "up") {
    $ComposeArgs += "-d"
}
if ($Service) {
    $ComposeArgs += $Service
}

# Execute action
switch ($Action) {
    "up" {
        Write-Host "Starting PikaLink services..." -ForegroundColor Blue
        $Command = "docker-compose up $($ComposeArgs -join ' ')"
    }
    "down" {
        Write-Host "Stopping PikaLink services..." -ForegroundColor Blue
        $Command = "docker-compose down"
    }
    "restart" {
        Write-Host "Restarting PikaLink services..." -ForegroundColor Blue
        $Command = "docker-compose restart $Service"
    }
    "logs" {
        Write-Host "Showing PikaLink logs..." -ForegroundColor Blue
        $Command = "docker-compose logs -f $Service"
    }
    "build" {
        Write-Host "Building PikaLink services..." -ForegroundColor Blue
        $Command = "docker-compose build $Service"
    }
}

Write-Host "Executing: $Command" -ForegroundColor Gray
Invoke-Expression $Command

if ($LASTEXITCODE -eq 0) {
    switch ($Action) {
        "up" {
            Write-Host "SUCCESS: PikaLink is running!" -ForegroundColor Green
            Write-Host "Access the application at: http://localhost:8080/admin" -ForegroundColor Green
        }
        "down" {
            Write-Host "SUCCESS: PikaLink services stopped" -ForegroundColor Green
        }
        default {
            Write-Host "SUCCESS: Action completed" -ForegroundColor Green
        }
    }
} else {
    Write-Host "ERROR: Command failed" -ForegroundColor Red
    exit 1
}