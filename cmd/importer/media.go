package main

import (
	"database/sql"
	"fmt"
	"log"
)

// migrateMedia はメディアデータを media テーブルへ移行します
func migrateMedia(srcDB *sql.DB, tx *sql.Tx) error {
	mediaRows, err := srcDB.Query("SELECT media_id, tweet_id, type, download_url, download_status, stash_scene_id, stash_image_id, width, height FROM media")
	if err != nil {
		return fmt.Errorf("query media failed: %w", err)
	}
	defer mediaRows.Close()

	stmt, err := tx.Prepare("INSERT INTO media (media_id, article_id, type, download_url, width, height, download_status, failed_reason, stash_scene_id, stash_image_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	seenSceneIDs := make(map[string]bool)
	seenImageIDs := make(map[string]bool)

	for mediaRows.Next() {
		var mediaID, tweetID string
		var mType, url, status, sScene, sImage sql.NullString
		var w, h sql.NullInt64

		if err := mediaRows.Scan(&mediaID, &tweetID, &mType, &url, &status, &sScene, &sImage, &w, &h); err != nil {
			return err
		}

		finalStatus := mapDownloadStatus(status.String)
		finalType := mapMediaType(mType.String)

		var finalSceneID interface{} = nil
		if sScene.Valid && sScene.String != "" && sScene.String != "0" {
			if !seenSceneIDs[sScene.String] {
				seenSceneIDs[sScene.String] = true
				finalSceneID = sScene.String
			}
		}

		var finalImageID interface{} = nil
		if sImage.Valid && sImage.String != "" && sImage.String != "0" {
			if !seenImageIDs[sImage.String] {
				seenImageIDs[sImage.String] = true
				finalImageID = sImage.String
			}
		}

		if _, err := stmt.Exec(mediaID, tweetID, finalType, url.String, w.Int64, h.Int64, finalStatus, "", finalSceneID, finalImageID); err != nil {
			return fmt.Errorf("media insert failed: %w", err)
		}
	}
	log.Println("[Importer] Media migrated successfully.")
	return nil
}

// migrateAuxiliary は URLリダイレクトとホワイトリストを移行します
func migrateAuxiliary(srcDB *sql.DB, tx *sql.Tx) error {
	// 1. URL Redirects
	urlRows, err := srcDB.Query("SELECT short_url, expanded_url, tweet_id FROM url_redirects")
	if err != nil {
		return err
	}
	defer urlRows.Close()

	urlStmt, err := tx.Prepare("INSERT INTO url_redirects (short_url, expanded_url, article_id) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer urlStmt.Close()

	for urlRows.Next() {
		var short, exp, tweetID sql.NullString
		if err := urlRows.Scan(&short, &exp, &tweetID); err != nil {
			return err
		}
		if _, err := urlStmt.Exec(short.String, exp.String, tweetID.String); err != nil {
			return err
		}
	}

	// 2. Whitelist
	wlRows, err := srcDB.Query("SELECT type, value, is_active FROM whitelist")
	if err != nil {
		return err
	}
	defer wlRows.Close()

	wlStmt, err := tx.Prepare("INSERT INTO whitelists (type, value, is_active) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer wlStmt.Close()

	for wlRows.Next() {
		var wType, wVal sql.NullString
		var isAct sql.NullBool
		if err := wlRows.Scan(&wType, &wVal, &isAct); err != nil {
			return err
		}
		if _, err := wlStmt.Exec(wType.String, wVal.String, isAct.Bool); err != nil {
			return err
		}
	}
	log.Println("[Importer] Auxiliary tables migrated successfully.")
	return nil
}
