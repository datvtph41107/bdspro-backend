package postgres

import (
	"context"
	"mdn4/modules/user/domain/user"

	"gorm.io/gorm"
)

// @bind: mdn4/internal/interface/repo.IDemoProductLogRepo
type DemoProductLogPostgres struct {
	DB *gorm.DB
}

func NewDemoProductLogPostgres(db *gorm.DB) *DemoProductLogPostgres {
	return &DemoProductLogPostgres{DB: db}
}

func (r *DemoProductLogPostgres) CreateBatch(ctx context.Context, logs []*user.DemoProductLog) error {
	if len(logs) == 0 {
		return nil
	}
	return r.DB.WithContext(ctx).Create(&logs).Error
}

func (r *DemoProductLogPostgres) ListProductIDsByProfileID(ctx context.Context, profileID uint64) ([]uint64, error) {
	var ids []uint64
	err := r.DB.WithContext(ctx).
		Model(&user.DemoProductLog{}).
		Where("profile_id = ?", profileID).
		Order("id ASC").
		Pluck("product_id", &ids).Error
	return ids, err
}

func (r *DemoProductLogPostgres) DeleteByProfileID(ctx context.Context, profileID uint64) error {
	return r.DB.WithContext(ctx).
		Where("profile_id = ?", profileID).
		Delete(&user.DemoProductLog{}).Error
}

func (r *DemoProductLogPostgres) DeleteByProfileIDAndProductIDs(ctx context.Context, profileID uint64, productIDs []uint64) error {
	if len(productIDs) == 0 {
		return nil
	}
	return r.DB.WithContext(ctx).
		Where("profile_id = ? AND product_id IN ?", profileID, productIDs).
		Delete(&user.DemoProductLog{}).Error
}

func (r *DemoProductLogPostgres) CreateContactBatch(ctx context.Context, logs []*user.DemoContactLog) error {
	if len(logs) == 0 {
		return nil
	}
	return r.DB.WithContext(ctx).Create(&logs).Error
}

func (r *DemoProductLogPostgres) ListContactIDsByProfileID(ctx context.Context, profileID uint64) ([]uint64, error) {
	var ids []uint64
	err := r.DB.WithContext(ctx).
		Model(&user.DemoContactLog{}).
		Where("profile_id = ?", profileID).
		Order("id ASC").
		Pluck("contact_id", &ids).Error
	return ids, err
}

func (r *DemoProductLogPostgres) DeleteContactsByProfileID(ctx context.Context, profileID uint64) error {
	return r.DB.WithContext(ctx).
		Where("profile_id = ?", profileID).
		Delete(&user.DemoContactLog{}).Error
}

func (r *DemoProductLogPostgres) DeleteContactsByProfileIDAndContactIDs(ctx context.Context, profileID uint64, contactIDs []uint64) error {
	if len(contactIDs) == 0 {
		return nil
	}
	return r.DB.WithContext(ctx).
		Where("profile_id = ? AND contact_id IN ?", profileID, contactIDs).
		Delete(&user.DemoContactLog{}).Error
}

func (r *DemoProductLogPostgres) CreateAppointmentBatch(ctx context.Context, logs []*user.DemoAppointmentLog) error {
	if len(logs) == 0 {
		return nil
	}
	return r.DB.WithContext(ctx).Create(&logs).Error
}

func (r *DemoProductLogPostgres) ListAppointmentIDsByProfileID(ctx context.Context, profileID uint64) ([]uint64, error) {
	var ids []uint64
	err := r.DB.WithContext(ctx).
		Model(&user.DemoAppointmentLog{}).
		Where("profile_id = ?", profileID).
		Order("id ASC").
		Pluck("appointment_id", &ids).Error
	return ids, err
}

func (r *DemoProductLogPostgres) DeleteAppointmentsByProfileID(ctx context.Context, profileID uint64) error {
	return r.DB.WithContext(ctx).
		Where("profile_id = ?", profileID).
		Delete(&user.DemoAppointmentLog{}).Error
}

func (r *DemoProductLogPostgres) DeleteAppointmentsByProfileIDAndIDs(ctx context.Context, profileID uint64, appointmentIDs []uint64) error {
	if len(appointmentIDs) == 0 {
		return nil
	}
	return r.DB.WithContext(ctx).
		Where("profile_id = ? AND appointment_id IN ?", profileID, appointmentIDs).
		Delete(&user.DemoAppointmentLog{}).Error
}

func (r *DemoProductLogPostgres) CreateDealBatch(ctx context.Context, logs []*user.DemoDealLog) error {
	if len(logs) == 0 {
		return nil
	}
	return r.DB.WithContext(ctx).Create(&logs).Error
}

func (r *DemoProductLogPostgres) ListDealIDsByProfileID(ctx context.Context, profileID uint64) ([]uint64, error) {
	var ids []uint64
	err := r.DB.WithContext(ctx).
		Model(&user.DemoDealLog{}).
		Where("profile_id = ?", profileID).
		Order("id ASC").
		Pluck("deal_id", &ids).Error
	return ids, err
}

func (r *DemoProductLogPostgres) DeleteDealsByProfileID(ctx context.Context, profileID uint64) error {
	return r.DB.WithContext(ctx).
		Where("profile_id = ?", profileID).
		Delete(&user.DemoDealLog{}).Error
}

func (r *DemoProductLogPostgres) DeleteDealsByProfileIDAndIDs(ctx context.Context, profileID uint64, dealIDs []uint64) error {
	if len(dealIDs) == 0 {
		return nil
	}
	return r.DB.WithContext(ctx).
		Where("profile_id = ? AND deal_id IN ?", profileID, dealIDs).
		Delete(&user.DemoDealLog{}).Error
}
