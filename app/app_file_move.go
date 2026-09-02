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

// moveFileSafe はファイルを移動します（同ボリュームならRename、失敗時はコピー後に元ファイルを削除）
func moveFileSafe(src, dst string) error {
	if src == dst {
		return nil
	}
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
	if err := copyFileSafe(src, dst); err != nil {
		return err
	}
	return os.Remove(src)
}
