package seo_strategy

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type QueryEvidenceSourceKind string

const (
	QueryEvidenceSourceFirstParty  QueryEvidenceSourceKind = "first_party"
	QueryEvidenceSourceThirdParty  QueryEvidenceSourceKind = "third_party"
	QueryEvidenceSourceQualitative QueryEvidenceSourceKind = "qualitative"
	QueryEvidenceSourceManual      QueryEvidenceSourceKind = "manual"
)

type QueryObservationType string

const (
	QueryObservationObserved  QueryObservationType = "observed"
	QueryObservationReported  QueryObservationType = "reported"
	QueryObservationSuggested QueryObservationType = "suggested"
	QueryObservationAssumed   QueryObservationType = "assumed"
)

type QueryDeviceCategory string

const (
	QueryDeviceUnknown QueryDeviceCategory = "unknown"
	QueryDeviceMobile  QueryDeviceCategory = "mobile"
	QueryDeviceDesktop QueryDeviceCategory = "desktop"
	QueryDeviceTablet  QueryDeviceCategory = "tablet"
	QueryDeviceAll     QueryDeviceCategory = "all"
)

type QueryEvidenceStatus string

const (
	QueryEvidenceStatusActive     QueryEvidenceStatus = "active"
	QueryEvidenceStatusSuperseded QueryEvidenceStatus = "superseded"
	QueryEvidenceStatusRejected   QueryEvidenceStatus = "rejected"
)

var (
	queryEvidenceSourceSystemPattern = regexp.MustCompile(`^[a-z][a-z0-9_]{1,63}$`)
	queryEvidenceCountryPattern      = regexp.MustCompile(`^[A-Z]{2}$`)
)

