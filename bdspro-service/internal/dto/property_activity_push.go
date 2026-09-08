package dto

import _enum "common/domain/enum"

type CreatePropertyActivityDTO struct {
	SubjectID   uint64                    `json:"subjectId"`
	SubjectType _enum.SubjectTypeProperty `json:"subjectType"`
	Action      _enum.PropertyAction      `json:"action"`
	ActorID     uint64                    `json:"actorId"`
	Description string                    `json:"description"`
}
