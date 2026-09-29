package usecase

import (
	"context"
	"organization/env"
	"organization/internal/custom_error"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	iusecase "organization/internal/interface"
	"organization/pkg/utils"
	sharepb "pb/types/shared"
)

type GroupDocumentUsecase interface {
	CreateGroupDocument(ctx context.Context, document *entity.GroupDocument) (*entity.GroupDocument, error)
	UpdateGroupDocument(ctx context.Context, document *entity.GroupDocument) (*entity.GroupDocument, error)
	DeleteGroupDocument(ctx context.Context, id uint32) error
	GetGroupDocumentByID(ctx context.Context, id uint32) (*entity.GroupDocument, error)
	GetGroupDocumentsByGroupID(ctx context.Context, groupID uint32, page, size int) ([]*entity.GroupDocument, uint32, error)
}

type groupDocumentUsecase struct {
	groupDocumentRepository repository.GroupDocumentRepository
	logActivityUsecase      GroupLogActivityUsecase
	groupAuthUsecase        GroupAuthUsecase
	userGrpcClient          iusecase.IUserClient
}

func NewGroupDocumentUsecase(groupDocumentRepository repository.GroupDocumentRepository,
	logActivityUsecase GroupLogActivityUsecase,
	groupAuthUsecase GroupAuthUsecase,
	userGrpcClient iusecase.IUserClient,
) GroupDocumentUsecase {
	return &groupDocumentUsecase{
		groupDocumentRepository: groupDocumentRepository,
		logActivityUsecase:      logActivityUsecase,
		groupAuthUsecase:        groupAuthUsecase,
		userGrpcClient:          userGrpcClient,
	}
}

func (u *groupDocumentUsecase) CreateGroupDocument(ctx context.Context, document *entity.GroupDocument) (*entity.GroupDocument, error) {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	err := u.groupAuthUsecase.IsMemberGroup(ctx, document.GroupId, currentUserId)
	if err != nil {
		return nil, custom_error.Forbidden("you are not allowed to create document")
	}
	document.CreatedBy = currentUserId
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

	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	err = u.groupAuthUsecase.IsMemberGroup(ctx, existedDocument.GroupId, currentUserId)
	if err != nil {
		return nil, custom_error.Forbidden("you are not allowed to update this document")
	}

	document.UpdatedBy = currentUserId

	return u.groupDocumentRepository.Update(ctx, document)
}

func (u *groupDocumentUsecase) DeleteGroupDocument(ctx context.Context, id uint32) error {
	document, err := u.groupDocumentRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if document == nil {
		return custom_error.RecordNotFound("document not found")
	}

	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	err = u.groupAuthUsecase.IsMemberGroup(ctx, document.GroupId, currentUserId)
	if err != nil {
		return custom_error.Forbidden("you are not allowed to delete document")
	}

	return u.groupDocumentRepository.Delete(ctx, id)
}

func (u *groupDocumentUsecase) GetGroupDocumentByID(ctx context.Context, id uint32) (*entity.GroupDocument, error) {
	document, err := u.groupDocumentRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if document == nil {
		return nil, custom_error.RecordNotFound("document not found")
	}

	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	err = u.groupAuthUsecase.IsMemberGroup(ctx, document.GroupId, currentUserId)
	if err != nil {
		return nil, custom_error.Forbidden("you are not allowed to get document")
	}

	return document, nil
}

func (u *groupDocumentUsecase) GetGroupDocumentsByGroupID(ctx context.Context, groupID uint32, page, size int) ([]*entity.GroupDocument, uint32, error) {
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	err := u.groupAuthUsecase.IsMemberGroup(ctx, groupID, currentUserId)
	if err != nil {
		return nil, 0, custom_error.Forbidden("you are not allowed to get documents")
	}

	documents, total, err := u.groupDocumentRepository.GetByGroupID(ctx, groupID, page, size)
	if err != nil {
		return nil, 0, err
	}
	profileIds := make([]uint64, len(documents))
	for i, document := range documents {
		profileIds[i] = uint64(document.CreatedBy)
	}
	profiles, err := u.userGrpcClient.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{
		Ids: profileIds,
	})
	if err != nil {
		return nil, 0, err
	}
	profileMap := make(map[uint64]string)
	for _, profile := range profiles.Profiles {
		profileMap[profile.Id] = profile.FullName
	}
	for _, document := range documents {
		document.Uploader = profileMap[uint64(document.CreatedBy)]
	}

	return documents, total, nil
}
