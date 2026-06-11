package seo_domain

type SeoInternalLink struct {
	ID          uint64     `gorm:"primaryKey;autoIncrement"`
	Link        string     `gorm:""`
	Title       string     `gorm:""`
	ParentSeoID uint64     `gorm:""`
	ChildSeoID  uint64     `gorm:""`
	ParentSeo   *SeoDomain `gorm:"foreignKey:ParentSeoID;references:ID"`
	ChildSeo    *SeoDomain `gorm:"foreignKey:ChildSeoID;references:ID"`
}

func (SeoInternalLink) TableName() string {
	return "seo_links"
}
