// package main

// import (
// 	_db "common/db"
// 	_jwt "common/jwt"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"relay/config"
// 	"relay/data"
// 	_ "relay/docs"
// 	"relay/models"
// 	"relay/routes"
// 	"relay/socket"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/spf13/viper"
// 	swaggerFile "github.com/swaggo/files"
// 	ginSwagger "github.com/swaggo/gin-swagger"
// )

// var publicRoutes = []string{
// 	"/chat/swagger",
// 	// "/chat",
// }

// func main() {
// 	config.LoadConfig()
// 	_db.NewDB()
// 	models.AutoMigrate()

// 	r := gin.Default()

// 	// Goroutine để ghi tin nhắn vào DB theo nhóm
// 	go func() {
// 		ticker := time.NewTicker(10 * time.Second)
// 		defer ticker.Stop()

// 		for range ticker.C {
// 			// Lấy các tin nhắn từ Redis và ghi vào DB
// 			messages, err := socket.GetMessagesFromQueue("chat:messages")
// 			socket.DeleteQueue("chat:messages")
// 			if err != nil {
// 				log.Println("Error fetching messages from Redis:", err)
// 				continue
// 			}

// 			var chatMessages []*data.ChatMessage
// 			for _, rawMsg := range messages {
// 				var chatMsg data.ChatMessage
// 				if err := json.Unmarshal([]byte(rawMsg), &chatMsg); err != nil {
// 					log.Println("Protobuf Decode Error:", err)
// 					continue
// 				}

// 				// var chatMsg data.ChatMessage
// 				// if err := json.Unmarshal(chatAction.Data, &chatMsg); err != nil {
// 				// 	log.Println("Failed to unmarshal ChatMessage:", err)
// 				// 	return
// 				// }

// 				chatMessages = append(chatMessages, &chatMsg)
// 			}

// 			err = models.SaveMessagesToDB(chatMessages)
// 			if err != nil {
// 				log.Println("Error saving messages to DB:", err)
// 			}
// 		}
// 	}()

// 	// WebSocket handler
// 	r.GET("/chat/connect/ws", socket.HandleWebSocket)

// 	http_r := r.Group("/chat")
// 	http_r.Use(_jwt.JWTAuthMiddleware(publicRoutes))
// 	http_r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFile.Handler))
// 	http_r.GET("/info", func(c *gin.Context) {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Chat service avaiable"})
// 	})

// 	http_r.GET("/message/history", routes.GetChatHistory)
// 	http_r.GET("/group/list", routes.GetGroupList)

// 	port := fmt.Sprintf(":%s", viper.GetString("server.port"))
// 	log.Println("Chat server is running on", port)
// 	r.Run(port)
// }

package main

import (
	"os"

	"relay/config"
	"relay/server"
)

func main() {
	os.Setenv("TZ", "Asia/Ho_Chi_Minh")
	config.InitConfig()

	server.NewWSServer(config.PORT).Start()
}
