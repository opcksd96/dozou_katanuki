// app_rpc_skin.go (100行以下)
package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"dozou_katanuki/models"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func cleanPlatformPath(platform string) string {
	if platform == "" { platform = "twitter" }
	cp := filepath.Clean(platform)
	if strings.Contains(cp, "..") || strings.ContainsAny(cp, "/\\") { return "twitter" }
	return cp
}

// GetSkinCSS はプラグインの skin/design.css (SPEC-PLUGIN-001) を取得する Wails バインドメソッドです
func (a *App) GetSkinCSS(platform string) (string, error) {
	skinPath := filepath.Join("plugins", cleanPlatformPath(platform), "skin", "design.css")
	data, err := os.ReadFile(skinPath)
	if err != nil {
		log.Printf("[Wails RPC] GetSkinCSS error reading %s: %v", skinPath, err)
		return "", err
	}
	return string(data), nil
}

// SaveSkinCSS はプラグインの skin/design.css (SPEC-PLUGIN-001) を上書き保存する Wails バインドメソッドです
func (a *App) SaveSkinCSS(platform, cssContent string) error {
	p := cleanPlatformPath(platform)
	skinDir := filepath.Join("plugins", p, "skin")
	if err := os.MkdirAll(skinDir, 0755); err != nil {
		return err
	}
	skinPath := filepath.Join(skinDir, "design.css")
	if err := os.WriteFile(skinPath, []byte(cssContent), 0644); err != nil {
		return err
	}
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, "skin:updated", map[string]string{"platform": p, "css": cssContent})
	}
	return nil
}

// GetSkinLayout はプラグインの skin/layout.yaml を取得する Wails バインドメソッドです
func (a *App) GetSkinLayout(platform string) (string, error) {
	layoutPath := filepath.Join("plugins", cleanPlatformPath(platform), "skin", "layout.yaml")
	data, err := os.ReadFile(layoutPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GetSkinController はプラグインの skin/controller.js を取得する Wails バインドメソッドです
func (a *App) GetSkinController(platform string) (string, error) {
	ctrlPath := filepath.Join("plugins", cleanPlatformPath(platform), "skin", "controller.js")
	data, err := os.ReadFile(ctrlPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// GetSkinPackage はプラグインの skin アセット一式を取得する Wails バインドメソッドです
func (a *App) GetSkinPackage(platform string) (*models.SkinPackage, error) {
	p := cleanPlatformPath(platform)
	layout, _ := a.GetSkinLayout(p)
	css, _ := a.GetSkinCSS(p)
	ctrl, _ := a.GetSkinController(p)
	return &models.SkinPackage{Platform: p, LayoutYAML: layout, DesignCSS: css, Controller: ctrl}, nil
}
