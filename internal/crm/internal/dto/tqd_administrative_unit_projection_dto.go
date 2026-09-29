package dto

const TqdAdministrativeUnitProjectionSchemaVersion = "tqd-administrative-unit-projection/v1"

type TqdAdministrativeUnitParent struct {
	EntityID           string `json:"entityId"`
	CanonicalKey       string `json:"canonicalKey"`
	Slug               string `json:"slug"`
	CanonicalPath      string `json:"canonicalPath"`
	UnitType           string `json:"unitType"`
	Name               string `json:"name"`
	AdministrativeCode string `json:"administrativeCode"`
}

type TqdAdministrativeUnitHierarchyItem struct {
	EntityID           string `json:"entityId"`
	UnitType           string `json:"unitType"`
	Name               string `json:"name"`
	AdministrativeCode string `json:"administrativeCode"`
	CanonicalPath      string `json:"canonicalPath"`
}

type TqdAdministrativePlanningContext struct {
	Summary              string                             `json:"summary"`
	PlanningProjectCount *int64                             `json:"planningProjectCount,omitempty"`
	DataAvailability     string                             `json:"dataAvailability"`
	RelatedProjects      []TqdAdministrativePlanningProject `json:"relatedProjects,omitempty"`
}

type TqdAdministrativePlanningProject struct {
	EntityID   string `json:"entityId"`
	Slug       string `json:"slug"`
	Code       string `json:"code,omitempty"`
	Name       string `json:"name"`
	Summary    string `json:"summary,omitempty"`
	PublicPath string `json:"publicPath"`
	MapPath    string `json:"mapPath"`
}

type TqdAdministrativeUnitProjection struct {
	SchemaVersion      string                               `json:"schemaVersion"`
	EntityID           string                               `json:"entityId"`
	CanonicalKey       string                               `json:"canonicalKey"`
	Slug               string                               `json:"slug"`
	CanonicalPath      string                               `json:"canonicalPath"`
	UnitType           string                               `json:"unitType"`
	UnitTypeName       string                               `json:"unitTypeName"`
	Name               string                               `json:"name"`
	DisplayName        string                               `json:"displayName"`
	AdministrativeCode string                               `json:"administrativeCode"`
	Parent             *TqdAdministrativeUnitParent         `json:"parent,omitempty"`
	Hierarchy          []TqdAdministrativeUnitHierarchyItem `json:"hierarchy"`
	Location           *TqdPlanningPoint                    `json:"location,omitempty"`
	GeometrySummary    TqdPlanningPreview                   `json:"geometrySummary"`
	PlanningContext    TqdAdministrativePlanningContext     `json:"planningContext"`
	Source             TqdPlanningProjectionSource          `json:"source"`
}
