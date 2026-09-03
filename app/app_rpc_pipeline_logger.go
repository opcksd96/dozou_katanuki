// app/app_rpc_pipeline_logger.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type PipelineLogEntry struct {
	Timestamp string `json:"timestamp"`
	Stage     string `json:"stage"` // REQUESTS, MOTRIX, THUNDER, STASH, SYSTEM
	Level     string `json:"level"` // INFO, WARN, ERROR, SUCCESS
	Message   string `json:"message"`
}

// AppendPipelineLog は アプリ別ログファイルにログを追記します
func (a *App) AppendPipelineLog(stage, level, msg string) {
	go func() {
		logDir := "logs"
		_ = os.MkdirAll(logDir, 0755)
		filePath := filepath.Join(logDir, fmt.Sprintf("%s.log", strings.ToLower(stage)))
		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return
		}
		defer f.Close()

		ts := time.Now().Format("2006-01-02 15:04:05")
		line := fmt.Sprintf("[%s] [%s] %s\n", ts, level, msg)
		_, _ = f.WriteString(line)
	}()
}

// GetPipelineLogs は 指定ステージ（空文字なら全ログ混合）の直近ログを取得します
func (a *App) GetPipelineLogs(stage string, limit int) ([]PipelineLogEntry, error) {
	if limit <= 0 {
		limit = 50
	}
	var entries []PipelineLogEntry

	stages := []string{"requests", "motrix", "thunder", "stash"}
	if stage != "" && stage != "all" {
		stages = []string{strings.ToLower(stage)}
	}

	for _, st := range stages {
		filePath := filepath.Join("logs", fmt.Sprintf("%s.log", st))
		f, err := os.Open(filePath)
		if err != nil {
			continue
		}

		var fileEntries []PipelineLogEntry
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "[") && strings.Contains(line, "]") {
				parts := strings.SplitN(line, "] ", 2)
				if len(parts) == 2 {
					header := strings.TrimPrefix(parts[0], "[")
					msg := strings.TrimSpace(parts[1])
					ts := header
					lvl := "INFO"
					if strings.Contains(msg, "[") && strings.Contains(msg, "]") {
						lvlParts := strings.SplitN(msg, "] ", 2)
						lvl = strings.Trim(lvlParts[0], "[]")
						if len(lvlParts) > 1 {
							msg = strings.TrimSpace(lvlParts[1])
						}
					}
					fileEntries = append(fileEntries, PipelineLogEntry{
						Timestamp: ts,
						Stage:     strings.ToUpper(st),
						Level:     lvl,
						Message:   msg,
					})
				}
			}
		}
		f.Close()

		if len(fileEntries) > limit {
			fileEntries = fileEntries[len(fileEntries)-limit:]
		}
		entries = append(entries, fileEntries...)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp < entries[j].Timestamp
	})

	if len(entries) > limit {
		entries = entries[len(entries)-limit:]
	}
	return entries, nil
}
