package handler

import (
	"context"

	_dto "common/domain/dto"
	"hub/infra/mapper"
	_usecase "hub/internal/usecase"
	hubpb "pb/types/hub"
)

type FAQHandler struct {
	hubpb.UnimplementedFAQServiceServer
	faqUsecase _usecase.IFAQUsecase
	faqMapper  *mapper.FAQMapper
}

func NewFAQHandler(
	faqUsecase _usecase.IFAQUsecase,
	faqMapper *mapper.FAQMapper,
) *FAQHandler {
	return &FAQHandler{
		faqUsecase: faqUsecase,
		faqMapper:  faqMapper,
	}
}

// @Summary Tạo FAQ mới
// @Description Tạo một FAQ mới với question, answer và groupKey
// @Tags FAQ
// @Accept json
// @Produce json
// @Param body body hubpb.CreateFAQRequest true "Thông tin FAQ"
// @Success 200 {object} hubpb.CreateFAQResponse
// @Router /v2/hub/faqs [post]
func (h *FAQHandler) CreateFAQ(ctx context.Context, req *hubpb.CreateFAQRequest) (*hubpb.CreateFAQResponse, error) {
	entity := h.faqMapper.CreateRequestToEntity(req)

	result, err := h.faqUsecase.Create(ctx, entity)
	if err != nil {
		return nil, err
	}

	return &hubpb.CreateFAQResponse{
		Data: h.faqMapper.EntityToDetailProto(result),
	}, nil
}

// @Summary Cập nhật FAQ
// @Description Cập nhật thông tin FAQ theo ID
// @Tags FAQ
// @Accept json
// @Produce json
// @Param id path uint64 true "FAQ ID"
// @Param body body hubpb.UpdateFAQRequest true "Thông tin FAQ cần cập nhật"
// @Success 200 {object} hubpb.UpdateFAQResponse
// @Router /v2/hub/faqs/{id} [put]
func (h *FAQHandler) UpdateFAQ(ctx context.Context, req *hubpb.UpdateFAQRequest) (*hubpb.UpdateFAQResponse, error) {
	// Get existing entity
	existing, err := h.faqUsecase.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	// Update entity with request data
	h.faqMapper.UpdateRequestToEntity(existing, req)

	// Save updated entity
	result, err := h.faqUsecase.Update(ctx, req.GetId(), existing)
	if err != nil {
		return nil, err
	}

	return &hubpb.UpdateFAQResponse{
		Data: h.faqMapper.EntityToDetailProto(result),
	}, nil
}

// @Summary Xóa FAQ
// @Description Xóa FAQ theo ID (soft delete)
// @Tags FAQ
// @Accept json
// @Produce json
// @Param id path uint64 true "FAQ ID"
// @Success 200 {object} hubpb.DeleteFAQResponse
// @Router /v2/hub/faqs/{id} [delete]
func (h *FAQHandler) DeleteFAQ(ctx context.Context, req *hubpb.DeleteFAQRequest) (*hubpb.DeleteFAQResponse, error) {
	_, err := h.faqUsecase.Delete(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &hubpb.DeleteFAQResponse{
		Success: true,
	}, nil
}

// @Summary Lấy chi tiết FAQ
// @Description Lấy thông tin chi tiết FAQ theo ID
// @Tags FAQ
// @Accept json
// @Produce json
// @Param id path uint64 true "FAQ ID"
// @Success 200 {object} hubpb.GetFAQResponse
// @Router /v2/hub/faqs/{id} [get]
func (h *FAQHandler) GetFAQ(ctx context.Context, req *hubpb.GetFAQRequest) (*hubpb.GetFAQResponse, error) {
	result, err := h.faqUsecase.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &hubpb.GetFAQResponse{
		Data: h.faqMapper.EntityToDetailProto(result),
	}, nil
}

// @Summary Lấy danh sách FAQs (Admin)
// @Description Lấy danh sách FAQs với tìm kiếm theo question, filter theo groupKey và phân trang
// @Tags FAQ/Admin
// @Accept json
// @Produce json
// @Param question query string false "Tìm kiếm theo question"
// @Param groupKey query string false "Lọc theo groupKey"
// @Param page query int false "Số trang (default: 1)"
// @Param size query int false "Số bản ghi trên trang (default: 20, max: 100)"
// @Success 200 {object} hubpb.GetFAQListResponse
// @Router /v2/hub/admin/faqs [get]
func (h *FAQHandler) GetFAQList(ctx context.Context, req *hubpb.GetFAQListRequest) (*hubpb.GetFAQListResponse, error) {
	// Create pagable from request
	pagable := &_dto.Pagable{
		Page: uint32(req.GetPage()),
		Size: uint32(req.GetSize()),
	}

	// Get filtered list
	data, total, err := h.faqUsecase.GetListWithFilter(
		ctx,
		req.GetQuestion(),
		req.GetGroupKey(),
		pagable,
	)
	if err != nil {
		return nil, err
	}

	return &hubpb.GetFAQListResponse{
		Data:  h.faqMapper.EntitiesToProto(data),
		Total: total,
	}, nil
}

// @Summary Lấy danh sách FAQs đơn giản cho user
// @Description API public trả về nội dung FAQ đầy đủ để ứng dụng và SEO render cùng một câu trả lời. Hỗ trợ search theo text và lọc groupKey.
// @Tags FAQ/Simple
// @Accept json
// @Produce json
// @Param text query string false "Search theo question hoặc answer"
// @Param groupKey query string false "Lọc theo groupKey"
// @Param page query int false "Số trang (default: 1)"
// @Param size query int false "Số bản ghi trên trang (default: 20, max: 100)"
// @Success 200 {object} hubpb.GetFAQSimpleResponse
// @Router /v2/hub/faqs [get]
func (h *FAQHandler) GetFAQSimple(ctx context.Context, req *hubpb.GetFAQSimpleRequest) (*hubpb.GetFAQSimpleResponse, error) {
	// Create pagable from request
	pagable := &_dto.Pagable{
		Page: uint32(req.GetPage()),
		Size: uint32(req.GetSize()),
	}

	// Get simple list with text search
	data, total, err := h.faqUsecase.GetSimpleList(
		ctx,
		req.Text,
		req.GetGroupKey(),
		pagable,
	)
	if err != nil {
		return nil, err
	}

	// Public consumers (including Next.js SEO rendering) need the complete answer.
	// Truncation belongs to presentation components, not to the source-of-truth API.
	var simpleItems []*hubpb.FAQSimple
	for _, item := range data {
		simpleItems = append(simpleItems, &hubpb.FAQSimple{
			Id:       item.ID,
			Question: h.faqMapper.SanitizeUTF8(item.Question),
			Answer:   h.faqMapper.SanitizeUTF8(item.Answer),
		})
	}

	return &hubpb.GetFAQSimpleResponse{
		Data:  simpleItems,
		Total: total,
	}, nil
}

// @Summary Lấy chi tiết FAQ đơn giản cho user
// @Description API public lấy chi tiết FAQ đầy đủ theo ID
// @Tags FAQ/Simple
// @Accept json
// @Produce json
// @Param id path uint64 true "FAQ ID"
// @Success 200 {object} hubpb.FAQDetail
// @Router /v2/hub/faqs/{id} [get]
func (h *FAQHandler) GetFAQSimpleDetail(ctx context.Context, req *hubpb.GetFAQRequest) (*hubpb.FAQDetail, error) {
	result, err := h.faqUsecase.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	data := h.faqMapper.EntityToDetailProto(result)
	return data, nil
}
