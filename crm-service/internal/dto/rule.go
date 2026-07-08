package dto

import _dto "common/domain/dto"

type RuleSearchDTO struct {
	_dto.Pagable
	RuleName string `json:"ruleName"`
}