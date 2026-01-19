package usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/provider"
	"bdspro/internal/repo"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"fmt"
	"net/url"
	"path"
	"sort"
	"strings"
	"time"
)

type ProductNoteUsecase struct {
	noteRepo  repo.ProductNoteRepository
	txManager provider.TransactionProvider
}

func NewProductNoteUsecase(
	noteRepo repo.ProductNoteRepository,
	tx provider.TransactionProvider,
) *ProductNoteUsecase {
	return &ProductNoteUsecase{
		noteRepo:  noteRepo,
		txManager: tx,
	}
}

func (u *ProductNoteUsecase) CreateNote(
	ctx context.Context,
	req *dto.CreateProductNoteRequest,
) (*domain.ProductNote, error) {
	authorID := _utils.GetProfileIdWithContext(ctx)
	if authorID == 0 {
		return nil, _errors.UnauthorizedException()
	}

	if strings.TrimSpace(req.Content) == "" {
		return nil, _errors.BadRequestException("note content is required")
	}

	if err := validateMentions(req.Content, req.Mentions); err != nil {
		return nil, err
	}

	var createdNote *domain.ProductNote

	err := u.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		note := &domain.ProductNote{
			ProductID: req.ProductID,
			AuthorID:  authorID,
			Content:   req.Content,
			PinnedAt:  nil,
		}

		if err := u.noteRepo.CreateNote(txCtx, note); err != nil {
			return err
		}

		// ===== Mentions =====
		if len(req.Mentions) > 0 {
			mentions := make([]domain.ProductNoteMention, 0, len(req.Mentions))

			for _, m := range req.Mentions {
				mentions = append(mentions, domain.ProductNoteMention{
					NoteID:   note.ID,
					UserID:   m.UserID,
					StartPos: m.StartPos,
					Length:   m.Length,
				})
			}

			if err := u.noteRepo.CreateMention(txCtx, mentions); err != nil {
				return err
			}
		}

		// ===== Files =====
		if len(req.Files) > 0 {
			files := make([]domain.ProductNoteFile, 0, len(req.Files))

			for _, f := range req.Files {
				files = append(files, domain.ProductNoteFile{
					NoteID:   note.ID,
					URL:      f.URL,
					FileName: f.FileName,
					FileType: f.FileType,
					Size:     f.Size,
				})
			}

			if err := u.noteRepo.CreateFiles(txCtx, files); err != nil {
				return err
			}
		}

		createdNote = note
		return nil
	})

	return createdNote, err
}

func extractFileName(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "unknown"
	}
	return path.Base(u.Path)
}

func (u *ProductNoteUsecase) TogglePin(ctx context.Context, noteID uint64, isPinned bool) error {
	note, err := u.noteRepo.GetByID(ctx, noteID)
	if err != nil {
		return err
	}
	if note == nil {
		return _errors.NotFoundException("note not found")
	}

	var pinnedAt *time.Time
	if isPinned {
		now := time.Now()
		pinnedAt = &now
	}

	return u.noteRepo.UpdatePinnedAt(ctx, noteID, pinnedAt)
}

func (u *ProductNoteUsecase) GetNotesByProduct(
	ctx context.Context,
	filter *dto.GetProductNotesFilter,
) ([]*dto.ProductNoteItemResponse, int, error) {
	resp, total, err := u.noteRepo.GetListByProductID(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return resp, int(total), nil
}

func (u *ProductNoteUsecase) DeleteNote(ctx context.Context, noteID uint64, requestUserID uint64) error {
	note, err := u.noteRepo.GetByID(ctx, noteID)
	if err != nil {
		return err
	}
	if note == nil {
		return _errors.NotFoundException("note not found")
	}

	profileId := _utils.GetProfileIdWithContext(ctx)
	if profileId != note.AuthorID {
		return _errors.ForbiddenException("permission denied")
	}

	return u.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := u.noteRepo.DeleteMentionByNoteID(txCtx, noteID); err != nil {
			return err
		}
		if err := u.noteRepo.DeleteFilesByNoteID(txCtx, noteID); err != nil {
			return err
		}
		return u.noteRepo.DeleteNote(txCtx, noteID)
	})
}

func (u *ProductNoteUsecase) UpdateNote(
	ctx context.Context,
	req *dto.UpdateProductNoteRequest,
) error {
	authorID := _utils.GetProfileIdWithContext(ctx)
	if authorID == 0 {
		return _errors.UnauthorizedException()
	}

	if err := validateMentions(req.Content, req.Mentions); err != nil {
		return err
	}

	note, err := u.noteRepo.GetByID(ctx, req.NoteID)
	if err != nil {
		return err
	}
	if note == nil {
		return _errors.NotFoundException("note not found")
	}
	if note.AuthorID != authorID {
		return _errors.ForbiddenException("permission denied")
	}

	return u.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := u.noteRepo.UpdateContent(txCtx, req.NoteID, req.Content); err != nil {
			return err
		}

		if err := u.noteRepo.DeleteMentionByNoteID(txCtx, req.NoteID); err != nil {
			return err
		}

		if len(req.Mentions) > 0 {
			mentions := make([]domain.ProductNoteMention, 0, len(req.Mentions))

			for _, m := range req.Mentions {
				mentions = append(mentions, domain.ProductNoteMention{
					NoteID:   req.NoteID,
					UserID:   m.UserID,
					StartPos: m.StartPos,
					Length:   m.Length,
				})
			}

			if err := u.noteRepo.CreateMention(txCtx, mentions); err != nil {
				return err
			}
		}

		// ===== Files =====
		if err := u.noteRepo.DeleteFilesByNoteID(txCtx, req.NoteID); err != nil {
			return err
		}

		if len(req.Files) > 0 {
			files := make([]domain.ProductNoteFile, 0, len(req.Files))

			for _, f := range req.Files {
				files = append(files, domain.ProductNoteFile{
					NoteID:   req.NoteID,
					URL:      f.URL,
					FileName: f.FileName,
					FileType: f.FileType,
					Size:     f.Size,
				})
			}

			if err := u.noteRepo.CreateFiles(txCtx, files); err != nil {
				return err
			}
		}

		return nil
	})
}

func validateMentions(content string, mentions []dto.MentionNoteDTO) error {
	contentRunes := []rune(content)
	contentLen := int32(len(contentRunes))

	if len(mentions) == 0 {
		return nil
	}
	// mentions := []Mention{
	// 	{StartPos: 10},
	// 	{StartPos: 3},
	// 	{StartPos: 7},
	// }
	sort.Slice(mentions, func(i, j int) bool {
		return mentions[i].StartPos < mentions[j].StartPos
	})
	// mentions := []Mention{
	// 	{StartPos: 3},
	// 	{StartPos: 7},
	// 	{StartPos: 10},
	// }

	var prevEnd int32 = 0
	for _, m := range mentions {
		if m.StartPos < 0 || m.Length <= 0 {
			return _errors.BadRequestException("invalid mention position")
		}

		end := m.StartPos + m.Length
		if end > contentLen {
			return _errors.BadRequestException("mention out of content range")
		}

		if m.StartPos < prevEnd {
			return _errors.BadRequestException(
				fmt.Sprintf(
					"mention overlap: userId=%d start=%d length=%d",
					m.UserID,
					m.StartPos,
					m.Length,
				),
			)
		}

		// text := string(contentRunes[m.StartPos:end])
		// if !strings.HasPrefix(text, "@") {
		// 	return _errors.BadRequestException("mention text must start with @")
		// }

		prevEnd = end
	}

	return nil
}
