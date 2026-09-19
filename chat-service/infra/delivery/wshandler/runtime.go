package wshandler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"chat/infra/redis"
	"chat/internal/constants"
	"chat/utils"
	"common/logging"

	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	upgrader             websocket.Upgrader
	redisClient          redisClient
	clients              map[string]*websocket.Conn          // Conn -> userID
	rooms                map[string]map[*websocket.Conn]bool // roomID -> list of connections
	mu                   sync.Mutex
	conversationUsercase conversationUsercase
}

type RedisMessage struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

func NewWebSocketHandler(conversationUsercase conversationUsercase) *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		redisClient:          redis.NewRedisClient(),
		clients:              make(map[string]*websocket.Conn),
		rooms:                make(map[string]map[*websocket.Conn]bool),
		conversationUsercase: conversationUsercase,
	}
}

func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	logger := logging.WithComponent(r.Context(), "ws_handler")
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to upgrade to WebSocket: %v", err))
		return
	}
	userId := utils.GetCurrentUserID(r.Context())
	if userId == 0 {
		conn.Close()
		logger.Error("Missing user_id")
		return
	}
	userIdStr := strconv.FormatUint(userId, 10)
	h.mu.Lock()
	h.clients[userIdStr] = conn
	h.mu.Unlock()

	logger.Info(fmt.Sprintf("User %d connected", userId))

	go h.handleDisconnect(conn, userIdStr)
}

func (h *WebSocketHandler) HandleRoomWebSocket(w http.ResponseWriter, r *http.Request) {
	logger := logging.WithComponent(r.Context(), "ws_handler")
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to upgrade to WebSocket: %v", err))
		return
	}

	userId := utils.GetCurrentUserID(r.Context())

	if userId == 0 {
		conn.Close()
		logger.Error("Missing user_id")
		return
	}
	userIdStr := strconv.FormatUint(userId, 10)

	roomID := r.URL.Query().Get("room_id")

	if roomID == "" {
		conn.Close()
		logger.Error("Missing  room_id")
		return
	}
	roomIDUint, err := strconv.ParseUint(roomID, 10, 64)
	if err != nil {
		conn.Close()
		logger.Error("Invalid room_id format")
		return
	}
	err = h.conversationUsercase.ValidateConversationAndCurrentUser(r.Context(), roomIDUint)
	if err != nil {
		conn.Close()
		logger.Error("Invalid room_id")
		return
	}
	h.mu.Lock()
	if h.rooms[roomID] == nil {
		h.rooms[roomID] = make(map[*websocket.Conn]bool)
	}
	h.rooms[roomID][conn] = true
	h.clients[userIdStr] = conn

	h.mu.Unlock()

	logger.Info(fmt.Sprintf("User %s joined room %s", userIdStr, roomID))

	go h.handleDisconnect(conn, userIdStr)
}

func (h *WebSocketHandler) SubscribeToRedisChannel(channel string) {
	logger := logging.WithComponent(context.Background(), "ws_handler")
	pubsub := h.redisClient.Subscribe(context.Background(), channel)
	messages := pubsub.Channel()

	go func() {
		for msg := range messages {
			logger.Info(fmt.Sprintf("Received message from channel %s: %s", channel, msg.Payload))
			var redisMsg RedisMessage
			if err := json.Unmarshal([]byte(msg.Payload), &redisMsg); err != nil {
				logger.Error(fmt.Sprintf("Invalid JSON format: %v", err))
				continue
			}

			h.broadcastMessage(redisMsg)
		}
	}()
}

func (h *WebSocketHandler) handleDisconnect(conn *websocket.Conn, userID string) {
	logger := logging.WithComponent(context.Background(), "ws_handler")
	defer conn.Close()
	<-context.Background().Done()

	h.mu.Lock()
	delete(h.clients, userID)
	for room, members := range h.rooms {
		if _, exists := members[conn]; exists {
			delete(members, conn)
			if len(members) == 0 {
				delete(h.rooms, room)
			}
		}
	}
	h.mu.Unlock()

	logger.Info(fmt.Sprintf("User %s disconnected", userID))
}

