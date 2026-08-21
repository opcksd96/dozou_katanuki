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
	db.Order("account_id ASC, avatar_seq ASC").Find(&histories)
	fmt.Printf("=== 世代履歴 (%d件) ===\n", len(histories))
	for _, h := range histories {
		hasB64 := len(h.AvatarBase64) > 0
		b64Len := len(h.AvatarBase64)
		fmt.Printf("AccountID: %s | Seq: %d | Key: %-25s | Base64格納: %v (%d bytes)\n",
			h.AccountID, h.AvatarSeq, h.AvatarVirtualKey, hasB64, b64Len)
	}
}