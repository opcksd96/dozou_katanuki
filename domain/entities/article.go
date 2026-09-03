// domain/entities/article.go (100行以下 - SPEC-PRINCIPLE-001)
package entities

import "time"

// Article は純粋なドメインエンティティです
type Article struct {
	ID             string
	AccountID      string
	ConversationID string
	ReplyToID      *string
	ReplyToHandle  *string
	CreatedAt      time.Time
	FullText       string
	Lang           string
	FullTextJA     *string
	FullTextEN     *string
	FullTextZH     *string
	Via            string
	IsRepost       bool
	IsLiked        bool
	WaybackURL     string
	SourceDomain   *string
	OriginalURL    *string
	SotweURL       *string
	NitterURL      *string
	TwistalkerURL  *string
	SourceName     *string
	IsTrash        bool
	TrashedBy      *string
	TrashReason    *string
	TrashedAt      *time.Time

	// Navigation / Associations (インフラ実装による遅延評価などを想定)
	Account *Account
	Media   []Media
}

// MarkAsTrash は記事をゴミ箱状態へ遷移させます
func (a *Article) MarkAsTrash(reason, by string) {
	a.IsTrash = true
	a.TrashReason = &reason
	a.TrashedBy = &by
	now := time.Now()
	a.TrashedAt = &now
}
