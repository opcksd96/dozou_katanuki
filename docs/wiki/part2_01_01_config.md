[[← 00 総合インデックス|part2_01_00_index]] | [[📚 目次 (Home)|Home]] | [[02 プロセス制御仕様 →|part2_01_02_lifecycle]]

# SPEC-CONFIG-001: 一元設定ポータル (config.json) 仕様

## 1. 設定構造ハイエラルキー (Single Source of Truth)
従来散逸していたYAMLや個別環境変数を全廃し、ルート直下の `config.json` をシステム唯一の設定の源泉（SSOT）とします。

* **system**: 実行環境（`env`）、デフォルトプラットフォーム、UI言語（`ja/en/zh`）。
* **network**: ポート定義（内部閉塞 `stash_port: 9999` 等）、バインドアドレス。
* **storage**: DBパス（`db_path`）、`stash_enabled`、物理メディア保存先。
* **scheduler**: ポーリング間隔、自動バックアップ周期、最大保持世代数。
* **broadcast**: 家庭内LANキャスト有効化フラグ、許可サブネット（CIDR）。
* **appearance**: 3言語対応の優先フォントファミリー定義。

## 2. Stash `config.yml` 透過同期仕様 (Zero-Config Stash Bridge)
`dozou_katanuki` の `config.json` を変更・保存（`SaveConfig` RPC）したタイミングで、内部に同包される `bin/config.yml` のコア設定（`host`, `port`, `dangerous_allow_public_without_auth`）が透過的・安全に自動同期されます。
* **セルフリブート保護**: Stash 起動時ではなく、dozou の設定保存時のみ同期を走らせることで、Stash 自身のセルフリブートや設定変更との競合・意図せぬ巻き戻りを完全に防止します。
* **LAN透過アクセス**: `host: 0.0.0.0`, `port: 9999`, `dangerous_allow_public_without_auth: "true"` が保証され、LAN内の別端末からでも直接 `http://<HostIP>:9999` で Stash WebUI にアクセス可能となります。

## 3. 「Stash使わんし！」モード (物理フォルダ保存ポリシー)
`storage.stash_enabled` が `false` の場合、システムは Stashapp を起動せず、軽量な物理フォルダダイレクトサーブモードへ自律移行します。

* **物理マッピング**: ダウンロード完了ファイルを `{local_media_dir}/{platform}/{username}/{media_id}` にフラット保存。
* **相対URL動的解決**: タイムライン構築時、メディアURLは `/media-local/{platform}/{username}/{media_id}` へ自動置換され、軽量Goバイナリ単体で動態保存が完結します。

---

[[← 00 総合インデックス|part2_01_00_index]] | [[📚 目次 (Home)|Home]] | [[02 プロセス制御仕様 →|part2_01_02_lifecycle]]
