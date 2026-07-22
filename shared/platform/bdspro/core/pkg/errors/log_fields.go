// platform/errors/log_fields.go
//
// LogFields chuyển AppError thành map[string]any cho structured logging.
//
// Tại sao không dùng zap.Field trực tiếp?
// Platform không muốn hard-depend vào zap ở error package.
// Logger package nhận map[string]any và convert sang zap.Field — decoupled.
package apperr

// LogFields trả về map đầy đủ thông tin để log.
// Khác với Fields() (strip Value), hàm này giữ nguyên để debug.
//
// Output ví dụ:
//
//	{
//	  "kind": "internal",
//	  "code": "INTERNAL",
//	  "op": "UserService.Register:UserRepo.Create:DB.Query",
//	  "message": "internal server error",
//	  "cause": "connection refused",
//	  "stack": "goroutine 1 [running]:\n...",   ← chỉ khi KindInternal
//	  "retryable": false,
//	}
func LogFields(err error) map[string]any {
	if err == nil {
		return nil
	}

	ae, ok := As(err)
	if !ok {
		// Non-AppError — vẫn log được, kind=unknown để filter
		return map[string]any{
			"kind":  "unknown",
			"code":  "UNKNOWN",
			"error": err.Error(),
		}
	}

	out := map[string]any{
		"kind":      ae.Kind().String(),
		"code":      string(ae.Code()),
		"message":   ae.Message(),
		"retryable": ae.Retryable(),
	}

	// Op chain: chỉ add nếu có — tránh empty string trong log
	if ae.Op() != "" {
		out["op"] = ae.Op()
	}

	// RetryAfter: chỉ add nếu > 0
	if ae.RetryAfter() > 0 {
		out["retry_after"] = ae.RetryAfter()
	}

	// Cause: root error message — giúp debug không cần stack
	if ae.Cause() != nil {
		out["cause"] = ae.Cause().Error()
	}

	// Fields: validation errors — đầy đủ kể cả Value (internal log)
	if ae.HasFields() {
		// Internal log giữ nguyên fields với Value
		out["fields"] = ae.fields // access private field trong cùng package
	}

	// Stack: CHỈ log cho KindInternal (F2)
	// Validation/Business/External errors không cần stack — không phải bug
	if ae.Kind() == KindInternal && len(ae.Stack()) > 0 {
		out["stack"] = string(ae.Stack())
	}

	return out
}
