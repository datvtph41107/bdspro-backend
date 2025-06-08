package pkg

import (
	"fmt"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func ForwardingRequest(g *Gateway, env string) func(c *gin.Context) {
	return func(c *gin.Context) {
		fmt.Println("ProxyRequest")
		serviceName := c.Param("service")
		path := c.Param("path")

		serviceIp := 0
		if serviceName == "config" {
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
		} else if serviceName == "search" {
			serviceIp = 8008
		} else if serviceName == "map" {
			serviceIp = 8101
		} else if serviceName == "bdspro" {
			serviceIp = 8102
		} else if serviceName == "tho247" {
			serviceIp = 8103
		} else if serviceName == "chat" {
			serviceIp = 8104
		}

		url := "localhost"
		if env == "release" {
			url = serviceName
		}
		targetURL := fmt.Sprintf(
			"http://%s:%d/%s%s?%s",
			url,
			serviceIp,
			serviceName,
			path,
			c.Request.URL.RawQuery,
		)

		fmt.Println(targetURL)
		fmt.Println(c.Request.Header)
		fmt.Println("c.Request.Method", c.Request.Method)

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
				// fmt.Println(key, value, " _HEADER_RESPONSE")
			}
		}

		// c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Không được dùng "*"
		// fmt.Println(resp.Body(), " _BODY_RESPONSE")
		// ✅ Gửi dữ liệu về FE
		// "Access-Control-Allow-Origin"
		// c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8102")
		// c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// fmt.Println(c.Writer.Header(), "HEADER_REPONSE.......")
		c.Data(resp.StatusCode(), resp.Header().Get("Content-Type"), resp.Body())

		// ✅ Đảm bảo CORS cho phép gửi Cookie

		// target, _ := url.Parse("http://userservice:8080") // Đổi thành URL của User Service
		// proxy := httputil.NewSingleHostReverseProxy(target)
		// c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
		// c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		// proxy.ServeHTTP(c.Writer, c.Request)

		// c.Data(resp.StatusCode(), contentType, resp.Body())
	}

}
