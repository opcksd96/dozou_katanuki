package main

import (
	"fmt"
	"dozou_katanuki/app"
)

func main() {
	tasks, err := app.GetTasksViaDirectCDP(9222)
	if err != nil {
		fmt.Printf("CDP error: %v\n", err)
		return
	}
	fmt.Printf("⚡ 迅雷内部アクティブタスク: 合計 %d 件\n", len(tasks))
	statusMap := make(map[string]int)
	for _, t := range tasks {
		statusMap[t.Status]++
	}
	for st, cnt := range statusMap {
		fmt.Printf("  %s: %d 件\n", st, cnt)
	}

	fmt.Println("\n【先頭10件の詳細】")
	for i, t := range tasks {
		if i >= 10 { break }
		fmt.Printf("[%d] %-16s | %s | %s | %s | %s (Err:%d, St:%d)\n",
			i+1, t.Status, t.FileName, t.ProgressText, t.SpeedText, t.DetailText, t.ErrorCode, t.TaskStatusCode)
	}
}
