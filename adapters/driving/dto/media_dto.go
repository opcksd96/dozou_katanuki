// adapters/driving/dto/media_dto.go (100行以下 - SPEC-PRINCIPLE-001)
package dto

import "time"

// RenderMediaURLs はメディアのURLリンク集の描画用モデルです
type RenderMediaURLs struct {
	Stream    string `json:"stream,omitempty"`
	Image     string `json:"image,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
	Preview   string `json:"preview,omitempty"`
	VTT       string `json:"vtt,omitempty"`
	Original  string `json:"original,omitempty"`
}

// RenderMediaVariant はM3U8などの複数解像度バリアント情報です
type RenderMediaVariant struct {
	VariantHash string `json:"variant_hash"`
	DownloadURL string `json:"download_url"`
	BitRate     int    `json:"bit_rate"`
}

// RenderMedia はメディア単体の描画用モデルです
type RenderMedia struct {
	ID             string               `json:"id"`
	Type           string               `json:"type"`
	DownloadStatus string               `json:"download_status"`
	FailedReason   string               `json:"failed_reason,omitempty"`
	URLs           RenderMediaURLs      `json:"urls"`
	Width          int                  `json:"width,omitempty"`
	Height         int                  `json:"height,omitempty"`
	StashSceneID   string               `json:"stash_scene_id,omitempty"`
	StashImageID   string               `json:"stash_image_id,omitempty"`
	IsBookmarked   bool                 `json:"is_bookmarked"`
	DownloadURL    string               `json:"download_url,omitempty"`
	FilePath       string               `json:"file_path,omitempty"`
	MediaQuality   string               `json:"media_quality,omitempty"`
	IsTrash        bool                 `json:"is_trash,omitempty"`
	TrashedBy      string               `json:"trashed_by,omitempty"`
	TrashReason    string               `json:"trash_reason,omitempty"`
	Variants       []RenderMediaVariant `json:"variants,omitempty"`
}

// MediaItemDetail はメディア管理ビュー用のレコードモデルです
type MediaItemDetail struct {
	RenderMedia
	MediaID     string    `json:"media_id"`
	ArticleID   string    `json:"article_id"`
	AccountID   string    `json:"account_id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"display_name"`
	AvatarURL   string    `json:"avatar_url"`
	RawStatus   string    `json:"raw_status"`
	HasStash    bool      `json:"has_stash"`
	CreatedAt   time.Time `json:"created_at"`
	Title       string    `json:"title"`
	FullText    string    `json:"full_text"`
	FullTextJA  string    `json:"full_text_ja,omitempty"`
	TweetDate   string    `json:"tweet_date"`
	WaybackURL  string    `json:"wayback_url,omitempty"`
}

// MediaSearchStats はメディア種別の統計カウントです
type MediaSearchStats struct {
	TotalCount int64 `json:"total_count"`
	ImageCount int64 `json:"image_count"`
	VideoCount int64 `json:"video_count"`
}

// MediaSearchResult はメディア一覧のページネーション結果です
type MediaSearchResult struct {
	Items []MediaItemDetail `json:"items"`
	Total int64             `json:"total"`
	Stats MediaSearchStats  `json:"stats"`
}

// DownloadStatusStats はダウンロードキュー全体のステータス別集計です
type DownloadStatusStats struct {
	Queued     int64 `json:"queued"`
	Completed  int64 `json:"completed"`
	Dead404    int64 `json:"dead_404"`
	Outsourced int64 `json:"outsourced"`
	Escalated  int64 `json:"escalated"`
	Retained   int64 `json:"retained"`
	Failed     int64 `json:"failed"`
	Total      int64 `json:"total"`
}
