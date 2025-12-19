package pkg

import (
	"common"
	"fmt"
	"log"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// ServiceConfig cấu hình cho từng service
type ServiceConfig struct {
	Port int
	Name string
}

// servicePorts map tên service với port tương ứng
var servicePorts = map[string]int{
	"config":       8003,
	"user":         8001,
	"file":         8002,
	"notification": 8004,
	"payment":      8005,
	"crm":          8006,
	"search":       8008,
	"map":          8101,
	"bdspro":       8102,
	"tho247":       8103,
	"chat":         8104,
}

// ForwardingRequest xử lý forwarding request đến các service
func ForwardingRequest(g *Gateway, env string) func(c *gin.Context) {
	return func(c *gin.Context) {
		serviceName := c.Param("service")
		path := c.Param("path")

		// Lấy port từ config
		port, exists := servicePorts[serviceName]
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Service %s not found", serviceName)})
			return
		}

		// Xác định host dựa trên environment
		host := common.GetPrefixProtobufUrl(env, serviceName)

		// Tạo target URL
		targetURL := fmt.Sprintf(
			"http://%s:%d/%s%s?%s",
			host,
			port,
			serviceName,
			path,
			c.Request.URL.RawQuery,
		)

		log.Printf("Forwarding request to: %s", targetURL)

		// Thực hiện request
		resp, err := g.Client.R().
			SetHeaderMultiValues(c.Request.Header).
			SetBody(c.Request.Body).
			Execute(c.Request.Method, targetURL)

		if err != nil {
			log.Printf("Forwarding error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Xác định content type
		ext := filepath.Ext(c.Request.URL.Path)
		contentType := mime.TypeByExtension(ext)
		if contentType == "" {
			contentType = "application/json"
		}

		// Copy headers từ response
		for key, values := range resp.Header() {
			for _, value := range values {
				c.Writer.Header().Add(key, value)
			}
		}

		// Trả về response
		c.Data(resp.StatusCode(), resp.Header().Get("Content-Type"), resp.Body())
	}
}
