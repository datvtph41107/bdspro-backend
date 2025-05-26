package models

type Like struct {
	ProfileID *uint64
	PostID    *uint64
}

func (Like) TableName() string {
	return "likes"
}