type SeoQueryEvidence struct {
	ID uint64

	ObjectiveID uint64
	UserNeedID  uint64

	RawQuery        string
	NormalizedQuery string
	QueryHash       string

	SourceSystem    string
	SourceKind      QueryEvidenceSourceKind
	ObservationType QueryObservationType
	SourceRecordKey string
	SourceReference string

	Locale         string
	CountryCode    string
	DeviceCategory QueryDeviceCategory

	ObservedFrom *time.Time
	ObservedTo   *time.Time
	ValidUntil   *time.Time

	Impressions  *int64
	Clicks       *int64
	CTR          *float64
	AvgPosition  *float64
	SearchVolume *int64

	Confidence ObjectiveConfidence
	Status     QueryEvidenceStatus

	Metadata json.RawMessage

	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateSeoQueryEvidence là domain factory runtime, không phải Wire provider.
func CreateSeoQueryEvidence(
	objectiveID uint64,
	userNeedID uint64,
	rawQuery string,
	sourceSystem string,
	sourceKind QueryEvidenceSourceKind,
	observationType QueryObservationType,
	sourceRecordKey string,
	now time.Time,
) (*SeoQueryEvidence, error) {
	rawQuery = strings.TrimSpace(rawQuery)
	normalizedQuery := NormalizeSearchQuery(rawQuery)

	evidence := &SeoQueryEvidence{
		ObjectiveID:     objectiveID,
		UserNeedID:      userNeedID,
		RawQuery:        rawQuery,
		NormalizedQuery: normalizedQuery,
		SourceSystem:    strings.TrimSpace(sourceSystem),
		SourceKind:      sourceKind,
		ObservationType: observationType,
		SourceRecordKey: strings.TrimSpace(sourceRecordKey),
		Locale:          "vi-VN",
		CountryCode:     "VN",
		DeviceCategory:  QueryDeviceUnknown,
		Confidence:      ObjectiveConfidenceExploratory,
		Status:          QueryEvidenceStatusActive,
		Metadata:        json.RawMessage(`{}`),
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	evidence.QueryHash = BuildQueryEvidenceHash(
		evidence.NormalizedQuery,
		evidence.Locale,
	)

	if err := ValidateQueryEvidence(evidence); err != nil {
		return nil, err
	}

	return evidence, nil
}

// NormalizeSearchQuery chỉ phục vụ lookup/deduplicate.
// RawQuery vẫn phải được giữ nguyên để bảo toàn evidence gốc.
func NormalizeSearchQuery(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ToLower(value)
	return strings.Join(strings.Fields(value), " ")
}

// BuildQueryEvidenceHash tạo fingerprint ổn định từ query + locale.
func BuildQueryEvidenceHash(normalizedQuery, locale string) string {
	payload := NormalizeSearchQuery(normalizedQuery) +
		"\x1f" +
		strings.ToLower(strings.TrimSpace(locale))

	sum := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(sum[:])
}

func ValidateQueryEvidence(evidence *SeoQueryEvidence) error {
	if evidence == nil {
		return errors.New("query evidence không tồn tại")
	}
	if evidence.ObjectiveID == 0 {
		return errors.New("objectiveId không hợp lệ")
	}
	if evidence.UserNeedID == 0 {
		return errors.New("userNeedId không hợp lệ")
	}

	evidence.RawQuery = strings.TrimSpace(evidence.RawQuery)
	if evidence.RawQuery == "" {
		return errors.New("raw query không được để trống")
	}
	if len([]rune(evidence.RawQuery)) > 1000 {
		return errors.New("raw query không được dài quá 1000 ký tự")
	}

	evidence.NormalizedQuery = NormalizeSearchQuery(evidence.RawQuery)
	if evidence.NormalizedQuery == "" {
		return errors.New("normalized query không được để trống")
	}

	evidence.SourceSystem = strings.TrimSpace(evidence.SourceSystem)
	if !queryEvidenceSourceSystemPattern.MatchString(evidence.SourceSystem) {
		return errors.New("source system phải là lowercase snake_case, dài từ 2 đến 64 ký tự")
	}
	if !IsValidQueryEvidenceSourceKind(evidence.SourceKind) {
		return fmt.Errorf("query evidence source kind không hợp lệ: %s", evidence.SourceKind)
	}
	if !IsValidQueryObservationType(evidence.ObservationType) {
		return fmt.Errorf("query observation type không hợp lệ: %s", evidence.ObservationType)
	}

	evidence.SourceRecordKey = strings.TrimSpace(evidence.SourceRecordKey)
	if evidence.SourceRecordKey == "" {
		return errors.New("source record key không được để trống")
	}
	if len([]rune(evidence.SourceRecordKey)) > 200 {
		return errors.New("source record key không được dài quá 200 ký tự")
	}
	evidence.SourceReference = strings.TrimSpace(evidence.SourceReference)

	evidence.Locale = strings.TrimSpace(evidence.Locale)
	if evidence.Locale == "" {
		return errors.New("locale không được để trống")
	}
	if len(evidence.Locale) > 16 {
		return errors.New("locale không được dài quá 16 ký tự")
	}

	evidence.CountryCode = strings.ToUpper(strings.TrimSpace(evidence.CountryCode))
	if !queryEvidenceCountryPattern.MatchString(evidence.CountryCode) {
		return errors.New("country code phải gồm đúng 2 chữ cái in hoa")
	}
	if !IsValidQueryDeviceCategory(evidence.DeviceCategory) {
		return fmt.Errorf("device category không hợp lệ: %s", evidence.DeviceCategory)
	}

	expectedHash := BuildQueryEvidenceHash(evidence.NormalizedQuery, evidence.Locale)
	if evidence.QueryHash == "" {
		evidence.QueryHash = expectedHash
	} else if evidence.QueryHash != expectedHash {
		return errors.New("query hash không khớp normalized query và locale")
	}

	if evidence.ObservedFrom != nil &&
		evidence.ObservedTo != nil &&
		evidence.ObservedTo.Before(*evidence.ObservedFrom) {
		return errors.New("observedTo không được trước observedFrom")
	}
	if evidence.Impressions != nil && *evidence.Impressions < 0 {
		return errors.New("impressions không được âm")
	}
	if evidence.Clicks != nil && *evidence.Clicks < 0 {
		return errors.New("clicks không được âm")
	}
	if evidence.Impressions != nil &&
		evidence.Clicks != nil &&
		*evidence.Clicks > *evidence.Impressions {
		return errors.New("clicks không được lớn hơn impressions")
	}
	if evidence.CTR != nil && (*evidence.CTR < 0 || *evidence.CTR > 1) {
		return errors.New("ctr phải nằm trong khoảng 0 đến 1")
	}
	if evidence.AvgPosition != nil && *evidence.AvgPosition <= 0 {
		return errors.New("avg position phải lớn hơn 0")
	}
	if evidence.SearchVolume != nil && *evidence.SearchVolume < 0 {
		return errors.New("search volume không được âm")
	}
	if !IsValidObjectiveConfidence(evidence.Confidence) {
		return fmt.Errorf("query evidence confidence không hợp lệ: %s", evidence.Confidence)
	}
	if !IsValidQueryEvidenceStatus(evidence.Status) {
		return fmt.Errorf("query evidence status không hợp lệ: %s", evidence.Status)
	}

	if len(evidence.Metadata) == 0 {
		evidence.Metadata = json.RawMessage(`{}`)
	}
	if !json.Valid(evidence.Metadata) {
		return errors.New("query evidence metadata không phải JSON hợp lệ")
	}

	return nil
}

func IsValidQueryEvidenceSourceKind(value QueryEvidenceSourceKind) bool {
	switch value {
	case QueryEvidenceSourceFirstParty,
		QueryEvidenceSourceThirdParty,
		QueryEvidenceSourceQualitative,
		QueryEvidenceSourceManual:
		return true
	default:
		return false
	}
}

func IsValidQueryObservationType(value QueryObservationType) bool {
	switch value {
	case QueryObservationObserved,
		QueryObservationReported,
		QueryObservationSuggested,
		QueryObservationAssumed:
		return true
	default:
		return false
	}
}

func IsValidQueryDeviceCategory(value QueryDeviceCategory) bool {
	switch value {
	case QueryDeviceUnknown,
		QueryDeviceMobile,
		QueryDeviceDesktop,
		QueryDeviceTablet,
		QueryDeviceAll:
		return true
	default:
		return false
	}
}

func IsValidQueryEvidenceStatus(value QueryEvidenceStatus) bool {
	switch value {
	case QueryEvidenceStatusActive,
		QueryEvidenceStatusSuperseded,
		QueryEvidenceStatusRejected:
		return true
	default:
		return false
	}
}

func ParseQueryEvidenceSourceKind(value string) (QueryEvidenceSourceKind, error) {
	result := QueryEvidenceSourceKind(strings.TrimSpace(value))
	if !IsValidQueryEvidenceSourceKind(result) {
		return "", fmt.Errorf("query evidence source kind không hợp lệ: %s", value)
	}
	return result, nil
}

func ParseQueryObservationType(value string) (QueryObservationType, error) {
	result := QueryObservationType(strings.TrimSpace(value))
	if !IsValidQueryObservationType(result) {
		return "", fmt.Errorf("query observation type không hợp lệ: %s", value)
	}
	return result, nil
}

func ParseQueryDeviceCategory(value string) (QueryDeviceCategory, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return QueryDeviceUnknown, nil
	}
	result := QueryDeviceCategory(value)
	if !IsValidQueryDeviceCategory(result) {
		return "", fmt.Errorf("query device category không hợp lệ: %s", value)
	}
	return result, nil
}

func ParseQueryEvidenceStatus(value string) (QueryEvidenceStatus, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return QueryEvidenceStatusActive, nil
	}
	result := QueryEvidenceStatus(value)
	if !IsValidQueryEvidenceStatus(result) {
		return "", fmt.Errorf("query evidence status không hợp lệ: %s", value)
	}
	return result, nil
}

// IsQueryEvidenceCurrent xác định evidence có được phép tham gia readiness.
func IsQueryEvidenceCurrent(evidence SeoQueryEvidence, now time.Time) bool {
	if evidence.Status != QueryEvidenceStatusActive {
		return false
	}
	if evidence.ValidUntil != nil && evidence.ValidUntil.Before(now) {
		return false
	}
	return true
}

// IsQueryEvidenceObserved phân biệt observation thực với suggestion/assumption.
func IsQueryEvidenceObserved(evidence SeoQueryEvidence) bool {
	switch evidence.ObservationType {
	case QueryObservationObserved, QueryObservationReported:
		return true
	default:
		return false
	}
}
