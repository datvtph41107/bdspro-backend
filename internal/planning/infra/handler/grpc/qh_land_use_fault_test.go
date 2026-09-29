package handler_grpc

import (
	"context"
	"errors"
	"strings"
	"testing"

	_dto "common/domain/dto"
	_errors "common/errors"

	tqdpb "pb/types/tqd"
	"tqd/internal"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/usecase"
)

type qhLandUseUsecaseFaultStub struct {
	usecase.QHLayerLandUseGroupUsecase
	listErr error
}

func (s *qhLandUseUsecaseFaultStub) List(
	context.Context,
	*_dto.Pagable,
	*uint64,
	*uint64,
) ([]qh_domain.QHLandUse, int64, error) {
	return nil, 0, s.listErr
}

func TestQHLandUseConflictMapsToAlreadyExists(t *testing.T) {
	err := _errors.ReturnError(
		service.LandUseDuplicate,
		_errors.WithPublicMessage("layer 7 already has land use 11"),
	)
	assertCanonicalStatus(t, mapLandUseGroupError(err), service.LandUseDuplicate)
}

func TestQHLandUseValidationMapsToInvalidArgument(t *testing.T) {
	err := _errors.ReturnError(service.LandUseNameRequired)
	assertCanonicalStatus(t, mapLandUseGroupError(err), service.LandUseNameRequired)
}

func TestQHLandUseNotFoundMapsToNotFound(t *testing.T) {
	err := _errors.ReturnError(
		service.LandUseNotFound,
		_errors.WithPublicMessage("land use 11 not found"),
	)
	assertCanonicalStatus(t, mapLandUseGroupError(err), service.LandUseNotFound)
}

func TestQHLandUseListDependencyFailureDoesNotLeak(t *testing.T) {
	dependencyErr := errors.New("postgres password=secret land-use list failed")
	h := &QHLandUseGrpcHandler{
		uc: &qhLandUseUsecaseFaultStub{listErr: dependencyErr},
	}

	_, err := h.ListLandUses(context.Background(), &tqdpb.ListLandUseRequest{})
	st := assertTechnicalStatus(t, err)
	if strings.Contains(st.Message(), "secret") || strings.Contains(st.Message(), "postgres") {
		t.Fatalf("dependency detail leaked: %q", st.Message())
	}
}
