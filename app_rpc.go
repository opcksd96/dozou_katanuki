package main

import (
	"encoding/json"
	"log"
	"os"

	"dozou_katanuki/models"
)

// GetAccounts は登録されている全アカウントのリストを供給する Wails バインドメソッドです
func (a *App) GetAccounts(platform string) ([]models.RenderAuthor, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	return a.timelineService.GetAccounts(platform)
}

// GetTimeline はフロントエンドへ RenderTree 配列を供給する Wails バインドメソッドです
func (a *App) GetTimeline(platform, accountID, filter string, limit, offset int) ([]models.RenderTree, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	res, err := a.timelineService.FetchTimeline(platform, accountID, filter, limit, offset)
	log.Printf("[Wails RPC] GetTimeline(platform=%s, accountID=%s, filter=%s) -> 取得件数: %d (err: %v)",
		platform, accountID, filter, len(res), err)
	return res, err
}

// GetSystemLanguage は config.json (SPEC-CONFIG-001) からシステム言語設定を取得します
func (a *App) GetSystemLanguage() string {
	type Config struct {
		System struct {
			Language string `json:"language"`
		} `json:"system"`
	}

	data, err := os.ReadFile("config.json")
	if err != nil {
		return "ja" // デフォルト言語
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil || cfg.System.Language == "" {
		return "ja"
	}
	return cfg.System.Language
}

// StartSalvageJob は自動サルベージジョブをキューに追加する Wails バインドメソッドです
func (a *App) StartSalvageJob(platform, account string, limit int) (*models.JobProgress, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	return a.jobOrchestrator.EnqueueSalvage(platform, account, limit)
}

// StartManualImportJob は手動 WARC インポートジョブをキューに追加する Wails バインドメソッドです
func (a *App) StartManualImportJob(warcPath string, offline bool) (*models.JobProgress, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	return a.jobOrchestrator.EnqueueManualImport(warcPath, offline)
}

// GetJobStatus は指定されたジョブのステータスを取得する Wails バインドメソッドです
func (a *App) GetJobStatus(jobID string) (*models.JobProgress, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	st := a.jobOrchestrator.GetStatus(jobID)
	return st, nil
}

// GetActiveJob は現在実行中のジョブを取得する Wails バインドメソッドです
func (a *App) GetActiveJob() (*models.JobProgress, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	return a.jobOrchestrator.GetActiveJob(), nil
}

// ListJobs は全ジョブの履歴ステータス一覧を取得する Wails バインドメソッドです
func (a *App) ListJobs() ([]*models.JobProgress, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	return a.jobOrchestrator.ListJobs(), nil
}

// CancelJob はジョブの実行を中断する Wails バインドメソッドです
func (a *App) CancelJob(jobID string) error {
	if err := a.waitForReady(); err != nil {
		return err
	}
	return a.jobOrchestrator.CancelJob(jobID)
}

// GetConfig は config.json (SPEC-CONFIG-001) の全設定を取得する Wails バインドメソッドです
func (a *App) GetConfig() (*models.AppConfig, error) {
	data, err := os.ReadFile("config.json")
	if err != nil {
		return nil, err
	}
	var cfg models.AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// SaveConfig は config.json (SPEC-CONFIG-001) を更新・保存する Wails バインドメソッドです
func (a *App) SaveConfig(cfg *models.AppConfig) error {
	if cfg == nil {
		return nil
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	// 改行を付与して保存
	data = append(data, '\n')
	if err := os.WriteFile("config.json", data, 0644); err != nil {
		return err
	}
	log.Printf("[Wails RPC] SaveConfig: config.json updated successfully")
	return nil
}
