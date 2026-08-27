[[← 05 二重化バックアップ仕様|part2_01_05_backup]] | [[📚 目次 (Home)|Home]] | [[07 DB健全性監査 →|part2_01_07_audit]]

# SPEC-RECOVERY-001: 災害復旧（完全オフライン自動リストア）手順

## 1. ゼロからの自動再構築設計
実稼働DB `archive.db` が完全破壊された場合、Layer 2 (生データ原本) のみから外部通信ゼロ・完全オフラインで動態保存状態を100%自動再構築します。

## 2. 対称リストアフロー

```mermaid
graph TD
    A[実稼働 archive.db が全損] --> B[1. 破損 archive.db を物理削除]
    B --> C[2. Wails Core 起動 ➔ GORM AutoMigrate で空DB生成]
    C --> D[3. Pythonサイドカーをリストアモードでキック]
    D --> E[4. dumps 配下の metadata.json と snapshot.warc.gz を一括走査]
    E --> F[5. Core API POST /api/articles へ無加工で順次投入]
    F --> G[6. メディア実体をStashへ再配置 ＆ UUID逆引きバインド]
    G --> H[リストア完了: 100%同一の timeline 導通状態が完全復帰]

    style A fill:#ffebee,stroke:#f44336,stroke-width:2px
    style H fill:#e8f5e9,stroke:#4caf50,stroke-width:2px
```

---

[[← 05 二重化バックアップ仕様|part2_01_05_backup]] | [[📚 目次 (Home)|Home]] | [[07 DB健全性監査 →|part2_01_07_audit]]
