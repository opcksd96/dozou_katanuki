// stash_config_sync.go (100行以下)
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"dozou_katanuki/models"
)

// SyncStashConfig は dozou の AppConfig を基に Stash の config.yml (bin/config.yml) を透過的に同期します
func SyncStashConfig(cfg *models.AppConfig) error {
	if cfg == nil {
		return nil
	}

	configYmlCandidates := []string{
		filepath.Join("bin", "config.yml"),
		filepath.Join(".", "config.yml"),
	}

	var targetPath string
	for _, p := range configYmlCandidates {
		if _, err := os.Stat(p); err == nil {
			targetPath = p
			break
		}
	}
	if targetPath == "" {
		targetPath = filepath.Join("bin", "config.yml")
	}

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

func updateStashConfigYAML(content string, host string, port int) (string, bool) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	var lines []string
	foundHost, foundPort, foundAuth, changed := false, false, false, false

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "host:") {
			newLine := fmt.Sprintf("host: %s", host)
			if line != newLine {
				line = newLine
				changed = true
			}
			foundHost = true
		} else if strings.HasPrefix(trimmed, "port:") {
			newLine := fmt.Sprintf("port: %d", port)
			if line != newLine {
				line = newLine
				changed = true
			}
			foundPort = true
		} else if strings.HasPrefix(trimmed, "dangerous_allow_public_without_auth:") {
			newLine := `dangerous_allow_public_without_auth: "true"`
			if line != newLine {
				line = newLine
				changed = true
			}
			foundAuth = true
		}
		lines = append(lines, line)
	}

	if !foundHost {
		lines = append(lines, fmt.Sprintf("host: %s", host))
		changed = true
	}
	if !foundPort {
		lines = append(lines, fmt.Sprintf("port: %d", port))
		changed = true
	}
	if !foundAuth {
		lines = append(lines, `dangerous_allow_public_without_auth: "true"`)
		changed = true
	}

	result := strings.Join(lines, "\n") + "\n"
	return result, changed
}
