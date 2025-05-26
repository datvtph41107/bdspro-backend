package logic_owner

import (
	logic_crud "common/logic/crud"
	"context"

	"google.golang.org/protobuf/types/known/anypb"
)

type OwnerHandler[T any, R logic_crud.IRepo[T]] struct {
	UC *OwnerUsecase[T, R]
}

func (*OwnerHandler[T, R]) RegisterHandler(c context.Context, body any) (*anypb.Any, error) {
	// switch expression {
	// case condition:

	// }
	return nil, nil
}
