// models/media.go (100行以下)
package models

import "database/sql"

// Media represents media table (Stashapp Mapping and 3-Stage Recovery status)
type Media struct {
	MediaID        string         `gorm:"primaryKey;column:media_id;type:text" json:"media_id"`
	ArticleID      string         `gorm:"index;column:article_id;type:text;not null" json:"article_id"`
	Type           string         `gorm:"column:type;type:text;not null" json:"type"`
	DownloadURL    string         `gorm:"column:download_url;type:text;not null" json:"download_url"`
	Width          int            `gorm:"column:width;type:integer;not null" json:"width"`
	Height         int            `gorm:"column:height;type:integer;not null" json:"height"`
	DownloadStatus string         `gorm:"column:download_status;type:text;not null;default:'QUEUED'" json:"download_status"`
	FailedReason   sql.NullString `gorm:"column:failed_reason;type:text" json:"failed_reason"`
	StashSceneID   sql.NullString `gorm:"index;column:stash_scene_id;type:text" json:"stash_scene_id"`
	StashImageID   sql.NullString `gorm:"index;column:stash_image_id;type:text" json:"stash_image_id"`
}

// BuildRenderMedia converts a DB Media record to frontend RenderMedia with normalized URLs and status
func BuildRenderMedia(m Media) RenderMedia {
	var mediaURLs RenderMediaURLs
	mediaURLs.Original = m.DownloadURL
	effectiveStatus := m.DownloadStatus
	if m.DownloadStatus == "COMPLETED" && (m.StashSceneID.Valid || m.StashImageID.Valid) {
		if m.Type == "video" || m.Type == "gif" {
			mediaURLs.Stream = "/stash-proxy/scene/" + m.StashSceneID.String + "/stream"
			mediaURLs.Thumbnail = "/stash-proxy/scene/" + m.StashSceneID.String + "/screenshot"
			mediaURLs.Preview = "/stash-proxy/scene/" + m.StashSceneID.String + "/preview"
			mediaURLs.VTT = "/stash-proxy/scene/" + m.StashSceneID.String + "/vtt"
		} else {
			mediaURLs.Image = "/stash-proxy/image/" + m.StashImageID.String + "/image"
			mediaURLs.Thumbnail = "/stash-proxy/image/" + m.StashImageID.String + "/thumbnail"
		}
	} else if effectiveStatus == "COMPLETED" {
		effectiveStatus = "DEAD_404"
	}
	return RenderMedia{
		ID: m.MediaID, Type: m.Type, DownloadStatus: effectiveStatus,
		FailedReason: m.FailedReason.String, URLs: mediaURLs, Width: m.Width, Height: m.Height,
		StashSceneID: m.StashSceneID.String, StashImageID: m.StashImageID.String,
	}
}

// MapMediaToRenderMedia converts slice of Media to slice of RenderMedia
func MapMediaToRenderMedia(mediaList []Media) []RenderMedia {
	result := make([]RenderMedia, 0, len(mediaList))
	for _, m := range mediaList {
		result = append(result, BuildRenderMedia(m))
	}
	return result
}
