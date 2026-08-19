//go:build windows
// +build windows

package driver

import (
	"os"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"
)

var (
	shell32          = syscall.NewLazyDLL("shell32.dll")
	shFileOperationW = shell32.NewProc("SHFileOperationW")
)

const (
	foDelete          = 0x0003
	fofAllowUndo      = 0x0040
	fofNoConfirmation = 0x0010
	fofSilent         = 0x0004
	fofNoErrorUI      = 0x0400
)

type shFileOpStructW struct {
	hwnd                  uintptr
	wFunc                 uint32
	pFrom                 *uint16
	pTo                   *uint16
	fFlags                uint16
	fAnyOperationsAborted int32
	hNameMappings         uintptr
	lpszProgressTitle     *uint16
}

// MoveToRecycleBin は Windows のごみ箱 (Recycle Bin) へ指定ファイルを安全に退避します
func MoveToRecycleBin(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil // 既に存在しない場合はスキップ
	}

	// ダブルNULL終端文字列を作成
	utf16Chars, err := syscall.UTF16FromString(absPath)
	if err != nil {
		return fallbackMoveToTrash(absPath)
	}
	// ダブルNULL終端を追加
	utf16Chars = append(utf16Chars, 0)

	op := shFileOpStructW{
		wFunc:  foDelete,
		pFrom:  &utf16Chars[0],
		fFlags: fofAllowUndo | fofNoConfirmation | fofSilent | fofNoErrorUI,
	}

	ret, _, _ := shFileOperationW.Call(uintptr(unsafe.Pointer(&op)))
	if ret != 0 {
		// API呼び出しに失敗した場合は安全に backups/trash へフォールバック移動
		return fallbackMoveToTrash(absPath)
	}

	return nil
}

func fallbackMoveToTrash(filePath string) error {
	trashDir := filepath.Join("backups", "trash", time.Now().Format("20060102"))
	if err := os.MkdirAll(trashDir, 0755); err != nil {
		return err
	}
	destPath := filepath.Join(trashDir, filepath.Base(filePath))
	return os.Rename(filePath, destPath)
}
