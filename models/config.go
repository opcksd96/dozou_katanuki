package models

// SystemConfig はシステム全般の動作環境設定を表します
type SystemConfig struct {
	Env              string `json:"env"`
	DefaultFramework string `json:"default_framework"`
	Language         string `json:"language"`
}

// NetworkConfig はネットワークおよびポート設定を表します
type NetworkConfig struct {
	FrontendPort        int      `json:"frontend_port,omitempty"`
	StashProxyPort      int      `json:"stash_proxy_port,omitempty"`
	MiddlewarePort      int      `json:"middleware_port,omitempty"`
	BackendPort         int      `json:"backend_port,omitempty"`
	StashPort           int      `json:"stash_port"`
	PublicBindAddress   string   `json:"public_bind_address,omitempty"`
	InternalBindAddress string   `json:"internal_bind_address,omitempty"`
}

// StorageConfig は各種データおよびメディア保存先設定を表します
type StorageConfig struct {
	DBPath             string `json:"db_path"`
	StashEnabled       bool   `json:"stash_enabled"`
	LocalMediaDir      string `json:"local_media_dir"`
	ThunderDownloadDir string `json:"thunder_download_dir,omitempty"` // 迅雷テンポラリダウンロードフォルダ
	StashDir           string `json:"stash_dir"`
	DumpsDir           string `json:"dumps_dir"`
}

// SchedulerConfig はバックグラウンド監視およびバックアップ間隔設定を表します
type SchedulerConfig struct {
	PollIntervalSec      int `json:"poll_interval_sec"`
	BackupIntervalHours  int `json:"backup_interval_hours"`
	MaxBackupGenerations int `json:"max_backup_generations"`
}

// BroadcastConfig はLANキャスト配信設定を表します
type BroadcastConfig struct {
	Enabled         bool     `json:"enabled"`
	AllowedNetworks []string `json:"allowed_networks"`
}

// AppearanceConfig はUI表示および多言語フォント・テーマ・文字倍率設定を表します
type AppearanceConfig struct {
	Theme        string  `json:"theme,omitempty"`      // "system", "dark", "light"
	FontScale    float64 `json:"font_scale,omitempty"` // 0.8 ~ 1.5 (default: 1.0)
	FontFamilyJa string  `json:"font_family_ja"`
	FontFamilyEn string  `json:"font_family_en"`
	FontFamilyZh string  `json:"font_family_zh"`
}

// TranslationConfig は翻訳APIの動作設定を表します
type TranslationConfig struct {
	Provider              string `json:"provider"`                 // "deepl", "google", "none"
	DeeplApiKey           string `json:"deepl_api_key"`            // DeepL API Key
	GoogleTranslateApiKey string `json:"google_translate_api_key"` // Google Translate API Key
}

// BroadcastStatus は現在のLANキャスト配信の稼働状態を表します
type BroadcastStatus struct {
	Enabled         bool     `json:"enabled"`
	Running         bool     `json:"running"`
	BindAddress     string   `json:"bind_address"`
	Port            int      `json:"port"`
	LocalIPs        []string `json:"local_ips"`
	DetectedSubnets []string `json:"detected_subnets"`
	AllowedNetworks []string `json:"allowed_networks"`
	CastURL         string   `json:"cast_url"`
}

// AppConfig は config.json 全体のルート設定モデル (SPEC-CONFIG-001) です
type AppConfig struct {
	System      SystemConfig      `json:"system"`
	Network     NetworkConfig     `json:"network"`
	Storage     StorageConfig     `json:"storage"`
	Scheduler   SchedulerConfig   `json:"scheduler"`
	Broadcast   BroadcastConfig   `json:"broadcast"`
	Appearance  AppearanceConfig  `json:"appearance"`
	Translation TranslationConfig `json:"translation"`
}


