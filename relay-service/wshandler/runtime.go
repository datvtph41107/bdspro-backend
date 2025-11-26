package wshandler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"

	"relay/constants"
	"relay/redis"
	"relay/utils"

	"github.com/gorilla/websocket"
	"github.com/hyperledger/fabric/common/flogging"
)

var wsLogger = flogging.MustGetLogger("ws_handler")

type WebSocketHandler struct {
	upgrader    websocket.Upgrader
	redisClient redisClient
	clients     map[uint64]*websocket.Conn // Conn -> userID
	// rooms       map[string]map[*websocket.Conn]bool // roomID -> list of connections
	rooms      map[uint64][]uint64 // roomID -> list of userIDs: gửi socket tới room
	users      map[uint64][]uint64 // userID -> list of roomIDs: check online
	mu         sync.Mutex
	chatClient ChatClient
	userClient UserClient
}

type RedisMessage struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

func NewWebSocketHandler(chatClient ChatClient, userClient UserClient) *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		redisClient: redis.NewRedisClient(),
		clients:     make(map[uint64]*websocket.Conn),
		// rooms:       make(map[string]map[*websocket.Conn]bool),
		rooms:      make(map[uint64][]uint64),
		users:      make(map[uint64][]uint64),
		chatClient: chatClient,
		userClient: userClient,
	}
}

// func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
// 	conn, err := h.upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		wsLogger.Errorf("Failed to upgrade to WebSocket: %v", err)
// 		return
// 	}
// 	userId := utils.GetCurrentUserID(r.Context())
// 	if userId == 0 {
// 		conn.Close()
// 		wsLogger.Error("Missing user_id")
// 		return
// 	}
// 	userIdStr := strconv.FormatUint(userId, 10)
// 	h.mu.Lock()
// 	h.clients[userIdStr] = conn
// 	h.mu.Unlock()

// 	wsLogger.Infof("User %s connected", userId)

// 	go h.handleDisconnect(conn, userIdStr)
// }

// func (h *WebSocketHandler) HandleRoomWebSocket(w http.ResponseWriter, r *http.Request) {
// 	conn, err := h.upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		wsLogger.Errorf("Failed to upgrade to WebSocket: %v", err)
// 		return
// 	}

// 	userId := utils.GetCurrentUserID(r.Context())

// 	if userId == 0 {
// 		conn.Close()
// 		wsLogger.Error("Missing user_id")
// 		return
// 	}
// 	userIdStr := strconv.FormatUint(userId, 10)

// 	roomID := r.URL.Query().Get("roomId")

// 	if roomID == "" {
// 		conn.Close()
// 		wsLogger.Error("Missing  roomId")
// 		return
// 	}
// 	roomIDUint, err := strconv.ParseUint(roomID, 10, 64)
// 	if err != nil {
// 		conn.Close()
// 		wsLogger.Error("Invalid roomId format")
// 		return
// 	}
// 	err = h.chatClient.ValidateConversationAndCurrentUser(r.Context(), roomIDUint)
// 	if err != nil {
// 		conn.Close()
// 		wsLogger.Error("Invalid roomId")
// 		return
// 	}
// 	h.mu.Lock()
// 	if h.rooms[roomID] == nil {
// 		h.rooms[roomID] = make(map[*websocket.Conn]bool)
// 	}
// 	h.rooms[roomID][conn] = true
// 	h.clients[userIdStr] = conn

// 	h.mu.Unlock()

// 	wsLogger.Infof("User %s joined room %s", userIdStr, roomID)

// 	go h.handleDisconnect(conn, userIdStr)
// }

