//go:build !windows
// +build !windows

package driver

import (
	"os"
	"path/filepath"
	"time"
)

// MoveToRecycleBin は非Windows環境で backups/trash/ へファイルを退避します
func MoveToRecycleBin(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil
	}

	trashDir := filepath.Join("backups", "trash", time.Now().Format("20060102"))
	if err := os.MkdirAll(trashDir, 0755); err != nil {
		return err
	}
	destPath := filepath.Join(trashDir, filepath.Base(absPath))
	return os.Rename(absPath, destPath)
}
