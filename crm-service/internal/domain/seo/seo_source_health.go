package seo_domain

// SeoSourceHealth represents runtime health information for a SEO source.
//
// This domain object is used by the SEO Source Registry / Source Health API
// to expose whether each SEO source is ready, degraded, disabled, or still
// coming soon.
//
// Example sources:
//   - manual
//   - adm_unit
//   - region
//   - parcel
//   - project
//   - document
//   - map
//   - report
//
// It is intentionally kept in the domain/seo package because source health
// is part of the SEO operating model, not a transport-only DTO.
type SeoSourceHealth struct {
	// SourceKey is the public key of the SEO source.
	//
	// Examples:
	//   - manual
	//   - project
	//   - parcel
	//   - region
	//   - adm_unit
	//   - document
	//   - map
	//   - report
	SourceKey string

	// SourceService is the upstream service that owns or provides the source data.
	//
	// Examples:
	//   - crm
	//   - bdspro
	//   - tqd
	//   - manual
	SourceService string

	// ResolverKey identifies which resolver implementation should be used for
	// searching/resolving this source.
	//
	// Examples:
	//   - manual
	//   - project
	//   - parcel
	//   - region
	//   - adm_unit
	ResolverKey string

	// Status describes the runtime/registry status of the source.
	//
	// Expected values usually align with SEO source registry statuses:
	//   - ready
	//   - coming_soon
	//   - disabled
	//   - degraded
	//   - error
	Status string

	// LatencyMs is the measured latency in milliseconds when checking or resolving
	// this source. It may be 0 when no active health check has been performed.
	LatencyMs int64

	// LastCheckedAt is an ISO/RFC3339 timestamp string indicating when this health
	// information was last checked.
	//
	// It is a string instead of time.Time to keep mapper/transport conversion simple
	// and consistent with existing DTO/protobuf response patterns.
	LastCheckedAt string

	// DisabledReason explains why the source is disabled or coming soon.
	//
	// Examples:
	//   - resolver_not_ready
	//   - template_not_ready
	//   - source_service_not_integrated
	//   - admin_disabled
	DisabledReason string

	// ErrorMessage contains the latest error message if the source health is
	// degraded or failed.
	ErrorMessage string
}
