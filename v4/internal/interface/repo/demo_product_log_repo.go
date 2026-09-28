package repo

import (
	"context"
	"userv4/internal/domain/user"
)

type IDemoProductLogRepo interface {
	CreateBatch(ctx context.Context, logs []*user.DemoProductLog) error
	ListProductIDsByProfileID(ctx context.Context, profileID uint64) ([]uint64, error)
	DeleteByProfileID(ctx context.Context, profileID uint64) error
	DeleteByProfileIDAndProductIDs(ctx context.Context, profileID uint64, productIDs []uint64) error

	CreateContactBatch(ctx context.Context, logs []*user.DemoContactLog) error
	ListContactIDsByProfileID(ctx context.Context, profileID uint64) ([]uint64, error)
	DeleteContactsByProfileID(ctx context.Context, profileID uint64) error
	DeleteContactsByProfileIDAndContactIDs(ctx context.Context, profileID uint64, contactIDs []uint64) error

	CreateAppointmentBatch(ctx context.Context, logs []*user.DemoAppointmentLog) error
	ListAppointmentIDsByProfileID(ctx context.Context, profileID uint64) ([]uint64, error)
	DeleteAppointmentsByProfileID(ctx context.Context, profileID uint64) error
	DeleteAppointmentsByProfileIDAndIDs(ctx context.Context, profileID uint64, appointmentIDs []uint64) error

	CreateDealBatch(ctx context.Context, logs []*user.DemoDealLog) error
	ListDealIDsByProfileID(ctx context.Context, profileID uint64) ([]uint64, error)
	DeleteDealsByProfileID(ctx context.Context, profileID uint64) error
	DeleteDealsByProfileIDAndIDs(ctx context.Context, profileID uint64, dealIDs []uint64) error
}
