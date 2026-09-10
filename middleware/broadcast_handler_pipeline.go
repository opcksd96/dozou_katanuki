// middleware/broadcast_handler_pipeline.go (100行以下 - SPEC-PRINCIPLE-001)
package middleware

import (
	"encoding/json"
	"net/http"

	"dozou_katanuki/domain/ports"
)

// handlePipelineCycleNowAPI handles POST /api/admin/pipeline/cycle-now
func (s *BroadcastService) handlePipelineCycleNowAPI(w http.ResponseWriter, r *http.Request) {
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

	res, err := s.adminUseCases.ExecutePipelineCycleNow()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(res)
}
