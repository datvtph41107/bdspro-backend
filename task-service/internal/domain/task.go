package domain

import "time"

type Task struct {
	Type     uint       `json:"type,omitempty"`
	TargetId uint64     `json:"targetId,omitempty"`
	Service  string     `json:"service,omitempty"`
	EndAt    *time.Time `json:"endAt,omitempty"`
	Delay    uint64     `json:"delay,omitempty"`
}
