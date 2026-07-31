package dto

type ColorRequest struct {
	Name    string `json:"name" validate:"required"`
	HexCode string `json:"hexCode" validate:"required"`
}

type ColorResponse struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	HexCode   string `json:"hexCode"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
