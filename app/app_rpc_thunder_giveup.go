// app/app_rpc_thunder_giveup.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import "fmt"

// GiveUpRetainedMedia はユーザーが明示的に諦めたタスクを DEAD_404 状態へ確定します
func (a *App) GiveUpRetainedMedia(mediaID string) (bool, error) {
	if mediaID == "" || a.Repo == nil { return false, fmt.Errorf("mediaID is required") }
	err := a.Repo.UpdateMediaMetadata(mediaID, "DEAD_404", "", "", "ユーザーによる探索諦め (GIVE_UP)")
	return err == nil, err
}
