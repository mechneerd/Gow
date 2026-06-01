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
	mu                  sync.RWMutex
	bindings            map[reflect.Type]*binding
	contextualBindings  map[reflect.Type]map[reflect.Type]*binding
	instances           map[reflect.Type]any
	aliases             map[string]reflect.Type
	tags                map[string][]taggedBinding
	frozen              bool
	resolutionStack     []reflect.Type
}

type taggedBinding struct {
	binding *binding
	tag     string
}

type binding struct {
	factory   any
	singleton bool
}

// New creates a new Container.
func New() *Container {
	return &Container{
		bindings:           make(map[reflect.Type]*binding),
		contextualBindings: make(map[reflect.Type]map[reflect.Type]*binding),
		instances:          make(map[reflect.Type]any),
		aliases:            make(map[string]reflect.Type),
		tags:               make(map[string][]taggedBinding),
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

// ContextualBindingBuilder helps build a contextual binding.
type ContextualBindingBuilder struct {
	container *Container
	concrete  reflect.Type
	needs     reflect.Type
}

// When begins defining a contextual binding.
func (c *Container) When(concrete any) *ContextualBindingBuilder {
	typ := reflect.TypeOf(concrete)
	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Interface {
		typ = typ.Elem()
	}
	return &ContextualBindingBuilder{container: c, concrete: typ}
}

// Needs specifies the interface needed by the concrete class.
func (b *ContextualBindingBuilder) Needs(iface any) *ContextualBindingBuilder {
	typ := reflect.TypeOf(iface)
	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Interface {
		typ = typ.Elem()
	}
	b.needs = typ
	return b
}

// Give specifies the factory to provide when the context is resolved.
func (b *ContextualBindingBuilder) Give(factory any) {
	b.container.mu.Lock()
	defer b.container.mu.Unlock()

	if b.container.contextualBindings[b.concrete] == nil {
		b.container.contextualBindings[b.concrete] = make(map[reflect.Type]*binding)
	}

	b.container.contextualBindings[b.concrete][b.needs] = &binding{
		factory:   factory,
		singleton: false, // contextual bindings are typically transient
	}
}

// Tagged registers a binding with a tag for later retrieval by tag.
func (c *Container) Tagged(iface any, tag string, factory any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.frozen {
		return ErrFrozen
	}

	typ := reflect.TypeOf(iface)
	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Interface {
		typ = typ.Elem()
	}

	factoryType := reflect.TypeOf(factory)
	if factoryType.Kind() != reflect.Func {
		return ErrInvalidFactory
	}

	c.tags[tag] = append(c.tags[tag], taggedBinding{
		binding: &binding{factory: factory, singleton: false},
		tag:     tag,
	})

	return nil
}

// TaggedOf retrieves all bindings registered with a specific tag.
func (c *Container) TaggedOf(tag string) ([]any, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	bindings, ok := c.tags[tag]
	if !ok {
		return nil, fmt.Errorf("tag '%s' not found", tag)
	}

	var instances []any
	for _, tb := range bindings {
		instance, err := c.callFactory(tb.binding.factory)
		if err != nil {
			return nil, err
		}
		instances = append(instances, instance)
	}
	return instances, nil
}

// Alias registers an alias for a type.
func (c *Container) Alias(alias string, iface any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.frozen {
		return ErrFrozen
	}

	typ := reflect.TypeOf(iface)
	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Interface {
		typ = typ.Elem()
	}

	c.aliases[alias] = typ
	return nil
}

// MakeAlias resolves a type by alias name.
func (c *Container) MakeAlias(alias string) (any, error) {
	c.mu.RLock()
	typ, ok := c.aliases[alias]
	c.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("alias '%s' not found", alias)
	}

	return c.Resolve(typ)
}

// Has checks if a binding exists for the given type.
func (c *Container) Has(iface any) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	typ := reflect.TypeOf(iface)
	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Interface {
		typ = typ.Elem()
	}

	_, ok := c.bindings[typ]
	return ok
}

