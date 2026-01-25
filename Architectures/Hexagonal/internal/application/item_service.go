package application

import (
	"context"
	"kiosco/internal/domain"
	"time"
)

// ItemService contiene los casos de uso del dominio Item
type ItemService struct {
	repository domain.ItemRepository
}

// NewItemService crea una nueva instancia del servicio de items
func NewItemService(repository domain.ItemRepository) *ItemService {
	return &ItemService{
		repository: repository,
	}
}

// GetAllItems obtiene todos los items
func (s *ItemService) GetAllItems(ctx context.Context) ([]*domain.Item, error) {
	return s.repository.GetAll(ctx)
}

// GetItemByID obtiene un item por su ID
func (s *ItemService) GetItemByID(ctx context.Context, id string) (*domain.Item, error) {
	if id == "" {
		return nil, domain.ErrItemNotFound
	}

	item, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if item == nil {
		return nil, domain.ErrItemNotFound
	}

	return item, nil
}

// CreateItem crea un nuevo item
func (s *ItemService) CreateItem(ctx context.Context, item *domain.Item) (*domain.Item, error) {
	// Validar el item
	if err := item.Validate(); err != nil {
		return nil, err
	}

	// Establecer timestamps
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	// Crear el item
	return s.repository.Create(ctx, item)
}

// UpdateItem actualiza un item existente
func (s *ItemService) UpdateItem(ctx context.Context, id string, item *domain.Item) (*domain.Item, error) {
	if id == "" {
		return nil, domain.ErrItemNotFound
	}

	// Validar el item
	if err := item.Validate(); err != nil {
		return nil, err
	}

	// Verificar que el item existe
	existingItem, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if existingItem == nil {
		return nil, domain.ErrItemNotFound
	}

	// Actualizar timestamps
	item.UpdatedAt = time.Now()
	item.CreatedAt = existingItem.CreatedAt // Preservar fecha de creación

	return s.repository.Update(ctx, id, item)
}

// DeleteItem elimina un item por su ID
func (s *ItemService) DeleteItem(ctx context.Context, id string) error {
	if id == "" {
		return domain.ErrItemNotFound
	}

	// Verificar que el item existe
	item, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if item == nil {
		return domain.ErrItemNotFound
	}

	return s.repository.Delete(ctx, id)
}
