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
		if data, err := os.ReadFile(filepath.Join(appdata, "com.motrix.next", "system.json")); err == nil {
			var sys map[string]interface{}
			if json.Unmarshal(data, &sys) == nil {
				port := fmt.Sprintf("%v", sys["rpc-listen-port"])
				sec := fmt.Sprintf("%v", sys["rpc-secret"])
				if port != "" && port != "<nil>" {
					if sec == "<nil>" {
						sec = ""
					}
					cfgs = append(cfgs, motrixConfig{Endpoint: "http://127.0.0.1:" + port + "/jsonrpc", Secret: sec})
				}
			}
		}
	}
	cfgs = append(cfgs,
		motrixConfig{Endpoint: "http://127.0.0.1:29100/jsonrpc", Secret: "YyCGmFuwCnHvF5Bi"},
		motrixConfig{Endpoint: "http://127.0.0.1:29100/jsonrpc", Secret: ""},
		motrixConfig{Endpoint: "http://127.0.0.1:16800/jsonrpc", Secret: ""},
		motrixConfig{Endpoint: "http://127.0.0.1:6800/jsonrpc", Secret: ""},
	)
	return cfgs
}

func callMotrixRPC(method string, params []interface{}) ([]byte, error) {
	client := &http.Client{Timeout: 3 * time.Second}
	var lastErr error
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
		if err != nil {
			lastErr = err
			continue
		}
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == 200 {
			var rpcResp struct {
				Result interface{} `json:"result"`
				Error  interface{} `json:"error"`
			}
			if json.Unmarshal(buf.Bytes(), &rpcResp) == nil {
				if rpcResp.Error != nil {
					lastErr = fmt.Errorf("rpc error: %v", rpcResp.Error)
					continue
				}
				if rpcResp.Result != nil {
					return buf.Bytes(), nil
				}
			}
		}
		lastErr = fmt.Errorf("status %d: %s", resp.StatusCode, buf.String())
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return nil, fmt.Errorf("motrix offline")
}
