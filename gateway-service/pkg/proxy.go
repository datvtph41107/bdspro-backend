package pkg

import (
	"fmt"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func (g *Gateway) ProxyRequest(c *gin.Context) {
	serviceName := c.Param("service")
	path := c.Param("path")

	// serviceKey := strings.ToUpper(serviceName + "-service")

	serviceIp := 0
	if serviceName == "config" {
		// serviceKey = "CONFIG-SERVER"
		serviceIp = 8003
	} else if serviceName == "user" {
		serviceIp = 8001
	} else if serviceName == "file" {
		serviceIp = 8002
	} else if serviceName == "notification" {
		serviceIp = 8004
	} else if serviceName == "payment" {
		serviceIp = 8005
	} else if serviceName == "crm" {
		serviceIp = 8006
	} else if serviceName == "membership" {
		serviceIp = 8007
	} else if serviceName == "map" {
		serviceIp = 8101
	} else if serviceName == "bdspro" {
		serviceIp = 8102
	} else if serviceName == "tho247" {
		serviceIp = 8103
	} else if serviceName == "chat" {
		serviceIp = 8104
	}

	// serviceRaw, ok := g.Services.Load(serviceKey)
	// if !ok {
	// 	log.Printf("Service %s not found", serviceName)
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
	// 	return
	// }

	// service, _ := serviceRaw.(EurekaService)
	// targetURL := fmt.Sprintf("http://%s/%s%s", serviceName, service.Port, serviceName, path)
	targetURL := fmt.Sprintf(
		"http://%s:%d/%s%s?%s",
		serviceName,
		serviceIp,
		serviceName,
		path,
		c.Request.URL.RawQuery,
	)

	fmt.Printf(targetURL)

	resp, err := g.Client.R().
		SetHeaderMultiValues(c.Request.Header).
		SetBody(c.Request.Body).
		Execute(c.Request.Method, targetURL)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ext := filepath.Ext(c.Request.URL.Path)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/json"
	}

	// ✅ Sao chép tất cả header từ BE -> FE (bao gồm Set-Cookie)
	for key, values := range resp.Header() {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}
	// ✅ Gửi dữ liệu về FE
	c.Data(resp.StatusCode(), resp.Header().Get("Content-Type"), resp.Body())

	// ✅ Đảm bảo CORS cho phép gửi Cookie
	// c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000") // Không được dùng "*"
	// c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	// target, _ := url.Parse("http://userservice:8080") // Đổi thành URL của User Service
	// proxy := httputil.NewSingleHostReverseProxy(target)
	// c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
	// c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	// proxy.ServeHTTP(c.Writer, c.Request)

	// c.Data(resp.StatusCode(), contentType, resp.Body())
}
