// app/app_rpc_whitelist_test.go (100行以下)
package app

import (
	"testing"

	"dozou_katanuki/driver"
)

func TestWhitelistCRUD(t *testing.T) {
	db, _ := setupTestDB(t)
	repo := driver.NewRepository(db)

	item, err := repo.AddWhitelist("account", "mashu_dev", "Team Chaldea", "@senpai")
	if err != nil || item.ID == 0 || item.Value != "mashu_dev" || item.GroupName != "Team Chaldea" || item.AliasOf != "@senpai" || !item.IsActive {
		t.Fatalf("AddWhitelist failed: %v", err)
	}

	list, err := repo.GetWhitelists()
	if err != nil || len(list) != 1 {
		t.Fatalf("GetWhitelists failed: %v", err)
	}

	if err := repo.ToggleWhitelist(item.ID); err != nil {
		t.Fatalf("ToggleWhitelist failed: %v", err)
	}
	list2, _ := repo.GetWhitelists()
	if list2[0].IsActive != false {
		t.Fatalf("Expected IsActive to be false")
	}

	if err := repo.UpdateWhitelist(item.ID, "keyword", "retro_famicom", "Retro Gaming", "", true); err != nil {
		t.Fatalf("UpdateWhitelist failed: %v", err)
	}
	list3, _ := repo.GetWhitelists()
	if list3[0].Type != "keyword" || list3[0].Value != "retro_famicom" || list3[0].GroupName != "Retro Gaming" {
		t.Fatalf("UpdateWhitelist failed: %+v", list3[0])
	}

	if err := repo.DeleteWhitelist(item.ID); err != nil {
		t.Fatalf("DeleteWhitelist failed: %v", err)
	}
	list4, _ := repo.GetWhitelists()
	if len(list4) != 0 {
		t.Fatalf("Expected len 0, got %d", len(list4))
	}
}
