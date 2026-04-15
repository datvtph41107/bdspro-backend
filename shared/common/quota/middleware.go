package quota

import (
	"context"
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor(qm QuotaManager, extractUserID func(ctx context.Context) string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		userID := extractUserID(ctx)
		if userID == "" {
			// Thử lấy từ metadata nếu extractUserID không có
			if md, ok := metadata.FromIncomingContext(ctx); ok {
				for _, key := range []string{"profileId", "originId", "authId", "user_id"} {
					if vals := md.Get(key); len(vals) > 0 && vals[0] != "" {
						userID = vals[0]
						log.Printf("[QUOTA] Extracted userID from metadata %s: %s", key, userID)
						break
					}
				}
			}
		}

		if userID == "" {
			log.Printf("[QUOTA] No userID found, skipping quota for method %s", info.FullMethod)
			return handler(ctx, req)
		}

		service, action := parseMethod(info.FullMethod)
		if service == "" {
			return handler(ctx, req)
		}

		// 1. Check quota
		quotaInfo, err := qm.CheckOnly(ctx, userID, service, action)
		if err != nil {
			log.Printf("[QUOTA] CheckOnly error: %v", err)
			return handler(ctx, req) // fail open
		}

		if quotaInfo.Remaining <= 0 {
			LogReject(ctx, userID, service, action, quotaInfo.Pool, quotaInfo.Used, quotaInfo.Limit, "quota exceeded")
			return nil, status.Errorf(codes.ResourceExhausted, "Quota exceeded for pool %s. Used %d/%d", quotaInfo.Pool, quotaInfo.Used, quotaInfo.Limit)
		}

		// 2. Gọi handler
		resp, err := handler(ctx, req)

		// 3. Nếu thành công, consume quota
		if err == nil {
			result, consumeErr := qm.Consume(ctx, userID, service, action)
			if consumeErr != nil {
				log.Printf("[QUOTA] Failed to consume quota: %v", consumeErr)
			} else if result.Success {
				log.Printf("[QUOTA] Consumed: user=%s, pool=%s, used=%d, remaining=%d", userID, result.QuotaInfo.Pool, result.QuotaInfo.Used, result.QuotaInfo.Remaining)
			}
		}

		return resp, err
	}
}

func parseMethod(fullMethod string) (service, action string) {
	if len(fullMethod) < 2 || fullMethod[0] != '/' {
		return "", ""
	}
	fullMethod = fullMethod[1:]
	idx := strings.Index(fullMethod, "/")
	if idx == -1 {
		return "", ""
	}
	servicePart := fullMethod[:idx]
	actionPart := fullMethod[idx+1:]

	dotIdx := strings.Index(servicePart, ".")
	if dotIdx != -1 {
		service = servicePart[:dotIdx]
	} else {
		service = servicePart
	}

	action = toSnakeCase(actionPart)
	return
}

func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}
