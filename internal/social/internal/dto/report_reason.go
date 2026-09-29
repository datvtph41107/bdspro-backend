package dto

import (
	_dto "common/domain/dto"
)

type ReportReasonRequest struct {
	_dto.Pagable
	ReasonName string `json:"reasonName"`
}
