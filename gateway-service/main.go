package main

import (
	"gateway/config"
	"gateway/pkg"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Middleware CORS
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("1GET_Request", c.Request.Header, c.Request.Method)
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Cho phép tất cả domain (không nên dùng * trong production)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Lang") // Thêm 'Lang' vào đây
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")                          // Nếu dùng credentials (cookies, auth)
		log.Println("WRITE4")

		log.Println("2GET_Request", c.Request.Method == "OPTIONS")

		// Xử lý preflight request (OPTIONS)
		if c.Request.Method == "OPTIONS" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			c.Writer.WriteHeader(http.StatusOK)
			return
		}

		log.Println("GET_NEXT", c.Request.Method)
		c.Next()
	}
}

func main() {
	config.LoadConfig()

	// apiKey := viper.GetString("api.key")
	port := viper.GetString("server.port")

	// url := fmt.Sprintf("http://%s:%s")
	gateway := pkg.NewGateway("url")
	// gateway.FetchServices()

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Cho phép gọi từ frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Set-Cookie"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{"message": "pong"})
	// })

	r.Any("/:service/*path", gateway.ProxyRequest)

	log.Println("API Gateway started on port " + port)
	r.Run(":" + port)
}
