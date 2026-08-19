package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// migrateAccounts は旧DBからアカウント情報およびプロフィール変更履歴を移行します
func migrateAccounts(srcDB *sql.DB, tx *sql.Tx) error {
	accRows, err := srcDB.Query("SELECT numeric_id, username, avatar_local_path FROM accounts")
	if err != nil {
		return fmt.Errorf("query accounts failed: %w", err)
	}
	defer accRows.Close()

	accountStmt, err := tx.Prepare("INSERT INTO accounts (numeric_id, username, display_name, avatar_url, updated_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer accountStmt.Close()

	histStmt, err := tx.Prepare("INSERT INTO account_profile_histories (account_id, display_name, avatar_original_url, avatar_seq, avatar_virtual_key, observed_at) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer histStmt.Close()

	for accRows.Next() {
		var numericID, username string
		var avatarLocalPath sql.NullString
		if err := accRows.Scan(&numericID, &username, &avatarLocalPath); err != nil {
			return err
		}

		newID := formatNumericID(numericID)

		hRows, err := srcDB.Query("SELECT display_name, observed_at FROM account_profile_history WHERE numeric_id = ? ORDER BY id ASC", numericID)
		if err != nil {
			return err
		}

		seq := 0
		latestDisplayName := username
		for hRows.Next() {
			seq++
			var dispName, obsAt sql.NullString
			if err := hRows.Scan(&dispName, &obsAt); err != nil {
				return err
			}
			dn := dispName.String
			if dn == "" {
				dn = username
			}
			latestDisplayName = dn
			virtualKey := fmt.Sprintf("%s_avatar_%03d", username, seq)

			var obs time.Time
			if obsAt.Valid && obsAt.String != "" {
				if t, err := time.Parse("2006-01-02 15:04:05", obsAt.String); err == nil {
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
				return err
			}
		}
		hRows.Close()

		if seq == 0 {
			seq = 1
			virtualKey := fmt.Sprintf("%s_avatar_%03d", username, seq)
			origAvatar := ""
			if avatarLocalPath.Valid {
				origAvatar = avatarLocalPath.String
			}
			if _, err := histStmt.Exec(newID, latestDisplayName, origAvatar, seq, virtualKey, time.Now()); err != nil {
				return err
			}
		}

		latestAvatarURL := fmt.Sprintf("%s_avatar_%03d", username, seq)
		if _, err := accountStmt.Exec(newID, username, latestDisplayName, latestAvatarURL, time.Now()); err != nil {
			return err
		}
	}
	log.Println("[Importer] Accounts & Profile History migrated successfully.")
	return nil
}
