// middleware/broadcast_security_test.go (100行以下)
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"dozou_katanuki/models"
)

func TestBroadcastSecurityIPAllowed(t *testing.T) {
	service := NewBroadcastService(
		models.NetworkConfig{MiddlewarePort: 5175, PublicBindAddress: "0.0.0.0"},
		models.BroadcastConfig{Enabled: true, AllowedNetworks: []string{"192.168.1.0/24", "10.0.0.0/8", "172.16.50.10"}},
		nil, nil, nil,
	)

	allowed := []string{"127.0.0.1", "::1", "192.168.1.1", "192.168.1.254", "10.0.0.1", "10.254.254.254", "172.16.50.10"}
	for _, ip := range allowed {
		if !service.isIPAllowed(ip) { t.Errorf("Expected IP %s to be allowed", ip) }
	}

	blocked := []string{"192.168.2.1", "172.16.50.11", "203.0.113.5", "8.8.8.8", "invalid-ip"}
	for _, ip := range blocked {
		if service.isIPAllowed(ip) { t.Errorf("Expected IP %s to be blocked", ip) }
	}
}

func TestBroadcastSecurityMiddleware(t *testing.T) {
	service := NewBroadcastService(
		models.NetworkConfig{MiddlewarePort: 5175},
		models.BroadcastConfig{Enabled: true, AllowedNetworks: []string{"192.168.1.0/24"}},
		nil, nil, nil,
	)
	handler := service.securityMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK); _, _ = w.Write([]byte("OK"))
	}))

	reqOK := httptest.NewRequest("GET", "/api/timeline", nil)
	reqOK.RemoteAddr = "192.168.1.50:12345"
	recOK := httptest.NewRecorder()
	handler.ServeHTTP(recOK, reqOK)
	if recOK.Code != http.StatusOK { t.Errorf("Expected 200, got %d", recOK.Code) }

	reqBlocked := httptest.NewRequest("GET", "/api/timeline", nil)
	reqBlocked.RemoteAddr = "203.0.113.50:12345"
	recBlocked := httptest.NewRecorder()
	handler.ServeHTTP(recBlocked, reqBlocked)
	if recBlocked.Code != http.StatusForbidden { t.Errorf("Expected 403, got %d", recBlocked.Code) }
}
