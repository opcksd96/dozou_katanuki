// middleware/broadcast_server_sse.go (100行以下 - SPEC-PRINCIPLE-001)
package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (s *BroadcastService) handleEventsAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	// CORS for development (allowing Vite dev server to connect)
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	clientChan := make(chan []byte, 100)
	s.mu.Lock()
	if s.clients == nil {
		s.clients = make(map[chan []byte]bool)
	}
	s.clients[clientChan] = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.clients, clientChan)
		s.mu.Unlock()
		close(clientChan)
		log.Println("[SSE] Client disconnected")
	}()

	log.Println("[SSE] New client connected")
	
	// Send initial dummy event to establish connection
	fmt.Fprintf(w, "event: connected\ndata: {}\n\n")
	flusher.Flush()

	for {
		select {
		case <-r.Context().Done():
			return
		case payload := <-clientChan:
			fmt.Fprintf(w, "data: %s\n\n", payload)
			flusher.Flush()
		}
	}
}

func (s *BroadcastService) PushEvent(eventName string, data ...interface{}) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.clients) == 0 {
		return
	}

	payload := map[string]interface{}{
		"name": eventName,
		"data": data,
	}
	
	bytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("[SSE] Error marshaling event payload: %v", err)
		return
	}

	for clientChan := range s.clients {
		select {
		case clientChan <- bytes:
		default:
			// Non-blocking send; if channel is full, skip
		}
	}
}
