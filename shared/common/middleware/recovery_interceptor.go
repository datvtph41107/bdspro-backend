package _middleware

import (
	"context"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorLogger is the minimal logger dependency for panic recovery.
type ErrorLogger interface {
	Errorf(format string, args ...any)
}

// StdLogger adapts log.Logger to ErrorLogger.
type StdLogger struct {
	*log.Logger
}

func (l StdLogger) Errorf(format string, args ...any) {
	if l.Logger == nil {
		log.Printf("ERROR: "+format, args...)
		return
	}
	l.Logger.Printf("ERROR: "+format, args...)
}

func NewStdLogger(name string) StdLogger {
	if name == "" {
		name = "grpc"
	}
	return StdLogger{log.New(os.Stdout, name+": ", log.LstdFlags)}
}

func UnaryRecoveryInterceptor(logger ErrorLogger) grpc.UnaryServerInterceptor {
	if logger == nil {
		logger = StdLogger{}
	}
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Errorf("panic recovered in unary interceptor method: %s, panic: %v", info.FullMethod, r)
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()
		return handler(ctx, req)
	}
}