func (h *WebSocketHandler) HandleConnect(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		wsLogger.Errorf("Failed to upgrade to WebSocket: %v", err)
		return
	}

	userId := utils.GetCurrentUserID(r.Context())
	print("userId: %d request: %s", userId, r.RequestURI)

	if userId == 0 {
		conn.Close()
		wsLogger.Error("Missing user_id")
		return
	}

	// userIdStr := strconv.FormatUint(userId, 10)

	// roomID := r.URL.Query().Get("roomId")

	// if roomID == "" {
	// 	conn.Close()
	// 	wsLogger.Error("Missing  roomId")
	// 	return
	// }
	// roomIDUint, err := strconv.ParseUint(roomID, 10, 64)
	// if err != nil {
	// 	conn.Close()
	// 	wsLogger.Error("Invalid roomId format")
	// 	return
	// }
	// err = h.chatClient.ValidateConversationAndCurrentUser(r.Context(), roomIDUint)
	// if err != nil {
	// 	conn.Close()
	// 	wsLogger.Error("Invalid roomId")
	// 	return
	// }
	h.mu.Lock()
	// if h.rooms[roomID] == nil {
	// 	h.rooms[roomID] = make(map[*websocket.Conn]bool)
	// }
	// h.rooms[roomID][conn] = true
	h.clients[userId] = conn
	h.loadRoomsOfMember(userId)

	h.mu.Unlock()

	// Lưu trạng thái online vào Redis
	h.saveUserOnlineStatus(context.Background(), userId, true)

	wsLogger.Infof("User %s connected", userId)

	// Subscribe to notification channel for this user
	go h.subscribeToNotificationChannel(conn, userId)

	go h.handleDisconnect(conn, userId)
}

func (h *WebSocketHandler) roomsOfUser(userId uint64) []uint64 {
	rooms := h.users[userId]
	for _, roomID := range rooms {
		if _, exists := h.rooms[roomID]; !exists {
			h.loadMemberOfRoom(roomID)
		}
	}
	return rooms
}

func (h *WebSocketHandler) SubscribeToRedisChannel(channel string) {
	pubsub := h.redisClient.Subscribe(context.Background(), channel)
	messages := pubsub.Channel()

	go func() {
		for msg := range messages {
			wsLogger.Infof("Received message from channel %s: %s", channel, msg.Payload)
			var redisMsg RedisMessage
			if err := json.Unmarshal([]byte(msg.Payload), &redisMsg); err != nil {
				wsLogger.Errorf("Invalid JSON format: %v", err)
				continue
			}

			h.broadcastMessage(redisMsg)
		}
	}()
}

func (h *WebSocketHandler) handleDisconnect(conn *websocket.Conn, userID uint64) {
	// Đọc từ connection để detect disconnect
	// Khi connection đóng, ReadMessage sẽ trả về error
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			// Connection đã đóng hoặc có lỗi
			break
		}
	}

	// Cleanup khi disconnect
	h.mu.Lock()
	delete(h.clients, userID)
	// for room, members := range h.rooms {
	// 	if _, exists := members[conn]; exists {
	// 		delete(members, conn)
	// 		if len(members) == 0 {
	// 			delete(h.rooms, room)
	// 		}
	// 	}
	// }
	// for room, memberIds := range h.roomMembers {
	// 	if _, exists := h.clients[userID]; exists {
	// 		delete(h.roomMembers[room], userID)
	// 		if len(h.roomMembers[room]) == 0 {
	// 			delete(h.rooms, room)
	// 		}
	// 	}
	// }
	h.mu.Unlock()

	// Lưu trạng thái offline vào Redis
	h.saveUserOnlineStatus(context.Background(), userID, false)

	// Cập nhật last seen trong user-service
	if h.userClient != nil {
		ctx := context.Background()
		if err := h.userClient.UpdateLastSeen(ctx, userID); err != nil {
			wsLogger.Errorf("Failed to update last seen for user %d: %v", userID, err)
		} else {
			wsLogger.Infof("Updated last seen for user %d", userID)
		}
	}

	wsLogger.Infof("User %d disconnected", userID)

	// Đóng connection
	conn.Close()
}

