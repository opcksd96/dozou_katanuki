package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUnifiedHandler_AvatarAndAssetServing(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "avatar_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// テスト用画像ファイルを作成
	testFile := filepath.Join(tempDir, "test_avatar.png")
	if err := os.WriteFile(testFile, []byte("fake-png-content"), 0644); err != nil {
		t.Fatalf("Failed to create test avatar: %v", err)
	}

	handler := NewUnifiedHandler(tempDir, nil)

	// 1. /avatars/ 経由で実体ファイルが存在する場合
	req := httptest.NewRequest("GET", "/avatars/test_avatar.png", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}
	if rec.Body.String() != "fake-png-content" {
		t.Errorf("Expected content 'fake-png-content', got '%s'", rec.Body.String())
	}

	// 2. /assets/ 経由で実体ファイルが存在する場合
	req = httptest.NewRequest("GET", "/assets/test_avatar.png", nil)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	// 3. ファイルが存在しない場合はデフォルトアバターSVGを返却
	req = httptest.NewRequest("GET", "/avatars/non_existent.png", nil)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200 for fallback avatar SVG, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "<svg") {
		t.Errorf("Expected SVG content, got '%s'", rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); ct != "image/svg+xml" {
		t.Errorf("Expected Content-Type image/svg+xml, got %s", ct)
	}
}
