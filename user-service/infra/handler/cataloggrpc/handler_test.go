package cataloggrpc

import (
	"context"
	"testing"
	"time"

	"common/callercontext"
	"common/identitycontext"
	"common/productcatalog"
	"common/transporttrust"
	userpb "pb/types/user"
	"user/internal/modules/catalog/admin"
	"user/internal/modules/catalog/domain"
	"user/internal/modules/catalog/publication"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handlerRepository struct {
	aggregate domain.PlanVersionAggregate
}

func (r *handlerRepository) GetPlanVersion(
	context.Context,
	uint64,
) (domain.PlanVersionAggregate, error) {
	return r.aggregate.Copy(), nil
}

func (r *handlerRepository) ListPlanVersions(
	context.Context,
	admin.Query,
) (admin.Page, error) {
	return admin.Page{PlanVersions: []domain.PlanVersionAggregate{r.aggregate}, Total: 1}, nil
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

type handlerPublisher struct{}

func (handlerPublisher) Publish(
	_ context.Context,
	command publication.Command,
) (publication.Result, error) {
	return publication.Result{PlanVersionID: command.PlanVersionID}, nil
}

func TestListPlanVersionsRequiresTrustedAdminAndMapsTerms(t *testing.T) {
	handler := NewHandler(catalogService(), allowPermissions(admin.PermissionPlanView))
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

func TestPublishPlanVersionValidatesEffectiveFrom(t *testing.T) {
	handler := NewHandler(catalogService(), allowPermissions(admin.PermissionPlanPublish))
	_, err := handler.PublishPlanVersion(trustedAdminContext(t), &userpb.PublishCatalogPlanVersionRequest{
		PlanVersionId: 9,
	})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("error = %v", err)
	}
}

func TestCatalogAuthorizationUsesPermissionNotRoleName(t *testing.T) {
	view := allowPermissions(admin.PermissionPlanView)
	handler := NewHandler(catalogService(), view)

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
	denied := NewHandler(catalogService(), allowPermissions())
	_, err := denied.ListPlanVersions(
		trustedActorContext(t, "ROLE_SUPER_ADMIN"),
		&userpb.ListCatalogPlanVersionsRequest{},
	)
	if status.Code(err) != codes.PermissionDenied {
		t.Fatalf("role-name bypass error = %v", err)
	}
}

func TestPublishRequiresPublishPermissionNotView(t *testing.T) {
	handler := NewHandler(catalogService(), allowPermissions(admin.PermissionPlanView))
	_, err := handler.PublishPlanVersion(
		trustedAdminContext(t),
		&userpb.PublishCatalogPlanVersionRequest{PlanVersionId: 9},
	)
	if status.Code(err) != codes.PermissionDenied {
		t.Fatalf("view-only publish error = %v", err)
	}
}

func TestCatalogAuthorizationFailsClosedWithoutAuthority(t *testing.T) {
	handler := NewHandler(catalogService(), nil)
	_, err := handler.ListPlanVersions(
		trustedAdminContext(t),
		&userpb.ListCatalogPlanVersionsRequest{},
	)
	if status.Code(err) != codes.Unavailable {
		t.Fatalf("missing authority error = %v", err)
	}
}

func catalogService() *admin.Service {
	from := time.Date(2026, 8, 15, 0, 0, 0, 0, time.UTC)
	repository := &handlerRepository{aggregate: domain.PlanVersionAggregate{
		ID:                 9,
		ProductCode:        "qhpro",
		ProductDisplayName: "QHPro",
		PlanCode:           "qhpro.pro",
		Version:            "1.0.0",
		DisplayName:        "Pro",
		Status:             productcatalog.StatusDraft,
		SubjectScope:       productcatalog.SubjectScopeProfile,
		EffectiveFrom:      &from,
		Entitlements: []productcatalog.Entitlement{{
			Code:      "workspace.report.usage",
			Kind:      productcatalog.EntitlementUsageAllowance,
			MeterCode: "workspace.report_generation.accepted",
			Amount:    10,
			Period:    productcatalog.PeriodCalendarMonth,
		}},
		Operations: []productcatalog.OperationPolicy{{
			Code:           "workspace.report.generate",
			FeatureCode:    "workspace.report.generate",
			MeterCode:      "workspace.report_generation.accepted",
			UnitsPerAction: 1,
		}},
	}}
	return admin.NewService(repository, handlerPublisher{})
}

func trustedAdminContext(t *testing.T) context.Context {
	return trustedActorContext(t, "ROLE_ADMIN")
}

func trustedActorContext(t *testing.T, role string) context.Context {
	t.Helper()
	ctx, err := transporttrust.BindCaller(context.Background(), transporttrust.Caller{
		ServiceID: "gateway-service",
		Level:     transporttrust.LevelVerified,
		IssuedAt:  time.Now().UTC(),
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = callercontext.BindCaller(ctx, callercontext.Caller{Kind: callercontext.KindUser})
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = identitycontext.BindActor(ctx, identitycontext.Actor{
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
