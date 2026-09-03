// app/app_rpc_thunder_launcher_test.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"testing"
)

func TestThunderLauncher_StatusCheck(t *testing.T) {
	running := isThunderProcessRunning()
	cdp := isThunderCDPListening()
	t.Logf("Thunder Process Running: %v, CDP Listening: %v", running, cdp)

	app := &App{}
	// CDP疎通未確認環境でも、EnsureThunderCDP のインターフェース整合性を検証
	if running && cdp {
		ok, err := app.EnsureThunderCDP()
		if err != nil || !ok {
			t.Errorf("expected EnsureThunderCDP to succeed when already listening, got ok=%v, err=%v", ok, err)
		}
	}
}
