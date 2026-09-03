package bankhandler

import (
	"context"
	bankuc "payment/internal/usecase/bank"
	paymentpb "pb/types/payment"
	sharepb "pb/types/shared"
)

type InternalHandler struct {
	paymentpb.UnimplementedInternalServiceServer
	BankUsecase *bankuc.BankUsecase
}

func NewInternalHandler(bankUsecase *bankuc.BankUsecase) *InternalHandler {
	return &InternalHandler{
		BankUsecase: bankUsecase,
	}
}

func (h *InternalHandler) GetBankById(ctx context.Context, req *sharepb.IdRequest) (*paymentpb.Bank, error) {
	bank, err := h.BankUsecase.GetByID(ctx, uint32(req.Id))
	if err != nil {
		return nil, err
	}
	return &paymentpb.Bank{
		Id:          bank.ID,
		Name:        bank.Name,
		Code:        bank.Code,
		Logo:        bank.Logo,
		Active:      bank.Active,
		Description: bank.Description,
	}, nil
}
