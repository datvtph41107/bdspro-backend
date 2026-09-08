package application

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"tqd/internal/config"
	"tqd/internal/domain/discovery/model"
	"tqd/internal/usecase/discovery/ports"
)

func newRequestID() string {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "discovery-request"
	}
	return hex.EncodeToString(raw[:])
}

type Service struct{ repo ports.Repository }

func NewService(repo ports.Repository) *Service { return &Service{repo: repo} }

func validatePoint(p domain.Point) error {
	if p.Latitude < -90 || p.Latitude > 90 || p.Longitude < -180 || p.Longitude > 180 {
		return fmt.Errorf("invalid latitude/longitude")
	}
	return nil
}

func (s *Service) Identify(ctx context.Context, req domain.IdentifyRequest) (*domain.IdentifyResponse, error) {
	if err := validatePoint(req.Point); err != nil {
		return nil, err
	}
	if req.Options.Limit <= 0 || req.Options.Limit > 50 {
		req.Options.Limit = 20
	}
	if req.Options.ToleranceMeters <= 0 {
		req.Options.ToleranceMeters = 25
	}

	type result struct {
		items []domain.EntityCandidate
		err   error
	}
	providers := []func(context.Context, domain.IdentifyRequest) ([]domain.EntityCandidate, error){
		s.repo.IdentifyParcels, s.repo.IdentifyRegions, s.repo.IdentifyAdministrativeUnits, s.repo.IdentifyPOIs,
	}
	ch := make(chan result, len(providers))
	var wg sync.WaitGroup
	for _, provider := range providers {
		wg.Add(1)
		go func(p func(context.Context, domain.IdentifyRequest) ([]domain.EntityCandidate, error)) {
			defer wg.Done()
			items, err := p(ctx, req)
			ch <- result{items, err}
		}(provider)
	}
	wg.Wait()
	close(ch)

	candidates := make([]domain.EntityCandidate, 0, 16)
	warnings := []string{}
	partial := false
	for r := range ch {
		if r.err != nil {
			partial = true
			warnings = append(warnings, r.err.Error())
			continue
		}
		candidates = append(candidates, r.items...)
	}
	candidates = dedupe(candidates)
	rankIdentify(candidates, req)
	sortCandidates(candidates)
	if len(candidates) > req.Options.Limit {
		candidates = candidates[:req.Options.Limit]
	}

	response := &domain.IdentifyResponse{RequestID: newRequestID(), Candidates: candidates, Partial: partial, Warnings: warnings}
	if len(candidates) == 0 {
		if partial && len(warnings) == len(providers) {
			return nil, fmt.Errorf("all discovery providers failed")
		}
		response.Resolution = domain.Resolution{Status: "empty", Confidence: 0, Strategy: "spatial_candidates"}
		return response, nil
	}
	response.Primary = &response.Candidates[0]
	confidence := math.Min(1, math.Max(0.1, response.Candidates[0].Rank.Score/150))
	ambiguous := len(candidates) > 1 && math.Abs(candidates[0].Rank.Score-candidates[1].Rank.Score) < 12
	status := "resolved"
	if ambiguous {
		status = "ambiguous"
	}
	if partial {
		status = "partial"
	}
	response.Resolution = domain.Resolution{Status: status, PrimaryKey: response.Primary.Entity.Key, Confidence: confidence, Ambiguous: ambiguous, Strategy: "contextual_spatial_ranking"}
	return response, nil
}

var coordinatePattern = regexp.MustCompile(`^\s*(-?\d{1,2}(?:\.\d+)?)\s*[,; ]\s*(-?\d{1,3}(?:\.\d+)?)\s*$`)
var parcelLandMapPattern = regexp.MustCompile(`(?i)(?:thửa|thua)\s*[:#-]?\s*([[:alnum:]_-]+)\D+(?:tờ|to)\s*[:#-]?\s*([[:alnum:]_-]+)`)
var parcelMapLandPattern = regexp.MustCompile(`(?i)(?:tờ|to)\s*[:#-]?\s*([[:alnum:]_-]+)\D+(?:thửa|thua)\s*[:#-]?\s*([[:alnum:]_-]+)`)
var parcelSlashPattern = regexp.MustCompile(`^\s*([[:alnum:]_-]+)\s*[/\-]\s*([[:alnum:]_-]+)\s*$`)

