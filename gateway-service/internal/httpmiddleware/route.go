package httpmiddleware

import "github.com/gin-gonic/gin"

const unmatchedRoute = "__unmatched__"

// normalizedRoute returns the bounded route identity used by
// Gateway HTTP telemetry.
//
// Matched requests use Gin's route template.
// Unmatched requests collapse to one sentinel instead of
// exposing raw URL paths as telemetry dimensions.
func normalizedRoute(
	ctx *gin.Context,
) string {
	route := ctx.FullPath()

	if route == "" {
		return unmatchedRoute
	}

	return route
}
