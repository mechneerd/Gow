package container

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

var (
	ErrFrozen = errors.New("container is frozen and cannot be mutated")
	ErrBindingNotFound = errors.New("binding not found")
	ErrInvalidFactory = errors.New("invalid factory function")
)

// Container is the IoC container.
type Container struct {
	mu        sync.RWMutex
	bindings  map[reflect.Type]*binding
	instances map[reflect.Type]any
	aliases   map[string]reflect.Type
	frozen    bool
}

type binding struct {
	factory   any
	singleton bool
}

// New creates a new Container.
func New() *Container {
	return &Container{
		bindings:  make(map[reflect.Type]*binding),
		instances: make(map[reflect.Type]any),
		aliases:   make(map[string]reflect.Type),
	}
}

// Freeze locks the container from further mutations.
func (c *Container) Freeze() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.frozen = true
}

// Bind registers a binding with the container.
func (c *Container) Bind(iface any, factory any) error {
	return c.register(iface, factory, false)
}

// Singleton registers a singleton binding with the container.
func (c *Container) Singleton(iface any, factory any) error {
	return c.register(iface, factory, true)
}

// Instance registers an existing instance as a singleton.
func (c *Container) Instance(iface any, instance any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.frozen {
		return ErrFrozen
	}

	typ := reflect.TypeOf(iface)
	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Interface {
		typ = typ.Elem()
	}

	c.instances[typ] = instance
	return nil
}

func (c *Container) register(iface any, factory any, singleton bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.frozen {
		return ErrFrozen
	}

	typ := reflect.TypeOf(iface)
	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Interface {
		typ = typ.Elem() // use the underlying interface type
	}

	factoryType := reflect.TypeOf(factory)
	if factoryType.Kind() != reflect.Func {
		return ErrInvalidFactory
	}

	c.bindings[typ] = &binding{
		factory:   factory,
		singleton: singleton,
	}

	return nil
}

// Make resolves a dependency from the container type-safely.
func Make[T any](c *Container) (T, error) {
	var target T
	targetType := reflect.TypeOf((*T)(nil)).Elem()

	instance, err := c.Resolve(targetType)
	if err != nil {
		return target, err
	}

	return instance.(T), nil
}

// Resolve resolves a dependency by its reflection type.
func (c *Container) Resolve(typ reflect.Type) (any, error) {
	c.mu.RLock()
	// Check instances first
	if instance, ok := c.instances[typ]; ok {
		c.mu.RUnlock()
		return instance, nil
	}

	// Check bindings
	b, ok := c.bindings[typ]
	c.mu.RUnlock()

	if !ok {
		// Attempt to resolve struct if it's concrete type
		if typ.Kind() == reflect.Struct {
			return c.build(typ)
		}
		return nil, fmt.Errorf("%w: %v", ErrBindingNotFound, typ)
	}

	// Build via factory
	instance, err := c.callFactory(b.factory)
	if err != nil {
		return nil, err
	}

	if b.singleton {
		c.mu.Lock()
		// Double-check instance wasn't created while acquiring lock
		if inst, ok := c.instances[typ]; ok {
			c.mu.Unlock()
			return inst, nil
		}
		c.instances[typ] = instance
		c.mu.Unlock()
	}

	return instance, nil
}

func (c *Container) callFactory(factory any) (any, error) {
	v := reflect.ValueOf(factory)
	t := v.Type()

	in := make([]reflect.Value, t.NumIn())
	for i := 0; i < t.NumIn(); i++ {
		paramType := t.In(i)
		
		// If the factory expects the container itself
		if paramType == reflect.TypeOf(c) {
			in[i] = reflect.ValueOf(c)
			continue
		}

		paramInstance, err := c.Resolve(paramType)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve parameter %v: %w", paramType, err)
		}
		in[i] = reflect.ValueOf(paramInstance)
	}

	out := v.Call(in)
	if len(out) == 0 {
		return nil, errors.New("factory returned no values")
	}

	// If factory returns (instance, error)
	if len(out) == 2 && !out[1].IsNil() {
		return nil, out[1].Interface().(error)
	}

	return out[0].Interface(), nil
}

// build attempts to instantiate a struct by injecting dependencies into its fields.
func (c *Container) build(typ reflect.Type) (any, error) {
	ptr := reflect.New(typ)
	val := ptr.Elem()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		// Only inject into tagged fields or all exported fields? 
		// For Laravel, constructor injection is preferred. In Go, struct injection is common.
		// Let's use `inject` tag as explicit marker for struct injection, 
		// otherwise we rely on constructors (factories).
		if field.Tag.Get("inject") != "" {
			dep, err := c.Resolve(field.Type)
			if err != nil {
				return nil, fmt.Errorf("failed to inject field %s: %w", field.Name, err)
			}
			if val.Field(i).CanSet() {
				val.Field(i).Set(reflect.ValueOf(dep))
			}
		}
	}

	return ptr.Interface(), nil
}
