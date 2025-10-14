package handler

import (
	"context"
	pb "pb/types/user"
	"user/infra/mapper"
	"user/internal/usecases"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PlanHandler struct {
	pb.UnimplementedPlanServiceServer
	planUsecase usecases.IPlanUsecase
	mapper      *mapper.PlanMapper
}

func NewPlanHandler(
	planUsecase usecases.IPlanUsecase,
	mapper *mapper.PlanMapper,
) *PlanHandler {
	return &PlanHandler{
		planUsecase: planUsecase,
		mapper:      mapper,
	}
}

// GetVisiblePlans gets all visible plans for users
// @Summary Lấy danh sách các gói membership hiển thị
// @Description Lấy tất cả các gói membership đang hiển thị cho user
// @Tags Plan
// @Accept json
// @Produce json
// @Success 200 {object} pb.GetVisiblePlansResponse
// @Failure 500 {object} map[string]string
// @Router /user/plans [get]
func (h *PlanHandler) GetVisiblePlans(ctx context.Context, req *pb.GetVisiblePlansRequest) (*pb.GetVisiblePlansResponse, error) {
	plans, errDTO := h.planUsecase.GetVisiblePlans(ctx)
	if errDTO != nil {
		return nil, status.Error(codes.Code(errDTO.Code), errDTO.Message)
	}

	return &pb.GetVisiblePlansResponse{
		Plans: h.mapper.ToProtoList(plans),
	}, nil
}

// CreatePlan creates a new plan
// @Summary Tạo gói membership mới
// @Description Tạo một gói membership mới (Admin only)
// @Tags Plan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body pb.CreatePlanRequest true "Thông tin gói"
// @Success 200 {object} pb.PlanResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/plans [post]
func (h *PlanHandler) CreatePlan(ctx context.Context, req *pb.CreatePlanRequest) (*pb.PlanResponse, error) {
	plan := h.mapper.FromCreateRequest(req)

	result, err := h.planUsecase.Create(ctx, plan)
	if err != nil {
		return nil, err
	}

	return h.mapper.ToProto(result), nil
}

// UpdatePlan updates an existing plan
// @Summary Cập nhật gói membership
// @Description Cập nhật thông tin gói membership (Admin only)
// @Tags Plan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body pb.UpdatePlanRequest true "Thông tin gói cập nhật"
// @Success 200 {object} pb.PlanResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/plans/{id} [put]
func (h *PlanHandler) UpdatePlan(ctx context.Context, req *pb.UpdatePlanRequest) (*pb.PlanResponse, error) {
	plan := h.mapper.FromUpdateRequest(req)

	result, err := h.planUsecase.Update(ctx, req.Id, plan)
	if err != nil {
		return nil, err
	}

	return h.mapper.ToProto(result), nil
}

// DeletePlan deletes a plan (soft delete)
// @Summary Xóa gói membership
// @Description Xóa mềm gói membership (Admin only)
// @Tags Plan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Plan ID"
// @Success 200 {object} pb.DeletePlanResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/plans/{id} [delete]
func (h *PlanHandler) DeletePlan(ctx context.Context, req *pb.DeletePlanRequest) (*pb.DeletePlanResponse, error) {
	_, err := h.planUsecase.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DeletePlanResponse{
		Success: true,
	}, nil
}

// GetPlanById gets a plan by ID
// @Summary Lấy thông tin gói membership theo ID
// @Description Lấy chi tiết một gói membership (Admin only)
// @Tags Plan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Plan ID"
// @Success 200 {object} pb.GetPlanByIdResponse
// @Failure 404 {object} map[string]string
// @Router /admin/plans/{id} [get]
func (h *PlanHandler) GetPlanById(ctx context.Context, req *pb.GetPlanByIdRequest) (*pb.GetPlanByIdResponse, error) {
	plan, err := h.planUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetPlanByIdResponse{
		Plan: h.mapper.ToProto(plan),
	}, nil
}

// GetAllPlans gets all plans
// @Summary Lấy tất cả gói membership
// @Description Lấy tất cả các gói membership (bao gồm ẩn) (Admin only)
// @Tags Plan
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} pb.GetAllPlansResponse
// @Failure 500 {object} map[string]string
// @Router /admin/plans [get]
func (h *PlanHandler) GetAllPlans(ctx context.Context, req *pb.GetAllPlansRequest) (*pb.GetAllPlansResponse, error) {
	plans, err := h.planUsecase.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.GetAllPlansResponse{
		Plans: h.mapper.ToProtoList(plans),
	}, nil
}
