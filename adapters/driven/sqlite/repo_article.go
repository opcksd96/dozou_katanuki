// adapters/driven/sqlite/repo_article.go (100行以下 - SPEC-PRINCIPLE-001)
package sqlite

import (
	"context"
	"dozou_katanuki/domain/entities"
	"dozou_katanuki/domain/ports"
	"gorm.io/gorm"
)

type ArticleRepositoryImpl struct {
	db *gorm.DB
}

var _ ports.ArticleRepository = (*ArticleRepositoryImpl)(nil)

func NewArticleRepositoryImpl(db *gorm.DB) *ArticleRepositoryImpl {
	return &ArticleRepositoryImpl{db: db}
}

func (r *ArticleRepositoryImpl) Save(ctx context.Context, article *entities.Article) error {
	schema := &ArticleSchema{}
	schema.FromEntity(article)
	return r.db.WithContext(ctx).Save(schema).Error
}

func (r *ArticleRepositoryImpl) GetByID(ctx context.Context, id string) (*entities.Article, error) {
	var schema ArticleSchema
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&schema).Error; err != nil {
		return nil, err
	}
	return schema.ToEntity(), nil
}

func (r *ArticleRepositoryImpl) GetThread(ctx context.Context, conversationID string) ([]*entities.Article, error) {
	var schemas []ArticleSchema
	if err := r.db.WithContext(ctx).Where("conversation_id = ?", conversationID).Order("created_at asc").Find(&schemas).Error; err != nil {
		return nil, err
	}
	
	result := make([]*entities.Article, len(schemas))
	for i, s := range schemas {
		schemaCopy := s
		result[i] = schemaCopy.ToEntity()
	}
	return result, nil
}
