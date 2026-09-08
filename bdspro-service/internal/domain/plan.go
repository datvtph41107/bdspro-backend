package domain

// not table
type MemberPlan struct {
	ID                 uint64   `json:"id"`
	PackageName        string   `json:"packageName"`
	Duration           int64    `json:"duration"` // day
	Price              int64    `json:"price"`
	CreateProduct      int      `json:"createProduct"`
	LimitPost          int64    `json:"limitPost"`
	PublicProductLimit int      `json:"publicProductLimit"`
	CreateProductAI    int      `json:"createProductAI"`
	FriendLimit        int      `json:"friendLimit"`
	CustomerLimit      int      `json:"customerLimit"`
	JoinOrg            int      `json:"joinOrg"`
	CreateOrganization int      `json:"createOrganization"`
	OrgMemberLimit     int      `json:"orgMemberLimit"`
	OrgProductLimit    int      `json:"orgProductLimit"`
	Color              string   `json:"color"`
	Features           []string `json:"features"`
	Visibility         bool     `json:"visibility"`
}
