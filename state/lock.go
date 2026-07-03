package state

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// LockStage 获取 Stage 的文件锁，防止并发运行同一 Stage。
// 返回的 *os.File 需要 defer UnlockStage(f) 释放。
func LockStage(root, stageID string) (*os.File, error) {
	stage, err := LoadStage(root, stageID)
	if err != nil {
		return nil, fmt.Errorf("lock: %w", err)
	}
	dir := filepath.Join(root, DirStages, stageDirName(stage))
	os.MkdirAll(dir, 0755)
	lockPath := filepath.Join(dir, ".lock")
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("lock file: %w", err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		f.Close()
		return nil, fmt.Errorf("stage %s 正在运行中（已有另一个进程锁定）", stageID)
	}
	return f, nil
}

// UnlockStage 释放文件锁。
func UnlockStage(f *os.File) {
	syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	f.Close()
}
