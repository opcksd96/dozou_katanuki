// app/app_rpc_motrix_client.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type motrixConfig struct {
	Endpoint string
	Secret   string
}

func getMotrixConfigs() []motrixConfig {
	var cfgs []motrixConfig
	appdata := os.Getenv("APPDATA")
	if appdata != "" {
		// 1. Motrix Next (system.json / config.json)
		if data, err := os.ReadFile(filepath.Join(appdata, "com.motrix.next", "system.json")); err == nil {
			var sys map[string]interface{}
			if json.Unmarshal(data, &sys) == nil {
				port := fmt.Sprintf("%v", sys["rpc-listen-port"])
				sec := fmt.Sprintf("%v", sys["rpc-secret"])
				if port != "" && port != "<nil>" {
					cfgs = append(cfgs, motrixConfig{Endpoint: "http://127.0.0.1:" + port + "/jsonrpc", Secret: sec})
				}
			}
		}
		// 2. Motrix (settings.json)
		if data, err := os.ReadFile(filepath.Join(appdata, "Motrix", "settings.json")); err == nil {
			var m struct {
				Engine struct {
					RpcPort   int    `json:"rpcPort"`
					RpcSecret string `json:"rpcSecret"`
				} `json:"engine"`
			}
			if json.Unmarshal(data, &m) == nil && m.Engine.RpcPort > 0 {
				cfgs = append(cfgs, motrixConfig{
					Endpoint: fmt.Sprintf("http://127.0.0.1:%d/jsonrpc", m.Engine.RpcPort),
					Secret:   m.Engine.RpcSecret,
				})
			}
		}
	}
	// デフォルトフォールバック (29100, 16800, 6800)
	cfgs = append(cfgs,
		motrixConfig{Endpoint: "http://127.0.0.1:29100/jsonrpc", Secret: ""},
		motrixConfig{Endpoint: "http://127.0.0.1:16800/jsonrpc", Secret: ""},
		motrixConfig{Endpoint: "http://127.0.0.1:6800/jsonrpc", Secret: ""},
	)
	return cfgs
}

func callMotrixRPC(method string, params []interface{}) ([]byte, error) {
	client := &http.Client{Timeout: 2 * time.Second}
	for _, cfg := range getMotrixConfigs() {
		var rpcParams []interface{}
		if cfg.Secret != "" && cfg.Secret != "<nil>" {
			rpcParams = append(rpcParams, "token:"+cfg.Secret)
		}
		if len(params) > 0 {
			rpcParams = append(rpcParams, params...)
		}
		payload, _ := json.Marshal(map[string]interface{}{
			"jsonrpc": "2.0", "id": "dozou", "method": method, "params": rpcParams,
		})
		resp, err := client.Post(cfg.Endpoint, "application/json", bytes.NewReader(payload))
		if err == nil && resp.StatusCode == 200 {
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			resp.Body.Close()
			return buf.Bytes(), nil
		}
		if resp != nil {
			resp.Body.Close()
		}
	}
	return nil, fmt.Errorf("motrix offline or unauthorized")
}
