// middleware/proxy_media.go (100行以下)
package middleware

import (
	"net/http"
	"os"
	"path/filepath"
)

func (h *UnifiedHandler) serveMediaLocal(w http.ResponseWriter, r *http.Request, relPath string) {
	relPath = filepath.Clean(relPath)
	foundPath := h.resolveMediaLocalPath(relPath)
	if foundPath != "" {
		h.serveFileContent(w, r, foundPath); return
	}
	h.serveMediaPlaceholder(w)
}

func (h *UnifiedHandler) serveMedia(w http.ResponseWriter, r *http.Request, mID string) {
	mID = filepath.Clean(mID)
	foundPath := h.resolveMediaPath(mID)
	if foundPath != "" {
		h.serveFileContent(w, r, foundPath); return
	}
	h.serveMediaPlaceholder(w)
}

func (h *UnifiedHandler) serveFileContent(w http.ResponseWriter, r *http.Request, path string) {
	f, err := os.Open(path)
	if err != nil { h.serveMediaPlaceholder(w); return }
	defer f.Close()
	fi, err := f.Stat()
	if err != nil || fi.IsDir() { h.serveMediaPlaceholder(w); return }
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", getContentType(path))
	http.ServeContent(w, r, fi.Name(), fi.ModTime(), f)
}

func (h *UnifiedHandler) resolveMediaLocalPath(rel string) string {
	if h.mediaDir == "" || rel == "" { return "" }
	p := filepath.Join(h.mediaDir, rel)
	if info, err := os.Stat(p); err == nil && !info.IsDir() { return p }
	exts := []string{".jpg", ".jpeg", ".png", ".webp", ".mp4", ".gif"}
	for _, ext := range exts {
		pExt := p + ext
		if info, err := os.Stat(pExt); err == nil && !info.IsDir() { return pExt }
	}
	return h.resolveMediaPath(filepath.Base(rel))
}

func (h *UnifiedHandler) resolveMediaPath(mID string) string {
	dirs := []string{h.mediaDir, `G:\Media_Storage\Influencers`, "blobs", "stash", "media_local"}
	for _, base := range dirs {
		if base == "" { continue }
		if hit := checkFileCandidates(base, mID); hit != "" { return hit }
		entries, err := os.ReadDir(base)
		if err != nil { continue }
		for _, e := range entries {
			if e.IsDir() {
				candDir := filepath.Join(base, e.Name(), "X(Twitter)", "_assets")
				if hit := checkFileCandidates(candDir, mID); hit != "" { return hit }
				candDir2 := filepath.Join(base, e.Name())
				if hit := checkFileCandidates(candDir2, mID); hit != "" { return hit }
			}
		}
	}
	return ""
}

func checkFileCandidates(dir, name string) string {
	if dir == "" || name == "" { return "" }
	p := filepath.Join(dir, name)
	if info, err := os.Stat(p); err == nil && !info.IsDir() { return p }
	exts := []string{".jpg", ".jpeg", ".png", ".webp", ".mp4", ".gif"}
	for _, ext := range exts {
		pExt := filepath.Join(dir, name+ext)
		if info, err := os.Stat(pExt); err == nil && !info.IsDir() { return pExt }
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
