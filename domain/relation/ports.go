package relation

import "context"

// RelationRepository is the outbound port for persisting account relations and evidences
type RelationRepository interface {
	// Relations
	SaveRelation(ctx context.Context, rel *AccountRelation) error
	GetRelationByID(ctx context.Context, id string) (*AccountRelation, error)
	GetRelationsBySource(ctx context.Context, sourceID string) ([]*AccountRelation, error)
	GetAllRelations(ctx context.Context) ([]*AccountRelation, error)

	// Evidences
	SaveEvidence(ctx context.Context, ev *ArticleRelationEvidence) error
	GetEvidencesByRelation(ctx context.Context, relationID string) ([]*ArticleRelationEvidence, error)
}
