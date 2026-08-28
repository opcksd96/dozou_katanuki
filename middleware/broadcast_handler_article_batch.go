// middleware/broadcast_handler_article_batch.go (Under 100 lines - SPEC-PRINCIPLE-001)
package middleware

import (
	"encoding/json"
	"net/http"
)

type BatchTrashRequest struct {
	IDs       []string `json:"ids"`
	TrashedBy string   `json:"trashed_by"`
	Reason    string   `json:"reason"`
}

type BatchIDsRequest struct {
	IDs []string `json:"ids"`
}

func (s *BroadcastService) handleArticleBatchTrashAPI(w http.ResponseWriter, r *http.Request) {
	if s.timelineService == nil { http.Error(w, `{"error":"Timeline service unavailable"}`, 503); return }
	if r.Method != http.MethodPost { http.Error(w, `{"error":"Method not allowed"}`, 405); return }
	var req BatchTrashRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.IDs) == 0 {
		http.Error(w, `{"error":"Invalid request payload, ids required"}`, 400); return
	}
	if req.TrashedBy == "" { req.TrashedBy = "admin" }
	if err := s.timelineService.BatchTrashArticles(req.IDs, req.TrashedBy, req.Reason); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8"); w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "count": len(req.IDs)})
}

func (s *BroadcastService) handleArticleBatchRestoreAPI(w http.ResponseWriter, r *http.Request) {
	if s.timelineService == nil { http.Error(w, `{"error":"Timeline service unavailable"}`, 503); return }
	if r.Method != http.MethodPost { http.Error(w, `{"error":"Method not allowed"}`, 405); return }
	var req BatchIDsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.IDs) == 0 {
		http.Error(w, `{"error":"Invalid request payload, ids required"}`, 400); return
	}
	if err := s.timelineService.BatchRestoreArticles(req.IDs); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8"); w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "count": len(req.IDs)})
}

func (s *BroadcastService) handleArticleBatchResetTranslationsAPI(w http.ResponseWriter, r *http.Request) {
	if s.timelineService == nil { http.Error(w, `{"error":"Timeline service unavailable"}`, 503); return }
	if r.Method != http.MethodPost { http.Error(w, `{"error":"Method not allowed"}`, 405); return }
	var req BatchIDsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.IDs) == 0 {
		http.Error(w, `{"error":"Invalid request payload, ids required"}`, 400); return
	}
	if err := s.timelineService.BatchResetTranslations(req.IDs); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8"); w.WriteHeader(500)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "count": len(req.IDs)})
}
