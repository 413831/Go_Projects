# Script PowerShell para generar claves seguras para la aplicación
# Uso: .\scripts\generate-secrets.ps1

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Generador de Claves Seguras" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""

# Generar ENCRYPTION_KEY (32 bytes en hexadecimal)
Write-Host "1. Generando ENCRYPTION_KEY (32 bytes)..." -ForegroundColor Yellow
$encryptionKeyBytes = New-Object byte[] 32
[System.Security.Cryptography.RandomNumberGenerator]::Fill($encryptionKeyBytes)
$ENCRYPTION_KEY = [System.BitConverter]::ToString($encryptionKeyBytes).Replace("-", "").ToLower()
Write-Host "   ENCRYPTION_KEY=$ENCRYPTION_KEY" -ForegroundColor Green
Write-Host ""

# Generar JWT_SECRET (64 caracteres)
Write-Host "2. Generando JWT_SECRET (64 caracteres)..." -ForegroundColor Yellow
$jwtSecretBytes = New-Object byte[] 48
[System.Security.Cryptography.RandomNumberGenerator]::Fill($jwtSecretBytes)
$JWT_SECRET = [Convert]::ToBase64String($jwtSecretBytes).Substring(0, 64)
Write-Host "   JWT_SECRET=$JWT_SECRET" -ForegroundColor Green
Write-Host ""

Write-Host "==========================================" -ForegroundColor Cyan
Write-Host "Agrega estas variables a tu archivo .env:" -ForegroundColor Cyan
Write-Host "==========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "ENCRYPTION_KEY=$ENCRYPTION_KEY" -ForegroundColor White
Write-Host "JWT_SECRET=$JWT_SECRET" -ForegroundColor White
Write-Host ""
Write-Host "IMPORTANTE: Guarda estas claves de forma segura." -ForegroundColor Red
Write-Host "No las compartas ni las subas al repositorio." -ForegroundColor Red
