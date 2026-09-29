package postgres

// type landUseRepo struct {
// 	db *gorm.DB
// }

// func NewLandUseRepository(db *gorm.DB) repo.LandUseRepository {
// 	return &landUseRepo{db: db}
// }

// func (r *landUseRepo) Create(ctx context.Context, row *qh_domain.LandUse) error {
// 	return r.db.WithContext(ctx).Create(row).Error
// }

// func (r *landUseRepo) GetByID(ctx context.Context, id uint64) (*qh_domain.LandUse, error) {
// 	if id == 0 {
// 		return nil, nil
// 	}
// 	var row qh_domain.LandUse
// 	err := r.db.WithContext(ctx).
// 		Where("id = ? AND deleted_at IS NULL", id).
// 		First(&row).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, nil
// 		}
// 		return nil, fmt.Errorf("get land_use by id: %w", err)
// 	}
// 	return &row, nil
// }

// func (r *landUseRepo) GetByCode(ctx context.Context, code string) (*qh_domain.LandUse, error) {
// 	if code == "" {
// 		return nil, nil
// 	}
// 	var row qh_domain.LandUse
// 	err := r.db.WithContext(ctx).
// 		Where("code = ? AND deleted_at IS NULL", code).
// 		First(&row).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, nil
// 		}
// 		return nil, fmt.Errorf("get land_use by code: %w", err)
// 	}
// 	return &row, nil
// }
