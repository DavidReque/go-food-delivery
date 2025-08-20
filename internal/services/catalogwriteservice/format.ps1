# PowerShell script to format Go code
# Uso: .\format.ps1

Write-Host "Formateando código Go..." -ForegroundColor Green

# Formatear líneas largas
Write-Host "1. Ejecutando golines..." -ForegroundColor Yellow
& golines -m 120 -w --ignore-generated .

# Formatear código con gofumpt
Write-Host "2. Ejecutando gofumpt..." -ForegroundColor Yellow
& gofumpt -l -w .

# Organizar imports (simplificado)
Write-Host "3. Ejecutando go fmt..." -ForegroundColor Yellow
& go fmt ./...

Write-Host "¡Formateo completado!" -ForegroundColor Green
