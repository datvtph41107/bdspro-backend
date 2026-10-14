package dto

import _dto "common/domain/dto"

// CatalogSearchRequest mirrors bdspropb.SearchQueryRequest for list catalog routes.
type CatalogSearchRequest struct {
	Page uint32 `form:"page"`
	Size uint32 `form:"size"`
	Text string `form:"text"`
}

func (r CatalogSearchRequest) Pagable() _dto.Pagable {
	return _dto.Pagable{Page: r.Page, Size: r.Size, Text: r.Text}
}
