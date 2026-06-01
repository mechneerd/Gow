package support

import (
	"context"
	"sync"
)

// contextKey is a private type for context keys defined in this package.
type contextKey string

const (
	// ContextKeyRequestID is the context key for the request ID.
	ContextKeyRequestID contextKey = "request_id"
	// ContextKeyUserID is the context key for the user ID.
	ContextKeyUserID contextKey = "user_id"
	// ContextKeyTenantID is the context key for the tenant ID.
	ContextKeyTenantID contextKey = "tenant_id"
	// ContextKeyLocale is the context key for the locale.
	ContextKeyLocale contextKey = "locale"
	// ContextKeyChannel is the context key for the channel (web, api, cli, etc.).
	ContextKeyChannel contextKey = "channel"
)

// ContextPropagator manages context values across goroutines.
type ContextPropagator struct {
	mu     sync.RWMutex
	values map[contextKey]any
}

// NewContextPropagator creates a new context propagator.
func NewContextPropagator() *ContextPropagator {
	return &ContextPropagator{
		values: make(map[contextKey]any),
	}
}

// Set sets a value in the propagator.
func (cp *ContextPropagator) Set(key contextKey, value any) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	cp.values[key] = value
}

// Get gets a value from the propagator.
func (cp *ContextPropagator) Get(key contextKey) (any, bool) {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	val, ok := cp.values[key]
	return val, ok
}

// WithContext creates a new context with the propagator's values.
func (cp *ContextPropagator) WithContext(ctx context.Context) context.Context {
	cp.mu.RLock()
	defer cp.mu.RUnlock()
	for k, v := range cp.values {
		ctx = context.WithValue(ctx, k, v)
	}
	return ctx
}

// FromContext extracts values from a context into the propagator.
func (cp *ContextPropagator) FromContext(ctx context.Context) {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	for _, k := range []contextKey{ContextKeyRequestID, ContextKeyUserID, ContextKeyTenantID, ContextKeyLocale, ContextKeyChannel} {
		if v := ctx.Value(k); v != nil {
			cp.values[k] = v
		}
	}
}

// ContextWithRequestID adds a request ID to the context.
func ContextWithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, ContextKeyRequestID, requestID)
}

// ContextWithUserID adds a user ID to the context.
func ContextWithUserID(ctx context.Context, userID any) context.Context {
	return context.WithValue(ctx, ContextKeyUserID, userID)
}

// ContextWithTenantID adds a tenant ID to the context.
func ContextWithTenantID(ctx context.Context, tenantID any) context.Context {
	return context.WithValue(ctx, ContextKeyTenantID, tenantID)
}

// ContextWithLocale adds a locale to the context.
func ContextWithLocale(ctx context.Context, locale string) context.Context {
	return context.WithValue(ctx, ContextKeyLocale, locale)
}

// ContextWithChannel adds a channel to the context.
func ContextWithChannel(ctx context.Context, channel string) context.Context {
	return context.WithValue(ctx, ContextKeyChannel, channel)
}

// ContextRequestID extracts the request ID from the context.
func ContextRequestID(ctx context.Context) string {
	if v, ok := ctx.Value(ContextKeyRequestID).(string); ok {
		return v
	}
	return ""
}

// ContextUserID extracts the user ID from the context.
func ContextUserID(ctx context.Context) any {
	return ctx.Value(ContextKeyUserID)
}

// ContextTenantID extracts the tenant ID from the context.
func ContextTenantID(ctx context.Context) any {
	return ctx.Value(ContextKeyTenantID)
}

// ContextLocale extracts the locale from the context.
func ContextLocale(ctx context.Context) string {
	if v, ok := ctx.Value(ContextKeyLocale).(string); ok {
		return v
	}
	return ""
}

// ContextChannel extracts the channel from the context.
func ContextChannel(ctx context.Context) string {
	if v, ok := ctx.Value(ContextKeyChannel).(string); ok {
		return v
	}
	return ""
}
