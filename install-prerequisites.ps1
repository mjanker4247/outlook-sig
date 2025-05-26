#!/usr/bin/env pwsh

Write-Host "Installing prerequisites for outlook-signature project..." -ForegroundColor Green

# Check if Scoop is installed, if not install it
if (!(Get-Command scoop -ErrorAction SilentlyContinue)) {
    Write-Host "Installing Scoop package manager..." -ForegroundColor Yellow
    Set-ExecutionPolicy RemoteSigned -Scope CurrentUser -Force
    Invoke-RestMethod get.scoop.sh | Invoke-Expression
}

# Add required buckets
scoop bucket add main
scoop bucket add versions
scoop bucket add extras

# Install Git if not present
if (!(Get-Command git -ErrorAction SilentlyContinue)) {
    Write-Host "Installing Git..." -ForegroundColor Yellow
    scoop install git
}

# Install Go if not present or wrong version
$goVersion = (go version 2>&1) | Out-String
if (!($goVersion -match "go1.24.2") -or !(Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "Installing Go 1.24.2..." -ForegroundColor Yellow
    scoop install go@1.24.2
    
    # Refresh environment variables
    $env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")
}

# Install Task if not present
if (!(Get-Command task -ErrorAction SilentlyContinue)) {
    Write-Host "Installing Task..." -ForegroundColor Yellow
    scoop install task
}

Write-Host "`nAll prerequisites have been installed!" -ForegroundColor Green
Write-Host "You can now run 'task build' to build the project." -ForegroundColor Green 