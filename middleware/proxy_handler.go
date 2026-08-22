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
	path := r.URL.Path
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK); return
	}
	if strings.HasPrefix(path, "/api/jobs/") { h.serveJobAPI(w, r); return }
	if strings.HasPrefix(path, "/avatars/") || strings.HasPrefix(path, "/assets/") {
		rel := strings.TrimPrefix(strings.TrimPrefix(path, "/avatars/"), "/assets/")
		absBase, err := filepath.Abs(h.avatarDir)
		if err != nil { http.NotFound(w, r); return }
		targetPath, err := filepath.Abs(filepath.Join(absBase, filepath.Clean(rel)))
		if err != nil || !strings.HasPrefix(targetPath, absBase) { http.Error(w, "Forbidden", http.StatusForbidden); return }
		if info, err := os.Stat(targetPath); err == nil && !info.IsDir() { http.ServeFile(w, r, targetPath); return }
		h.serveDefaultAvatar(w); return
	}
	if strings.HasPrefix(path, "/media/") {
		mID := filepath.Clean(strings.TrimPrefix(path, "/media/"))
		foundPath := h.resolveMediaPath(mID)
		if foundPath != "" {
			f, err := os.Open(foundPath)
			if err == nil {
				defer f.Close()
				if fi, err := f.Stat(); err == nil && !fi.IsDir() {
					w.Header().Set("Access-Control-Allow-Origin", "*")
					w.Header().Set("Content-Type", getContentType(foundPath))
					http.ServeContent(w, r, fi.Name(), fi.ModTime(), f)
					return
				}
			}
		}
		h.serveMediaPlaceholder(w); return
	}
	if strings.HasPrefix(path, "/stash-proxy/") {
		if h.stashProxy != nil {
			r.URL.Path = strings.TrimPrefix(path, "/stash-proxy")
			if r.URL.RawPath != "" { r.URL.RawPath = strings.TrimPrefix(r.URL.RawPath, "/stash-proxy") }
			h.stashProxy.ServeHTTP(w, r); return
		}
		h.serveMediaPlaceholder(w); return
	}
	http.NotFound(w, r)
}

func (h *UnifiedHandler) resolveMediaPath(mID string) string {
	dirs := []string{h.mediaDir, `G:\Media_Storage\Influencers`, "blobs", "stash", "media_local"}
	for _, base := range dirs {
		if base == "" { continue }
		p := filepath.Join(base, mID)
		if info, err := os.Stat(p); err == nil && !info.IsDir() { return p }
		entries, err := os.ReadDir(base)
		if err != nil { continue }
		for _, e := range entries {
			if e.IsDir() {
				cand := filepath.Join(base, e.Name(), "X(Twitter)", "_assets", mID)
				if info, err := os.Stat(cand); err == nil && !info.IsDir() { return cand }
				cand2 := filepath.Join(base, e.Name(), mID)
				if info, err := os.Stat(cand2); err == nil && !info.IsDir() { return cand2 }
			}
		}
	}
	return ""
}

func (h *UnifiedHandler) serveDefaultAvatar(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "image/svg+xml"); w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="64" height="64"><rect width="100%" height="100%" fill="#334155" rx="32"/><circle cx="32" cy="24" r="12" fill="#94a3b8"/><path d="M16 54c0-8.837 7.163-16 16-16s16 7.163 16 16" fill="#94a3b8"/></svg>`))
}

func (h *UnifiedHandler) serveMediaPlaceholder(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "image/svg+xml"); w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="400" height="250"><rect width="100%" height="100%" fill="#1e293b" rx="8"/><text x="50%" y="50%" dominant-baseline="middle" text-anchor="middle" fill="#64748b" font-family="sans-serif">Attached Media Preview</text></svg>`))
}
