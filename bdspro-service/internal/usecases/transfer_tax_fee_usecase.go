package usecases

import (
	"bdspro/internal/dto"
	"context"
	"fmt"
	"math"
)

type TransferTaxFeeUsecase struct{}

func NewTransferTaxFeeUsecase() *TransferTaxFeeUsecase {
	return &TransferTaxFeeUsecase{}
}

// CalculateTransferTaxFee Tính toán thuế và phí chuyển nhượng bất động sản
func (uc *TransferTaxFeeUsecase) CalculateTransferTaxFee(ctx context.Context, req *dto.CalculateTransferTaxFeeRequest) (*dto.CalculateTransferTaxFeeResponse, error) {
	transferValue := req.TransferValue

	// Thuế thu nhập cá nhân (TNCN): 2% trên giá trị chuyển nhượng (nếu là cá nhân)
	var personalIncomeTax float64
	if !req.IsBusiness {
		// Nếu có giá gốc, tính thuế trên phần chênh lệch
		if req.OriginalPrice != nil && *req.OriginalPrice > 0 && *req.OriginalPrice < transferValue {
			profit := transferValue - *req.OriginalPrice
			personalIncomeTax = profit * 0.02 // 2% trên phần lợi nhuận
		} else {
			// Nếu không có giá gốc hoặc giá gốc >= giá chuyển nhượng, tính 2% trên toàn bộ giá trị
			personalIncomeTax = transferValue * 0.02
		}
	}

	// Thuế thu nhập doanh nghiệp: 2% trên giá trị chuyển nhượng (nếu là doanh nghiệp)
	var corporateTax float64
	if req.IsBusiness {
		corporateTax = transferValue * 0.02
	}

	// Lệ phí trước bạ: 0.5% trên giá trị chuyển nhượng
	registrationFee := transferValue * 0.005

	// Phí công chứng: tính theo bậc thang
	// - Dưới 50 triệu: 50.000 VND
	// - Từ 50 triệu đến 100 triệu: 100.000 VND
	// - Từ 100 triệu đến 1 tỷ: 0.1% trên giá trị
	// - Trên 1 tỷ: 0.05% trên giá trị (tối đa 10 triệu)
	notarizationFee := uc.calculateNotarizationFee(transferValue)

	// Phí thẩm định hồ sơ: cố định theo quy định
	// Đất ở, đất nông nghiệp: 190.000 VND (dưới 100m²), 260.000 VND (từ 100m² trở lên)
	// Đất sản xuất kinh doanh: 260.000 VND (dưới 100m²), 320.000 VND (từ 100m² trở lên)
	// Mặc định dùng giá trị phổ biến: 260.000 VND
	dossierReviewFee := 260_000

	// Lệ phí địa chính: cố định
	// Phường/thị trấn: 15.000 VND, khu vực khác: 5.000 VND
	// Mặc định: 15.000 VND
	cadastralFee := 15_000

	// Lệ phí cấp GCN (Giấy chứng nhận): cố định
	// Thường từ 50.000 - 100.000 VND, mặc định: 100.000 VND
	certificateFee := 100_000

	// Phí đo vẽ/địa chính: thường tính theo diện tích
	// Mặc định: 200.000 VND (có thể điều chỉnh theo diện tích thực tế)
	surveyFee := 200_000

	// Tổng thuế và phí
	totalTaxFee := personalIncomeTax + corporateTax + registrationFee + notarizationFee +
		float64(dossierReviewFee) + float64(cadastralFee) + float64(certificateFee) + float64(surveyFee)

	// Số tiền thực nhận sau khi trừ thuế phí
	netAmount := transferValue - totalTaxFee

	// Tạo breakdown với ghi chú
	breakdown := dto.TaxFeeBreakdown{
		PersonalIncomeTaxNote: uc.getPersonalIncomeTaxNote(personalIncomeTax, req),
		CorporateTaxNote:      uc.getCorporateTaxNote(corporateTax),
		RegistrationFeeNote:   fmt.Sprintf("Lệ phí trước bạ: 0.5%% trên giá trị chuyển nhượng = %s VND", uc.formatCurrency(registrationFee)),
		NotarizationFeeNote:   uc.getNotarizationFeeNote(transferValue, notarizationFee),
		DossierReviewFeeNote:  fmt.Sprintf("Phí thẩm định hồ sơ: %s VND (mặc định cho đất ở từ 100m²)", uc.formatCurrency(float64(dossierReviewFee))),
		CadastralFeeNote:      fmt.Sprintf("Lệ phí địa chính: %s VND (mặc định cho phường/thị trấn)", uc.formatCurrency(float64(cadastralFee))),
		CertificateFeeNote:    fmt.Sprintf("Lệ phí cấp GCN: %s VND", uc.formatCurrency(float64(certificateFee))),
		SurveyFeeNote:         fmt.Sprintf("Phí đo vẽ/địa chính: %s VND (có thể thay đổi theo diện tích)", uc.formatCurrency(float64(surveyFee))),
	}

	return &dto.CalculateTransferTaxFeeResponse{
		TransferValue:     transferValue,
		PersonalIncomeTax: personalIncomeTax,
		CorporateTax:      corporateTax,
		RegistrationFee:   registrationFee,
		NotarizationFee:   notarizationFee,
		DossierReviewFee:  float64(dossierReviewFee),
		CadastralFee:      float64(cadastralFee),
		CertificateFee:    float64(certificateFee),
		SurveyFee:         float64(surveyFee),
		TotalTaxFee:       totalTaxFee,
		NetAmount:         netAmount,
		Breakdown:         breakdown,
	}, nil
}

