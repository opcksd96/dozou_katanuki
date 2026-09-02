// app/app_rpc_stash_reconcile_media.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"
)

// ReconcileStashMedia は Scenes と Images の両方を照合してバインド総数を返します
func (a *App) ReconcileStashMedia() (int, error) {
	if err := a.WaitForReady(); err != nil {
		return 0, err
	}
	scenesCount, _ := a.ReconcileStashScenes()
	imagesCount, _ := a.ReconcileStashImages()
	total := scenesCount + imagesCount

	if total > 0 {
		a.AppendPipelineLog("STASH", "SUCCESS", fmt.Sprintf("Stash照合完了: 計%d件バインド (Scene:%d, Image:%d)", total, scenesCount, imagesCount))
	}
	return total, nil
}
