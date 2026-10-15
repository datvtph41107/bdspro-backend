package entity

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// PropertyLineage is the property aggregate root (property_lineage).
// ORM::TABLE
type PropertyLineage struct {
	ID                 uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt          *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at;index"`
	OriginProfileID    *uint64        `gorm:"column:origin_profile_id;index"`
	LineageStatus      uint32         `gorm:"column:lineage_status;default:10"`
	PropertyIdentifyID *uint64        `gorm:"column:property_identify_id"`
	LocationID         *uint64        `gorm:"column:location_id"`
	PropertyInfoID     *uint64        `gorm:"column:property_info_id"`
	LandInfoID         *uint64        `gorm:"column:land_info_id"`
	BuildingInfoID     *uint64        `gorm:"column:building_info_id"`
	EdvidenceID        *uint64        `gorm:"column:edvidence_id"`
	StatisticID        *uint64        `gorm:"column:statistic_id"`
	NationalID         string         `gorm:"column:national_id;size:50"`

	Info         *PropertyInfo         `gorm:"foreignKey:PropertyInfoID;references:ID"`
	Location     *PropertyLocation     `gorm:"foreignKey:LocationID;references:ID"`
	LandInfo     *PropertyLandInfo     `gorm:"foreignKey:LandInfoID;references:ID"`
	BuildingInfo *PropertyBuildingInfo `gorm:"foreignKey:BuildingInfoID;references:ID"`
	Statistic    *PropertyStatistic    `gorm:"foreignKey:StatisticID;references:ID"`
	Identify     *PropertyIdentify     `gorm:"foreignKey:PropertyIdentifyID;references:ID"`
	MediaList    []PropertyMedia       `gorm:"foreignKey:LineageID;references:ID"`
}

func (PropertyLineage) TableName() string { return "property_lineage" }

// PropertyInfo holds title / type for a lineage.
type PropertyInfo struct {
	ID                 uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt          *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at;index"`
	PropertyIdentifyID *uint64        `gorm:"column:property_identify_id"`
	Title              string         `gorm:"column:title;size:255"`
	Identifier         string         `gorm:"column:identifier;size:100"`
	UnitCode           string         `gorm:"column:unit_code;size:50"`
	Level              string         `gorm:"column:level;size:100"`
	Note               string         `gorm:"column:note;size:255"`
	SourceType         uint32         `gorm:"column:source_type;not null;default:20"`
	LegalStatus        uint32         `gorm:"column:legal_status"`
	Scope              uint32         `gorm:"column:scope;default:10"`
	PropertyTypeID     *uint64        `gorm:"column:property_type_id"`
	ProjectID          *uint64        `gorm:"column:project_id"`
	AvatarID           *uint64        `gorm:"column:avatar_id"`

	PropertyType *PropertyType  `gorm:"foreignKey:PropertyTypeID;references:ID"`
	Avatar       *PropertyMedia `gorm:"foreignKey:AvatarID;references:ID"`
}

func (PropertyInfo) TableName() string { return "property_info" }

// PropertyLocation holds address fields for a lineage.
type PropertyLocation struct {
	ID                 uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt          *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at;index"`
	PropertyIdentifyID *uint64        `gorm:"column:property_identify_id"`
	AddressDetail      string         `gorm:"column:address_detail;size:500"`
	RegionID           *uint64        `gorm:"column:region_id"`
	ProvinceID         *uint64        `gorm:"column:province_id"`
	DistrictID         *uint64        `gorm:"column:district_id"`
	WardID             *uint64        `gorm:"column:ward_id"`
	Latitude           *float64       `gorm:"column:latitude"`
	Longitude          *float64       `gorm:"column:longitude"`
	MapURL             string         `gorm:"column:map_url;size:500"`
	PositionType       uint32         `gorm:"column:position_type;not null;default:10"`

	Province *ProvinceV2 `gorm:"foreignKey:ProvinceID;references:ID"`
	Ward     *WardV2     `gorm:"foreignKey:WardID;references:ID"`
}

func (PropertyLocation) TableName() string { return "property_location" }

// ProvinceV2 / WardV2 are thin name lookups on the v4 DSN.
type ProvinceV2 struct {
	ID   uint64 `gorm:"primaryKey"`
	Name string `gorm:"column:name"`
}

func (ProvinceV2) TableName() string { return "province_v2" }

type WardV2 struct {
	ID   uint64 `gorm:"primaryKey"`
	Name string `gorm:"column:name"`
}

func (WardV2) TableName() string { return "ward_v2" }

