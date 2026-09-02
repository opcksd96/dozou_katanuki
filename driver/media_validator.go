// driver/media_validator.go (100行以下 - SPEC-PRINCIPLE-001)
package driver

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// ValidateMediaFile はダウンロードされたファイルの実体を検証し、HTMLやテキスト等の偽ファイルを排除します
func ValidateMediaFile(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("ファイルオープン失敗: %w", err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return fmt.Errorf("ファイル情報取得失敗: %w", err)
	}
	if fi.Size() == 0 {
		return fmt.Errorf("空ファイル (0 bytes)")
	}

	buf := make([]byte, 512)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return fmt.Errorf("ファイル読み込み失敗: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("読み込みデータなし (0 bytes)")
	}

	sample := buf[:n]
	mimeType := http.DetectContentType(sample)
	lowerStr := strings.ToLower(string(sample))

	// HTML / XML / JSONエラー等のテキスト偽ファイルを検知
	if strings.HasPrefix(mimeType, "text/") ||
		strings.Contains(lowerStr, "<!doctype") ||
		strings.Contains(lowerStr, "<html") ||
		strings.Contains(lowerStr, "<head") ||
		strings.Contains(lowerStr, "<body") ||
		strings.Contains(lowerStr, "<?xml") ||
		strings.Contains(lowerStr, "{\"errors\"") ||
		strings.Contains(lowerStr, "{\"error\"") {
		return fmt.Errorf("偽ファイル検知 (HTML/テキストデータ: %s)", mimeType)
	}

	if fi.Size() < 256 && !strings.HasPrefix(mimeType, "image/") {
		return fmt.Errorf("破損ファイル (サイズ不足: %d bytes)", fi.Size())
	}

	return nil
}
