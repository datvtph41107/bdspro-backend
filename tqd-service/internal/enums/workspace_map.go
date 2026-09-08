package enums

import "strings"

// =====================================================
// WORKSPACE ENTITY TYPE
// =====================================================

type WorkspaceEntityType uint32

const (
	WorkspaceEntityTypeUnknown WorkspaceEntityType = 0

	WorkspaceEntityTypeParcel   WorkspaceEntityType = 1
	WorkspaceEntityTypeLocation WorkspaceEntityType = 2
	WorkspaceEntityTypeRegion   WorkspaceEntityType = 3
	WorkspaceEntityTypeReport   WorkspaceEntityType = 4
)

func (e WorkspaceEntityType) Uint32() uint32 {
	return uint32(e)
}

func (e WorkspaceEntityType) String() string {
	switch e {
	case WorkspaceEntityTypeParcel:
		return "parcel"
	case WorkspaceEntityTypeLocation:
		return "location"
	case WorkspaceEntityTypeRegion:
		return "region"
	case WorkspaceEntityTypeReport:
		return "report"
	default:
		return "unknown"
	}
}

func WorkspaceEntityTypeFromUint32(v uint32) WorkspaceEntityType {
	e := WorkspaceEntityType(v)
	if e.IsValid() {
		return e
	}
	return WorkspaceEntityTypeUnknown
}

func WorkspaceEntityTypeFromString(v string) WorkspaceEntityType {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "parcel":
		return WorkspaceEntityTypeParcel
	case "location", "point":
		return WorkspaceEntityTypeLocation
	case "region", "planning_region":
		return WorkspaceEntityTypeRegion
	case "report", "generated_report":
		return WorkspaceEntityTypeReport
	default:
		return WorkspaceEntityTypeUnknown
	}
}

func (e WorkspaceEntityType) IsValid() bool {
	switch e {
	case WorkspaceEntityTypeParcel,
		WorkspaceEntityTypeLocation,
		WorkspaceEntityTypeRegion,
		WorkspaceEntityTypeReport:
		return true
	default:
		return false
	}
}

func (e WorkspaceEntityType) SupportsViewHistory() bool {
	return e == WorkspaceEntityTypeParcel ||
		e == WorkspaceEntityTypeLocation ||
		e == WorkspaceEntityTypeRegion ||
		e == WorkspaceEntityTypeReport
}

func (e WorkspaceEntityType) SupportsFocus() bool {
	return e == WorkspaceEntityTypeParcel ||
		e == WorkspaceEntityTypeLocation ||
		e == WorkspaceEntityTypeRegion
}

func (e WorkspaceEntityType) SupportsFollow() bool {
	return e == WorkspaceEntityTypeParcel
}

func (e WorkspaceEntityType) SupportsGeneratedReportTarget() bool {
	return e == WorkspaceEntityTypeParcel ||
		e == WorkspaceEntityTypeRegion
}

// =====================================================
// VIEW HISTORY SOURCE
// =====================================================

type ViewHistorySource uint32

const (
	ViewHistorySourceSearch   ViewHistorySource = 1
	ViewHistorySourceMapClick ViewHistorySource = 2
	ViewHistorySourceReport   ViewHistorySource = 3
	ViewHistorySourceShare    ViewHistorySource = 4
	ViewHistorySourceDetail   ViewHistorySource = 5
	ViewHistorySourceFollow   ViewHistorySource = 6
	ViewHistorySourceHistory  ViewHistorySource = 7
	ViewHistorySourceUnknown  ViewHistorySource = 99
)

func (s ViewHistorySource) Uint32() uint32 {
	return uint32(s)
}

func (s ViewHistorySource) String() string {
	switch s {
	case ViewHistorySourceSearch:
		return "search"
	case ViewHistorySourceMapClick:
		return "mapClick"
	case ViewHistorySourceReport:
		return "report"
	case ViewHistorySourceShare:
		return "share"
	case ViewHistorySourceDetail:
		return "detail"
	case ViewHistorySourceFollow:
		return "follow"
	case ViewHistorySourceHistory:
		return "history"
	default:
		return "unknown"
	}
}

func ViewHistorySourceFromUint32(v uint32) ViewHistorySource {
	s := ViewHistorySource(v)
	if s.IsValid() {
		return s
	}
	return ViewHistorySourceUnknown
}

func (s ViewHistorySource) IsValid() bool {
	switch s {
	case ViewHistorySourceSearch,
		ViewHistorySourceMapClick,
		ViewHistorySourceReport,
		ViewHistorySourceShare,
		ViewHistorySourceDetail,
		ViewHistorySourceFollow,
		ViewHistorySourceHistory,
		ViewHistorySourceUnknown:
		return true
	default:
		return false
	}
}

func (s ViewHistorySource) CanCountView() bool {
	return s != ViewHistorySourceUnknown
}

// =====================================================
// GENERATED REPORT TYPE
// =====================================================

type GeneratedReportType uint32

const (
	GeneratedReportTypeUnknown    GeneratedReportType = 0
	GeneratedReportTypePlanning   GeneratedReportType = 1
	GeneratedReportTypeSnapshot   GeneratedReportType = 2
	GeneratedReportTypeComparison GeneratedReportType = 3
	GeneratedReportTypeAnalysis   GeneratedReportType = 4
)

func (t GeneratedReportType) Uint32() uint32 {
	return uint32(t)
}

func (t GeneratedReportType) String() string {
	switch t {
	case GeneratedReportTypePlanning:
		return "planning"
	case GeneratedReportTypeSnapshot:
		return "snapshot"
	case GeneratedReportTypeComparison:
		return "comparison"
	case GeneratedReportTypeAnalysis:
		return "analysis"
	default:
		return "unknown"
	}
}

