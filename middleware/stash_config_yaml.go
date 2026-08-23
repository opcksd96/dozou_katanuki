// middleware/stash_config_yaml.go (100行以下)
package middleware

import (
	"bufio"
	"fmt"
	"strings"
)

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

	return strings.Join(lines, "\n") + "\n", changed
}
