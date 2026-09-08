package grpc

import (
	"context"
	"testing"
	"time"

	"common/identity"
	userpb "pb/types/user"
	catalogdomain "user/internal/domain/plan"
	"user/internal/usecase/plan/admin"
	"user/internal/usecase/plan/publish"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type planHandlerRepository struct {
	aggregate catalogdomain.PlanVersionAggregate
}

func (r *planHandlerRepository) GetPlanVersion(
	context.Context,
	uint64,
) (catalogdomain.PlanVersionAggregate, error) {
	return r.aggregate.Copy(), nil
}

func (r *planHandlerRepository) ListPlanVersions(
	context.Context,
	admin.Query,
) (admin.Page, error) {
	return admin.Page{PlanVersions: []catalogdomain.PlanVersionAggregate{r.aggregate}, Total: 1}, nil
}

func (r *planHandlerRepository) CreatePlanVersionDraft(_ context.Context, record admin.DraftRecord) (uint64, error) {
	r.aggregate = handlerAggregateFromDraft(9, record)
	return 9, nil
}

func (r *planHandlerRepository) UpdatePlanVersionDraft(_ context.Context, id uint64, record admin.DraftRecord) error {
	r.aggregate = handlerAggregateFromDraft(id, record)
	return nil
}

func (r *planHandlerRepository) DeletePlanVersionDraft(context.Context, uint64) error { return nil }

func (r *planHandlerRepository) RetirePlanVersion(_ context.Context, _ uint64, _ uint64, retiredAt time.Time) error {
	r.aggregate.Status = catalogdomain.StatusRetired
	r.aggregate.EffectiveUntil = &retiredAt
	return nil
}

type fakePermissionAuthorizer struct {
	allowed map[string]bool
	err     error
	actorID uint64
	code    string
}

func (f *fakePermissionAuthorizer) HasPermission(
	_ context.Context,
	actorID uint64,
	permissionCode string,
) (bool, error) {
	f.actorID = actorID
	f.code = permissionCode
	if f.err != nil {
		return false, f.err
	}
	return f.allowed[permissionCode], nil
}

func allowPermissions(codes ...string) *fakePermissionAuthorizer {
	allowed := make(map[string]bool, len(codes))
	for _, code := range codes {
		allowed[code] = true
	}
	return &fakePermissionAuthorizer{allowed: allowed}
}

type planHandlerPublisher struct{}

func (planHandlerPublisher) Publish(
	_ context.Context,
	command publish.Command,
) (publish.Result, error) {
	return publish.Result{
		PlanVersionID: command.PlanVersionID,
		PublishedAt:   command.PublishedAt,
	}, nil
}

func TestListPlanVersionsRequiresTrustedAdminAndMapsTerms(t *testing.T) {
	handler := NewPlanVersionAdminHandler(planVersionService(), allowPermissions(admin.PermissionPlanView))
	if _, err := handler.ListPlanVersions(context.Background(), nil); status.Code(err) != codes.PermissionDenied {
		t.Fatalf("unsigned error = %v", err)
	}

	response, err := handler.ListPlanVersions(trustedAdminContext(t), &userpb.ListCatalogPlanVersionsRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if response.Total != 1 || len(response.PlanVersions) != 1 ||
		response.PlanVersions[0].GetEntitlements()[0].GetMeterCode() != "workspace.report_generation.accepted" ||
		len(response.PlanVersions[0].GetOperations()) != 1 ||
		response.PlanVersions[0].GetOperations()[0].GetUnitsPerAction() != 1 {
		t.Fatalf("response = %#v", response)
	}
}

func TestPublishPlanVersionUsesServerTimeWhenEffectiveFromIsOmitted(t *testing.T) {
	handler := NewPlanVersionAdminHandler(planVersionService(), allowPermissions(admin.PermissionPlanPublish))
	response, err := handler.PublishPlanVersion(trustedAdminContext(t), &userpb.PublishCatalogPlanVersionRequest{
		PlanVersionId: 9,
	})
	if err != nil {
		t.Fatalf("PublishPlanVersion() error = %v", err)
	}
	if response.GetPublishedAt() == nil || response.GetPublishedAt().AsTime().IsZero() {
		t.Fatalf("published_at = %v", response.GetPublishedAt())
	}
}

func TestCatalogAuthorizationUsesPermissionNotRoleName(t *testing.T) {
	view := allowPermissions(admin.PermissionPlanView)
	handler := NewPlanVersionAdminHandler(planVersionService(), view)

	// A non-admin role with the durable VIEW assignment is authorized.
	if _, err := handler.ListPlanVersions(
		trustedActorContext(t, "ROLE_SUPPORT"),
		&userpb.ListCatalogPlanVersionsRequest{},
	); err != nil {
		t.Fatalf("explicit permission rejected because of role name: %v", err)
	}
	if view.actorID != 42 || view.code != admin.PermissionPlanView {
		t.Fatalf("authorization call = actor:%d code:%q", view.actorID, view.code)
	}

	// A SUPER_ADMIN role name without the durable assignment does not bypass IAM.
	denied := NewPlanVersionAdminHandler(planVersionService(), allowPermissions())
	_, err := denied.ListPlanVersions(
		trustedActorContext(t, "ROLE_SUPER_ADMIN"),
		&userpb.ListCatalogPlanVersionsRequest{},
	)
	if status.Code(err) != codes.PermissionDenied {
		t.Fatalf("role-name bypass error = %v", err)
	}
}

func TestPublishRequiresPublishPermissionNotView(t *testing.T) {
	handler := NewPlanVersionAdminHandler(planVersionService(), allowPermissions(admin.PermissionPlanView))
	_, err := handler.PublishPlanVersion(
		trustedAdminContext(t),
		&userpb.PublishCatalogPlanVersionRequest{PlanVersionId: 9},
	)
	if status.Code(err) != codes.PermissionDenied {
		t.Fatalf("view-only publish error = %v", err)
	}
}

func TestCatalogAuthorizationFailsClosedWithoutAuthority(t *testing.T) {
	handler := NewPlanVersionAdminHandler(planVersionService(), nil)
	_, err := handler.ListPlanVersions(
		trustedAdminContext(t),
		&userpb.ListCatalogPlanVersionsRequest{},
	)
	if status.Code(err) != codes.Unavailable {
		t.Fatalf("missing authority error = %v", err)
	}
}

func planVersionService() *admin.Service {
	from := time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)
	repository := &planHandlerRepository{aggregate: catalogdomain.PlanVersionAggregate{
		ID:                 9,
		ProductCode:        "qhpro",
		ProductDisplayName: "QHPro",
		PlanCode:           "qhpro.pro",
		Version:            "1.0.0",
		DisplayName:        "Pro",
		Status:             catalogdomain.StatusDraft,
		SubjectScope:       catalogdomain.SubjectScopeProfile,
		EffectiveFrom:      &from,
		Entitlements: []catalogdomain.Entitlement{{
			Code:      "workspace.report.usage",
			Kind:      catalogdomain.EntitlementUsageAllowance,
			MeterCode: "workspace.report_generation.accepted",
			Amount:    10,
			Period:    catalogdomain.PeriodCalendarMonth,
		}},
		Operations: []catalogdomain.OperationBinding{{
			Code:           "workspace.report.generate",
			FeatureCode:    "workspace.report.generate",
			MeterCode:      "workspace.report_generation.accepted",
			UnitsPerAction: 1,
		}},
	}}
	return admin.NewService(repository, planHandlerPublisher{})
}

func trustedAdminContext(t *testing.T) context.Context {
	return trustedActorContext(t, "ROLE_ADMIN")
}

func trustedActorContext(t *testing.T, role string) context.Context {
	t.Helper()
	ctx, err := identity.BindServiceCaller(context.Background(), identity.ServiceCaller{ServiceID: "gateway-service"})
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = identity.BindCaller(ctx, identity.Caller{Kind: identity.CallerUser})
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = identity.BindActor(ctx, identity.Actor{
		AuthID:    7,
		ProfileID: 42,
		SessionID: 8,
		Role:      role,
		TokenType: "ACCESS",
	})
	if err != nil {
		t.Fatal(err)
	}
	return ctx
}

func handlerAggregateFromDraft(id uint64, record admin.DraftRecord) catalogdomain.PlanVersionAggregate {
	return catalogdomain.PlanVersionAggregate{
		ID: id, ProductCode: record.Terms.ProductCode, ProductDisplayName: record.ProductDisplayName,
		PlanCode: record.Terms.PlanCode, Version: record.Terms.Version, DisplayName: record.Terms.DisplayName,
		Status: catalogdomain.StatusDraft, SubjectScope: record.Terms.SubjectScope,
		SubscriptionTermDays: record.Terms.SubscriptionTermDays, TierRank: record.TierRank,
		TermsChecksum: record.TermsChecksum, Entitlements: record.Terms.Entitlements,
		Operations: record.Terms.Operations, Prices: record.Terms.Prices,
	}
}
