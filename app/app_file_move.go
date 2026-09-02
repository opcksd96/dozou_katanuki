// app/app_file_move.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"io"
	"os"
)

// copyFileSafe はソースファイルをコピーします（元ファイルは削除しません）
func copyFileSafe(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()

	_, err = io.Copy(d, s)
	return err
}
