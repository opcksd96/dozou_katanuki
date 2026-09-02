// driver/media_validator_test.go (100行以下 - SPEC-PRINCIPLE-001)
package driver

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateMediaFile(t *testing.T) {
	tempDir := t.TempDir()

	// 1. HTML偽ファイル (text/html)
	htmlPath := filepath.Join(tempDir, "fake.jpg")
	htmlContent := []byte("<!DOCTYPE html><html><head><title>404 Not Found</title></head><body>Error</body></html>")
	if err := os.WriteFile(htmlPath, htmlContent, 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	err := ValidateMediaFile(htmlPath)
	if err == nil {
		t.Fatalf("expected error for HTML fake file, but got nil")
	}
	t.Logf("HTML fake file correctly rejected: %v", err)

	// 2. 空ファイル
	emptyPath := filepath.Join(tempDir, "empty.jpg")
	_ = os.WriteFile(emptyPath, []byte{}, 0644)
	if err := ValidateMediaFile(emptyPath); err == nil {
		t.Fatalf("expected error for empty file, but got nil")
	}

	// 3. 有効なJPEG画像ヘッダー
	validJpegPath := filepath.Join(tempDir, "valid.jpg")
	validJpegContent := append([]byte("\xFF\xD8\xFF\xE0\x00\x10JFIF\x00\x01\x01\x01\x00`\x00`\x00\x00"), make([]byte, 500)...)
	_ = os.WriteFile(validJpegPath, validJpegContent, 0644)
	if err := ValidateMediaFile(validJpegPath); err != nil {
		t.Fatalf("expected valid JPEG to pass, but got: %v", err)
	}
	t.Logf("Valid JPEG correctly accepted")
}
