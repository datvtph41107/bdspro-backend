package transformer

import (
	"time"

	organizationpb "pb/types/organization"

	"organization/internal/domain/entity"
)

type GroupDocumentTransformer interface {
	CreateGroupDocumentRequestToEntity(request *organizationpb.CreateGroupDocumentRequest) *entity.GroupDocument
	UpdateGroupDocumentRequestToEntity(request *organizationpb.UpdateGroupDocumentRequest) *entity.GroupDocument
	EntityToCreateGroupDocumentResponse(document *entity.GroupDocument) *organizationpb.CreateGroupDocumentResponse
	EntityToUpdateGroupDocumentResponse(document *entity.GroupDocument) *organizationpb.UpdateGroupDocumentResponse
	EntityToGetGroupDocumentResponse(document *entity.GroupDocument) *organizationpb.GetGroupDocumentResponse
	EntityToGetGroupDocumentsResponse(documents []*entity.GroupDocument, total uint32) *organizationpb.GetGroupDocumentsResponse
}

type groupDocumentTransformer struct{}

func NewGroupDocumentTransformer() GroupDocumentTransformer {
	return &groupDocumentTransformer{}
}

func (t *groupDocumentTransformer) CreateGroupDocumentRequestToEntity(request *organizationpb.CreateGroupDocumentRequest) *entity.GroupDocument {
	return &entity.GroupDocument{
		GroupId:     request.GroupId,
		Name:        request.Name,
		Description: request.Description,
		FileUrl:     request.FileUrl,
		FileType:    request.FileType,
		FileSize:    request.FileSize,
	}
}

func (t *groupDocumentTransformer) UpdateGroupDocumentRequestToEntity(request *organizationpb.UpdateGroupDocumentRequest) *entity.GroupDocument {
	return &entity.GroupDocument{
		Id:          request.Id,
		Name:        request.Name,
		Description: request.Description,
		FileUrl:     request.FileUrl,
		FileType:    request.FileType,
		FileSize:    request.FileSize,
	}
}

func (t *groupDocumentTransformer) EntityToCreateGroupDocumentResponse(document *entity.GroupDocument) *organizationpb.CreateGroupDocumentResponse {
	return &organizationpb.CreateGroupDocumentResponse{
		Id: document.Id,
	}
}

func (t *groupDocumentTransformer) EntityToUpdateGroupDocumentResponse(document *entity.GroupDocument) *organizationpb.UpdateGroupDocumentResponse {
	return &organizationpb.UpdateGroupDocumentResponse{
		Id: document.Id,
	}
}

func (t *groupDocumentTransformer) EntityToGetGroupDocumentResponse(document *entity.GroupDocument) *organizationpb.GetGroupDocumentResponse {
	return &organizationpb.GetGroupDocumentResponse{
		Document: &organizationpb.GroupDocument{
			Id:          document.Id,
			GroupId:     document.GroupId,
			Name:        document.Name,
			Description: document.Description,
			FileUrl:     document.FileUrl,
			FileType:    document.FileType,
			FileSize:    document.FileSize,
			CreatedAt:   document.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   document.UpdatedAt.Format(time.RFC3339),
		},
	}
}

func (t *groupDocumentTransformer) EntityToGetGroupDocumentsResponse(documents []*entity.GroupDocument, total uint32) *organizationpb.GetGroupDocumentsResponse {
	documentsResponse := make([]*organizationpb.GroupDocument, len(documents))
	for i, document := range documents {
		documentsResponse[i] = &organizationpb.GroupDocument{
			Id:          document.Id,
			GroupId:     document.GroupId,
			Name:        document.Name,
			Description: document.Description,
			FileUrl:     document.FileUrl,
			FileType:    document.FileType,
			FileSize:    document.FileSize,
			CreatedAt:   document.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   document.UpdatedAt.Format(time.RFC3339),
			Uploader:    document.Uploader,
		}
	}
	return &organizationpb.GetGroupDocumentsResponse{
		Data:  documentsResponse,
		Total: total,
	}
}