// PropertyLandInfo is land / GCN details.
type PropertyLandInfo struct {
	ID                 uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt          *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at;index"`
	PropertyIdentifyID *uint64        `gorm:"column:property_identify_id"`
	DocumentType       uint32         `gorm:"column:document_type;default:10"`
	DocumentNo         string         `gorm:"column:document_no;size:50"`
	IssuringAuth       string         `gorm:"column:issuring_auth;size:50"`
	Plot               uint32         `gorm:"column:plot"`
	Sheet              uint32         `gorm:"column:sheet"`
	AreaTotal          *float64       `gorm:"column:area_total;type:decimal(15,2)"`
	AreaLand           *float64       `gorm:"column:area_land;type:decimal(15,2)"`
	AreaPlant          *float64       `gorm:"column:area_plant;type:decimal(15,2)"`
	PurposeUsed        uint32         `gorm:"column:purpose_used;default:10"`
	Note               string         `gorm:"column:note"`
	FrontWidth         *float64       `gorm:"column:front_width;type:decimal(10,2)"`
	Depth              *float64       `gorm:"column:depth;type:decimal(10,2)"`
	StreetWidth        *float64       `gorm:"column:street_width;type:decimal(10,2)"`
	LandNote           string         `gorm:"column:land_note;type:text"`
}

func (PropertyLandInfo) TableName() string { return "property_land_info" }

// PropertyBuildingInfo is building / house details on a lineage.
type PropertyBuildingInfo struct {
	ID                 uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt          *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at;index"`
	PropertyIdentifyID *uint64        `gorm:"column:property_identify_id"`
	AreaActual         *float64       `gorm:"column:area_actual;type:decimal(15,2)"`
	AreaFloor          *float64       `gorm:"column:area_floor;type:decimal(15,2)"`
	AreaConstruction   *float64       `gorm:"column:area_construction;type:decimal(15,2)"`
	Floors             *uint32        `gorm:"column:floors"`
	RoomNumber         *uint32        `gorm:"column:room_number"`
	Bedrooms           *uint32        `gorm:"column:bedrooms"`
	Bathrooms          *uint32        `gorm:"column:bathrooms"`
	Note               string         `gorm:"column:note;type:text"`
	BuildStatus        uint32         `gorm:"column:build_status;not null;default:10"`
	BuildingType       uint32         `gorm:"column:building_type;not null;default:10"`
	Direction          uint32         `gorm:"column:direction"`
	BalconyDirection   uint32         `gorm:"column:balcony_direction"`
}

func (PropertyBuildingInfo) TableName() string { return "property_building_info" }

// PropertyStatistic holds aggregate counts for a lineage.
type PropertyStatistic struct {
	ID                 uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt          *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at;index"`
	PropertyIdentifyID *uint64        `gorm:"column:property_identify_id"`
	ProductCount       int64          `gorm:"column:product_count;default:0"`
	AssetCount         int64          `gorm:"column:asset_count;default:0"`
	ListingCount       int64          `gorm:"column:listing_count;default:0"`
}

func (PropertyStatistic) TableName() string { return "property_statistic" }

// PropertyMedia is image/video attached to a property lineage.
type PropertyMedia struct {
	ID                 uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt          *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at;index"`
	PropertyIdentifyID *uint64        `gorm:"column:property_identify_id"`
	LineageID          uint64         `gorm:"column:lineage_id;index"`
	MediaType          string         `gorm:"column:media_type;size:20"`
	MediaURL           string         `gorm:"column:media_url;size:500"`
	ThumbURL           string         `gorm:"column:thumb_url;size:500"`
	SortOrder          int32          `gorm:"column:sort_order;default:0"`
}

func (PropertyMedia) TableName() string { return "property_media" }

// PropertyRelation links a property lineage to a product or asset.
type PropertyRelation struct {
	PropertyID   uint64         `gorm:"column:property_id;primaryKey"`
	RelationID   *uint64        `gorm:"column:relation_id;index"`
	RelationType uint32         `gorm:"column:relation_type;index"`
	UpdatedAt    *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (PropertyRelation) TableName() string { return "property_relation" }

// PropertyUser links a profile to a property lineage for /property/me.
type PropertyUser struct {
	ID                uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt         *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt         *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt         gorm.DeletedAt `gorm:"column:deleted_at;index"`
	OwnerOriginID     uint64         `gorm:"column:owner_origin_id;index;not null"`
	PropertyLineageID uint64         `gorm:"column:property_lineage_id;index;not null"`
	RoleID            uint64         `gorm:"column:role_id;default:0"`
	OriginID          *uint64        `gorm:"column:origin_id"`
	OwneredAt         *time.Time     `gorm:"column:ownered_at"`
	ArchivedAt        *time.Time     `gorm:"column:archived_at"`
	HiddenAt          *time.Time     `gorm:"column:hidden_at"`
}

func (PropertyUser) TableName() string { return "property_user" }

// PropertyIdentify is the identity root for verified properties.
type PropertyIdentify struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
	PID       uint64         `gorm:"column:p_id;not null"`
	Version   uint32         `gorm:"column:version;not null"`
	LineageID *uint64        `gorm:"column:lineage_id"`
}

