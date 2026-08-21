package main

import (
	"fmt"

	"dozou_katanuki/driver"
	"dozou_katanuki/models"
)

func main() {
	db, err := driver.InitDB("archive.db")
	if err != nil {
		fmt.Println("DB Init Error:", err)
		return
	}

	var accounts []models.Account
	db.Find(&accounts)
	fmt.Printf("=== アカウント数 (%d件) ===\n", len(accounts))
	for _, a := range accounts {
		var count int64
		db.Model(&models.Article{}).Where("account_id = ?", a.NumericID).Count(&count)
		fmt.Printf("Account: %s (@%s) -> 記事数: %d件\n", a.DisplayName, a.Username, count)
	}

	var whitelists []models.Whitelist
	db.Find(&whitelists)
	fmt.Printf("\n=== Whitelist (%d件) ===\n", len(whitelists))
	for _, w := range whitelists {
		fmt.Printf("ID: %d | Type: %s | Value: %s | Active: %v\n", w.ID, w.Type, w.Value, w.IsActive)
	}
}