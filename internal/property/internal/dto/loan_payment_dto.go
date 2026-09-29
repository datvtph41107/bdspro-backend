package dto

// CalculateLoanPaymentRequest DTO cho request tính lịch trả nợ vay ngân hàng
type CalculateLoanPaymentRequest struct {
	InterestRate         float64  `json:"interestRate" binding:"required,gt=0"`  // Lãi suất năm đầu (%/năm) - thường là lãi suất ưu đãi
	PropertyValue        float64  `json:"propertyValue" binding:"required,gt=0"` // Giá trị tài sản (VND)
	LoanAmount           float64  `json:"loanAmount" binding:"required,gt=0"`    // Số tiền vay (VND)
	Term                 uint32   `json:"term" binding:"required,gt=0"`          // Kỳ hạn (tháng)
	PaymentMethod        uint32   `json:"paymentMethod" binding:"required"`      // Phương thức trả nợ: 10=dư nợ giảm dần, 20=trả lãi định kỳ/gốc cuối kỳ
	FloatingInterestRate *float64 `json:"floatingInterestRate,omitempty"`        // Lãi suất thả nổi từ năm 2 (%/năm) - thường 10-11%
}

// CalculateLoanPaymentResponse DTO cho response tính lịch trả nợ
type CalculateLoanPaymentResponse struct {
	Schedule      []LoanPaymentScheduleItem `json:"schedule"`      // Lịch trả nợ theo tháng
	TotalInterest float64                   `json:"totalInterest"` // Tổng tiền lãi
	TotalPayment  float64                   `json:"totalPayment"`  // Tổng tiền trả
}

// LoanPaymentScheduleItem Chi tiết trả nợ từng tháng
type LoanPaymentScheduleItem struct {
	Month              uint32  `json:"month"`              // Tháng
	RemainingPrincipal float64 `json:"remainingPrincipal"` // Dư nợ còn lại
	PrincipalPayment   float64 `json:"principalPayment"`   // Tiền gốc trả trong tháng
	InterestPayment    float64 `json:"interestPayment"`    // Tiền lãi trả trong tháng
	TotalPayment       float64 `json:"totalPayment"`       // Tổng trả trong tháng
	InterestRate       float64 `json:"interestRate"`       // Lãi suất áp dụng cho tháng này (%/năm)
}
