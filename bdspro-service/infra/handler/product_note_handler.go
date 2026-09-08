package handler

import (
	"bdspro/infra/client"
	"bdspro/infra/mapper"
	"bdspro/internal/dto"
	"bdspro/internal/usecases"
	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"
	"context"

	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
)

type ProductNoteHandler struct {
	bdspropb.UnimplementedProductNoteServiceServer
	ProductNoteUsecase *usecases.ProductNoteUsecase
	ProductNoteMapper  *mapper.ProductNoteMapper
	UserClient         *client.UserClient
	SyncProvider       *_utils.SyncUtil
}

func NewProductNoteHandler(
	usecase *usecases.ProductNoteUsecase,
	mapper *mapper.ProductNoteMapper,
	userClient *client.UserClient,
	SyncProvider *_utils.SyncUtil,
) *ProductNoteHandler {
	return &ProductNoteHandler{
		ProductNoteUsecase: usecase,
		ProductNoteMapper:  mapper,
		UserClient:         userClient,
		SyncProvider:       SyncProvider,
	}
}

func (h *ProductNoteHandler) ListProductNotes(
	ctx context.Context,
	req *bdspropb.ListProductNotesRequest,
) (*bdspropb.ListProductNotesResponse, error) {
	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyProductNotes, req.ProductId)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	// client mặc định sẽ truyền page = 0 GET /v2/deals?page=0
	if !updated && req.Page == 0 {
		return &bdspropb.ListProductNotesResponse{}, nil
	}
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
	// if req.Page == 0 && len(items) > 0 {
	// 	time := time.Now().UnixMilli() - 24*60*60*1000
	// 	if len(items) > 0 {
	// 		time = items[0].Note.ProductNoteUpdatedAt.UnixMilli()
	// 	}
	// 	h.SyncProvider.PutTimestamp(ctx, key, time)
	// 	// h.SyncProvider.PutTimestamp(ctx, fmt.Sprintf(key, req.ProductId, userId), items[0].Note.ProductNoteUpdatedAt.UnixMilli())
	// }
	h.enrichUsers(ctx, items)
	return &bdspropb.ListProductNotesResponse{
		Data:  h.ProductNoteMapper.ProductNoteListToProto(items),
		Total: int64(total),
	}, nil
}

func (h *ProductNoteHandler) CreateProductNote(
	ctx context.Context,
	req *bdspropb.CreateProductNoteRequest,
) (*bdspropb.Response, error) {
	if req.ProductId == 0 {
		return nil, _errors.BadRequestException("productId is required")
	}

	result, err := h.ProductNoteUsecase.CreateNote(ctx, &dto.CreateProductNoteRequest{
		ProductID: req.ProductId,
		Content:   req.Content,
		Mentions:  mapProtoMentionsToDTO(req.Mentions),
		Files:     mapProtoFilesToDTO(req.Files),
	})
	if err != nil {
		return nil, err
	}
	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyProductNotes, req.ProductId)
	h.SyncProvider.PutTimestamp(ctx, key, result.UpdatedAt.UnixMilli())
	return &bdspropb.Response{
		Id:      req.ProductId,
		Message: "success",
	}, nil
}

func (h *ProductNoteHandler) UpdateProductNote(
	ctx context.Context,
	req *bdspropb.UpdateProductNoteReqClient,
) (*bdspropb.Response, error) {
	if req.Id == 0 {
		return nil, _errors.BadRequestException("noteId is required")
	}

	err := h.ProductNoteUsecase.UpdateNote(ctx, &dto.UpdateProductNoteRequest{
		NoteID:   req.Id,
		Content:  req.Content,
		Mentions: mapProtoMentionsToDTO(req.Mentions),
		Files:    mapProtoFilesToDTO(req.Files),
	})
	if err != nil {
		return nil, err
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
	if req.NoteId == 0 {
		return nil, _errors.BadRequestException("noteId is required")
	}

	if err := h.ProductNoteUsecase.DeleteNote(ctx, req.NoteId); err != nil {
		return nil, err
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
	if req.Id == 0 {
		return nil, _errors.BadRequestException("noteId is required")
	}

	if err := h.ProductNoteUsecase.TogglePin(ctx, req.Id, req.IsPinned); err != nil {
		return nil, err
	}

	return &bdspropb.Response{
		Id:      req.Id,
		Message: "success",
	}, nil
}

func (h *ProductNoteHandler) enrichUsers(
	ctx context.Context,
	notes []*dto.ProductNote,
) {
	userIDs := make(map[uint64]struct{})

	for _, n := range notes {
		if n.Note.AuthorId != 0 {
			userIDs[n.Note.AuthorId] = struct{}{}
		}
		for _, m := range n.Mentions {
			userIDs[m.UserId] = struct{}{}
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
		if p, ok := profiles[n.Note.AuthorId]; ok {
			n.Note.Author = p
		}

		for _, m := range n.Mentions {
			if p, ok := profiles[m.UserId]; ok {
				m.Name = p.FullName
				m.Avatar = p.Avatar
				m.Role = "member"
			}
		}
	}
}

func mapProtoMentionsToDTO(
	list []*bdspropb.ProductNoteMention,
) []dto.MentionNoteDTO {

	if len(list) == 0 {
		return nil
	}

	res := make([]dto.MentionNoteDTO, 0, len(list))
	for _, m := range list {
		res = append(res, dto.MentionNoteDTO{
			UserID:   m.UserId,
			StartPos: m.StartPos,
			Length:   m.Length,
		})
	}
	return res
}

func mapProtoFilesToDTO(
	list []*bdspropb.ProductNoteFileInput,
) []dto.ProductNoteFileDTO {
	if len(list) == 0 {
		return nil
	}

	res := make([]dto.ProductNoteFileDTO, 0, len(list))
	for _, f := range list {
		res = append(res, dto.ProductNoteFileDTO{
			URL:      f.Url,
			FileName: f.FileName,
			FileType: f.FileType,
			Size:     f.Size,
		})
	}
	return res
}

func mapKeys(m map[uint64]struct{}) []uint64 {
	res := make([]uint64, 0, len(m))
	for k := range m {
		res = append(res, k)
	}
	return res
}
