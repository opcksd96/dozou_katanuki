package relation

import "time"

// EvidenceType defines why this relationship was formed
type EvidenceType string

const (
	EvidenceModelSubject   EvidenceType = "model_subject"
	EvidenceBotMention     EvidenceType = "bot_mention"
	EvidenceQuoteTweet     EvidenceType = "quote_tweet"
	EvidenceMediaDuplicate EvidenceType = "media_duplicate"
)

// ArticleRelationEvidence represents a specific tweet/post that proves the relation
type ArticleRelationEvidence struct {
	ID              string       `json:"id"`
	RelationID      string       `json:"relation_id"`
	SourceArticleID string       `json:"source_article_id"`
	TargetArticleID *string      `json:"target_article_id,omitempty"` // pointer for NULL
	Type            EvidenceType `json:"evidence_type"`
	ContextSnippet  string       `json:"context_snippet"`
	MediaID         *string      `json:"media_id,omitempty"` // pointer for NULL
	IsSalvaged      bool         `json:"is_salvaged"`
	ObservedAt      time.Time    `json:"observed_at"`
}
