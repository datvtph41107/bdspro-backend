package application

import (
	"tqd/internal/usecase/generatedreport/processing"
	"tqd/internal/usecase/usage"
)

// Acceptance mô tả đầy đủ durable facts mà create-report
// yêu cầu PostgreSQL commit như một business acceptance.
//
// Report application tạo intent này.
// Persistence chỉ arbitrate command và lưu các facts trong một transaction.
type Acceptance struct {
	Report Report
	Usage  *usage.Event
	Job    processing.Job
}

// AcceptanceResult trả durable facts sau transaction arbitration.
//
// DurableUsage không phải cờ điều khiển Redis. Nó là usage event
// thực sự thuộc command đã thắng. Current Redis-D bridge dùng evidence
// này để quyết reservation của request hiện tại cần Commit hay Cancel.
type AcceptanceResult struct {
	Report       Report
	Created      bool
	DurableUsage *usage.Event
}
