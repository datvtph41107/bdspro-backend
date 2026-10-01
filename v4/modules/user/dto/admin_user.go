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

// LockUserRequest is POST /v2/user/admin/users/{profileId}/lock.
type LockUserRequest struct {
	ProfileID uint64 `path:"profileId" json:"profileId"`
	Reason    string `json:"reason"`
	LockType  uint32 `json:"lockType"`
}

// UnLockUserRequest is POST /v2/user/admin/users/{profileId}/unlock.
type UnLockUserRequest struct {
	ProfileID uint64 `path:"profileId" json:"profileId"`
	Reason    string `json:"reason"`
}

// ApproveUserRequest is POST /v2/user/admin/users/{profileId}/approve.
type ApproveUserRequest struct {
	ProfileID uint64 `path:"profileId" json:"profileId"`
	Note      string `json:"note"`
}

// GetUserStatsRequest is GET /v2/user/stats.
type GetUserStatsRequest struct {
	FirstDate  string `form:"firstDate"`
	SecondDate string `form:"secondDate"`
}

// GetUserGrowthDailyStatsRequest is GET /v2/user/admin/users/dashboard/daily.
type GetUserGrowthDailyStatsRequest struct {
	LastDates uint32 `form:"lastDates"`
	Gap       string `form:"gap"`
}

// UpdateUserRequest is PUT /v2/user/admin/users/{profileId}.
type UpdateUserRequest struct {
	ProfileID    uint64 `path:"profileId" json:"profileId"`
	FullName     string `json:"fullName"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Avatar       string `json:"avatar"`
	RoleID       uint64 `json:"roleId"`
	Gender       uint32 `json:"gender"`
	Address      string `json:"address"`
	Birth        string `json:"birth"`
	Status       uint32 `json:"status"`
	IDVerifiedAt string `json:"idVerifiedAt"`
}

// DeleteUserRequest is DELETE /v2/user/admin/users/{profileId}.
type DeleteUserRequest struct {
	ProfileID uint64 `path:"profileId" json:"profileId"`
	Reason    string `json:"reason"`
}

// AdminListRequest is GET /v2/user/admin-profile/list.
type AdminListRequest struct {
	Page   uint32 `form:"page"`
	Size   uint32 `form:"size"`
	RoleID uint64 `form:"roleId"`
	Name   string `form:"name"`
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
