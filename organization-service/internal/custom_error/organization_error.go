package custom_error

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/protoadapt"
)


var (
	InvalidRequestCode codes.Code = 3001
	ExistingOrganizationNameCode codes.Code = 3002
	InvalidTaxCodeCode codes.Code = 3003
	UserAlreadyInOrganizationCode codes.Code = 3004
	InternalServerErrorCode codes.Code = 3005
	InvalidOrganizationIDCode codes.Code = 3006
)

func InvalidRequest(details ...protoadapt.MessageV1 ) error {
	return NewError(int32(InvalidRequestCode), errors.New("invalid request"), details...)
}

func ExistingOrganizationName() error {
	return NewError(int32(ExistingOrganizationNameCode), errors.New("existing organization name"))
}

func InvalidTaxCode() error {
	return NewError(int32(InvalidTaxCodeCode), errors.New("invalid tax code"))
}

func UserAlreadyInOrganization() error {	
	return NewError(int32(UserAlreadyInOrganizationCode), errors.New("user already in organization"))
}

func InternalServerError() error {
	return NewError(int32(InternalServerErrorCode), errors.New("internal server error"))
}

func InvalidOrganizationID() error {
	return NewError(int32(InvalidOrganizationIDCode), errors.New("invalid organization id"))
}

func Forbidden(message string) error {
	return NewError(int32(codes.PermissionDenied), errors.New(message))
}

func RecordNotFound(message string) error {
	return NewError(int32(codes.NotFound), errors.New(message))
}