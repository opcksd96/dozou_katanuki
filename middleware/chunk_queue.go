// middleware/chunk_queue.go (100行以下)
package middleware

import (
	"dozou_katanuki/driver"
	"dozou_katanuki/models"
)

const maxChunkSize = 50

// FetchInChunks は 50件チャンク分割・反復取得アルゴリズムを実行します
func FetchInChunks(repo *driver.Repository, accountID, filter string, targetLimit, baseOffset int) ([]models.Article, error) {
	var accumulated []models.Article
	currentOffset := baseOffset
	remaining := targetLimit

	for remaining > 0 {
		chunkSize := remaining
		if chunkSize > maxChunkSize {
			chunkSize = maxChunkSize
		}

		articles, err := repo.FetchArticles(accountID, filter, chunkSize, currentOffset)
		if err != nil {
			return nil, err
		}

		accumulated = append(accumulated, articles...)
		fetchedCount := len(articles)

		// DB枯渇 (EOF) または 要求チャンク未満で終了
		if fetchedCount < chunkSize {
			break
		}

		remaining -= fetchedCount
		currentOffset += fetchedCount
	}

	return accumulated, nil
}
