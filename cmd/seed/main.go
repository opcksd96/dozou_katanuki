package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"dozou_katanuki/driver"
	"dozou_katanuki/models"
)

func main() {
	// プロジェクトルートの archive.db に接続
	db, err := driver.InitDB("archive.db")
	if err != nil {
		log.Fatalf("DB接続失敗: %v", err)
	}
	repo := driver.NewRepository(db)

	now := time.Now()

	// テスト用シードデータ一覧
	seeds := []*models.Article{
		{
			ID:             "art_10001",
			ConversationID: "conv_10001",
			AccountID:      "acc_mashu",
			FullText:       "先輩、おはようございます！本日もカルデアでのデータ収集任務を開始します。",
			FullTextJA:     sql.NullString{String: "先輩、おはようございます！本日もカルデアでのデータ収集任務を開始します。", Valid: true},
			FullTextEN:     sql.NullString{String: "Good morning, Senpai! Starting today's data collection mission.", Valid: true},
			CreatedAt:      now.Add(-2 * time.Hour),
			IsLiked:        true,
			WaybackURL:     "https://web.archive.org/web/test1",
			Account: models.Account{
				NumericID:   "acc_mashu",
				Username:    "mash_kyrielight",
				DisplayName: "マシュ・キリエライト",
				AvatarURL:   "mash_avatar_01",
			},
			Media: []models.Media{
				{
					MediaID:      "med_20001",
					Type:         "image",
					Width:        1200,
					Height:       800,
					StashImageID: sql.NullString{String: "img_shield_01", Valid: true},
				},
			},
		},
		{
			ID:             "art_10002",
			ConversationID: "conv_10002",
			AccountID:      "acc_mashu",
			FullText:       "Wails 4.0.0 の内部結合テストは極めて順調です。SQLite WALモードの同期を確認しました。",
			FullTextJA:     sql.NullString{String: "Wails 4.0.0 の内部結合テストは極めて順調です。", Valid: true},
			CreatedAt:      now.Add(-10 * time.Minute),
			IsLiked:        false,
			WaybackURL:     "https://web.archive.org/web/test2",
			Account: models.Account{
				NumericID:   "acc_mashu",
				Username:    "mash_kyrielight",
				DisplayName: "マシュ・キリエライト",
				AvatarURL:   "mash_avatar_02", // アバター更新による世代監査テスト
			},
		},
	}

	for _, art := range seeds {
		if err := repo.UpsertArticleTx(art); err != nil {
			log.Fatalf("シード投入エラー (%s): %v", art.ID, err)
		}
		fmt.Printf("✔ 記事投入完了: %s (%s)\n", art.ID, art.Account.DisplayName)
	}

	log.Println("すべてのシードデータ投入が完了しました！")
}