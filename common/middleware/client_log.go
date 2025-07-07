package _middleware

import (
	_constants "common/constants"
	"context"
	"log"
	"strconv"
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
	appendUint64IfExists := func(ctxKey _constants.CtxKey, mdKey string) {
		if val, ok := ctx.Value(ctxKey).(uint64); ok {
			md.Append(mdKey, strconv.FormatUint(val, 10))
		}
	}

	appendStringIfExists := func(ctxKey _constants.CtxKey, mdKey string) {
		if val, ok := ctx.Value(ctxKey).(string); ok && val != "" {
			md.Append(mdKey, val)
		}
	}

	appendTimeIfExists := func(ctxKey _constants.CtxKey, mdKey string) {
		if val, ok := ctx.Value(ctxKey).(time.Time); ok && !val.IsZero() {
			md.Append(mdKey, val.Format(time.RFC3339))
		}
	}

	// Gắn các giá trị cần truyền sang service B
	appendUint64IfExists(_constants.ProfileIDKey, "profileId")
	appendUint64IfExists(_constants.OrganizationIDKey, "organizationId")
	appendUint64IfExists(_constants.AuthIDKey, "authId")
	appendUint64IfExists(_constants.SessionKey, "session")
	appendUint64IfExists(_constants.PlanIDKey, "planId")

	appendStringIfExists(_constants.RoleKey, "role")
	appendStringIfExists(_constants.TypeKey, "type")

	appendTimeIfExists(_constants.PlanFromKey, "planFrom")

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
