package mapper

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"

	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
)

type ProductNoteMapper struct{}

func NewProductNoteMapper() *ProductNoteMapper {
	return &ProductNoteMapper{}
}

func (m *ProductNoteMapper) ProductNoteListToProto(
	list []*dto.ProductNote,
) []*bdspropb.ProductNoteEntity {
	if len(list) == 0 {
		return nil
	}

	res := make([]*bdspropb.ProductNoteEntity, 0, len(list))
	for _, n := range list {
		res = append(res, m.ProductNoteToProto(n))
	}
	return res
}

func (m *ProductNoteMapper) ProductNoteToProto(
	n *dto.ProductNote,
) *bdspropb.ProductNoteEntity {

	if n == nil || n.Note == nil {
		return nil
	}

	note := n.Note

	return &bdspropb.ProductNoteEntity{
		Id:        note.ID,
		ProductId: note.ProductID,
		Content:   note.Content,
		IsPinned:  note.IsPinned,
		CreatedAt: note.CreatedAt,
		Author:    mapAuthorToProto(note.Author),
		Mentions:  mapMentionsToProto(n.Mentions),
		Files:     mapFilesToProto(n.Files),
	}
}

func mapAuthorToProto(
	p *sharepb.ProfileItem,
) *bdspropb.ProductNoteMember {
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
	list []*dto.ProductNoteMentionItem,
) []*bdspropb.ProductNoteMention {
	if len(list) == 0 {
		return nil
	}

	res := make([]*bdspropb.ProductNoteMention, 0, len(list))
	for _, m := range list {
		if m == nil {
			continue
		}

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
	list []*domain.ProductNoteFile,
) []*bdspropb.ProductNoteFile {
	if len(list) == 0 {
		return nil
	}

	res := make([]*bdspropb.ProductNoteFile, 0, len(list))
	for _, f := range list {
		if f == nil {
			continue
		}

		res = append(res, &bdspropb.ProductNoteFile{
			Url:      f.URL,
			FileName: f.FileName,
			FileType: f.FileType,
			Size:     f.Size,
		})
	}
	return res
}
