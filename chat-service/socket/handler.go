package socket

import (
	"chat/models"
	chatpb "chat/proto"
	"common/token"
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// type Client struct {
// 	Conn *websocket.Conn
// 	User string
// }

var clients = make(map[uint64]*websocket.Conn)

// Lưu nhóm chat groupId -> danh sách userId
var groups = make(map[uint64][]uint64)

// var connectLock sync.Mutex

func HandleWebSocket(c *gin.Context) {
	authStr := c.Query("token")
	claims, err := token.ParseJWT(authStr)
	userId := claims.ProfileId
	// userId, err := strconv.ParseUint(claims.ProfileId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	if userId == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user parameter"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket Upgrade Error:", err)
		return
	}
	defer conn.Close()

	clients[userId] = conn
	groupIds, _ := models.GetGroupIDsByMemberID(userId)
	log.Printf("groupIds", groupIds)
	for _, groupId := range groupIds {
		members, _ := models.GetMemberIDsByGroupID(groupId)
		log.Printf("member", groupId, members)
		groups[groupId] = members
	}

	// Lắng nghe Redis để nhận tin nhắn
	go listenToRedis(conn, userId)

	// Nhận và gửi tin nhắn từ client
	for {
		log.Printf("TIN nhắn tới")
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("ReadMessage Error:", err)
			// userId, _ := strconv.ParseUint(user, 10, 64)
			delete(clients, userId)
			break
		}
		handleIncomingMessage(msg)
	}
}

func listenToRedis(conn *websocket.Conn, user uint64) {
	log.Printf("Chạy ngay lập tức")
	pubsub := redisClient.Subscribe(context.Background(), "chat:"+strconv.FormatUint(user, 10))
	ch := pubsub.Channel()
	for msg := range ch {
		log.Printf("TIN nhắn đi")
		var chatMsg chatpb.ChatMessage
		if err := proto.Unmarshal([]byte(msg.Payload), &chatMsg); err != nil {
			log.Println("Protobuf Decode Error:", err)
			continue
		}

		encodedMsg, _ := proto.Marshal(&chatMsg)
		if err := conn.WriteMessage(websocket.BinaryMessage, encodedMsg); err != nil {
			log.Println("WriteMessage Error:", err)
			return
		}
	}
}

func handleIncomingMessage(msg []byte) {
	var chatMsg chatpb.ChatMessage
	if err := proto.Unmarshal(msg, &chatMsg); err != nil {
		log.Println("Protobuf Decode Error:", err)
		return
	}

	encodedMsg, _ := proto.Marshal(&chatMsg)

	if len(chatMsg.To) > 5 && chatMsg.To[:5] == "group" {
		groupIdStr := chatMsg.To[6:]
		groupId, _ := strconv.ParseUint(groupIdStr, 10, 64)
		memberIds := groups[groupId]
		for _, memId := range memberIds {
			// clients[memId].WriteMessage()
			target := "chat:" + strconv.FormatUint(memId, 32)
			redisClient.Publish(context.Background(), target, encodedMsg)
		}
		PushMessageToQueue("chat:messages", encodedMsg)
		return
	}
	// Gửi tin nhắn qua Redis

	redisClient.Publish(context.Background(), "chat:"+chatMsg.To, encodedMsg)
	PushMessageToQueue("chat:messages", encodedMsg)
}
