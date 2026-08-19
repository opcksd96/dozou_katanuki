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

	var histories []models.AccountProfileHistory
	db.Find(&histories)

	fmt.Printf("=== アバター履歴 (%d件) ===\n", len(histories))
	for _, h := range histories {
		fmt.Printf("AccountID: %s | VirtualKey: %s | OriginalURL: %s\n", h.AccountID, h.AvatarVirtualKey, h.AvatarOriginalURL)
	}
}