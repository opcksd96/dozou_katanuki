package domain

import (
	"context"
	"errors"
)

// Scope defines the permission level of the current execution context
type Scope string

const (
	// ScopeReadOnly indicates the operation is invoked from a public/LAN interface
	// and should only be allowed to perform read operations.
	ScopeReadOnly Scope = "read_only"

	// ScopeAdmin indicates the operation is invoked from a privileged interface (e.g. Wails)
	// and is allowed to perform any operation.
	ScopeAdmin Scope = "admin"
)

type contextKey string

const scopeKey contextKey = "execution_scope"

var ErrPermissionDenied = errors.New("permission denied: insufficient scope")

// WithScope returns a new context with the given execution scope.
func WithScope(ctx context.Context, scope Scope) context.Context {
	return context.WithValue(ctx, scopeKey, scope)
}

// GetScope returns the execution scope from the context.
// If no scope is defined, it conservatively defaults to ScopeReadOnly.
func GetScope(ctx context.Context) Scope {
	if scope, ok := ctx.Value(scopeKey).(Scope); ok {
		return scope
	}
	return ScopeReadOnly
}

// EnsureAdmin returns an error if the context does not have the Admin scope.
func EnsureAdmin(ctx context.Context) error {
	if GetScope(ctx) != ScopeAdmin {
		return ErrPermissionDenied
	}
	return nil
}
