package provider

// SearchAnalytics records text-search keywords / latency into hub tables (db_v4).
// Nil-safe: callers may invoke on a nil receiver via NoopSearchAnalytics.
type SearchAnalytics interface {
	TrackKeyword(keyword, source string)
	TrackLatency(method, path, service string, durationMs int64)
}

// NoopSearchAnalytics is used when analytics wiring is unavailable.
type NoopSearchAnalytics struct{}

func (NoopSearchAnalytics) TrackKeyword(string, string)                  {}
func (NoopSearchAnalytics) TrackLatency(string, string, string, int64) {}
