// middleware/broadcast_handler_thunder.go (100行以下 - SPEC-PRINCIPLE-001)
package middleware

import (
	"encoding/json"
	"net/http"

	"dozou_katanuki/domain/ports"
)

// handleLaunchThunderAPI handles POST /api/admin/pipeline/launch-thunder
func (s *BroadcastService) handleLaunchThunderAPI(w http.ResponseWriter, r *http.Request) {
	if ports.GetScope(r.Context()) != ports.ScopeAdmin {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	if s.adminUseCases == nil {
		http.Error(w, "Admin use cases not initialized", http.StatusInternalServerError)
		return
	}

	ok, err := s.adminUseCases.LaunchThunder()
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]bool{"success": ok})
}
