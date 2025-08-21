package middleware

import (
	_utils "common/utils"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AdminAuthInterceptor middleware để kiểm tra quyền admin
func AdminAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Chỉ áp dụng cho admin APIs
	if isAdminAPI(info.FullMethod) {
		// Lấy profile ID từ context
		profileID := _utils.GetProfileIdWithContext(ctx)
		if profileID == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "unauthorized: missing profile ID")
		}

		// TODO: Kiểm tra quyền admin từ database hoặc cache
		// Hiện tại tạm thời cho phép tất cả user có profile ID
		// Trong thực tế cần kiểm tra role từ database
		if !isAdmin(profileID) {
			return nil, status.Errorf(codes.PermissionDenied, "unauthorized: admin access required")
		}
	}

	return handler(ctx, req)
}

// isAdminAPI kiểm tra xem API có phải là admin API không
func isAdminAPI(method string) bool {
	adminMethods := []string{
		"/userpb.ProfileService/ListAllUsers",
		// Thêm các admin API khác ở đây
	}

	for _, adminMethod := range adminMethods {
		if method == adminMethod {
			return true
		}
	}
	return false
}

// isAdmin kiểm tra xem user có quyền admin không
// TODO: Implement logic kiểm tra quyền admin từ database
func isAdmin(profileID uint64) bool {
	// Tạm thời cho phép tất cả user có profile ID > 0
	// Trong thực tế cần kiểm tra role từ database
	return profileID > 0
}
