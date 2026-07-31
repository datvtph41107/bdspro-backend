package client

import (
	"context"
	"pb/clients"
	"user/env"
	"user/internal/dto"
	"user/internal/interface/providers"
)

type PaymentClient struct {
	*clients.PaymentClient
}

func NewPaymentClient() providers.IPaymentProvider {
	paymentClient := clients.MakePaymentClient(env.ENV_RUNTIME)
	return &PaymentClient{
		PaymentClient: paymentClient,
		// client:        resty.New(),
		// url:           viper.GetString("payment.url"),
		// path:          viper.GetString("payment.path"),
	}
}

func (p *PaymentClient) CreateOrder(ctx context.Context, req dto.OrderRequest) (*dto.OrderResponse, error) {
	// targetURL := fmt.Sprintf("http://%s/%s", p.url, p.path)

	// resp, err := p.client.R().
	// 	SetContext(ctx).
	// 	SetBody(req).
	// 	Post(targetURL)

	// if err != nil {
	// 	return nil, err
	// }

	// if resp.IsError() {
	// 	return nil, fmt.Errorf("payment service error: %s", resp.String())
	// }

	// var result dto.OrderResponse
	// if err := json.Unmarshal(resp.Body(), &result); err != nil {
	// 	return nil, err
	// }

	// return &result, nil
	return nil, nil
}
