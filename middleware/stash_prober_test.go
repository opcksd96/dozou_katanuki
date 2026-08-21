package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestStashProberSuccess(t *testing.T) {
	var mu sync.Mutex
	readyEmitted := false
	toastEmitted := false

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}))
	defer ts.Close()

	emitter := func(event string, data ...interface{}) {
		mu.Lock()
		defer mu.Unlock()
		if event == "stash:ready" {
			readyEmitted = true
		}
		if event == "toast:notify" {
			toastEmitted = true
		}
	}

	prober := NewStashProber("./dummy.exe", ts.URL, emitter)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	prober.Start(ctx)
	time.Sleep(500 * time.Millisecond)

	mu.Lock()
	if !readyEmitted || !toastEmitted {
		t.Errorf("expected readyEmitted=true, toastEmitted=true, got ready=%v, toast=%v", readyEmitted, toastEmitted)
	}
	mu.Unlock()

	prober.Stop()
}
