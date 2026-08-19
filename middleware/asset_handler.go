// middleware/asset_handler.go (100行以下)
package middleware

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// NewAssetFileHandler は /avatars/ 宛てのリクエストを安全に実体サーブします
func NewAssetFileHandler(assetsBaseDir string) http.Handler {
	absBase, err := filepath.Abs(assetsBaseDir)
	if err != nil {
		absBase = assetsBaseDir
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if strings.HasPrefix(path, "/avatars/") {
			relPath := strings.TrimPrefix(path, "/avatars/")
			targetPath := filepath.Join(absBase, filepath.Clean(relPath))

			info, err := os.Stat(targetPath)
			if err == nil && !info.IsDir() {
				http.ServeFile(w, r, targetPath)
				return
			}
			http.NotFound(w, r)
			return
		}

		http.NotFound(w, r)
	})
}
