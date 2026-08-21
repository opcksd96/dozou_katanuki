// driver/audit_integrity.go (100行以下)
package driver

import (
	"database/sql"
	"fmt"
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
