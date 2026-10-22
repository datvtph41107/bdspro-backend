package controller

import (
	"context"
	"net/http"
	"strconv"

	_dto "common/domain/dto"
	_errors "common/errors"
	"common/httpapi"
	"mdn4/modules/tqd/dto"
	"mdn4/modules/tqd/services"
)

// OpenHourController owns /v2/tqd/open-hours*.
type OpenHourController struct {
	uc services.OpenHourUsecase
}

func NewOpenHourController(uc services.OpenHourUsecase) *OpenHourController {
	return &OpenHourController{uc: uc}
}

func (c *OpenHourController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.POST(mux, "/v2/tqd/open-hours", c.create)
	httpapi.GET(mux, "/v2/tqd/open-hours/times", c.list)
	httpapi.GET(mux, "/v2/tqd/open-hours/by-day/{dayOfWeek}", c.byDay)
	httpapi.GET(mux, "/v2/tqd/open-hours/{id}", c.get)
	httpapi.PUT(mux, "/v2/tqd/open-hours/{id}", c.update)
	httpapi.DELETE(mux, "/v2/tqd/open-hours/{id}", c.delete)
}

type openHourCreateReq struct {
	POIID       uint64  `json:"poiId"`
	DayOfWeek   int     `json:"dayOfWeek"`
	OpenTime    string  `json:"openTime"`
	CloseTime   string  `json:"closeTime"`
	Note        string  `json:"note"`
	IsOpen      bool    `json:"isOpen"`
	IsRecurring bool    `json:"isRecurring"`
	StartDate   *string `json:"startDate"`
	EndDate     *string `json:"endDate"`
	Type        int     `json:"type"`
}

type openHourUpdateReq struct {
	ID          uint64  `path:"id" json:"-"`
	DayOfWeek   *int    `json:"dayOfWeek"`
	OpenTime    *string `json:"openTime"`
	CloseTime   *string `json:"closeTime"`
	Note        *string `json:"note"`
	IsOpen      *bool   `json:"isOpen"`
	IsRecurring *bool   `json:"isRecurring"`
	StartDate   *string `json:"startDate"`
	EndDate     *string `json:"endDate"`
	Status      *int    `json:"status"`
	Type        *int    `json:"type"`
}

type listOpenHourReq struct {
	Page      uint32 `form:"page"`
	Size      uint32 `form:"size"`
	POIID     uint64 `form:"poiId"`
	DayOfWeek *int   `form:"dayOfWeek"`
	IsOpen    *bool  `form:"isOpen"`
	Status    *int   `form:"status"`
	Type      *int   `form:"type"`
}

type byDayReq struct {
	DayOfWeek int    `path:"dayOfWeek"`
	POIID     uint64 `form:"poiId"`
}

func (c *OpenHourController) create(ctx context.Context, req openHourCreateReq) (*httpapi.Result, error) {
	out, err := c.uc.Create(ctx, &dto.CreateOpenHourRequest{
		POIID: req.POIID, DayOfWeek: req.DayOfWeek, OpenTime: req.OpenTime, CloseTime: req.CloseTime,
		Note: req.Note, IsOpen: req.IsOpen, IsRecurring: req.IsRecurring, StartDate: req.StartDate,
		EndDate: req.EndDate, Type: req.Type,
	})
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *OpenHourController) get(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	out, err := c.uc.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *OpenHourController) update(ctx context.Context, req openHourUpdateReq) (*httpapi.Result, error) {
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	out, err := c.uc.Update(ctx, req.ID, &dto.UpdateOpenHourRequest{
		DayOfWeek: req.DayOfWeek, OpenTime: req.OpenTime, CloseTime: req.CloseTime, Note: req.Note,
		IsOpen: req.IsOpen, IsRecurring: req.IsRecurring, StartDate: req.StartDate, EndDate: req.EndDate,
		Status: req.Status, Type: req.Type,
	})
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *OpenHourController) delete(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	if err := c.uc.Delete(ctx, req.ID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *OpenHourController) list(ctx context.Context, req listOpenHourReq) (*httpapi.Result, error) {
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	out, err := c.uc.List(ctx, &dto.OpenHourFilter{
		Pagable: pagable, POIID: req.POIID, DayOfWeek: req.DayOfWeek, IsOpen: req.IsOpen, Status: req.Status, Type: req.Type,
	})
	if err != nil {
		return nil, err
	}
	if out == nil {
		return jsonListResult([]dto.OpenHourResponse{}, 0, pagable.GetPage(), pagable.GetSize())
	}
	return jsonListResult(out.Data, out.Total, pagable.GetPage(), pagable.GetSize())
}

func (c *OpenHourController) byDay(ctx context.Context, req byDayReq) (*httpapi.Result, error) {
	out, err := c.uc.GetTimesByDay(ctx, &dto.GetTimesByDayRequest{DayOfWeek: strconv.Itoa(req.DayOfWeek), POIID: req.POIID})
	if err != nil {
		return nil, err
	}
	if out == nil {
		return jsonListResult([]any{}, 0, 1, 0)
	}
	return jsonListResult(out, int64(len(out)), 1, uint32(len(out)))
}
