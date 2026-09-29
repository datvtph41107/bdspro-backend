package mapper

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
	bdspropb "pb/types/bdspro"
)

type DealCommissionTransformer interface {
	TransformUpdateCommissionRequest(req *bdspropb.UpdateDealMemberCommissionRequest) (float64, enums.CommissionType, string)
	TransformUpdateCommissionResponse(req *bdspropb.UpdateDealMemberCommissionRequest) *bdspropb.UpdateDealMemberCommissionResponse
	TransformUpdateNoteRequest(req *bdspropb.UpdateDealMemberNoteRequest) string
	TransformUpdateNoteResponse(req *bdspropb.UpdateDealMemberNoteRequest) *bdspropb.UpdateDealMemberNoteResponse
	TransformUpdateMembersCommissionRequest(req *bdspropb.UpdateDealMembersCommissionRequest) []*domain.DealMember
	TransformUpdateMembersCommissionResponse(req *bdspropb.UpdateDealMembersCommissionRequest) *bdspropb.UpdateDealMembersCommissionResponse
	TransformCommissionStatsResponse(stats *domain.CommissionStats) *bdspropb.GetDealCommissionStatsResponse
}

type dealCommissionTransformer struct{}

func NewDealCommissionTransformer() DealCommissionTransformer {
	return &dealCommissionTransformer{}
}

func (t *dealCommissionTransformer) TransformUpdateCommissionRequest(req *bdspropb.UpdateDealMemberCommissionRequest) (float64, enums.CommissionType, string) {
	commissionType := enums.CommissionType(req.CommissionType)
	if commissionType == 0 {
		commissionType = enums.CommissionTypePercent // Default to percent
	}

	return req.CommissionValue, commissionType, req.Note
}

func (t *dealCommissionTransformer) TransformUpdateCommissionResponse(req *bdspropb.UpdateDealMemberCommissionRequest) *bdspropb.UpdateDealMemberCommissionResponse {
	return &bdspropb.UpdateDealMemberCommissionResponse{
		DealId:          req.DealId,
		MemberId:        req.MemberId,
		CommissionValue: req.CommissionValue,
		CommissionType:  req.CommissionType,
		Note:            req.Note,
	}
}

func (t *dealCommissionTransformer) TransformUpdateNoteRequest(req *bdspropb.UpdateDealMemberNoteRequest) string {
	return req.Note
}

func (t *dealCommissionTransformer) TransformUpdateNoteResponse(req *bdspropb.UpdateDealMemberNoteRequest) *bdspropb.UpdateDealMemberNoteResponse {
	return &bdspropb.UpdateDealMemberNoteResponse{
		DealId:   req.DealId,
		MemberId: req.MemberId,
		Note:     req.Note,
	}
}

func (t *dealCommissionTransformer) TransformUpdateMembersCommissionRequest(req *bdspropb.UpdateDealMembersCommissionRequest) []*domain.DealMember {
	members := make([]*domain.DealMember, len(req.Members))
	for i, memberReq := range req.Members {
		commissionType := enums.CommissionType(memberReq.CommissionType)
		if commissionType == 0 {
			commissionType = enums.CommissionTypePercent
		}

		members[i] = &domain.DealMember{
			MemberID:        memberReq.MemberId,
			CommissionValue: memberReq.CommissionValue,
			CommissionType:  commissionType,
			Note:            memberReq.Note,
		}
	}
	return members
}

func (t *dealCommissionTransformer) TransformUpdateMembersCommissionResponse(req *bdspropb.UpdateDealMembersCommissionRequest) *bdspropb.UpdateDealMembersCommissionResponse {
	return &bdspropb.UpdateDealMembersCommissionResponse{
		DealId:  req.DealId,
		Members: req.Members,
	}
}

func (t *dealCommissionTransformer) TransformCommissionStatsResponse(stats *domain.CommissionStats) *bdspropb.GetDealCommissionStatsResponse {
	return &bdspropb.GetDealCommissionStatsResponse{
		DealId:               stats.DealID,
		TargetProfit:         stats.TargetProfit,
		TotalCommission:      stats.TotalCommission,
		MemberCount:          stats.MemberCount,
		RemainingProfit:      stats.RemainingProfit,
		CommissionPercentage: stats.CommissionPercentage,
	}
}
