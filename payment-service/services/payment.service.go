package services

import (
	_db "common/db"
	_jwt "common/jwt"
	"net/http"
	"payment/models"
	"payment/repo"

	"github.com/gin-gonic/gin"
	"github.com/payOSHQ/payos-lib-golang"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type PaymentService struct {
	DB *gorm.DB
	// PayOS           PayOSClient
	RedirectSuccess string
	RedirectFail    string
	SuccessCode     string
	PaymentRepo     *repo.PaymentRepo
}

// type PayOSClient interface {
// 	// CreatePaymentLink(data *models.Order) (*payos.CheckoutResponseDataType, error)
// 	GetPaymentLinkInformation(orderID uint) (*payos.PaymentLinkDataType, error)
// 	CancelPaymentLink(orderID uint, reason string) (*payos.PaymentLinkDataType, error)
// 	ConfirmWebhook(webhookURL string) (string, error)
// }

// type Order struct {
// 	ID          uint `gorm:"primaryKey"`
// 	Amount      float64
// 	Description string
// 	Status      string
// 	CreatedBy   uint
// 	CreatedAt   time.Time
// 	UpdatedAt   time.Time
// }

func NewPaymentService(repo *repo.PaymentRepo, db *gorm.DB) *PaymentService {
	return &PaymentService{
		PaymentRepo:     repo,
		DB:              db,
		RedirectSuccess: viper.GetString("payment.url.success-redirect"),
		RedirectFail:    viper.GetString("payment.url.fail-redirect"),
		SuccessCode:     "00",
		// PayOS:       *payos,
	}
}

// 1. Tạo link thanh toán
// @Summary Tạo link thanh toán
// @Description Tạo link thanh toán và lưu vào DB
// @Tags PayOS
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order body models.OrderRequest true "Order thông tin"
// @Router /payos/create [post]
func (s *PaymentService) CreatePaymentLink(c *gin.Context) {
	order := &models.Order{}

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		// return nil, err
		return
	}

	order.Status = "PENDING"
	if err := s.PaymentRepo.Save(c, &order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	paymentData := payos.CheckoutRequestType{
		OrderCode:   int64(order.ID),
		Amount:      order.Amount,
		Description: order.Description,
		ReturnUrl:   s.RedirectSuccess,
		CancelUrl:   s.RedirectFail,
	}

	response, err := payos.CreatePaymentLink(paymentData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	mapCheckoutToEntity(response, order)
	if err := _db.SaveWithAudit(c, &order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"code": 0, "data": response})
}

// 2. Lấy thông tin link thanh toán
// @Summary Lấy thông tin link thanh toán
// @Description Lấy thông tin link thanh toán từ PayOS
// @Tags PayOS
// @Param orderId path int true "ID của đơn hàng"
// @Security BearerAuth
// @Router /payos/info/{orderId} [get]
func (s *PaymentService) GetPaymentLinkInformation(c *gin.Context) {
	orderId := c.Param("orderId")

	data, _ := payos.GetPaymentLinkInformation(orderId)

	c.JSON(http.StatusAccepted, gin.H{"data": data})
}

// 3. Hủy link thanh toán
// @Summary Hủy link thanh toán
// @Description Hủy link thanh toán và cập nhật DB
// @Tags PayOS
// @Param orderId path int true "ID của đơn hàng"
// @Param reason query string false "Lý do hủy"
// @Security BearerAuth
// @Router /payos/cancel/{orderId} [post]
func (s *PaymentService) CancelPaymentLink(c *gin.Context) {
	reason, _ := c.GetQuery("reason")
	orderId := c.Param("orderId")
	userId := _jwt.GetProfileId(c)

	// var webhookData payos.WebhookType
	response, err := payos.CancelPaymentLink(orderId, &reason)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 400, "message": err.Error()})
		return
	}

	var order models.Order
	if err := s.PaymentRepo.DB.Where("id = ? AND created_by = ?", orderId, userId).First(&order).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Đơn hàng không tồn tại hoặc không phải của bạn"})
		return
	}

	if order.Status == "PENDING" {
		order.Status = "CANCELLED"
		order.Reason = reason
		_db.SaveWithAudit(c, &order)
	} else {
		// return nil, errors.New("order has already been processed")
		c.JSON(http.StatusBadRequest, gin.H{"message": "Đơn đã được xử lý"})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"code": 0, "data": response})
	// return response, nil
}

func (s *PaymentService) ConfirmWebhook(webhookURL string) (string, error) {
	return payos.ConfirmWebhook(webhookURL)
}

// 5. Xác minh Webhook từ PayOS
// @Summary Xác minh Webhook
// @Description Xác minh Webhook sau khi thanh toán thành công
// @Tags PayOS
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /payos/webhook/verify [post]
func (s *PaymentService) VerifyPaymentWebhookData(c *gin.Context) {
	var webhookData payos.WebhookType

	if err := c.ShouldBindJSON(&webhookData); err != nil {
		c.JSON(http.StatusAccepted, gin.H{"error": "Invalid request body"})
		// return nil, err
		return
	}

	var order models.Order
	if err := s.PaymentRepo.DB.Where("id = ?", webhookData.Data.OrderCode).First(&order).Error; err != nil {
		c.JSON(http.StatusAccepted, gin.H{"error": err.Error})
		return
		// return nil, errors.New("order not found")
	}

	status := "FAILED"
	if webhookData.Code == s.SuccessCode {
		status = "COMPLETED"
	}
	order.Status = status

	order.MapWebhookToOrder(webhookData)
	s.PaymentRepo.Save(c, &order)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": order})
	// return webhookData.Data, nil
}

func mapCheckoutToEntity(response *payos.CheckoutResponseDataType, order *models.Order) {
	order.Amount = response.Amount
	order.Description = response.Description
	order.Status = response.Status
	order.ID = uint64(response.OrderCode)
	order.AccountName = response.AccountName
	order.AccountNumber = response.AccountNumber
	order.Bin = response.Bin
	order.OrderCode = uint(response.OrderCode)
	order.Currency = response.Currency
	order.PaymentLinkID = response.PaymentLinkId
	order.CheckoutURL = response.CheckoutUrl
	order.QRCode = response.QRCode
}
