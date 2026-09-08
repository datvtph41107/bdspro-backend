package httpmiddleware

import (
	"strings"

	"github.com/gin-contrib/cors"
)

func WithRequestIDCORS(
	config cors.Config,
) cors.Config {
	config.AllowHeaders =
		appendUniqueHeader(
			config.AllowHeaders,
			RequestIDHeader,
		)

	config.ExposeHeaders =
		appendUniqueHeader(
			config.ExposeHeaders,
			RequestIDHeader,
		)

	return config
}

/**
 * WithOperationIDCORS authorizes client propagation and response visibility
 * for the public Operation-ID contract.
 */
func WithOperationIDCORS(
	config cors.Config,
) cors.Config {
	config.AllowHeaders =
		appendUniqueHeader(
			config.AllowHeaders,
			OperationIDHeader,
		)

	config.ExposeHeaders =
		appendUniqueHeader(
			config.ExposeHeaders,
			OperationIDHeader,
		)

	return config
}

// WithIdempotencyKeyCORS authorizes caller propagation of the optional
// command identity. The key is not a response contract and is not exposed.
func WithIdempotencyKeyCORS(
	config cors.Config,
) cors.Config {
	config.AllowHeaders =
		appendUniqueHeader(
			config.AllowHeaders,
			IdempotencyKeyHeader,
		)

	return config
}

func appendUniqueHeader(
	headers []string,
	header string,
) []string {
	for _, existing := range headers {
		if strings.EqualFold(existing, header) {
			return headers
		}
	}

	return append(headers, header)
}
