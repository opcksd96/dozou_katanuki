package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"dozou_katanuki/models"
)

// BroadcastService は家庭内LANデバイス向けにメディアやタイムラインをHTTP配信する中継サーバーです (SPEC-SCHEDULER-001)
type BroadcastService struct {
	netCfg          models.NetworkConfig
	bcastCfg        models.BroadcastConfig
	unifiedHandler  *UnifiedHandler
	timelineService *TimelineService
	emitter         EventEmitter

	server   *http.Server
	listener net.Listener
	mu       sync.RWMutex
	running  bool
}

// NewBroadcastService は新しい BroadcastService を生成します
func NewBroadcastService(
	netCfg models.NetworkConfig,
	bcastCfg models.BroadcastConfig,
	handler *UnifiedHandler,
	timeline *TimelineService,
	emitter EventEmitter,
) *BroadcastService {
	if netCfg.MiddlewarePort <= 0 {
		netCfg.MiddlewarePort = 5175
	}
	if netCfg.PublicBindAddress == "" {
		netCfg.PublicBindAddress = "0.0.0.0"
	}
	if len(bcastCfg.AllowedNetworks) == 0 {
		bcastCfg.AllowedNetworks = []string{"192.168.0.0/16", "10.0.0.0/8", "172.16.0.0/12", "127.0.0.1/32", "::1/128"}
	}

	return &BroadcastService{
		netCfg:          netCfg,
		bcastCfg:        bcastCfg,
		unifiedHandler:  handler,
		timelineService: timeline,
		emitter:         emitter,
	}
}

// Start はサーバーを開始します（enabled が true の場合）
func (s *BroadcastService) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return nil
	}

	if !s.bcastCfg.Enabled {
		log.Println("[Broadcast] LAN Broadcast is currently disabled in config.")
		return nil
	}

	return s.startServerLocked()
}

// Stop はサーバーを停止します
func (s *BroadcastService) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running || s.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.server.Shutdown(ctx)
	s.running = false
	s.server = nil
	s.listener = nil
	log.Println("[Broadcast] LAN Broadcast Server stopped.")
	return err
}

// UpdateConfig は設定を更新し、必要に応じてサーバーを再起動します
func (s *BroadcastService) UpdateConfig(netCfg models.NetworkConfig, bcastCfg models.BroadcastConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.netCfg = netCfg
	s.bcastCfg = bcastCfg

	if s.running {
		if s.server != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			_ = s.server.Shutdown(ctx)
			cancel()
			s.server = nil
			s.listener = nil
			s.running = false
		}
	}

	if s.bcastCfg.Enabled {
		return s.startServerLocked()
	}
	return nil
}

func (s *BroadcastService) startServerLocked() error {
	addr := fmt.Sprintf("%s:%d", s.netCfg.PublicBindAddress, s.netCfg.MiddlewarePort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("[Broadcast] Failed to bind address %s: %v", addr, err)
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/timeline", s.handleTimelineAPI)
	mux.HandleFunc("/api/article", s.handleArticleAPI)
	mux.HandleFunc("/api/accounts", s.handleAccountsAPI)
	mux.HandleFunc("/api/broadcast/status", s.handleStatusAPI)
	mux.HandleFunc("/", s.handleRoot)

	// セキュリティゲートウェイ ＆ CORS ミドルウェアでラップ
	handler := s.securityMiddleware(s.corsMiddleware(mux))

	server := &http.Server{
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	s.server = server
	s.listener = listener
	s.running = true

	localIPs := GetLocalIPv4s()
	log.Printf("[Broadcast] 📡 LAN Broadcast Server running on http://%s (Port: %d)", addr, s.netCfg.MiddlewarePort)
	for _, ip := range localIPs {
		log.Printf("[Broadcast]  ➔ Cast URL for local devices: http://%s:%d", ip, s.netCfg.MiddlewarePort)
	}

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Printf("[Broadcast] Server encountered error: %v", err)
			s.mu.Lock()
			s.running = false
			s.mu.Unlock()
		}
	}()

	return nil
}

