package wallethandler

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"math"
	commercedomain "payment/internal/domain/payment"
	settlementuc "payment/internal/usecase/settlement"
	paymentpb "pb/types/payment"
	sharepb "pb/types/shared"
	"regexp"
	"strconv"
	"strings"
	"time"

	"payment/internal/job"
	walletuc "payment/internal/usecase/wallet"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type PaymentHandler struct {
	paymentpb.UnimplementedPaymentServiceServer
	paymentUsecase           walletuc.PaymentUsecase
	PaymentDashboardStatsJob *job.PaymentDashboardStatsJob
	commerceSettlement       *settlementuc.Service
	sepayAPIKey              string
}

func NewPaymentHandler(
	paymentUsecase walletuc.PaymentUsecase,
	paymentDashboardStatsJob *job.PaymentDashboardStatsJob,
	commerceSettlement *settlementuc.Service,
	sepayAPIKey string,
) *PaymentHandler {
	return &PaymentHandler{
		paymentUsecase:           paymentUsecase,
		PaymentDashboardStatsJob: paymentDashboardStatsJob,
		commerceSettlement:       commerceSettlement,
		sepayAPIKey:              strings.TrimSpace(sepayAPIKey),
	}
}

// @Summary Lấy dashboard ví
// @Description Lấy dashboard ví
// @Tags Payment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param empty body paymentpb.Empty true "Empty request"
// @Success 200 {object} paymentpb.WalletDashboard "Wallet dashboard"
// @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// @Failure 500 {object} status.Error "Lỗi server"
// @Router /wallets/dashboard [get]
func (h *PaymentHandler) GetWalletDashboard(ctx context.Context, req *paymentpb.Empty) (*paymentpb.WalletDashboard, error) {
	return h.paymentUsecase.GetWalletDashboard(ctx)
}

// @Summary Lấy danh sách giao dịch ví
// @Description Lấy danh sách giao dịch ví
// @Tags Payment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param getWalletTransactionsRequest body paymentpb.GetWalletTransactionsRequest true "Request body"
// @Success 200 {object} paymentpb.WalletTransactionList "Danh sách giao dịch ví"
// @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// @Failure 500 {object} status.Error "Lỗi server"
// @Router /wallets/transactions [get]
func (h *PaymentHandler) GetWalletTransactions(ctx context.Context, req *paymentpb.GetWalletTransactionsRequest) (*paymentpb.WalletTransactionList, error) {
	return h.paymentUsecase.GetWalletTransactions(ctx, req)
}

func (h *PaymentHandler) GetWalletTransactionsByWalletId(ctx context.Context, req *paymentpb.GetWalletTransactionsByWalletIdRequest) (*paymentpb.WalletTransactionList, error) {
	return h.paymentUsecase.GetWalletTransactionsByWalletId(ctx, req)
}

func (h *PaymentHandler) TransferBetweenWallets(ctx context.Context, req *paymentpb.TransferBetweenWalletsRequest) (*paymentpb.TransferBetweenWalletsResponse, error) {
	return h.paymentUsecase.TransferBetweenWallets(ctx, req)
}

// @Summary Nạp tiền
// @Description Nạp tiền
// @Tags Payment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param depositRequest body paymentpb.DepositRequest true "Request body"
// @Success 200 {object} paymentpb.WalletTransaction "Giao dịch nạp tiền"
// @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// @Failure 500 {object} status.Error "Lỗi server"
// @Router /wallets/deposit [post]
func (h *PaymentHandler) Deposit(ctx context.Context, req *paymentpb.DepositRequest) (*paymentpb.WalletTransaction, error) {
	return h.paymentUsecase.Deposit(ctx, req)
}

// func (h *PaymentHandler) Withdraw(ctx context.Context, req *paymentpb.WithdrawRequest) (*paymentpb.WithdrawalRequest, error) {
// 	return h.paymentUsecase.Withdraw(ctx, req)
// }

