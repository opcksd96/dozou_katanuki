package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

// migrateTweets は旧ツイートデータを articles テーブルへ移行します
func migrateTweets(srcDB *sql.DB, tx *sql.Tx) error {
	tweetRows, err := srcDB.Query("SELECT tweet_id, numeric_id, conversation_id, created_at, full_text, reply_to_tweet_id, is_retweet, wayback_url, is_liked FROM tweets")
	if err != nil {
		return fmt.Errorf("query tweets failed: %w", err)
	}
	defer tweetRows.Close()

	stmt, err := tx.Prepare("INSERT INTO articles (id, account_id, conversation_id, reply_to_id, reply_to_handle, created_at, full_text, lang, full_text_ja, full_text_en, full_text_zh, via, is_repost, is_liked, wayback_url) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for tweetRows.Next() {
		var id, numID string
		var convID, fullText, replyTo, wayback sql.NullString
		var createdAt sql.NullString
		var isRetweet, isLiked sql.NullBool

		if err := tweetRows.Scan(&id, &numID, &convID, &createdAt, &fullText, &replyTo, &isRetweet, &wayback, &isLiked); err != nil {
			return err
		}

		newNumID := formatNumericID(numID)
		txt := fullText.String
		lang := detectLanguage(txt)

		var cAt time.Time
		if createdAt.Valid && createdAt.String != "" {
			if t, err := time.Parse("2006-01-02 15:04:05", createdAt.String); err == nil {
				cAt = t
			} else if t2, err2 := time.Parse(time.RFC3339, createdAt.String); err2 == nil {
				cAt = t2
			} else {
				cAt = time.Now()
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

		if _, err := stmt.Exec(id, newNumID, cID, replyTo.String, "", cAt, txt, lang, fullJa, fullEn, fullZh, "Twitter for Web", isRepost, isLiked.Bool, wayback.String); err != nil {
			return err
		}
	}
	log.Println("[Importer] Articles (Tweets) migrated successfully.")
	return nil
}
