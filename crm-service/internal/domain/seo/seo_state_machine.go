package seo_domain

import "time"

func (s *SeoDomain) Publish(now time.Time) {
	if s == nil {
		return
	}
	s.PageStatus = SeoPageStatusPublished
	s.Published = true
	if s.PublishedAt == nil || s.PublishedAt.IsZero() {
		s.PublishedAt = &now
	}
	s.MarkRenderPending()
}

func (s *SeoDomain) Unpublish() {
	if s == nil {
		return
	}
	s.PageStatus = SeoPageStatusDraft
	s.Published = false
	s.IsSiteMap = false
	s.MarkRenderPending()
}

func (s *SeoDomain) Archive() {
	if s == nil {
		return
	}
	s.PageStatus = SeoPageStatusArchived
	s.Published = false
	s.IsSiteMap = false
	s.IsIndex = false
	s.NeedGenerate = false
	s.RenderStatus = SeoRenderStatusNone
	s.LastRenderError = ""
}

func (s *SeoDomain) Restore() {
	if s == nil {
		return
	}
	s.PageStatus = SeoPageStatusDraft
	s.Published = false
	s.IsSiteMap = false
	s.MarkRenderPending()
}

func (s *SeoDomain) MarkSourceStale(sourceUpdatedAt time.Time) {
	if s == nil {
		return
	}
	s.SourceStatus = SeoSourceStatusStale
	s.SourceUpdatedAt = &sourceUpdatedAt
	s.MarkRenderPending()
}

func (s *SeoDomain) MarkSourceMissing(sourceUpdatedAt time.Time) {
	if s == nil {
		return
	}
	s.SourceStatus = SeoSourceStatusMissing
	s.RefMissing = true
	s.SourceUpdatedAt = &sourceUpdatedAt
	s.MarkRenderPending()
}
