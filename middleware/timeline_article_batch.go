// middleware/timeline_article_batch.go (Under 100 lines - SPEC-PRINCIPLE-001)
package middleware

import "dozou_katanuki/models"

func (s *TimelineService) BatchTrashArticles(ids []string, trashedBy, reason string) error {
	return s.repo.BatchTrashArticles(ids, trashedBy, reason)
}

func (s *TimelineService) BatchRestoreArticles(ids []string) error {
	return s.repo.BatchRestoreArticles(ids)
}

func (s *TimelineService) BatchResetTranslations(ids []string) error {
	return s.repo.BatchResetTranslations(ids)
}

func (s *TimelineService) GetArticlesByIDs(ids []string) ([]models.RenderTree, error) {
	articles, err := s.repo.GetArticlesByIDs(ids)
	if err != nil { return nil, err }
	items := make([]models.RenderTree, len(articles))
	for i, art := range articles { items[i] = ToRenderTree(art, "twitter") }
	return items, nil
}
