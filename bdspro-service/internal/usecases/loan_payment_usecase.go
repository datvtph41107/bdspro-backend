package usecases

import (
	"bdspro/internal/dto"
	"context"
	"math"
)

type LoanPaymentUsecase struct{}

func NewLoanPaymentUsecase() *LoanPaymentUsecase {
	return &LoanPaymentUsecase{}
}

// CalculateLoanPayment Tính toán lịch trả nợ vay ngân hàng
func (uc *LoanPaymentUsecase) CalculateLoanPayment(ctx context.Context, req *dto.CalculateLoanPaymentRequest) (*dto.CalculateLoanPaymentResponse, error) {
	loanAmount := req.LoanAmount
	interestRatePercent := req.InterestRate  // Lưu giá trị % gốc
	interestRate := req.InterestRate / 100.0 // Chuyển từ % sang số thập phân
	term := int(req.Term)

	// Tính lãi suất thả nổi nếu có
	var floatingInterestRatePercent float64 // Lưu giá trị % gốc
	var floatingInterestRate float64
	hasFloatingRate := req.FloatingInterestRate != nil && *req.FloatingInterestRate > 0
	if hasFloatingRate {
		floatingInterestRatePercent = *req.FloatingInterestRate
		floatingInterestRate = *req.FloatingInterestRate / 100.0
	}

	var schedule []dto.LoanPaymentScheduleItem

	switch req.PaymentMethod {
	case 10: // Dư nợ giảm dần
		schedule = uc.calculateDecreasingBalance(loanAmount, interestRate, interestRatePercent, floatingInterestRate, floatingInterestRatePercent, hasFloatingRate, term)
	case 20: // Trả lãi định kỳ/gốc cuối kỳ
		schedule = uc.calculateInterestOnly(loanAmount, interestRate, interestRatePercent, floatingInterestRate, floatingInterestRatePercent, hasFloatingRate, term)
	default:
		// Mặc định dùng dư nợ giảm dần
		schedule = uc.calculateDecreasingBalance(loanAmount, interestRate, interestRatePercent, floatingInterestRate, floatingInterestRatePercent, hasFloatingRate, term)
	}

	// Tính tổng lãi và tổng trả
	var totalInterest float64
	var totalPayment float64
	for _, item := range schedule {
		totalInterest += item.InterestPayment
		totalPayment += item.TotalPayment
	}

	return &dto.CalculateLoanPaymentResponse{
		Schedule:      schedule,
		TotalInterest: totalInterest,
		TotalPayment:  totalPayment,
	}, nil
}

// calculateDecreasingBalance Tính lịch trả nợ theo phương thức dư nợ giảm dần
// Mỗi tháng trả gốc cố định, lãi tính trên dư nợ còn lại
// Năm đầu dùng interestRate, từ năm 2 (tháng 13) trở đi dùng floatingInterestRate nếu có
func (uc *LoanPaymentUsecase) calculateDecreasingBalance(loanAmount, interestRate, interestRatePercent, floatingInterestRate, floatingInterestRatePercent float64, hasFloatingRate bool, term int) []dto.LoanPaymentScheduleItem {
	schedule := make([]dto.LoanPaymentScheduleItem, term)
	monthlyPrincipal := loanAmount / float64(term) // Gốc trả mỗi tháng
	remainingPrincipal := loanAmount

	for i := 0; i < term; i++ {
		// Xác định lãi suất tháng dựa trên tháng hiện tại
		var monthlyInterestRate float64
		var currentInterestRate float64 // Lãi suất năm áp dụng cho tháng này (%/năm)
		if hasFloatingRate && i >= 12 {
			// Từ tháng 13 trở đi (năm 2+) dùng lãi suất thả nổi
			monthlyInterestRate = floatingInterestRate / 12.0
			currentInterestRate = floatingInterestRatePercent
		} else {
			// Tháng 1-12 (năm đầu) dùng lãi suất ưu đãi
			monthlyInterestRate = interestRate / 12.0
			currentInterestRate = interestRatePercent
		}

		interestPayment := remainingPrincipal * monthlyInterestRate
		principalPayment := monthlyPrincipal

		// Tháng cuối cùng, trả hết phần dư nợ còn lại
		if i == term-1 {
			principalPayment = remainingPrincipal
		}

		totalPayment := principalPayment + interestPayment
		remainingPrincipal -= principalPayment

		// Đảm bảo dư nợ không âm do làm tròn
		if remainingPrincipal < 0.01 {
			remainingPrincipal = 0
		}

		schedule[i] = dto.LoanPaymentScheduleItem{
			Month:              uint32(i + 1),
			RemainingPrincipal: math.Round(remainingPrincipal*100) / 100,
			PrincipalPayment:   math.Round(principalPayment*100) / 100,
			InterestPayment:    math.Round(interestPayment*100) / 100,
			TotalPayment:       math.Round(totalPayment*100) / 100,
			InterestRate:       math.Round(currentInterestRate*100) / 100,
		}
	}

	return schedule
}

// calculateInterestOnly Tính lịch trả nợ theo phương thức trả lãi định kỳ/gốc cuối kỳ
// Mỗi tháng chỉ trả lãi, gốc trả cuối kỳ
// Năm đầu dùng interestRate, từ năm 2 (tháng 13) trở đi dùng floatingInterestRate nếu có
func (uc *LoanPaymentUsecase) calculateInterestOnly(loanAmount, interestRate, interestRatePercent, floatingInterestRate, floatingInterestRatePercent float64, hasFloatingRate bool, term int) []dto.LoanPaymentScheduleItem {
	schedule := make([]dto.LoanPaymentScheduleItem, term)
	remainingPrincipal := loanAmount

	for i := 0; i < term; i++ {
		// Xác định lãi suất tháng dựa trên tháng hiện tại
		var monthlyInterestRate float64
		var currentInterestRate float64 // Lãi suất năm áp dụng cho tháng này (%/năm)
		if hasFloatingRate && i >= 12 {
			// Từ tháng 13 trở đi (năm 2+) dùng lãi suất thả nổi
			monthlyInterestRate = floatingInterestRate / 12.0
			currentInterestRate = floatingInterestRatePercent
		} else {
			// Tháng 1-12 (năm đầu) dùng lãi suất ưu đãi
			monthlyInterestRate = interestRate / 12.0
			currentInterestRate = interestRatePercent
		}

		interestPayment := remainingPrincipal * monthlyInterestRate
		var principalPayment float64
		var totalPayment float64

		if i == term-1 {
			// Tháng cuối cùng trả cả gốc và lãi
			principalPayment = remainingPrincipal
			totalPayment = principalPayment + interestPayment
			remainingPrincipal = 0
		} else {
			// Các tháng khác chỉ trả lãi, dư nợ không đổi
			principalPayment = 0
			totalPayment = interestPayment
			// remainingPrincipal vẫn giữ nguyên loanAmount
		}

		schedule[i] = dto.LoanPaymentScheduleItem{
			Month:              uint32(i + 1),
			RemainingPrincipal: math.Round(remainingPrincipal*100) / 100,
			PrincipalPayment:   math.Round(principalPayment*100) / 100,
			InterestPayment:    math.Round(interestPayment*100) / 100,
			TotalPayment:       math.Round(totalPayment*100) / 100,
			InterestRate:       math.Round(currentInterestRate*100) / 100,
		}
	}

	return schedule
}