// @Summary Tạo giao dịch
// @Description Tạo giao dịch
// @Tags Payment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param paymentRequest body paymentpb.PaymentRequest true "Request body"
// @Success 200 {object} paymentpb.WalletTransaction "Giao dịch"
// @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// @Failure 500 {object} status.Error "Lỗi server"
// @Router /wallets/payment [post]
func (h *PaymentHandler) MakePayment(ctx context.Context, req *paymentpb.PaymentRequest) (*paymentpb.WalletTransaction, error) {
	return h.paymentUsecase.MakePayment(ctx, req)
}

// @Summary Lấy danh sách phương thức thanh toán
// @Description Lấy danh sách phương thức thanh toán
// @Tags Payment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param getPaymentMethodsRequest body paymentpb.GetPaymentMethodsRequest true "Request body"
// @Success 200 {object} paymentpb.PaymentMethodList "Danh sách phương thức thanh toán"
// @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// @Failure 500 {object} status.Error "Lỗi server"
// @Router /payment-methods [get]
func (h *PaymentHandler) GetPaymentMethods(ctx context.Context, req *paymentpb.GetPaymentMethodsRequest) (*paymentpb.PaymentMethodList, error) {
	return h.paymentUsecase.GetPaymentMethods(ctx, req)
}

// @Summary Tạo phương thức thanh toán
// @Description Tạo phương thức thanh toán
// @Tags Payment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param paymentMethod body paymentpb.PaymentMethod true "Phương thức thanh toán"
// @Success 200 {object} paymentpb.PaymentMethod "Phương thức thanh toán"
// @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// @Failure 500 {object} status.Error "Lỗi server"
// @Router /payment-methods [post]
func (h *PaymentHandler) CreatePaymentMethod(ctx context.Context, req *paymentpb.PaymentMethod) (*paymentpb.PaymentMethod, error) {
	return h.paymentUsecase.CreatePaymentMethod(ctx, req)
}

// @Summary Cập nhật phương thức thanh toán
// @Description Cập nhật phương thức thanh toán
// @Tags Payment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param paymentMethod body paymentpb.PaymentMethod true "Phương thức thanh toán"
// @Success 200 {object} paymentpb.PaymentMethod "Phương thức thanh toán"
// @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// @Failure 500 {object} status.Error "Lỗi server"
// @Router /payment-methods [put]
func (h *PaymentHandler) UpdatePaymentMethod(ctx context.Context, req *paymentpb.PaymentMethod) (*paymentpb.PaymentMethod, error) {
	return h.paymentUsecase.UpdatePaymentMethod(ctx, req)
}

// @Summary Xóa phương thức thanh toán
// @Description Xóa phương thức thanh toán
// @Tags Payment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param deletePaymentMethodRequest body paymentpb.DeletePaymentMethodRequest true "Request body"
// @Success 200 {object} paymentpb.Empty "Phương thức thanh toán"
// @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// @Failure 500 {object} status.Error "Lỗi server"
// @Router /payment-methods [delete]
func (h *PaymentHandler) DeletePaymentMethod(ctx context.Context, req *paymentpb.DeletePaymentMethodRequest) (*paymentpb.Empty, error) {
	return h.paymentUsecase.DeletePaymentMethod(ctx, req)
}

// @Summary Lấy báo cáo ví
// @Description Lấy báo cáo ví
// @Tags Payment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param walletId path string true "ID ví"
// @Param reportRequest body paymentpb.ReportRequest true "Request body"
// @Success 200 {object} paymentpb.ReportResponse "Báo cáo ví"
// @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// @Failure 500 {object} status.Error "Lỗi server"
// @Router /wallets/{walletId}/report [get]
func (h *PaymentHandler) GetWalletReport(ctx context.Context, req *paymentpb.ReportRequest) (*paymentpb.ReportResponse, error) {
	return h.paymentUsecase.GetWalletReport(ctx, req)
}

