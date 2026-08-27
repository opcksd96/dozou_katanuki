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
	mediaDir   string
	stashProxy *httputil.ReverseProxy
	jobOrch    *JobOrchestrator
}

func NewUnifiedHandler(avatarDir string, stashURL *url.URL) *UnifiedHandler {
	_ = os.MkdirAll(avatarDir, 0755)
	var proxy *httputil.ReverseProxy
	if stashURL != nil {
		proxy = httputil.NewSingleHostReverseProxy(stashURL)
		orig := proxy.Director
		proxy.Director = func(req *http.Request) { orig(req); req.Host = stashURL.Host }
		proxy.ModifyResponse = func(resp *http.Response) error {
			resp.Header.Set("Access-Control-Allow-Origin", "*")
			resp.Header.Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS, POST")
			return nil
		}
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			w.Header().Set("Content-Type", "image/svg+xml"); w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="400" height="250"><rect width="100%" height="100%" fill="#1e293b" rx="8"/><text x="50%" y="50%" dominant-baseline="middle" text-anchor="middle" fill="#ef4444" font-family="sans-serif">Stash Offline (502)</text></svg>`))
		}
	}
	return &UnifiedHandler{avatarDir: avatarDir, stashProxy: proxy}
}

func (h *UnifiedHandler) SetJobOrchestrator(orch *JobOrchestrator) { h.jobOrch = orch }
func (h *UnifiedHandler) SetMediaDir(dir string)                  { h.mediaDir = dir }

func (h *UnifiedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, HEAD")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, Range, X-Requested-With")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Range, Accept-Ranges")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK); return
	}
	path := r.URL.Path
	if strings.HasPrefix(path, "/api/jobs/") { h.serveJobAPI(w, r); return }
	if strings.HasPrefix(path, "/avatars/") || strings.HasPrefix(path, "/assets/") {
		h.serveAvatarOrAsset(w, r, path); return
	}
	if strings.HasPrefix(path, "/media-local/") {
		h.serveMediaLocal(w, r, strings.TrimPrefix(path, "/media-local/")); return
	}
	if strings.HasPrefix(path, "/media/") {
		h.serveMedia(w, r, strings.TrimPrefix(path, "/media/")); return
	}
	if strings.HasPrefix(path, "/stash-proxy/") {
		h.serveStashProxy(w, r, path); return
	}
	http.NotFound(w, r)
}

func (h *UnifiedHandler) serveAvatarOrAsset(w http.ResponseWriter, r *http.Request, path string) {
	rel := strings.TrimPrefix(strings.TrimPrefix(path, "/avatars/"), "/assets/")
	absBase, err := filepath.Abs(h.avatarDir)
	if err != nil { http.NotFound(w, r); return }
	targetPath, err := filepath.Abs(filepath.Join(absBase, filepath.Clean(rel)))
	if err != nil || !strings.HasPrefix(targetPath, absBase) {
		http.Error(w, "Forbidden", http.StatusForbidden); return
	}
	if info, err := os.Stat(targetPath); err == nil && !info.IsDir() {
		http.ServeFile(w, r, targetPath); return
	}
	h.serveDefaultAvatar(w)
}

func (h *UnifiedHandler) serveStashProxy(w http.ResponseWriter, r *http.Request, path string) {
	if h.stashProxy != nil {
		r.URL.Path = strings.TrimPrefix(path, "/stash-proxy")
		if r.URL.RawPath != "" { r.URL.RawPath = strings.TrimPrefix(r.URL.RawPath, "/stash-proxy") }
		h.stashProxy.ServeHTTP(w, r); return
	}
	h.serveMediaPlaceholder(w)
}