// GetStatus は現在の配信ステータスを取得します
func (s *BroadcastService) GetStatus() *models.BroadcastStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	localIPs := GetLocalIPv4s()
	castURL := ""
	if len(localIPs) > 0 && s.netCfg.MiddlewarePort > 0 {
		castURL = fmt.Sprintf("http://%s:%d", localIPs[0], s.netCfg.MiddlewarePort)
	}

	return &models.BroadcastStatus{
		Enabled:         s.bcastCfg.Enabled,
		Running:         s.running,
		BindAddress:     s.netCfg.PublicBindAddress,
		Port:            s.netCfg.MiddlewarePort,
		LocalIPs:        localIPs,
		AllowedNetworks: s.bcastCfg.AllowedNetworks,
		CastURL:         castURL,
	}
}

// securityMiddleware は IP / CIDR サブネット検証を行うセキュリティゲートウェイです
func (s *BroadcastService) securityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := s.extractClientIP(r)

		if !s.isIPAllowed(clientIP) {
			log.Printf("[Broadcast Security] 403 Forbidden: Blocked unauthorized client IP %s for %s", clientIP, r.URL.Path)
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error":     "Forbidden",
				"message":   fmt.Sprintf("Access denied by broadcast security policy for IP: %s", clientIP),
				"client_ip": clientIP,
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

// corsMiddleware は CORS 制限を中和するヘッダーを自動付与します
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

// extractClientIP はリクエストからクライアントのIPアドレスを抽出します
func (s *BroadcastService) extractClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			ip := strings.TrimSpace(parts[0])
			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// isIPAllowed は指定されたIPが許可ネットワーク（CIDR/IP）に含まれるかを検証します
func (s *BroadcastService) isIPAllowed(ipStr string) bool {
	parsedIP := net.ParseIP(ipStr)
	if parsedIP == nil {
		return false
	}

	// ループバックは常に安全に通す
	if parsedIP.IsLoopback() {
		return true
	}

	for _, cidr := range s.bcastCfg.AllowedNetworks {
		cidr = strings.TrimSpace(cidr)
		if cidr == "" {
			continue
		}

		// CIDR 形式 (例: 192.168.1.0/24)
		if strings.Contains(cidr, "/") {
			_, ipNet, err := net.ParseCIDR(cidr)
			if err == nil && ipNet.Contains(parsedIP) {
				return true
			}
		} else {
			// 単一 IP 形式 (例: 192.168.1.50)
			targetIP := net.ParseIP(cidr)
			if targetIP != nil && targetIP.Equal(parsedIP) {
				return true
			}
		}
	}

	return false
}

func (s *BroadcastService) handleRoot(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// メディア・アバター・Stashプロキシ・JobAPI の中継
	if strings.HasPrefix(path, "/stash-proxy/") ||
		strings.HasPrefix(path, "/avatars/") ||
		strings.HasPrefix(path, "/assets/") ||
		strings.HasPrefix(path, "/api/jobs/") {
		if s.unifiedHandler != nil {
			s.unifiedHandler.ServeHTTP(w, r)
			return
		}
	}

	// ルート案内
	if path == "/" || path == "/cast" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		status := s.GetStatus()
		html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="ja">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Dozou Katanuki LAN Broadcast</title>
	<style>
		body { background: #020617; color: #e2e8f0; font-family: sans-serif; display: flex; justify-content: center; align-items: center; min-height: 100vh; margin: 0; padding: 20px; box-sizing: border-box; }
		.card { background: #0f172a; border: 1px solid #1e293b; border-radius: 16px; padding: 28px; max-width: 540px; width: 100%%; box-shadow: 0 10px 25px rgba(0,0,0,0.5); }
		h1 { font-size: 20px; margin-top: 0; color: #38bdf8; display: flex; align-items: center; gap: 8px; }
		.badge { display: inline-block; padding: 4px 10px; background: #064e3b; color: #34d399; border-radius: 9999px; font-size: 12px; font-weight: bold; }
		.item { margin: 14px 0; font-size: 14px; }
		.label { color: #94a3b8; font-size: 12px; margin-bottom: 4px; }
		.val { background: #020617; padding: 8px 12px; border-radius: 8px; border: 1px solid #334155; font-family: monospace; word-break: break-all; }
		a { color: #38bdf8; text-decoration: none; }
		a:hover { text-decoration: underline; }
	</style>
</head>
<body>
	<div class="card">
		<h1>📡 Dozou Katanuki Broadcast</h1>
		<p><span class="badge">● Streaming Online</span></p>
		<div class="item">
			<div class="label">Cast Stream URL</div>
			<div class="val">%s</div>
		</div>
		<div class="item">
			<div class="label">Allowed Subnets</div>
			<div class="val">%s</div>
		</div>
		<div class="item">
			<div class="label">Available Endpoints</div>
			<div class="val">
				• <a href="/api/timeline">/api/timeline</a><br>
				• <a href="/api/accounts">/api/accounts</a><br>
				• <a href="/api/broadcast/status">/api/broadcast/status</a><br>
				• /stash-proxy/* (Stash Video Streaming)<br>
				• /avatars/* (Author Icons)
			</div>
		</div>
	</div>
</body>
</html>`, status.CastURL, strings.Join(status.AllowedNetworks, ", "))
		_, _ = w.Write([]byte(html))
		return
	}

	http.NotFound(w, r)
}

func (s *BroadcastService) handleStatusAPI(w http.ResponseWriter, r *http.Request) {
	status := s.GetStatus()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(status)
}

func (s *BroadcastService) handleTimelineAPI(w http.ResponseWriter, r *http.Request) {
	if s.timelineService == nil {
		http.Error(w, `{"error":"Timeline service not available"}`, http.StatusServiceUnavailable)
		return
	}

	q := r.URL.Query()
	platform := q.Get("platform")
	if platform == "" {
		platform = "twitter"
	}
	accountID := q.Get("account_id")
	filter := q.Get("filter")
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))

	trees, err := s.timelineService.FetchTimeline(platform, accountID, filter, limit, offset)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(trees)
}

func (s *BroadcastService) handleArticleAPI(w http.ResponseWriter, r *http.Request) {
	if s.timelineService == nil {
		http.Error(w, `{"error":"Timeline service not available"}`, http.StatusServiceUnavailable)
		return
	}

	q := r.URL.Query()
	platform := q.Get("platform")
	if platform == "" {
		platform = "twitter"
	}
	id := q.Get("id")
	if id == "" {
		http.Error(w, `{"error":"Article id query parameter is required"}`, http.StatusBadRequest)
		return
	}

	detail, err := s.timelineService.GetArticleDetail(platform, id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(detail)
}

func (s *BroadcastService) handleAccountsAPI(w http.ResponseWriter, r *http.Request) {
	if s.timelineService == nil {
		http.Error(w, `{"error":"Timeline service not available"}`, http.StatusServiceUnavailable)
		return
	}

	platform := r.URL.Query().Get("platform")
	if platform == "" {
		platform = "twitter"
	}

	accounts, err := s.timelineService.GetAccounts(platform)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(accounts)
}

// GetLocalIPv4s はホストマシンの非ループバック IPv4 アドレス一覧を取得します
func GetLocalIPv4s() []string {
	var ips []string
	ifaces, err := net.Interfaces()
	if err != nil {
		return []string{"127.0.0.1"}
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue
			}
			ips = append(ips, ip.String())
		}
	}

	if len(ips) == 0 {
		ips = append(ips, "127.0.0.1")
	}
	return ips
}