func (PropertyIdentify) TableName() string { return "property_identify" }

// PropertyType is the catalog of property kinds.
type PropertyType struct {
	ID               uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt        *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt        *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at;index"`
	Name             string         `gorm:"column:name;size:100"`
	Active           bool           `gorm:"column:active;default:true"`
	ClassifyProperty uint32         `gorm:"column:classify_property;default:10"`
	FieldStrs        string         `gorm:"column:field_strs;type:text"`
}

func (PropertyType) TableName() string { return "property_type" }

// PropertyEdvidence is evidence / proof attached to a lineage.
type PropertyEdvidence struct {
	ID                 uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt          *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at;index"`
	PropertyIdentifyID *uint64        `gorm:"column:property_identify_id"`
	Title              string         `gorm:"column:title;size:255"`
	FileID             *uint64        `gorm:"column:file_id"`
	Description        string         `gorm:"column:description;type:text"`
}

func (PropertyEdvidence) TableName() string { return "property_edvidence" }

// PropertyExternalRef is an external system reference.
type PropertyExternalRef struct {
	ID                 uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt          *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at;index"`
	LineageID          uint64         `gorm:"column:lineage_id;index;not null"`
	PropertyIdentifyID *uint64        `gorm:"column:property_identify_id"`
	ExternalRefID      uint64         `gorm:"column:external_ref_id;not null"`
	SourceSystem       uint32         `gorm:"column:source_system"`
	SourceCode         string         `gorm:"column:source_code;size:255"`
	Confidence         uint32         `gorm:"column:confidence;default:0"`
	SyncStatus         uint32         `gorm:"column:sync_status;default:10"`
	Note               string         `gorm:"column:note;type:text"`
}

func (PropertyExternalRef) TableName() string { return "property_external_ref" }

// PropertyProduct links products to a property (legacy join table).
type PropertyProduct struct {
	ID         uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt  *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index"`
	PropertyID uint64         `gorm:"column:property_id;index;not null"`
	ProductID  uint64         `gorm:"column:product_id;index;not null"`
}

func (PropertyProduct) TableName() string { return "property_product" }

// PropertyTagLink links amenity/feature tags to a property.
type PropertyTagLink struct {
	PropertyID uint64 `gorm:"column:property_id;primaryKey"`
	TagID      uint64 `gorm:"column:tag_id;primaryKey"`
}

func (PropertyTagLink) TableName() string { return "property_tag_link" }

// PropertyReport is a user report against a lineage.
type PropertyReport struct {
	ID                 uint64         `gorm:"primaryKey;autoIncrement"`
	CreatedAt          *time.Time     `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt          *time.Time     `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt          gorm.DeletedAt `gorm:"column:deleted_at;index"`
	LineageID          uint64         `gorm:"column:lineage_id;index;not null"`
	ReporterOriginID   uint64         `gorm:"column:reporter_origin_id;index;not null"`
	ReportType         uint32         `gorm:"column:report_type;not null"`
	SubIssues          string         `gorm:"column:sub_issues;type:text"`
	Status             uint32         `gorm:"column:status;not null;default:10"`
	Description        string         `gorm:"column:description;type:text"`
	AttachmentMediaIDs string         `gorm:"column:attachment_media_ids;type:text"`
	AttachmentFileIDs  string         `gorm:"column:attachment_file_ids;type:text"`
	RelatedLineageID   *uint64        `gorm:"column:related_lineage_id"`
	AdminNote          string         `gorm:"column:admin_note;type:text"`
	ReviewerOriginID   *uint64        `gorm:"column:reviewer_origin_id"`
}

func (PropertyReport) TableName() string { return "property_report" }

