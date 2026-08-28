// app/app_thunder_urls.go (100行以下 - SPEC-PRINCIPLE-001)
package app

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// EncodeThunderURL は生URLを迅雷専用の thunder:// スキームへBase64エンコードします
func EncodeThunderURL(rawURL string) string {
	if rawURL == "" {
		return ""
	}
	return "thunder://" + base64.StdEncoding.EncodeToString([]byte("AA"+rawURL+"ZZ"))
}

// BuildThunderFallbackURLs はあらゆる解像度と Wayback プレフィックスを展開します
func BuildThunderFallbackURLs(rawURL string) []string {
	if rawURL == "" {
		return nil
	}
	base := rawURL
	for _, sfx := range []string{":orig", ":large", ":medium", ":small", ":thumb", ":tiny"} {
		if strings.HasSuffix(base, sfx) {
			base = strings.TrimSuffix(base, sfx)
			break
		}
	}
	cleanBase := base
	if idx := strings.Index(base, "?"); idx != -1 {
		cleanBase = base[:idx]
	}

	var list []string
	seen := make(map[string]bool)
	add := func(u string) {
		if u != "" && !seen[u] {
			seen[u] = true
			list = append(list, u)
		}
	}

	if strings.Contains(base, "video") || strings.Contains(base, ".mp4") {
		for _, tag := range []string{"?tag=14", "?tag=12", "?tag=10"} {
			add(cleanBase + tag)
			add(fmt.Sprintf("https://web.archive.org/web/2id_/%s%s", cleanBase, tag))
		}
		add(cleanBase)
		add(fmt.Sprintf("https://web.archive.org/web/2id_/%s", cleanBase))
		add(rawURL)
		return list
	}

	cleanNoExt := cleanBase
	for _, ext := range []string{".jpg", ".jpeg", ".png", ".webp"} {
		if strings.HasSuffix(strings.ToLower(cleanNoExt), ext) {
			cleanNoExt = cleanNoExt[:len(cleanNoExt)-len(ext)]
			break
		}
	}

	for _, q := range []string{"orig", "large", "medium", "small", "thumb"} {
		add(fmt.Sprintf("%s?format=jpg&name=%s", cleanNoExt, q))
		add(fmt.Sprintf("%s:%s", cleanBase, q))
		add(fmt.Sprintf("https://web.archive.org/web/2id_/%s?format=jpg&name=%s", cleanNoExt, q))
		add(fmt.Sprintf("https://web.archive.org/web/2id_/%s:%s", cleanBase, q))
	}
	add(cleanBase)
	add(fmt.Sprintf("https://web.archive.org/web/2id_/%s", cleanBase))
	add(rawURL)
	return list
}
