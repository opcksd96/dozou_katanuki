// middleware/proxy_handler.go (100行以下)
package middleware

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type UnifiedHandler struct {
	avatarDir  string
	stashProxy *httputil.ReverseProxy
}

func NewUnifiedHandler(avatarDir string, stashURL *url.URL) *UnifiedHandler {
	_ = os.MkdirAll(avatarDir, 0755)

	var proxy *httputil.ReverseProxy
	if stashURL != nil {
		proxy = httputil.NewSingleHostReverseProxy(stashURL)
		// Stash へのリクエストヘッダー調整（CORS / Host維持など）
		originalDirector := proxy.Director
		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			req.Host = stashURL.Host
		}
	}

	return &UnifiedHandler{
		avatarDir:  avatarDir,
		stashProxy: proxy,
	}
}

func (h *UnifiedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// 1. アバター画像 (/avatars/ または /assets/) へのリクエスト解決
	if strings.HasPrefix(path, "/avatars/") || strings.HasPrefix(path, "/assets/") {
		relPath := path
		if strings.HasPrefix(path, "/avatars/") {
			relPath = strings.TrimPrefix(path, "/avatars/")
		} else {
			relPath = strings.TrimPrefix(path, "/assets/")
		}

		filePath := filepath.Join(h.avatarDir, filepath.Clean(relPath))
		if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
			http.ServeFile(w, r, filePath)
			return
		}

		// ファイルが存在しない場合はデフォルトアバターSVGを返却
		h.serveDefaultAvatar(w)
		return
	}

	// 2. Stash サーバーへのリバースプロキシ (/stash-proxy/...)
	if strings.HasPrefix(path, "/stash-proxy/") {
		if h.stashProxy != nil {
			// 例: /stash-proxy/scene/123/stream -> /scene/123/stream
			r.URL.Path = strings.TrimPrefix(path, "/stash-proxy")
			if r.URL.RawPath != "" {
				r.URL.RawPath = strings.TrimPrefix(r.URL.RawPath, "/stash-proxy")
			}
			h.stashProxy.ServeHTTP(w, r)
			return
		}

		// Stashプロキシ未設定時のプレースホルダー
		h.serveMediaPlaceholder(w)
		return
	}

	http.NotFound(w, r)
}

func (h *UnifiedHandler) serveDefaultAvatar(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 64 64"><rect width="100%" height="100%" fill="#334155" rx="32"/><circle cx="32" cy="24" r="12" fill="#94a3b8"/><path d="M16 54c0-8.837 7.163-16 16-16s16 7.163 16 16" fill="#94a3b8"/></svg>`))
}

func (h *UnifiedHandler) serveMediaPlaceholder(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="400" height="250" viewBox="0 0 400 250"><rect width="100%" height="100%" fill="#1e293b" rx="8"/><text x="50%" y="50%" dominant-baseline="middle" text-anchor="middle" fill="#64748b" font-family="sans-serif" font-size="14">Attached Media Preview</text></svg>`))
}
