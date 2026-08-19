package services

type RenderContent struct {
	Original string `json:"original"`
	JA       string `json:"ja"`
	EN       string `json:"en"`
	ZH       string `json:"zh"`
}

type RenderAuthor struct {
	NumericID   string `json:"numeric_id"`
	Handle      string `json:"handle"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	Bio         string `json:"bio"`
}

type RenderMediaURLs struct {
	Stream string `json:"stream"`
	Image  string `json:"image"`
}

type RenderMediaItem struct {
	MediaID      string          `json:"media_id"`
	Type         string          `json:"type"`
	URLs         RenderMediaURLs `json:"urls"`
	Width        int             `json:"width"`
	Height       int             `json:"height"`
	FailedReason string          `json:"failed_reason,omitempty"`
}

type RenderMetrics struct {
	Replies   int `json:"replies"`
	Retweets  int `json:"retweets"`
	Likes     int `json:"likes"`
	Bookmarks int `json:"bookmarks"`
	Quotes    int `json:"quotes"`
}

type RenderTree struct {
	ID             string            `json:"id"`
	ConversationID string            `json:"conversation_id"`
	CreatedAt      string            `json:"created_at"`
	Content        RenderContent     `json:"content"`
	Author         RenderAuthor      `json:"author"`
	Media          []RenderMediaItem `json:"media"`
	Metrics        RenderMetrics     `json:"metrics"`
	IsLiked        bool              `json:"is_liked"`
	SourceURL      string            `json:"source_url"`
}
