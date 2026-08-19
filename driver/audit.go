package driver

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"dozou_katanuki/models"

	"gorm.io/gorm"
)

// RunIntegrityCheck は PRAGMA integrity_check を実行し、SQLite ページの破損を検査します
func RunIntegrityCheck(db *gorm.DB) (bool, []string, error) {
	rows, err := db.Raw("PRAGMA integrity_check;").Rows()
	if err != nil {
		return false, nil, fmt.Errorf("failed to execute integrity_check: %w", err)
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var line string
		if err := rows.Scan(&line); err != nil {
			return false, nil, fmt.Errorf("failed to scan integrity_check result: %w", err)
		}
		results = append(results, line)
	}

	if len(results) == 1 && strings.EqualFold(strings.TrimSpace(results[0]), "ok") {
		return true, results, nil
	}

	return false, results, nil
}

// RunForeignKeyCheck は PRAGMA foreign_key_check を実行し、外部キー制約違反を検出します
func RunForeignKeyCheck(db *gorm.DB) ([]models.ForeignKeyViolation, error) {
	rows, err := db.Raw("PRAGMA foreign_key_check;").Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to execute foreign_key_check: %w", err)
	}
	defer rows.Close()

	var violations []models.ForeignKeyViolation
	for rows.Next() {
		var v models.ForeignKeyViolation
		var fkid sql.NullInt32
		if err := rows.Scan(&v.Table, &v.RowID, &v.ParentTable, &fkid); err != nil {
			return nil, fmt.Errorf("failed to scan foreign_key_check result: %w", err)
		}
		if fkid.Valid {
			v.FkID = int(fkid.Int32)
		}
		violations = append(violations, v)
	}

	return violations, nil
}

// FindOrphanDBMedia は親記事が存在しないレコードや、完了状態なのに実ファイルが見当たらない孤立 DB メディアを検出します
func FindOrphanDBMedia(db *gorm.DB, stashDir, blobsDir string) ([]models.OrphanDBMedia, error) {
	var orphans []models.OrphanDBMedia

	// 1. 親記事 (articles) が存在しない media レコードの検出
	var missingParentMedia []models.Media
	err := db.Raw(`
		SELECT m.* FROM media m
		LEFT JOIN articles a ON m.article_id = a.id
		WHERE a.id IS NULL
	`).Scan(&missingParentMedia).Error
	if err != nil {
		return nil, fmt.Errorf("failed to query orphan media without parent article: %w", err)
	}

	for _, m := range missingParentMedia {
		orphans = append(orphans, models.OrphanDBMedia{
			MediaID:      m.MediaID,
			ArticleID:    m.ArticleID,
			Type:         m.Type,
			DownloadURL:  m.DownloadURL,
			Status:       m.DownloadStatus,
			StashSceneID: m.StashSceneID.String,
			StashImageID: m.StashImageID.String,
			Reason:       fmt.Sprintf("親記事 (%s) が存在しません", m.ArticleID),
		})
	}

	// 2. 実ファイル確認（COMPLETED なのに実ファイルがストレージ上に一切ない場合）
	var completedMedia []models.Media
	if err := db.Where("download_status = ?", "COMPLETED").Find(&completedMedia).Error; err == nil {
		for _, m := range completedMedia {
			hasFile := false

			// blobs ディレクトリ内の検索
			if blobsDir != "" {
				matches, _ := filepath.Glob(filepath.Join(blobsDir, m.MediaID+".*"))
				if len(matches) > 0 {
					hasFile = true
				}
			}

			// stash ディレクトリ内の検索
			if !hasFile && stashDir != "" {
				if m.StashSceneID.Valid && m.StashSceneID.String != "" {
					matches, _ := filepath.Glob(filepath.Join(stashDir, "scenes", "*"+m.StashSceneID.String+"*"))
					if len(matches) > 0 {
						hasFile = true
					}
				}
				if !hasFile && m.StashImageID.Valid && m.StashImageID.String != "" {
					matches, _ := filepath.Glob(filepath.Join(stashDir, "images", "*"+m.StashImageID.String+"*"))
					if len(matches) > 0 {
						hasFile = true
					}
				}
			}

			if !hasFile && (blobsDir != "" || stashDir != "") {
				// 既に親なしとしてリストに入っていなければ追加
				alreadyAdded := false
				for _, o := range orphans {
					if o.MediaID == m.MediaID {
						alreadyAdded = true
						break
					}
				}
				if !alreadyAdded {
					orphans = append(orphans, models.OrphanDBMedia{
						MediaID:      m.MediaID,
						ArticleID:    m.ArticleID,
						Type:         m.Type,
						DownloadURL:  m.DownloadURL,
						Status:       m.DownloadStatus,
						StashSceneID: m.StashSceneID.String,
						StashImageID: m.StashImageID.String,
						Reason:       "ステータスはCOMPLETEDですが実ファイルが見つかりません",
					})
				}
			}
		}
	}

	return orphans, nil
}

