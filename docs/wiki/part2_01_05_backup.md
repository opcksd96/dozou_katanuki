[[← 04 Admin Board仕様|part2_01_04_admin_board]] | [[📚 目次 (Home)|Home]] | [[06 災害復旧手順 →|part2_01_06_recovery]]

# SPEC-BACKUP-001: 二重化バックアップ（Dual-Source DR）仕様

## 1. 二重系統データ保全アーキテクチャ

```mermaid
flowchart TD
    DB["SQLite3 (archive.db)<br>実稼働マスター"]
    
    subgraph L1 ["Layer 1: Fast Path (バイナリ復元)"]
        F1["VACUUM INTO コピー<br>backups/database/archive_*.db"]
        F2["・ミリ秒即時復旧 (RTO極小)<br>・リレーション完全維持"]
    end
    
    subgraph L2 ["Layer 2: Deep Path (生データ原本)"]
        D1["原本魚拓 ＆ メタデータ<br>backups/dumps/"]
        D2["・ISO 28500 warc.gz 原本<br>・metadata.json 共通中間表現"]
    end
    
    DB --> L1
    DB --> L2
```

## 2. アバター保全・隔離ポリシー (Avatar Isolation)
* **ライブラリ汚染防止**: Stashの画像グリッド混入を防ぐため、アバター実ファイルはすべて `backups/dumps/{platform}/{username}/avatars/` に物理隔離保存します。
* **原本性保証**: `metadata.json` 内には生のアバターオリジナルURLを不変データとして保持します。

---

[[← 04 Admin Board仕様|part2_01_04_admin_board]] | [[📚 目次 (Home)|Home]] | [[06 災害復旧手順 →|part2_01_06_recovery]]
