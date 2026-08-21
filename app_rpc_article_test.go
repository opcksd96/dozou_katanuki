// app_rpc_article_test.go (100行以下)
package main

import (
	"database/sql"
	"testing"
	"time"

	"dozou_katanuki/driver"
	"dozou_katanuki/models"
)

func TestArticleSearchAndTranslation(t *testing.T) {
	db, _ := setupTestDB(t)
	repo := driver.NewRepository(db)

	acc := models.Account{NumericID: "1001", Username: "senpai_apu", DisplayName: "先輩マスター", UpdatedAt: time.Now()}
	db.Create(&acc)

	art := models.Article{
		ID: "art_001", AccountID: "1001", ConversationID: "art_001", CreatedAt: time.Now(),
		FullText: "NES sound register $4000 test", Lang: "ja",
		FullTextJA: sql.NullString{String: "ファミコン音源レジスタ $4000 テスト", Valid: true},
	}
	db.Create(&art)

	res, total, err := repo.SearchArticles("sound", "all", "all", 10, 0)
	if err != nil || total != 1 || len(res) != 1 { t.Fatalf("SearchArticles failed: %v", err) }

	res2, total2, err := repo.SearchArticles("ファミコン", "all", "all", 10, 0)
	if err != nil || total2 != 1 || len(res2) != 1 { t.Fatalf("SearchArticles by translation failed: %v", err) }

	err = repo.UpdateArticleTranslations("art_001", "ファミコン更新", "NES updated", "红白机更新")
	if err != nil { t.Fatalf("UpdateArticleTranslations failed: %v", err) }

	upArt, err := repo.GetArticleByID("art_001")
	if err != nil || upArt.FullTextJA.String != "ファミコン更新" { t.Fatalf("Translations not updated: %+v", upArt) }
}
