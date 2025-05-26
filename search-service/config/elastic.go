package config

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

var ES *elasticsearch.Client

func InitElastic() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			AppProperties.Elastic.Search.Url,
		},
		Username: AppProperties.Elastic.Search.User,
		Password: AppProperties.Elastic.Search.Pass,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	ES = es
}
