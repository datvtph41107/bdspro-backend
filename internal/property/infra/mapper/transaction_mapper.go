package mapper

import (
	"bdspro/internal/dto"
	bdspropb "pb/types/bdspro"
)

// MapCalculateTransferTaxFeeRequestToDTO Convert proto request to DTO
func MapCalculateTransferTaxFeeRequestToDTO(req *bdspropb.CalculateTransferTaxFeeRequest) *dto.CalculateTransferTaxFeeRequest {
	dtoReq := &dto.CalculateTransferTaxFeeRequest{
		TransferValue: req.TransferValue,
		IsBusiness:    req.IsBusiness,
	}

	if req.OriginalPrice != nil {
		originalPrice := *req.OriginalPrice
		dtoReq.OriginalPrice = &originalPrice
	}

	if req.ProvinceId != nil {
		provinceID := *req.ProvinceId
		dtoReq.ProvinceID = &provinceID
	}

	return dtoReq
}

// MapCalculateTransferTaxFeeResponseToProto Convert DTO response to proto
func MapCalculateTransferTaxFeeResponseToProto(resp *dto.CalculateTransferTaxFeeResponse) *bdspropb.CalculateTransferTaxFeeResponse {
	return &bdspropb.CalculateTransferTaxFeeResponse{
		TransferValue:     resp.TransferValue,
		PersonalIncomeTax: resp.PersonalIncomeTax,
		CorporateTax:      resp.CorporateTax,
		RegistrationFee:   resp.RegistrationFee,
		NotarizationFee:   resp.NotarizationFee,
		DossierReviewFee:  resp.DossierReviewFee,
		CadastralFee:      resp.CadastralFee,
		CertificateFee:    resp.CertificateFee,
		SurveyFee:         resp.SurveyFee,
		TotalTaxFee:       resp.TotalTaxFee,
		NetAmount:         resp.NetAmount,
		Breakdown: &bdspropb.TaxFeeBreakdown{
			PersonalIncomeTaxNote: resp.Breakdown.PersonalIncomeTaxNote,
			CorporateTaxNote:      resp.Breakdown.CorporateTaxNote,
			RegistrationFeeNote:   resp.Breakdown.RegistrationFeeNote,
			NotarizationFeeNote:   resp.Breakdown.NotarizationFeeNote,
			DossierReviewFeeNote:  resp.Breakdown.DossierReviewFeeNote,
			CadastralFeeNote:      resp.Breakdown.CadastralFeeNote,
			CertificateFeeNote:    resp.Breakdown.CertificateFeeNote,
			SurveyFeeNote:         resp.Breakdown.SurveyFeeNote,
		},
	}
}

// MapCalculateLoanPaymentRequestToDTO Convert proto request to DTO
func MapCalculateLoanPaymentRequestToDTO(req *bdspropb.CalculateLoanPaymentRequest) *dto.CalculateLoanPaymentRequest {
	dtoReq := &dto.CalculateLoanPaymentRequest{
		InterestRate:  req.InterestRate,
		PropertyValue: req.PropertyValue,
		LoanAmount:    req.LoanAmount,
		Term:          req.Term,
		PaymentMethod: req.PaymentMethod,
	}

	if req.FloatingInterestRate != nil {
		floatingRate := *req.FloatingInterestRate
		dtoReq.FloatingInterestRate = &floatingRate
	}

	return dtoReq
}

// MapCalculateLoanPaymentResponseToProto Convert DTO response to proto
func MapCalculateLoanPaymentResponseToProto(resp *dto.CalculateLoanPaymentResponse) *bdspropb.CalculateLoanPaymentResponse {
	schedule := make([]*bdspropb.LoanPaymentScheduleItem, len(resp.Schedule))
	for i, item := range resp.Schedule {
		schedule[i] = &bdspropb.LoanPaymentScheduleItem{
			Month:              item.Month,
			RemainingPrincipal: item.RemainingPrincipal,
			PrincipalPayment:   item.PrincipalPayment,
			InterestPayment:    item.InterestPayment,
			TotalPayment:       item.TotalPayment,
			InterestRate:       item.InterestRate,
		}
	}

	return &bdspropb.CalculateLoanPaymentResponse{
		Schedule:      schedule,
		TotalInterest: resp.TotalInterest,
		TotalPayment:  resp.TotalPayment,
	}
}

// type TransactionTransformer interface {
// 	CreateTransactionRequestToEntity(req *bdspropb.CreateTransactionRequest) *tx_domain.Tx
// 	UpdateTransactionRequestToEntity(req *bdspropb.UpdateTransactionRequest) *tx_domain.Tx
// 	EntityToTransactionResponse(tx *tx_domain.Tx) *bdspropb.Transaction
// 	EntityToListTransactionsResponse(txs []*tx_domain.Tx, total int32) *bdspropb.ListTransactionsResponse
// 	EntityToTransactionTypesResponse(types []string) *bdspropb.GetTransactionTypesResponse
// }

