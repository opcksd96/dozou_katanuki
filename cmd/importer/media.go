// cmd/importer/media.go (100行以下)
package main

import (
	"database/sql"
	"fmt"
	"log"
)

func migrateMedia(srcDB *sql.DB, tx *sql.Tx) error {
	rows, err := srcDB.Query("SELECT media_id, tweet_id, type, download_url, download_status, stash_scene_id, stash_image_id, width, height FROM media")
	if err != nil { return fmt.Errorf("query media failed: %w", err) }
	defer rows.Close()

	stmt, err := tx.Prepare("INSERT INTO media (media_id, article_id, type, download_url, width, height, download_status, failed_reason, stash_scene_id, stash_image_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil { return err }
	defer stmt.Close()

	seenScenes, seenImages := make(map[string]bool), make(map[string]bool)
	for rows.Next() {
		var mediaID, tweetID string
		var mType, url, status, sScene, sImage sql.NullString
		var w, h sql.NullInt64
		if err := rows.Scan(&mediaID, &tweetID, &mType, &url, &status, &sScene, &sImage, &w, &h); err != nil { return err }

		var finalScene interface{} = nil
		if sScene.Valid && sScene.String != "" && sScene.String != "0" && !seenScenes[sScene.String] {
			seenScenes[sScene.String] = true; finalScene = sScene.String
		}
		var finalImage interface{} = nil
		if sImage.Valid && sImage.String != "" && sImage.String != "0" && !seenImages[sImage.String] {
			seenImages[sImage.String] = true; finalImage = sImage.String
		}
		if _, err := stmt.Exec(mediaID, tweetID, mapMediaType(mType.String), url.String, w.Int64, h.Int64, mapDownloadStatus(status.String), "", finalScene, finalImage); err != nil {
			return fmt.Errorf("media insert failed: %w", err)
		}
	}
	log.Println("[Importer] Media migrated successfully.")
	return nil
}

func migrateAuxiliary(srcDB *sql.DB, tx *sql.Tx) error {
	uRows, err := srcDB.Query("SELECT short_url, expanded_url, tweet_id FROM url_redirects")
	if err != nil { return err }
	defer uRows.Close()

	uStmt, err := tx.Prepare("INSERT INTO url_redirects (short_url, expanded_url, article_id) VALUES (?, ?, ?)")
	if err != nil { return err }
	defer uStmt.Close()
	for uRows.Next() {
		var short, exp, tweetID sql.NullString
		if err := uRows.Scan(&short, &exp, &tweetID); err != nil { return err }
		if _, err := uStmt.Exec(short.String, exp.String, tweetID.String); err != nil { return err }
	}

	wRows, err := srcDB.Query("SELECT type, value, is_active FROM whitelist")
	if err != nil { return err }
	defer wRows.Close()

	wStmt, err := tx.Prepare("INSERT INTO whitelists (type, value, is_active) VALUES (?, ?, ?)")
	if err != nil { return err }
	defer wStmt.Close()
	for wRows.Next() {
		var wType, wVal sql.NullString
		var isAct sql.NullBool
		if err := wRows.Scan(&wType, &wVal, &isAct); err != nil { return err }
		if _, err := wStmt.Exec(wType.String, wVal.String, isAct.Bool); err != nil { return err }
	}
	log.Println("[Importer] Auxiliary tables migrated successfully.")
	return nil
}
