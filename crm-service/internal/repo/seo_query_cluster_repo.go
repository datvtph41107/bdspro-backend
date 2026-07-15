package repo

import (
	"context"
	"errors"

	seo_strategy "crm/internal/domain/seo_strategy"
)

var (
	ErrSeoQueryClusterNotFound           = errors.New("seo query cluster không tồn tại")
	ErrSeoQueryClusterConflict           = errors.New("seo query cluster đã được cập nhật")
	ErrSeoQueryClusterInvalid            = errors.New("seo query cluster không hợp lệ")
	ErrSeoQueryClusterMemberNotFound     = errors.New("query evidence chưa là thành viên của query cluster")
	ErrSeoQueryClusterMembershipConflict = errors.New("query cluster membership xung đột")
	ErrSeoQueryClusterMemberMismatch     = errors.New("query cluster và query evidence không cùng user need/objective")
)

type SeoQueryClusterFilter struct {
	ObjectiveID        uint64
	UserNeedID         uint64
	IntentHypothesisID uint64

	Keyword    string
	Status     seo_strategy.QueryClusterStatus
	Confidence seo_strategy.ObjectiveConfidence

	Page uint32
	Size uint32
}

type SeoQueryClusterRepo interface {
	Create(ctx context.Context, cluster *seo_strategy.SeoQueryCluster) error

	GetByID(ctx context.Context, id uint64) (*seo_strategy.SeoQueryCluster, error)

	GetByCode(
		ctx context.Context,
		userNeedID uint64,
		code string,
	) (*seo_strategy.SeoQueryCluster, error)

	List(
		ctx context.Context,
		filter SeoQueryClusterFilter,
	) ([]seo_strategy.SeoQueryCluster, int64, error)

	UpdateDefinition(
		ctx context.Context,
		cluster *seo_strategy.SeoQueryCluster,
		expectedVersion uint32,
	) error

	UpdateStatus(
		ctx context.Context,
		id uint64,
		status seo_strategy.QueryClusterStatus,
		expectedVersion uint32,
	) error

	UpsertMember(
		ctx context.Context,
		member *seo_strategy.SeoQueryClusterMember,
		expectedVersion uint32,
	) error

	RemoveMember(
		ctx context.Context,
		clusterID uint64,
		queryEvidenceID uint64,
		expectedVersion uint32,
	) error

	ListMembers(
		ctx context.Context,
		clusterID uint64,
	) ([]seo_strategy.LinkedQueryClusterMember, error)
}
