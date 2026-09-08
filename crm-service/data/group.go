package data

// GroupDTO represents the Data Transfer Object for FriendGroup
type GroupDTO struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	TextColor       string `json:"textColor"`
	BackgroundColor string `json:"backgroundColor"`
}

// NewGroupDTO is a constructor for creating a new FriendGroupDTO instance
// func NewGroupDTO(id uint, name, textColor, backgroundColor string) *GroupDTO {
// 	return &GroupDTO{
// 		ID:              id,
// 		Name:            name,
// 		TextColor:       textColor,
// 		BackgroundColor: backgroundColor,
// 	}
// }