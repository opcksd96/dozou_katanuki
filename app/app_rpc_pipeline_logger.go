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

	"dozou_katanuki/middleware"
)

type PipelineLogEntry struct {
	Timestamp string `json:"timestamp"`
	Stage     string `json:"stage"` // REQUESTS, MOTRIX, THUNDER, STASH, SYSTEM
	Level     string `json:"level"` // INFO, WARN, ERROR, SUCCESS
	Message   string `json:"message"`
}

func mapStageToComponent(stage string) string {
	switch strings.ToLower(stage) {
	case "stash": return "stash"
	case "system": return "system"
	case "requests": return "crawler"
	case "scraper": return "scraper"
	default: return "downloader"
	}
}

// AppendPipelineLog は アプリ別ログファイルにログを追記し、システムジャーナルへも同時記録します
func (a *App) AppendPipelineLog(stage, level, msg string) {
	lvl := strings.ToUpper(level)
	if lvl == "SUCCESS" { lvl = "INFO" }
	middleware.GetGlobalJournal().Record(
		mapStageToComponent(stage), lvl, strings.ToLower(stage)+"_event",
		msg, map[string]interface{}{"stage": stage, "level": level},
	)

	go func() {
		_ = os.MkdirAll("logs", 0755)
		p := filepath.Join("logs", fmt.Sprintf("%s.log", strings.ToLower(stage)))
		f, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil { return }
		defer f.Close()
		ts := time.Now().Format("2006-01-02 15:04:05")
		_, _ = f.WriteString(fmt.Sprintf("[%s] [%s] %s\n", ts, level, msg))
	}()
}

// GetPipelineLogs は 指定ステージ（空文字なら全ログ混合）の直近ログを取得します
func (a *App) GetPipelineLogs(stage string, limit int) ([]PipelineLogEntry, error) {
	if limit <= 0 { limit = 50 }
	var entries []PipelineLogEntry
	stages := []string{"scraper", "requests", "motrix", "thunder", "stash"}
	if stage != "" && stage != "all" { stages = []string{strings.ToLower(stage)} }

	for _, st := range stages {
		f, err := os.Open(filepath.Join("logs", fmt.Sprintf("%s.log", st)))
		if err != nil { continue }
		var fileEntries []PipelineLogEntry
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "[") && strings.Contains(line, "]") {
				parts := strings.SplitN(line, "] ", 2)
				if len(parts) == 2 {
					ts, msg, lvl := strings.TrimPrefix(parts[0], "["), strings.TrimSpace(parts[1]), "INFO"
					if strings.Contains(msg, "[") && strings.Contains(msg, "]") {
						lvlParts := strings.SplitN(msg, "] ", 2)
						lvl = strings.Trim(lvlParts[0], "[]")
						if len(lvlParts) > 1 { msg = strings.TrimSpace(lvlParts[1]) }
					}
					fileEntries = append(fileEntries, PipelineLogEntry{Timestamp: ts, Stage: strings.ToUpper(st), Level: lvl, Message: msg})
				}
			}
		}
		_ = scanner.Err()
		_ = f.Close()
		if len(fileEntries) > limit { fileEntries = fileEntries[len(fileEntries)-limit:] }
		entries = append(entries, fileEntries...)
	}

	sort.Slice(entries, func(i, j int) bool { return entries[i].Timestamp < entries[j].Timestamp })
	if len(entries) > limit { entries = entries[len(entries)-limit:] }
	return entries, nil
}
