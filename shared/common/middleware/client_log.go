package _middleware

import (
	_enums "common/domain/enum"
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryClientInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	start := time.Now()

	// Tạo metadata từ context hiện tại
	md := metadata.MD{}

	// Các hàm helper để lấy giá trị từ context và gắn vào metadata
	appendUint64IfExists := func(ctxKey _enums.ContextKey, mdKey string) {
		if val, ok := ctx.Value(ctxKey).(uint64); ok {
			md.Append(mdKey, strconv.FormatUint(val, 10))
		}
	}

	appendStringIfExists := func(ctxKey _enums.ContextKey, mdKey string) {
		if val, ok := ctx.Value(ctxKey).(string); ok && val != "" {
			md.Append(mdKey, val)
		}
	}

	appendTimeIfExists := func(ctxKey _enums.ContextKey, mdKey string) {
		if val, ok := ctx.Value(ctxKey).(time.Time); ok && !val.IsZero() {
			md.Append(mdKey, val.Format(time.RFC3339))
		}
	}

	// Gắn các giá trị cần truyền sang service B
	appendUint64IfExists(_enums.ProfileIDKey, "profileId")
	appendUint64IfExists(_enums.OriginIDKey, "originId")
	appendUint64IfExists(_enums.OrganizationIDKey, "organizationId")
	appendUint64IfExists(_enums.AuthIDKey, "authId")
	appendUint64IfExists(_enums.SessionKey, "session")
	appendUint64IfExists(_enums.PlanIDKey, "planId")

	appendStringIfExists(_enums.RoleKey, "role")
	appendUint64SliceIfExists := func(ctxKey _enums.ContextKey, mdKey string) {
		roleIds, ok := ctx.Value(ctxKey).([]uint64)
		if !ok || len(roleIds) == 0 {
			return
		}
		parts := make([]string, 0, len(roleIds))
		for _, id := range roleIds {
			parts = append(parts, strconv.FormatUint(id, 10))
		}
		md.Append(mdKey, strings.Join(parts, ","))
	}
	appendUint64SliceIfExists(_enums.RoleIDsKey, "roleIds")
	appendStringIfExists(_enums.TypeKey, "type")

	appendTimeIfExists(_enums.PlanFromKey, "planFrom")

	// Gắn metadata vào context
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Tiếp tục gọi RPC
	err := invoker(ctx, method, req, reply, cc, opts...)
	elapsed := time.Since(start)

	log.Printf("[gRPC Client] Method: %s, Duration: %s, Error: %v", method, elapsed, err)
	return err
}

func MonitorConnection(ctx context.Context, name string, conn *grpc.ClientConn) {
	for {
		state := conn.GetState()
		log.Printf("[gRPC] Connection state of %s: %v", name, state)

		// Wait until the state changes before continuing the loop
		if !conn.WaitForStateChange(ctx, state) {
			log.Println("[gRPC] Context canceled or connection closed")
			return
		}
	}
}
