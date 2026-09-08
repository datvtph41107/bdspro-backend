package modules

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
	"context"
	"fmt"
)

type ComparisonResult struct {
	EntityType       string               `json:"entityType"`
	EntityID         uint64               `json:"entityId"`
	MismatchedFields []FieldMismatch      `json:"mismatchedFields"`
	MismatchLevel    enums.ImpactLevel    `json:"mismatchLevel"`
	ProjectedState   enums.EffectiveState `json:"projectedState"`
}

type FieldMismatch struct {
	FieldPath   string            `json:"fieldPath"`
	SnapshotVal interface{}       `json:"snapshotVal"`
	CurrentVal  interface{}       `json:"currentVal"`
	Severity    enums.ImpactLevel `json:"severity"`
	Description string            `json:"description"`
}

type SnapshotComparator interface {
	Compare(ctx context.Context, snapshot *domain.PropertyImpactVersion, property *domain.PropertyLineage) (*ComparisonResult, error)
}

type snapshotComparatorImpl struct {
	ruleEngine RuleEngine
}

func NewSnapshotComparator(ruleEngine RuleEngine) SnapshotComparator {
	return &snapshotComparatorImpl{ruleEngine: ruleEngine}
}

func (c *snapshotComparatorImpl) Compare(ctx context.Context, snapshot *domain.PropertyImpactVersion, property *domain.PropertyLineage) (*ComparisonResult, error) {
	var snapshotData map[string]interface{}
	if err := snapshot.SnapshotData.Unmarshal(&snapshotData); err != nil {
		return nil, err
	}

	propertyFields := c.extractRelevantFields(property, snapshot.EntityType)
	mismatches := c.findMismatches(snapshotData, propertyFields, snapshot.EntityType)
	mismatchLevel := c.calculateMismatchLevel(mismatches)
	projectedState := c.projectState(mismatchLevel, snapshot.EntityType)

	return &ComparisonResult{
		EntityType:       snapshot.EntityType,
		EntityID:         snapshot.EntityID,
		MismatchedFields: mismatches,
		MismatchLevel:    mismatchLevel,
		ProjectedState:   projectedState,
	}, nil
}

func (c *snapshotComparatorImpl) extractRelevantFields(property *domain.PropertyLineage, entityType string) map[string]interface{} {
	fields := make(map[string]interface{})

	if property.PropertyInfo != nil {
		fields["info.title"] = property.PropertyInfo.Title
		fields["info.property_type_id"] = property.PropertyInfo.PropertyTypeID
		fields["info.legal_status"] = property.PropertyInfo.LegalStatus
	}

	if property.Location != nil {
		fields["location.latitude"] = property.Location.Latitude
		fields["location.longitude"] = property.Location.Longitude
		fields["location.province_id"] = property.Location.ProvinceID
		fields["location.ward_id"] = property.Location.WardID
	}

	if property.LandInfo != nil {
		fields["land_info.area_total"] = property.LandInfo.AreaTotal
		fields["land_info.plot"] = property.LandInfo.Plot
	}

	if property.BuildingInfo != nil {
		fields["building_info.area_actual"] = property.BuildingInfo.AreaActual
		fields["building_info.floors"] = property.BuildingInfo.Floors
	}

	return fields
}

func (c *snapshotComparatorImpl) findMismatches(snapshotData, propertyFields map[string]interface{}, entityType string) []FieldMismatch {
	var mismatches []FieldMismatch

	for fieldPath, propertyVal := range propertyFields {
		snapshotVal, exists := snapshotData[fieldPath]
		if !exists {
			continue
		}

		if !valuesEqual(snapshotVal, propertyVal) {
			mismatches = append(mismatches, FieldMismatch{
				FieldPath:   fieldPath,
				SnapshotVal: snapshotVal,
				CurrentVal:  propertyVal,
				Severity:    c.getFieldSeverity(fieldPath, entityType),
				Description: fmt.Sprintf("%s thay đổi từ %v thành %v", fieldPath, snapshotVal, propertyVal),
			})
		}
	}

	return mismatches
}

func (c *snapshotComparatorImpl) calculateMismatchLevel(mismatches []FieldMismatch) enums.ImpactLevel {
	if len(mismatches) == 0 {
		return enums.ImpactLevelNone
	}
	maxLevel := enums.ImpactLevelNone
	for _, mismatch := range mismatches {
		if mismatch.Severity > maxLevel {
			maxLevel = mismatch.Severity
		}
	}
	return maxLevel
}

func (c *snapshotComparatorImpl) projectState(mismatchLevel enums.ImpactLevel, entityType string) enums.EffectiveState {
	switch mismatchLevel {
	case enums.ImpactLevelNone, enums.ImpactLevelMinor:
		return enums.EffectiveStateActive
	case enums.ImpactLevelModerate:
		return enums.EffectiveStateRestricted
	case enums.ImpactLevelMajor:
		return enums.EffectiveStateFrozen
	case enums.ImpactLevelCritical:
		return enums.EffectiveStateArchived
	default:
		return enums.EffectiveStateActive
	}
}

func (c *snapshotComparatorImpl) getFieldSeverity(fieldPath, entityType string) enums.ImpactLevel {
	rules := c.ruleEngine.GetRulesForField(fieldPath)
	if len(rules) > 0 {
		return rules[0].CriticalLevel
	}
	return enums.ImpactLevelMinor
}

func valuesEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}
