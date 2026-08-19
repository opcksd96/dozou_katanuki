package main

import (
	"crypto/md5"
	"fmt"
	"strings"
	"unicode"
)

// detectLanguage はテキストから使用言語（ja, zh, en）を判定します
func detectLanguage(text string) string {
	hasJa := false
	hasZh := false
	for _, r := range text {
		if unicode.In(r, unicode.Hiragana, unicode.Katakana) {
			hasJa = true
			break
		}
		if unicode.In(r, unicode.Han) {
			hasZh = true
		}
	}
	if hasJa {
		return "ja"
	}
	if hasZh {
		return "zh"
	}
	return "en"
}

// formatNumericID は UNKNOWN_ で始まるIDをMD5ハッシュGUID形式に正規化します
func formatNumericID(id string) string {
	if strings.HasPrefix(id, "UNKNOWN_") {
		hash := md5.Sum([]byte(id))
		return fmt.Sprintf("%x-%x-%x-%x-%x", hash[0:4], hash[4:6], hash[6:8], hash[8:10], hash[10:16])
	}
	return id
}

// mapDownloadStatus はダウンロード状態文字列を正規化します
func mapDownloadStatus(status string) string {
	switch strings.ToUpper(status) {
	case "SUCCESS", "COMPLETED":
		return "COMPLETED"
	case "DEAD_404":
		return "DEAD_404"
	case "OUTSOURCED":
		return "OUTSOURCED"
	case "RETAINED":
		return "RETAINED"
	default:
		return "QUEUED"
	}
}

// mapMediaType はメディア種別を正規化します
func mapMediaType(mType string) string {
	switch strings.ToLower(mType) {
	case "photo", "image":
		return "image"
	case "video":
		return "video"
	case "animated_gif", "gif":
		return "gif"
	default:
		return "image"
	}
}
