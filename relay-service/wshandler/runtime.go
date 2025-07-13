package wshandler

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

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
}

type RedisMessage struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

func NewWebSocketHandler(chatClient ChatClient) *WebSocketHandler {
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

	wsLogger.Infof("User %s connected", userId)

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
	defer conn.Close()
	<-context.Background().Done()

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

	wsLogger.Infof("User %s disconnected", userID)
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
		// if roomName, ok := m["roomId"].(string); ok {
		// 	h.rooms[roomName] = make(map[*websocket.Conn]bool)
		// }
		if rawMembers, ok := m["members"].([]any); ok {
			members := make([]uint64, 0, len(rawMembers))
			for _, v := range rawMembers {
				if userID, ok := v.(float64); ok {
					members = append(members, uint64(userID))
				}
			}
			h.sendToUser(members, data)
		}

	case constants.JOIN_ROOM_TYPE:
		if roomID, ok := m["roomId"].(float64); ok {
			if userID, ok := m["userId"].(float64); ok {
				h.addUserToRoom(uint64(roomID), uint64(userID))
				h.sendToRoom(uint64(roomID), data)
			}
		}

	case constants.LEAVE_ROOM_TYPE:
		if roomID, ok := m["roomId"].(float64); ok {
			if userID, ok := m["userId"].(float64); ok {
				h.sendToRoom(uint64(roomID), data)

				h.removeUserFromRoom(uint64(roomID), uint64(userID))
			}
		}

	case constants.DELETE_ROOM_TYPE:
		if roomID, ok := m["roomId"].(float64); ok {
			h.sendToRoom(uint64(roomID), data)

			delete(h.rooms, uint64(roomID))
		}

	case constants.REMOVE_MEMBER_TYPE:
		if roomID, ok := m["roomId"].(float64); ok {
			if userID, ok := m["userId"].(float64); ok {
				h.sendToRoom(uint64(roomID), data)

				h.removeUserFromRoom(uint64(roomID), uint64(userID))

			}
		}
	case constants.UPDATE_ROOM_TYPE:
		if userIds, ok := m["members"].([]any); ok {
			for _, v := range userIds {
				if userID, ok := v.(float64); ok {
					h.sendToUser([]uint64{uint64(userID)}, data)
				}
			}
		}

	case constants.UPDATE_SETTINGS_TYPE:
		// Truyền tin nhắn đến phòng theo loại sự kiện
		if conversationID, ok := m["conversationId"].(float64); ok {
			h.sendToRoom(uint64(conversationID), data)
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
		if roomID, ok := m["roomId"].(float64); ok {
			h.sendToRoom(uint64(roomID), data)
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
