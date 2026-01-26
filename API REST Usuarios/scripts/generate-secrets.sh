#!/bin/bash

# Script para generar claves seguras para la aplicación
# Uso: ./scripts/generate-secrets.sh

echo "=========================================="
echo "Generador de Claves Seguras"
echo "=========================================="
echo ""

# Verificar si openssl está disponible
if ! command -v openssl &> /dev/null; then
    echo "Error: openssl no está instalado"
    echo "Por favor, instala openssl para usar este script"
    exit 1
fi

echo "1. Generando ENCRYPTION_KEY (32 bytes)..."
ENCRYPTION_KEY=$(openssl rand -hex 32)
echo "   ENCRYPTION_KEY=$ENCRYPTION_KEY"
echo ""

echo "2. Generando JWT_SECRET (64 caracteres)..."
JWT_SECRET=$(openssl rand -base64 48 | tr -d '\n' | head -c 64)
echo "   JWT_SECRET=$JWT_SECRET"
echo ""

echo "=========================================="
echo "Agrega estas variables a tu archivo .env:"
echo "=========================================="
echo ""
echo "ENCRYPTION_KEY=$ENCRYPTION_KEY"
echo "JWT_SECRET=$JWT_SECRET"
echo ""
echo "IMPORTANTE: Guarda estas claves de forma segura."
echo "No las compartas ni las subas al repositorio."
