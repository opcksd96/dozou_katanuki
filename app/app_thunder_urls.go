package app

import (
	"encoding/base64"
	"fmt"
)

// EncodeThunderURL は生URLを迅雷専用の thunder:// スキームへBase64エンコードします
func EncodeThunderURL(rawURL string) string {
	if rawURL == "" {
		return ""
	}
	return "thunder://" + base64.StdEncoding.EncodeToString([]byte("AA"+rawURL+"ZZ"))
}

// BuildThunderFallbackURLs はあらゆる解像度 (orig, large, medium, small, thumb) と Wayback プレフィックスを展開します
func BuildThunderFallbackURLs(rawURL string) []string {
	if rawURL == "" {
		return nil
	}
	base := rawURL
	for _, sfx := range []string{":orig", ":large", ":medium", ":small", ":thumb", ":tiny"} {
		if len(base) > len(sfx) && base[len(base)-len(sfx):] == sfx {
			base = base[:len(base)-len(sfx)]
			break
		}
	}
	qualities := []string{"orig", "large", "medium", "small", "thumb"}
	var list []string
	seen := make(map[string]bool)
	add := func(u string) {
		if u != "" && !seen[u] {
			seen[u] = true
			list = append(list, u)
		}
	}

	for _, q := range qualities {
		add(fmt.Sprintf("%s?format=jpg&name=%s", base, q))
		add(fmt.Sprintf("%s:%s", base, q))
		add(fmt.Sprintf("https://web.archive.org/web/2id_/%s?format=jpg&name=%s", base, q))
		add(fmt.Sprintf("https://web.archive.org/web/2id_/%s:%s", base, q))
	}
	add(base)
	add(fmt.Sprintf("https://web.archive.org/web/2id_/%s", base))
	add(rawURL)
	return list
}
