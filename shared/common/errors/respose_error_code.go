package _errors

func (e ErrorCode) GetCode() int32 {
	return e.Code
}

func (e ErrorCode) GetMessage() string {
	return e.Message
}

var (
	NOT_ENTERED_USER          = ErrorCode{400, "User information is required"}
	EXPIRED_VERIFICATION_CODE = ErrorCode{400, "Verification code has expired"}
	INVALID_VERIFICATION_CODE = ErrorCode{400, "Invalid verification code"}

	// 409 Conflict
	ALREADY_EXIST_USER  = ErrorCode{409, "User already exists"}
	ALREADY_EXIST_EMAIL = ErrorCode{409, "Email already exists"}

	// 500 Internal Server Error
	NOT_EXIST_TAX_RATE = ErrorCode{500, "Tax rate does not exist"}
	PAYMENT_ERROR      = ErrorCode{500, "Payment failed"}
)
