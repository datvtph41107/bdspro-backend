package enums

type WorkspaceEntityType uint32

const (
	WorkspaceEntityTypeParcel   WorkspaceEntityType = 1
	WorkspaceEntityTypeLocation WorkspaceEntityType = 2
	WorkspaceEntityTypeReport   WorkspaceEntityType = 3
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
	case WorkspaceEntityTypeReport:
		return "report"
	default:
		return "unknown"
	}
}

type ViewHistorySource uint32

const (
	ViewHistorySourceSearch    ViewHistorySource = 1
	ViewHistorySourceMapClick  ViewHistorySource = 2
	ViewHistorySourceReport    ViewHistorySource = 3
	ViewHistorySourceShareLink ViewHistorySource = 4
	ViewHistorySourceUnknown   ViewHistorySource = 99
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
	case ViewHistorySourceShareLink:
		return "shareLink"
	case ViewHistorySourceUnknown:
		return "unknown"
	default:
		return "unknown"
	}
}

type GeneratedReportType uint32

const (
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

type GeneratedReportStatus uint32

const (
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

type PreviewMode uint32

const (
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

type FocusType uint32

const (
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
