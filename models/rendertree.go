// models/rendertree.go (100行以下)
package models

type RenderContent struct {
	Original string `json:"original"`
	JA       string `json:"ja,omitempty"`
	EN       string `json:"en,omitempty"`
	ZH       string `json:"zh,omitempty"`
}

type RenderAuthor struct {
	NumericID   string `json:"numeric_id"`
	Handle      string `json:"handle"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	Bio         string `json:"bio,omitempty"`
}

type RenderMediaURLs struct {
	Stream    string `json:"stream,omitempty"`
	Image     string `json:"image,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
	Original  string `json:"original,omitempty"`
}

type RenderMedia struct {
	ID             string          `json:"id"`
	Type           string          `json:"type"`
	DownloadStatus string          `json:"download_status"`
	FailedReason   string          `json:"failed_reason,omitempty"`
	URLs           RenderMediaURLs `json:"urls"`
	Width          int             `json:"width,omitempty"`
	Height         int             `json:"height,omitempty"`
	StashSceneID   string          `json:"stash_scene_id,omitempty"`
	StashImageID   string          `json:"stash_image_id,omitempty"`
}

type RenderMetrics struct {
	Replies  int `json:"replies"`
	Retweets int `json:"retweets"`
	Likes    int `json:"likes"`
	Views    int `json:"views,omitempty"`
}

// RenderTree はフロントエンド描画用の完成データ構造体です
type RenderTree struct {
	ID             string        `json:"id"`
	ConversationID string        `json:"conversation_id"`
	CreatedAt      string        `json:"created_at"`
	Content        RenderContent `json:"content"`
	Author         RenderAuthor  `json:"author"`
	Media          []RenderMedia `json:"media"`
	Metrics        RenderMetrics `json:"metrics"`
	IsLiked        bool          `json:"is_liked"`
	IsPinned       bool          `json:"is_pinned"`
	SourceURL      string        `json:"source_url"`
	ParentID       string        `json:"parent_id,omitempty"`
}

// ArticleSearchResult は記事検索のページネーション付き結果構造体です
type ArticleSearchResult struct {
	Items []RenderTree `json:"items"`
	Total int64        `json:"total"`
}

