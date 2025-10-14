package models

import (
	_models "common/domain/entity"
)

type Plan struct {
	_models.BaseEntity
	PackageName        string   `gorm:"column:package_name;not null" json:"packageName"`
	Duration           int64    `gorm:"column:duration;not null" json:"duration"` // day
	Price              int64    `gorm:"column:price;not null" json:"price"`
	CreateProduct      int      `gorm:"column:create_product;default:0" json:"createProduct"`
	LimitPost          int64    `gorm:"column:post_limit;default:0" json:"limitPost"`
	PublicProductLimit int      `gorm:"column:public_product_limit;default:0" json:"publicProductLimit"`
	CreateProductAI    int      `gorm:"column:create_product_ai;default:0" json:"createProductAI"`
	FriendLimit        int      `gorm:"column:friend_limit;default:0" json:"friendLimit"`
	CustomerLimit      int      `gorm:"column:customer_limit;default:0" json:"customerLimit"`
	JoinOrg            int      `gorm:"column:join_org;default:0" json:"joinOrg"`
	CreateOrganization int      `gorm:"column:create_organization;default:0" json:"createOrganization"`
	OrgMemberLimit     int      `gorm:"column:org_member_limit;default:0" json:"orgMemberLimit"`
	OrgProductLimit    int      `gorm:"column:org_product_limit;default:0" json:"orgProductLimit"`
	Color              string   `gorm:"column:color" json:"color"`
	Features           []string `gorm:"type:text[];column:features" json:"features"`
	Visibility         bool     `gorm:"column:visibility;default:true" json:"visibility"`
}

// TableName overrides the table name used by GORM
func (Plan) TableName() string {
	return "plans"
}

