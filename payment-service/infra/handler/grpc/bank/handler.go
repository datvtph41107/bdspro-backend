package bankhandler

import (
	"context"
	paymentpb "pb/types/payment"
	sharepb "pb/types/shared"

	bankuc "payment/internal/usecase/bank"
)

type BankHandler struct {
	paymentpb.UnimplementedBankServiceServer
	BankUsecase *bankuc.BankUsecase
}

func NewBankHandler(bankUsecase *bankuc.BankUsecase) *BankHandler {
	return &BankHandler{
		BankUsecase: bankUsecase,
	}
}

// @Summary Lấy danh sách ngân hàng
// @Description Lấy danh sách ngân hàng
// @Tags Bank
// @Accept json
// @Produce json
// @Success 200 {object} paymentpb.BankList
// @Router /bank/all [get]
func (h *BankHandler) GetAll(ctx context.Context, req *sharepb.Empty) (*paymentpb.BankList, error) {
	banks, err := h.BankUsecase.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	bankList := make([]*paymentpb.Bank, len(banks))
	for i, bank := range banks {
		bankList[i] = &paymentpb.Bank{
			Id:          bank.ID,
			Name:        bank.Name,
			Code:        bank.Code,
			Logo:        bank.Logo,
			Active:      bank.Active,
			Description: bank.Description,
		}
	}
	return &paymentpb.BankList{
		Data: bankList,
	}, nil
}
