// app/app_rpc_stash.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (a *App) getStashGraphQLEndpoint() string {
	port := 9999
	if cfg, err := a.GetConfig(); err == nil && cfg != nil && cfg.Network.StashPort > 0 {
		port = cfg.Network.StashPort
	}
	return fmt.Sprintf("http://127.0.0.1:%d/graphql", port)
}

func (a *App) queryStashGraphQL(query string, variables map[string]interface{}) (map[string]interface{}, error) {
	endpoint := a.getStashGraphQLEndpoint()
	payload := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}
	body, err := json.Marshal(payload)
	if err != nil { return nil, err }

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(endpoint, "application/json", bytes.NewReader(body))
	if err != nil { return nil, fmt.Errorf("stash graphql connection error: %w", err) }
	defer resp.Body.Close()

	var result struct {
		Data   map[string]interface{} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil { return nil, err }
	if len(result.Errors) > 0 { return nil, fmt.Errorf("graphql error: %s", result.Errors[0].Message) }
	return result.Data, nil
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok { return v }
	return ""
}

func getFloat(m map[string]interface{}, key string) float64 {
	if v, ok := m[key].(float64); ok { return v }
	return 0
}

func (a *App) isStashServerOnline() bool {
	return a.IsStashReady()
}
