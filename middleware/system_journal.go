// middleware/system_journal.go (100行以下 - SPEC-PRINCIPLE-001)
package middleware

import (
	"fmt"
	"sync"
	"time"

	"dozou_katanuki/models"
)

const maxJournalEntries = 200

// SystemJournalManager はインフラ常駐処理用のオンメモリ・リングバッファジャーナルを管理します
type SystemJournalManager struct {
	mu      sync.RWMutex
	entries []models.SystemJournalEntry
	seq     int64
}

var (
	globalJournal *SystemJournalManager
	journalOnce   sync.Once
)

// GetGlobalJournal はシングルトンのジャーナルマネージャーを返します
func GetGlobalJournal() *SystemJournalManager {
	journalOnce.Do(func() {
		globalJournal = &SystemJournalManager{
			entries: make([]models.SystemJournalEntry, 0, maxJournalEntries),
		}
	})
	return globalJournal
}

// Record は新しい構造化ジャーナルエントリをオンメモリに追加します
func (jm *SystemJournalManager) Record(component, level, event, message string, payload map[string]interface{}) {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	jm.seq++
	entry := models.SystemJournalEntry{
		ID:        fmt.Sprintf("JRN-%06d", jm.seq),
		Timestamp: time.Now(),
		Component: component,
		Level:     level,
		Event:     event,
		Message:   message,
		Payload:   payload,
	}

	if len(jm.entries) >= maxJournalEntries {
		jm.entries = jm.entries[1:]
	}
	jm.entries = append(jm.entries, entry)
}

// GetEntries は最新のジャーナルエントリ一覧（降順）を返します
func (jm *SystemJournalManager) GetEntries(limit int) []models.SystemJournalEntry {
	jm.mu.RLock()
	defer jm.mu.RUnlock()

	total := len(jm.entries)
	if limit <= 0 || limit > total {
		limit = total
	}

	result := make([]models.SystemJournalEntry, 0, limit)
	for i := total - 1; i >= total-limit; i-- {
		result = append(result, jm.entries[i])
	}
	return result
}
