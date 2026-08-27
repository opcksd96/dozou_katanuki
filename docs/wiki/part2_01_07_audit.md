[[← 06 災害復旧手順|part2_01_06_recovery]] | [[📚 目次 (Home)|Home]] | [[08 常駐スケジューラー →|part2_01_08_scheduler]]

# SPEC-AUDIT-001: SQLite3 整合性監査＆パージプロトコル

## 1. SQLite3 整合性監査 (PRAGMA Audit)
* **`PRAGMA integrity_check;`**: データページ、B-Tree、インデックスの破損を徹底スキャン。破損時は Layer 1 / 2 からの復旧アラートを発行。
* **`PRAGMA foreign_key_check;`**: `accounts` ➔ `articles` ➔ `media` 間の孤立外部キーエラーが0件であることを保証。

## 2. 孤立メディア・ゾンビキャッシュの自動パージ
* **SQLite3 孤立メディア検出**: DBに存在するがStash側にないレコードを検知・削除。
* **Stash 孤立ファイルパージ**: `stash/scenes/` 内を自動走査し、DBの `media_id` と一致しない未紐付けファイルをOSのゴミ箱へ自動退避。

---

[[← 06 災害復旧手順|part2_01_06_recovery]] | [[📚 目次 (Home)|Home]] | [[08 常駐スケジューラー →|part2_01_08_scheduler]]