// @Summary Xuất dữ liệu ví
// @Description Xuất dữ liệu ví
// @Tags Payment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param walletId path string true "ID ví"
// @Param reportRequest body paymentpb.ReportRequest true "Request body"
// @Success 200 {object} paymentpb.ExportResponse "Dữ liệu ví"
// @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// @Failure 500 {object} status.Error "Lỗi server"
// @Router /wallets/{walletId}/export [get]
func (h *PaymentHandler) ExportWalletData(ctx context.Context, req *paymentpb.ReportRequest) (*paymentpb.ExportResponse, error) {
	return h.paymentUsecase.ExportWalletData(ctx, req)
}

// @Summary Lấy danh sách loại giao dịch
// @Description Lấy danh sách loại giao dịch
// @Tags Payment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param empty body paymentpb.Empty true "Empty request"
// @Success 200 {object} paymentpb.TransactionTypeList "Danh sách loại giao dịch"
// @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// @Failure 500 {object} status.Error "Lỗi server"
// @Router /transaction-types [get]
func (h *PaymentHandler) GetTransactionTypes(ctx context.Context, req *paymentpb.Empty) (*paymentpb.TransactionTypeList, error) {
	return h.paymentUsecase.GetTransactionTypes(ctx, req)
}

// @Summary Tạo loại giao dịch
// @Description Tạo loại giao dịch
// @Tags Payment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param transactionType body paymentpb.TransactionType true "Loại giao dịch"
// @Success 200 {object} paymentpb.TransactionType "Loại giao dịch"
// @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// @Failure 500 {object} status.Error "Lỗi server"
// @Router /transaction-types [post]
func (h *PaymentHandler) CreateTransactionType(ctx context.Context, req *paymentpb.TransactionType) (*paymentpb.TransactionType, error) {
	return h.paymentUsecase.CreateTransactionType(ctx, req)
}

// @Summary Cập nhật loại giao dịch
// @Description Cập nhật loại giao dịch
// @Tags Payment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param transactionType body paymentpb.TransactionType true "Loại giao dịch"
// @Success 200 {object} paymentpb.TransactionType "Loại giao dịch"
// @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// @Failure 500 {object} status.Error "Lỗi server"
// @Router /transaction-types [put]
func (h *PaymentHandler) UpdateTransactionType(ctx context.Context, req *paymentpb.TransactionType) (*paymentpb.TransactionType, error) {
	return h.paymentUsecase.UpdateTransactionType(ctx, req)
}

// @Summary Xóa loại giao dịch
// @Description Xóa loại giao dịch
// @Tags Payment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param deleteTransactionTypeRequest body paymentpb.DeleteTransactionTypeRequest true "Request body"
// @Success 200 {object} paymentpb.Empty "Loại giao dịch"
// @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// @Failure 500 {object} status.Error "Lỗi server"
// @Router /transaction-types [delete]
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

