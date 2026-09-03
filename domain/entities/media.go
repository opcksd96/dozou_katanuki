// domain/entities/media.go (100行以下 - SPEC-PRINCIPLE-001)
package entities

import "time"

// MediaVariant は解像度などのバリアント情報を持つバリューオブジェクトです
type MediaVariant struct {
	VariantHash string
	MediaID     string
	ArticleID   string
	DownloadURL string
	BitRate     int
}

// Media は純粋なメディアドメインエンティティです
type Media struct {
	MediaID        string
	ArticleID      string
	AccountID      string
	Type           string
	DownloadURL    string
	Width          int
	Height         int
	DownloadStatus string // "QUEUED", "COMPLETED", "FAILED" etc.
	FailedReason   *string
	StashSceneID   *string
	StashImageID   *string
	IsBookmarked   bool
	MediaQuality   string
	IsTrash        bool
	TrashedBy      string
	TrashReason    string
	TrashedAt      *time.Time
	Variants       []MediaVariant
}

// IsCompleted はダウンロードが完了しているか判定します
func (m *Media) IsCompleted() bool {
	return m.DownloadStatus == "COMPLETED"
}

// MarkAsTrash はメディアをゴミ箱状態へ遷移させます
func (m *Media) MarkAsTrash(reason, by string) {
	m.IsTrash = true
	m.TrashReason = reason
	m.TrashedBy = by
	now := time.Now()
	m.TrashedAt = &now
}