func (h *WebSocketHandler) broadcastMessage(msg RedisMessage) {
	h.mu.Lock()
	defer h.mu.Unlock()

	data, err := json.Marshal(msg)
	if err != nil {
		wsLogger.Errorf("Failed to marshal message: %v", err)
		return
	}

	m, ok := msg.Data.(map[string]any)
	if !ok {
		wsLogger.Errorf("Invalid message format: %+v", msg.Data)
		return
	}

	switch msg.Type {

	case constants.CREATE_ROOM_TYPE:
		if rawMembers, ok := m["members"].([]any); ok {
			members, err := utils.ParseUint64Slice(rawMembers)
			if err != nil {
				wsLogger.Errorf("Failed to parse members: %v", err)
				return
			}
			h.sendToUser(members, data)
		}

	case constants.JOIN_ROOM_TYPE:
		roomID, okRoom := utils.ParseUint64(m["roomId"])
		userID, okUser := utils.ParseUint64(m["userId"])
		if okRoom && okUser {
			h.addUserToRoom(roomID, userID)
			h.sendToRoom(roomID, data)
		} else {
			wsLogger.Errorf("Failed to parse roomId or userId in JOIN_ROOM_TYPE")
		}

	case constants.LEAVE_ROOM_TYPE:
		roomID, okRoom := utils.ParseUint64(m["roomId"])
		userID, okUser := utils.ParseUint64(m["userId"])
		if okRoom && okUser {
			h.sendToRoom(roomID, data)
			h.removeUserFromRoom(roomID, userID)
		} else {
			wsLogger.Errorf("Failed to parse roomId or userId in LEAVE_ROOM_TYPE")
		}

	case constants.DELETE_ROOM_TYPE:
		if roomID, ok := utils.ParseUint64(m["roomId"]); ok {
			h.sendToRoom(roomID, data)
			delete(h.rooms, roomID)
		} else {
			wsLogger.Errorf("Failed to parse roomId in DELETE_ROOM_TYPE")
		}

	case constants.REMOVE_MEMBER_TYPE:
		roomID, okRoom := utils.ParseUint64(m["roomId"])
		userID, okUser := utils.ParseUint64(m["userId"])
		if okRoom && okUser {
			h.sendToRoom(roomID, data)
			h.removeUserFromRoom(roomID, userID)
		} else {
			wsLogger.Errorf("Failed to parse roomId or userId in REMOVE_MEMBER_TYPE")
		}
	case constants.UPDATE_ROOM_TYPE:
		if rawMembers, ok := m["members"].([]any); ok {
			members, err := utils.ParseUint64Slice(rawMembers)
			if err != nil {
				wsLogger.Errorf("Failed to parse members in UPDATE_ROOM_TYPE: %v", err)
				return
			}
			for _, userID := range members {
				h.sendToUser([]uint64{userID}, data)
			}
		}

	case constants.UPDATE_SETTINGS_TYPE:
		// Truyền tin nhắn đến phòng theo loại sự kiện
		if conversationID, ok := utils.ParseUint64(m["conversationId"]); ok {
			h.sendToRoom(conversationID, data)
		} else {
			wsLogger.Errorf("Failed to parse conversationId in UPDATE_SETTINGS_TYPE")
		}

	case constants.RECALL_MESSAGE_TYPE,
		constants.EDIT_MESSAGE_TYPE,
		constants.PIN_MESSAGE_TYPE,
		constants.REACT_MESSAGE_TYPE,
		constants.SEEN_TYPE,
		constants.DELETE_VIOLATION_MESSAGE_TYPE,
		constants.READ_TYPE,
		constants.MESSAGE_TYPE:

		// Truyền tin nhắn đến phòng theo loại sự kiện
		if roomID, ok := utils.ParseUint64(m["roomId"]); ok {
			h.sendToRoom(roomID, data)
		} else {
			wsLogger.Errorf("Failed to parse roomId in message type: %s", msg.Type)
		}

	case constants.NOTIFICATION_TYPE:
		// Xử lý notification event - gửi đến user cụ thể
		if ownerID, ok := utils.ParseUint64(m["ownerId"]); ok {
			h.sendToUser([]uint64{ownerID}, data)
		} else {
			wsLogger.Errorf("Failed to parse ownerId in NOTIFICATION_TYPE")
		}

	default:
		wsLogger.Warnf("Unhandled message type: %s", msg.Type)
	}
}

// func (h *WebSocketHandler) sendToRoom(roomID string, data []byte) {
// 	if roomClients, exists := h.rooms[roomID]; exists {
// 		for conn := range roomClients {
// 			err := conn.WriteMessage(websocket.TextMessage, data)
// 			if err != nil {
// 				wsLogger.Errorf("Failed to send message to room %s: %v", roomID, err)
// 			}

// 		}
// 	} else {
// 		wsLogger.Warnf("Room %s not found, message not delivered", roomID)
// 	}
// }

