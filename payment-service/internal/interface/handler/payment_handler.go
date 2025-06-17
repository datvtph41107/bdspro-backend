package handler

import (
	paymentpb "pb/types/payment"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


type PaymentHandler struct {
	 paymentpb.UnimplementedPaymentServiceServer
}

func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{}
}

func (h *PaymentHandler) GetWalletDashboard(ctx context.Context, req *paymentpb.Wallet) (*paymentpb.WalletDashboard, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetWalletDashboard not implemented")
}

func (h *PaymentHandler) GetWalletTransactions(ctx context.Context, req *paymentpb.GetWalletTransactionsRequest) (*paymentpb.WalletTransactionList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetWalletTransactions not implemented")
}

func (h *PaymentHandler) Deposit(ctx context.Context, req *paymentpb.DepositRequest) (*paymentpb.WalletTransaction, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Deposit not implemented")
}

func (h *PaymentHandler) Withdraw(ctx context.Context, req *paymentpb.WithdrawRequest) (*paymentpb.WithdrawalRequest, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Withdraw not implemented")
}

func (h *PaymentHandler) MakePayment(ctx context.Context, req *paymentpb.PaymentRequest) (*paymentpb.WalletTransaction, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MakePayment not implemented")
}

func (h *PaymentHandler) GetPaymentMethods(ctx context.Context, req *paymentpb.GetPaymentMethodsRequest) (*paymentpb.PaymentMethodList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPaymentMethods not implemented")
}

func (h *PaymentHandler) CreatePaymentMethod(ctx context.Context, req *paymentpb.PaymentMethod) (*paymentpb.PaymentMethod, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePaymentMethod not implemented")
}

func (h *PaymentHandler) UpdatePaymentMethod(ctx context.Context, req *paymentpb.PaymentMethod) (*paymentpb.PaymentMethod, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePaymentMethod not implemented")
}

func (h *PaymentHandler) DeletePaymentMethod(ctx context.Context, req *paymentpb.DeletePaymentMethodRequest) (*paymentpb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePaymentMethod not implemented")
}

func (h *PaymentHandler) GetWalletReport(ctx context.Context, req *paymentpb.ReportRequest) (*paymentpb.ReportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetWalletReport not implemented")
}

func (h *PaymentHandler) ExportWalletData(ctx context.Context, req *paymentpb.ReportRequest) (*paymentpb.ExportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExportWalletData not implemented")
}

func (h *PaymentHandler) GetDealTransactionsDashboard(ctx context.Context, req *paymentpb.Empty) (*paymentpb.DealTransactionsDashboard, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDealTransactionsDashboard not implemented")
}

func (h *PaymentHandler) GetDealTransactions(ctx context.Context, req *paymentpb.GetDealTransactionsRequest) (*paymentpb.DealTransactionList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDealTransactions not implemented")
}

func (h *PaymentHandler) GetDealTransaction(ctx context.Context, req *paymentpb.GetDealTransactionRequest) (*paymentpb.DealTransaction, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDealTransaction not implemented")
}

func (h *PaymentHandler) CreateDealTransaction(ctx context.Context, req *paymentpb.DealTransaction) (*paymentpb.DealTransaction, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateDealTransaction not implemented")
}

func (h *PaymentHandler) UpdateDealTransaction(ctx context.Context, req *paymentpb.DealTransaction) (*paymentpb.DealTransaction, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateDealTransaction not implemented")
}

func (h *PaymentHandler) ApproveDealTransaction(ctx context.Context, req *paymentpb.ApproveDealTransactionRequest) (*paymentpb.DealTransaction, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ApproveDealTransaction not implemented")
}

func (h *PaymentHandler) GetTransactionTypes(ctx context.Context, req *paymentpb.Empty) (*paymentpb.TransactionTypeList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTransactionTypes not implemented")
}

func (h *PaymentHandler) CreateTransactionType(ctx context.Context, req *paymentpb.TransactionType) (*paymentpb.TransactionType, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTransactionType not implemented")
}

func (h *PaymentHandler) UpdateTransactionType(ctx context.Context, req *paymentpb.TransactionType) (*paymentpb.TransactionType, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTransactionType not implemented")
}

func (h *PaymentHandler) DeleteTransactionType(ctx context.Context, req *paymentpb.DeleteTransactionTypeRequest) (*paymentpb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTransactionType not implemented")
}