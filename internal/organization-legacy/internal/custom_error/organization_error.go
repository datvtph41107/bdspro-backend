package custom_error

import (
	_errors "common/errors"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/protoadapt"
)

var (
	InvalidRequestCode            codes.Code = 3001
	ExistingOrganizationNameCode  codes.Code = 3002
	InvalidTaxCodeCode            codes.Code = 3003
	UserAlreadyInOrganizationCode codes.Code = 3004
	InternalServerErrorCode       codes.Code = 3005
	InvalidOrganizationIDCode     codes.Code = 3006
	UserOnlyOneOrganizationCode   codes.Code = 3007
)

func InvalidRequest(details ...protoadapt.MessageV1) error {
	return _errors.ThrowError(int32(InvalidRequestCode), errors.New("invalid request"), details...)
}

func ExistingOrganizationName() error {
	return _errors.ThrowError(int32(ExistingOrganizationNameCode), errors.New("existing organization name"))
}

func InvalidTaxCode() error {
	return _errors.ThrowError(int32(InvalidTaxCodeCode), errors.New("invalid tax code"))
}

func UserAlreadyInOrganization() error {
	return _errors.ThrowError(int32(UserAlreadyInOrganizationCode), errors.New("user already in organization"))
}

func InternalServerError() error {
	return _errors.ThrowError(int32(InternalServerErrorCode), errors.New("internal server error"))
}

func InvalidOrganizationID() error {
	return _errors.ThrowError(int32(InvalidOrganizationIDCode), errors.New("invalid organization id"))
}

func Forbidden(message string) error {
	return _errors.ThrowError(int32(codes.PermissionDenied), errors.New(message))
}

func RecordNotFound(message string) error {
	return _errors.ThrowError(int32(codes.NotFound), errors.New(message))
}

func UserOnlyOneOrganization() error {
	return _errors.ThrowError(int32(UserOnlyOneOrganizationCode), errors.New("user only one organization"))
}
