// middleware/mission_reporter.go (100行以下 - SPEC-PRINCIPLE-001)
package middleware

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"dozou_katanuki/models"
)

const maxReportGenerations = 5
const reportsDir = "reports"

// SaveMissionReport は 5W1H Markdown レポートをファイル保存し、5世代ローテーションを実施します
func SaveMissionReport(rep *models.MissionReport) (string, error) {
	if err := os.MkdirAll(reportsDir, 0755); err != nil { return "", err }
	sanitizedAcc := strings.TrimPrefix(rep.Account, "@")
	if sanitizedAcc == "" { sanitizedAcc = "unknown" }

	ts := rep.FinishedAt.Format("20060102_150405")
	rep.ID = fmt.Sprintf("REP-%s", ts)
	rep.FileName = fmt.Sprintf("mission_%s_%s.md", ts, sanitizedAcc)
	targetPath := filepath.Join(reportsDir, rep.FileName)

	if err := os.WriteFile(targetPath, []byte(rep.MarkdownText), 0644); err != nil {
		return "", err
	}
	_ = rotateReports(maxReportGenerations)
	return targetPath, nil
}

// rotateReports は reports ディレクトリ内のファイルを最新 N 世代に保ちます
func rotateReports(maxKeep int) error {
	entries, err := os.ReadDir(reportsDir)
	if err != nil { return err }

	var files []os.DirEntry
	for _, e := range entries {
		if !e.IsDir() && strings.HasPrefix(e.Name(), "mission_") && strings.HasSuffix(e.Name(), ".md") {
			files = append(files, e)
		}
	}
	if len(files) <= maxKeep { return nil }

	sort.Slice(files, func(i, j int) bool {
		iInfo, _ := files[i].Info(); jInfo, _ := files[j].Info()
		if iInfo == nil || jInfo == nil { return false }
		return iInfo.ModTime().After(jInfo.ModTime())
	})

	for i := maxKeep; i < len(files); i++ {
		_ = os.Remove(filepath.Join(reportsDir, files[i].Name()))
	}
	return nil
}

// GetLatestReport は最新のミッションレポート（Markdown）を読み込んで返します
func GetLatestReport() (*models.MissionReport, error) {
	entries, err := os.ReadDir(reportsDir)
	if err != nil || len(entries) == 0 { return nil, fmt.Errorf("no reports found") }

	var validFiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasPrefix(e.Name(), "mission_") && strings.HasSuffix(e.Name(), ".md") {
			validFiles = append(validFiles, e.Name())
		}
	}
	if len(validFiles) == 0 { return nil, fmt.Errorf("no reports found") }
	sort.Sort(sort.Reverse(sort.StringSlice(validFiles)))

	latestName := validFiles[0]
	content, err := os.ReadFile(filepath.Join(reportsDir, latestName))
	if err != nil { return nil, err }

	return &models.MissionReport{
		ID:           latestName,
		FileName:     latestName,
		MarkdownText: string(content),
		FinishedAt:   time.Now(),
	}, nil
}
