package requestcontext

import "net/http"

/**
 * IdempotencyKeyFromRequest parses the optional HTTP transport value.
 * It is read-only and never originates a key.
 */
func IdempotencyKeyFromRequest(
	request *http.Request,
) (
	string,
	bool,
	error,
) {
	if request == nil {
		panic("requestcontext: nil HTTP request")
	}

	return ParseIdempotencyKey(
		request.Header.Values(IdempotencyKeyHeader),
	)
}

/**
 * PrepareHTTPIdempotencyKey validates, binds, and canonicalizes an optional
 * caller-owned key without partially mutating the request on failure.
 */
func PrepareHTTPIdempotencyKey(
	request *http.Request,
) (
	string,
	bool,
	error,
) {
	key, ok, err := IdempotencyKeyFromRequest(request)
	if err != nil || !ok {
		return key, ok, err
	}

	boundContext, err := BindIdempotencyKey(
		request.Context(),
		key,
	)
	if err != nil {
		return "", false, err
	}

	preparedRequest := request.WithContext(boundContext)
	preparedRequest.Header = request.Header.Clone()
	if preparedRequest.Header == nil {
		preparedRequest.Header = make(http.Header)
	}

	preparedRequest.Header.Del(IdempotencyKeyHeader)
	preparedRequest.Header.Set(IdempotencyKeyHeader, key)

	*request = *preparedRequest

	return key, true, nil
}
