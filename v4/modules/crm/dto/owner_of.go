package dto

import (
	"fmt"
	"strconv"
	"strings"

	base_enum "base/enum"
	_errors "common/errors"
	sharepb "pb/types/shared"
)

// ParseOwnerOf reads a `{ownerOf}` path segment the way grpc-gateway binds sharepb.OwnerOf:
// either the enum name ("user", "organization", …) or its number ("10", "30", …).
func ParseOwnerOf(raw string) (base_enum.EOwnerOf, error) {
	raw = strings.TrimSpace(raw)
	if n, err := strconv.ParseInt(raw, 10, 32); err == nil {
		return base_enum.EOwnerOf(n), nil
	}
	if n, ok := sharepb.OwnerOf_value[strings.ToLower(raw)]; ok {
		return base_enum.EOwnerOf(n), nil
	}
	return 0, _errors.BadRequestException(fmt.Sprintf("invalid ownerOf: %q", raw))
}
