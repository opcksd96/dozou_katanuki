// middleware/broadcast_service.go (100行以下)
package middleware

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"dozou_katanuki/models"
)

type BroadcastService struct {
	netCfg          models.NetworkConfig
	bcastCfg        models.BroadcastConfig
	unifiedHandler  *UnifiedHandler
	timelineService *TimelineService
	emitter         EventEmitter
	server          *http.Server
	listener        net.Listener
	mu              sync.RWMutex
	running         bool
}

func NewBroadcastService(netCfg models.NetworkConfig, bcastCfg models.BroadcastConfig, handler *UnifiedHandler, timeline *TimelineService, emitter EventEmitter) *BroadcastService {
	if netCfg.MiddlewarePort <= 0 { netCfg.MiddlewarePort = 5175 }
	if netCfg.PublicBindAddress == "" { netCfg.PublicBindAddress = "0.0.0.0" }
	if len(bcastCfg.AllowedNetworks) == 0 {
		bcastCfg.AllowedNetworks = []string{"192.168.0.0/16", "10.0.0.0/8", "172.16.0.0/12", "127.0.0.1/32", "::1/128"}
	}
	return &BroadcastService{netCfg: netCfg, bcastCfg: bcastCfg, unifiedHandler: handler, timelineService: timeline, emitter: emitter}
}

func (s *BroadcastService) Start(ctx context.Context) error {
	s.mu.Lock(); defer s.mu.Unlock()
	if s.running { return nil }
	if !s.bcastCfg.Enabled {
		log.Println("[Broadcast] LAN Broadcast is currently disabled in config.")
		return nil
	}
	return s.startServerLocked()
}

func (s *BroadcastService) Stop() error {
	s.mu.Lock(); defer s.mu.Unlock()
	if !s.running || s.server == nil { return nil }
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second); defer cancel()
	err := s.server.Shutdown(ctx)
	s.running, s.server, s.listener = false, nil, nil
	log.Println("[Broadcast] LAN Broadcast Server stopped.")
	return err
}

func (s *BroadcastService) UpdateConfig(netCfg models.NetworkConfig, bcastCfg models.BroadcastConfig) error {
	s.mu.Lock(); defer s.mu.Unlock()
	s.netCfg, s.bcastCfg = netCfg, bcastCfg
	if s.running && s.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		_ = s.server.Shutdown(ctx); cancel()
		s.server, s.listener, s.running = nil, nil, false
	}
	if s.bcastCfg.Enabled { return s.startServerLocked() }
	return nil
}

func (s *BroadcastService) GetStatus() *models.BroadcastStatus {
	s.mu.RLock(); defer s.mu.RUnlock()
	localIPs := GetLocalIPv4s()
	castURL := ""
	if len(localIPs) > 0 && s.netCfg.MiddlewarePort > 0 {
		castURL = fmt.Sprintf("http://%s:%d", localIPs[0], s.netCfg.MiddlewarePort)
	}
	return &models.BroadcastStatus{
		Enabled: s.bcastCfg.Enabled, Running: s.running, BindAddress: s.netCfg.PublicBindAddress,
		Port: s.netCfg.MiddlewarePort, LocalIPs: localIPs, AllowedNetworks: s.bcastCfg.AllowedNetworks, CastURL: castURL,
	}
}
