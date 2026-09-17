package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"syscall"
	"time"
)

// Khóa process trên thư mục riêng. append/truncate không cho hai process ghi/xóa cùng file.
// Native/Docker backend chạy Linux; flock tự giải phóng cả khi process crash.
func lockDirectory(directory string) (func() error, error) {
	if err := os.MkdirAll(directory, 0700); err != nil {
		return nil, err
	}
	file, err := os.OpenFile(filepath.Join(directory, ".writer.lock"), os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	if err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		file.Close()
		return nil, fmt.Errorf("log directory already has a writer: %s: %w", directory, err)
	}
	return file.Close, nil
}

var runName = regexp.MustCompile(`^\d{8}T\d{6}\.\d{9}Z-\d+$`)

// Chỉ dọn thư mục run mang tên do logger tạo; kiểm tra lock trước khi xóa.
// Không động tới file khác, run hiện tại hoặc run của process còn sống.
func pruneRuns(root, current string, keep, days int) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Name() > entries[j].Name() })
	retained := 1
	for _, entry := range entries {
		path := filepath.Join(root, entry.Name())
		if !entry.IsDir() || !runName.MatchString(entry.Name()) || path == current {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if retained < keep && time.Since(info.ModTime()) <= time.Duration(days)*24*time.Hour {
			retained++
			continue
		}
		release, err := lockDirectory(path)
		if err != nil {
			continue
		}
		_ = os.RemoveAll(path)
		_ = release()
	}
}

// truncate xóa nội dung mọi active file đã biết ngay khi start, kể cả category
// chưa có event ở phiên mới; tránh dev đọc nhầm ERROR cũ là lỗi của phiên hiện tại.
func truncateActiveFiles(directory string) error {
	for _, level := range []string{"debug", "info", "warn", "error"} {
		for _, category := range categories {
			path := filepath.Join(directory, level, category+".jsonl")
			info, err := os.Lstat(path)
			if os.IsNotExist(err) {
				continue
			}
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return fmt.Errorf("refuse to truncate non-regular log: %s", path)
			}
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0600)
			if err != nil {
				return err
			}
			if err = f.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}
