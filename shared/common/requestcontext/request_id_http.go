package requestcontext

import "net/http"

// EnsureRequestIDDecision canonicalize HTTP request header
func EnsureRequestIDDecision(
	request *http.Request,
) RequestIDDecision {
	if request == nil {
		return ResolveRequestID(nil)
	}

	if request.Header == nil {
		request.Header = make(http.Header)
	}

	// request.Header.Values -> return []string{"abc", "def"}
	/**
	 * Description: Header.Values / Header.Get
	 * Header.Get chỉ cung cấp một string / header nhiều values
	 ** X-Request-ID: req_A
	 ** X-Request-ID: req_B
	 */
	decision := ResolveRequestID(
		// request.Header.Values("X-Request-ID")
		// nil               không có
		// []string{"A"}     một
		// []string{"A","B"} nhiều
		request.Header.Values(
			RequestIDHeader,
		),
	)

	// Xóa values cũ tránh duplicate/stale headers
	/**
	 ** X-Request-ID: req_A
	 ** X-Request-ID: req_B
	 */
	request.Header.Del(RequestIDHeader)

	// Set lại 1 canonical value
	request.Header.Set(
		RequestIDHeader,
		decision.ID,
	)

	return decision
}

// EnsureRequestID giữ API tương thích source
func EnsureRequestID(
	request *http.Request,
) string {
	return EnsureRequestIDDecision(request).ID
}

// SetResponseHeader trả canonical request ID cho client
func SetResponseHeader(
	writer http.ResponseWriter,
	requestID string,
) {
	if writer == nil {
		return
	}

	if !IsValidRequestID(requestID) {
		return
	}

	writer.Header().Set(
		RequestIDHeader,
		requestID,
	)
}

// PrepareHttpRequestDecision bind canonical
// vào header req, res header và req context
func PrepareHTTPRequestDecision(
	writer http.ResponseWriter,
	request *http.Request,
) RequestIDDecision {
	if request == nil {
		panic("requestcontext: nil HTTP request")
	}

	decision := EnsureRequestIDDecision(request)

	SetResponseHeader(
		writer,
		decision.ID,
	)

	preparedRequest := request.WithContext(
		WithRequestID(
			request.Context(),
			decision.ID,
		),
	)
	// giữ parent context
	// thay value bằng request có context mới.
	*request = *preparedRequest

	return decision
}

func PrepareHTTPRequest(
	writer http.ResponseWriter,
	request *http.Request,
) string {
	return PrepareHTTPRequestDecision(writer, request).ID
}

// gán *reqd = *prepare -> đổi từ request cha truyền vào params không gán lại biến bên ngoài sử dụng
