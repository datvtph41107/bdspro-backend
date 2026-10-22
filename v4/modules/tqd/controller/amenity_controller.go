package controller

import (
	"context"
	"net/http"

	_dto "common/domain/dto"
	_errors "common/errors"
	"common/httpapi"
	"mdn4/modules/tqd/dto"
	"mdn4/modules/tqd/mapper"
	"mdn4/modules/tqd/services"
)

// AmenityController owns /v2/tqd/amenities*.
type AmenityController struct {
	uc     *services.AmenityUsecase
	mapper *mapper.AmenityMapper
}

func NewAmenityController(uc *services.AmenityUsecase) *AmenityController {
	return &AmenityController{uc: uc, mapper: mapper.NewAmenityMapper()}
}

func (c *AmenityController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.POST(mux, "/v2/tqd/amenities", c.create)
	httpapi.GET(mux, "/v2/tqd/amenities", c.list)
	httpapi.GET(mux, "/v2/tqd/amenities/{id}", c.get)
	httpapi.PUT(mux, "/v2/tqd/amenities/{id}", c.update)
	httpapi.DELETE(mux, "/v2/tqd/amenities/{id}", c.delete)
}

type amenityBody struct {
	ID          uint64 `path:"id" json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Category    string `json:"category"`
	SortOrder   int    `json:"sortOrder"`
	IsActive    bool   `json:"isActive"`
}

type listAmenityReq struct {
	Page     uint32 `form:"page"`
	Size     uint32 `form:"size"`
	Search   string `form:"search"`
	Category string `form:"category"`
}

func (c *AmenityController) create(ctx context.Context, req amenityBody) (*httpapi.Result, error) {
	entity := c.mapper.ToDomainFromCreateRequest(&dto.CreateAmenityRequestDTO{
		Name: req.Name, Description: req.Description, Icon: req.Icon,
		Category: req.Category, SortOrder: req.SortOrder, IsActive: req.IsActive,
	})
	out, err := c.uc.Create(ctx, entity)
	if err != nil {
		return nil, err
	}
	return protoResult(c.mapper.ToProto(out))
}

func (c *AmenityController) get(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	out, err := c.uc.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return protoResult(c.mapper.ToProto(out))
}

func (c *AmenityController) update(ctx context.Context, req amenityBody) (*httpapi.Result, error) {
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	entity := c.mapper.ToDomainFromUpdateRequest(&dto.UpdateAmenityRequestDTO{
		Name: req.Name, Description: req.Description, Icon: req.Icon,
		Category: req.Category, SortOrder: req.SortOrder, IsActive: req.IsActive,
	}, req.ID)
	out, err := c.uc.Update(ctx, req.ID, entity)
	if err != nil {
		return nil, err
	}
	return protoResult(c.mapper.ToProto(out))
}

func (c *AmenityController) delete(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	if _, err := c.uc.Delete(ctx, req.ID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true, "message": "Amenity deleted successfully"})
}

func (c *AmenityController) list(ctx context.Context, req listAmenityReq) (*httpapi.Result, error) {
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	filter := &dto.AmenityFilterDTO{Pagable: pagable, Search: req.Search, Category: req.Category}
	rows, total, err := c.uc.ListAmenities(ctx, filter)
	if err != nil {
		return nil, err
	}
	data := make([]any, 0, len(rows))
	for i := range rows {
		pb := c.mapper.ToProto(&rows[i])
		v, err := protoToAny(pb)
		if err != nil {
			return nil, err
		}
		data = append(data, v)
	}
	return httpapi.OK(data, httpapi.WithPagination(&httpapi.Pagination{
		Page: int32(pagable.GetPage()), Size: int32(pagable.GetSize()), Total: total,
	})), nil
}
