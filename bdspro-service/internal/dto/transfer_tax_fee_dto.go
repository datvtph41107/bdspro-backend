package dto

// CalculateTransferTaxFeeRequest DTO cho request tính thuế và phí chuyển nhượng
type CalculateTransferTaxFeeRequest struct {
	TransferValue float64  `json:"transferValue" binding:"required,gt=0"` // Giá trị chuyển nhượng (VND)
	IsBusiness    bool     `json:"isBusiness"`                            // true nếu là doanh nghiệp, false nếu là cá nhân
	OriginalPrice *float64 `json:"originalPrice,omitempty"`               // Giá gốc mua vào (để tính thuế TNCN nếu có)
	ProvinceID    *uint64  `json:"provinceId,omitempty"`                  // ID tỉnh/thành phố (để tính lệ phí trước bạ theo từng địa phương)
}

// CalculateTransferTaxFeeResponse DTO cho response tính thuế và phí chuyển nhượng
type CalculateTransferTaxFeeResponse struct {
	TransferValue     float64         `json:"transferValue"`     // Giá trị chuyển nhượng
	PersonalIncomeTax float64         `json:"personalIncomeTax"` // Thuế thu nhập cá nhân (TNCN) - 2% nếu là cá nhân
	CorporateTax      float64         `json:"corporateTax"`      // Thuế thu nhập doanh nghiệp - 2% nếu là doanh nghiệp
	RegistrationFee   float64         `json:"registrationFee"`   // Lệ phí trước bạ - 0.5%
	NotarizationFee   float64         `json:"notarizationFee"`   // Phí công chứng
	DossierReviewFee  float64         `json:"dossierReviewFee"`  // Phí thẩm định hồ sơ
	CadastralFee      float64         `json:"cadastralFee"`      // Lệ phí địa chính
	CertificateFee    float64         `json:"certificateFee"`    // Lệ phí cấp GCN (Giấy chứng nhận)
	SurveyFee         float64         `json:"surveyFee"`         // Phí đo vẽ/địa chính
	TotalTaxFee       float64         `json:"totalTaxFee"`       // Tổng thuế và phí
	NetAmount         float64         `json:"netAmount"`         // Số tiền thực nhận sau khi trừ thuế phí
	Breakdown         TaxFeeBreakdown `json:"breakdown"`         // Chi tiết từng khoản
}

// TaxFeeBreakdown Chi tiết từng khoản thuế và phí
type TaxFeeBreakdown struct {
	PersonalIncomeTaxNote string `json:"personalIncomeTaxNote"` // Ghi chú thuế TNCN
	CorporateTaxNote      string `json:"corporateTaxNote"`      // Ghi chú thuế doanh nghiệp
	RegistrationFeeNote   string `json:"registrationFeeNote"`   // Ghi chú lệ phí trước bạ
	NotarizationFeeNote   string `json:"notarizationFeeNote"`   // Ghi chú phí công chứng
	DossierReviewFeeNote  string `json:"dossierReviewFeeNote"`  // Ghi chú phí thẩm định hồ sơ
	CadastralFeeNote      string `json:"cadastralFeeNote"`      // Ghi chú lệ phí địa chính
	CertificateFeeNote    string `json:"certificateFeeNote"`    // Ghi chú lệ phí cấp GCN
	SurveyFeeNote         string `json:"surveyFeeNote"`         // Ghi chú phí đo vẽ/địa chính
}
