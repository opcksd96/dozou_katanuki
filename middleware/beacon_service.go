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
	mu        sync.RWMutex
}

func NewBeaconService(pushEvent func(string, ...interface{})) *BeaconService {
	return &BeaconService{
		pushEvent: pushEvent,
	}
}

func (s *BeaconService) GetState() BeaconState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state
}

func (s *BeaconService) Start(ctx context.Context) {
	fmt.Println("[BeaconService] 起動: 外部プロセスのポート監視を開始します...")
	go s.loop(ctx)
}

func (s *BeaconService) loop(ctx context.Context) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	s.checkAll()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("[BeaconService] 停止")
			return
		case <-ticker.C:
			s.checkAll()
		}
	}
}

func (s *BeaconService) checkAll() {
	stashOK := checkPort(9999)
	thunderOK := checkPort(9222)
	motrixOK := checkPort(16800)

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.state.StashReady != stashOK {
		s.state.StashReady = stashOK
		s.pushEvent("stash:ready", stashOK)
		if stashOK {
			s.pushEvent("toast:notify", map[string]string{"type": "success", "message": "🟢 Stash 接続完了 (Beacon)！"})
		}
	}

	if s.state.ThunderReady != thunderOK {
		s.state.ThunderReady = thunderOK
		s.pushEvent("thunder:ready", thunderOK)
		if thunderOK {
			s.pushEvent("toast:notify", map[string]string{"type": "success", "message": "⚡ 迅雷 CDP 接続完了 (Beacon)！"})
		}
	}

	if s.state.MotrixReady != motrixOK {
		s.state.MotrixReady = motrixOK
		s.pushEvent("motrix:ready", motrixOK)
		if motrixOK {
			s.pushEvent("toast:notify", map[string]string{"type": "success", "message": "🚀 Motrix 接続完了 (Beacon)！"})
		}
	}
}

func checkPort(port int) bool {
	address := fmt.Sprintf("127.0.0.1:%d", port)
	conn, err := net.DialTimeout("tcp", address, 500*time.Millisecond)
	if err != nil {
		return false
	}
	if conn != nil {
		_ = conn.Close()
		return true
	}
	return false
}
