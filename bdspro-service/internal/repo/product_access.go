package repo

// type ProductAccessRepo interface {
// 	ExistsByUserAndProduct(userId, productId uint64) (bool, error)
// 	GetByUserAndProduct(userId, productId uint64) (*domain.SharingAccess, error)
// 	GetByTargetAndProduct(userId, target, productId uint64) (*domain.SharingAccess, error)
// 	GetAll() ([]domain.SharingAccess, error)
// 	OwnerSearchProduct(c context.Context, id uint64, dto *dto.ProductAccessSearch) (*[]domain.ProductAccessQuery, error)
// 	UpdateCommission(c context.Context,
// 		profileId uint64,
// 		dto dto.ProductAccessUpdate,
// 	) error
// 	GetByID(id uint64) (*domain.SharingAccess, error)
// 	Delete(c context.Context, id uint64) error
// 	// GetProductAccessByOrgID(ctx context.Context, orgID int64, dto dto.ProductAccessSearch) ([]*domain.ProductAccess, error)
// 	// BulkSave(c context.Context, id uint64, dto dto.ProductAccessBulkRequest) (*[]domain.ProductAccess, error)
// }
