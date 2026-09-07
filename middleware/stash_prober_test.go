// middleware/stash_prober_test.go (100行以下 - SPEC-PRINCIPLE-001)
package middleware

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestBeaconServiceLifecycle(t *testing.T) {
	var mu sync.Mutex
	events := make(map[string]int)

	emitter := func(event string, data ...interface{}) {
		mu.Lock()
		defer mu.Unlock()
		events[event]++
	}

	service := NewBeaconService(emitter)
	state := service.GetState()
	if state.StashReady || state.ThunderReady || state.MotrixReady {
		t.Errorf("expected initial state to be false, got %+v", state)
	}

	ctx, cancel := context.WithCancel(context.Background())
	service.Start(ctx)
	time.Sleep(50 * time.Millisecond)
	cancel()
}
