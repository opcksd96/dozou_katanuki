// ignition.go (100行以下 - SPEC-PRINCIPLE-001)
// Wails のファイルウォッチャーに変更を検知させ、ホットリロード＆バインディング再生成を強制点火(Ignition)するためのトリガーファイルです。
package main

import "time"

// LastIgnitionTimestamp は Wails リビルドをトリガーした最終更新時刻です
var LastIgnitionTimestamp = "2026-08-29T07:56:00+09:00"

// TouchIgnition はイグニッション更新用ヘルパー関数です
func TouchIgnition() string {
	return time.Now().Format(time.RFC3339)
}
