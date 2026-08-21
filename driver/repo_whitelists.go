// driver/repo_whitelists.go (100行以下)
package driver

import (
	"dozou_katanuki/models"
)

// GetWhitelists は登録されているすべてのホワイトリストを取得します
func (r *Repository) GetWhitelists() ([]models.Whitelist, error) {
	var list []models.Whitelist
	err := r.db.Order("id ASC").Find(&list).Error
	return list, err
}

// AddWhitelist はホワイトリスト項目を追加します
func (r *Repository) AddWhitelist(itemType, value string) (*models.Whitelist, error) {
	item := models.Whitelist{Type: itemType, Value: value, IsActive: true}
	if err := r.db.Create(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// UpdateWhitelist はホワイトリスト項目の内容を更新します
func (r *Repository) UpdateWhitelist(id uint, itemType, value string, isActive bool) error {
	return r.db.Model(&models.Whitelist{}).Where("id = ?", id).Updates(map[string]interface{}{
		"type": itemType, "value": value, "is_active": isActive,
	}).Error
}

// DeleteWhitelist はホワイトリスト項目を削除します
func (r *Repository) DeleteWhitelist(id uint) error {
	return r.db.Delete(&models.Whitelist{}, id).Error
}

// ToggleWhitelist はホワイトリスト項目の有効/無効を切り替えます
func (r *Repository) ToggleWhitelist(id uint) error {
	var item models.Whitelist
	if err := r.db.First(&item, id).Error; err != nil {
		return err
	}
	item.IsActive = !item.IsActive
	return r.db.Save(&item).Error
}
