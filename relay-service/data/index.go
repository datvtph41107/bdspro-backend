package data

type ChatMessage struct {
	From    uint64 `json:"from"`
	To      uint64 `json:"to"`
	Message string `json:"message"`
}

type CreateRoom struct {
	RoomName string   `json:"room_name"`
	Members  []uint64 `json:"members"`
}

type ChatAction struct {
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}
