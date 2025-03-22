package services

import (
	"fmt"

	payment "payment/models"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

type PaymentClient struct {
	Client *resty.Client
	Url    string
	Path   string
}

func NewPaymentClient() *PaymentClient {
	return &PaymentClient{
		Client: resty.New(),
		Url:    viper.GetString("payment.url"),
		Path:   viper.GetString("payment.path"),
	}
}

func (p *PaymentClient) RequestOrder(c *gin.Context, body payment.OrderRequest) (any, error) {
	targetURL := fmt.Sprintf("http://%s/%s", p.Url, p.Path)
	resp, err := p.Client.R().
		SetHeaderMultiValues(c.Request.Header).
		SetBody(body).
		Execute(c.Request.Method, targetURL)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
