package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/glebarez/go-sqlite"
)

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

	// 移行先のテーブルをクリーンアップ
	tables := []string{"account_profile_histories", "media", "url_redirects", "articles", "accounts", "whitelists"}
	for _, t := range tables {
		if _, err := tx.Exec("DELETE FROM " + t); err != nil {
			log.Fatalf("Failed to clean table %s: %v", t, err)
		}
	}

	// 各マイグレーションを実行
	if err := migrateAccounts(srcDB, tx); err != nil {
		log.Fatalf("Accounts migration failed: %v", err)
	}
	if err := migrateTweets(srcDB, tx); err != nil {
		log.Fatalf("Tweets migration failed: %v", err)
	}
	if err := migrateMedia(srcDB, tx); err != nil {
		log.Fatalf("Media migration failed: %v", err)
	}
	if err := migrateAuxiliary(srcDB, tx); err != nil {
		log.Fatalf("Auxiliary migration failed: %v", err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatal("Transaction commit failed:", err)
	}

	fmt.Println("Clean migration completed successfully.")
}
