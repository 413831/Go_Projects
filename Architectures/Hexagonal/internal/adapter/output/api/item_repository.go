package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"kiosco/internal/domain"
	"net/http"
	"time"
)

// HTTPItemRepository es el adaptador de salida que se comunica con la API externa
type HTTPItemRepository struct {
	baseURL    string
	httpClient *http.Client
}

// NewHTTPItemRepository crea una nueva instancia del repositorio HTTP
func NewHTTPItemRepository(baseURL string) *HTTPItemRepository {
	return &HTTPItemRepository{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetAll obtiene todos los items desde la API externa
func (r *HTTPItemRepository) GetAll(ctx context.Context) ([]*domain.Item, error) {
	url := fmt.Sprintf("%s/items", r.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error al crear request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al realizar request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error en la API externa (status %d): %s", resp.StatusCode, string(body))
	}
	
	var items []*domain.Item
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, fmt.Errorf("error al decodificar respuesta: %w", err)
	}
	
	return items, nil
}

// GetByID obtiene un item por su ID desde la API externa
func (r *HTTPItemRepository) GetByID(ctx context.Context, id string) (*domain.Item, error) {
	url := fmt.Sprintf("%s/items/%s", r.baseURL, id)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error al crear request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al realizar request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusNotFound {
		return nil, domain.ErrItemNotFound
	}
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error en la API externa (status %d): %s", resp.StatusCode, string(body))
	}
	
	var item domain.Item
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return nil, fmt.Errorf("error al decodificar respuesta: %w", err)
	}
	
	return &item, nil
}

// Create crea un nuevo item en la API externa
func (r *HTTPItemRepository) Create(ctx context.Context, item *domain.Item) (*domain.Item, error) {
	url := fmt.Sprintf("%s/items", r.baseURL)
	
	body, err := json.Marshal(item)
	if err != nil {
		return nil, fmt.Errorf("error al serializar item: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error al crear request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al realizar request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error en la API externa (status %d): %s", resp.StatusCode, string(body))
	}
	
	var createdItem domain.Item
	if err := json.NewDecoder(resp.Body).Decode(&createdItem); err != nil {
		return nil, fmt.Errorf("error al decodificar respuesta: %w", err)
	}
	
	return &createdItem, nil
}

// Update actualiza un item existente en la API externa
func (r *HTTPItemRepository) Update(ctx context.Context, id string, item *domain.Item) (*domain.Item, error) {
	url := fmt.Sprintf("%s/items/%s", r.baseURL, id)
	
	body, err := json.Marshal(item)
	if err != nil {
		return nil, fmt.Errorf("error al serializar item: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error al crear request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al realizar request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusNotFound {
		return nil, domain.ErrItemNotFound
	}
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error en la API externa (status %d): %s", resp.StatusCode, string(body))
	}
	
	var updatedItem domain.Item
	if err := json.NewDecoder(resp.Body).Decode(&updatedItem); err != nil {
		return nil, fmt.Errorf("error al decodificar respuesta: %w", err)
	}
	
	return &updatedItem, nil
}

// Delete elimina un item de la API externa
func (r *HTTPItemRepository) Delete(ctx context.Context, id string) error {
	url := fmt.Sprintf("%s/items/%s", r.baseURL, id)
	
	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("error al crear request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error al realizar request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusNotFound {
		return domain.ErrItemNotFound
	}
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error en la API externa (status %d): %s", resp.StatusCode, string(body))
	}
	
	return nil
}
