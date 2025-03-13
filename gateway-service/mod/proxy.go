package mod

import (
	"fmt"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func (g *Gateway) ProxyRequest(c *gin.Context) {
	serviceName := c.Param("service")
	path := c.Param("path")

	serviceKey := strings.ToUpper(serviceName + "-service")
	if serviceName == "config" {
		serviceKey = "CONFIG-SERVER"
	}

	serviceRaw, ok := g.Services.Load(serviceKey)
	if !ok {
		log.Printf("Service %s not found", serviceName)
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	service, _ := serviceRaw.(EurekaService)
	targetURL := fmt.Sprintf("http://%s:%d/%s%s", service.IP, service.Port, serviceName, path)

	resp, err := g.Client.R().
		SetHeaderMultiValues(c.Request.Header).
		SetBody(c.Request.Body).
		Execute(c.Request.Method, targetURL)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to proxy request"})
		return
	}

	ext := filepath.Ext(c.Request.URL.Path)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/json"
	}

	c.Data(resp.StatusCode(), contentType, resp.Body())
}
