package _errors

// type authErrorKeyType string

// const authErrorKey authErrorKeyType = "auth_error"

// func ErrorToGatewayResponse(ctx context.Context, w http.ResponseWriter, err error) {
// 	w.Header().Set("Content-Type", "application/json")

// 	// Nếu context có lỗi auth -> xử lý ở đây
// 	if ae, ok := ctx.Value(authErrorKey).(AuthError); ok {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		json.NewEncoder(w).Encode(map[string]string{
// 			"error": ae.Message,
// 		})
// 		return
// 	}

// 	// Nếu là lỗi gRPC bình thường
// 	grpcErr, ok := status.FromError(err)
// 	if !ok {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(w).Encode(map[string]string{
// 			"error": "Internal Server Error",
// 		})
// 		return
// 	}

// 	httpStatus := httpStatusFromCode(grpcErr.Code())
// 	w.WriteHeader(httpStatus)
// 	json.NewEncoder(w).Encode(map[string]string{
// 		"error": grpcErr.Message(),
// 	})
// }

// func httpStatusFromCode(code codes.Code) int {
// 	switch code {
// 	case codes.InvalidArgument:
// 		return http.StatusBadRequest
// 	case codes.Unauthenticated:
// 		return http.StatusUnauthorized
// 	case codes.PermissionDenied:
// 		return http.StatusForbidden
// 	case codes.NotFound:
// 		return http.StatusNotFound
// 	case codes.AlreadyExists:
// 		return http.StatusConflict
// 	case codes.Unimplemented:
// 		return http.StatusNotImplemented
// 	case codes.Unavailable:
// 		return http.StatusServiceUnavailable
// 	default:
// 		return http.StatusInternalServerError
// 	}
// }
