package domain

const AdministrativeUnitProjectionSchemaVersion = "tqd-administrative-unit-projection/v1"

type AdministrativeUnitParent struct {
	EntityID           string `json:"entityId"`
	CanonicalKey       string `json:"canonicalKey"`
	Slug               string `json:"slug"`
	CanonicalPath      string `json:"canonicalPath"`
	UnitType           string `json:"unitType"`
	Name               string `json:"name"`
	AdministrativeCode string `json:"administrativeCode"`
}

type AdministrativeUnitHierarchyItem struct {
	EntityID           string `json:"entityId"`
	UnitType           string `json:"unitType"`
	Name               string `json:"name"`
	AdministrativeCode string `json:"administrativeCode"`
	CanonicalPath      string `json:"canonicalPath"`
}

type AdministrativePlanningContext struct {
	Summary              string                          `json:"summary"`
	PlanningProjectCount *int64                          `json:"planningProjectCount,omitempty"`
	DataAvailability     string                          `json:"dataAvailability"`
	RelatedProjects      []AdministrativePlanningProject `json:"relatedProjects,omitempty"`
}

type AdministrativePlanningProject struct {
	EntityID   string `json:"entityId"`
	Slug       string `json:"slug"`
	Code       string `json:"code,omitempty"`
	Name       string `json:"name"`
	Summary    string `json:"summary,omitempty"`
	PublicPath string `json:"publicPath"`
	MapPath    string `json:"mapPath"`
}

type AdministrativeUnitProjection struct {
	SchemaVersion      string                            `json:"schemaVersion"`
	EntityID           string                            `json:"entityId"`
	CanonicalKey       string                            `json:"canonicalKey"`
	Slug               string                            `json:"slug"`
	CanonicalPath      string                            `json:"canonicalPath"`
	UnitType           string                            `json:"unitType"`
	UnitTypeName       string                            `json:"unitTypeName"`
	Name               string                            `json:"name"`
	DisplayName        string                            `json:"displayName"`
	AdministrativeCode string                            `json:"administrativeCode"`
	Parent             *AdministrativeUnitParent         `json:"parent,omitempty"`
	Hierarchy          []AdministrativeUnitHierarchyItem `json:"hierarchy"`
	Location           *Point                            `json:"location,omitempty"`
	GeometrySummary    PreviewContext                    `json:"geometrySummary"`
	PlanningContext    AdministrativePlanningContext     `json:"planningContext"`
	Source             ProjectionSource                  `json:"source"`
}
