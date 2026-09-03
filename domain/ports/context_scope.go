// domain/ports/context_scope.go (100行以下 - SPEC-PRINCIPLE-001)
package ports

import "context"

type contextKey string

const scopeKey contextKey = "dozou_scope"

// 定義されたスコープ（権限）
const (
	// ScopeReadOnly は閲覧のみを許可するスコープです（LANからのリモートアクセス等）
	ScopeReadOnly = "read_only"
	
	// ScopeAdmin はすべての操作（削除、エクスプローラー起動等）を許可するスコープです（ローカルWails等）
	ScopeAdmin = "admin"
)

// WithScope は指定されたスコープを付与したコンテキストを返します
func WithScope(ctx context.Context, scope string) context.Context {
	return context.WithValue(ctx, scopeKey, scope)
}

// GetScope はコンテキストから現在のスコープを取得します。
// 設定されていない場合はデフォルトでより安全な ScopeReadOnly を返します。
func GetScope(ctx context.Context) string {
	val := ctx.Value(scopeKey)
	if scope, ok := val.(string); ok && scope != "" {
		return scope
	}
	return ScopeReadOnly
}

// IsAdmin はコンテキストが管理者権限を持っているかを判定します
func IsAdmin(ctx context.Context) bool {
	return GetScope(ctx) == ScopeAdmin
}
