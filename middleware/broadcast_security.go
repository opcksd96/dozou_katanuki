// middleware/broadcast_security.go (100行以下)
package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

func (s *BroadcastService) securityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := s.extractClientIP(r)
		if !s.isIPAllowed(clientIP) {
			log.Printf("[Broadcast Security] 403 Forbidden: Blocked %s for %s", clientIP, r.URL.Path)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "Forbidden", "message": fmt.Sprintf("Access denied for IP: %s", clientIP), "client_ip": clientIP,
			})
			return
		}
		log.Printf("[Broadcast Access] 200 OK: %s %s from %s", r.Method, r.URL.Path, clientIP)
		next.ServeHTTP(w, r)
	})
}

func (s *BroadcastService) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, HEAD")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, Range, X-Requested-With")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Range, Accept-Ranges")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *BroadcastService) extractClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			ip := strings.TrimSpace(parts[0])
			if net.ParseIP(ip) != nil { return ip }
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil { return r.RemoteAddr }
	return host
}

func (s *BroadcastService) isIPAllowed(ipStr string) bool {
	parsedIP := net.ParseIP(ipStr)
	if parsedIP == nil { return false }
	if ip4 := parsedIP.To4(); ip4 != nil {
		parsedIP = ip4
	}
	if parsedIP.IsLoopback() || parsedIP.IsPrivate() {
		return true
	}

	for _, cidr := range s.bcastCfg.AllowedNetworks {
		cidr = strings.TrimSpace(cidr)
		if cidr == "" { continue }
		if strings.Contains(cidr, "/") {
			_, ipNet, err := net.ParseCIDR(cidr)
			if err == nil && ipNet.Contains(parsedIP) { return true }
		} else {
			targetIP := net.ParseIP(cidr)
			if targetIP != nil {
				if t4 := targetIP.To4(); t4 != nil { targetIP = t4 }
				if targetIP.Equal(parsedIP) { return true }
			}
		}
	}
	return false
}
