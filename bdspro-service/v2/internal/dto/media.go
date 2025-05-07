package dto

type MediaItem struct {
	MediaURL  string `json:"mediaUrl"`
	MediaType string `json:"mediaType"`
	IsMain    bool   `json:"isMain"`
	Order     int    `json:"order"`
}
