package rpc

import "common/trust"

const (
	RequestIDMetadataKey      = "x-request-id"
	OperationIDMetadataKey    = "x-operation-id"
	IdempotencyKeyMetadataKey = "idempotency-key"

	CallerKindMetadataKey = "x-qhpro-caller-kind"
	APIKeyIDMetadataKey   = "x-qhpro-api-key-id"
	APIKeyAppMetadataKey  = "x-qhpro-api-key-app"

	AuthIDMetadataKey         = "authid"
	ProfileIDMetadataKey      = "profileid"
	OriginIDMetadataKey       = "originid"
	OrganizationIDMetadataKey = "organizationid"
	SessionIDMetadataKey      = "session"
	RoleMetadataKey           = "role"
	TokenTypeMetadataKey      = "type"
)

var qhproMetadataContract = trust.NewMetadataContract(
	[]string{
		CallerKindMetadataKey,
		APIKeyIDMetadataKey,
		APIKeyAppMetadataKey,
		"x-qhpro-quota-subject-kind",
		"x-qhpro-quota-subject-id",
		"x-qhpro-quota-subject-source",
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
		"x-qhpro-quota-subject-kind",
		"x-qhpro-quota-subject-id",
		"x-qhpro-quota-subject-source",
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

// QHPROMetadataContract returns the immutable application metadata contract
// protected by the canonical QHPRO RPC transport. The trust package protects
// the supplied keys without interpreting their business meaning.
func QHPROMetadataContract() trust.MetadataContract {
	return qhproMetadataContract
}
