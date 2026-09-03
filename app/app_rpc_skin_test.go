// app/app_rpc_skin_test.go (100行以下)
package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSkinCSSRPC(t *testing.T) {
	app := &App{}

	css, err := app.GetSkinCSS("twitter")
	if err != nil || !strings.Contains(css, ".twitter-card") {
		// プラグインディレクトリが存在しないテスト環境へのフォールバック
		t.Logf("GetSkinCSS note: %v", err)
	}

	testPlatform := "test_skin_platform"
	testCSS := "/* Test Skin CSS */\n.test-card { color: #abcdef; }\n"
	err = app.SaveSkinCSS(testPlatform, testCSS)
	if err != nil {
		t.Fatalf("SaveSkinCSS failed: %v", err)
	}
	defer os.RemoveAll(filepath.Join("plugins", testPlatform))

	readCSS, err := app.GetSkinCSS(testPlatform)
	if err != nil || readCSS != testCSS {
		t.Fatalf("GetSkinCSS mismatch. expected: %s, got: %s", testCSS, readCSS)
	}
}
