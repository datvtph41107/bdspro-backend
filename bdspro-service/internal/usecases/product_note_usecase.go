package usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/provider"
	"bdspro/internal/repo"
	_errors "common/errors"
	_utils "common/utils"
	"context"
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

func (u *ProductNoteUsecase) CreateNote(ctx context.Context, req *dto.CreateProductNoteRequest) (*domain.ProductNote, error) {
	authorID := _utils.GetProfileIdWithContext(ctx)
	if authorID == 0 {
		return nil, _errors.UnauthorizedException()
	}
	if req.Content == "" {
		return nil, _errors.BadRequestException("note content is required")
	}

	var createdNote *domain.ProductNote
	err := u.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
		note := &domain.ProductNote{
			ProductID: req.ProductID,
			AuthorID:  authorID,
			Content:   req.Content,
			IsPinned:  false,
		}
		if err := u.noteRepo.CreateNote(txCtx, note); err != nil {
			return err
		}

		mentionIds := req.MentionIDs
		if len(mentionIds) > 0 {
			mentions := make([]domain.ProductNoteMention, 0, len(mentionIds))
			for _, uid := range mentionIds {
				mentions = append(mentions, domain.ProductNoteMention{
					NoteID: note.ID,
					UserID: uid,
				})
			}
			if err := u.noteRepo.CreateMention(txCtx, mentions); err != nil {
				return err
			}
		}
		if len(req.Files) > 0 {
			if err := u.noteRepo.AttachFilesToNote(txCtx, note.ID, req.Files); err != nil {
				return err
			}
		}
		createdNote = note
		return nil
	})

	return createdNote, err
}

func (u *ProductNoteUsecase) TogglePin(ctx context.Context, noteID uint64, isPinned bool) error {
	note, err := u.noteRepo.GetByID(ctx, noteID)
	if err != nil {
		return err
	}
	if note == nil {
		return _errors.NotFoundException("note not found")
	}

	return u.noteRepo.UpdatePin(ctx, noteID, isPinned)
}

func (u *ProductNoteUsecase) GetNotesByProduct(
	ctx context.Context,
	filter *dto.GetProductNotesFilter,
) ([]*dto.ProductNoteItemRespone, int, error) {
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

	if note.AuthorID != requestUserID {
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

func (u *ProductNoteUsecase) UpdateNote(ctx context.Context, req *dto.UpdateProductNoteRequest) error {
	authorID := _utils.GetProfileIdWithContext(ctx)
	if authorID == 0 {
		return _errors.UnauthorizedException()
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
		if len(req.MentionIDs) > 0 {
			if err := u.noteRepo.DeleteMentionByNoteID(txCtx, req.NoteID); err != nil {
				return err
			}
			mentions := make([]domain.ProductNoteMention, 0, len(req.MentionIDs))
			for _, uid := range req.MentionIDs {
				mentions = append(mentions, domain.ProductNoteMention{
					NoteID: req.NoteID,
					UserID: uid,
				})
			}
			if err := u.noteRepo.CreateMention(txCtx, mentions); err != nil {
				return err
			}
		}

		if len(req.Files) > 0 {
			if err := u.noteRepo.DeleteFilesByNoteID(txCtx, req.NoteID); err != nil {
				return err
			}

			if err := u.noteRepo.AttachFilesToNote(txCtx, req.NoteID, req.Files); err != nil {
				return err
			}
		}

		return nil
	})
}
