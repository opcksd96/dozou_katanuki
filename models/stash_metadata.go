// models/stash_metadata.go (100行以下)
package models

// StashFileDetails represents media file metadata in Stash
type StashFileDetails struct {
	Path       string  `json:"path"`
	Size       int64   `json:"size"`
	Duration   float64 `json:"duration"`
	VideoCodec string  `json:"video_codec"`
	AudioCodec string  `json:"audio_codec"`
	Width      int     `json:"width"`
	Height     int     `json:"height"`
	FrameRate  float64 `json:"frame_rate"`
	BitRate    int64   `json:"bit_rate"`
}

// StashTag represents a tag in Stash
type StashTag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// StashMetadataResult represents detailed scene/image metadata from Stash GraphQL
type StashMetadataResult struct {
	ID        string             `json:"id"`
	IsScene   bool               `json:"is_scene"`
	Title     string             `json:"title"`
	Details   string             `json:"details"`
	URL       string             `json:"url"`
	Date      string             `json:"date"`
	Rating100 int                `json:"rating100"`
	Files     []StashFileDetails `json:"files"`
	Tags      []StashTag         `json:"tags"`
	Studio    string             `json:"studio,omitempty"`
}
