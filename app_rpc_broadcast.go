// app_rpc_broadcast.go (100行以下)
package main

import (
	"dozou_katanuki/middleware"
	"dozou_katanuki/models"
)

// GetBroadcastStatus は LAN Broadcast ＆ キャスト配信の稼働状態を取得する Wails バインドメソッドです
func (a *App) GetBroadcastStatus() (*models.BroadcastStatus, error) {
	if a.broadcastService == nil {
		cfg, _ := a.GetConfig()
		netCfg := models.NetworkConfig{MiddlewarePort: 5175, PublicBindAddress: "0.0.0.0"}
		bcastCfg := models.BroadcastConfig{}
		if cfg != nil {
			netCfg = cfg.Network
			bcastCfg = cfg.Broadcast
		}
		return &models.BroadcastStatus{
			Enabled:         bcastCfg.Enabled,
			Running:         false,
			BindAddress:     netCfg.PublicBindAddress,
			Port:            netCfg.MiddlewarePort,
			LocalIPs:        middleware.GetLocalIPv4s(),
			AllowedNetworks: bcastCfg.AllowedNetworks,
		}, nil
	}
	return a.broadcastService.GetStatus(), nil
}

// ToggleBroadcast は LAN Broadcast 配信の有効/無効を即座に切り替える Wails バインドメソッドです
func (a *App) ToggleBroadcast(enabled bool) (*models.BroadcastStatus, error) {
	cfg, err := a.GetConfig()
	if err != nil {
		return nil, err
	}
	cfg.Broadcast.Enabled = enabled
	if err := a.SaveConfig(cfg); err != nil {
		return nil, err
	}
	return a.GetBroadcastStatus()
}
