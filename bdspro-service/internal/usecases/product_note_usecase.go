package usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/provider"
	"bdspro/internal/repo"
	_enum "common/domain/enum"
	_errors "common/errors"
	"common/logging"
	_utils "common/utils"
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"path"
	"sort"
	"strings"
	"time"
)

type ProductNoteUsecase struct {
	noteRepo             repo.ProductNoteRepository
	txManager            provider.TransactionProvider
	NotifyProvider       provider.NotificationProvider
	PropertyRelationRepo repo.PropertyRelationRepo
	ProductRepo          repo.ProductRepo
	// NotificationClient   provider.NotificationProvider
}

func NewProductNoteUsecase(
	noteRepo repo.ProductNoteRepository,
	tx provider.TransactionProvider,
	notifyProvider provider.NotificationProvider,
	propertyRelationRepo repo.PropertyRelationRepo,
	productRepo repo.ProductRepo,
	// notificationClient provider.NotificationProvider,
) *ProductNoteUsecase {
	return &ProductNoteUsecase{
		noteRepo:             noteRepo,
		txManager:            tx,
		NotifyProvider:       notifyProvider,
		PropertyRelationRepo: propertyRelationRepo,
		ProductRepo:          productRepo,
		// NotificationClient:   notificationClient,
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

		// Lấy PropertyID từ ProductID
		var propertyID uint64
		if u.PropertyRelationRepo != nil {
			product, err := u.ProductRepo.GetProductByID(ctx, &req.ProductID)
			if err == nil && product != nil && product.PropertyID != nil {
				propertyID = *product.PropertyID
			}
		}

		// Nếu không có PropertyID, sử dụng ProductID làm fallback
		subjectID := propertyID
		if subjectID == 0 {
			subjectID = req.ProductID
		}

		logger := logging.FromContext(ctx)

		go func() {
			contextTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			err := u.NotifyProvider.RegistedEventProperty(
				contextTimeout,
				&dto.CreatePropertyActivityDTO{
					SubjectID:   subjectID,
					Action:      _enum.ACTION_ADD_PRODUCT_TODEAL,
					ActorID:     authorID,
					Description: "Đã có ghi chú mới được thêm",
				},
			)
			if err != nil {
				logger.Warn(
					"register property event failed",
					slog.Any("error", err),
				)
			}
		}()
		// go func() {
		// 	contextTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		// 	defer cancel()

		// 	err := u.NotificationClient.CreateNotification(
		// 		contextTimeout,
		// 		infoEntity.Avatar,
		// 		"Chào mừng bạn đến với BDSPro",
		// 		[]string{"Thiết lập tài khoản thành công"},
		// 		_enum.NotificationSystem,
		// 		&infoEntity.ProfileID,
		// 		infoEntity.ProfileID,
		// 		_enum.EOwnerOfMember,
		// 		[]string{},
		// 	)
		// 	if err != nil {
		// 		log.Printf("Failed to send welcome notification: %v", err)
		// 	} else {
		// 		log.Printf("Successfully sent welcome notification to user %d", infoEntity.ProfileID)
		// 	}
		// }()
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
) ([]*dto.ProductNote, int, error) {
	rows, total, err := u.noteRepo.GetListByProductID(ctx, filter)
	if err != nil || len(rows) == 0 {
		return nil, int(total), err
	}

	noteIDs := make([]uint64, 0, len(rows))
	noteMap := make(map[uint64]*dto.ProductNote)

	for _, r := range rows {
		noteIDs = append(noteIDs, r.ID)
		noteMap[r.ID] = &dto.ProductNote{
			Note:     r,
			Mentions: []*dto.ProductNoteMentionItem{},
			Files:    []*domain.ProductNoteFile{},
		}
	}

	mentions, _ := u.noteRepo.ListMentionsByNoteIDs(ctx, noteIDs)
	for _, m := range mentions {
		noteMap[m.NoteID].Mentions = append(
			noteMap[m.NoteID].Mentions,
			m,
		)
	}

	files, _ := u.noteRepo.ListByNoteIDs(ctx, noteIDs)
	for _, f := range files {
		noteMap[f.NoteID].Files = append(noteMap[f.NoteID].Files, f)
	}

	result := make([]*dto.ProductNote, 0, len(rows))
	for _, r := range rows {
		result = append(result, noteMap[r.ID])
	}

	return result, int(total), nil
}

func (u *ProductNoteUsecase) DeleteNote(
	ctx context.Context,
	noteID uint64,
) error {
	requestUserID := _utils.GetProfileIdWithContext(ctx)
	if requestUserID == 0 {
		return _errors.UnauthorizedException()
	}

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
