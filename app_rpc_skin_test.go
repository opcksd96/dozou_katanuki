package main

import (
	"testing"
)

func TestApp_GetSkinPackage(t *testing.T) {
	app := &App{}

	// plugins/twitter/skin フォルダとファイルが存在することを確認
	pkg, err := app.GetSkinPackage("twitter")
	if err != nil {
		t.Fatalf("GetSkinPackage failed: %v", err)
	}

	if pkg == nil {
		t.Fatal("expected non-nil SkinPackage")
	}

	if pkg.Platform != "twitter" {
		t.Errorf("expected platform twitter, got %s", pkg.Platform)
	}

	if pkg.LayoutYAML == "" {
		t.Error("expected non-empty LayoutYAML")
	}

	if pkg.DesignCSS == "" {
		t.Error("expected non-empty DesignCSS")
	}

	if pkg.Controller == "" {
		t.Error("expected non-empty Controller")
	}
}

func TestApp_SkinSecurity_Traversal(t *testing.T) {
	app := &App{}

	// 不正なパス（パストラバーサル）を指定しても安全にフォールバックされること
	pkg, err := app.GetSkinPackage("../../../etc")
	if err != nil {
		t.Fatalf("GetSkinPackage traversal check failed: %v", err)
	}
	if pkg.Platform != "twitter" {
		t.Errorf("expected fallback to twitter platform, got %s", pkg.Platform)
	}
}
