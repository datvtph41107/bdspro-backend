package socket

import (
	_db "common/db"
	_jwt "common/jwt"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"relay/data"
	"relay/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	claims, err := _jwt.ParseJWT(authStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}
	userId := claims.ProfileId
	// userId, err := strconv.ParseUint(claims.ProfileId, 10, 64)
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
	log.Println("groupIds", groupIds)
	for _, groupId := range groupIds {
		members, _ := models.GetMemberIDsByGroupID(groupId)
		log.Println("member", groupId, members)
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
		handleIncomingMessage(msg, userId)
	}
}

func listenToRedis(conn *websocket.Conn, user uint64) {
	log.Printf("Nhận sự kiện từ REDIS")
	pubsub := redisClient.Subscribe(context.Background(), "chat:"+strconv.FormatUint(user, 10))
	ch := pubsub.Channel()
	for msg := range ch {
		log.Printf("Gửi tin nhắn")

		// chatMsg.From = user
		// encodedMsg, _ := json.Marshal(chatMsg)
		if err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
			log.Println("WriteMessage Error:", err)
			return
		}
	}
}

func handleIncomingMessage(msg []byte, userId uint64) {
	var actionMsg data.ChatAction
	if err := json.Unmarshal(msg, &actionMsg); err != nil {
		log.Println("Protobuf Decode Error:", err)
		return
	}

	dataBytes, err := json.Marshal(actionMsg.Data) // Chuyển interface{} về JSON []byte
	if err != nil {
		log.Println("Protobuf Decode Error:", err)
		return
	}

	// var chatReturn chatpb.ChatReturn
	switch actionMsg.Type {
	case "create_room":
		var createRoom data.CreateRoom

		if err := json.Unmarshal([]byte(dataBytes), &createRoom); err != nil {
			log.Println("Failed to unmarshal CreateRoom:", err)
			return
		}

		// Tạo một cuộc trò chuyện mới
		conversation := models.ConversationModel{
			CreatedBy: userId, // Hoặc lấy từ user gửi request
			Title:     createRoom.RoomName,
			Avatar:    "", // Có thể cập nhật sau
		}

		// Lưu vào database
		if err := _db.DB.Create(&conversation).Error; err != nil {
			log.Println("Failed to create conversation:", err)
			return
		}

		actionMsg.Data = conversation
		encodedMsg, _ := json.Marshal(&actionMsg)

		var members []models.ParticipantModel
		for _, memId := range createRoom.Members {
			members = append(members, models.ParticipantModel{
				MemberId:  memId,
				GroupId:   conversation.ID,
				CreatedBy: userId,
			})
			redisClient.Publish(context.Background(), "chat:"+strconv.FormatUint(memId, 32), encodedMsg)
		}
		_db.DB.CreateInBatches(members, 10)
		groups[conversation.ID] = createRoom.Members

		return
	case "message":
		var chatMsg data.ChatMessage
		if err := json.Unmarshal(dataBytes, &chatMsg); err != nil {
			returnError(userId, err.Error())
			return
		}

		// groupId, _ := strconv.ParseUint(, 10, 32)
		memberIds := groups[chatMsg.To]
		if memberIds == nil {
			returnError(userId, "Phòng không tồn tại")
			return
		}

		chatMsg.From = userId
		actionMsg.Data = chatMsg
		encodedMsg, _ := json.Marshal(&actionMsg)

		for _, memId := range memberIds {
			// clients[memId].WriteMessage()
			target := "chat:" + strconv.FormatUint(memId, 32)
			redisClient.Publish(context.Background(), target, encodedMsg)
		}

		// Gửi tin nhắn qua Redis
		redisClient.Publish(context.Background(), "chat:"+strconv.FormatUint(chatMsg.To, 10), encodedMsg)
		encodedChatMsg, _ := json.Marshal(&chatMsg)
		PushMessageToQueue("chat:messages", encodedChatMsg)
		return
	}

}

func returnError(userId uint64, message string) {
	encodeErr, _ := json.Marshal(data.ChatAction{
		Type:    "error",
		Data:    nil,
		Message: message,
	})
	redisClient.Publish(context.Background(),
		"chat:"+strconv.FormatUint(userId, 10),
		encodeErr,
	)
}
