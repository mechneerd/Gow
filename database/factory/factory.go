package factory

// Factory is an interface for model factories.
type Factory[T any] interface {
	Definition() *T
	Create(count int) ([]*T, error)
}

// BaseFactory provides utilities for factories.
type BaseFactory struct {
}

// NewBaseFactory initializes a new BaseFactory.
func NewBaseFactory() *BaseFactory {
	return &BaseFactory{}
}

