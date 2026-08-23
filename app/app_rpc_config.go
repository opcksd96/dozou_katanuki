// app/app_rpc_config.go (100行以下)
package app

import (
	"encoding/json"
	"log"
	"os"

	"dozou_katanuki/middleware"
	"dozou_katanuki/models"
)

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
	data = append(data, '\n')
	if err := os.WriteFile("config.json", data, 0644); err != nil {
		return err
	}
	log.Printf("[Wails RPC] SaveConfig: config.json updated successfully")

	if a.JobOrchestrator != nil {
		a.JobOrchestrator.SetStorageConfig(cfg.Storage)
	}
	if a.UnifiedHandler != nil {
		a.UnifiedHandler.SetMediaDir(cfg.Storage.LocalMediaDir)
	}
	if a.BroadcastService != nil {
		if err := a.BroadcastService.UpdateConfig(cfg.Network, cfg.Broadcast); err != nil {
			log.Printf("[Wails RPC] SaveConfig broadcast service update warning: %v", err)
		}
	}
	if err := middleware.SyncStashConfig(cfg); err != nil {
		log.Printf("[Wails RPC] SaveConfig Stash config sync warning: %v", err)
	}
	return nil
}

// getTranslationEnv は config.json の翻訳設定から Python サブプロセス用の環境変数を生成します
func (a *App) getTranslationEnv() map[string]string {
	cfg, err := a.GetConfig()
	if err != nil || cfg == nil {
		return nil
	}
	env := make(map[string]string)
	if cfg.Translation.DeeplApiKey != "" {
		env["DEEPL_API_KEY"] = cfg.Translation.DeeplApiKey
	}
	if cfg.Translation.GoogleTranslateApiKey != "" {
		env["GOOGLE_TRANSLATE_API_KEY"] = cfg.Translation.GoogleTranslateApiKey
	}
	return env
}
