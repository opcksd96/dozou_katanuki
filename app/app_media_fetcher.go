// app/app_media_fetcher.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

type FetchResult struct {
	Success    bool
	StatusCode int
	Bytes      int64
	ErrorMsg   string
}

// tryDirectFetchDetailed は 詳細なHTTPレスポンス情報とログ付きで直接ダウンロードを試行します
func (a *App) tryDirectFetchDetailed(client *http.Client, url, destPath string) FetchResult {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return FetchResult{Success: false, StatusCode: 0, ErrorMsg: err.Error()}
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return FetchResult{Success: false, StatusCode: 0, ErrorMsg: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return FetchResult{Success: false, StatusCode: resp.StatusCode, ErrorMsg: fmt.Sprintf("HTTP %d", resp.StatusCode)}
	}

	out, err := os.Create(destPath)
	if err != nil {
		return FetchResult{Success: false, StatusCode: resp.StatusCode, ErrorMsg: fmt.Sprintf("ファイル作成失敗: %v", err)}
	}
	defer out.Close()

	n, err := io.Copy(out, resp.Body)
	if err != nil || n < 512 {
		_ = os.Remove(destPath)
		return FetchResult{Success: false, StatusCode: resp.StatusCode, Bytes: n, ErrorMsg: "データ不完全/空ファイル"}
	}

	return FetchResult{Success: true, StatusCode: 200, Bytes: n}
}
