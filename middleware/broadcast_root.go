// middleware/broadcast_root.go (100行以下)
package middleware

import (
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	_ = mime.AddExtensionType(".js", "application/javascript; charset=utf-8")
	_ = mime.AddExtensionType(".mjs", "application/javascript; charset=utf-8")
	_ = mime.AddExtensionType(".css", "text/css; charset=utf-8")
	_ = mime.AddExtensionType(".html", "text/html; charset=utf-8")
	_ = mime.AddExtensionType(".json", "application/json; charset=utf-8")
	_ = mime.AddExtensionType(".wasm", "application/wasm")
	_ = mime.AddExtensionType(".woff2", "font/woff2")
	_ = mime.AddExtensionType(".svg", "image/svg+xml")
}

func findDistDir() string {
	candidates := []string{
		filepath.Join("frontend", "dist"),
		filepath.Join("..", "frontend", "dist"),
		filepath.Join("..", "..", "frontend", "dist"),
	}
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), "frontend", "dist"))
	}
	for _, c := range candidates {
		if fi, err := os.Stat(c); err == nil && fi.IsDir() {
			if _, err2 := os.Stat(filepath.Join(c, "index.html")); err2 == nil {
				return c
			}
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
	f, err := os.Open(filePath)
	if err != nil { return false }
	defer f.Close()
	fi, err := f.Stat()
	if err != nil || fi.IsDir() { return false }

	w.Header().Set("Content-Type", getContentType(filePath))
	http.ServeContent(w, r, fi.Name(), fi.ModTime(), f)
	return true
}

func (s *BroadcastService) serveFromFS(w http.ResponseWriter, r *http.Request, cleanPath string) bool {
	s.mu.RLock()
	dfs := s.distFS
	s.mu.RUnlock()
	if dfs == nil { return false }

	prefixes := []string{"", "frontend/dist/", "dist/", "frontend/"}
	for _, pfx := range prefixes {
		lookup := filepath.ToSlash(filepath.Clean(pfx + cleanPath))
		lookup = strings.TrimPrefix(lookup, "/")
		f, err := dfs.Open(lookup)
		if err == nil {
			defer f.Close()
			fi, err := f.Stat()
			if err == nil && !fi.IsDir() {
				w.Header().Set("Content-Type", getContentType(lookup))
				if rs, ok := f.(io.ReadSeeker); ok {
					http.ServeContent(w, r, fi.Name(), fi.ModTime(), rs)
				} else {
					_, _ = io.Copy(w, f)
				}
				return true
			}
		}
	}
	return false
}

func (s *BroadcastService) handleRoot(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	cleanRel := strings.TrimPrefix(filepath.Clean(path), string(filepath.Separator))
	cleanRel = strings.TrimPrefix(cleanRel, "/")
	distDir := findDistDir()

	// 1. 静的ファイル（JS, CSS, Font, HTML等）配信
	if cleanRel != "" {
		diskFile := filepath.Join(distDir, cleanRel)
		if serveDistFile(w, r, diskFile) { return }
		if s.serveFromFS(w, r, cleanRel) { return }
	}

	// 2. フロントエンド専用アセット拡張子の場合は 404 (アバター誤爆防止)
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".js" || ext == ".mjs" || ext == ".css" || ext == ".woff2" || ext == ".wasm" || ext == ".map" {
		http.NotFound(w, r); return
	}

	// 3. プラグイン
	if strings.HasPrefix(path, "/plugins/") {
		http.StripPrefix("/plugins/", http.FileServer(http.Dir("plugins"))).ServeHTTP(w, r); return
	}

	// 4. Stash プロキシ / アバター
	if strings.HasPrefix(path, "/stash-proxy/") || strings.HasPrefix(path, "/avatars/") ||
		strings.HasPrefix(path, "/assets/") || strings.HasPrefix(path, "/api/jobs/") {
		if s.unifiedHandler != nil { s.unifiedHandler.ServeHTTP(w, r); return }
	}

	// 5. SPA ルートフォールバック (index.html)
	indexDisk := filepath.Join(distDir, "index.html")
	if serveDistFile(w, r, indexDisk) { return }
	if s.serveFromFS(w, r, "index.html") { return }

	http.NotFound(w, r)
}

func GetLocalIPv4s() []string {
	var ips []string
	ifaces, err := net.Interfaces()
	if err != nil { return []string{"127.0.0.1"} }
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 { continue }
		addrs, err := iface.Addrs()
		if err != nil { continue }
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet: ip = v.IP
			case *net.IPAddr: ip = v.IP
			}
			if ip == nil || ip.IsLoopback() { continue }
			if ip4 := ip.To4(); ip4 != nil { ips = append(ips, ip4.String()) }
		}
	}
	if len(ips) == 0 { ips = append(ips, "127.0.0.1") }
	return ips
}

func GetLocalSubnets() []string {
	seen := make(map[string]bool)
	var subnets []string
	ifaces, err := net.Interfaces()
	if err != nil { return []string{"127.0.0.1/32"} }
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 { continue }
		addrs, err := iface.Addrs()
		if err != nil { continue }
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				if ipNet.IP == nil || ipNet.IP.IsLoopback() { continue }
				if ip4 := ipNet.IP.To4(); ip4 != nil {
					maskedIP := ip4.Mask(ipNet.Mask)
					ones, _ := ipNet.Mask.Size()
					cidr := fmt.Sprintf("%s/%d", maskedIP.String(), ones)
					if !seen[cidr] { seen[cidr] = true; subnets = append(subnets, cidr) }
				}
			}
		}
	}
	if !seen["127.0.0.1/32"] { subnets = append(subnets, "127.0.0.1/32") }
	return subnets
}
