package usecase

import (
	"context"

	"github.com/tuandq2112/organization-service/internal/custom_error"
	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
)

type GroupDocumentUsecase interface {
	CreateGroupDocument(ctx context.Context, document *entity.GroupDocument) (*entity.GroupDocument, error)
	UpdateGroupDocument(ctx context.Context, document *entity.GroupDocument) (*entity.GroupDocument, error)
	DeleteGroupDocument(ctx context.Context, id uint32) error
	GetGroupDocumentByID(ctx context.Context, id uint32) (*entity.GroupDocument, error)
	GetGroupDocumentsByGroupID(ctx context.Context, groupID uint32) ([]*entity.GroupDocument, error)
}

type groupDocumentUsecase struct {
	groupDocumentRepository repository.GroupDocumentRepository
}

func NewGroupDocumentUsecase(groupDocumentRepository repository.GroupDocumentRepository) GroupDocumentUsecase {
	return &groupDocumentUsecase{groupDocumentRepository: groupDocumentRepository}
}

func (u *groupDocumentUsecase) CreateGroupDocument(ctx context.Context, document *entity.GroupDocument) (*entity.GroupDocument, error) {
	return u.groupDocumentRepository.Create(ctx, document)
}

func (u *groupDocumentUsecase) UpdateGroupDocument(ctx context.Context, document *entity.GroupDocument) (*entity.GroupDocument, error) {
	existedDocument, err := u.groupDocumentRepository.GetByID(ctx, document.Id)
	if err != nil {
		return nil, err
	}

	if existedDocument == nil {
		return nil, custom_error.RecordNotFound("document not found")
	}

	if document.CreatedBy != existedDocument.CreatedBy {
		return nil, custom_error.Forbidden("you are not allowed to update this document")
	}

	return u.groupDocumentRepository.Update(ctx, document)
}

func (u *groupDocumentUsecase) DeleteGroupDocument(ctx context.Context, id uint32) error {
	return u.groupDocumentRepository.Delete(ctx, id)
}

func (u *groupDocumentUsecase) GetGroupDocumentByID(ctx context.Context, id uint32) (*entity.GroupDocument, error) {
	return u.groupDocumentRepository.GetByID(ctx, id)
}

func (u *groupDocumentUsecase) GetGroupDocumentsByGroupID(ctx context.Context, groupID uint32) ([]*entity.GroupDocument, error) {
	return u.groupDocumentRepository.GetByGroupID(ctx, groupID)
} 