func normalizeQuery(q string) string {
	return strings.Join(strings.Fields(strings.ToLower(strings.TrimSpace(q))), " ")
}

func interpret(q string) domain.SearchInterpretation {
	normalized := normalizeQuery(q)
	out := domain.SearchInterpretation{OriginalQuery: q, NormalizedQuery: normalized, DetectedIntents: []string{"free_text"}, StructuredTerms: map[string]string{}}
	if m := coordinatePattern.FindStringSubmatch(normalized); len(m) == 3 {
		out.DetectedIntents = []string{"coordinate_lookup"}
		out.StructuredTerms["latitude"] = m[1]
		out.StructuredTerms["longitude"] = m[2]
	} else if m := parcelLandMapPattern.FindStringSubmatch(normalized); len(m) == 3 {
		out.DetectedIntents = []string{"parcel_lookup"}
		out.StructuredTerms["landNumber"] = m[1]
		out.StructuredTerms["mapNumber"] = m[2]
	} else if m := parcelMapLandPattern.FindStringSubmatch(normalized); len(m) == 3 {
		out.DetectedIntents = []string{"parcel_lookup"}
		out.StructuredTerms["mapNumber"] = m[1]
		out.StructuredTerms["landNumber"] = m[2]
	} else if m := parcelSlashPattern.FindStringSubmatch(normalized); len(m) == 3 {
		// Compact input follows the placeholder contract: map-number/land-number.
		out.DetectedIntents = []string{"parcel_lookup"}
		out.StructuredTerms["mapNumber"] = m[1]
		out.StructuredTerms["landNumber"] = m[2]
	} else if strings.Contains(normalized, "phường") || strings.Contains(normalized, "xã") || strings.Contains(normalized, "tỉnh") || strings.Contains(normalized, "quận") || strings.Contains(normalized, "huyện") {
		out.DetectedIntents = []string{"administrative_lookup", "free_text"}
	} else if strings.Contains(normalized, "quy hoạch") || strings.Contains(normalized, "quy hoach") {
		out.DetectedIntents = []string{"planning_lookup", "free_text"}
	}
	return out
}

func (s *Service) Search(ctx context.Context, req domain.SearchRequest) (*domain.SearchResponse, error) {
	req.Query = strings.TrimSpace(req.Query)
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 20
	}
	interpretation := interpret(req.Query)
	if interpretation.DetectedIntents[0] == "coordinate_lookup" {
		lat, _ := strconv.ParseFloat(interpretation.StructuredTerms["latitude"], 64)
		lon, _ := strconv.ParseFloat(interpretation.StructuredTerms["longitude"], 64)
		identified, err := s.Identify(ctx, domain.IdentifyRequest{Point: domain.Point{Latitude: lat, Longitude: lon}, Context: req.MapContext, Options: domain.IdentifyOptions{Limit: req.Limit, ToleranceMeters: 25}})
		if err != nil {
			return nil, err
		}
		return groupSearch(newRequestID(), interpretation, filterKinds(identified.Candidates, req.EntityKinds), identified.Partial, identified.Warnings), nil
	}

	providerRequest := req
	if len(interpretation.DetectedIntents) > 0 && interpretation.DetectedIntents[0] == "parcel_lookup" {
		// Repository exact lookup uses "map_number land_number" while users normally type "thửa <land> tờ <map>".
		providerRequest.Query = strings.TrimSpace(interpretation.StructuredTerms["mapNumber"] + " " + interpretation.StructuredTerms["landNumber"])
	}

	type result struct {
		items []domain.EntityCandidate
		err   error
	}
	providers := []func(context.Context, domain.SearchRequest) ([]domain.EntityCandidate, error){s.repo.SearchParcels, s.repo.SearchRegions, s.repo.SearchPlanningProjects, s.repo.SearchAdministrativeUnits, s.repo.SearchPOIs}
	ch := make(chan result, len(providers))
	var wg sync.WaitGroup
	for _, provider := range providers {
		wg.Add(1)
		go func(p func(context.Context, domain.SearchRequest) ([]domain.EntityCandidate, error)) {
			defer wg.Done()
			items, err := p(ctx, providerRequest)
			ch <- result{items, err}
		}(provider)
	}
	wg.Wait()
	close(ch)
	candidates := []domain.EntityCandidate{}
	warnings := []string{}
	partial := false
	for r := range ch {
		if r.err != nil {
			partial = true
			warnings = append(warnings, r.err.Error())
			continue
		}
		candidates = append(candidates, r.items...)
	}
	if len(candidates) == 0 && len(warnings) == len(providers) {
		return nil, fmt.Errorf("all search providers failed")
	}
	candidates = filterKinds(dedupe(candidates), req.EntityKinds)
	rankSearch(candidates, req, interpretation)
	sortCandidates(candidates)
	if len(candidates) > req.Limit {
		candidates = candidates[:req.Limit]
	}
	return groupSearch(newRequestID(), interpretation, candidates, partial, warnings), nil
}

