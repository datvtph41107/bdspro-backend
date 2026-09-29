package _err

type ErrorDTO struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Cài đặt phương thức Error() string để thỏa interface error
func (e *ErrorDTO) Error() string {
	return e.Message
}

func ErrorReturn(code int, message string) *ErrorDTO {
	return &ErrorDTO{
		Code:    code,
		Message: message,
	}
}
