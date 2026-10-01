package controller

import (
	"context"
	"net/http"

	base_enum "base/enum"
	_errors "common/errors"
	"userv4/internal/httpapi"
	"userv4/modules/crm/dto"
	"userv4/modules/crm/services"
)

type LeadController struct{ uc *services.LeadUsecase }

func NewLeadController(uc *services.LeadUsecase) *LeadController { return &LeadController{uc: uc} }

func (c *LeadController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/crm/lead/detail/{id}", c.detail)
	httpapi.GET(mux, "/v2/crm/lead/by-contact/{id}", c.byContact)
	httpapi.POST(mux, "/v2/crm/lead/{ownerOf}/{ownerId}", c.create)
	httpapi.PUT(mux, "/v2/crm/lead/{id}", c.update)
	httpapi.DELETE(mux, "/v2/crm/lead/{id}", c.delete)
	httpapi.POST(mux, "/v2/crm/lead/assign-person", c.assignPerson)
	httpapi.POST(mux, "/v2/crm/lead/assign-crm", c.assignCRM)
	httpapi.GET(mux, "/v2/crm/lead/manage/{ownerOf}/{ownerId}/{pipelineId}", c.manage)
	httpapi.GET(mux, "/v2/crm/lead/care-expired", c.careExpired)
	httpapi.PUT(mux, "/v2/crm/lead/stage", c.switchStage)
}

func (c *LeadController) require() error {
	if c == nil || c.uc == nil {
		return _errors.ReturnError(501, "lead service not configured")
	}
	return nil
}

type leadCreateReq struct {
	OwnerOf  int32  `path:"ownerOf"`
	OwnerID  uint64 `path:"ownerId"`
	dto.LeadSaveDTO
}
type leadManageReq struct {
	OwnerOf    int32  `path:"ownerOf"`
	OwnerID    uint64 `path:"ownerId"`
	PipelineID uint64 `path:"pipelineId"`
	Page       uint32 `form:"page"`
	Size       uint32 `form:"size"`
}
type leadAssignPersonReq struct {
	ID             uint64 `json:"id"`
	ChargePersonID uint64 `json:"chargePersonId"`
	AssignNote     string `json:"assignNote"`
}
type leadAssignCRMReq struct {
	ContactID      uint64 `json:"contactId"`
	PersonChargeID uint64 `json:"personChargeId"`
	StageID        uint64 `json:"stageId"`
	Note           string `json:"note"`
}
type leadSwitchStageReq struct {
	ID      uint64  `json:"id"`
	StageID *uint64 `json:"stageId"`
	Note    string  `json:"note"`
}
type leadUpdateReq struct {
	ID uint64 `path:"id" json:"id"`
	dto.LeadSaveDTO
}

func (c *LeadController) detail(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	if err := c.require(); err != nil { return nil, err }
	out, err := c.uc.GetByID(ctx, req.ID)
	if err != nil { return nil, err }
	return httpapi.OK(out), nil
}
func (c *LeadController) byContact(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	if err := c.require(); err != nil { return nil, err }
	out, err := c.uc.GetByContact(ctx, req.ID)
	if err != nil { return nil, err }
	return httpapi.OK(out), nil
}
func (c *LeadController) create(ctx context.Context, req leadCreateReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil { return nil, err }
	in := req.LeadSaveDTO
	in.OwnerID = req.OwnerID
	in.OwnerOf = base_enum.EOwnerOf(req.OwnerOf)
	out, err := c.uc.Create(ctx, &in)
	if err != nil { return nil, err }
	return httpapi.OK(out), nil
}
func (c *LeadController) update(ctx context.Context, req leadUpdateReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil { return nil, err }
	out, err := c.uc.Update(ctx, req.ID, &req.LeadSaveDTO)
	if err != nil { return nil, err }
	return httpapi.OK(out), nil
}
func (c *LeadController) delete(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	if err := c.require(); err != nil { return nil, err }
	if err := c.uc.Delete(ctx, req.ID); err != nil { return nil, err }
	return submitResult(req.ID, "deleted")
}
func (c *LeadController) assignPerson(ctx context.Context, req leadAssignPersonReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil { return nil, err }
	out, err := c.uc.Assign(ctx, req.ID, &dto.LeadDTO{ChargePersonId: &req.ChargePersonID, AssignNote: req.AssignNote})
	if err != nil { return nil, err }
	return httpapi.OK(out), nil
}
func (c *LeadController) assignCRM(ctx context.Context, req leadAssignCRMReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil { return nil, err }
	out, err := c.uc.AssignCrm(ctx, req.ContactID, req.PersonChargeID, req.StageID, req.Note)
	if err != nil { return nil, err }
	return httpapi.OK(out), nil
}
func (c *LeadController) manage(ctx context.Context, req leadManageReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil { return nil, err }
	items, total, err := c.uc.GetManager(ctx, dto.LeadSearchDTO{
		OwnerOf: base_enum.EOwnerOf(req.OwnerOf), OwnerID: req.OwnerID, PipelineID: &req.PipelineID,
	})
	if err != nil { return nil, err }
	_ = req.Page; _ = req.Size
	return httpapi.OK(items, httpapi.WithPagination(&httpapi.Pagination{Total: total})), nil
}
func (c *LeadController) careExpired(ctx context.Context, _ httpapi.Empty) (*httpapi.Result, error) {
	if err := c.require(); err != nil { return nil, err }
	n, err := c.uc.GetCareExpiredTime(ctx)
	if err != nil { return nil, err }
	return httpapi.OK(map[string]any{"careExpired": n}), nil
}
func (c *LeadController) switchStage(ctx context.Context, req leadSwitchStageReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil { return nil, err }
	out, err := c.uc.SwitchStage(ctx, req.ID, req.StageID, req.Note)
	if err != nil { return nil, err }
	return httpapi.OK(out), nil
}
