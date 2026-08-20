package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dozou_katanuki/models"
)

func TestBroadcastSecurityIPAllowed(t *testing.T) {
	service := NewBroadcastService(
		models.NetworkConfig{
			MiddlewarePort:    5175,
			PublicBindAddress: "0.0.0.0",
		},
		models.BroadcastConfig{
			Enabled:         true,
			AllowedNetworks: []string{"192.168.1.0/24", "10.0.0.0/8", "172.16.50.10"},
		},
		nil,
		nil,
		nil,
	)

	// 許可されるべきIP
	allowedIPs := []string{
		"127.0.0.1",       // ループバック
		"::1",             // IPv6 ループバック
		"192.168.1.1",     // 192.168.1.0/24 内
		"192.168.1.254",   // 192.168.1.0/24 内
		"10.0.0.1",       // 10.0.0.0/8 内
		"10.254.254.254", // 10.0.0.0/8 内
		"172.16.50.10",    // 単一IP一致
	}

	for _, ip := range allowedIPs {
		if !service.isIPAllowed(ip) {
			t.Errorf("Expected IP %s to be allowed, but was rejected", ip)
		}
	}

	// 拒否されるべきIP
	blockedIPs := []string{
		"192.168.2.1",     // 192.168.1.0/24 外
		"172.16.50.11",    // 単一IP不一致
		"203.0.113.5",     // パブリックIP
		"8.8.8.8",         // 外部IP
		"invalid-ip",      // 不正なIP形式
	}

	for _, ip := range blockedIPs {
		if service.isIPAllowed(ip) {
			t.Errorf("Expected IP %s to be blocked, but was allowed", ip)
		}
	}
}

func TestBroadcastSecurityMiddleware(t *testing.T) {
	service := NewBroadcastService(
		models.NetworkConfig{MiddlewarePort: 5175},
		models.BroadcastConfig{
			Enabled:         true,
			AllowedNetworks: []string{"192.168.1.0/24"},
		},
		nil,
		nil,
		nil,
	)

	handler := service.securityMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}))

	// Case 1: 許可IP (192.168.1.50)
	reqOK := httptest.NewRequest("GET", "/api/timeline", nil)
	reqOK.RemoteAddr = "192.168.1.50:12345"
	recOK := httptest.NewRecorder()
	handler.ServeHTTP(recOK, reqOK)

	if recOK.Code != http.StatusOK {
		t.Errorf("Expected status 200 for allowed IP, got %d", recOK.Code)
	}

	// Case 2: 拒否IP (203.0.113.50)
	reqBlocked := httptest.NewRequest("GET", "/api/timeline", nil)
	reqBlocked.RemoteAddr = "203.0.113.50:12345"
	recBlocked := httptest.NewRecorder()
	handler.ServeHTTP(recBlocked, reqBlocked)

	if recBlocked.Code != http.StatusForbidden {
		t.Errorf("Expected status 403 Forbidden for blocked IP, got %d", recBlocked.Code)
	}
}

func TestBroadcastCORSMiddleware(t *testing.T) {
	service := NewBroadcastService(
		models.NetworkConfig{MiddlewarePort: 5175},
		models.BroadcastConfig{Enabled: true},
		nil,
		nil,
		nil,
	)

	handler := service.corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("DATA"))
	}))

	// OPTIONS プリフライト
	reqOpt := httptest.NewRequest("OPTIONS", "/api/timeline", nil)
	recOpt := httptest.NewRecorder()
	handler.ServeHTTP(recOpt, reqOpt)

	if recOpt.Code != http.StatusOK {
		t.Errorf("Expected status 200 for OPTIONS, got %d", recOpt.Code)
	}
	if recOpt.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("Expected Access-Control-Allow-Origin: *, got %s", recOpt.Header().Get("Access-Control-Allow-Origin"))
	}

	// GET リクエスト
	reqGet := httptest.NewRequest("GET", "/api/timeline", nil)
	recGet := httptest.NewRecorder()
	handler.ServeHTTP(recGet, reqGet)

	if recGet.Code != http.StatusOK {
		t.Errorf("Expected status 200 for GET, got %d", recGet.Code)
	}
	if recGet.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("Expected Access-Control-Allow-Origin: *, got %s", recGet.Header().Get("Access-Control-Allow-Origin"))
	}
}

func TestBroadcastStatusAPI(t *testing.T) {
	service := NewBroadcastService(
		models.NetworkConfig{MiddlewarePort: 5175, PublicBindAddress: "0.0.0.0"},
		models.BroadcastConfig{
			Enabled:         true,
			AllowedNetworks: []string{"192.168.0.0/16"},
		},
		nil,
		nil,
		nil,
	)

	req := httptest.NewRequest("GET", "/api/broadcast/status", nil)
	rec := httptest.NewRecorder()
	service.handleStatusAPI(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", rec.Code)
	}

	var status models.BroadcastStatus
	if err := json.Unmarshal(rec.Body.Bytes(), &status); err != nil {
		t.Fatalf("Failed to parse status JSON: %v", err)
	}

	if !status.Enabled {
		t.Errorf("Expected status.Enabled to be true")
	}
	if status.Port != 5175 {
		t.Errorf("Expected port 5175, got %d", status.Port)
	}
	if len(status.AllowedNetworks) != 1 || status.AllowedNetworks[0] != "192.168.0.0/16" {
		t.Errorf("Unexpected AllowedNetworks: %v", status.AllowedNetworks)
	}
}
