package main

import (
	chatpb "chat/proto"
	"encoding/base64"
	"fmt"
	"log"

	"google.golang.org/protobuf/proto"
)

func main() {
	msg := &chatpb.ChatMessage{
		From:    "Alo",
		To:      "group:1",
		Message: "Hello test redis",
	}

	// Chuyển đổi tin nhắn thành byte array
	data, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal("Error marshaling message:", err)
	}
	base64Data := base64.StdEncoding.EncodeToString(data)

	// Lưu ra file để gửi qua Postman
	// err = os.WriteFile("message.bin", data, 0644)
	// if err != nil {
	// 	log.Fatal("Error saving file:", err)
	// }

	fmt.Println("Message saved to message: " + base64Data)
}