// GetKnownMediaIdentifiers は全メディアIDおよび Stash Scene/Image ID のインデックスマップを作成します
func GetKnownMediaIdentifiers(db *gorm.DB) (map[string]bool, error) {
	known := make(map[string]bool)

	type idRecord struct {
		MediaID      string         `gorm:"column:media_id"`
		StashSceneID sql.NullString `gorm:"column:stash_scene_id"`
		StashImageID sql.NullString `gorm:"column:stash_image_id"`
	}

	var records []idRecord
	if err := db.Model(&models.Media{}).Select("media_id, stash_scene_id, stash_image_id").Scan(&records).Error; err != nil {
		return nil, err
	}

	for _, r := range records {
		if r.MediaID != "" {
			known[strings.ToLower(r.MediaID)] = true
		}
		if r.StashSceneID.Valid && r.StashSceneID.String != "" {
			known[strings.ToLower(r.StashSceneID.String)] = true
		}
		if r.StashImageID.Valid && r.StashImageID.String != "" {
			known[strings.ToLower(r.StashImageID.String)] = true
		}
	}

	return known, nil
}

// ScanOrphanFiles は stashDir および blobsDir 内の全ファイルを走査し、DB に登録されていない未紐付け孤立ファイルを検出します
func ScanOrphanFiles(stashDir, blobsDir string, knownKeys map[string]bool) ([]models.OrphanFile, error) {
	var orphanFiles []models.OrphanFile

	scanDir := func(dirPath, category string) {
		if dirPath == "" {
			return
		}
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			return
		}

		_ = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil || info == nil || info.IsDir() {
				return nil
			}

			// 除外対象（システムファイル・設定ファイルなど）
			baseName := info.Name()
			if strings.HasPrefix(baseName, ".") || strings.HasSuffix(baseName, ".json") || strings.HasSuffix(baseName, ".yml") || strings.HasSuffix(baseName, ".yaml") {
				return nil
			}

			ext := filepath.Ext(baseName)
			nameWithoutExt := strings.TrimSuffix(baseName, ext)
			lowerName := strings.ToLower(nameWithoutExt)

			// knownKeys との照合
			isKnown := false
			if knownKeys[lowerName] {
				isKnown = true
			} else {
				// 部分一致（Stashファイル名命名規則対応: {scene_id} - title 等）
				for k := range knownKeys {
					if len(k) >= 4 && strings.Contains(lowerName, k) {
						isKnown = true
						break
					}
				}
			}

			if !isKnown {
				orphanFiles = append(orphanFiles, models.OrphanFile{
					Path:     path,
					FileName: baseName,
					FileSize: info.Size(),
					Category: category,
				})
			}

			return nil
		})
	}

	if stashDir != "" {
		scanDir(filepath.Join(stashDir, "scenes"), "stash_scene")
		scanDir(filepath.Join(stashDir, "images"), "stash_image")
		scanDir(filepath.Join(stashDir, "generated"), "stash_generated")
	}

	if blobsDir != "" {
		scanDir(blobsDir, "blob")
	}

	return orphanFiles, nil
}

// DeleteDBMediaByIDs は指定された media_id の DB レコードを削除します
func DeleteDBMediaByIDs(db *gorm.DB, mediaIDs []string) (int, error) {
	if len(mediaIDs) == 0 {
		return 0, nil
	}
	res := db.Where("media_id IN ?", mediaIDs).Delete(&models.Media{})
	return int(res.RowsAffected), res.Error
}
