// adapters/driven/sqlite/repo_media.go (100行以下 - SPEC-PRINCIPLE-001)
package sqlite

import (
	"context"
	"dozou_katanuki/domain/entities"
	"dozou_katanuki/domain/ports"
	"gorm.io/gorm"
)

type MediaRepositoryImpl struct {
	db *gorm.DB
}

var _ ports.MediaRepository = (*MediaRepositoryImpl)(nil)

func NewMediaRepositoryImpl(db *gorm.DB) *MediaRepositoryImpl {
	return &MediaRepositoryImpl{db: db}
}

func (r *MediaRepositoryImpl) Save(ctx context.Context, media *entities.Media) error {
	schema := &MediaSchema{}
	schema.FromEntity(media)
	return r.db.WithContext(ctx).Save(schema).Error
}

func (r *MediaRepositoryImpl) GetByID(ctx context.Context, mediaID string) (*entities.Media, error) {
	var schema MediaSchema
	if err := r.db.WithContext(ctx).Preload("Variants").Where("media_id = ?", mediaID).First(&schema).Error; err != nil {
		return nil, err
	}
	return schema.ToEntity(), nil
}

func (r *MediaRepositoryImpl) GetByArticleID(ctx context.Context, articleID string) ([]*entities.Media, error) {
	var schemas []MediaSchema
	if err := r.db.WithContext(ctx).Preload("Variants").Where("article_id = ?", articleID).Find(&schemas).Error; err != nil {
		return nil, err
	}
	
	result := make([]*entities.Media, len(schemas))
	for i, s := range schemas {
		schemaCopy := s
		result[i] = schemaCopy.ToEntity()
	}
	return result, nil
}

func (r *MediaRepositoryImpl) ListByStatus(ctx context.Context, status string) ([]*entities.Media, error) {
	var schemas []MediaSchema
	if err := r.db.WithContext(ctx).Preload("Variants").Where("download_status = ?", status).Find(&schemas).Error; err != nil {
		return nil, err
	}
	
	result := make([]*entities.Media, len(schemas))
	for i, s := range schemas {
		schemaCopy := s
		result[i] = schemaCopy.ToEntity()
	}
	return result, nil
}
