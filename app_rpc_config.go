// app_rpc_config.go (100行以下)
package main

import (
	"encoding/json"
	"log"
	"os"

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

	if a.broadcastService != nil {
		if err := a.broadcastService.UpdateConfig(cfg.Network, cfg.Broadcast); err != nil {
			log.Printf("[Wails RPC] SaveConfig broadcast service update warning: %v", err)
		}
	}
	return nil
}
