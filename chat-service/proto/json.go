package chatpb

import (
	"google.golang.org/protobuf/types/known/anypb"
)

// type ChatMessage struct {
// 	From    uint64 `protobuf:"varint,1,opt,name=from,proto3" json:"from,omitempty"`
// 	To      string `protobuf:"bytes,2,opt,name=to,proto3" json:"to,omitempty"`
// 	Message string `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
// }

// type CreateRoom struct {
// 	RoomName string   `protobuf:"bytes,1,opt,name=room_name,json=roomName,proto3" json:"room_name,omitempty"`
// 	Members  []uint64 `protobuf:"varint,2,rep,packed,name=members,proto3" json:"members,omitempty"`
// }

// type ChatAction struct {
// 	Type string     `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
// 	Data *anypb.Any `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
// }

// ChatMessage đại diện cho tin nhắn trò chuyện
type ChatMessage struct {
	From    uint64 `json:"from,omitempty" protobuf:"varint,1,opt,name=from,proto3"`
	To      string `json:"to,omitempty" protobuf:"bytes,2,opt,name=to,proto3"`
	Message string `json:"message,omitempty" protobuf:"bytes,3,opt,name=message,proto3"`
}

// CreateRoom dùng để tạo nhóm trò chuyện
type CreateRoom struct {
	RoomName string   `json:"room_name,omitempty" protobuf:"bytes,1,opt,name=room_name,json=roomName,proto3"`
	Members  []uint64 `json:"members,omitempty" protobuf:"varint,2,rep,packed,name=members,proto3"`
}

// ChatAction đại diện cho một hành động chat
type ChatAction struct {
	Type string     `json:"type,omitempty" protobuf:"bytes,1,opt,name=type,proto3"`
	Data *anypb.Any `json:"data,omitempty" protobuf:"bytes,2,opt,name=data,proto3"`
}
