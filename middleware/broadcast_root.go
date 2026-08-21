// middleware/broadcast_root.go (100行以下)
package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

func (s *BroadcastService) handleRoot(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.HasPrefix(path, "/stash-proxy/") || strings.HasPrefix(path, "/avatars/") ||
		strings.HasPrefix(path, "/assets/") || strings.HasPrefix(path, "/api/jobs/") {
		if s.unifiedHandler != nil {
			s.unifiedHandler.ServeHTTP(w, r)
			return
		}
	}
	if path == "/" || path == "/cast" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		st := s.GetStatus()
		html := fmt.Sprintf(`<!DOCTYPE html><html lang="ja"><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width, initial-scale=1.0"><title>Dozou Katanuki Broadcast</title><style>body{background:#020617;color:#e2e8f0;font-family:sans-serif;display:flex;justify-content:center;align-items:center;min-height:100vh;margin:0;padding:20px;box-sizing:border-box}.card{background:#0f172a;border:1px solid #1e293b;border-radius:16px;padding:28px;max-width:540px;width:100%%;box-shadow:0 10px 25px rgba(0,0,0,0.5)}h1{font-size:20px;margin-top:0;color:#38bdf8;display:flex;align-items:center;gap:8px}.badge{display:inline-block;padding:4px 10px;background:#064e3b;color:#34d399;border-radius:9999px;font-size:12px;font-weight:bold}.item{margin:14px 0;font-size:14px}.label{color:#94a3b8;font-size:12px;margin-bottom:4px}.val{background:#020617;padding:8px 12px;border-radius:8px;border:1px solid #334155;font-family:monospace;word-break:break-all}a{color:#38bdf8;text-decoration:none}a:hover{text-decoration:underline}</style></head><body><div class="card"><h1>📡 Dozou Katanuki Broadcast</h1><p><span class="badge">● Streaming Online</span></p><div class="item"><div class="label">Cast Stream URL</div><div class="val">%s</div></div><div class="item"><div class="label">Allowed Subnets</div><div class="val">%s</div></div><div class="item"><div class="label">Available Endpoints</div><div class="val">• <a href="/api/timeline">/api/timeline</a><br>• <a href="/api/accounts">/api/accounts</a><br>• <a href="/api/broadcast/status">/api/broadcast/status</a><br>• /stash-proxy/* (Stash Video Streaming)<br>• /avatars/* (Author Icons)</div></div></div></body></html>`, st.CastURL, strings.Join(st.AllowedNetworks, ", "))
		_, _ = w.Write([]byte(html))
		return
	}
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
