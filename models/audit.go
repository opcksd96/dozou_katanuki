package models

import "time"

// ForeignKeyViolation は PRAGMA foreign_key_check の違反行を表します
type ForeignKeyViolation struct {
	Table       string `json:"table"`
	RowID       int64  `json:"row_id"`
	ParentTable string `json:"parent_table"`
	FkID        int    `json:"fk_id"`
}

// OrphanDBMedia は実ファイルが存在しない、または親記事が存在しない DB メディアレコードを表します
type OrphanDBMedia struct {
	MediaID      string `json:"media_id"`
	ArticleID    string `json:"article_id"`
	Type         string `json:"type"`
	DownloadURL  string `json:"download_url"`
	Status       string `json:"status"`
	StashSceneID string `json:"stash_scene_id,omitempty"`
	StashImageID string `json:"stash_image_id,omitempty"`
	Reason       string `json:"reason"`
}

// OrphanFile は DB に登録されていないストレージ上の未紐付け孤立ファイルを表します
type OrphanFile struct {
	Path     string `json:"path"`
	FileName string `json:"file_name"`
	FileSize int64  `json:"file_size"`
	Category string `json:"category"` // "stash_scene", "stash_image", "blob", "other"
}

// AuditReport は SPEC-AUDIT-001 に基づくデータベースおよびストレージの総合監査結果を表します
type AuditReport struct {
	ExecutedAt         time.Time             `json:"executed_at"`
	IntegrityOK        bool                  `json:"integrity_ok"`
	IntegrityErrors    []string              `json:"integrity_errors"`
	ForeignKeyOK       bool                  `json:"foreign_key_ok"`
	ForeignKeyErrors   []ForeignKeyViolation `json:"foreign_key_errors"`
	OrphanDBMedia      []OrphanDBMedia       `json:"orphan_db_media"`
	OrphanFiles        []OrphanFile          `json:"orphan_files"`
	PurgedFileCount    int                   `json:"purged_file_count"`
	PurgedDBMediaCount int                   `json:"purged_db_media_count"`
	Summary            string                `json:"summary"`
}
