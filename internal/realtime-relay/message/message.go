package main

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"
	chatpb "relay/proto"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/proto"
)

func main() {
	msgData, err := ptypes.MarshalAny(&chatpb.ChatMessage{
		From:    1,
		To:      1,
		Message: "xin chào cả nhà hiện tại là gửi cho group 1 có mem là 3 và 4",
	})

	if err != nil {
		slog.Error(
			"marshal Relay ChatMessage",
			slog.String("component", "message.utility"),
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	msg := &chatpb.ChatAction{
		Type: "message",
		Data: msgData, // ✅ Đúng kiểu google.protobuf.Any
	}

	// Chuyển đổi tin nhắn thành byte array
	data, err := proto.Marshal(msg)
	if err != nil {
		slog.Error(
			"marshal Relay message",
			slog.String("component", "message.utility"),
			slog.Any("error", err),
		)
		os.Exit(1)
	}
	base64Data := base64.StdEncoding.EncodeToString(data)

	// Lưu ra file để gửi qua Postman
	// err = os.WriteFile("message.bin", data, 0644)
	// if err != nil {
	// 	slog.Error("save Relay message file", slog.Any("error", err))
	// }

	fmt.Println("Message plain is: " + msg.String())
	fmt.Println("Message base64: " + base64Data)
}
