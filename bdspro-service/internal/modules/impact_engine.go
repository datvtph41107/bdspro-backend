package modules

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	property_repo "bdspro/internal/repo/property"
	"context"
	"fmt"
	"time"
)

type EvaluateRequest struct {
	PropertyID      uint64
	CurrentVersion  *domain.PropertyLineage
	ProposedChanges *dto.UpdatePropertyDTO
	ActorID         uint64
	ActorRole       enums.ActorRole
	PreviewMode     bool
	Confirmed       bool
}

type ImpactEvaluation struct {
	ChangeDetection     *ChangeDetectionResult `json:"changeDetection"`
	SnapshotComparisons []*ComparisonResult    `json:"snapshotComparisons"`
	ImpactLevel         enums.ImpactLevel      `json:"impactLevel"`
	AffectedSummary     *AffectedSummary       `json:"affectedSummary"`
	ProjectedStates     *ProjectedStateSummary `json:"projectedStates"`
	ActionPolicy        enums.ActionPolicy     `json:"actionPolicy"`
	RequiresConfirm     bool                   `json:"requiresConfirm"`
	BlockReason         string                 `json:"blockReason,omitempty"`
	Warnings            []string               `json:"warnings"`
	Suggestions         []string               `json:"suggestions"`
	EvaluatedAt         time.Time              `json:"evaluatedAt"`
	EvaluationTimeMs    int64                  `json:"evaluationTimeMs"`
}

type ChangeDetectionResult struct {
	ChangedFields    []FieldChange `json:"changedFields"`
	CriticalFields   []string      `json:"criticalFields"`
	SnapshotMismatch bool          `json:"snapshotMismatch"`
}

type AffectedSummary struct {
	Products int `json:"products"`
	Listings int `json:"listings"`
	Assets   int `json:"assets"`
	Deals    int `json:"deals"`
	CrmNotes int `json:"crmNotes"`
	Total    int `json:"total"`
}

type ProjectedStateSummary struct {
	Active     int `json:"active"`
	Restricted int `json:"restricted"`
	Frozen     int `json:"frozen"`
	Archived   int `json:"archived"`
}

type ImpactDetectionEngine interface {
	Evaluate(ctx context.Context, req *EvaluateRequest) (*ImpactEvaluation, error)
	Preview(ctx context.Context, req *EvaluateRequest) (*ImpactEvaluation, error)
}

type impactEngineImpl struct {
	snapshotRepo       property_repo.PropertyImpactRepository
	ruleEngine         RuleEngine
	snapshotComparator SnapshotComparator
	propertyRepo       property_repo.PropertyLineageRepository
}

func NewImpactDetectionEngine(
	snapshotRepo property_repo.PropertyImpactRepository,
	ruleEngine RuleEngine,
	snapshotComparator SnapshotComparator,
	propertyRepo property_repo.PropertyLineageRepository,
) ImpactDetectionEngine {
	return &impactEngineImpl{
		snapshotRepo:       snapshotRepo,
		ruleEngine:         ruleEngine,
		snapshotComparator: snapshotComparator,
		propertyRepo:       propertyRepo,
	}
}

func (e *impactEngineImpl) Evaluate(ctx context.Context, req *EvaluateRequest) (*ImpactEvaluation, error) {
	startTime := time.Now()

	evaluation := &ImpactEvaluation{EvaluatedAt: startTime}

	changes, err := e.detectChanges(req.CurrentVersion, req.ProposedChanges)
	if err != nil {
		return nil, err
	}
	evaluation.ChangeDetection = changes

	if len(changes.ChangedFields) == 0 {
		evaluation.ImpactLevel = enums.ImpactLevelNone
		evaluation.ActionPolicy = enums.ActionPolicyAllowDirect
		evaluation.EvaluationTimeMs = time.Since(startTime).Milliseconds()
		return evaluation, nil
	}

	affectedEntities, err := e.snapshotRepo.GetAffectedEntities(ctx, req.PropertyID)
	if err != nil {
		return nil, err
	}

	evaluation.AffectedSummary = &AffectedSummary{
		Products: len(affectedEntities.Products),
		Listings: len(affectedEntities.Listings),
		Assets:   len(affectedEntities.Assets),
		Deals:    len(affectedEntities.Deals),
		CrmNotes: len(affectedEntities.CrmNotes),
		Total: len(affectedEntities.Products) + len(affectedEntities.Listings) +
			len(affectedEntities.Assets) + len(affectedEntities.Deals) +
			len(affectedEntities.CrmNotes),
	}

	var comparisons []*ComparisonResult
	for _, productID := range affectedEntities.Products {
		snapshot, err := e.snapshotRepo.GetCurrent(ctx, "product", productID)
		if err == nil && snapshot != nil {
			comp, err := e.snapshotComparator.Compare(ctx, snapshot, req.CurrentVersion)
			if err == nil {
				comparisons = append(comparisons, comp)
			}
		}
	}
	evaluation.SnapshotComparisons = comparisons

	ruleResult, err := e.ruleEngine.EvaluateImpact(changes.ChangedFields, req.ActorRole)
	if err != nil {
		return nil, err
	}
	evaluation.ImpactLevel = ruleResult.ImpactLevel
	evaluation.Warnings = ruleResult.Warnings

	evaluation.ProjectedStates = e.projectStates(comparisons)
	evaluation.ActionPolicy = e.determineActionPolicy(evaluation.ImpactLevel, req.ActorRole, req.PreviewMode, req.Confirmed)
	evaluation.RequiresConfirm = evaluation.ActionPolicy == enums.ActionPolicyRequireConfirm

	if evaluation.ActionPolicy == enums.ActionPolicyBlock {
		evaluation.BlockReason = e.getBlockReason(evaluation.ImpactLevel, ruleResult.CriticalChanges)
	}

	evaluation.Suggestions = e.generateSuggestions(evaluation)
	evaluation.EvaluationTimeMs = time.Since(startTime).Milliseconds()

	return evaluation, nil
}