// type transactionTransformer struct {
// 	DepositeMapper *DepositeMapper
// }

// func NewTransactionTransformer(
// 	depositeMapper *DepositeMapper,
// ) TransactionTransformer {
// 	return &transactionTransformer{
// 		DepositeMapper: depositeMapper,
// 	}
// }

// // func (t *transactionTransformer) GetStatusTransaction(e *domain.Transaction) uint32 {
// // 	if e.DepositeID != nil {
// // 		return 20
// // 	}
// // 	return 10
// // }

// func (t *transactionTransformer) CreateTransactionRequestToEntity(req *bdspropb.CreateTransactionRequest) *tx_domain.Tx {
// 	transactionDate, _ := time.Parse(time.RFC3339, req.TransactionDate)
// 	return &tx_domain.Tx{
// 		// OrganizationId:       req.OrganizationId,
// 		OwnerId:              uint64(req.OwnerId),
// 		Amount:               req.Amount,
// 		Currency:             req.Currency,
// 		TransactionName:      req.TransactionName,
// 		Description:          req.Description,
// 		TransactionType:      enums.TransactionType(req.TransactionType),
// 		CategoryId:           req.CategoryId,
// 		PaymentMethodId:      req.PaymentMethodId,
// 		RelatedDealId:        req.RelatedDealId,
// 		RelatedTransactionId: req.RelatedTransactionId,
// 		TransactionDate:      transactionDate,
// 		TransactionStatus:    enums.TransactionStatusDraft,
// 		ApprovalStatus:       enums.ApprovalStatusPending,
// 	}
// }

// func (t *transactionTransformer) UpdateTransactionRequestToEntity(req *bdspropb.UpdateTransactionRequest) *tx_domain.Tx {
// 	transactionDate, _ := time.Parse(time.RFC3339, req.TransactionDate)
// 	return &tx_domain.Tx{
// 		BaseEntity: _models.BaseEntity{
// 			ID: uint64(req.Id),
// 		},
// 		Amount:               req.Amount,
// 		Currency:             req.Currency,
// 		TransactionName:      req.TransactionName,
// 		Description:          req.Description,
// 		CategoryId:           req.CategoryId,
// 		PaymentMethodId:      req.PaymentMethodId,
// 		RelatedDealId:        req.RelatedDealId,
// 		RelatedTransactionId: req.RelatedTransactionId,
// 		TransactionDate:      transactionDate,
// 	}
// }

// func (t *transactionTransformer) EntityToTransactionResponse(tx *tx_domain.Tx) *bdspropb.Transaction {
// 	if tx == nil {
// 		return nil
// 	}
// 	var approvalDate *string
// 	if tx.ApprovalDate != nil && !tx.ApprovalDate.IsZero() {
// 		date := tx.ApprovalDate.Format(time.RFC3339)
// 		approvalDate = &date
// 	}

// 	result := &bdspropb.Transaction{
// 		Id:      tx.ID,
// 		OwnerId: uint32(tx.OwnerId),
// 		// OrganizationId:       tx.OrganizationId,
// 		Amount:               tx.Amount,
// 		Currency:             tx.Currency,
// 		TransactionName:      tx.TransactionName,
// 		Description:          tx.Description,
// 		TransactionType:      string(tx.TransactionType),
// 		CategoryId:           tx.CategoryId,
// 		PaymentMethodId:      tx.PaymentMethodId,
// 		RelatedDealId:        tx.RelatedDealId,
// 		RelatedTransactionId: tx.RelatedTransactionId,
// 		ApprovalStatus:       string(tx.ApprovalStatus),
// 		ApprovedBy:           tx.ApprovedBy,
// 		ApprovalDate:         approvalDate,
// 		TransactionStatus:    uint32(tx.TransactionStatus),
// 		TransactionDate:      tx.TransactionDate.Format(time.RFC3339),
// 		DepositeAmount:       tx.DepositeAmount,
// 		DepositeNote:         tx.DepositeNote,
// 		ContactId:            tx.ContactID,
// 	}

// 	return result
// }

// func (t *transactionTransformer) EntityToListTransactionsResponse(txs []*tx_domain.Tx, total int32) *bdspropb.ListTransactionsResponse {
// 	transactions := make([]*bdspropb.Transaction, 0, len(txs))
// 	for _, tx := range txs {
// 		transactions = append(transactions, t.EntityToTransactionResponse(tx))
// 	}

// 	return &bdspropb.ListTransactionsResponse{
// 		Data:  transactions,
// 		Total: total,
// 	}
// }

// func (t *transactionTransformer) EntityToTransactionTypesResponse(types []string) *bdspropb.GetTransactionTypesResponse {
// 	return &bdspropb.GetTransactionTypesResponse{
// 		Types: types,
// 	}
// }
