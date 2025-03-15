package main

import (
	"chat/config"
	_ "chat/docs"
	"chat/models"
	chatpb "chat/proto"
	"chat/routes"
	"chat/socket"
	_db "common/db"
	"common/token"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFile "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/protobuf/proto"
)

var publicRoutes = []string{
	"/chat/swagger",
	// "/chat",
}

func main() {
	config.LoadConfig()
	_db.NewDB()
	models.AutoMigrate()

	r := gin.Default()

	// Goroutine để ghi tin nhắn vào DB theo nhóm
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			// Lấy các tin nhắn từ Redis và ghi vào DB
			messages, err := socket.GetMessagesFromQueue("chat:messages")
			socket.DeleteQueue("chat:messages")
			if err != nil {
				log.Println("Error fetching messages from Redis:", err)
				continue
			}

			var chatMessages []*chatpb.ChatMessage
			for _, rawMsg := range messages {
				var chatMsg chatpb.ChatMessage
				if err := proto.Unmarshal([]byte(rawMsg), &chatMsg); err != nil {
					log.Println("Protobuf Decode Error:", err)
					continue
				}
				chatMessages = append(chatMessages, &chatMsg)
			}

			err = models.SaveMessagesToDB(chatMessages)
			if err != nil {
				log.Println("Error saving messages to DB:", err)
			}
		}
	}()

	// WebSocket handler
	r.GET("/chat/connect/ws", socket.HandleWebSocket)

	http_r := r.Group("/chat")
	http_r.Use(token.JWTAuthMiddleware(publicRoutes))
	http_r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFile.Handler))
	http_r.GET("/info", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Chat service avaiable"})
	})

	http_r.GET("/message/history", routes.GetChatHistory)
	http_r.GET("/group/list", routes.GetGroupList)

	port := fmt.Sprintf(":%s", viper.GetString("server.port"))
	log.Println("Chat server is running on", port)
	r.Run(port)
}