func (e *impactEngineImpl) Preview(ctx context.Context, req *EvaluateRequest) (*ImpactEvaluation, error) {
	req.PreviewMode = true
	req.Confirmed = false
	return e.Evaluate(ctx, req)
}

func (e *impactEngineImpl) detectChanges(current *domain.PropertyLineage, proposed *dto.UpdatePropertyDTO) (*ChangeDetectionResult, error) {
	changes := &ChangeDetectionResult{
		ChangedFields:  []FieldChange{},
		CriticalFields: []string{},
	}

	currentMap := e.buildCurrentMap(current)

	if proposed.Info != nil {
		e.detectInfoChanges(currentMap, proposed.Info, changes)
	}
	if proposed.Location != nil {
		e.detectLocationChanges(currentMap, proposed.Location, changes)
	}
	if proposed.LandInfo != nil {
		e.detectLandInfoChanges(currentMap, proposed.LandInfo, changes)
	}
	if proposed.LineageMeta != nil {
		e.detectLineageMetaChanges(currentMap, proposed.LineageMeta, changes)
	}

	for _, change := range changes.ChangedFields {
		if e.isCriticalField(change.FieldPath) {
			changes.CriticalFields = append(changes.CriticalFields, change.FieldPath)
			changes.SnapshotMismatch = true
		}
	}

	return changes, nil
}

func (e *impactEngineImpl) buildCurrentMap(lineage *domain.PropertyLineage) map[string]interface{} {
	m := make(map[string]interface{})

	if lineage.PropertyInfo != nil {
		m["info.title"] = lineage.PropertyInfo.Title
		m["info.note"] = lineage.PropertyInfo.Note
		m["info.scope"] = lineage.PropertyInfo.Scope
		m["info.property_type_id"] = lineage.PropertyInfo.PropertyTypeID
	}

	if lineage.Location != nil {
		m["location.latitude"] = lineage.Location.Latitude
		m["location.longitude"] = lineage.Location.Longitude
		m["location.province_id"] = lineage.Location.ProvinceID
		m["location.ward_id"] = lineage.Location.WardID
	}

	if lineage.LandInfo != nil {
		m["land_info.area_total"] = lineage.LandInfo.AreaTotal
		m["land_info.plot"] = lineage.LandInfo.Plot
	}

	m["lineage_status"] = lineage.LineageStatus

	return m
}

func (e *impactEngineImpl) detectInfoChanges(current map[string]interface{}, info *dto.UpdateInfoDTO, changes *ChangeDetectionResult) {
	if info.Title != nil {
		e.addChangeIfChanged(changes, "info.title", current["info.title"], *info.Title)
	}
	if info.PropertyTypeID != nil {
		e.addChangeIfChanged(changes, "info.property_type_id", current["info.property_type_id"], *info.PropertyTypeID)
	}
}

func (e *impactEngineImpl) detectLocationChanges(current map[string]interface{}, loc *dto.UpdateLocationDTO, changes *ChangeDetectionResult) {
	if loc.Latitude != nil {
		e.addChangeIfChanged(changes, "location.latitude", current["location.latitude"], *loc.Latitude)
	}
	if loc.Longitude != nil {
		e.addChangeIfChanged(changes, "location.longitude", current["location.longitude"], *loc.Longitude)
	}
	if loc.ProvinceID != nil {
		e.addChangeIfChanged(changes, "location.province_id", current["location.province_id"], *loc.ProvinceID)
	}
}

