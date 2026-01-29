package container

import (
	"fmt"
	"reflect"
	"sync"
)

// Container define la interfaz para el contenedor de inyección de dependencias
type Container interface {
	Register(name string, factory interface{}) error
	RegisterSingleton(name string, factory interface{}) error
	Get(name string) (interface{}, error)
	GetAs(name string, target interface{}) error
	Resolve(target interface{}) error
}

// DIContainer implementa un contenedor de inyección de dependencias simple
type DIContainer struct {
	factories    map[string]interface{}
	singletons   map[string]interface{}
	singletonSet map[string]bool
	mu           sync.RWMutex
}

// NewDIContainer crea un nuevo contenedor de inyección de dependencias
func NewDIContainer() Container {
	return &DIContainer{
		factories:    make(map[string]interface{}),
		singletons:   make(map[string]interface{}),
		singletonSet: make(map[string]bool),
	}
}

// Register registra una fábrica para un tipo
func (c *DIContainer) Register(name string, factory interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.validateFactory(factory); err != nil {
		return fmt.Errorf("fábrica inválida para %s: %w", name, err)
	}

	c.factories[name] = factory
	return nil
}

// RegisterSingleton registra una fábrica para un singleton
func (c *DIContainer) RegisterSingleton(name string, factory interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.validateFactory(factory); err != nil {
		return fmt.Errorf("fábrica inválida para %s: %w", name, err)
	}

	c.factories[name] = factory
	c.singletonSet[name] = true
	return nil
}

// Get obtiene una instancia del tipo registrado
func (c *DIContainer) Get(name string) (interface{}, error) {
	c.mu.RLock()
	
	// Verificar si es singleton y ya fue creado
	if c.singletonSet[name] {
		if instance, exists := c.singletons[name]; exists {
			c.mu.RUnlock()
			return instance, nil
		}
	}
	
	factory, exists := c.factories[name]
	if !exists {
		c.mu.RUnlock()
		return nil, fmt.Errorf("no se encontró fábrica para %s", name)
	}
	c.mu.RUnlock()

	// Ejecutar fábrica
	result, err := c.executeFactory(factory)
	if err != nil {
		return nil, fmt.Errorf("error ejecutando fábrica para %s: %w", name, err)
	}

	// Guardar singleton si corresponde
	c.mu.Lock()
	if c.singletonSet[name] {
		c.singletons[name] = result
	}
	c.mu.Unlock()

	return result, nil
}

// GetAs obtiene una instancia y la asigna al target
func (c *DIContainer) GetAs(name string, target interface{}) error {
	instance, err := c.Get(name)
	if err != nil {
		return err
	}

	return c.assignToTarget(instance, target)
}

// Resolve resuelve dependencias automáticamente usando reflection
func (c *DIContainer) Resolve(target interface{}) error {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.IsNil() {
		return fmt.Errorf("target debe ser un puntero no nulo")
	}

	targetType := targetValue.Elem().Type()
	
	// Buscar fábrica por tipo
	for name, factory := range c.factories {
		factoryType := reflect.TypeOf(factory)
		if factoryType.Out(0) == targetType {
			return c.GetAs(name, target)
		}
	}

	return fmt.Errorf("no se encontró fábrica para el tipo %s", targetType)
}

// validateFactory valida que la fábrica tenga la firma correcta
func (c *DIContainer) validateFactory(factory interface{}) error {
	factoryType := reflect.TypeOf(factory)
	
	if factoryType.Kind() != reflect.Func {
		return fmt.Errorf("la fábrica debe ser una función")
	}
	
	if factoryType.NumOut() != 2 {
		return fmt.Errorf("la fábrica debe retornar (interface{}, error)")
	}
	
	if !factoryType.Out(0).Implements(reflect.TypeOf((*interface{})(nil)).Elem()) {
		return fmt.Errorf("el primer retorno debe ser interface{}")
	}
	
	errorType := reflect.TypeOf((*error)(nil)).Elem()
	if !factoryType.Out(1).Implements(errorType) {
		return fmt.Errorf("el segundo retorno debe ser error")
	}
	
	return nil
}

// executeFactory ejecuta una fábrica
func (c *DIContainer) executeFactory(factory interface{}) (interface{}, error) {
	factoryValue := reflect.ValueOf(factory)
	
	// Preparar argumentos (por ahora solo soportamos fábricas sin argumentos)
	var args []reflect.Value
	
	results := factoryValue.Call(args)
	
	if !results[1].IsNil() {
		return nil, results[1].Interface().(error)
	}
	
	return results[0].Interface(), nil
}

// assignToTarget asigna una instancia al target usando reflection
func (c *DIContainer) assignToTarget(instance interface{}, target interface{}) error {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.IsNil() {
		return fmt.Errorf("target debe ser un puntero no nulo")
	}

	instanceValue := reflect.ValueOf(instance)
	targetType := targetValue.Elem().Type()
	
	if !instanceValue.Type().ConvertibleTo(targetType) {
		return fmt.Errorf("no se puede convertir %s a %s", instanceValue.Type(), targetType)
	}
	
	targetValue.Elem().Set(instanceValue.Convert(targetType))
	return nil
}

// ServiceLocator proporciona acceso global al contenedor
type ServiceLocator struct {
	container Container
}

var (
	defaultLocator *ServiceLocator
	locatorOnce    sync.Once
)

// GetDefaultLocator retorna el localizador de servicios por defecto
func GetDefaultLocator() *ServiceLocator {
	locatorOnce.Do(func() {
		defaultLocator = &ServiceLocator{
			container: NewDIContainer(),
		}
	})
	return defaultLocator
}

// SetDefaultContainer establece el contenedor por defecto
func SetDefaultContainer(container Container) {
	GetDefaultLocator().container = container
}

// GetContainer retorna el contenedor por defecto
func GetContainer() Container {
	return GetDefaultLocator().container
}

// Register registra un servicio en el contenedor por defecto
func Register(name string, factory interface{}) error {
	return GetContainer().Register(name, factory)
}

// RegisterSingleton registra un singleton en el contenedor por defecto
func RegisterSingleton(name string, factory interface{}) error {
	return GetContainer().RegisterSingleton(name, factory)
}

// Get obtiene un servicio del contenedor por defecto
func Get(name string) (interface{}, error) {
	return GetContainer().Get(name)
}

// GetAs obtiene un servicio y lo asigna al target
func GetAs(name string, target interface{}) error {
	return GetContainer().GetAs(name, target)
}

// Resolve resuelve dependencias automáticamente
func Resolve(target interface{}) error {
	return GetContainer().Resolve(target)
}
