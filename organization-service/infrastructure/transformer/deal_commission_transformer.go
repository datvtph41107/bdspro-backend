package transformer

import (
	"organization/internal/domain/entity"
	"organization/internal/enums"
	organizationpb "pb/types/organization"
)

type DealCommissionTransformer interface {
	TransformUpdateCommissionRequest(req *organizationpb.UpdateDealMemberCommissionRequest) (float64, enums.CommissionType, string)
	TransformUpdateCommissionResponse(req *organizationpb.UpdateDealMemberCommissionRequest) *organizationpb.UpdateDealMemberCommissionResponse
	TransformUpdateNoteRequest(req *organizationpb.UpdateDealMemberNoteRequest) string
	TransformUpdateNoteResponse(req *organizationpb.UpdateDealMemberNoteRequest) *organizationpb.UpdateDealMemberNoteResponse
	TransformUpdateMembersCommissionRequest(req *organizationpb.UpdateDealMembersCommissionRequest) []*entity.DealMember
	TransformUpdateMembersCommissionResponse(req *organizationpb.UpdateDealMembersCommissionRequest) *organizationpb.UpdateDealMembersCommissionResponse
	TransformCommissionStatsResponse(stats *entity.CommissionStats) *organizationpb.GetDealCommissionStatsResponse
}

type dealCommissionTransformer struct{}

func NewDealCommissionTransformer() DealCommissionTransformer {
	return &dealCommissionTransformer{}
}

func (t *dealCommissionTransformer) TransformUpdateCommissionRequest(req *organizationpb.UpdateDealMemberCommissionRequest) (float64, enums.CommissionType, string) {
	commissionType := enums.CommissionType(req.CommissionType)
	if commissionType == 0 {
		commissionType = enums.CommissionTypePercent // Default to percent
	}

	return req.CommissionValue, commissionType, req.Note
}

func (t *dealCommissionTransformer) TransformUpdateCommissionResponse(req *organizationpb.UpdateDealMemberCommissionRequest) *organizationpb.UpdateDealMemberCommissionResponse {
	return &organizationpb.UpdateDealMemberCommissionResponse{
		DealId:         req.DealId,
		MemberId:       req.MemberId,
		CommissionValue: req.CommissionValue,
		CommissionType:  req.CommissionType,
		Note:           req.Note,
	}
}

func (t *dealCommissionTransformer) TransformUpdateNoteRequest(req *organizationpb.UpdateDealMemberNoteRequest) string {
	return req.Note
}

func (t *dealCommissionTransformer) TransformUpdateNoteResponse(req *organizationpb.UpdateDealMemberNoteRequest) *organizationpb.UpdateDealMemberNoteResponse {
	return &organizationpb.UpdateDealMemberNoteResponse{
		DealId:   req.DealId,
		MemberId: req.MemberId,
		Note:     req.Note,
	}
}

func (t *dealCommissionTransformer) TransformUpdateMembersCommissionRequest(req *organizationpb.UpdateDealMembersCommissionRequest) []*entity.DealMember {
	members := make([]*entity.DealMember, len(req.Members))
	for i, memberReq := range req.Members {
		commissionType := enums.CommissionType(memberReq.CommissionType)
		if commissionType == 0 {
			commissionType = enums.CommissionTypePercent
		}
		
		members[i] = &entity.DealMember{
			MemberID:        memberReq.MemberId,
			CommissionValue: memberReq.CommissionValue,
			CommissionType:  commissionType,
			Note:            memberReq.Note,
		}
	}
	return members
}

func (t *dealCommissionTransformer) TransformUpdateMembersCommissionResponse(req *organizationpb.UpdateDealMembersCommissionRequest) *organizationpb.UpdateDealMembersCommissionResponse {
	return &organizationpb.UpdateDealMembersCommissionResponse{
		DealId:   req.DealId,
		Members:  req.Members,
	}
}

func (t *dealCommissionTransformer) TransformCommissionStatsResponse(stats *entity.CommissionStats) *organizationpb.GetDealCommissionStatsResponse {
	return &organizationpb.GetDealCommissionStatsResponse{
		DealId:              stats.DealID,
		TargetProfit:        stats.TargetProfit,
		TotalCommission:     stats.TotalCommission,
		MemberCount:         stats.MemberCount,
		RemainingProfit:     stats.RemainingProfit,
		CommissionPercentage: stats.CommissionPercentage,
	}
} 