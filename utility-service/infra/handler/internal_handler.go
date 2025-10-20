package handler

import (
	"context"
	sharepb "pb/types/shared"
	utilitypb "pb/types/utility"
)

type UtilityInternalService struct {
	utilitypb.UnimplementedUtilityInternalServiceServer
}

func NewUtilityInternalService() *UtilityInternalService {
	return &UtilityInternalService{}
}

// @Summary Health check
// @Description Kiểm tra trạng thái service
// @Tags Internal
// @Accept json
// @Produce json
// @Router /internal/health [get]
func (s *UtilityInternalService) HealthCheck(ctx context.Context, req *sharepb.Empty) (*sharepb.SubmitResponse, error) {
	return &sharepb.SubmitResponse{
		Message: "Utility service is running",
	}, nil
}
