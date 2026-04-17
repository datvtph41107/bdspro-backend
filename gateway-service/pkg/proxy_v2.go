package pkg

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"gateway/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func NewProxyHandler(cfg *config.GatewayConfig) gin.HandlerFunc {
	client := resty.New().SetTimeout(time.Duration(cfg.Gateway.UpstreamTimeout) * time.Second)
	return func(c *gin.Context) {
		serviceName := c.Param("service")
		path := c.Param("path")

		ep, ok := cfg.Services.Endpoints[serviceName]
		if !ok || len(ep.Addresses) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("service %s not found", serviceName)})
			return
		}

		// Chọn địa chỉ đầu tiên (có thể nâng cấp round‑robin)
		addr := ep.Addresses[0]
		targetURL := fmt.Sprintf("http://%s/%s%s?%s", addr, serviceName, path, c.Request.URL.RawQuery)
		log.Printf("[proxy] %s %s", c.Request.Method, targetURL)

		resp, err := client.R().
			SetHeaderMultiValues(c.Request.Header).
			SetBody(c.Request.Body).
			Execute(c.Request.Method, targetURL)

		if err != nil {
			log.Printf("[proxy] error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "upstream error"})
			return
		}

		for k, v := range resp.Header() {
			c.Writer.Header()[k] = v
		}
		c.Data(resp.StatusCode(), resp.Header().Get("Content-Type"), resp.Body())
	}
}
