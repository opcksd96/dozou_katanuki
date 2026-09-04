package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	"dozou_katanuki/models"
)

func TestHandleRoot(t *testing.T) {
	service := NewBroadcastService(
		models.NetworkConfig{MiddlewarePort: 5175, PublicBindAddress: "0.0.0.0"},
		models.BroadcastConfig{Enabled: true},
		NewUnifiedHandler("./assets", nil), nil, nil, nil,
	)

	// Mock FS
	mockFS := fstest.MapFS{
		"frontend/dist/index.html": {
			Data: []byte("<!DOCTYPE html><html><body><div id=\"app\"></div></body></html>"),
		},
		"frontend/dist/assets/index.js": {
			Data: []byte("console.log('test');"),
		},
		"frontend/dist/assets/style.css": {
			Data: []byte("body { background: #000; }"),
		},
	}
	service.SetDistFS(mockFS)

	// 1. Root "/" should return HTML
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	service.handleRoot(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for /, got %d", rec.Code)
	}
	if !strings.Contains(rec.Header().Get("Content-Type"), "text/html") {
		t.Fatalf("Expected text/html, got %s", rec.Header().Get("Content-Type"))
	}

	// 2. JS bundle "/assets/index.js"
	reqJS := httptest.NewRequest("GET", "/assets/index.js", nil)
	recJS := httptest.NewRecorder()
	service.handleRoot(recJS, reqJS)
	if recJS.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for JS, got %d", recJS.Code)
	}
	if !strings.Contains(recJS.Header().Get("Content-Type"), "application/javascript") {
		t.Errorf("Expected application/javascript, got %s", recJS.Header().Get("Content-Type"))
	}

	// 3. CSS "/assets/style.css"
	reqCSS := httptest.NewRequest("GET", "/assets/style.css", nil)
	recCSS := httptest.NewRecorder()
	service.handleRoot(recCSS, reqCSS)
	if recCSS.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for CSS, got %d", recCSS.Code)
	}
	if !strings.Contains(recCSS.Header().Get("Content-Type"), "text/css") {
		t.Errorf("Expected text/css, got %s", recCSS.Header().Get("Content-Type"))
	}

	// 4. Missing JS should return 404, NEVER avatar SVG
	reqMissJS := httptest.NewRequest("GET", "/assets/missing-bundle.js", nil)
	recMissJS := httptest.NewRecorder()
	service.handleRoot(recMissJS, reqMissJS)
	if recMissJS.Code != http.StatusNotFound {
		t.Fatalf("Expected 404 for missing JS, got %d (Content-Type: %s)", recMissJS.Code, recMissJS.Header().Get("Content-Type"))
	}
}
