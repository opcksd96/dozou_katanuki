// app/app_thunder_cdp_ws.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

// EvaluateCDPExpression は WebSocket (gorilla/websocket) 経由で Runtime.evaluate を実行します
func EvaluateCDPExpression(wsURL, expr string, timeout time.Duration) (string, error) {
	return EvaluateCDP(wsURL, expr, timeout)
}

// EvaluateCDP は WebSocket (gorilla/websocket) 経由で Runtime.evaluate を実行します
func EvaluateCDP(wsURL, expr string, timeout time.Duration) (string, error) {
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	dialer := websocket.Dialer{HandshakeTimeout: timeout}
	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to dial CDP websocket: %w", err)
	}
	defer conn.Close()

	_ = conn.SetWriteDeadline(time.Now().Add(timeout))
	req := map[string]interface{}{
		"id": 1, "method": "Runtime.evaluate",
		"params": map[string]interface{}{"expression": expr, "returnByValue": true},
	}
	if err := conn.WriteJSON(req); err != nil {
		return "", fmt.Errorf("failed to write CDP request: %w", err)
	}

	_ = conn.SetReadDeadline(time.Now().Add(timeout))
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return "", fmt.Errorf("failed to read CDP response: %w", err)
		}
		var base struct {
			ID int `json:"id"`
		}
		if err := json.Unmarshal(msg, &base); err == nil && base.ID == 1 {
			return string(msg), nil
		}
	}
}
