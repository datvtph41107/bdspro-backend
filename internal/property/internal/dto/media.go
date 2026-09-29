package dto

type MediaItem struct {
	ID        uint64 `json:"id"`
	MediaURL  string `json:"mediaUrl"`
	MediaType string `json:"mediaType"`
	IsMain    bool   `json:"isMain"`
	Order     int    `json:"order"`
}
