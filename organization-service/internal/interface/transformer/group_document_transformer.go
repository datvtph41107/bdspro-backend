package transformer

import (
	"time"

	grouppb "pb/types/group"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
)

type GroupDocumentTransformer interface {
	CreateGroupDocumentRequestToEntity(request *grouppb.CreateGroupDocumentRequest) *entity.GroupDocument
	UpdateGroupDocumentRequestToEntity(request *grouppb.UpdateGroupDocumentRequest) *entity.GroupDocument
	EntityToCreateGroupDocumentResponse(document *entity.GroupDocument) *grouppb.CreateGroupDocumentResponse
	EntityToUpdateGroupDocumentResponse(document *entity.GroupDocument) *grouppb.UpdateGroupDocumentResponse
	EntityToGetGroupDocumentResponse(document *entity.GroupDocument) *grouppb.GetGroupDocumentResponse
	EntityToGetGroupDocumentsResponse(documents []*entity.GroupDocument) *grouppb.GetGroupDocumentsResponse
}

type groupDocumentTransformer struct{}

func NewGroupDocumentTransformer() GroupDocumentTransformer {
	return &groupDocumentTransformer{}
}

func (t *groupDocumentTransformer) CreateGroupDocumentRequestToEntity(request *grouppb.CreateGroupDocumentRequest) *entity.GroupDocument {
	return &entity.GroupDocument{
		GroupId:     request.GroupId,
		Name:        request.Name,
		Description: request.Description,
		FileUrl:     request.FileUrl,
		FileType:    request.FileType,
		FileSize:    request.FileSize,
	}
}

func (t *groupDocumentTransformer) UpdateGroupDocumentRequestToEntity(request *grouppb.UpdateGroupDocumentRequest) *entity.GroupDocument {
	return &entity.GroupDocument{
		Id:          request.Id,
		Name:        request.Name,
		Description: request.Description,
		FileUrl:     request.FileUrl,
		FileType:    request.FileType,
		FileSize:    request.FileSize,
	}
}

func (t *groupDocumentTransformer) EntityToCreateGroupDocumentResponse(document *entity.GroupDocument) *grouppb.CreateGroupDocumentResponse {
	return &grouppb.CreateGroupDocumentResponse{
		Id: document.Id,
	}
}

func (t *groupDocumentTransformer) EntityToUpdateGroupDocumentResponse(document *entity.GroupDocument) *grouppb.UpdateGroupDocumentResponse {
	return &grouppb.UpdateGroupDocumentResponse{
		Id: document.Id,
	}
}

func (t *groupDocumentTransformer) EntityToGetGroupDocumentResponse(document *entity.GroupDocument) *grouppb.GetGroupDocumentResponse {
	return &grouppb.GetGroupDocumentResponse{
		Document: &grouppb.GroupDocument{
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

func (t *groupDocumentTransformer) EntityToGetGroupDocumentsResponse(documents []*entity.GroupDocument) *grouppb.GetGroupDocumentsResponse {
	documentsResponse := make([]*grouppb.GroupDocument, len(documents))
	for i, document := range documents {
		documentsResponse[i] = &grouppb.GroupDocument{
			Id:          document.Id,
			GroupId:     document.GroupId,
			Name:        document.Name,
			Description: document.Description,
			FileUrl:     document.FileUrl,
			FileType:    document.FileType,
			FileSize:    document.FileSize,
			CreatedAt:   document.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   document.UpdatedAt.Format(time.RFC3339),
		}
	}
	return &grouppb.GetGroupDocumentsResponse{
		Documents: documentsResponse,
		Total:     uint32(len(documents)),
	}
} 