func (s *Service) GetEntity(ctx context.Context, key string) (*domain.EntityCandidate, error) {
	ref, err := domain.ParseEntityKey(key)
	if err != nil {
		return nil, err
	}
	return s.repo.GetEntity(ctx, ref)
}

func dedupe(in []domain.EntityCandidate) []domain.EntityCandidate {
	seen := map[string]bool{}
	out := make([]domain.EntityCandidate, 0, len(in))
	for _, c := range in {
		if c.Entity.Key == "" || seen[c.Entity.Key] {
			continue
		}
		seen[c.Entity.Key] = true
		out = append(out, c)
	}
	return out
}
func sortCandidates(items []domain.EntityCandidate) {
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Rank.Score == items[j].Rank.Score {
			return items[i].Entity.Key < items[j].Entity.Key
		}
		return items[i].Rank.Score > items[j].Rank.Score
	})
}
func contains(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}
func allowedKind(list []domain.EntityKind, kind domain.EntityKind) bool {
	if len(list) == 0 {
		return true
	}
	for _, value := range list {
		if value == kind {
			return true
		}
	}
	return false
}

func filterKinds(items []domain.EntityCandidate, kinds []domain.EntityKind) []domain.EntityCandidate {
	if len(kinds) == 0 {
		return items
	}
	out := make([]domain.EntityCandidate, 0, len(items))
	for _, item := range items {
		if allowedKind(kinds, item.Entity.Kind) {
			out = append(out, item)
		}
	}
	return out
}

func preferred(list []domain.EntityKind, kind domain.EntityKind) bool {
	for _, v := range list {
		if v == kind {
			return true
		}
	}
	return false
}

const maxIdentifyProviderScore = 50

func normalizedIdentifyProviderScore(raw float64) (float64, bool) {
	if raw < 0 {
		return 0, true
	}
	if raw > maxIdentifyProviderScore {
		return maxIdentifyProviderScore, true
	}
	return raw, false
}

func rankIdentify(items []domain.EntityCandidate, req domain.IdentifyRequest) {
	highDetailZoom := req.Context.Zoom >= float64(config.ParcelPolygonMinZ)
	lowDetailZoom := req.Context.Zoom > 0 && !highDetailZoom

	for i := range items {
		score, normalized := normalizedIdentifyProviderScore(items[i].Rank.Score)
		reasons := append([]string{}, items[i].Rank.ReasonCodes...)
		if normalized {
			reasons = append(reasons, "provider_score_normalized")
		}

		switch items[i].Entity.Kind {
		case domain.EntityKindParcel:
			score += 65

			// At parcel-detail zoom, exact parcel containment is more specific than
			// a planning region covering the same point. Planning mode remains an
			// explicit override because users may intentionally inspect a layer.
			if highDetailZoom && req.Context.MapMode != "planning" {
				score += 120
				reasons = append(reasons, "high_detail_zoom")
			}
			if req.Context.MapMode == "parcel" {
				score += 80
				reasons = append(reasons, "map_mode_parcel")
			}

		case domain.EntityKindPlanningRegion:
			score += 55
			if lowDetailZoom {
				score += 35
				reasons = append(reasons, "overview_zoom")
			}
			if req.Context.MapMode == "planning" {
				score += 120
				reasons = append(reasons, "map_mode_planning")
			}
			if layer, ok := items[i].Attributes["layerId"].(string); ok && contains(req.Context.ActiveLayerIDs, layer) {
				score += 45
				reasons = append(reasons, "active_layer")
			}

		case domain.EntityKindAdministrativeUnit:
			score += 25

		case domain.EntityKindPOI:
			score += 20
			if req.Context.MapMode == "poi" {
				score += 50
				reasons = append(reasons, "map_mode_poi")
			}
		}

		if preferred(req.Context.PreferredKinds, items[i].Entity.Kind) {
			score += 30
			reasons = append(reasons, "preferred_kind")
		}
		items[i].Rank = domain.Rank{Score: score, ReasonCodes: reasons}
	}
}

