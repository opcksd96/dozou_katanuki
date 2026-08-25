// middleware/broadcast_root.go (100行以下)
package middleware

import (
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	types := map[string]string{
		".js": "application/javascript; charset=utf-8", ".mjs": "application/javascript; charset=utf-8",
		".css": "text/css; charset=utf-8", ".html": "text/html; charset=utf-8",
		".json": "application/json; charset=utf-8", ".wasm": "application/wasm",
		".woff2": "font/woff2", ".svg": "image/svg+xml",
	}
	for ext, ct := range types { _ = mime.AddExtensionType(ext, ct) }
}

func findDistDir() string {
	candidates := []string{filepath.Join("frontend", "dist"), filepath.Join("..", "frontend", "dist"), filepath.Join("..", "..", "frontend", "dist")}
	if exe, err := os.Executable(); err == nil { candidates = append(candidates, filepath.Join(filepath.Dir(exe), "frontend", "dist")) }
	for _, c := range candidates {
		if fi, err := os.Stat(c); err == nil && fi.IsDir() {
			if _, err2 := os.Stat(filepath.Join(c, "index.html")); err2 == nil { return c }
		}
	}
	return filepath.Join("frontend", "dist")
}

func getContentType(name string) string {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".js", ".mjs": return "application/javascript; charset=utf-8"
	case ".css": return "text/css; charset=utf-8"
	case ".html": return "text/html; charset=utf-8"
	case ".json": return "application/json; charset=utf-8"
	case ".svg": return "image/svg+xml"
	case ".wasm": return "application/wasm"
	case ".woff2": return "font/woff2"
	case ".png": return "image/png"
	case ".jpg", ".jpeg": return "image/jpeg"
	default:
		ct := mime.TypeByExtension(filepath.Ext(name))
		if ct == "" { ct = "application/octet-stream" }
		return ct
	}
}

func serveDistFile(w http.ResponseWriter, r *http.Request, filePath string) bool {
	f, err := os.Open(filePath); if err != nil { return false }
	defer f.Close()
	fi, err := f.Stat(); if err != nil || fi.IsDir() { return false }
	w.Header().Set("Content-Type", getContentType(filePath))
	http.ServeContent(w, r, fi.Name(), fi.ModTime(), f)
	return true
}

func (s *BroadcastService) serveFromFS(w http.ResponseWriter, r *http.Request, cleanPath string) bool {
	s.mu.RLock(); dfs := s.distFS; s.mu.RUnlock()
	if dfs == nil { return false }
	for _, pfx := range []string{"", "frontend/dist/", "dist/", "frontend/"} {
		lookup := strings.TrimPrefix(filepath.ToSlash(filepath.Clean(pfx+cleanPath)), "/")
		if f, err := dfs.Open(lookup); err == nil {
			defer f.Close()
			if fi, err := f.Stat(); err == nil && !fi.IsDir() {
				w.Header().Set("Content-Type", getContentType(lookup))
				if rs, ok := f.(io.ReadSeeker); ok { http.ServeContent(w, r, fi.Name(), fi.ModTime(), rs) } else { _, _ = io.Copy(w, f) }
				return true
			}
		}
	}
	return false
}

func (s *BroadcastService) handleRoot(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.HasPrefix(path, "/plugins/") {
		http.StripPrefix("/plugins/", http.FileServer(http.Dir("plugins"))).ServeHTTP(w, r)
		return
	}
	distDir := findDistDir()
	cleanRel := strings.TrimPrefix(strings.TrimPrefix(filepath.Clean(path), string(filepath.Separator)), "/")
	if cleanRel != "" {
		if serveDistFile(w, r, filepath.Join(distDir, cleanRel)) || s.serveFromFS(w, r, cleanRel) {
			return
		}
	}
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".js" || ext == ".mjs" || ext == ".css" || ext == ".woff2" || ext == ".wasm" || ext == ".map" {
		http.NotFound(w, r)
		return
	}
	if strings.HasPrefix(path, "/stash-proxy/") || strings.HasPrefix(path, "/avatars/") || strings.HasPrefix(path, "/assets/") || strings.HasPrefix(path, "/media/") || strings.HasPrefix(path, "/media-local/") || strings.HasPrefix(path, "/api/jobs/") {
		if s.unifiedHandler != nil { s.unifiedHandler.ServeHTTP(w, r); return }
	}
	if serveDistFile(w, r, filepath.Join(distDir, "index.html")) || s.serveFromFS(w, r, "index.html") { return }
	http.NotFound(w, r)
}
