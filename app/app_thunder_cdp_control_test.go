// app/app_thunder_cdp_control_test.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"testing"
)

func TestThunderCDP_ControlTask(t *testing.T) {
	a := &App{}
	st := a.GetThunderCDPStatus()
	if !st.IsConnected || len(st.CapturedTasks) == 0 {
		t.Skip("Thunder CDP not connected or no tasks found")
		return
	}

	targetTask := st.CapturedTasks[0].FileName
	t.Logf("Testing CDP Pause on: %s", targetTask)

	ok, err := a.ControlThunderTaskViaCDP(targetTask, "pause")
	t.Logf("Pause Result: ok=%v, err=%v", ok, err)

	t.Logf("Testing CDP Resume on: %s", targetTask)
	okResume, errResume := a.ControlThunderTaskViaCDP(targetTask, "resume")
	t.Logf("Resume Result: ok=%v, err=%v", okResume, errResume)

	t.Logf("Testing CDP Task Detail on: %s", targetTask)
	dt, errDt := a.GetThunderTaskDetailViaCDP(targetTask)
	if dt != nil {
		t.Logf("Detail Result: success=%v, url=%s, savePath=%s, err=%v", dt.Success, dt.DownloadURL, dt.SavePath, errDt)
	} else {
		t.Logf("Detail Result: nil (err=%v)", errDt)
	}

	t.Logf("Testing CDP Task Restore from Recycle Bin")
	okRestore, errRestore := a.ControlThunderTaskViaCDP("", "restore")
	t.Logf("Restore Result: ok=%v, err=%v", okRestore, errRestore)
}