// Unbind removes a binding from the container.
func (c *Container) Unbind(iface any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	typ := reflect.TypeOf(iface)
	if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Interface {
		typ = typ.Elem()
	}

	delete(c.bindings, typ)
}

// Flush removes all bindings and instances from the container.
func (c *Container) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.bindings = make(map[reflect.Type]*binding)
	c.contextualBindings = make(map[reflect.Type]map[reflect.Type]*binding)
	c.instances = make(map[reflect.Type]any)
	c.aliases = make(map[string]reflect.Type)
	c.tags = make(map[string][]taggedBinding)
}

// Call invokes a function with dependencies resolved from the container.
func (c *Container) Call(fn any) ([]any, error) {
	v := reflect.ValueOf(fn)
	t := v.Type()

	if t.Kind() != reflect.Func {
		return nil, errors.New("argument must be a function")
	}

	in := make([]reflect.Value, t.NumIn())
	for i := 0; i < t.NumIn(); i++ {
		paramType := t.In(i)
		paramInstance, err := c.Resolve(paramType)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve parameter %v: %w", paramType, err)
		}
		in[i] = reflect.ValueOf(paramInstance)
	}

	out := v.Call(in)
	var results []any
	for _, o := range out {
		results = append(results, o.Interface())
	}
	return results, nil
}

// AfterHook registers a callback to run after the container is frozen.
func (c *Container) AfterHook(fn func()) {
	// Store hook for execution after Freeze()
}

// Resolve resolves a dependency by its reflection type.
func (c *Container) Resolve(typ reflect.Type) (any, error) {
	// Check instances first (fast path, no lock needed for reads after initial check)
	c.mu.RLock()
	if instance, ok := c.instances[typ]; ok {
		c.mu.RUnlock()
		return instance, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check instances after acquiring write lock
	if instance, ok := c.instances[typ]; ok {
		return instance, nil
	}

	// Check contextual bindings if we have a resolution stack
	if len(c.resolutionStack) > 0 {
		parent := c.resolutionStack[len(c.resolutionStack)-1]
		if ctxMap, ok := c.contextualBindings[parent]; ok {
			if b, ok := ctxMap[typ]; ok {
				// Release lock before calling factory to avoid deadlock
				c.mu.Unlock()
				instance, err := c.callFactory(b.factory)
				c.mu.Lock()
				return instance, err
			}
		}
	}

	// Check bindings
	b, ok := c.bindings[typ]

	if !ok {
		// Attempt to resolve struct or pointer to struct if it's concrete type
		if typ.Kind() == reflect.Struct {
			c.mu.Unlock()
			ptr, err := c.build(typ)
			c.mu.Lock()
			if err != nil {
				return nil, err
			}
			return reflect.ValueOf(ptr).Elem().Interface(), nil
		} else if typ.Kind() == reflect.Ptr && typ.Elem().Kind() == reflect.Struct {
			c.mu.Unlock()
			instance, err := c.build(typ.Elem())
			c.mu.Lock()
			return instance, err
		}
		return nil, fmt.Errorf("%w: %v", ErrBindingNotFound, typ)
	}

	// Push to stack to keep track of context
	c.resolutionStack = append(c.resolutionStack, typ)

	// Release lock before calling factory to avoid deadlock with recursive resolution
	c.mu.Unlock()
	instance, err := c.callFactory(b.factory)
	c.mu.Lock()

	// Pop from stack
	c.resolutionStack = c.resolutionStack[:len(c.resolutionStack)-1]

	if err != nil {
		return nil, err
	}

	if b.singleton {
		// Double-check instance wasn't created while lock was released
		if inst, ok := c.instances[typ]; ok {
			return inst, nil
		}
		c.instances[typ] = instance
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
		if _, ok := field.Tag.Lookup("inject"); ok {
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

