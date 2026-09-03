package sqlite

import (
	"context"
	"database/sql"
	"dozou_katanuki/domain/relation"
)

type RelationRepoImpl struct {
	db *sql.DB
}

func NewRelationRepoImpl(db *sql.DB) *RelationRepoImpl {
	return &RelationRepoImpl{db: db}
}

func (r *RelationRepoImpl) SaveRelation(ctx context.Context, rel *relation.AccountRelation) error {
	query := `INSERT INTO account_relations (id, source_account_id, target_account_id, target_handle, relation_type, direction, weight, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	          ON CONFLICT(id) DO UPDATE SET target_handle=excluded.target_handle, weight=excluded.weight, updated_at=excluded.updated_at`
	_, err := r.db.ExecContext(ctx, query, rel.ID, rel.SourceAccountID, rel.TargetAccountID, rel.TargetHandle, rel.Type, rel.Direction, rel.Weight, rel.CreatedAt, rel.UpdatedAt)
	return err
}

func (r *RelationRepoImpl) GetRelationByID(ctx context.Context, id string) (*relation.AccountRelation, error) {
	// TODO: Get relation by ID implementation
	return nil, nil
}

func (r *RelationRepoImpl) GetRelationsBySource(ctx context.Context, sourceID string) ([]*relation.AccountRelation, error) {
	// TODO: Get relations by source implementation
	return nil, nil
}

func (r *RelationRepoImpl) GetAllRelations(ctx context.Context) ([]*relation.AccountRelation, error) {
	query := `SELECT id, source_account_id, target_account_id, target_handle, relation_type, direction, weight, created_at, updated_at FROM account_relations`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rels []*relation.AccountRelation
	for rows.Next() {
		var rel relation.AccountRelation
		if err := rows.Scan(&rel.ID, &rel.SourceAccountID, &rel.TargetAccountID, &rel.TargetHandle, &rel.Type, &rel.Direction, &rel.Weight, &rel.CreatedAt, &rel.UpdatedAt); err != nil {
			return nil, err
		}
		rels = append(rels, &rel)
	}
	return rels, nil
}

func (r *RelationRepoImpl) SaveEvidence(ctx context.Context, ev *relation.ArticleRelationEvidence) error {
	// TODO: Save evidence implementation
	return nil
}

func (r *RelationRepoImpl) GetEvidencesByRelation(ctx context.Context, relationID string) ([]*relation.ArticleRelationEvidence, error) {
	// TODO: Get evidences by relation implementation
	return nil, nil
}
