// app/app_rpc_config_test.go (100行以下)
package app

import (
	"os"
	"testing"
)

func TestConfigGetAndSave(t *testing.T) {
	// 元の config.json を退避してテスト後に復元
	originalData, err := os.ReadFile("config.json")
	if err != nil {
		originalData, _ = os.ReadFile("../config.json")
	}
	defer func() {
		if len(originalData) > 0 {
			_ = os.WriteFile("config.json", originalData, 0644)
		}
	}()

	app := &App{}

	// 1. GetConfig のテスト
	cfg, err := app.GetConfig()
	if err != nil {
		// ルートの config.json を参照できるようにする
		_ = os.WriteFile("config.json", originalData, 0644)
		cfg, err = app.GetConfig()
		if err != nil {
			t.Fatalf("GetConfig failed: %v", err)
		}
	}
	if cfg.System.Language == "" {
		t.Errorf("Expected language to not be empty, got %s", cfg.System.Language)
	}

	// 2. SaveConfig のテスト（変更して保存）
	originalLang := cfg.System.Language
	cfg.System.Language = "zh"
	cfg.Storage.StashEnabled = false
	cfg.Translation.DeeplApiKey = "test-deepl-key:fx"
	cfg.Translation.GoogleTranslateApiKey = "test-google-key"

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
	if updatedCfg.Translation.DeeplApiKey != "test-deepl-key:fx" {
		t.Errorf("Expected deepl key 'test-deepl-key:fx', got '%s'", updatedCfg.Translation.DeeplApiKey)
	}
	env := app.getTranslationEnv()
	if env["DEEPL_API_KEY"] != "test-deepl-key:fx" || env["GOOGLE_TRANSLATE_API_KEY"] != "test-google-key" {
		t.Errorf("Unexpected translation env: %v", env)
	}

	// 4. 元に戻す確認
	cfg.System.Language = originalLang
	cfg.Storage.StashEnabled = true
	cfg.Translation.DeeplApiKey = ""
	cfg.Translation.GoogleTranslateApiKey = ""
	if err := app.SaveConfig(cfg); err != nil {
		t.Fatalf("SaveConfig reset failed: %v", err)
	}
}
