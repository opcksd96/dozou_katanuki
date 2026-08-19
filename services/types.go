package services

// RenderContent represents localized and decorated text blocks
type RenderContent struct {
	Original string `json:"original"`
	JA       string `json:"ja"`
	EN       string `json:"en"`
	ZH       string `json:"zh"`
}

// RenderAuthor represents sanitized author profile presentation
type RenderAuthor struct {
	NumericID   string `json:"numeric_id"`
	Handle      string `json:"handle"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	Bio         string `json:"bio"`
}

// RenderMediaURLs represents media stream and direct image URLs
type RenderMediaURLs struct {
	Stream string `json:"stream"`
	Image  string `json:"image"`
}

// RenderMediaItem represents media resource presentation
type RenderMediaItem struct {
	MediaID      string          `json:"media_id"`
	Type         string          `json:"type"` // "image" | "video" | "gif"
	URLs         RenderMediaURLs `json:"urls"`
	Width        int             `json:"width"`
	Height       int             `json:"height"`
	FailedReason string          `json:"failed_reason,omitempty"`
}

// RenderTree represents timeline UI render unit (SSOT Presentation DTO)
type RenderTree struct {
	ID             string            `json:"id"`
	ConversationID string            `json:"conversation_id"`
	CreatedAt      string            `json:"created_at"`
	Content        RenderContent     `json:"content"`
	Author         RenderAuthor      `json:"author"`
	Media          []RenderMediaItem `json:"media"`
	IsLiked        bool              `json:"is_liked"`
	SourceURL      string            `json:"source_url"`
}