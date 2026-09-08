package tqdmultipart

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func firstNonEmpty(c *gin.Context, snake, camel string) string {
	if value := strings.TrimSpace(c.PostForm(snake)); value != "" {
		return value
	}
	return strings.TrimSpace(c.PostForm(camel))
}

func parseUint64(c *gin.Context, snake, camel string) (uint64, bool, error) {
	raw := firstNonEmpty(c, snake, camel)
	if raw == "" {
		return 0, false, nil
	}
	value, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, true, fmt.Errorf("%s invalid", snake)
	}
	return value, true, nil
}

func parseUint32(c *gin.Context, snake, camel string) (uint32, bool, error) {
	raw := firstNonEmpty(c, snake, camel)
	if raw == "" {
		return 0, false, nil
	}
	value, err := strconv.ParseUint(raw, 10, 32)
	if err != nil {
		return 0, true, fmt.Errorf("%s invalid", snake)
	}
	return uint32(value), true, nil
}

func parseBool(c *gin.Context, snake, camel string) bool {
	raw := strings.ToLower(firstNonEmpty(c, snake, camel))
	return raw == "1" || raw == "true" || raw == "yes" || raw == "on"
}
