# scripts/db_test.ps1
# SQLite CLI またはワンライナーで確認
go run -e "
package main
import ('fmt'; 'dozou_katanuki/driver'; 'dozou_katanuki/models')
func main() {
    db, _ := driver.InitDB(\"archive.db\")
    var histories []models.AccountProfileHistory
    db.Find(&histories)
    for _, h := range histories {
        fmt.Printf(\"AccountID: %s, Key: %s, URL: %s\n\", h.AccountID, h.AvatarVirtualKey, h.AvatarOriginalURL)
    }
}
"
