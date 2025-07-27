package handler

import (
	"context"
	paymentpb "pb/types/payment"
	"regexp"

	"github.com/tuandq2112/payment-service/appconfig"
	"github.com/tuandq2112/payment-service/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type PaymentHandler struct {
	paymentpb.UnimplementedPaymentServiceServer
	paymentUsecase usecase.PaymentUsecase
}

func NewPaymentHandler(paymentUsecase usecase.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{
		paymentUsecase: paymentUsecase,
	}
}

func (h *PaymentHandler) GetWalletDashboard(ctx context.Context, req *paymentpb.Empty) (*paymentpb.WalletDashboard, error) {
	return h.paymentUsecase.GetWalletDashboard(ctx)
}

func (h *PaymentHandler) GetWalletTransactions(ctx context.Context, req *paymentpb.GetWalletTransactionsRequest) (*paymentpb.WalletTransactionList, error) {
	return h.paymentUsecase.GetWalletTransactions(ctx, req)
}

func (h *PaymentHandler) GetWalletTransactionsByWalletId(ctx context.Context, req *paymentpb.GetWalletTransactionsByWalletIdRequest) (*paymentpb.WalletTransactionList, error) {
	return h.paymentUsecase.GetWalletTransactionsByWalletId(ctx, req)
}

func (h *PaymentHandler) TransferBetweenWallets(ctx context.Context, req *paymentpb.TransferBetweenWalletsRequest) (*paymentpb.TransferBetweenWalletsResponse, error) {
	return h.paymentUsecase.TransferBetweenWallets(ctx, req)
}

func (h *PaymentHandler) Deposit(ctx context.Context, req *paymentpb.DepositRequest) (*paymentpb.WalletTransaction, error) {
	return h.paymentUsecase.Deposit(ctx, req)
}

// func (h *PaymentHandler) Withdraw(ctx context.Context, req *paymentpb.WithdrawRequest) (*paymentpb.WithdrawalRequest, error) {
// 	return h.paymentUsecase.Withdraw(ctx, req)
// }

func (h *PaymentHandler) MakePayment(ctx context.Context, req *paymentpb.PaymentRequest) (*paymentpb.WalletTransaction, error) {
	return h.paymentUsecase.MakePayment(ctx, req)
}

func (h *PaymentHandler) GetPaymentMethods(ctx context.Context, req *paymentpb.GetPaymentMethodsRequest) (*paymentpb.PaymentMethodList, error) {
	return h.paymentUsecase.GetPaymentMethods(ctx, req)
}

func (h *PaymentHandler) CreatePaymentMethod(ctx context.Context, req *paymentpb.PaymentMethod) (*paymentpb.PaymentMethod, error) {
	return h.paymentUsecase.CreatePaymentMethod(ctx, req)
}

func (h *PaymentHandler) UpdatePaymentMethod(ctx context.Context, req *paymentpb.PaymentMethod) (*paymentpb.PaymentMethod, error) {
	return h.paymentUsecase.UpdatePaymentMethod(ctx, req)
}

func (h *PaymentHandler) DeletePaymentMethod(ctx context.Context, req *paymentpb.DeletePaymentMethodRequest) (*paymentpb.Empty, error) {
	return h.paymentUsecase.DeletePaymentMethod(ctx, req)
}

func (h *PaymentHandler) GetWalletReport(ctx context.Context, req *paymentpb.ReportRequest) (*paymentpb.ReportResponse, error) {
	return h.paymentUsecase.GetWalletReport(ctx, req)
}

func (h *PaymentHandler) ExportWalletData(ctx context.Context, req *paymentpb.ReportRequest) (*paymentpb.ExportResponse, error) {
	return h.paymentUsecase.ExportWalletData(ctx, req)
}



func (h *PaymentHandler) GetTransactionTypes(ctx context.Context, req *paymentpb.Empty) (*paymentpb.TransactionTypeList, error) {
	return h.paymentUsecase.GetTransactionTypes(ctx, req)
}

func (h *PaymentHandler) CreateTransactionType(ctx context.Context, req *paymentpb.TransactionType) (*paymentpb.TransactionType, error) {
	return h.paymentUsecase.CreateTransactionType(ctx, req)
}

func (h *PaymentHandler) UpdateTransactionType(ctx context.Context, req *paymentpb.TransactionType) (*paymentpb.TransactionType, error) {
	return h.paymentUsecase.UpdateTransactionType(ctx, req)
}

func (h *PaymentHandler) DeleteTransactionType(ctx context.Context, req *paymentpb.DeleteTransactionTypeRequest) (*paymentpb.Empty, error) {
	return h.paymentUsecase.DeleteTransactionType(ctx, req)
}

func (h *PaymentHandler) CreateWallet(ctx context.Context, req *paymentpb.CreateWalletRequest) (*paymentpb.CreateWalletResponse, error) {
	return h.paymentUsecase.CreateWallet(ctx, req)
}

func (h *PaymentHandler) GetWalletByUserId(ctx context.Context, req *paymentpb.GetWalletByUserIdRequest) (*paymentpb.Wallet, error) {
	return h.paymentUsecase.GetWalletByUserId(ctx, req)
}

func (h *PaymentHandler) GetWalletByWalletId(ctx context.Context, req *paymentpb.GetWalletByWalletIdRequest) (*paymentpb.Wallet, error) {
	return h.paymentUsecase.GetWalletByWalletId(ctx, req)
}

func (h *PaymentHandler) SepayWebhook(ctx context.Context, req *paymentpb.SepayWebhookRequest) (*paymentpb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}
	authHeader := md.Get("grpcgateway-authorization")
	
	if authHeader[0] != appconfig.SepayAPIKey {
		return nil, status.Error(codes.Unauthenticated, "invalid api key")
	}

	depositCode := ExtractDepositCode(req.Content)
	_, err := h.paymentUsecase.HandleDepositWebhook(ctx, depositCode, req.TransferAmount)
	if err != nil {
		return nil, err
	}
	return &paymentpb.Empty{}, nil
}

func ExtractDepositCode(content string) string {
    re := regexp.MustCompile(`DEP\d+`)
    match := re.FindString(content)
    return match
}