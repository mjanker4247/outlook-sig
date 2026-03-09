#!/usr/bin/env pwsh

Write-Host "Installing prerequisites for outlook-signature project..." -ForegroundColor Green

# ── Scoop ──────────────────────────────────────────────────────────────────
if (!(Get-Command scoop -ErrorAction SilentlyContinue)) {
    Write-Host "Installing Scoop package manager..." -ForegroundColor Yellow
    Set-ExecutionPolicy RemoteSigned -Scope CurrentUser -Force
    Invoke-RestMethod get.scoop.sh | Invoke-Expression
}

scoop bucket add main    2>$null
scoop bucket add versions 2>$null
scoop bucket add extras  2>$null

# ── Core tools ─────────────────────────────────────────────────────────────
if (!(Get-Command git -ErrorAction SilentlyContinue)) {
    Write-Host "Installing Git..." -ForegroundColor Yellow
    scoop install git
}

$goVersion = (go version 2>&1) | Out-String
if (!($goVersion -match "go1\.24") -or !(Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "Installing Go 1.24.2..." -ForegroundColor Yellow
    scoop install go@1.24.2
    $env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" +
                [System.Environment]::GetEnvironmentVariable("Path","User")
}

if (!(Get-Command task -ErrorAction SilentlyContinue)) {
    Write-Host "Installing Task..." -ForegroundColor Yellow
    scoop install task
}

# ── Linting & formatting ───────────────────────────────────────────────────
if (!(Get-Command golangci-lint -ErrorAction SilentlyContinue)) {
    Write-Host "Installing golangci-lint..." -ForegroundColor Yellow
    scoop install golangci-lint
}

Write-Host "Installing goimports..." -ForegroundColor Yellow
go install golang.org/x/tools/cmd/goimports@latest

# ── Cross-compilation (Docker + fyne-cross) ────────────────────────────────
if (!(Get-Command docker -ErrorAction SilentlyContinue)) {
    Write-Host "Installing Docker Desktop..." -ForegroundColor Yellow
    scoop install extras/docker-desktop
    Write-Host "  NOTE: Start Docker Desktop and complete setup before using cross-compilation tasks." -ForegroundColor Cyan
}

Write-Host "Installing fyne-cross..." -ForegroundColor Yellow
go install github.com/fyne-io/fyne-cross@latest

# ── Signing note ───────────────────────────────────────────────────────────
Write-Host ""
Write-Host "NOTE: Code signing (task sign-windows) requires signtool.exe from the Windows SDK." -ForegroundColor Cyan
Write-Host "      Install via: Visual Studio Build Tools -> Individual Components -> Windows 10/11 SDK" -ForegroundColor Cyan

Write-Host ""
Write-Host "All prerequisites installed!" -ForegroundColor Green
Write-Host "Run 'task build' to build the project." -ForegroundColor Green