func (h *WebSocketHandler) broadcastMessage(msg RedisMessage) {
	logger := logging.WithComponent(context.Background(), "ws_handler")
	h.mu.Lock()
	defer h.mu.Unlock()

	data, err := json.Marshal(msg)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to marshal message: %v", err))
		return
	}

	m, ok := msg.Data.(map[string]any)
	if !ok {
		logger.Error(fmt.Sprintf("Invalid message format: %+v", msg.Data))
		return
	}

	switch msg.Type {

	case constants.CREATE_ROOM_TYPE:
		if roomName, ok := m["room_id"].(string); ok {
			h.rooms[roomName] = make(map[*websocket.Conn]bool)
		}
		if rawMembers, ok := m["members"].([]any); ok {
			members := make([]string, 0, len(rawMembers))
			for _, v := range rawMembers {
				if str, ok := v.(string); ok {
					members = append(members, str)
				}
			}
			h.sendToUser(members, data)
		}

	case constants.JOIN_ROOM_TYPE:
		if roomID, ok := m["room_id"].(string); ok {
			if userID, ok := m["user_id"].(string); ok {
				h.sendToRoom(roomID, data)

				h.addUserToRoom(roomID, userID)
			}
		}

	case constants.LEAVE_ROOM_TYPE:
		if roomID, ok := m["room_id"].(string); ok {
			if userID, ok := m["user_id"].(string); ok {
				h.sendToRoom(roomID, data)

				h.removeUserFromRoom(roomID, userID)
			}
		}

	case constants.DELETE_ROOM_TYPE:
		if roomID, ok := m["room_id"].(string); ok {
			h.sendToRoom(roomID, data)

			delete(h.rooms, roomID)

		}

	case constants.REMOVE_MEMBER_TYPE:
		if roomID, ok := m["room_id"].(string); ok {
			if userID, ok := m["user_id"].(string); ok {
				h.sendToRoom(roomID, data)

				h.removeUserFromRoom(roomID, userID)

			}
		}
	case constants.UPDATE_ROOM_TYPE:
		if userIds, ok := m["members"].([]any); ok {
			for _, v := range userIds {
				if userID, ok := v.(string); ok {
					h.sendToUser([]string{userID}, data)
				}
			}
		}
	case constants.TYPING_MESSAGE_TYPE:
		var roomID string
		if id, ok := m["roomId"].(string); ok {
			roomID = id
		}
		if roomID != "" {
			h.sendToRoom(roomID, data)
		}

	case constants.RECALL_MESSAGE_TYPE,
		constants.EDIT_MESSAGE_TYPE,
		constants.PIN_MESSAGE_TYPE,
		constants.REACT_MESSAGE_TYPE,
		constants.SEEN_TYPE,
		constants.UPDATE_SETTINGS_TYPE,
		constants.DELETE_VIOLATION_MESSAGE_TYPE,
		constants.READ_TYPE,
		constants.MESSAGE_TYPE:

		// Truyền tin nhắn đến phòng theo loại sự kiện
		// Hỗ trợ cả room_id (snake_case) và roomId (camelCase)
		var roomID string
		if id, ok := m["room_id"].(string); ok {
			roomID = id
		} else if id, ok := m["roomId"].(string); ok {
			roomID = id
		}
		if roomID != "" {
			h.sendToRoom(roomID, data)
		}

	default:
		logger.Warn(fmt.Sprintf("Unhandled message type: %s", msg.Type))
	}
}

func (h *WebSocketHandler) sendToRoom(roomID string, data []byte) {
	logger := logging.WithComponent(context.Background(), "ws_handler")
	if roomClients, exists := h.rooms[roomID]; exists {
		for conn := range roomClients {
			err := conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to send message to room %s: %v", roomID, err))
			}

		}
	} else {
		logger.Warn(fmt.Sprintf("Room %s not found, message not delivered", roomID))
	}
}

func (h *WebSocketHandler) addUserToRoom(roomID, userID string) {
	for uid, conn := range h.clients {
		if uid == userID {
			if _, exists := h.rooms[roomID]; !exists {
				h.rooms[roomID] = make(map[*websocket.Conn]bool)
			}
			h.rooms[roomID][conn] = true
		}
	}
}

func (h *WebSocketHandler) removeUserFromRoom(roomID, userID string) {
	for uid, conn := range h.clients {
		if uid == userID {
			delete(h.rooms[roomID], conn)
		}
	}
}

func (h *WebSocketHandler) sendToUser(members []string, data []byte) {
	logger := logging.WithComponent(context.Background(), "ws_handler")
	for _, userId := range members {
		if h.clients[userId] != nil {
			conn := h.clients[userId]
			err := conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				logger.Error(fmt.Sprintf("Failed to send message to user %s: %v", userId, err))
			}
		}
	}
}
