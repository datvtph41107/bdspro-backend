package services

import (
	"context"
	"strconv"

	_utils "common/utils"
	"mdn4/internal/interface/providers"
	"mdn4/modules/bdspro/dto"
	"mdn4/modules/bdspro/repository"
	bdspropb "pb/types/bdspro"
)

type ProductNoteService interface {
	ListByProduct(ctx context.Context, productID uint64, req dto.SyncByIDRequest) (*bdspropb.ListProductNotesResponse, error)
}

// @bind: mdn4/modules/bdspro/services.ProductNoteService
type LocalProductNoteService struct {
	sync    *_utils.SyncUtil
	repo    repository.ProductNoteRepository
	profile providers.ProfileProvider
}

func NewLocalProductNoteService(
	sync *_utils.SyncUtil,
	repo repository.ProductNoteRepository,
	profile providers.ProfileProvider,
) *LocalProductNoteService {
	return &LocalProductNoteService{sync: sync, repo: repo, profile: profile}
}

func (s *LocalProductNoteService) ListByProduct(ctx context.Context, productID uint64, req dto.SyncByIDRequest) (*bdspropb.ListProductNotesResponse, error) {
	if productID == 0 {
		return &bdspropb.ListProductNotesResponse{}, nil
	}
	if s.sync != nil {
		key := s.sync.GetKey(ctx, _utils.SyncKeyProductNotes, productID)
		if !s.sync.HasUpdated(ctx, key, req.Timestamp) && req.Page == 0 && req.Timestamp > 0 {
			return &bdspropb.ListProductNotesResponse{}, nil
		}
	}

	rows, total, err := s.repo.ListByProductID(ctx, productID, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return &bdspropb.ListProductNotesResponse{Total: total}, nil
	}

	noteIDs := make([]uint64, 0, len(rows))
	profileIDs := make([]uint64, 0)
	seen := map[uint64]struct{}{}
	for _, r := range rows {
		noteIDs = append(noteIDs, r.ID)
		if r.AuthorID > 0 {
			if _, ok := seen[r.AuthorID]; !ok {
				seen[r.AuthorID] = struct{}{}
				profileIDs = append(profileIDs, r.AuthorID)
			}
		}
	}

	mentions, _ := s.repo.ListMentionsByNoteIDs(ctx, noteIDs)
	files, _ := s.repo.ListFilesByNoteIDs(ctx, noteIDs)
	for _, m := range mentions {
		if m.UserID > 0 {
			if _, ok := seen[m.UserID]; !ok {
				seen[m.UserID] = struct{}{}
				profileIDs = append(profileIDs, m.UserID)
			}
		}
	}

	profiles := map[uint64]*bdspropb.ProductNoteMember{}
	if s.profile != nil && len(profileIDs) > 0 {
		if m, err := s.profile.GetMapProfileByIds(ctx, profileIDs); err == nil {
			for id, p := range m {
				if p == nil {
					continue
				}
				profiles[id] = &bdspropb.ProductNoteMember{
					Id: id, Name: p.FullName, Avatar: p.Avatar, Role: "member",
				}
			}
		}
	}

	mentionsByNote := map[uint64][]*bdspropb.ProductNoteMention{}
	for _, m := range mentions {
		item := &bdspropb.ProductNoteMention{
			UserId: m.UserID, StartPos: m.StartPos, Length: m.Length, Role: "member",
		}
		if p := profiles[m.UserID]; p != nil {
			item.Name, item.Avatar = p.Name, p.Avatar
		}
		mentionsByNote[m.NoteID] = append(mentionsByNote[m.NoteID], item)
	}
	filesByNote := map[uint64][]*bdspropb.ProductNoteFile{}
	for _, f := range files {
		filesByNote[f.NoteID] = append(filesByNote[f.NoteID], &bdspropb.ProductNoteFile{
			Url: f.URL, FileName: f.FileName, FileType: f.FileType, Size: f.Size,
		})
	}

	data := make([]*bdspropb.ProductNoteEntity, 0, len(rows))
	for _, r := range rows {
		item := &bdspropb.ProductNoteEntity{
			Id: r.ID, ProductId: r.ProductID, Content: r.Content, IsPinned: r.IsPinned,
			Author: profiles[r.AuthorID], Mentions: mentionsByNote[r.ID], Files: filesByNote[r.ID],
		}
		if r.CreatedAt != nil {
			item.CreatedAt = strconv.FormatInt(r.CreatedAt.UnixMilli(), 10)
		}
		data = append(data, item)
	}
	return &bdspropb.ListProductNotesResponse{Data: data, Total: total}, nil
}
