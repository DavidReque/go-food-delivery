# PowerShell script to format Go code
# Usage: .\format.ps1

Write-Host "Formatting Go code..." -ForegroundColor Green

# Format long lines
Write-Host "1. Executing golines..." -ForegroundColor Yellow
& golines -m 120 -w --ignore-generated .

# Format code with gofumpt
Write-Host "2. Executing gofumpt..." -ForegroundColor Yellow
& gofumpt -l -w .

# Organize imports (simplified)
Write-Host "3. Executing go fmt..." -ForegroundColor Yellow
& go fmt ./...

Write-Host "¡Formatting completed!" -ForegroundColor Green
