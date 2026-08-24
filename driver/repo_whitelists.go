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
func (r *Repository) AddWhitelist(itemType, value, groupName, aliasOf string) (*models.Whitelist, error) {
	item := models.Whitelist{Type: itemType, Value: value, GroupName: groupName, AliasOf: aliasOf, IsActive: true}
	if err := r.db.Create(&item).Error; err != nil {
		return nil, err
	}
	if itemType == "account" {
		_ = r.db.Model(&models.Account{}).Where("username = ?", value).Updates(map[string]interface{}{
			"group_name": groupName, "alias_of": aliasOf,
		})
	}
	return &item, nil
}

// UpdateWhitelist はホワイトリスト項目の内容を更新します
func (r *Repository) UpdateWhitelist(id uint, itemType, value, groupName, aliasOf string, isActive bool) error {
	err := r.db.Model(&models.Whitelist{}).Where("id = ?", id).Updates(map[string]interface{}{
		"type": itemType, "value": value, "group_name": groupName, "alias_of": aliasOf, "is_active": isActive,
	}).Error
	if err == nil && itemType == "account" {
		_ = r.db.Model(&models.Account{}).Where("username = ?", value).Updates(map[string]interface{}{
			"group_name": groupName, "alias_of": aliasOf,
		})
	}
	return err
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

