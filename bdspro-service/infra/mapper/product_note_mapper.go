package mapper

import (
	"bdspro/internal/dto"

	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
)

type ProductNoteMapper struct{}

func NewProductNoteMapper() *ProductNoteMapper {
	return &ProductNoteMapper{}
}

func (m *ProductNoteMapper) ProductNoteListToProto(
	list []*dto.ProductNoteItemResponse,
) []*bdspropb.ProductNoteEntity {

	if len(list) == 0 {
		return nil
	}

	res := make([]*bdspropb.ProductNoteEntity, 0, len(list))
	for _, n := range list {
		res = append(res, m.ProductNoteItemToProto(n))
	}
	return res
}

func (m *ProductNoteMapper) ProductNoteItemToProto(
	n *dto.ProductNoteItemResponse,
) *bdspropb.ProductNoteEntity {
	return &bdspropb.ProductNoteEntity{
		Id:        n.ID,
		ProductId: n.ProductID,
		Content:   n.Content,
		IsPinned:  n.IsPinned,
		Author:    mapAuthorToProto(n.Author),
		Mentions:  mapMentionsToProto(n.Mentions),
		Files:     mapFilesToProto(n.Files),
		CreatedAt: n.CreatedAt,
	}
}

func mapAuthorToProto(p *sharepb.ProfileItem) *bdspropb.ProductNoteMember {
	if p == nil {
		return nil
	}

	return &bdspropb.ProductNoteMember{
		Id:     p.Id,
		Name:   p.FullName,
		Avatar: p.Avatar,
		Role:   "member",
	}
}

func mapMentionsToProto(
	list []dto.ProductNoteMentionItem,
) []*bdspropb.ProductNoteMention {
	if len(list) == 0 {
		return nil
	}

	res := make([]*bdspropb.ProductNoteMention, 0, len(list))
	for _, m := range list {
		res = append(res, &bdspropb.ProductNoteMention{
			UserId:   m.UserId,
			StartPos: m.StartPos,
			Length:   m.Length,
			Name:     m.Name,
			Avatar:   m.Avatar,
			Role:     m.Role,
		})
	}
	return res
}

func mapFilesToProto(
	list []dto.ProductNoteFileDTO,
) []*bdspropb.ProductNoteFile {
	if len(list) == 0 {
		return nil
	}

	res := make([]*bdspropb.ProductNoteFile, 0, len(list))
	for _, f := range list {
		res = append(res, &bdspropb.ProductNoteFile{
			Url:      f.URL,
			FileName: f.FileName,
			FileType: f.FileType,
			Size:     f.Size,
		})
	}
	return res
}
