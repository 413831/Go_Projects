package domain

import "errors"

var (
	ErrItemNotFound      = errors.New("item no encontrado")
	ErrInvalidItemName   = errors.New("el nombre del item es requerido")
	ErrInvalidItemPrice  = errors.New("el precio del item debe ser mayor o igual a cero")
	ErrInvalidItemStock  = errors.New("el stock del item debe ser mayor o igual a cero")
	ErrItemAlreadyExists = errors.New("el item ya existe")
)
