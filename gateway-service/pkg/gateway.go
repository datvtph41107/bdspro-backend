package pkg

import (
	"sync"

	"github.com/go-resty/resty/v2"
)

// NewGateway tạo instance mới của Gateway
func NewGateway(eurekaURL string) *Gateway {
	return &Gateway{
		EurekaURL: eurekaURL,
		Client:    resty.New(),
		Services:  sync.Map{},
	}
}
