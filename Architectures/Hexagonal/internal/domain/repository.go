package domain

import "context"

// ItemRepository define el puerto (interfaz) para el repositorio de items
// Este es el contrato que debe cumplir cualquier adaptador de salida
type ItemRepository interface {
	// GetAll obtiene todos los items
	GetAll(ctx context.Context) ([]*Item, error)
	
	// GetByID obtiene un item por su ID
	GetByID(ctx context.Context, id string) (*Item, error)
	
	// Create crea un nuevo item
	Create(ctx context.Context, item *Item) (*Item, error)
	
	// Update actualiza un item existente
	Update(ctx context.Context, id string, item *Item) (*Item, error)
	
	// Delete elimina un item por su ID
	Delete(ctx context.Context, id string) error
}
