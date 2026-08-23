// middleware/stash_config_sync.go (100行以下)
package middleware

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"dozou_katanuki/models"
)

// SyncStashConfig は dozou の AppConfig を基に Stash の config.yml (bin/stash/config.yml) を透過的に同期します
func SyncStashConfig(cfg *models.AppConfig) error {
	if cfg == nil {
		return nil
	}

	targetPath := findStashConfigPath()
	data, err := os.ReadFile(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			initial := fmt.Sprintf("host: %s\nport: %d\ndangerous_allow_public_without_auth: \"true\"\n",
				getStashHost(cfg), getStashPort(cfg))
			_ = os.MkdirAll(filepath.Dir(targetPath), 0755)
			return os.WriteFile(targetPath, []byte(initial), 0644)
		}
		return err
	}

	updated, changed := updateStashConfigYAML(string(data), getStashHost(cfg), getStashPort(cfg))
	if changed {
		if err := os.WriteFile(targetPath, []byte(updated), 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", targetPath, err)
		}
		log.Printf("[StashSync] %s successfully synced (host: %s, port: %d)", targetPath, getStashHost(cfg), getStashPort(cfg))
	}
	return nil
}

func findStashConfigPath() string {
	candidates := []string{
		filepath.Join("bin", "stash", "config.yml"),
		filepath.Join("bin", "config.yml"),
		filepath.Join(".", "config.yml"),
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return filepath.Join("bin", "stash", "config.yml")
}

func getStashHost(cfg *models.AppConfig) string {
	if cfg.Network.PublicBindAddress != "" {
		return cfg.Network.PublicBindAddress
	}
	return "0.0.0.0"
}

func getStashPort(cfg *models.AppConfig) int {
	if cfg.Network.StashPort > 0 {
		return cfg.Network.StashPort
	}
	return 9999
}
