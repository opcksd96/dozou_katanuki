// models/rendertree.go (100行以下)
package models

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

type RenderMedia struct {
	Type         string   `json:"type"`
	URLs         []string `json:"urls"`
	FailedReason string   `json:"failed_reason,omitempty"`
}

// RenderMetrics は投稿への反応統計データです
type RenderMetrics struct {
	Replies int `json:"replies"`
	Reposts int `json:"reposts"`
	Likes   int `json:"likes"`
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
	SourceURL      string        `json:"source_url"`
}
