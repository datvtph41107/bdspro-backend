package controller

import (
	"context"
	"net/http"

	crmpb "pb/types/crm"
	sharepb "pb/types/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"common/httpapi"
	"mdn4/modules/crm/services"
)

// EnumController serves CRM enum lookups (crm/enum.proto).
type EnumController struct {
	enum *services.EnumUsecase
}

func NewEnumController(enum *services.EnumUsecase) *EnumController {
	return &EnumController{enum: enum}
}

func (c *EnumController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.HandleProto(mux, http.MethodGet, "/v2/enums/crm/{name}", c.GetEnumByName)
}

func (c *EnumController) GetEnumByName(ctx context.Context, req *crmpb.GetEnumRequest) (*sharepb.ItemV3PageProto, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "Tên enum không được để trống")
	}
	items, err := c.enum.GetEnumByName(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return &sharepb.ItemV3PageProto{Data: items, Total: uint32(len(items))}, nil
}
