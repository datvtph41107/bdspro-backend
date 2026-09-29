package dto

// GroupDTO represents the Data Transfer Object for FriendGroup
type GroupDTO struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	TextColor       string `json:"textColor"`
	BackgroundColor string `json:"backgroundColor"`
}

// MakeGroupDTO is a constructor for creating a new FriendGroupDTO instance
func MakeGroupDTO(id uint, name, textColor, backgroundColor string) *GroupDTO {
	return &GroupDTO{
		ID:              id,
		Name:            name,
		TextColor:       textColor,
		BackgroundColor: backgroundColor,
	}
}
