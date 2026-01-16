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

func (m *ProductNoteMapper) ProductNoteListToProto(list []*dto.ProductNoteItemRespone) []*bdspropb.ProductNoteEntity {
	if len(list) == 0 {
		return nil
	}

	res := make([]*bdspropb.ProductNoteEntity, 0, len(list))
	for _, n := range list {
		res = append(res, m.ProductNoteItemToProto(n))
	}
	return res
}

func (m *ProductNoteMapper) ProductNoteItemToProto(n *dto.ProductNoteItemRespone) *bdspropb.ProductNoteEntity {
	return &bdspropb.ProductNoteEntity{
		Id:       n.ID,
		Content:  n.Content,
		IsPinned: n.IsPinned,
		Author:   mapAuthorToProto(n.Author),
		Mentions: mapMentionsToProto(n.Mentions),
		// Files:     n.Files,
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

func mapMentionsToProto(list []dto.ProductNoteMentionItem) []*bdspropb.ProductNoteMember {
	if len(list) == 0 {
		return nil
	}

	res := make([]*bdspropb.ProductNoteMember, 0, len(list))
	for _, m := range list {
		res = append(res, &bdspropb.ProductNoteMember{
			Id:     m.UserId,
			Name:   m.Name,
			Avatar: m.Avatar,
			Role:   m.Role,
		})
	}
	return res
}

// func mapFilesToProto(list []domain.ProductNoteFile) []*bdspropb.ProductNoteFile {
func mapFilesToProto(list []string) []string {
	if len(list) == 0 {
		return nil
	}

	// res := make([]*bdspropb.ProductNoteFile, 0, len(list))
	res := make([]string, 0, len(list))
	for _, f := range list {
		// res = append(res, &bdspropb.ProductNoteFile{
		// 	Id:       f.ID,
		// 	FileName: f.FileName,
		// 	FileUrl:  f.URL,
		// 	FileSize: f.Size,
		// })
		res = append(res, f)
	}
	return res
}
