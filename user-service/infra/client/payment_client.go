package client

import (
	"context"
	"encoding/json"
	"fmt"
	"user/internal/dto"
	"user/internal/interface/providers"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

type PaymentClient struct {
	client *resty.Client
	url    string
	path   string
}

func NewPaymentClient() providers.IPaymentProvider {
	return &PaymentClient{
		client: resty.New(),
		url:    viper.GetString("payment.url"),
		path:   viper.GetString("payment.path"),
	}
}

func (p *PaymentClient) CreateOrder(ctx context.Context, req dto.OrderRequest) (*dto.OrderResponse, error) {
	targetURL := fmt.Sprintf("http://%s/%s", p.url, p.path)

	resp, err := p.client.R().
		SetContext(ctx).
		SetBody(req).
		Post(targetURL)

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("payment service error: %s", resp.String())
	}

	var result dto.OrderResponse
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, err
	}

	return &result, nil
}