func rankSearch(items []domain.EntityCandidate, req domain.SearchRequest, interpretation domain.SearchInterpretation) {
	q := interpretation.NormalizedQuery
	for i := range items {
		score := items[i].Rank.Score
		reasons := append([]string{}, items[i].Rank.ReasonCodes...)
		title := normalizeQuery(items[i].Presentation.Title)
		if title == q {
			score += 100
			reasons = append(reasons, "exact_title")
		} else if strings.HasPrefix(title, q) {
			score += 60
			reasons = append(reasons, "prefix_title")
		} else if strings.Contains(title, q) {
			score += 35
			reasons = append(reasons, "contains_title")
		}
		if len(interpretation.DetectedIntents) > 0 && interpretation.DetectedIntents[0] == "parcel_lookup" && items[i].Entity.Kind == domain.EntityKindParcel {
			score += 80
			reasons = append(reasons, "parcel_intent")
		}
		if len(interpretation.DetectedIntents) > 0 && interpretation.DetectedIntents[0] == "administrative_lookup" && items[i].Entity.Kind == domain.EntityKindAdministrativeUnit {
			score += 65
			reasons = append(reasons, "administrative_intent")
		}
		if len(interpretation.DetectedIntents) > 0 && interpretation.DetectedIntents[0] == "planning_lookup" && (items[i].Entity.Kind == domain.EntityKindPlanningRegion || items[i].Entity.Kind == domain.EntityKindPlanningProject) {
			score += 65
			reasons = append(reasons, "planning_intent")
		}
		items[i].Rank = domain.Rank{Score: score, ReasonCodes: reasons}
	}
}

func groupSearch(requestID string, interpretation domain.SearchInterpretation, candidates []domain.EntityCandidate, partial bool, warnings []string) *domain.SearchResponse {
	labels := map[domain.EntityKind]string{
		domain.EntityKindAdministrativeUnit: "Địa bàn hành chính",
		domain.EntityKindParcel:             "Thửa đất",
		domain.EntityKindPlanningRegion:     "Vùng quy hoạch",
		domain.EntityKindPlanningProject:    "Đồ án quy hoạch",
		domain.EntityKindPlanningMap:        "Bản đồ quy hoạch",
		domain.EntityKindMapLayer:           "Lớp bản đồ",
		domain.EntityKindPOI:                "Điểm quan tâm",
		domain.EntityKindLegalDocument:      "Văn bản pháp lý",
		domain.EntityKindNewsArticle:        "Tin quy hoạch",
		domain.EntityKindAnalysisReport:     "Báo cáo và phân tích",
	}

	order := []domain.EntityKind{
		domain.EntityKindAdministrativeUnit,
		domain.EntityKindParcel,
		domain.EntityKindPlanningRegion,
		domain.EntityKindPlanningProject,
		domain.EntityKindPlanningMap,
		domain.EntityKindMapLayer,
		domain.EntityKindPOI,
		domain.EntityKindLegalDocument,
		domain.EntityKindNewsArticle,
		domain.EntityKindAnalysisReport,
	}

	byKind := map[domain.EntityKind][]domain.EntityCandidate{}
	for _, candidate := range candidates {
		byKind[candidate.Entity.Kind] = append(byKind[candidate.Entity.Kind], candidate)
	}

	groups := make([]domain.SearchGroup, 0, len(byKind))
	for _, kind := range order {
		items := byKind[kind]
		if len(items) == 0 {
			continue
		}

		label := labels[kind]
		if label == "" {
			label = string(kind)
		}

		groups = append(groups, domain.SearchGroup{
			Kind:  kind,
			Label: label,
			Total: len(items),
			Items: items,
		})
	}

	return &domain.SearchResponse{
		RequestID:      requestID,
		Interpretation: interpretation,
		Groups:         groups,
		Partial:        partial,
		Warnings:       warnings,
	}
}