// @Summary Webhook thanh toán Sepay
// @Description Webhook thanh toán Sepay
// @Tags Payment
// @Accept json
// @Produce json
// @Param sepayWebhookRequest body paymentpb.SepayWebhookRequest true "Request body"
// @Success 200 {object} paymentpb.Empty "Webhook thanh toán Sepay"
// @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// @Failure 500 {object} status.Error "Lỗi server"
// @Router /sepay/webhook [post]
func (h *PaymentHandler) SepayWebhook(ctx context.Context, req *paymentpb.SepayWebhookRequest) (*paymentpb.Empty, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}
	digestValues := md.Get("x-qhpro-provider-credential-sha256")
	if len(digestValues) != 1 || h.sepayAPIKey == "" {
		return nil, status.Error(codes.Unauthenticated, "invalid provider credential")
	}
	expectedRaw := sha256.Sum256([]byte(h.sepayAPIKey))
	expected := hex.EncodeToString(expectedRaw[:])
	if subtle.ConstantTimeCompare([]byte(digestValues[0]), []byte(expected)) != 1 {
		return nil, status.Error(codes.Unauthenticated, "invalid provider credential")
	}
	if h.commerceSettlement != nil {
		if ref, ok := ExtractCommercialOrderReference(req.GetContent(), req.GetReferenceCode()); ok {
			if strings.ToLower(strings.TrimSpace(req.GetTransferType())) != "in" {
				return nil, status.Error(codes.InvalidArgument, "commercial webhook must be incoming transfer")
			}
			amount := req.GetTransferAmount()
			if math.IsNaN(amount) || math.IsInf(amount, 0) || amount <= 0 || amount != math.Trunc(amount) {
				return nil, status.Error(codes.InvalidArgument, "invalid VND amount")
			}
			occurred, err := time.ParseInLocation("2006-01-02 15:04:05", req.GetTransactionDate(), time.FixedZone("Asia/Ho_Chi_Minh", 7*3600))
			if err != nil {
				return nil, status.Error(codes.InvalidArgument, "invalid transaction date")
			}
			_, err = h.commerceSettlement.RecordProviderEvent(ctx, commercedomain.ProviderEvent{Provider: "sepay", TransactionID: strconv.FormatUint(uint64(req.GetId()), 10), Reference: ref, Amount: commercedomain.Money{Currency: commercedomain.CurrencyVND, AmountMinor: int64(amount)}, OccurredAt: occurred, Type: commercedomain.ProviderFundsConfirmed, EvidenceHash: fmt.Sprintf("sepay:%d", req.GetId())})
			if err != nil {
				return nil, status.Error(codes.FailedPrecondition, err.Error())
			}
			return &paymentpb.Empty{}, nil
		}
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

// @Summary Lấy thống kê số lượng giao dịch đang xử lý
// @Description Lấy thống kê số lượng giao dịch đang xử lý theo 2 ngày
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param firstDate query string true "Ngày đầu tiên (YYYY-MM-DD)"
// @Param secondDate query string true "Ngày thứ hai (YYYY-MM-DD)"
// @Security BearerAuth
// @Success 200 {object} sharepb.GetStatsResponse "Thành công"
// @Router /v2/payment/stats [get]
func (h *PaymentHandler) GetPaymentStats(ctx context.Context, req *sharepb.GetStatsRequest) (*sharepb.GetStatsResponse, error) {
	firstDate, err := time.Parse(time.DateOnly, req.FirstDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid firstDate format")
	}

	secondDate, err := time.Parse(time.DateOnly, req.SecondDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid secondDate format")
	}

	stats, stats2, err := h.PaymentDashboardStatsJob.GetStats(ctx, firstDate, secondDate)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &sharepb.GetStatsResponse{
		FirstDate:       stats.CalculateTime.Format(time.DateOnly),
		SecondDate:      stats2.CalculateTime.Format(time.DateOnly),
		FirstDateCount:  uint32(stats.Count),
		SecondDateCount: uint32(stats2.Count),
	}, nil
}

func (h *PaymentHandler) GetDashboardMetrics(ctx context.Context, req *paymentpb.GetDashboardMetricsRequest) (*paymentpb.GetDashboardMetricsResponse, error) {
	fromDate, err := time.Parse(time.DateOnly, req.FromDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid fromDate format")
	}
	toDate, err := time.Parse(time.DateOnly, req.ToDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Invalid toDate format")
	}
	metrics, err := h.paymentUsecase.GetDashboardMetrics(ctx, fromDate, toDate)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &paymentpb.GetDashboardMetricsResponse{
		Data: metrics,
	}, nil
}

func ExtractCommercialOrderReference(values ...string) (string, bool) {
	re := regexp.MustCompile(`(?i)QHP-[0-9A-F]{16}`)
	for _, v := range values {
		if x := re.FindString(v); x != "" {
			return strings.ToUpper(x), true
		}
	}
	return "", false
}
