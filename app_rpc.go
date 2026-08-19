package main

import (
	"encoding/json"
	"log"
	"os"

	"dozou_katanuki/middleware"
	"dozou_katanuki/models"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// GetAccounts は登録されている全アカウントのリストを供給する Wails バインドメソッドです
func (a *App) GetAccounts(platform string) ([]models.RenderAuthor, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	return a.timelineService.GetAccounts(platform)
}

// GetTimeline はフロントエンドへ RenderTree 配列を供給する Wails バインドメソッドです
func (a *App) GetTimeline(platform, accountID, filter string, limit, offset int) ([]models.RenderTree, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	res, err := a.timelineService.FetchTimeline(platform, accountID, filter, limit, offset)
	log.Printf("[Wails RPC] GetTimeline(platform=%s, accountID=%s, filter=%s) -> 取得件数: %d (err: %v)",
		platform, accountID, filter, len(res), err)
	return res, err
}

// GetArticleDetail は指定された個別記事およびスレッド会話ツリーを取得する Wails バインドメソッドです
func (a *App) GetArticleDetail(platform, id string) (*models.ArticleDetailResult, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	return a.timelineService.GetArticleDetail(platform, id)
}

// GetSystemLanguage は config.json (SPEC-CONFIG-001) からシステム言語設定を取得します
func (a *App) GetSystemLanguage() string {
	type Config struct {
		System struct {
			Language string `json:"language"`
		} `json:"system"`
	}

	data, err := os.ReadFile("config.json")
	if err != nil {
		return "ja" // デフォルト言語
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil || cfg.System.Language == "" {
		return "ja"
	}
	return cfg.System.Language
}

// StartSalvageJob は自動サルベージジョブをキューに追加する Wails バインドメソッドです
func (a *App) StartSalvageJob(platform, account string, limit int) (*models.JobProgress, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	return a.jobOrchestrator.EnqueueSalvage(platform, account, limit)
}

// StartManualImportJob は手動 WARC インポートジョブをキューに追加する Wails バインドメソッドです
func (a *App) StartManualImportJob(warcPath string, offline bool) (*models.JobProgress, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	return a.jobOrchestrator.EnqueueManualImport(warcPath, offline)
}

// GetJobStatus は指定されたジョブのステータスを取得する Wails バインドメソッドです
func (a *App) GetJobStatus(jobID string) (*models.JobProgress, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	st := a.jobOrchestrator.GetStatus(jobID)
	return st, nil
}

// GetActiveJob は現在実行中のジョブを取得する Wails バインドメソッドです
func (a *App) GetActiveJob() (*models.JobProgress, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	return a.jobOrchestrator.GetActiveJob(), nil
}

// ListJobs は全ジョブの履歴ステータス一覧を取得する Wails バインドメソッドです
func (a *App) ListJobs() ([]*models.JobProgress, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	return a.jobOrchestrator.ListJobs(), nil
}

// CancelJob はジョブの実行を中断する Wails バインドメソッドです
func (a *App) CancelJob(jobID string) error {
	if err := a.waitForReady(); err != nil {
		return err
	}
	return a.jobOrchestrator.CancelJob(jobID)
}

// GetConfig は config.json (SPEC-CONFIG-001) の全設定を取得する Wails バインドメソッドです
func (a *App) GetConfig() (*models.AppConfig, error) {
	data, err := os.ReadFile("config.json")
	if err != nil {
		return nil, err
	}
	var cfg models.AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// SaveConfig は config.json (SPEC-CONFIG-001) を更新・保存する Wails バインドメソッドです
func (a *App) SaveConfig(cfg *models.AppConfig) error {
	if cfg == nil {
		return nil
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	// 改行を付与して保存
	data = append(data, '\n')
	if err := os.WriteFile("config.json", data, 0644); err != nil {
		return err
	}
	log.Printf("[Wails RPC] SaveConfig: config.json updated successfully")
	return nil
}

// GetWhitelists は全ホワイトリスト項目を取得する Wails バインドメソッドです
func (a *App) GetWhitelists() ([]models.Whitelist, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	return a.repo.GetWhitelists()
}

// AddWhitelist はホワイトリスト項目を追加する Wails バインドメソッドです
func (a *App) AddWhitelist(itemType, value string) (*models.Whitelist, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	item, err := a.repo.AddWhitelist(itemType, value)
	if err != nil {
		log.Printf("[Wails RPC] AddWhitelist error: %v", err)
		return nil, err
	}
	log.Printf("[Wails RPC] AddWhitelist added: [%s] %s (id: %d)", item.Type, item.Value, item.ID)
	return item, nil
}

// UpdateWhitelist はホワイトリスト項目を更新する Wails バインドメソッドです
func (a *App) UpdateWhitelist(id uint, itemType, value string, isActive bool) error {
	if err := a.waitForReady(); err != nil {
		return err
	}
	err := a.repo.UpdateWhitelist(id, itemType, value, isActive)
	if err != nil {
		log.Printf("[Wails RPC] UpdateWhitelist error (id: %d): %v", id, err)
		return err
	}
	log.Printf("[Wails RPC] UpdateWhitelist updated: (id: %d) [%s] %s (active: %v)", id, itemType, value, isActive)
	return nil
}

// DeleteWhitelist はホワイトリスト項目を削除する Wails バインドメソッドです
func (a *App) DeleteWhitelist(id uint) error {
	if err := a.waitForReady(); err != nil {
		return err
	}
	err := a.repo.DeleteWhitelist(id)
	if err != nil {
		log.Printf("[Wails RPC] DeleteWhitelist error (id: %d): %v", id, err)
		return err
	}
	log.Printf("[Wails RPC] DeleteWhitelist deleted: (id: %d)", id)
	return nil
}

// ToggleWhitelist はホワイトリスト項目の有効/無効を切り替える Wails バインドメソッドです
func (a *App) ToggleWhitelist(id uint) error {
	if err := a.waitForReady(); err != nil {
		return err
	}
	err := a.repo.ToggleWhitelist(id)
	if err != nil {
		log.Printf("[Wails RPC] ToggleWhitelist error (id: %d): %v", id, err)
		return err
	}
	log.Printf("[Wails RPC] ToggleWhitelist toggled: (id: %d)", id)
	return nil
}

// SearchArticles は保存済み記事の検索および一覧取得を行う Wails バインドメソッドです
func (a *App) SearchArticles(query, accountID, filter string, limit, offset int) (*models.ArticleSearchResult, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	articles, total, err := a.repo.SearchArticles(query, accountID, filter, limit, offset)
	if err != nil {
		log.Printf("[Wails RPC] SearchArticles error: %v", err)
		return nil, err
	}

	items := make([]models.RenderTree, len(articles))
	for i, art := range articles {
		items[i] = middleware.ToRenderTree(art, "twitter")
	}

	log.Printf("[Wails RPC] SearchArticles(query=%s, accountID=%s, filter=%s) -> 取得: %d件 / 全%d件",
		query, accountID, filter, len(items), total)
	return &models.ArticleSearchResult{
		Items: items,
		Total: total,
	}, nil
}

// GetArticle は指定されたIDの記事詳細を取得する Wails バインドメソッドです
func (a *App) GetArticle(id string) (*models.RenderTree, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	art, err := a.repo.GetArticleByID(id)
	if err != nil {
		log.Printf("[Wails RPC] GetArticle error (id=%s): %v", id, err)
		return nil, err
	}
	renderTree := middleware.ToRenderTree(*art, "twitter")
	return &renderTree, nil
}

// UpdateArticleTranslations は記事の日本語・英語・中国語翻訳テキストを更新する Wails バインドメソッドです
func (a *App) UpdateArticleTranslations(id, ja, en, zh string) error {
	if err := a.waitForReady(); err != nil {
		return err
	}
	err := a.repo.UpdateArticleTranslations(id, ja, en, zh)
	if err != nil {
		log.Printf("[Wails RPC] UpdateArticleTranslations error (id=%s): %v", id, err)
		return err
	}
	log.Printf("[Wails RPC] UpdateArticleTranslations updated: (id=%s) JA:%d字 EN:%d字 ZH:%d字",
		id, len(ja), len(en), len(zh))
	return nil
}

// RetryMediaDownload は指定されたメディアのダウンロードステータスをリセットし再試行ジョブをキックする Wails バインドメソッドです
func (a *App) RetryMediaDownload(mediaID string) error {
	if err := a.waitForReady(); err != nil {
		return err
	}
	if a.repo.ResetMediaStatus(mediaID) != nil {
		return a.repo.ResetMediaStatus(mediaID)
	}
	_, err := a.jobOrchestrator.EnqueueMediaDownload("twitter", mediaID)
	if err != nil {
		log.Printf("[Wails RPC] RetryMediaDownload EnqueueMediaDownload error: %v", err)
		return err
	}
	log.Printf("[Wails RPC] RetryMediaDownload triggered successfully for media: %s", mediaID)
	runtime.EventsEmit(a.ctx, "media:retried", map[string]string{"media_id": mediaID})
	return nil
}

// TriggerBackup は即時オンラインバックアップ (VACUUM INTO) を実行する Wails バインドメソッドです
func (a *App) TriggerBackup() (string, error) {
	if err := a.waitForReady(); err != nil {
		return "", err
	}
	if a.scheduler == nil {
		return "", nil
	}
	return a.scheduler.TriggerBackup()
}

// TriggerPoll は即時第3段階ポーリング (Motrix ➔ Stash 自動回収) を実行する Wails バインドメソッドです
func (a *App) TriggerPoll() (*models.JobProgress, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	if a.scheduler == nil {
		return nil, nil
	}
	return a.scheduler.TriggerPoll()
}

// RunAudit は SQLite3 整合性監査 (PRAGMA) および孤立ファイルの検出・パージを実行する Wails バインドメソッドです
func (a *App) RunAudit(purgeFiles, purgeDB bool) (*models.AuditReport, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	if a.auditService == nil {
		return nil, nil
	}

	stashDir := "./stash"
	blobsDir := "./blobs"
	cfg, err := a.GetConfig()
	if err == nil && cfg != nil {
		if cfg.Storage.StashDir != "" {
			stashDir = cfg.Storage.StashDir
		}
		if cfg.Storage.LocalMediaDir != "" {
			blobsDir = cfg.Storage.LocalMediaDir
		}
	}

	return a.auditService.RunAudit(a.ctx, stashDir, blobsDir, purgeFiles, purgeDB)
}

// PurgeOrphanFiles は指定された孤立ファイルをOSのごみ箱へ退避する Wails バインドメソッドです
func (a *App) PurgeOrphanFiles(paths []string) (int, error) {
	if err := a.waitForReady(); err != nil {
		return 0, err
	}
	if a.auditService == nil {
		return 0, nil
	}
	return a.auditService.PurgeOrphanFiles(paths)
}

// PurgeOrphanDBMedia は指定された media_id の DB レコードを削除する Wails バインドメソッドです
func (a *App) PurgeOrphanDBMedia(mediaIDs []string) (int, error) {
	if err := a.waitForReady(); err != nil {
		return 0, err
	}
	if a.auditService == nil {
		return 0, nil
	}
	return a.auditService.PurgeOrphanDBMedia(mediaIDs)
}

// TriggerRestore は Layer 2 (dumps/) からの完全オフライン自動リストア (SPEC-RECOVERY-001) をキックする Wails バインドメソッドです
func (a *App) TriggerRestore(dumpsDir string, resetDB bool) (*models.JobProgress, error) {
	if err := a.waitForReady(); err != nil {
		return nil, err
	}
	if dumpsDir == "" {
		dumpsDir = "./backups/dumps"
		cfg, err := a.GetConfig()
		if err == nil && cfg != nil && cfg.Storage.DumpsDir != "" {
			dumpsDir = cfg.Storage.DumpsDir
		}
	}

	if resetDB && a.repo != nil {
		log.Println("[Wails RPC] TriggerRestore: Resetting database before restore...")
		if err := a.repo.ResetDatabase(); err != nil {
			log.Printf("[Wails RPC] TriggerRestore: Database reset warning: %v", err)
		}
	}

	log.Printf("[Wails RPC] TriggerRestore: Starting restore job from %s (resetDB: %v)", dumpsDir, resetDB)
	return a.jobOrchestrator.EnqueueRestore(dumpsDir)
}

