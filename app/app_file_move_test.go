// app/app_file_move_test.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMoveFileSafe(t *testing.T) {
	tempDir := t.TempDir()

	srcFile := filepath.Join(tempDir, "source.txt")
	dstFile := filepath.Join(tempDir, "dest.txt")

	content := []byte("hello thunder move")
	if err := os.WriteFile(srcFile, content, 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// 1. 通常の移動テスト
	if err := moveFileSafe(srcFile, dstFile); err != nil {
		t.Fatalf("moveFileSafe failed: %v", err)
	}

	// 元ファイルが削除されていることを確認
	if _, err := os.Stat(srcFile); !os.IsNotExist(err) {
		t.Errorf("Source file still exists after move")
	}

	// 先ファイルの内容確認
	got, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("Failed to read dest file: %v", err)
	}
	if string(got) != string(content) {
		t.Errorf("Content mismatch: got %s, want %s", string(got), string(content))
	}

	// 2. 同一パスの場合
	if err := moveFileSafe(dstFile, dstFile); err != nil {
		t.Errorf("moveFileSafe same path should return nil, got: %v", err)
	}

	// 3. 上書き移動テスト
	newSrc := filepath.Join(tempDir, "new_source.txt")
	newContent := []byte("overwritten content")
	_ = os.WriteFile(newSrc, newContent, 0644)

	if err := moveFileSafe(newSrc, dstFile); err != nil {
		t.Fatalf("moveFileSafe overwrite failed: %v", err)
	}
	if _, err := os.Stat(newSrc); !os.IsNotExist(err) {
		t.Errorf("New source file still exists after overwrite move")
	}
	gotOverwritten, _ := os.ReadFile(dstFile)
	if string(gotOverwritten) != string(newContent) {
		t.Errorf("Overwritten content mismatch: got %s, want %s", string(gotOverwritten), string(newContent))
	}
}
