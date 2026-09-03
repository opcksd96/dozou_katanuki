// middleware/broadcast_server_beacon.go (100行以下 - SPEC-PRINCIPLE-001)
package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"dozou_katanuki/adapters/driving/dto"
)

func (s *BroadcastService) handleBeaconAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// セキュリティ: localhost (127.0.0.1 または ::1) からのアクセスのみ許可
	remoteIP := r.RemoteAddr
	if strings.Contains(remoteIP, ":") {
		remoteIP = strings.Split(remoteIP, ":")[0]
	}
	if remoteIP != "127.0.0.1" && remoteIP != "::1" && remoteIP != "[" {
		http.Error(w, `{"error":"Forbidden: Internal API"}`, http.StatusForbidden)
		return
	}

	var req dto.BeaconRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Bad request"}`, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Callback を経由して App 側へ伝達
	s.mu.RLock()
	cb := s.beaconCallback
	s.mu.RUnlock()

	if cb != nil {
		cb(req)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
