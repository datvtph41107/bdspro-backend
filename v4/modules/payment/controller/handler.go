package controller

import (
	"net/http"

	"userv4/modules/payment/services"
)

// Handler is the modules/payment HTTP entry.
type Handler struct {
	bank *BankController
}

func NewHandler(banks services.BankService) *Handler {
	return &Handler{
		bank: NewBankController(banks),
	}
}

// RegisterHTTP mounts payment routes on the v4 HTTP server.
func (h *Handler) RegisterHTTP(mux *http.ServeMux) {
	h.bank.RegisterHTTP(mux)
}
