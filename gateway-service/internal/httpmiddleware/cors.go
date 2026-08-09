package httpmiddleware

import (
	_requestcontext "common/requestcontext"
	"strings"

	"github.com/gin-contrib/cors"
)

func WithRequestIDCORS(
	config cors.Config,
) cors.Config {
	config.AllowHeaders =
		appendUniqueHeader(
			config.AllowHeaders,
			_requestcontext.RequestIDHeader,
		)

	config.ExposeHeaders =
		appendUniqueHeader(
			config.ExposeHeaders,
			_requestcontext.RequestIDHeader,
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
