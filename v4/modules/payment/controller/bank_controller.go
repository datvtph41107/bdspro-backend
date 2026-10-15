package controller

import (
	"context"
	"net/http"

	"userv4/internal/httpapi"
	"userv4/modules/payment/dto"
	"userv4/modules/payment/services"
)

// BankController owns GET /v2/payment/bank/all (gateway parity).
type BankController struct {
	banks services.BankService
}

func NewBankController(banks services.BankService) *BankController {
	return &BankController{banks: banks}
}

func (c *BankController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/payment/bank/all", c.listAll)
}

func (c *BankController) listAll(ctx context.Context, req dto.BankListRequest) (*httpapi.Result, error) {
	result, err := c.banks.ListAll(ctx, req.Timestamp)
	if err != nil {
		return nil, err
	}
	return bankListResult(result)
}
