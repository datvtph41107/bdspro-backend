package provider

import "context"

// ChatProvider interface cho việc tạo conversation/room trong chat service
type ChatProvider interface {
	// CreateConversation tạo một conversation mới
	// name: tên conversation
	// createdBy: ID của user tạo conversation
	// members: danh sách ID của các member (không bao gồm createdBy)
	CreateConversation(ctx context.Context, name string, createdBy uint64, members []uint64) (uint64, error)
}
