package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
	"unicode"

	_ "github.com/glebarez/go-sqlite"
)

// Language detection
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

// Convert UNKNOWN_ ID to MD5 hash GUID-like string
func formatNumericID(id string) string {
	if strings.HasPrefix(id, "UNKNOWN_") {
		hash := md5.Sum([]byte(id))
		return fmt.Sprintf("%x-%x-%x-%x-%x", hash[0:4], hash[4:6], hash[6:8], hash[8:10], hash[10:16])
	}
	return id
}

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

func main() {
	srcDB, err := sql.Open("sqlite", "archive_org.db")
	if err != nil {
		log.Fatal("Failed to open src DB:", err)
	}
	defer srcDB.Close()

	destDB, err := sql.Open("sqlite", "archive.db")
	if err != nil {
		log.Fatal("Failed to open dest DB:", err)
	}
	defer destDB.Close()

	tx, err := destDB.Begin()
	if err != nil {
		log.Fatal("Failed to begin dest TX:", err)
	}
	defer tx.Rollback()

	// Clear destination tables for clean migration
	tables := []string{"account_profile_histories", "media", "url_redirects", "articles", "accounts", "whitelists"}
	for _, t := range tables {
		if _, err := tx.Exec("DELETE FROM " + t); err != nil {
			log.Fatalf("Failed to clean table %s: %v", t, err)
		}
	}

	// 1. Migrate Accounts and Profile History
	accRows, err := srcDB.Query("SELECT numeric_id, username, avatar_local_path FROM accounts")
	if err != nil {
		log.Fatal("Query accounts failed:", err)
	}
	defer accRows.Close()

	accountStmt, err := tx.Prepare("INSERT INTO accounts (numeric_id, username, display_name, avatar_url, updated_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer accountStmt.Close()

	histStmt, err := tx.Prepare("INSERT INTO account_profile_histories (account_id, display_name, avatar_original_url, avatar_seq, avatar_virtual_key, observed_at) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer histStmt.Close()

	for accRows.Next() {
		var numericID, username string
		var avatarLocalPath sql.NullString
		if err := accRows.Scan(&numericID, &username, &avatarLocalPath); err != nil {
			log.Fatal(err)
		}

		newID := formatNumericID(numericID)

		// Fetch profile histories for this account in chronological order
		hRows, err := srcDB.Query("SELECT display_name, observed_at FROM account_profile_history WHERE numeric_id = ? ORDER BY id ASC", numericID)
		if err != nil {
			log.Fatal(err)
		}

		seq := 0
		latestDisplayName := username
		for hRows.Next() {
			seq++
			var dispName sql.NullString
			var obsAt sql.NullString
			if err := hRows.Scan(&dispName, &obsAt); err != nil {
				log.Fatal(err)
			}
			dn := dispName.String
			if dn == "" {
				dn = username
			}
			latestDisplayName = dn

			virtualKey := fmt.Sprintf("%s_avatar_%03d", username, seq)

			var obs time.Time
			if obsAt.Valid && obsAt.String != "" {
				t, err := time.Parse("2006-01-02 15:04:05", obsAt.String)
				if err == nil {
					obs = t
				} else {
					obs = time.Now()
				}
			} else {
				obs = time.Now()
			}

			origAvatar := ""
			if avatarLocalPath.Valid {
				origAvatar = avatarLocalPath.String
			}

			if _, err := histStmt.Exec(newID, dn, origAvatar, seq, virtualKey, obs); err != nil {
				log.Fatal(err)
			}
		}
		hRows.Close()

		if seq == 0 {
			// Fallback if no history exists
			seq = 1
			virtualKey := fmt.Sprintf("%s_avatar_%03d", username, seq)
			origAvatar := ""
			if avatarLocalPath.Valid {
				origAvatar = avatarLocalPath.String
			}
			if _, err := histStmt.Exec(newID, latestDisplayName, origAvatar, seq, virtualKey, time.Now()); err != nil {
				log.Fatal(err)
			}
		}

		latestAvatarURL := fmt.Sprintf("%s_avatar_%03d", username, seq)
		if _, err := accountStmt.Exec(newID, username, latestDisplayName, latestAvatarURL, time.Now()); err != nil {
			log.Fatal(err)
		}
	}

	// 2. Migrate Articles (Tweets)
	tweetRows, err := srcDB.Query("SELECT tweet_id, numeric_id, conversation_id, created_at, full_text, reply_to_tweet_id, is_retweet, wayback_url, is_liked FROM tweets")
	if err != nil {
		log.Fatal(err)
	}
	defer tweetRows.Close()

	articleStmt, err := tx.Prepare("INSERT INTO articles (id, account_id, conversation_id, reply_to_id, reply_to_handle, created_at, full_text, lang, full_text_ja, full_text_en, full_text_zh, via, is_repost, is_liked, wayback_url) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer articleStmt.Close()

	for tweetRows.Next() {
		var id, numID string
		var convID, fullText, replyTo, wayback sql.NullString
		var createdAt sql.NullString
		var isRetweet, isLiked sql.NullBool

		if err := tweetRows.Scan(&id, &numID, &convID, &createdAt, &fullText, &replyTo, &isRetweet, &wayback, &isLiked); err != nil {
			log.Fatal(err)
		}

		newNumID := formatNumericID(numID)
		txt := fullText.String
		lang := detectLanguage(txt)

		var cAt time.Time
		if createdAt.Valid && createdAt.String != "" {
			t, err := time.Parse("2006-01-02 15:04:05", createdAt.String)
			if err == nil {
				cAt = t
			} else {
				t2, err2 := time.Parse(time.RFC3339, createdAt.String)
				if err2 == nil {
					cAt = t2
				} else {
					cAt = time.Now()
				}
			}
		} else {
			cAt = time.Now()
		}

		isRepost := isRetweet.Bool || strings.HasPrefix(strings.TrimSpace(txt), "RT @")

		var fullJa, fullEn, fullZh sql.NullString
		if lang == "ja" {
			fullJa = sql.NullString{String: txt, Valid: true}
		} else if lang == "zh" {
			fullZh = sql.NullString{String: txt, Valid: true}
		} else {
			fullEn = sql.NullString{String: txt, Valid: true}
		}

		cID := id
		if convID.Valid && convID.String != "" {
			cID = convID.String
		}

		_, err = articleStmt.Exec(id, newNumID, cID, replyTo.String, "", cAt, txt, lang, fullJa, fullEn, fullZh, "Twitter for Web", isRepost, isLiked.Bool, wayback.String)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 3. Migrate Media
	mediaRows, err := srcDB.Query("SELECT media_id, tweet_id, type, download_url, download_status, stash_scene_id, stash_image_id, width, height FROM media")
	if err != nil {
		log.Fatal(err)
	}
	defer mediaRows.Close()

	mediaStmt, err := tx.Prepare("INSERT INTO media (media_id, article_id, type, download_url, width, height, download_status, failed_reason, stash_scene_id, stash_image_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer mediaStmt.Close()

	seenSceneIDs := make(map[string]bool)
	seenImageIDs := make(map[string]bool)

	for mediaRows.Next() {
		var mediaID, tweetID string
		var mType, url, status sql.NullString
		var sScene, sImage sql.NullString
		var w, h sql.NullInt64

		if err := mediaRows.Scan(&mediaID, &tweetID, &mType, &url, &status, &sScene, &sImage, &w, &h); err != nil {
			log.Fatal(err)
		}

		finalStatus := mapDownloadStatus(status.String)
		finalType := mapMediaType(mType.String)

		var finalSceneID interface{} = nil
		if sScene.Valid && sScene.String != "" && sScene.String != "0" {
			if !seenSceneIDs[sScene.String] {
				seenSceneIDs[sScene.String] = true
				finalSceneID = sScene.String
			} else {
				log.Printf("Warning: Duplicate stash_scene_id '%s' on media '%s' skipped", sScene.String, mediaID)
			}
		}

		var finalImageID interface{} = nil
		if sImage.Valid && sImage.String != "" && sImage.String != "0" {
			if !seenImageIDs[sImage.String] {
				seenImageIDs[sImage.String] = true
				finalImageID = sImage.String
			} else {
				log.Printf("Warning: Duplicate stash_image_id '%s' on media '%s' skipped", sImage.String, mediaID)
			}
		}

		_, err = mediaStmt.Exec(mediaID, tweetID, finalType, url.String, w.Int64, h.Int64, finalStatus, "", finalSceneID, finalImageID)
		if err != nil {
			log.Fatal("media insert failed:", err)
		}
	}

	// 4. Migrate url_redirects
	urlRows, err := srcDB.Query("SELECT short_url, expanded_url, tweet_id FROM url_redirects")
	if err != nil {
		log.Fatal(err)
	}
	defer urlRows.Close()

	urlStmt, err := tx.Prepare("INSERT INTO url_redirects (short_url, expanded_url, article_id) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer urlStmt.Close()

	for urlRows.Next() {
		var short, exp, tweetID sql.NullString
		if err := urlRows.Scan(&short, &exp, &tweetID); err != nil {
			log.Fatal(err)
		}
		_, err = urlStmt.Exec(short.String, exp.String, tweetID.String)
		if err != nil {
			log.Fatal(err)
		}
	}

	// 5. Migrate whitelist
	wlRows, err := srcDB.Query("SELECT type, value, is_active FROM whitelist")
	if err != nil {
		log.Fatal(err)
	}
	defer wlRows.Close()

	wlStmt, err := tx.Prepare("INSERT INTO whitelists (type, value, is_active) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer wlStmt.Close()

	for wlRows.Next() {
		var wType, wVal sql.NullString
		var isAct sql.NullBool
		if err := wlRows.Scan(&wType, &wVal, &isAct); err != nil {
			log.Fatal(err)
		}
		_, err = wlStmt.Exec(wType.String, wVal.String, isAct.Bool)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Clean migration completed successfully.")
}
