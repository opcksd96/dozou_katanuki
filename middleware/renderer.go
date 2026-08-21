// middleware/renderer.go (SPEC-MIDDLEWARE-001-2 / 100行以下)
package middleware

import (
	"fmt"
	"html"
	"regexp"
	"strings"

	"dozou_katanuki/models"
)

var (
	reURL     = regexp.MustCompile(`(https?://[^\s<]+)`)
	reTag     = regexp.MustCompile(`(^|[^\w#&;])#([a-zA-Z0-9_\p{L}\p{N}]+)`)
	reMention = regexp.MustCompile(`(^|[^\w@&;])@([a-zA-Z0-9_]{1,30})`)
	reFiller  = regexp.MustCompile(`https?://t\.co/[a-zA-Z0-9_]+\s*$`)
)

func decorateText(text, platform string, redirects []models.UrlRedirect, hasMedia bool) string {
	if text == "" { return "" }
	res := text
	for _, r := range redirects {
		if r.ShortURL != "" && r.ExpandedURL != "" && strings.Contains(res, r.ShortURL) {
			res = strings.ReplaceAll(res, r.ShortURL, r.ExpandedURL)
		}
	}
	if hasMedia { res = reFiller.ReplaceAllString(res, "") }
	res = strings.TrimSpace(res)

	safe := html.EscapeString(res)
	safe = reURL.ReplaceAllString(safe, `<a href="$1" target="_blank" rel="noopener noreferrer" class="external-link">$1</a>`)
	safe = reTag.ReplaceAllString(safe, fmt.Sprintf(`$1<a href="/%s/search?q=$2" class="hashtag-link" data-tag="$2">#$2</a>`, platform))
	safe = reMention.ReplaceAllString(safe, fmt.Sprintf(`$1<a href="/%s/$2" class="mention-link" data-mention="$2">@$2</a>`, platform))
	return strings.ReplaceAll(safe, "\n", "<br/>")
}

func ToRenderTree(item models.Article, platform string) models.RenderTree {
	if platform == "" { platform = "twitter" }
	hasMedia := len(item.Media) > 0
	avatarURL := AuditAndResolveAvatar(platform, item.CreatedAt, item.Account.ProfileHistory)

	orig := decorateText(item.FullText, platform, item.UrlRedirects, hasMedia)
	ja := item.FullTextJA.String
	if ja == "" { ja = item.FullText }
	jaDecorated := decorateText(ja, platform, item.UrlRedirects, hasMedia)

	en := item.FullTextEN.String
	if en == "" { en = item.FullText }
	enDecorated := decorateText(en, platform, item.UrlRedirects, hasMedia)

	zh := item.FullTextZH.String
	if zh == "" { zh = item.FullText }
	zhDecorated := decorateText(zh, platform, item.UrlRedirects, hasMedia)

	return models.RenderTree{
		ID:             item.ID,
		ConversationID: item.ConversationID,
		CreatedAt:      item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Content: models.RenderContent{
			Original: orig,
			JA:       jaDecorated,
			EN:       enDecorated,
			ZH:       zhDecorated,
		},
		Author: models.RenderAuthor{
			NumericID:   item.Account.NumericID,
			Handle:      item.Account.Username,
			DisplayName: item.Account.DisplayName,
			AvatarURL:   avatarURL,
		},
		Media: mapMediaToRenderMedia(item.Media),
		Metrics: models.RenderMetrics{
			Replies: 0, Retweets: 0, Likes: 0,
		},
		IsLiked:       item.IsLiked,
		IsPinned:      false,
		SourceURL:     item.WaybackURL,
		ParentID:      item.ReplyToID.String,
		ReplyToHandle: item.ReplyToHandle.String,
	}
}

func mapMediaToRenderMedia(mediaList []models.Media) []models.RenderMedia {
	result := make([]models.RenderMedia, 0, len(mediaList))
	for _, m := range mediaList {
		var mediaURLs models.RenderMediaURLs
		mediaURLs.Original = m.DownloadURL
		effectiveStatus := m.DownloadStatus
		if m.DownloadStatus == "COMPLETED" && (m.StashSceneID.Valid || m.StashImageID.Valid) {
			if m.Type == "video" || m.Type == "gif" {
				mediaURLs.Stream = fmt.Sprintf("/stash-proxy/scene/%s/stream", m.StashSceneID.String)
			} else {
				mediaURLs.Image = fmt.Sprintf("/stash-proxy/image/%s/image", m.StashImageID.String)
				mediaURLs.Thumbnail = fmt.Sprintf("/stash-proxy/image/%s/thumbnail", m.StashImageID.String)
			}
		} else if effectiveStatus == "COMPLETED" {
			effectiveStatus = "DEAD_404"
		}
		result = append(result, models.RenderMedia{
			ID: m.MediaID, Type: m.Type, DownloadStatus: effectiveStatus,
			FailedReason: m.FailedReason.String, URLs: mediaURLs, Width: m.Width, Height: m.Height,
			StashSceneID: m.StashSceneID.String, StashImageID: m.StashImageID.String,
		})
	}
	return result
}
