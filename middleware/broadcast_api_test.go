// middleware/broadcast_api_test.go (100行以下)
package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dozou_katanuki/models"
)

func TestBroadcastCORSMiddleware(t *testing.T) {
	service := NewBroadcastService(models.NetworkConfig{MiddlewarePort: 5175}, models.BroadcastConfig{Enabled: true}, nil, nil, nil)
	handler := service.corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK); _, _ = w.Write([]byte("DATA"))
	}))

	reqOpt := httptest.NewRequest("OPTIONS", "/api/timeline", nil)
	recOpt := httptest.NewRecorder()
	handler.ServeHTTP(recOpt, reqOpt)
	if recOpt.Code != http.StatusOK || recOpt.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("Expected CORS headers on OPTIONS")
	}

	reqGet := httptest.NewRequest("GET", "/api/timeline", nil)
	recGet := httptest.NewRecorder()
	handler.ServeHTTP(recGet, reqGet)
	if recGet.Code != http.StatusOK || recGet.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("Expected CORS headers on GET")
	}
}

func TestBroadcastStatusAPI(t *testing.T) {
	service := NewBroadcastService(
		models.NetworkConfig{MiddlewarePort: 5175, PublicBindAddress: "0.0.0.0"},
		models.BroadcastConfig{Enabled: true, AllowedNetworks: []string{"192.168.0.0/16"}},
		nil, nil, nil,
	)
	req := httptest.NewRequest("GET", "/api/broadcast/status", nil)
	rec := httptest.NewRecorder()
	service.handleStatusAPI(rec, req)
	if rec.Code != http.StatusOK { t.Fatalf("Expected 200 OK, got %d", rec.Code) }

	var status models.BroadcastStatus
	if err := json.Unmarshal(rec.Body.Bytes(), &status); err != nil || !status.Enabled || status.Port != 5175 {
		t.Fatalf("Unexpected status: %+v", status)
	}
}

func TestGenerateSelfSignedCert(t *testing.T) {
	cert, err := GenerateSelfSignedCert([]string{"192.168.10.1", "192.168.3.202", "127.0.0.1"})
	if err != nil {
		t.Fatalf("GenerateSelfSignedCert failed: %v", err)
	}
	if len(cert.Certificate) == 0 {
		t.Fatalf("Expected valid certificate chain")
	}
}

func TestBroadcastLifecycle(t *testing.T) {
	service := NewBroadcastService(
		models.NetworkConfig{MiddlewarePort: 0, PublicBindAddress: "127.0.0.1"},
		models.BroadcastConfig{Enabled: true},
		nil, nil, nil,
	)
	if err := service.Start(t.Context()); err != nil {
		t.Fatalf("Failed to start broadcast service: %v", err)
	}
	if !service.running {
		t.Fatalf("Expected service to be running")
	}

	// Stop must return immediately and clean up resources without hanging
	if err := service.Stop(); err != nil {
		t.Fatalf("Failed to stop broadcast service: %v", err)
	}
	if service.running {
		t.Fatalf("Expected service to be stopped")
	}
}

