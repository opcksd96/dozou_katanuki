package dto

import "time"

// TimelineItemDTO represents a single post/article in the timeline.
// This is the common contract between the backend, WebUI, and Wails Desktop.
type TimelineItemDTO struct {
	ID          string            `json:"id"`
	AccountID   string            `json:"account_id"`
	AccountName string            `json:"account_name"`
	AvatarURL   string            `json:"avatar_url"`
	Content     string            `json:"content"`
	PostedAt    time.Time         `json:"posted_at"`
	MediaFiles  []MediaFileDTO    `json:"media_files"`
	SourceURL   string            `json:"source_url"`
	IsScraped   bool              `json:"is_scraped"`
}

// MediaFileDTO represents a media attachment inside a TimelineItemDTO.
type MediaFileDTO struct {
	ID           string `json:"id"`
	MediaType    string `json:"media_type"` // "image" or "video"
	ThumbnailURL string `json:"thumbnail_url"`
	OriginalURL  string `json:"original_url"`
	LocalPath    string `json:"local_path,omitempty"`
}

// TimelineResponseDTO wraps the timeline results and pagination metadata.
type TimelineResponseDTO struct {
	Items      []TimelineItemDTO `json:"items"`
	TotalCount int64             `json:"total_count"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
}
