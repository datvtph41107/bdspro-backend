package mod

import (
	"sync"

	"github.com/go-resty/resty/v2"
)

type EurekaService struct {
	Name string `json:"name"`
	IP   string `json:"ipAddr"`
	Port int    `json:"port"`
}

type Gateway struct {
	EurekaURL string
	Client    *resty.Client
	Services  sync.Map
}

type EurekaResponse struct {
	Applications struct {
		Applications []struct {
			Name      string `json:"name"`
			Instances []struct {
				IPAddr string `json:"ipAddr"`
				Port   struct {
					Value int `json:"$"`
				} `json:"port"`
			} `json:"instance"`
		} `json:"application"`
	} `json:"applications"`
}
