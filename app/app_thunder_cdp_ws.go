// app/app_thunder_cdp_ws.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"
)

// EvaluateCDPExpression は 軽量 WebSocket 通信で CDP Runtime.evaluate を実行します
func EvaluateCDPExpression(wsURL string, expr string, timeout time.Duration) (string, error) {
	// ws://localhost:9222/devtools/page/...
	hostPath := strings.TrimPrefix(wsURL, "ws://")
	parts := strings.SplitN(hostPath, "/", 2)
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid ws url: %s", wsURL)
	}
	host, path := parts[0], "/"+parts[1]

	conn, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(timeout))

	keyBytes := make([]byte, 16)
	_, _ = rand.Read(keyBytes)
	key := base64.StdEncoding.EncodeToString(keyBytes)

	handshake := fmt.Sprintf("GET %s HTTP/1.1\r\nHost: %s\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Key: %s\r\nSec-WebSocket-Version: 13\r\n\r\n", path, host, key)
	if _, err := conn.Write([]byte(handshake)); err != nil {
		return "", err
	}

	buf := make([]byte, 2048)
	n, err := conn.Read(buf)
	if err != nil || !strings.Contains(string(buf[:n]), "101") {
		return "", fmt.Errorf("websocket handshake failed: %v", err)
	}

	payload, _ := json.Marshal(map[string]interface{}{
		"id": 1, "method": "Runtime.evaluate",
		"params": map[string]interface{}{"expression": expr, "returnByValue": true},
	})

	frame := buildWSFrame(payload)
	if _, err := conn.Write(frame); err != nil {
		return "", err
	}

	resBuf := make([]byte, 65536)
	rn, err := conn.Read(resBuf)
	if err != nil {
		return "", err
	}

	return extractWSPayload(resBuf[:rn]), nil
}

func buildWSFrame(payload []byte) []byte {
	l := len(payload)
	frame := []byte{0x81}
	if l <= 125 {
		frame = append(frame, byte(0x80|l))
	} else {
		frame = append(frame, 0x80|126, byte(l>>8), byte(l&0xFF))
	}
	mask := make([]byte, 4)
	_, _ = rand.Read(mask)
	frame = append(frame, mask...)
	for i := 0; i < l; i++ {
		frame = append(frame, payload[i]^mask[i%4])
	}
	return frame
}

func extractWSPayload(data []byte) string {
	if len(data) < 2 {
		return ""
	}
	headerLen := 2
	b2 := data[1]
	pLen := int(b2 & 0x7F)
	if pLen == 126 && len(data) >= 4 {
		pLen = int(data[2])<<8 | int(data[3])
		headerLen = 4
	}
	if len(data) >= headerLen+pLen {
		return string(data[headerLen : headerLen+pLen])
	}
	return string(data[headerLen:])
}
