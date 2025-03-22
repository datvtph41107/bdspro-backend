package models

import _models "common/models"

type Plan struct {
	_models.BaseEntity
	ID                 uint     `gorm:"primaryKey;autoIncrement" json:"id"`
	PackageName        string   `gorm:"column:package_name" json:"packageName"`
	Duration           int64    `gorm:"column:duration" json:"duration"` // day
	Price              int64    `gorm:"column:price" json:"price"`
	CreateProduct      int      `gorm:"column:create_product" json:"createProduct"`
	PublicProductLimit int      `gorm:"column:public_product_limit" json:"publicProductLimit"`
	CreateProductAI    int      `gorm:"column:create_product_ai" json:"createProductAI"`
	FriendLimit        int      `gorm:"column:friend_limit" json:"friendLimit"`
	CustomerLimit      int      `gorm:"column:customer_limit" json:"customerLimit"`
	JoinOrg            int      `gorm:"column:join_org" json:"joinOrg"`
	CreateOrganization int      `gorm:"column:create_organization" json:"createOrganization"`
	OrgMemberLimit     int      `gorm:"column:org_member_limit" json:"orgMemberLimit"`
	OrgProductLimit    int      `gorm:"column:org_product_limit" json:"orgProductLimit"`
	Color              string   `gorm:"column:color" json:"color"`
	Features           []string `gorm:"type:text[];column:features" json:"features"`
	Visibility         bool     `json:"visibility"`
}

// TableName overrides the table name used by GORM
func (Plan) TableName() string {
	return "plans"
}
