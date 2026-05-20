package mapper

// type RegionLabelMapper struct{}

// func NewRegionLabelMapper() *RegionLabelMapper {
// 	return &RegionLabelMapper{}
// }

// // ToProtoBatchAssignResponse - Convert batch assign result to proto response
// func (m *RegionLabelMapper) ToProtoBatchAssignResponse(successCount, failureCount int, failures []dto.BatchAssignFailure, message string) *tqdpb.BatchAssignLabelsResponse {
// 	protoFailures := make([]*tqdpb.BatchAssignFailure, len(failures))
// 	for i, f := range failures {
// 		protoFailures[i] = &tqdpb.BatchAssignFailure{
// 			RegionId: f.RegionID,
// 			Reason:   f.Reason,
// 		}
// 	}

// 	return &tqdpb.BatchAssignLabelsResponse{
// 		SuccessCount: int32(successCount),
// 		FailureCount: int32(failureCount),
// 		Message:      message,
// 		Failures:     protoFailures,
// 	}
// }

// // ToProtoAssignResponse - Convert assign result to proto response
// func (m *RegionLabelMapper) ToProtoAssignResponse(regionId uint64, labelIds []uint64, message string) *tqdpb.AssignLabelsToRegionResponse {
// 	return &tqdpb.AssignLabelsToRegionResponse{
// 		Success:  true,
// 		Message:  message,
// 		RegionId: regionId,
// 		LabelIds: labelIds,
// 	}
// }

// // ToProtoAddResponse - Convert add labels result to proto response
// func (m *RegionLabelMapper) ToProtoAddResponse(regionId uint64, labelIds []uint64, message string) *tqdpb.AddLabelsToRegionResponse {
// 	return &tqdpb.AddLabelsToRegionResponse{
// 		Success:  true,
// 		Message:  message,
// 		RegionId: regionId,
// 		LabelIds: labelIds,
// 	}
// }

// // ToProtoRemoveResponse - Convert remove label result to proto response
// func (m *RegionLabelMapper) ToProtoRemoveResponse(regionId, labelId uint64, message string) *tqdpb.RemoveLabelFromRegionResponse {
// 	return &tqdpb.RemoveLabelFromRegionResponse{
// 		Success:  true,
// 		Message:  message,
// 		RegionId: regionId,
// 		LabelId:  labelId,
// 	}
// }

// // ToProtoGetLabelsResponse - Convert labels list to proto response
// func (m *RegionLabelMapper) ToProtoGetLabelsResponse(labels []qh_domain.QHLabel) *tqdpb.GetLabelsByRegionResponse {
// 	data := make([]*tqdpb.QHLabelResponse, len(labels))
// 	for i, l := range labels {
// 		data[i] = NewLabelMapper().ToProtoQHLabelResponse(&l)
// 	}

// 	return &tqdpb.GetLabelsByRegionResponse{
// 		Data:  data,
// 		Total: int32(len(labels)),
// 	}
// }
