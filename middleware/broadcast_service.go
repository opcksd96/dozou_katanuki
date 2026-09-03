// middleware/broadcast_service.go (100行以下)
package middleware

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"sync"

	"dozou_katanuki/adapters/driving/dto"
	"dozou_katanuki/models"
)

type BroadcastService struct {
	netCfg          models.NetworkConfig
	bcastCfg        models.BroadcastConfig
	unifiedHandler  *UnifiedHandler
	timelineService *TimelineService
	emitter         EventEmitter
	beaconCallback  func(req dto.BeaconRequestDTO)
	clients         map[chan []byte]bool
	distFS          fs.FS
	server          *http.Server
	listener        net.Listener
	useTLS          bool
	mu              sync.RWMutex
	running         bool
	adminUseCases   AdminUseCases
}

func NewBroadcastService(netCfg models.NetworkConfig, bcastCfg models.BroadcastConfig, handler *UnifiedHandler, timeline *TimelineService, emitter EventEmitter) *BroadcastService {
	if netCfg.MiddlewarePort <= 0 { netCfg.MiddlewarePort = 5175 }
	if netCfg.PublicBindAddress == "" { netCfg.PublicBindAddress = "0.0.0.0" }
	if len(bcastCfg.AllowedNetworks) == 0 {
		bcastCfg.AllowedNetworks = append(GetLocalSubnets(), "127.0.0.1/32", "::1/128")
	}
	return &BroadcastService{netCfg: netCfg, bcastCfg: bcastCfg, unifiedHandler: handler, timelineService: timeline, emitter: emitter}
}

func (s *BroadcastService) SetDistFS(dfs fs.FS) {
	s.mu.Lock(); defer s.mu.Unlock(); s.distFS = dfs
}

type AdminUseCases interface {
	TogglePipelineAutoEngine(enable bool) (bool, error)
	IsPipelineAutoEngineRunning() bool
	GetPipelineOverview() (interface{}, error) // Will use interface{} for now to avoid circular deps
	GetPipelineLogs(stage string, limit int) (interface{}, error)
	SyncThunderDownloads(req string) (interface{}, error)
	ResetAllToQueuedAndBootstrap() (interface{}, error)
	IgnitePipeline() (interface{}, error)
}

func (s *BroadcastService) SetBeaconCallback(cb func(req dto.BeaconRequestDTO)) {
	s.mu.Lock(); defer s.mu.Unlock(); s.beaconCallback = cb
}

func (s *BroadcastService) SetAdminUseCases(uc AdminUseCases) {
	s.mu.Lock(); defer s.mu.Unlock(); s.adminUseCases = uc
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
	if !s.running { return nil }
	s.running = false
	var err error
	if s.server != nil { err = s.server.Close(); s.server = nil }
	if s.listener != nil { _ = s.listener.Close(); s.listener = nil }
	log.Println("[Broadcast] LAN Broadcast Server stopped.")
	return err
}

func (s *BroadcastService) UpdateConfig(netCfg models.NetworkConfig, bcastCfg models.BroadcastConfig) error {
	s.mu.Lock(); defer s.mu.Unlock()
	s.netCfg, s.bcastCfg = netCfg, bcastCfg
	if s.running {
		if s.server != nil { _ = s.server.Close() }
		if s.listener != nil { _ = s.listener.Close() }
		s.server, s.listener, s.running = nil, nil, false
	}
	if s.bcastCfg.Enabled { return s.startServerLocked() }
	return nil
}

func (s *BroadcastService) GetStatus() *models.BroadcastStatus {
	s.mu.RLock(); defer s.mu.RUnlock()
	localIPs := GetLocalIPv4s()
	scheme := "http"
	if s.useTLS { scheme = "https" }
	castURL := ""
	if len(localIPs) > 0 && s.netCfg.MiddlewarePort > 0 {
		castURL = fmt.Sprintf("%s://%s:%d", scheme, localIPs[0], s.netCfg.MiddlewarePort)
	}
	return &models.BroadcastStatus{
		Enabled: s.bcastCfg.Enabled, Running: s.running, BindAddress: s.netCfg.PublicBindAddress,
		Port: s.netCfg.MiddlewarePort, LocalIPs: localIPs, DetectedSubnets: GetLocalSubnets(),
		AllowedNetworks: s.bcastCfg.AllowedNetworks, CastURL: castURL,
	}
}
