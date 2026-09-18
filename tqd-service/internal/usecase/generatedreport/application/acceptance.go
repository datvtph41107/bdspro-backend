package application

import (
	"tqd/internal/usecase/generatedreport/processing"
	"tqd/internal/usecase/usage"
)

// Acceptance là yêu cầu ghi Report + Job và Usage (khi có quota) trong cùng
// transaction PostgreSQL. Commit xác nhận hệ thống đã nhận trách nhiệm tạo
// báo cáo; chưa có nghĩa PDF đã sẵn sàng hoặc client đã tải được file.
// Application quyết định nghiệp vụ; persistence phân xử command trùng và ghi
// dữ liệu nguyên tử. Worker retry/regenerate dùng lại job đã được chấp nhận.
type Acceptance struct {
	Report Report
	Usage  *usage.Event
	Job    processing.Job
}

// AcceptanceResult là kết quả sau khi database phân xử command idempotency.
// Created=false nghĩa là retry đã gặp report cũ và phải dùng lại kết quả đó.
//
// DurableUsage không phải cờ điều khiển Redis. Nó là usage event
// thực sự thuộc command đã thắng. Application dùng bằng chứng này để quyết
// định reservation của request hiện tại cần Commit hay Cancel.
type AcceptanceResult struct {
	Report       Report
	Created      bool
	DurableUsage *usage.Event
}
