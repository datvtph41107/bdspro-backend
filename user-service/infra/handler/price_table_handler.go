package handler

import (
	_dto "common/domain/dto"
	_err "common/domain/err"
	"context"
	pb "pb/types/user"
	"user/infra/mapper"
	"user/internal/usecases"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PriceTableHandler struct {
	pb.UnimplementedPriceTableServiceServer
	priceTableUsecase usecases.IPriceTableUsecase
	mapper            *mapper.PriceTableMapper
}

// @bind: internal/usecases.IPriceTableUsecase
// @bind: infra/mapper.PriceTableMapper
func NewPriceTableHandler(
	priceTableUsecase usecases.IPriceTableUsecase,
	mapper *mapper.PriceTableMapper,
) *PriceTableHandler {
	return &PriceTableHandler{
		priceTableUsecase: priceTableUsecase,
		mapper:            mapper,
	}
}

// CreatePriceTable creates a new price table
// @Summary Tạo bảng giá mới
// @Description Tạo một bảng giá mới (Admin only)
// @Tags PriceTable
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body pb.CreatePriceTableRequest true "Thông tin bảng giá"
// @Success 200 {object} pb.PriceTableResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/price-tables [post]
func (h *PriceTableHandler) CreatePriceTable(ctx context.Context, req *pb.CreatePriceTableRequest) (*pb.PriceTableResponse, error) {
	entity := h.mapper.FromCreateRequest(req)
	result, err := h.priceTableUsecase.Create(ctx, entity)
	if err != nil {
		if errDTO, ok := err.(*_err.ErrorDTO); ok {
			return nil, status.Error(codes.Code(errDTO.Code), errDTO.Message)
		}
		return nil, err
	}
	return h.mapper.ToProto(result), nil
}

// UpdatePriceTable updates an existing price table
// @Summary Cập nhật bảng giá
// @Description Cập nhật thông tin bảng giá (Admin only)
// @Tags PriceTable
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Price Table ID"
// @Param request body pb.UpdatePriceTableRequest true "Thông tin bảng giá cập nhật"
// @Success 200 {object} pb.PriceTableResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/price-tables/{id} [put]
func (h *PriceTableHandler) UpdatePriceTable(ctx context.Context, req *pb.UpdatePriceTableRequest) (*pb.PriceTableResponse, error) {
	entity := h.mapper.FromUpdateRequest(req)
	result, err := h.priceTableUsecase.Update(ctx, req.Id, entity)
	if err != nil {
		if errDTO, ok := err.(*_err.ErrorDTO); ok {
			return nil, status.Error(codes.Code(errDTO.Code), errDTO.Message)
		}
		return nil, err
	}
	return h.mapper.ToProto(result), nil
}

// DeletePriceTable deletes a price table
// @Summary Xóa bảng giá
// @Description Xóa mềm bảng giá (Admin only)
// @Tags PriceTable
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Price Table ID"
// @Success 200 {object} pb.DeletePriceTableResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/price-tables/{id} [delete]
func (h *PriceTableHandler) DeletePriceTable(ctx context.Context, req *pb.DeletePriceTableRequest) (*pb.DeletePriceTableResponse, error) {
	_, err := h.priceTableUsecase.Delete(ctx, req.Id)
	if err != nil {
		if errDTO, ok := err.(*_err.ErrorDTO); ok {
			return nil, status.Error(codes.Code(errDTO.Code), errDTO.Message)
		}
		return nil, err
	}
	return &pb.DeletePriceTableResponse{Success: true}, nil
}

// GetPriceTableById gets a price table by ID
// @Summary Lấy thông tin bảng giá theo ID
// @Description Lấy chi tiết một bảng giá (Admin only)
// @Tags PriceTable
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Price Table ID"
// @Success 200 {object} pb.GetPriceTableByIdResponse
// @Failure 404 {object} map[string]string
// @Router /admin/price-tables/{id} [get]
func (h *PriceTableHandler) GetPriceTableById(ctx context.Context, req *pb.GetPriceTableByIdRequest) (*pb.GetPriceTableByIdResponse, error) {
	result, err := h.priceTableUsecase.GetByID(ctx, req.Id)
	if err != nil {
		if errDTO, ok := err.(*_err.ErrorDTO); ok {
			return nil, status.Error(codes.Code(errDTO.Code), errDTO.Message)
		}
		return nil, err
	}
	return &pb.GetPriceTableByIdResponse{
		PriceTable: h.mapper.ToProto(result),
	}, nil
}

// GetPriceTableList gets a list of price tables
// @Summary Lấy danh sách bảng giá
// @Description Lấy danh sách bảng giá với phân trang (Admin only)
// @Tags PriceTable
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Trang hiện tại"
// @Param size query int false "Số lượng item trên mỗi trang"
// @Param search query string false "Tìm kiếm theo tên"
// @Success 200 {object} pb.GetPriceTableListResponse
// @Failure 500 {object} map[string]string
// @Router /admin/price-tables [get]
func (h *PriceTableHandler) GetPriceTableList(ctx context.Context, req *pb.GetPriceTableListRequest) (*pb.GetPriceTableListResponse, error) {
	pagable := &_dto.Pagable{
		Page: uint32(req.Page),
		Size: uint32(req.Size),
	}

	entities, total, err := h.priceTableUsecase.GetList(ctx, pagable)
	if err != nil {
		if errDTO, ok := err.(*_err.ErrorDTO); ok {
			return nil, status.Error(codes.Code(errDTO.Code), errDTO.Message)
		}
		return nil, err
	}

	return &pb.GetPriceTableListResponse{
		Data:  h.mapper.ToProtoList(entities),
		Total: total,
	}, nil
}