func (h *WebSocketHandler) addUserToRoom(roomID uint64, userID uint64) {
	if _, exists := h.rooms[roomID]; !exists {
		h.rooms[roomID] = make([]uint64, 0)
	}
	h.rooms[roomID] = append(h.rooms[roomID], userID)
	// for uid, conn := range h.clients {
	// 	if uid == userID {

	// 	}
	// }
}

func (h *WebSocketHandler) removeUserFromRoom(roomID uint64, userID uint64) {
	members := h.rooms[roomID]
	for i, member := range members {
		if member == userID {
			h.rooms[roomID] = append(members[:i], members[i+1:]...)
			break
		}
	}
	// for uid, conn := range h.clients {
	// 	if uid == userID {
	// 		delete(h.rooms[roomID], conn)
	// 	}
	// }
}

func (h *WebSocketHandler) sendToUser(members []uint64, data []byte) {
	for _, userId := range members {
		if h.clients[userId] != nil {
			conn := h.clients[userId]
			err := conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				wsLogger.Errorf("Failed to send message to user %s: %v", userId, err)
			}
		}
	}
}

func (h *WebSocketHandler) loadMemberOfRoom(roomID uint64) {
	// roomIDUint, err := strconv.ParseUint(roomID, 10, 64)
	// if err != nil {
	// 	wsLogger.Errorf("Invalid roomId format: %v", err)
	// 	return
	// }
	members, err := h.chatClient.GetMembersOfRoom(context.Background(), roomID)
	if err != nil {
		wsLogger.Errorf("Failed to get room members: %v", err)
		return
	}
	h.rooms[roomID] = members
}

func (h *WebSocketHandler) loadRoomsOfMember(userId uint64) {
	// roomIDUint, err := strconv.ParseUint(roomID, 10, 64)
	// if err != nil {
	// 	wsLogger.Errorf("Invalid roomId format: %v", err)
	// 	return
	// }
	members, err := h.chatClient.GetRoomsOfMember(context.Background(), userId)
	if err != nil {
		wsLogger.Errorf("Failed to get room members: %v", err)
		return
	}
	h.rooms[userId] = members
}

func (h *WebSocketHandler) sendToRoom(roomID uint64, data []byte) {
	_, ok := h.rooms[roomID]
	wsLogger.Infof("Room %s: %v, %v", roomID, ok, string(data))
	if !ok {
		h.loadMemberOfRoom(roomID)
	}
	members, ok := h.rooms[roomID]
	if !ok {
		wsLogger.Warnf("Room %s not found, message not delivered2", roomID)
		return
	}

	for _, clientId := range members {
		conn := h.clients[clientId]
		if conn == nil {
			continue
		}
		err := conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			wsLogger.Errorf("Failed to send message to all users: %v", err)
		}
	}
}

// subscribeToNotificationChannel subscribes to notification channel for a specific user
func (h *WebSocketHandler) subscribeToNotificationChannel(conn *websocket.Conn, userID uint64) {
	channel := "notification:" + strconv.FormatUint(userID, 10)
	pubsub := h.redisClient.Subscribe(context.Background(), channel)
	messages := pubsub.Channel()

	wsLogger.Infof("Subscribed to notification channel: %s for user: %d", channel, userID)

	for msg := range messages {
		wsLogger.Infof("Received notification from channel %s: %s", channel, msg.Payload)
		// Send notification directly to the WebSocket connection
		if err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
			wsLogger.Errorf("Failed to send notification to user %d: %v", userID, err)
			return
		}
	}
}

// saveUserOnlineStatus lưu trạng thái online/offline của user vào Redis
// Format: <true/false>|<timestamp>
func (h *WebSocketHandler) saveUserOnlineStatus(ctx context.Context, userID uint64, isOnline bool) {
	key := constants.USER_ONLINE_STATUS_KEY_PREFIX + strconv.FormatUint(userID, 10)
	timestamp := time.Now().Format(time.RFC3339)

	var status string
	if isOnline {
		status = "true|" + timestamp
	} else {
		status = "false|" + timestamp
	}

	if err := h.redisClient.Set(ctx, key, status); err != nil {
		wsLogger.Errorf("Failed to save user online status for user %d: %v", userID, err)
	} else {
		wsLogger.Infof("Saved user online status for user %d: %s", userID, status)
	}
}
