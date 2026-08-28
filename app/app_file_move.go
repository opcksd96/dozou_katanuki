// app/app_file_move.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"io"
	"os"
)

func moveFileSafe(src, dst string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
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

	if _, err := io.Copy(d, s); err != nil {
		return err
	}
	s.Close()
	_ = os.Remove(src)
	return nil
}
