package dto

// AdminProfileIDPath is path `{profileId}` for admin user routes.
type AdminProfileIDPath struct {
	ProfileID uint64 `path:"profileId"`
}

// AdminIDPath is path `{id}` for admin-profile routes.
type AdminIDPath struct {
	ID uint64 `path:"id"`
}

// ListAllUsersRequest is GET /v2/user/admin/users.
type ListAllUsersRequest struct {
	Search    string `form:"search"`
	Page      uint32 `form:"page"`
	Size      uint32 `form:"size"`
	SortBy    string `form:"sortBy"`
	SortOrder string `form:"sortOrder"`
	Status    uint32 `form:"status"`
	RoleType  string `form:"roleType"`
	StartDate string `form:"startDate"`
	EndDate   string `form:"endDate"`
	Bookmark  *bool  `form:"bookmark"`
	RoleID    uint64 `form:"roleId"`
	Sort      string `form:"sort"`
}

// GetUserStatsRequest is GET /v2/user/stats.
type GetUserStatsRequest struct {
	FirstDate  string `form:"firstDate"`
	SecondDate string `form:"secondDate"`
}

// AdminProfileWriteRequest is POST create / PUT update admin-profile.
type AdminProfileWriteRequest struct {
	ID               uint64   `path:"id" json:"id"`
	FullName         string   `json:"fullName"`
	Email            string   `json:"email"`
	Phone            string   `json:"phone"`
	Avatar           string   `json:"avatar"`
	Address          string   `json:"address"`
	Gender           *uint32  `json:"gender"`
	Birth            string   `json:"birth"`
	Status           uint32   `json:"status"`
	RoleType         string   `json:"roleType"`
	RoleRealEstate   string   `json:"roleRealEstate"`
	Position         string   `json:"position"`
	DepartmentID     uint64   `json:"departmentId"`
	Bookmark         bool     `json:"bookmark"`
	TickVerified     bool     `json:"tickVerified"`
	Attachments      []string `json:"attachments"`
	InternalNotes    string   `json:"internalNotes"`
	SendNotification bool     `json:"sendNotification"`
	JobTitle         string   `json:"jobTitle"`
	WorkAt           string   `json:"workAt"`
	RoleKey          string   `json:"roleKey"`
	RoleID           *uint64  `json:"roleId"`
	Username         string   `json:"username"`
	Password         string   `json:"password"`
	SuperAdmin       uint32   `json:"superAdmin"`
}
