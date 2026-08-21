package trust

import (
	"sort"
	"strings"

	"google.golang.org/grpc/metadata"
)

// MetadataContract declares which application metadata fields participate in
// transport trust. The trust package does not interpret those fields; it only
// protects the exact keys supplied by the transport adapter.
type MetadataContract struct {
	privilegedKeys  map[string]struct{}
	signedKeys      []string
	unsignedMethods map[string]struct{}
}

// NewMetadataContract creates an immutable transport metadata contract.
func NewMetadataContract(
	privilegedKeys []string,
	signedKeys []string,
	unsignedMethods []string,
) MetadataContract {
	privileged := make(map[string]struct{}, len(privilegedKeys))
	for _, key := range privilegedKeys {
		if key = normalizeMetadataKey(key); key != "" {
			privileged[key] = struct{}{}
		}
	}

	signedSet := make(map[string]struct{}, len(signedKeys))
	for _, key := range signedKeys {
		if key = normalizeMetadataKey(key); key != "" {
			signedSet[key] = struct{}{}
		}
	}
	// The assertion fields themselves are always signed. Callers should only
	// provide application metadata keys.
	for _, key := range assertionMetadataKeys {
		signedSet[key] = struct{}{}
	}

	signed := make([]string, 0, len(signedSet))
	for key := range signedSet {
		signed = append(signed, key)
	}
	sort.Strings(signed)

	unsigned := make(map[string]struct{}, len(unsignedMethods))
	for _, method := range unsignedMethods {
		if method = strings.TrimSpace(method); method != "" {
			unsigned[method] = struct{}{}
		}
	}

	return MetadataContract{
		privilegedKeys:  privileged,
		signedKeys:      signed,
		unsignedMethods: unsigned,
	}
}

func (c MetadataContract) hasPrivilegedMetadata(md metadata.MD) bool {
	for key := range c.privilegedKeys {
		if len(md.Get(key)) > 0 {
			return true
		}
	}
	return false
}

func (c MetadataContract) allowsUnsignedMethod(method string) bool {
	_, ok := c.unsignedMethods[strings.TrimSpace(method)]
	return ok
}

func normalizeMetadataKey(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}
