package models

// BookmarkUserEntity đại diện cho bookmark user của admin
type BookmarkUserEntity struct {
	AdminID uint64 `gorm:"primaryKey;column:admin_id;not null" json:"adminId"`
	UserID  uint64 `gorm:"primaryKey;column:user_id;not null" json:"userId"`

	// Relations
	User *UserProfileEntity `gorm:"foreignKey:UserID;references:ProfileID" json:"user,omitempty"`
}

// TableName đặt tên bảng trong DB
func (BookmarkUserEntity) TableName() string {
	return "bookmark_user"
}

// BookmarkUserListRequest request để lấy danh sách bookmark user
type BookmarkUserListRequest struct {
	AdminID uint64 `json:"adminId" binding:"required"`
	Page    int    `json:"page" binding:"min=1"`
	Size    int    `json:"size" binding:"min=1,max=100"`
}

// BookmarkUserUpdateRequest request để cập nhật bookmark user
type BookmarkUserUpdateRequest struct {
	AdminID uint64   `json:"adminId" binding:"required"`
	UserIDs []uint64 `json:"userIds" binding:"required"`
}
