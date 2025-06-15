package dto

import _dto "common/dto"

type RuleSearchDTO struct {
	_dto.Pagable
	RuleName string `json:"ruleName"`
}
