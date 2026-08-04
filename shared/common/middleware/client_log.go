package _middleware

import (
	_transporttrust "common/transporttrust"

	"context"
	"errors"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryClientInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	return unaryClientInterceptorWithTrust(
		_transporttrust.ConfigFromEnv(),
		ctx, method, req, reply, cc, invoker, opts...,
	)
}

func unaryClientInterceptorWithTrust(
	trustConfig _transporttrust.Config,
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	start := time.Now()

	outgoingCtx, err := outgoingContextWithTrust(trustConfig, ctx, method)
	if err != nil {
		return err
	}

	err = invoker(outgoingCtx, method, req, reply, cc, opts...)
	log.Printf("[gRPC Client] Method: %s, Duration: %s, Error: %v", method, time.Since(start), err)
	return err
}

/**
 * StreamClientInterceptor propagates signed metadata to streaming calls.
 */
func StreamClientInterceptor(
	ctx context.Context,
	desc *grpc.StreamDesc,
	cc *grpc.ClientConn,
	method string,
	streamer grpc.Streamer,
	opts ...grpc.CallOption,
) (grpc.ClientStream, error) {
	return streamClientInterceptorWithTrust(
		_transporttrust.ConfigFromEnv(),
		ctx, desc, cc, method, streamer, opts...,
	)
}

func streamClientInterceptorWithTrust(
	trustConfig _transporttrust.Config,
	ctx context.Context,
	desc *grpc.StreamDesc,
	cc *grpc.ClientConn,
	method string,
	streamer grpc.Streamer,
	opts ...grpc.CallOption,
) (grpc.ClientStream, error) {
	start := time.Now()

	outgoingCtx, err := outgoingContextWithTrust(trustConfig, ctx, method)
	if err != nil {
		return nil, err
	}

	stream, err := streamer(outgoingCtx, desc, cc, method, opts...)
	log.Printf("[gRPC Stream Client] Method: %s, Duration: %s, Error: %v", method, time.Since(start), err)
	return stream, err
}

func outgoingContextWithTrust(
	trustConfig _transporttrust.Config,
	ctx context.Context,
	method string,
) (context.Context, error) {
	canonicalMD := GrpcClientMetadataFromContext(ctx)
	mergedMD := metadata.MD{}
	if existingMD, ok := metadata.FromOutgoingContext(ctx); ok {
		mergedMD = existingMD.Copy()
	}
	for key, values := range canonicalMD {
		mergedMD.Set(key, values...)
	}

	if err := _transporttrust.Sign(method, mergedMD, trustConfig); err != nil {
		if trustConfig.Mode == _transporttrust.ModeEnforce {
			return nil, status.Error(codes.FailedPrecondition, "internal metadata signing is unavailable")
		}
		if !errors.Is(err, _transporttrust.ErrNotConfigured) {
			log.Printf("[gRPC Client] Trusted metadata signing skipped: %v", err)
		}
	}

	return metadata.NewOutgoingContext(ctx, mergedMD), nil
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
