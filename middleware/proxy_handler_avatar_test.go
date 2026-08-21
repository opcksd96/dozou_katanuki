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

	// 3. 個別アバターが存在しない場合は デフォルトSVG (200 OK) を返却して404を防止
	req = httptest.NewRequest("GET", "/avatars/twitter/non_existent.png", nil)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200 (default avatar SVG) for missing avatar, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "<svg") {
		t.Errorf("Expected SVG default avatar for missing avatar, got '%s'", rec.Body.String())
	}

	// 4. パストラバーサルの試行は 403 Forbidden または 404 で遮断
	req = httptest.NewRequest("GET", "/avatars/../../secret.txt", nil)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden && rec.Code != http.StatusNotFound {
		t.Errorf("Expected status 403 or 404 for path traversal attempt, got %d", rec.Code)
	}

	// 5. 明示的なデフォルトアバター要求時はデフォルトアバターSVGを返却
	req = httptest.NewRequest("GET", "/avatars/default_avatar.jpg", nil)
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200 for default avatar SVG, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "<svg") {
		t.Errorf("Expected SVG content, got '%s'", rec.Body.String())
	}
}
