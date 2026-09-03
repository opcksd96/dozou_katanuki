// application/timeline/usecase.go (100行以下 - SPEC-PRINCIPLE-001)
package timeline

import (
	"context"
	"dozou_katanuki/adapters/driving/dto"
	"dozou_katanuki/domain/ports"
)

type TimelineUseCase interface {
	GetThread(ctx context.Context, conversationID string) ([]dto.RenderTree, error)
}

type timelineUseCaseImpl struct {
	articleRepo ports.ArticleRepository
	mediaRepo   ports.MediaRepository
}

func NewTimelineUseCase(ar ports.ArticleRepository, mr ports.MediaRepository) TimelineUseCase {
	return &timelineUseCaseImpl{articleRepo: ar, mediaRepo: mr}
}

func (u *timelineUseCaseImpl) GetThread(ctx context.Context, conversationID string) ([]dto.RenderTree, error) {
	articles, err := u.articleRepo.GetThread(ctx, conversationID)
	if err != nil {
		return nil, err
	}
	
	trees := make([]dto.RenderTree, 0, len(articles))
	for _, art := range articles {
		medias, _ := u.mediaRepo.GetByArticleID(ctx, art.ID)
		
		var renderMedias []dto.RenderMedia
		for _, m := range medias {
			var variants []dto.RenderMediaVariant
			for _, v := range m.Variants {
				variants = append(variants, dto.RenderMediaVariant{
					VariantHash: v.VariantHash, DownloadURL: v.DownloadURL, BitRate: v.BitRate,
				})
			}
			renderMedias = append(renderMedias, dto.RenderMedia{
				ID: m.MediaID, Type: m.Type, DownloadStatus: m.DownloadStatus,
				Width: m.Width, Height: m.Height, IsBookmarked: m.IsBookmarked,
				MediaQuality: m.MediaQuality, IsTrash: m.IsTrash, Variants: variants,
			})
		}
		
		rt := dto.RenderTree{
			ID: art.ID, ConversationID: art.ConversationID,
			CreatedAt: art.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			Content: dto.RenderContent{Original: art.FullText},
			Author: dto.RenderAuthor{NumericID: art.AccountID}, // Minimal mapping for now
			IsLiked: art.IsLiked, WaybackURL: art.WaybackURL, IsTrash: art.IsTrash,
			Media: renderMedias,
		}
		
		trees = append(trees, rt)
	}
	
	return trees, nil
}
