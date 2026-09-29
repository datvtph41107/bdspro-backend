package repository

// type DealInvitationRepository interface {
// 	Create(ctx context.Context, invitation *entity.DealMember) (*entity.DealMember, error)
// 	Update(ctx context.Context, invitation *entity.DealMember) (*entity.DealMember, error)
// 	GetByID(ctx context.Context, id uint64) (*entity.DealMember, error)
// 	GetByDealIDAndInviteeID(ctx context.Context, dealID, inviteeID uint64) (*entity.DealMember, error)
// 	GetByDealID(ctx context.Context, dealID uint64, page, size int) ([]*entity.DealMember, uint32, error)
// 	GetPendingByInviteeID(ctx context.Context, inviteeID uint64, page, size int) ([]*entity.DealMember, uint32, error)
// 	GetAcceptedByDealID(ctx context.Context, dealID uint64) ([]*entity.DealMember, error)
// 	CountByDealID(ctx context.Context, dealID uint64) (uint32, error)
// 	CountPendingByInviteeID(ctx context.Context, inviteeID uint64) (uint32, error)
// 	Delete(ctx context.Context, id uint64) error
// 	ExistsByDealIDAndInviteeID(ctx context.Context, dealID, inviteeID uint64) (bool, error)
// 	GetByIds(ctx context.Context, dealId uint64, ids []uint64) ([]*entity.DealMember, error)
// }
