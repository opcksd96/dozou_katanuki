// middleware/broadcast_handler_article_trash.go (Under 100 lines - SPEC-PRINCIPLE-001)
package middleware

import (
	"encoding/json"
	"net/http"
	
	"dozou_katanuki/domain/ports"
)

type TrashArticleRequest struct {
	ID        string `json:"id"`
	TrashedBy string `json:"trashed_by"`
	Reason    string `json:"reason"`
}

type RestoreArticleRequest struct {
	ID string `json:"id"`
}

func (s *BroadcastService) handleArticleTrashAPI(w http.ResponseWriter, r *http.Request) {
	if !ports.IsAdmin(r.Context()) { http.Error(w, `{"error":"Forbidden: Requires Admin scope"}`, http.StatusForbidden); return }
	if s.timelineService == nil {
		http.Error(w, `{"error":"Timeline service not available"}`, http.StatusServiceUnavailable)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var req TrashArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ID == "" {
		http.Error(w, `{"error":"Invalid request payload, id is required"}`, http.StatusBadRequest)
		return
	}
	if req.TrashedBy == "" { req.TrashedBy = "admin" }

	if err := s.timelineService.TrashArticle(req.ID, req.TrashedBy, req.Reason); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "id": req.ID})
}

func (s *BroadcastService) handleArticleRestoreAPI(w http.ResponseWriter, r *http.Request) {
	if !ports.IsAdmin(r.Context()) { http.Error(w, `{"error":"Forbidden: Requires Admin scope"}`, http.StatusForbidden); return }
	if s.timelineService == nil {
		http.Error(w, `{"error":"Timeline service not available"}`, http.StatusServiceUnavailable)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var req RestoreArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ID == "" {
		http.Error(w, `{"error":"Invalid request payload, id is required"}`, http.StatusBadRequest)
		return
	}
	if err := s.timelineService.RestoreArticle(req.ID); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "id": req.ID})
}
