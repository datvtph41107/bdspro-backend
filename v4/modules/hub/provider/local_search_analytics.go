package provider

import (
	hubsvc "mdn4/modules/hub/services"
	tqdprovider "mdn4/modules/tqd/provider"
)

// LocalSearchAnalytics adapts hub keyword/latency services for modules/tqd.
// @bind: mdn4/modules/tqd/provider.SearchAnalytics
type LocalSearchAnalytics struct {
	keywords hubsvc.KeywordSearchService
	latency  hubsvc.SearchLatencyService
}

func NewLocalSearchAnalytics(
	keywords hubsvc.KeywordSearchService,
	latency hubsvc.SearchLatencyService,
) *LocalSearchAnalytics {
	return &LocalSearchAnalytics{keywords: keywords, latency: latency}
}

var _ tqdprovider.SearchAnalytics = (*LocalSearchAnalytics)(nil)

func (a *LocalSearchAnalytics) TrackKeyword(keyword, source string) {
	if a == nil || a.keywords == nil {
		return
	}
	a.keywords.Track(keyword, source, 1)
}

func (a *LocalSearchAnalytics) TrackLatency(method, path, service string, durationMs int64) {
	if a == nil || a.latency == nil {
		return
	}
	a.latency.Track(method, path, service, durationMs)
}
