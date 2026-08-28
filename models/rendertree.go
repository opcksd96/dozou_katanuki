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
	Username    string `json:"username,omitempty"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	Bio         string `json:"bio,omitempty"`
	GroupName   string `json:"group_name,omitempty"`
	AliasOf     string `json:"alias_of,omitempty"`
}

type RenderMediaURLs struct {
	Stream    string `json:"stream,omitempty"`
	Image     string `json:"image,omitempty"`
	Thumbnail string `json:"thumbnail,omitempty"`
	Preview   string `json:"preview,omitempty"`
	VTT       string `json:"vtt,omitempty"`
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
	IsBookmarked   bool            `json:"is_bookmarked"`
	DownloadURL    string          `json:"download_url,omitempty"`
	FilePath       string          `json:"file_path,omitempty"`
	MediaQuality   string          `json:"media_quality,omitempty"`
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
	SourceDomain   string        `json:"source_domain,omitempty"`
	OriginalURL    string        `json:"original_url,omitempty"`
	WaybackURL     string        `json:"wayback_url,omitempty"`
	SotweURL       string        `json:"sotwe_url,omitempty"`
	NitterURL      string        `json:"nitter_url,omitempty"`
	TwistalkerURL  string        `json:"twistalker_url,omitempty"`
	ParentID       string        `json:"parent_id,omitempty"`
	ReplyToHandle  string        `json:"reply_to_handle,omitempty"`
	IsTrash        bool          `json:"is_trash,omitempty"`
	TrashedBy      string        `json:"trashed_by,omitempty"`
	TrashReason    string        `json:"trash_reason,omitempty"`
}

// ArticleSearchResult は記事検索のページネーション付き結果構造体です
type ArticleSearchResult struct {
	Items []RenderTree `json:"items"`
	Total int64        `json:"total"`
}

// ArticleDetailResult は個別記事およびスレッド会話ツリーの返却構造体です
type ArticleDetailResult struct {
	Article RenderTree   `json:"article"`
	Thread  []RenderTree `json:"thread"`
}

// SkinPackage はプラットフォーム別プレゼンテーションスキンのアセット群です (SPEC-PLUGIN-001)
type SkinPackage struct {
	Platform   string `json:"platform"`
	LayoutYAML string `json:"layout_yaml"`
	DesignCSS  string `json:"design_css"`
	Controller string `json:"controller_js"`
}