// calculateNotarizationFee Tính phí công chứng theo bậc thang
func (uc *TransferTaxFeeUsecase) calculateNotarizationFee(transferValue float64) float64 {
	switch {
	case transferValue < 50_000_000:
		return 50_000
	case transferValue < 100_000_000:
		return 100_000
	case transferValue < 1_000_000_000:
		return transferValue * 0.001 // 0.1%
	default:
		fee := transferValue * 0.0005 // 0.05%
		if fee > 10_000_000 {
			return 10_000_000 // Tối đa 10 triệu
		}
		return fee
	}
}

// getPersonalIncomeTaxNote Tạo ghi chú cho thuế TNCN
func (uc *TransferTaxFeeUsecase) getPersonalIncomeTaxNote(tax float64, req *dto.CalculateTransferTaxFeeRequest) string {
	if req.IsBusiness {
		return "Không áp dụng (người chuyển nhượng là doanh nghiệp)"
	}
	if req.OriginalPrice != nil && *req.OriginalPrice > 0 && *req.OriginalPrice < req.TransferValue {
		// profit := req.TransferValue - *req.OriginalPrice
		return fmt.Sprintf("Thuế TNCN: 2%% trên phần lợi nhuận (%s - %s) = %s VND",
			uc.formatCurrency(req.TransferValue),
			uc.formatCurrency(*req.OriginalPrice),
			uc.formatCurrency(tax))
	}
	return fmt.Sprintf("Thuế TNCN: 2%% trên giá trị chuyển nhượng = %s VND", uc.formatCurrency(tax))
}

// getCorporateTaxNote Tạo ghi chú cho thuế doanh nghiệp
func (uc *TransferTaxFeeUsecase) getCorporateTaxNote(tax float64) string {
	if tax == 0 {
		return "Không áp dụng (người chuyển nhượng là cá nhân)"
	}
	return fmt.Sprintf("Thuế thu nhập doanh nghiệp: 2%% trên giá trị chuyển nhượng = %s VND", uc.formatCurrency(tax))
}

// getNotarizationFeeNote Tạo ghi chú cho phí công chứng
func (uc *TransferTaxFeeUsecase) getNotarizationFeeNote(transferValue, fee float64) string {
	switch {
	case transferValue < 50_000_000:
		return fmt.Sprintf("Phí công chứng: Dưới 50 triệu = %s VND", uc.formatCurrency(fee))
	case transferValue < 100_000_000:
		return fmt.Sprintf("Phí công chứng: Từ 50-100 triệu = %s VND", uc.formatCurrency(fee))
	case transferValue < 1_000_000_000:
		return fmt.Sprintf("Phí công chứng: 0.1%% trên giá trị chuyển nhượng = %s VND", uc.formatCurrency(fee))
	default:
		return fmt.Sprintf("Phí công chứng: 0.05%% trên giá trị chuyển nhượng (tối đa 10 triệu) = %s VND", uc.formatCurrency(fee))
	}
}

// formatCurrency Định dạng số tiền với dấu phẩy ngăn cách hàng nghìn
func (uc *TransferTaxFeeUsecase) formatCurrency(amount float64) string {
	// Làm tròn đến hàng đơn vị
	rounded := math.Round(amount)

	// Chuyển sang int64 để format
	intAmount := int64(rounded)

	// Format với dấu phẩy
	if intAmount < 1000 {
		return fmt.Sprintf("%d", intAmount)
	}

	result := ""
	remainder := intAmount
	count := 0

	for remainder > 0 {
		if count > 0 && count%3 == 0 {
			result = "," + result
		}
		digit := remainder % 10
		result = fmt.Sprintf("%d", digit) + result
		remainder = remainder / 10
		count++
	}

	return result
}
