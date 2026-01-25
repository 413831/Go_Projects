package domain

import "time"

// Item representa la entidad de dominio de un producto del kiosco
type Item struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validate valida que los campos requeridos del Item estén presentes
func (i *Item) Validate() error {
	if i.Name == "" {
		return ErrInvalidItemName
	}
	if i.Price < 0 {
		return ErrInvalidItemPrice
	}
	if i.Stock < 0 {
		return ErrInvalidItemStock
	}
	return nil
}
