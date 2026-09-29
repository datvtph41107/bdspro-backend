package grpcserver

import (
	"time"

	usergrpc "payment/infra/client/user"
)

func newUserClient(address string, timeout time.Duration) (*usergrpc.Client, func(), error) {
	return usergrpc.New(address, timeout)
}
