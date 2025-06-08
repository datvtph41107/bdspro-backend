package dto

import (
	_dto "common/dto"
)

type ReportReasonRequest struct {
	_dto.Page
	ReasonName string `json:"reasonName"`
}
