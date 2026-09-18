package custom_error

import (
	_errors "common/errors"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/protoadapt"
)

func InvalidRequest(details ...protoadapt.MessageV1) error {
	return organizationError(codes.InvalidArgument, "ORGANIZATION_REQUEST_INVALID", errors.New("invalid request"), details...)
}

func ExistingOrganizationName() error {
	return organizationError(codes.AlreadyExists, "ORGANIZATION_NAME_EXISTS", errors.New("existing organization name"))
}

func InvalidTaxCode() error {
	return organizationError(codes.InvalidArgument, "ORGANIZATION_TAX_CODE_INVALID", errors.New("invalid tax code"))
}

func UserAlreadyInOrganization() error {
	return organizationError(codes.AlreadyExists, "ORGANIZATION_MEMBER_EXISTS", errors.New("user already in organization"))
}

func InternalServerError() error {
	return organizationError(codes.Internal, "ORGANIZATION_INTERNAL_ERROR", errors.New("internal server error"))
}

func InvalidOrganizationID() error {
	return organizationError(codes.InvalidArgument, "ORGANIZATION_ID_INVALID", errors.New("invalid organization id"))
}

func Forbidden(message string) error {
	return _errors.ThrowError(int32(codes.PermissionDenied), errors.New(message))
}

func RecordNotFound(message string) error {
	return _errors.ThrowError(int32(codes.NotFound), errors.New(message))
}

func UserOnlyOneOrganization() error {
	return organizationError(codes.FailedPrecondition, "ORGANIZATION_MEMBERSHIP_LIMIT", errors.New("user only one organization"))
}
