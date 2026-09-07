// middleware/beacon_service.go (100行以下 - SPEC-PRINCIPLE-001)
package middleware

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

type BeaconState struct {
	StashReady   bool
	ThunderReady bool
	MotrixReady  bool
}

type BeaconService struct {
	pushEvent func(string, ...interface{})
	state     BeaconState
	pollCount int
	mu        sync.RWMutex
}

func NewBeaconService(pushEvent func(string, ...interface{})) *BeaconService {
	return &BeaconService{pushEvent: pushEvent}
}

func (s *BeaconService) GetState() BeaconState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state
}

func (s *BeaconService) Start(ctx context.Context) {
	go s.loop(ctx)
}

func (s *BeaconService) loop(ctx context.Context) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()
	s.checkAll()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.checkAll()
		}
	}
}

func (s *BeaconService) checkAll() {
	stashOK, thunderOK, motrixOK := checkPort(9999), checkPort(9222), checkPort(16800)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pollCount++

	if s.state.StashReady != stashOK || s.pollCount%10 == 1 {
		level := "INFO"
		if !stashOK { level = "WARN" }
		event := "stash_beacon_poll"
		if s.state.StashReady != stashOK { event = "stash_status_change" }
		GetGlobalJournal().Record("stash", level, event, fmt.Sprintf("Stash 疎通確認 (port 9999): ready=%v", stashOK), map[string]interface{}{"port": 9999, "ready": stashOK})
	}

	if s.state.StashReady != stashOK {
		s.state.StashReady = stashOK
		s.pushEvent("stash:ready", stashOK)
		if stashOK { s.pushEvent("toast:notify", map[string]string{"type": "success", "message": "🟢 Stash 接続完了 (Beacon)！"}) }
	}
	if s.state.ThunderReady != thunderOK {
		s.state.ThunderReady = thunderOK
		s.pushEvent("thunder:ready", thunderOK)
		if thunderOK { s.pushEvent("toast:notify", map[string]string{"type": "success", "message": "⚡ 迅雷 CDP 接続完了 (Beacon)！"}) }
		GetGlobalJournal().Record("downloader", "INFO", "thunder_cdp_status", fmt.Sprintf("迅雷 CDP 疎通状態変更: ready=%v", thunderOK), map[string]interface{}{"port": 9222, "ready": thunderOK})
	}
	if s.state.MotrixReady != motrixOK {
		s.state.MotrixReady = motrixOK
		s.pushEvent("motrix:ready", motrixOK)
		if motrixOK { s.pushEvent("toast:notify", map[string]string{"type": "success", "message": "🚀 Motrix 接続完了 (Beacon)！"}) }
	}
}

func checkPort(port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 500*time.Millisecond)
	if err != nil { return false }
	_ = conn.Close()
	return true
}
