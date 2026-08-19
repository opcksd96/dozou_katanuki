// models/url_redirect.go (100行以下)
package models

type UrlRedirect struct {
	ShortURL    string `gorm:"primaryKey;column:short_url;type:text" json:"short_url"`
	ExpandedURL string `gorm:"column:expanded_url;type:text;not null" json:"expanded_url"`
	ArticleID   string `gorm:"index;column:article_id;type:text;not null" json:"article_id"`
}
