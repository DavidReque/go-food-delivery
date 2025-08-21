# PowerShell script to format Go code
# Uso: .\format.ps1

Write-Host "Formatting Go code..." -ForegroundColor Green

# Formatear líneas largas
Write-Host "1. Executing golines..." -ForegroundColor Yellow
& golines -m 120 -w --ignore-generated .

# Formatear código con gofumpt
Write-Host "2. Executing gofumpt..." -ForegroundColor Yellow
& gofumpt -l -w .

# Organizar imports (simplificado)
Write-Host "3. Executing go fmt..." -ForegroundColor Yellow
& go fmt ./...

Write-Host "¡Formatting completed!" -ForegroundColor Green