func (t GeneratedReportType) IsValid() bool {
	switch t {
	case GeneratedReportTypePlanning,
		GeneratedReportTypeSnapshot,
		GeneratedReportTypeComparison,
		GeneratedReportTypeAnalysis:
		return true
	default:
		return false
	}
}

func (t GeneratedReportType) RequiresComparisonInput() bool {
	return t == GeneratedReportTypeComparison
}

func (t GeneratedReportType) SupportsEntityType(entity WorkspaceEntityType) bool {
	switch t {
	case GeneratedReportTypePlanning,
		GeneratedReportTypeComparison,
		GeneratedReportTypeAnalysis:
		return entity == WorkspaceEntityTypeParcel ||
			entity == WorkspaceEntityTypeRegion
	case GeneratedReportTypeSnapshot:
		return entity == WorkspaceEntityTypeParcel ||
			entity == WorkspaceEntityTypeRegion ||
			entity == WorkspaceEntityTypeLocation
	default:
		return false
	}
}

func (t GeneratedReportType) DefaultFormat() string {
	switch t {
	case GeneratedReportTypeSnapshot,
		GeneratedReportTypePlanning,
		GeneratedReportTypeComparison,
		GeneratedReportTypeAnalysis:
		// Both canonical renderers currently produce a PDF. Format describes
		// the persisted artifact, not a report-type presentation hint.
		return "pdf"
	default:
		return ""
	}
}

// =====================================================
// GENERATED REPORT STATUS
// =====================================================

type GeneratedReportStatus uint32

const (
	GeneratedReportStatusUnknown GeneratedReportStatus = 0

	GeneratedReportStatusReady      GeneratedReportStatus = 10
	GeneratedReportStatusProcessing GeneratedReportStatus = 20
	GeneratedReportStatusExpired    GeneratedReportStatus = 30
	GeneratedReportStatusFailed     GeneratedReportStatus = 40
)

func (s GeneratedReportStatus) Uint32() uint32 {
	return uint32(s)
}

func (s GeneratedReportStatus) String() string {
	switch s {
	case GeneratedReportStatusReady:
		return "ready"
	case GeneratedReportStatusProcessing:
		return "processing"
	case GeneratedReportStatusExpired:
		return "expired"
	case GeneratedReportStatusFailed:
		return "failed"
	default:
		return "unknown"
	}
}

func (s GeneratedReportStatus) IsValid() bool {
	switch s {
	case GeneratedReportStatusReady,
		GeneratedReportStatusProcessing,
		GeneratedReportStatusExpired,
		GeneratedReportStatusFailed:
		return true
	default:
		return false
	}
}

func (s GeneratedReportStatus) IsInProgress() bool {
	return s == GeneratedReportStatusProcessing
}

func (s GeneratedReportStatus) IsTerminal() bool {
	return s == GeneratedReportStatusReady ||
		s == GeneratedReportStatusExpired ||
		s == GeneratedReportStatusFailed
}

func (s GeneratedReportStatus) CanDownload() bool {
	return s == GeneratedReportStatusReady
}

func (s GeneratedReportStatus) CanShare() bool {
	return s == GeneratedReportStatusReady
}

func (s GeneratedReportStatus) CanRegenerate() bool {
	return s == GeneratedReportStatusReady ||
		s == GeneratedReportStatusExpired ||
		s == GeneratedReportStatusFailed
}

func (s GeneratedReportStatus) CanRemove() bool {
	return s.IsValid()
}

func (s GeneratedReportStatus) ShouldPoll() bool {
	return s == GeneratedReportStatusProcessing
}

func (s GeneratedReportStatus) CanTransitionTo(next GeneratedReportStatus) bool {
	if !next.IsValid() {
		return false
	}

	switch s {
	case GeneratedReportStatusUnknown:
		return next == GeneratedReportStatusProcessing ||
			next == GeneratedReportStatusFailed
	case GeneratedReportStatusProcessing:
		return next == GeneratedReportStatusReady ||
			next == GeneratedReportStatusFailed ||
			next == GeneratedReportStatusExpired
	case GeneratedReportStatusReady:
		return next == GeneratedReportStatusProcessing ||
			next == GeneratedReportStatusExpired
	case GeneratedReportStatusExpired:
		return next == GeneratedReportStatusProcessing
	case GeneratedReportStatusFailed:
		return next == GeneratedReportStatusProcessing
	default:
		return false
	}
}

// =====================================================
// PREVIEW MODE
// =====================================================

type PreviewMode uint32

const (
	PreviewModeUnknown   PreviewMode = 0
	PreviewModeGeometry  PreviewMode = 1
	PreviewModeThumbnail PreviewMode = 2
	PreviewModeBounds    PreviewMode = 3
	PreviewModeStatus    PreviewMode = 4
)

func (m PreviewMode) Uint32() uint32 {
	return uint32(m)
}

func (m PreviewMode) String() string {
	switch m {
	case PreviewModeGeometry:
		return "geometry"
	case PreviewModeThumbnail:
		return "thumbnail"
	case PreviewModeBounds:
		return "bounds"
	case PreviewModeStatus:
		return "status"
	default:
		return "unknown"
	}
}

// =====================================================
// FOCUS TYPE
// =====================================================

type FocusType uint32

const (
	FocusTypeUnknown  FocusType = 0
	FocusTypeBounds   FocusType = 1
	FocusTypePoint    FocusType = 2
	FocusTypeGeometry FocusType = 3
)

func (f FocusType) Uint32() uint32 {
	return uint32(f)
}

func (f FocusType) String() string {
	switch f {
	case FocusTypeBounds:
		return "bounds"
	case FocusTypePoint:
		return "point"
	case FocusTypeGeometry:
		return "geometry"
	default:
		return "unknown"
	}
}
