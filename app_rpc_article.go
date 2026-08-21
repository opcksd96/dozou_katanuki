// app_rpc_article.go (100行以下)
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"dozou_katanuki/middleware"
	"dozou_katanuki/models"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)


// SearchArticles は保存済み記事の検索および一覧取得を行う Wails バインドメソッドです
func (a *App) SearchArticles(query, accountID, filter string, limit, offset int) (*models.ArticleSearchResult, error) {
	if err := a.waitForReady(); err != nil { return nil, err }
	articles, total, err := a.repo.SearchArticles(query, accountID, filter, limit, offset)
	if err != nil { log.Printf("[Wails RPC] SearchArticles error: %v", err); return nil, err }
	items := make([]models.RenderTree, len(articles))
	for i, art := range articles { items[i] = middleware.ToRenderTree(art, "twitter") }
	return &models.ArticleSearchResult{Items: items, Total: total}, nil
}

// GetArticle は指定されたIDの記事詳細を取得する Wails バインドメソッドです
func (a *App) GetArticle(id string) (*models.RenderTree, error) {
	if err := a.waitForReady(); err != nil { return nil, err }
	art, err := a.repo.GetArticleByID(id)
	if err != nil { log.Printf("[Wails RPC] GetArticle error (id=%s): %v", id, err); return nil, err }
	renderTree := middleware.ToRenderTree(*art, "twitter")
	return &renderTree, nil
}

// UpdateArticleTranslations は記事の日本語・英語・中国語翻訳テキストを更新する Wails バインドメソッドです
func (a *App) UpdateArticleTranslations(id, ja, en, zh string) error {
	if err := a.waitForReady(); err != nil { return err }
	if err := a.repo.UpdateArticleTranslations(id, ja, en, zh); err != nil {
		log.Printf("[Wails RPC] UpdateArticleTranslations error (id=%s): %v", id, err)
		return err
	}
	return nil
}

// AutoTranslateArticle は指定記事を翻訳し、DB保存を行わずに翻訳結果の下書き（RenderTree）を返します
func (a *App) AutoTranslateArticle(id string) (*models.RenderTree, error) {
	if err := a.waitForReady(); err != nil { return nil, err }
	art, err := a.GetArticle(id)
	if err != nil { return nil, err }
	cmd := exec.Command("python", "plugins/twitter/scraper/main.py", "-m", "translate", "--article-id", id, "--dry-run")
	if env := a.getTranslationEnv(); len(env) > 0 {
		cmd.Env = os.Environ()
		for k, v := range env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[Wails RPC] AutoTranslateArticle error: %v, out: %s", err, string(out))
		return nil, fmt.Errorf("translate failed: %w", err)
	}
	for _, l := range strings.Split(string(out), "\n") {
		l = strings.TrimSpace(l)
		if strings.HasPrefix(l, "JSON:") {
			var tRes map[string]string
			if err := json.Unmarshal([]byte(strings.TrimPrefix(l, "JSON:")), &tRes); err == nil {
				if v, ok := tRes["ja"]; ok && v != "" { art.Content.JA = v }
				if v, ok := tRes["en"]; ok && v != "" { art.Content.EN = v }
				if v, ok := tRes["zh"]; ok && v != "" { art.Content.ZH = v }
				return art, nil
			}
		}
	}
	return art, nil
}

// RetryMediaDownload は指定されたメディアのダウンロードステータスをリセットし再試行ジョブをキックします
func (a *App) RetryMediaDownload(mediaID string) error {
	if err := a.waitForReady(); err != nil { return err }
	if a.repo.ResetMediaStatus(mediaID) != nil { return a.repo.ResetMediaStatus(mediaID) }
	if _, err := a.jobOrchestrator.EnqueueMediaDownload("twitter", mediaID); err != nil {
		log.Printf("[Wails RPC] RetryMediaDownload error: %v", err)
		return err
	}
	runtime.EventsEmit(a.ctx, "media:retried", map[string]string{"media_id": mediaID})
	return nil
}
