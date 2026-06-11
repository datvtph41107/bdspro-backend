package seo_domain

type SeoRelative struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	ParentSeoID uint64 `gorm:"not null"`

	ChildSeoID uint64 `gorm:"not null"`

	ChildSeo *SeoDomain `gorm:"foreignKey:ChildSeoID;references:ID"`
}

func (SeoRelative) TableName() string {
	return "seo_relative"
}
