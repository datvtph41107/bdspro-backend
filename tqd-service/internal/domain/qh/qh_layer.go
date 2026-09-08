package qh_domain

import (
	"encoding/json"
	"time"
	"tqd/internal/domain/jsonb"
	"tqd/internal/enums"

	"github.com/lib/pq"
)

type QHLayer struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"type:varchar(100);" json:"name"`
	DisplayName string `gorm:"type:varchar(255);not null" json:"displayName"`
	Description string `gorm:"type:text" json:"description,omitempty"`

	PlanningProjectID *uint64            `gorm:"column:planning_project_id;index" json:"planningProjectId,omitempty"`
	PlanningProject   *QHPlanningProject `gorm:"foreignKey:PlanningProjectID;references:ID"`

	Type   enums.LayerType   `gorm:"type:int;not null;index" json:"type"`
	Status enums.LayerStatus `gorm:"type:int;default:10;index" json:"status"`

	DisplayOrder int   `gorm:"column:display_order;default:0;index" json:"displayOrder"`
	AnalyseOrder int   `gorm:"column:analyse_order;default:10;index" json:"analyseOrder"`
	Visible      int32 `gorm:"default:0;index" json:"visible"`

	MinZoom uint32 `gorm:"default:6" json:"minZoom"`
	MaxZoom uint32 `gorm:"default:16" json:"maxZoom"`

	Avatar       string `gorm:"type:text" json:"avatar,omitempty"`
	ImageURL     string `gorm:"type:text" json:"imageUrl,omitempty"`
	ThumbnailURL string `gorm:"type:text" json:"thumbnailUrl,omitempty"`

	ImportStatus  enums.LayerImportStatus `gorm:"column:import_status;type:int;default:20;index" json:"importStatus"`
	ImportBatchID *string                 `gorm:"column:import_batch_id;type:varchar(64);index" json:"importBatchId,omitempty"`

	SourceCode  string      `gorm:"type:varchar(100);index" json:"sourceCode,omitempty"`
	SourceType  string      `gorm:"type:varchar(20);not null;default:raster" json:"sourceType,omitempty"`
	LayerURL    string      `gorm:"type:text" json:"layerUrl,omitempty"`
	StyleConfig jsonb.JSONB `gorm:"type:jsonb;default:'[]'" json:"styleConfig,omitempty"`

	InsuanceDate  *time.Time `gorm:"type:date;index" json:"insuanceDate,omitempty"`  // ngày ban hành
	EffectiveDate *time.Time `gorm:"type:date;index" json:"effectiveDate,omitempty"` // ngày hiệu lực
	ExpiryDate    *time.Time `gorm:"type:date;index" json:"expiryDate,omitempty"`    // ngày hết hiệu lực

	CreatedBy uint64     `gorm:"default:0" json:"createdBy"`
	UpdatedBy uint64     `gorm:"default:0" json:"updatedBy"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"deletedAt,omitempty"`

	Labels []QHLabel `gorm:"many2many:qh_label_layers;" json:"labels,omitempty"`

	TrustValue   float32           `gorm:"default:0" json:"trustValue"`
	LegalStatus  enums.LegalStatus `gorm:"default:0" json:"legalStatus"`
	ReplacedByID *uint64           `gorm:"column:replaced_by_id;index" json:"replaceLayerId,omitempty"`

	LegalLevel  enums.LegalLevel `gorm:"default:0" json:"legalLevel"`
	LegalNumber string           `gorm:"type:varchar(50)" json:"legalNumber,omitempty"` // số/ký hiệu văn bản
	LegalDoc    string           `gorm:"type:text" json:"legalDoc,omitempty"`           // văn bản pháp lý (nội dung/mô tả)

	DefaultVisible bool `gorm:"default:false" json:"defaultVisible"`

	AuthorityIssuringID *uint64              `gorm:"type:bigint" json:"authorityIssuring,omitempty"` // cơ quan ban hành
	AuthorityIssuring   *QHAuthorityIssuring `gorm:"foreignKey:AuthorityIssuringID;references:ID"`
	PublishScopes       pq.Int32Array        `gorm:"type:integer[]" json:"publishScopes,omitempty"` // phạm vi áp dụng

	FamilyID *uint64        `gorm:"column:family_id;index" json:"familyId,omitempty"`
	Family   *QHLayerFamily `gorm:"foreignKey:FamilyID;references:ID"`
	ParentID *uint64        `gorm:"column:parent_id;index" json:"parentId,omitempty"`
	Parent   *QHLayer       `gorm:"foreignKey:ParentID;references:ID"`
}

func (QHLayer) TableName() string {
	return "qh_layers"
}

func (l *QHLayer) IsActive() bool {
	if l.DeletedAt != nil {
		return false
	}
	if l.Status != enums.LayerStatusActive {
		return false
	}
	if l.Visible != 1 {
		return false
	}
	now := time.Now()
	if l.EffectiveDate != nil && l.EffectiveDate.After(now) {
		return false
	}
	if l.ExpiryDate != nil && l.ExpiryDate.Before(now) {
		return false
	}
	return true
}

func (l *QHLayer) GetActiveLabelsQuery() string {
	return "layer_id = ? AND status = 10 AND deleted_at IS NULL"
}

func (l *QHLayer) GetCanonicalLegalStatus() string {
	return l.LegalStatus.String()
}

type StyleConfig struct {
	FillColor     string  `json:"fillColor,omitempty"`
	FillOpacity   float64 `json:"fillOpacity,omitempty"`
	StrokeColor   string  `json:"strokeColor,omitempty"`
	StrokeWidth   int     `json:"strokeWidth,omitempty"`
	LineColor     string  `json:"lineColor,omitempty"`
	LineWidth     int     `json:"lineWidth,omitempty"`
	LineOpacity   float64 `json:"lineOpacity,omitempty"`
	CircleColor   string  `json:"circleColor,omitempty"`
	CircleRadius  int     `json:"circleRadius,omitempty"`
	CircleOpacity float64 `json:"circleOpacity,omitempty"`
}

func (s StyleConfig) ToJSONB() jsonb.JSONB {
	data, _ := json.Marshal(s)
	var result jsonb.JSONB
	json.Unmarshal(data, &result)
	return result
}

func GetDefaultStyleConfig() json.RawMessage {
	return json.RawMessage(`[
	  {
	    "id": "layer68_fill",
	    "type": "fill",
	    "paint": {
	      "fill-color": "rgba(255,0,0,0.4)",
	      "fill-opacity": 0.2
	    },
	    "source": "layer68",
	    "source-layer": "layer68"
	  },
	  {
	    "id": "layer68_outline",
	    "type": "line",
	    "paint": {
	      "line-color": "rgba(255,0,0,0.9)",
	      "line-width": 2
	    },
	    "source": "layer68",
	    "source-layer": "layer68"
	  }
	]`)
}

func (l *QHLayer) IsClientVisible() bool {
	return l.IsActive() && l.Visible == 1
}