func (e *impactEngineImpl) detectLandInfoChanges(current map[string]interface{}, land *dto.UpdateLandInfoDTO, changes *ChangeDetectionResult) {
	if land.AreaTotal != nil {
		e.addChangeIfChanged(changes, "land_info.area_total", current["land_info.area_total"], *land.AreaTotal)
	}
	if land.Plot != nil {
		e.addChangeIfChanged(changes, "land_info.plot", current["land_info.plot"], *land.Plot)
	}
}

func (e *impactEngineImpl) detectLineageMetaChanges(current map[string]interface{}, meta *dto.UpdateLineageMetaDTO, changes *ChangeDetectionResult) {
	if meta.LineageStatus != nil {
		e.addChangeIfChanged(changes, "lineage_status", current["lineage_status"], *meta.LineageStatus)
	}
}

func (e *impactEngineImpl) addChangeIfChanged(changes *ChangeDetectionResult, fieldPath string, oldVal, newVal interface{}) {
	if !valuesEqual(oldVal, newVal) {
		changes.ChangedFields = append(changes.ChangedFields, FieldChange{
			FieldPath:  fieldPath,
			OldValue:   oldVal,
			NewValue:   newVal,
			ChangeType: "update",
		})
	}
}

func (e *impactEngineImpl) isCriticalField(fieldPath string) bool {
	criticalFields := []string{
		"property_identify_id",
		"info.property_type_id",
		"location.latitude",
		"location.longitude",
		"land_info.plot",
		"lineage_status",
	}
	for _, cf := range criticalFields {
		if fieldPath == cf {
			return true
		}
	}
	return false
}

func (e *impactEngineImpl) projectStates(comparisons []*ComparisonResult) *ProjectedStateSummary {
	summary := &ProjectedStateSummary{}
	for _, comp := range comparisons {
		switch comp.ProjectedState {
		case enums.EffectiveStateActive:
			summary.Active++
		case enums.EffectiveStateRestricted:
			summary.Restricted++
		case enums.EffectiveStateFrozen:
			summary.Frozen++
		case enums.EffectiveStateArchived:
			summary.Archived++
		}
	}
	return summary
}

func (e *impactEngineImpl) determineActionPolicy(impactLevel enums.ImpactLevel, actorRole enums.ActorRole, previewMode, confirmed bool) enums.ActionPolicy {
	if previewMode {
		return enums.ActionPolicyAllowDirect
	}

	if impactLevel.RequiresConfirmation() && !confirmed {
		return enums.ActionPolicyRequireConfirm
	}

	if impactLevel.RequiresBlock() && !actorRole.HasFullAccess() {
		return enums.ActionPolicyBlock
	}

	if impactLevel == enums.ImpactLevelCritical && !confirmed {
		return enums.ActionPolicyRequireReview
	}

	return enums.ActionPolicyAllowDirect
}

func (e *impactEngineImpl) getBlockReason(impactLevel enums.ImpactLevel, criticalChanges []FieldChange) string {
	reason := fmt.Sprintf("Không thể lưu thay đổi do mức độ ảnh hưởng %s.\n", impactLevel.String())
	if len(criticalChanges) > 0 {
		reason += "Các thay đổi quan trọng:\n"
		for _, change := range criticalChanges {
			reason += fmt.Sprintf("- %s: từ %v thành %v\n", change.FieldPath, change.OldValue, change.NewValue)
		}
	}
	reason += "Vui lòng liên hệ admin để xử lý."
	return reason
}

func (e *impactEngineImpl) generateSuggestions(evaluation *ImpactEvaluation) []string {
	suggestions := []string{}

	if evaluation.ImpactLevel >= enums.ImpactLevelModerate {
		suggestions = append(suggestions, "Vui lòng rà soát lại các sản phẩm/tin đăng liên quan sau khi lưu thay đổi")
	}

	if evaluation.ProjectedStates.Restricted > 0 {
		suggestions = append(suggestions, fmt.Sprintf("Có %d sản phẩm sẽ chuyển sang trạng thái hạn chế, cần xử lý bổ sung", evaluation.ProjectedStates.Restricted))
	}

	if evaluation.ProjectedStates.Frozen > 0 {
		suggestions = append(suggestions, fmt.Sprintf("Có %d sản phẩm sẽ bị đóng băng, không thể tiếp tục giao dịch", evaluation.ProjectedStates.Frozen))
	}

	if len(evaluation.Warnings) > 0 {
		suggestions = append(suggestions, evaluation.Warnings...)
	}

	return suggestions
}
