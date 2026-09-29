package domain

// AdminOrder là projection chỉ đọc của toàn bộ chuỗi bằng chứng thương mại.
// Payment-service là canonical owner nên admin không cần ghép dữ liệu từ DB
// của service khác và cũng không được sửa trực tiếp các bảng này.
type AdminOrder struct {
	Order              Order
	Attempts           []PaymentAttempt
	Settlement         *Settlement
	Fulfillment        *Fulfillment
	ManualConfirmation *CommandEffect
}

// AdminFulfillment gắn fulfillment với order để admin có đủ ngữ cảnh xử lý
// retry/review mà không phải suy diễn từ subject hoặc frontend state.
type AdminFulfillment struct {
	Fulfillment Fulfillment
	Order       AdminOrder
}
