package rpc

import (
	"sort"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	RequestIDMetadataKey      = "x-request-id"
	OperationIDMetadataKey    = "x-operation-id"
	IdempotencyKeyMetadataKey = "idempotency-key"

	CallerKindMetadataKey               = "x-qhpro-caller-kind"
	APIKeyIDMetadataKey                 = "x-qhpro-api-key-id"
	APIKeyAppMetadataKey                = "x-qhpro-api-key-app"
	ProviderCredentialDigestMetadataKey = "x-qhpro-provider-credential-sha256"

	AuthIDMetadataKey         = "authid"
	ProfileIDMetadataKey      = "profileid"
	OriginIDMetadataKey       = "originid"
	OrganizationIDMetadataKey = "organizationid"
	SessionIDMetadataKey      = "session"
	RoleMetadataKey           = "role"
	TokenTypeMetadataKey      = "type"
)

type metadataContract struct {
	privilegedKeys  map[string]struct{}
	signedKeys      []string
	unsignedMethods map[string]struct{}
}

func newMetadataContract(
	privilegedKeys []string,
	signedKeys []string,
	unsignedMethods []string,
) metadataContract {
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

	// Assertion fields are always protected by the assertion itself.
	for _, key := range serviceAssertionMetadataKeys() {
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

	return metadataContract{
		privilegedKeys:  privileged,
		signedKeys:      signed,
		unsignedMethods: unsigned,
	}
}

func (c metadataContract) hasPrivilegedMetadata(md metadata.MD) bool {
	for key := range c.privilegedKeys {
		if len(md.Get(key)) > 0 {
			return true
		}
	}
	return false
}

func (c metadataContract) allowsUnsigned(method string) bool {
	_, ok := c.unsignedMethods[strings.TrimSpace(method)]
	return ok
}

func normalizeMetadataKey(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

var qhproMetadataContract = newMetadataContract(
	[]string{
		CallerKindMetadataKey,
		APIKeyIDMetadataKey,
		APIKeyAppMetadataKey,
		ProviderCredentialDigestMetadataKey,
		"x-api-key",
		AuthIDMetadataKey,
		ProfileIDMetadataKey,
		OriginIDMetadataKey,
		OrganizationIDMetadataKey,
		SessionIDMetadataKey,
		RoleMetadataKey,
		TokenTypeMetadataKey,
		"planid",
		"planfrom",
	},
	[]string{
		CallerKindMetadataKey,
		APIKeyIDMetadataKey,
		APIKeyAppMetadataKey,
		ProviderCredentialDigestMetadataKey,
		"x-api-key",
		AuthIDMetadataKey,
		ProfileIDMetadataKey,
		OriginIDMetadataKey,
		OrganizationIDMetadataKey,
		SessionIDMetadataKey,
		RoleMetadataKey,
		TokenTypeMetadataKey,
		"planid",
		"planfrom",
		RequestIDMetadataKey,
		OperationIDMetadataKey,
		IdempotencyKeyMetadataKey,
	},
	[]string{
		"/grpc.health.v1.Health/Check",
		"/grpc.health.v1.Health/Watch",
	},
)
