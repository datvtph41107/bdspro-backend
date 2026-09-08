package dto

import (
	"bdspro/internal/enums"
	_dto "common/domain/dto"
	"time"
)

// CreateAssetForUserDTO DTO cho tạo asset cho user
type CreateAssetForUserDTO struct {
	Asset     AssetSaveRequestDTO `json:"asset"`
	ProfileID uint64              `json:"profileId"`
	IsOwner   bool                `json:"isOwner"`
	RoleID    uint64              `json:"roleId"`
}

// CreateAssetForOrgDTO DTO cho tạo asset cho organization
type CreateAssetForOrgDTO struct {
	Asset          AssetSaveRequestDTO `json:"asset"`
	OrganizationID uint64              `json:"organizationId"`
	IsOwner        bool                `json:"isOwner"`
	RoleID         uint64              `json:"roleId"`
}

// AssetSaveRequestDTO DTO cho thông tin asset cần tạo
type AssetSaveRequestDTO struct {
	Name             string             `json:"name" binding:"required"`
	WardID           *uint64            `json:"wardId,omitempty"`
	DistrictID       *uint64            `json:"districtId" binding:"required"`
	ProvinceID       *uint64            `json:"provinceId" binding:"required"`
	ProductID        *uint64            `json:"productId,omitempty"`
	ParentAssetID    *uint64            `json:"parentAssetId,omitempty"`
	SplitMergeStatus string             `json:"splitMergeStatus,omitempty"`
	PurchasePrice    float64            `json:"purchasePrice" binding:"required"`
	PurchaseDate     *time.Time         `json:"purchaseDate,omitempty"`
	LegalStatus      enums.EDocType     `json:"legalStatus,omitempty"`
	Archived         bool               `json:"archived"`
	RentStatus       enums.EAssetStatus `json:"rentStatus,omitempty"`
	Address          string             `json:"address,omitempty"`
	Description      string             `json:"description,omitempty"`
	Area             float64            `json:"area" binding:"required"`
	PropertyTypeID   *uint64            `json:"propertyTypeId,omitempty"`
	OwnerType        enums.EOwnerOf     `json:"ownerType"`
	LegalItems       []AssetLegalDTO    `json:"legalItems,omitempty"`
	ID               uint64             `json:"id"`
}

// AssetLegalDTO DTO cho thông tin pháp lý tài sản
type AssetLegalDTO struct {
	DocumentType string     `json:"documentType"`
	DocumentNo   string     `json:"documentNo"`
	IssueDate    *time.Time `json:"issueDate,omitempty"`
	IssueBy      string     `json:"issueBy,omitempty"`
	ExpireDate   *time.Time `json:"expireDate,omitempty"`
	Status       string     `json:"status,omitempty"`
	Note         string     `json:"note,omitempty"`
}

// CreateAssetResponseDTO DTO cho response tạo asset
type CreateAssetResponseDTO struct {
	AssetID uint64 `json:"assetId"`
	Message string `json:"message"`
}

// AssetListRequestDTO DTO cho request danh sách tài sản
type AssetListRequestDTO struct {
	_dto.Pagable
	OrganizationID *uint64             `json:"organizationId,omitempty"`
	ProfileID      *uint64             `json:"profileId,omitempty"`
	OwnerType      *enums.EOwnerOf     `json:"ownerType,omitempty"`
	PropertyTypeID *uint64             `json:"propertyTypeId,omitempty"`
	ProvinceID     *uint64             `json:"provinceId,omitempty"`
	DistrictID     *uint64             `json:"districtId,omitempty"`
	WardID         *uint64             `json:"wardId,omitempty"`
	Archived       *bool               `json:"archived,omitempty"`
	RentStatus     *enums.EAssetStatus `json:"rentStatus,omitempty"`
	Search         string              `json:"search,omitempty"`
}

// AssetListResponseDTO DTO cho response danh sách tài sản
type AssetListResponseDTO struct {
	_dto.Pagable
	Data []AssetDetailDTO `json:"data"`
}

// AssetDetailDTO DTO cho chi tiết tài sản
type AssetDetailDTO struct {
	ID               uint64             `json:"id"`
	Name             string             `json:"name"`
	WardID           *uint64            `json:"wardId"`
	WardName         string             `json:"wardName,omitempty"`
	DistrictID       *uint64            `json:"districtId,omitempty"`
	DistrictName     string             `json:"districtName,omitempty"`
	ProvinceID       *uint64            `json:"provinceId,omitempty"`
	ProvinceName     string             `json:"provinceName,omitempty"`
	ProductID        *uint64            `json:"productId,omitempty"`
	ParentAssetID    *uint64            `json:"parentAssetId,omitempty"`
	SplitMergeStatus string             `json:"splitMergeStatus,omitempty"`
	PurchasePrice    float64            `json:"purchasePrice"`
	PurchaseDate     *time.Time         `json:"purchaseDate,omitempty"`
	LegalStatus      enums.EDocType     `json:"legalStatus,omitempty"`
	Archived         bool               `json:"archived"`
	RentStatus       enums.EAssetStatus `json:"rentStatus,omitempty"`
	Address          string             `json:"address,omitempty"`
	Description      string             `json:"description,omitempty"`
	Area             float64            `json:"area"`
	PropertyTypeID   *uint64            `json:"propertyTypeId,omitempty"`
	PropertyTypeName string             `json:"propertyTypeName,omitempty"`
	OwnerID          *uint64            `json:"ownerId,omitempty"`
	OwnerType        enums.EOwnerOf     `json:"ownerType"`
	LegalItems       []AssetLegalDTO    `json:"legalItems,omitempty"`
	CreatedAt        *time.Time         `json:"createdAt,omitempty"`
	UpdatedAt        *time.Time         `json:"updatedAt,omitempty"`
	Permissions      []string           `json:"permissions,omitempty"`
}
