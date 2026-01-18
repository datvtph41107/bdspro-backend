package handler

import (
	"bdspro/infra/client"
	"bdspro/infra/mapper"
	"bdspro/internal/dto"
	"bdspro/internal/usecases"
	_dto "common/domain/dto"
	"context"

	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductNoteHandler struct {
	bdspropb.UnimplementedProductNoteServiceServer
	ProductNoteUsecase *usecases.ProductNoteUsecase
	ProductNoteMapper  *mapper.ProductNoteMapper
	UserClient         *client.UserClient
}

func NewProductNoteHandler(
	usecase *usecases.ProductNoteUsecase,
	mapper *mapper.ProductNoteMapper,
	userClient *client.UserClient,
) *ProductNoteHandler {
	return &ProductNoteHandler{
		ProductNoteUsecase: usecase,
		ProductNoteMapper:  mapper,
		UserClient:         userClient,
	}
}

// LIST
func (h *ProductNoteHandler) ListProductNotes(
	ctx context.Context,
	req *bdspropb.ListProductNotesRequest,
) (*bdspropb.ListProductNotesResponse, error) {
	filter := &dto.GetProductNotesFilter{
		ProductID:  req.ProductId,
		OnlyPinned: req.OnlyPinned,
		AuthorID:   req.AuthorId,
		Query:      &req.Query,
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
	}

	items, total, err := h.ProductNoteUsecase.GetNotesByProduct(ctx, filter)
	if err != nil {
		return nil, err
	}
	if items == nil {
		return nil, nil
	}

	h.enrichUsers(ctx, items)

	return &bdspropb.ListProductNotesResponse{
		Data:  h.ProductNoteMapper.ProductNoteListToProto(items),
		Total: int64(total),
	}, nil
}

// CREATE
func (h *ProductNoteHandler) CreateProductNote(
	ctx context.Context,
	req *bdspropb.CreateProductNoteRequest,
) (*bdspropb.Response, error) {
	_, err := h.ProductNoteUsecase.CreateNote(ctx, &dto.CreateProductNoteRequest{
		ProductID:  req.ProductId,
		Content:    req.Content,
		MentionIDs: req.MentionIds,
		Files:      req.Files,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &bdspropb.Response{
		Id:      req.ProductId,
		Message: "success",
	}, nil
}

// UPDATE
func (h *ProductNoteHandler) UpdateProductNote(
	ctx context.Context,
	req *bdspropb.UpdateProductNoteReqClient,
) (*bdspropb.Response, error) {
	// files := make([]dto.ProductNoteFileItem, 0, len(req.Files))
	// for _, f := range req.Files {
	// 	files = append(files, dto.ProductNoteFileItem{
	// 		ID:       f.Id,
	// 		FileName: f.FileName,
	// 		URL:      f.FileUrl,
	// 		Type:     f.FileType,
	// 		Size:     f.FileSize,
	// 	})
	// }

	err := h.ProductNoteUsecase.UpdateNote(ctx, &dto.UpdateProductNoteRequest{
		NoteID:     req.Id,
		Content:    req.Content,
		MentionIDs: req.MentionIds,
		Files:      req.Files,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &bdspropb.Response{
		Id:      req.Id,
		Message: "success",
	}, nil
}

func (h *ProductNoteHandler) DeleteProductNote(
	ctx context.Context,
	req *bdspropb.DeleteProductNoteRequest,
) (*bdspropb.Response, error) {

	if err := h.ProductNoteUsecase.DeleteNote(ctx, req.NoteId, req.RequestUserID); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &bdspropb.Response{
		Id:      req.NoteId,
		Message: "success",
	}, nil
}

func (h *ProductNoteHandler) TogglePinProductNote(
	ctx context.Context,
	req *bdspropb.TogglePinProductNoteRequest,
) (*bdspropb.Response, error) {

	if err := h.ProductNoteUsecase.TogglePin(ctx, req.Id, req.IsPinned); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &bdspropb.Response{
		Id:      req.Id,
		Message: "success",
	}, nil
}

func (h *ProductNoteHandler) enrichUsers(
	ctx context.Context,
	notes []*dto.ProductNoteItemRespone,
) {
	userIDs := make(map[uint64]struct{})

	for _, n := range notes {
		if n.AuthorId != 0 {
			userIDs[n.AuthorId] = struct{}{}
		}
		for _, id := range n.MentionIDs {
			userIDs[uint64(id)] = struct{}{}
		}
	}

	if len(userIDs) == 0 {
		return
	}

	profiles, err := h.UserClient.GetMapProfileByIds(
		ctx,
		&sharepb.GetProfileByIdsRequest{Ids: mapKeys(userIDs)},
	)
	if err != nil {
		return
	}

	for _, n := range notes {
		if p, ok := profiles[n.AuthorId]; ok {
			n.Author = p
		}

		for _, id := range n.MentionIDs {
			uid := uint64(id)
			if p, ok := profiles[uid]; ok {
				n.Mentions = append(n.Mentions, dto.ProductNoteMentionItem{
					UserId: uid,
					Name:   p.FullName,
					Avatar: p.Avatar,
					Role:   "member",
				})
			}
		}
	}
}

func (h *ProductNoteHandler) getAuthorProfile(ctx context.Context, userID uint64) *sharepb.ProfileItem {
	if userID == 0 {
		return nil
	}

	resp, err := h.UserClient.GetProfileByIds(
		ctx,
		&sharepb.GetProfileByIdsRequest{Ids: []uint64{userID}},
	)
	if err != nil || len(resp.Profiles) == 0 {
		return nil
	}

	return resp.Profiles[0]
}

func mapKeys(m map[uint64]struct{}) []uint64 {
	res := make([]uint64, 0, len(m))
	for k := range m {
		res = append(res, k)
	}
	return res
}
