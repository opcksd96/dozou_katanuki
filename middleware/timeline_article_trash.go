// middleware/timeline_article_trash.go (Under 100 lines - SPEC-PRINCIPLE-001)
package middleware

func (s *TimelineService) TrashArticle(id, trashedBy, reason string) error {
	return s.repo.TrashArticle(id, trashedBy, reason)
}

func (s *TimelineService) RestoreArticle(id string) error {
	return s.repo.RestoreArticle(id)
}
