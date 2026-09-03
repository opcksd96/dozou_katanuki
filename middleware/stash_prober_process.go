// middleware/stash_prober_process.go (100行以下 - SPEC-PRINCIPLE-001)
package middleware

import (
	"fmt"
)

// NOTE: v4.5.0 アーキテクチャ原則（PLAN-01）に基づき、
// Stash プロセスの起動・停止・強制終了はすべて外装の PS1 Top Conductor (scripts/dev.ps1 / run.ps1) が担当します。
// Go ランタイム内部から tasklist や taskkill を発行する処理は完全に無効化・撤廃されました。

func (p *StashProber) isStashProcessAlive() bool {
	// Go 内部での OS プロセス監視は行わず、HTTP ポートの接続状態 (p.connected) のみを参照
	return p.connected
}

func (p *StashProber) kickStashAsync() {
	// Go 内部からの直接起動は禁止（死の多重起動ループ防止）
	fmt.Println("[StashProber] kickStashAsync は廃止されました。外部スクリプト(dev.ps1/run.ps1)による管理に委譲しています。")
}

func (p *StashProber) Stop() {
	if p.cancelFunc != nil {
		p.cancelFunc()
	}
	p.running = false
	fmt.Println("[StashProber] StashProber 停止完了 (プロセス管理は外部Conductorに委譲)")
}
