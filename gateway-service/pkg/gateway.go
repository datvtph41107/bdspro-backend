package pkg

import (
	"sync"

	"github.com/go-resty/resty/v2"
)

func NewGateway(eurekaURL string) *Gateway {
	return &Gateway{
		EurekaURL: eurekaURL,
		Client:    resty.New(), // Dùng 1 client duy nhất
		Services:  sync.Map{},
	}
}
