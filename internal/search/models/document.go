package models

type Document struct {
	Text       string `json:"text"`
	TargetID   uint64 `json:"targetId"`
	Type       uint   `json:"type"`
	Popularity int64  `json:"popularity"`
}
