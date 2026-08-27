[[← 01 一元設定仕様|part2_01_01_config]] | [[📚 目次 (Home)|Home]] | [[03 インメモリプロキシ仕様 →|part2_01_03_proxy]]

# SPEC-LIFECYCLE-001: Wails-Stash プロセスライフサイクル制御仕様

## 1. 起動・終了シーケンス (Non-blocking Prober & Lifeline Sync)

```mermaid
sequenceDiagram
    autonumber
    participant UI as Frontend (Vue 3 Dumb UI)
    participant Wails as Wails Entry (app.go)
    participant Prober as Middleware StashProber
    participant Stash as Stash (stash-win.exe)

    Note over UI, Wails: 【1. 0.5秒 爆速起動 (完全ノンブロッキング)】
    Wails->>UI: DOMレンダリング開始 (即座に WindowShow)
    Wails->>Prober: Start() (バックグラウンド起動)

    Note over Prober, Stash: 【2. 200ms 非同期プロービング & 自律キック】
    loop 200ms 間隔
        Prober->>Stash: GET http://127.0.0.1:9999/
    end

    alt 2秒以上未接続の場合 (未起動検知)
        Prober->>UI: 🟡 トースト「📦 Stash をバックグラウンド起動中...」
        Prober->>Stash: exec(stash-win.exe) を非同期キックして即座離脱
    end

    Stash-->>Prober: HTTP 200 OK (疎通完了)
    Prober->>UI: 🟢 トースト「🟢 Stash 接続完了！ (Port 9999)」 ＆ stash:ready 発火
    UI->>UI: メディア画像・動画DOMをスッと更新・フェードイン！

    Note over UI, Stash: --- Wails終了時 ---
    Wails->>Prober: Stop()
    Prober->>Stash: taskkill /F /IM stash-win.exe (Stash道連れ完全終了)
```

## 2. プロセス制御規約
* **完全ノンブロッキング**: Wails の `startup` / `domReady` では Stash 起動を一切同期待機せず、0.5秒未満で即座にUIを展開。
* **ミドルウェア層プローバー**: `middleware.StashProber` が 200ms 間隔でポート 9999 を非同期プロービングし、未起動時は自律キック、セルフリブート時も自動追従。
* **トースト＆DOM更新連携**: 接続完了・切断・自動起動ステータスをUIトーストにリアルタイム通知し、メディアDOMを安全に再描画。
* **道連れ終了**: 親ウィンドウ終了時（`OnShutdown`）に確実に `taskkill /F /IM stash-win.exe` を発行してゾンビプロセスを完全根絶。

---

[[← 01 一元設定仕様|part2_01_01_config]] | [[📚 目次 (Home)|Home]] | [[03 インメモリプロキシ仕様 →|part2_01_03_proxy]]
