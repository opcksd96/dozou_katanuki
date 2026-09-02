// app/app_rpc_motrix_launch.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var candidateMotrixPaths = []string{
	`C:\Program Files\MotrixNext\motrix-next.exe`,
	`C:\Program Files\Motrix\Motrix.exe`,
	`C:\Program Files (x86)\MotrixNext\motrix-next.exe`,
	`C:\Program Files (x86)\Motrix\Motrix.exe`,
}

// LaunchMotrix は Motrix Next / Motrix の実行可能ファイルを探索して起動します
func (a *App) LaunchMotrix() (bool, error) {
	paths := append([]string{}, candidateMotrixPaths...)
	if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
		paths = append(paths,
			filepath.Join(localAppData, "Programs", "MotrixNext", "motrix-next.exe"),
			filepath.Join(localAppData, "Programs", "Motrix", "Motrix.exe"),
			filepath.Join(localAppData, "Programs", "motrix-next", "motrix-next.exe"),
		)
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			cmd := exec.Command(p)
			if err := cmd.Start(); err == nil {
				return true, nil
			}
		}
	}

	return false, fmt.Errorf("motrix executable not found in candidate paths")
}
