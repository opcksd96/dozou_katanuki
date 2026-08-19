package main

import (
	"os"
	"testing"
)

func TestConfigGetAndSave(t *testing.T) {
	// 元の config.json を退避してテスト後に復元
	originalData, err := os.ReadFile("config.json")
	if err != nil {
		t.Fatalf("Failed to read original config.json: %v", err)
	}
	defer func() {
		_ = os.WriteFile("config.json", originalData, 0644)
	}()

	app := &App{}

	// 1. GetConfig のテスト
	cfg, err := app.GetConfig()
	if err != nil {
		t.Fatalf("GetConfig failed: %v", err)
	}
	if cfg.System.Language == "" {
		t.Errorf("Expected language to not be empty, got %s", cfg.System.Language)
	}
	if cfg.Network.StashPort == 0 {
		t.Errorf("Expected StashPort to be non-zero, got %d", cfg.Network.StashPort)
	}

	// 2. SaveConfig のテスト（変更して保存）
	originalLang := cfg.System.Language
	cfg.System.Language = "zh"
	cfg.Storage.StashEnabled = false

	if err := app.SaveConfig(cfg); err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}

	// 3. 再度読み込んで検証
	updatedCfg, err := app.GetConfig()
	if err != nil {
		t.Fatalf("GetConfig after save failed: %v", err)
	}
	if updatedCfg.System.Language != "zh" {
		t.Errorf("Expected language 'zh', got '%s'", updatedCfg.System.Language)
	}
	if updatedCfg.Storage.StashEnabled != false {
		t.Errorf("Expected stash_enabled false, got %v", updatedCfg.Storage.StashEnabled)
	}

	// 4. 元に戻す確認
	cfg.System.Language = originalLang
	cfg.Storage.StashEnabled = true
	if err := app.SaveConfig(cfg); err != nil {
		t.Fatalf("SaveConfig reset failed: %v", err)
	}
}
