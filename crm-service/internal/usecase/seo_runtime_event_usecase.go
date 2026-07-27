package usecase

import (
	_errors "common/errors"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"
	"crm/internal/repo"

	"gorm.io/datatypes"
)

var seoRuntimeEventNames = map[string]struct{}{
	"seo_page_view":             {},
	"seo_map_cta_click":         {},
	"seo_related_link_click":    {},
	"seo_faq_open":              {},
	"seo_search":                {},
	"seo_map_identify":          {},
	"rn_activation":             {},
	"entity_detail_open":        {},
	"canonical_page_navigation": {},
	"navigation_error":          {},
	"invalid_canonical_url":     {},
	"related_content_open":      {},
	"map_entity_focus":          {},
	"web_vitals":                {},
}

type SeoRuntimeEventUsecase struct {
	repo repo.SeoRuntimeEventRepo
	now  func() time.Time
}

func NewSeoRuntimeEventUsecase(repository repo.SeoRuntimeEventRepo) *SeoRuntimeEventUsecase {
	return &SeoRuntimeEventUsecase{repo: repository, now: time.Now}
}

func (u *SeoRuntimeEventUsecase) Ingest(
	ctx context.Context,
	req *dto.SeoRuntimeEventIngestRequest,
) (*dto.SeoRuntimeEventIngestResponse, error) {
	if u == nil || u.repo == nil {
		return nil, _errors.ReturnError(503, "SEO runtime measurement repository chưa sẵn sàng")
	}
	if req == nil {
		return nil, _errors.ReturnError(400, "SEO runtime event request không hợp lệ")
	}

	normalized, err := normalizeSeoRuntimeEvent(req, u.now().UTC())
	if err != nil {
		return nil, err
	}
	created, err := u.repo.CreateIdempotent(ctx, normalized)
	if err != nil {
		return nil, fmt.Errorf("persist SEO runtime event: %w", err)
	}

	return &dto.SeoRuntimeEventIngestResponse{
		SchemaVersion: dto.SeoRuntimeEventSchemaVersion,
		Accepted:      true,
		Persisted:     created,
		Idempotent:    !created,
		EventKey:      normalized.EventKey,
		ReceivedAt:    normalized.ReceivedAt.Format(time.RFC3339Nano),
	}, nil
}

func normalizeSeoRuntimeEvent(req *dto.SeoRuntimeEventIngestRequest, receivedAt time.Time) (*seo_domain.SeoRuntimeEvent, error) {
	if req.SchemaVersion != 1 {
		return nil, _errors.ReturnError(400, "SEO runtime event schemaVersion không được hỗ trợ")
	}
	eventName := strings.ToLower(strings.TrimSpace(req.Event))
	if _, ok := seoRuntimeEventNames[eventName]; !ok {
		return nil, _errors.ReturnError(400, "SEO runtime event không được hỗ trợ")
	}

	source := boundedRuntimeString(req.Source, 32)
	if source != "nextjs" && source != "react-native" {
		return nil, _errors.ReturnError(400, "SEO runtime event source không hợp lệ")
	}
	sessionID := boundedRuntimeString(req.SessionID, 120)
	moduleID := boundedRuntimeString(req.ModuleID, 100)
	canonicalPath := strings.TrimSpace(req.CanonicalPath)
	canonicalURL := strings.TrimSpace(req.CanonicalURL)
	if sessionID == "" || moduleID == "" || !isSafeRuntimeCanonicalPath(canonicalPath) || !isSafeRuntimeCanonicalURL(canonicalURL, canonicalPath) {
		return nil, _errors.ReturnError(400, "SEO runtime event identity không hợp lệ")
	}

	occurredAt, err := time.Parse(time.RFC3339Nano, strings.TrimSpace(req.Timestamp))
	if err != nil {
		return nil, _errors.ReturnError(400, "SEO runtime event timestamp không hợp lệ")
	}
	occurredAt = occurredAt.UTC()
	if occurredAt.After(receivedAt.Add(10*time.Minute)) || occurredAt.Before(receivedAt.AddDate(0, -3, 0)) {
		return nil, _errors.ReturnError(400, "SEO runtime event timestamp nằm ngoài cửa sổ cho phép")
	}

	metrics := normalizeSeoRuntimeMetrics(req.Metrics)
	metadata := normalizeSeoRuntimeMetadata(req.Metadata)
	metricsJSON, _ := json.Marshal(metrics)
	metadataJSON, _ := json.Marshal(metadata)

	action := boundedRuntimeString(req.Action, 100)
	target := boundedRuntimeString(req.Target, 1000)
	environment := boundedRuntimeString(req.Environment, 40)
	userAgent := boundedRuntimeString(req.UserAgent, 300)
	eventKey := seoRuntimeEventKey(eventName, source, sessionID, canonicalURL, occurredAt, action, target)

	return &seo_domain.SeoRuntimeEvent{
		EventKey:      eventKey,
		SchemaVersion: 1,
		Source:        source,
		EventName:     eventName,
		SessionID:     sessionID,
		CanonicalURL:  canonicalURL,
		CanonicalPath: canonicalPath,
		ModuleID:      moduleID,
		Environment:   environment,
		Indexable:     req.Indexable,
		Action:        action,
		Target:        target,
		OccurredAt:    occurredAt,
		ReceivedAt:    receivedAt,
		Metrics:       datatypes.JSON(metricsJSON),
		Metadata:      datatypes.JSON(metadataJSON),
		UserAgent:     userAgent,
	}, nil
}

func isSafeRuntimeCanonicalPath(value string) bool {
	return strings.HasPrefix(value, "/") && !strings.HasPrefix(value, "//") && !strings.Contains(value, "..") && !strings.ContainsAny(value, "\\?#")
}

func isSafeRuntimeCanonicalURL(value string, canonicalPath string) bool {
	parsed, err := url.Parse(value)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || strings.TrimSpace(parsed.Host) == "" {
		return false
	}
	return parsed.EscapedPath() == canonicalPath || parsed.Path == canonicalPath
}

func boundedRuntimeString(value string, limit int) string {
	value = strings.TrimSpace(value)
	if len(value) > limit {
		value = value[:limit]
	}
	return value
}

func normalizeSeoRuntimeMetrics(input map[string]float64) map[string]float64 {
	output := make(map[string]float64)
	count := 0
	for key, value := range input {
		if count >= 20 || value != value || value > 1e12 || value < -1e12 {
			continue
		}
		key = boundedRuntimeString(key, 80)
		if key == "" {
			continue
		}
		output[key] = value
		count++
	}
	return output
}

func normalizeSeoRuntimeMetadata(input map[string]interface{}) map[string]interface{} {
	output := make(map[string]interface{})
	count := 0
	for rawKey, rawValue := range input {
		if count >= 20 {
			break
		}
		key := boundedRuntimeString(rawKey, 80)
		if key == "" {
			continue
		}
		switch value := rawValue.(type) {
		case nil:
			output[key] = nil
		case bool:
			output[key] = value
		case float64:
			if value == value && value <= 1e12 && value >= -1e12 {
				output[key] = value
			}
		case string:
			output[key] = boundedRuntimeString(value, 300)
		default:
			continue
		}
		count++
	}
	return output
}

func seoRuntimeEventKey(eventName, source, sessionID, canonicalURL string, occurredAt time.Time, action, target string) string {
	payload := strings.Join([]string{
		"v1",
		eventName,
		source,
		sessionID,
		canonicalURL,
		occurredAt.Format(time.RFC3339Nano),
		action,
		target,
	}, "\x1f")
	hash := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(hash[:])
}
