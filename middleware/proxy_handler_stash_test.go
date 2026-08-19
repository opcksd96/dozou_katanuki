package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestUnifiedHandler_StashProxy(t *testing.T) {
	// モック Stash サーバーを起動
	stashServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/scene/123/stream" {
			w.Header().Set("Content-Type", "video/mp4")
			w.Header().Set("Accept-Ranges", "bytes")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("mock-video-stream-data"))
			return
		}
		if r.URL.Path == "/image/456/image" {
			w.Header().Set("Content-Type", "image/jpeg")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("mock-image-data"))
			return
		}
		http.NotFound(w, r)
	}))
	defer stashServer.Close()

	stashURL, err := url.Parse(stashServer.URL)
	if err != nil {
		t.Fatalf("Failed to parse stash URL: %v", err)
	}

	handler := NewUnifiedHandler("./non_existent_avatar_dir", stashURL)

	// 1. /stash-proxy/scene/123/stream が /scene/123/stream へ中継されること
	req := httptest.NewRequest("GET", "/stash-proxy/scene/123/stream", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	resp := rec.Result()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	if string(body) != "mock-video-stream-data" {
		t.Errorf("Expected 'mock-video-stream-data', got '%s'", string(body))
	}
	if cors := resp.Header.Get("Access-Control-Allow-Origin"); cors != "*" {
		t.Errorf("Expected CORS header '*', got '%s'", cors)
	}

	// 2. /stash-proxy/image/456/image が /image/456/image へ中継されること
	req = httptest.NewRequest("GET", "/stash-proxy/image/456/image", nil)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	resp = rec.Result()
	body, _ = io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	if string(body) != "mock-image-data" {
		t.Errorf("Expected 'mock-image-data', got '%s'", string(body))
	}
}

func TestUnifiedHandler_StashProxyOffline(t *testing.T) {
	// 停止しているStash URL
	stashURL, _ := url.Parse("http://127.0.0.1:59999")
	handler := NewUnifiedHandler("./temp", stashURL)

	req := httptest.NewRequest("GET", "/stash-proxy/scene/1/stream", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Errorf("Expected status 502 Bad Gateway for offline Stash, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "<svg") {
		t.Errorf("Expected SVG fallback for offline Stash, got '%s'", rec.Body.String())
	}
}
