package relation

import "time"

// RelationType defines the nature of the relationship
type RelationType string

const (
	RelTypePhotographer RelationType = "photographer"
	RelTypeBotSave      RelationType = "bot_save"
	RelTypeReupload     RelationType = "reupload"
	RelTypeReply        RelationType = "reply"
	RelTypeQuote        RelationType = "quote"
	RelTypeAlias        RelationType = "alias"
)

// Direction defines how the edge connects accounts
type Direction string

const (
	DirOutbound      Direction = "outbound"
	DirInbound       Direction = "inbound"
	DirBidirectional Direction = "bidirectional"
)

// AccountRelation represents a directed edge between two accounts
type AccountRelation struct {
	ID              string       `json:"id"`
	SourceAccountID string       `json:"source_account_id"`
	TargetAccountID string       `json:"target_account_id"`
	TargetHandle    string       `json:"target_handle"`
	Type            RelationType `json:"relation_type"`
	Direction       Direction    `json:"direction"`
	Weight          float64      `json:"weight"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}