// PropertyAuditLog records property mutation impact.
type PropertyAuditLog struct {
	ID               uint64           `gorm:"primaryKey;autoIncrement"`
	PropertyID       uint64           `gorm:"column:property_id;index;not null"`
	UserID           uint64           `gorm:"column:user_id;not null"`
	Action           string           `gorm:"column:action;size:50;not null"`
	ImpactLevel      int32            `gorm:"column:impact_level"`
	AffectedProducts int              `gorm:"column:affected_products"`
	AffectedListings int              `gorm:"column:affected_listings"`
	AffectedAssets   int              `gorm:"column:affected_assets"`
	Changes          AuditChangesJSON `gorm:"column:changes;type:jsonb"`
	ConfirmedAt      *time.Time       `gorm:"column:confirmed_at"`
	SessionID        string           `gorm:"column:session_id;size:100"`
	BlockedReason    string           `gorm:"column:blocked_reason;type:text"`
	ProposeLineageID *uint64          `gorm:"column:propose_lineage_id"`
	CreatedAt        time.Time        `gorm:"column:created_at;autoCreateTime"`
}

func (PropertyAuditLog) TableName() string { return "property_audit_logs" }

// PropertyImpactVersion stores entity snapshots for impact analysis.
type PropertyImpactVersion struct {
	ID           uint64     `gorm:"primaryKey;autoIncrement"`
	EntityType   string     `gorm:"column:entity_type;size:50;not null"`
	EntityID     uint64     `gorm:"column:entity_id;not null"`
	PropertyID   uint64     `gorm:"column:property_id;not null"`
	SnapshotData JSONB      `gorm:"column:snapshot_data;type:jsonb;not null"`
	Version      uint32     `gorm:"column:version;default:1"`
	ValidFrom    time.Time  `gorm:"column:valid_from;not null"`
	ValidTo      *time.Time `gorm:"column:valid_to"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime"`
	CreatedBy    uint64     `gorm:"column:created_by;not null"`
}

func (PropertyImpactVersion) TableName() string { return "property_impact_versions" }

// PropertyLineageHistory keeps replaced lineage snapshots after merge.
type PropertyLineageHistory struct {
	ID                     uint64    `gorm:"primaryKey;autoIncrement"`
	LineageID              uint64    `gorm:"column:lineage_id;not null;index"`
	ContributeID           uint64    `gorm:"column:contribute_id;not null;index"`
	ReplacedByLineageID    uint64    `gorm:"column:replaced_by_lineage_id;not null"`
	SnapshotPropertyInfoID *uint64   `gorm:"column:snapshot_property_info_id"`
	SnapshotLocationID     *uint64   `gorm:"column:snapshot_location_id"`
	SnapshotLandInfoID     *uint64   `gorm:"column:snapshot_land_info_id"`
	SnapshotBuildingInfoID *uint64   `gorm:"column:snapshot_building_info_id"`
	SnapshotEdvidenceID    *uint64   `gorm:"column:snapshot_edvidence_id"`
	MergedByID             uint64    `gorm:"column:merged_by_id;not null"`
	MergedAt               time.Time `gorm:"column:merged_at;not null"`
	MergeNote              string    `gorm:"column:merge_note;type:text"`
	CreatedAt              time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (PropertyLineageHistory) TableName() string { return "property_lineage_history" }

// CountryIdentifier is the national property identifier record.
type CountryIdentifier struct {
	ID             string         `gorm:"column:id;primaryKey;size:40"`
	LandParcelCode string         `gorm:"column:land_parcel_code;size:30"`
	ProjectCode    string         `gorm:"column:project_code;size:30"`
	ProvinceID     string         `gorm:"column:province_id;size:5;not null"`
	DistrictID     string         `gorm:"column:district_id;size:7;not null"`
	WardID         string         `gorm:"column:ward_id;size:11;not null"`
	Latitude       float64        `gorm:"column:latitude"`
	Longitude      float64        `gorm:"column:longitude"`
	Type           string         `gorm:"column:type;size:3;not null"`
	LegalStatus    string         `gorm:"column:legal_status;size:3"`
	CurrentOwnerID string         `gorm:"column:current_owner_id;size:20"`
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (CountryIdentifier) TableName() string { return "property_identifiers" }

// Alias kept for existing sketches.
type Property = PropertyLineage

// Relation type constants (property_relation.relation_type).
const (
	RelationTypeProduct uint32 = 10
	RelationTypeAsset   uint32 = 20
)

// ProductMediaInput is request-only media carried through Product.Create (not a table).
type ProductMediaInput struct {
	MediaURL  string
	MediaType string
	IsMain    bool
	Order     int32
}

// JSONB / AuditChangesJSON helpers for audit + impact tables.
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

type AuditChangesJSON map[string]interface{}

func (j AuditChangesJSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *AuditChangesJSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}
