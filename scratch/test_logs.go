package main

import (
	"encoding/json"
	"fmt"
	"os"

	"dozou_katanuki/app"
)

func main() {
	a := &app.App{}
	logs, err := a.GetPipelineLogs("all", 5)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	bytes, _ := json.MarshalIndent(logs, "", "  ")
	fmt.Println(string(bytes))
}
