package processing

import (
	"common/operation"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

/**
 * Status cho biết background job đang ở bước nào.
 */
type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
)

/**
 * Job là công việc tạo file report chạy sau khi request đã accepted.
 */
type Job struct {
	ID       string
	ReportID uint64
	UserID   uint64

	Operation   operation.Code
	OperationID string
	CommandKey  string

	Status   Status
	Attempts int

	AvailableAt  time.Time
	LockedAt     *time.Time
	LockedBy     string
	ClaimVersion int64
	LastError    string

	CreatedAt time.Time
	UpdatedAt time.Time
}

/**
 * Output là file/result worker đã tạo xong cho report.
 */
type Output struct {
	ThumbnailURL string
	ImageURL     string
	PDFURL       string
	ShareURL     string
	FileSize     uint64
	Format       string
}

/**
 * BuildID tạo job ID ổn định cho cùng user + command.
 */
func BuildID(userID uint64, commandKey string) string {
	raw := []byte(fmt.Sprintf("%d\n%s", userID, commandKey))
	sum := sha256.Sum256(raw)
	return "report_job_" + hex.EncodeToString(sum[:16])
}
