package dto

import (
	_dto "common/dto"
)

type ReportReasonRequest struct {
	_dto.Pagable
	ReasonName string `json:"reasonName"`
}
