package usecase

import (
	_errors "common/errors"
	"context"
	"errors"
	"strings"
	"time"

	seo_strategy "crm/internal/domain/seo_strategy"
	"crm/internal/dto"
	"crm/internal/repo"
)

type seoQueryClusterUserNeedReader interface {
	GetByID(ctx context.Context, id uint64) (*seo_strategy.SeoUserNeed, error)
}

type seoQueryClusterIntentReader interface {
	GetByID(ctx context.Context, id uint64) (*seo_strategy.SeoIntentHypothesis, error)
	ListLinkedEvidence(ctx context.Context, intentID uint64) ([]seo_strategy.LinkedIntentQueryEvidence, error)
}

type seoQueryClusterEvidenceReader interface {
	GetByID(ctx context.Context, id uint64) (*seo_strategy.SeoQueryEvidence, error)
}

type SeoQueryClusterUsecase struct {
	repository          repo.SeoQueryClusterRepo
	userNeedReader      seoQueryClusterUserNeedReader
	intentReader        seoQueryClusterIntentReader
	queryEvidenceReader seoQueryClusterEvidenceReader
	now                 func() time.Time
}

func NewSeoQueryClusterUsecase(
	repository repo.SeoQueryClusterRepo,
	userNeedRepository repo.SeoUserNeedRepo,
	intentRepository repo.SeoIntentHypothesisRepo,
	queryEvidenceRepository repo.SeoQueryEvidenceRepo,
) *SeoQueryClusterUsecase {
	return &SeoQueryClusterUsecase{
		repository:          repository,
		userNeedReader:      userNeedRepository,
		intentReader:        intentRepository,
		queryEvidenceReader: queryEvidenceRepository,
		now:                 time.Now,
	}
}

func (u *SeoQueryClusterUsecase) currentTime() time.Time {
	if u.now == nil {
		return time.Now().UTC()
	}
	return u.now().UTC()
}

func (u *SeoQueryClusterUsecase) Create(
	ctx context.Context,
	req *dto.CreateSeoQueryClusterRequest,
) (*seo_strategy.SeoQueryClusterAggregate, error) {
	if req == nil || req.IntentHypothesisID == 0 {
		return nil, _errors.ReturnError(400, "intentHypothesisId không hợp lệ")
	}

	intent, err := u.intentReader.GetByID(ctx, req.IntentHypothesisID)
	if err != nil {
		return nil, mapQueryClusterIntentError(err)
	}
	if intent == nil {
		return nil, _errors.ReturnError(404, "seo intent hypothesis không tồn tại")
	}
	if intent.Status != seo_strategy.IntentHypothesisStatusActive {
		return nil, _errors.ReturnError(409, "intent hypothesis phải active trước khi tạo query cluster")
	}

	need, err := u.userNeedReader.GetByID(ctx, intent.UserNeedID)
	if err != nil {
		return nil, mapQueryClusterUserNeedError(err)
	}
	if need == nil {
		return nil, _errors.ReturnError(404, "seo user need không tồn tại")
	}
	if need.Status != seo_strategy.UserNeedStatusActive {
		return nil, _errors.ReturnError(409, "user need phải active trước khi tạo query cluster")
	}
	if need.ObjectiveID != intent.ObjectiveID {
		return nil, _errors.ReturnError(409, "intent hypothesis không thuộc đúng objective của user need")
	}

	cluster, err := seo_strategy.CreateSeoQueryClusterDraft(
		intent.ObjectiveID,
		intent.UserNeedID,
		intent.ID,
		req.Code,
		req.Name,
		u.currentTime(),
	)
	if err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}

	existing, err := u.repository.GetByCode(ctx, cluster.UserNeedID, cluster.Code)
	if err == nil && existing != nil {
		return nil, _errors.ReturnError(409, "query cluster code đã tồn tại trong user need")
	}
	if err != nil && !errors.Is(err, repo.ErrSeoQueryClusterNotFound) {
		return nil, mapQueryClusterRepoError(err)
	}

	if err := u.repository.Create(ctx, cluster); err != nil {
		return nil, mapQueryClusterRepoError(err)
	}
	return u.Get(ctx, cluster.ID)
}

func (u *SeoQueryClusterUsecase) Update(
	ctx context.Context,
	req *dto.UpdateSeoQueryClusterRequest,
) (*seo_strategy.SeoQueryClusterAggregate, error) {
	if req == nil || req.ID == 0 || req.ExpectedVersion == 0 {
		return nil, _errors.ReturnError(400, "query cluster update request không hợp lệ")
	}

	cluster, err := u.repository.GetByID(ctx, req.ID)
	if err != nil {
		return nil, mapQueryClusterRepoError(err)
	}
	if cluster.Status == seo_strategy.QueryClusterStatusArchived {
		return nil, _errors.ReturnError(409, "không thể sửa query cluster đã archived")
	}

	confidence := seo_strategy.ObjectiveConfidence(strings.TrimSpace(req.Confidence))
	if confidence == "" {
		confidence = seo_strategy.ObjectiveConfidenceExploratory
	}
	if !seo_strategy.IsValidObjectiveConfidence(confidence) {
		return nil, _errors.ReturnError(400, "query cluster confidence không hợp lệ")
	}

	cluster.Name = strings.TrimSpace(req.Name)
	cluster.SearchPromise = strings.TrimSpace(req.SearchPromise)
	cluster.InclusionRule = strings.TrimSpace(req.InclusionRule)
	cluster.ExclusionRule = strings.TrimSpace(req.ExclusionRule)
	cluster.Confidence = confidence
	cluster.OwnerCode = strings.TrimSpace(req.OwnerCode)
	cluster.ReviewDueAt = req.ReviewDueAt
	cluster.UpdatedAt = u.currentTime()

	if err := seo_strategy.ValidateQueryClusterDefinition(cluster); err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}
	if err := u.repository.UpdateDefinition(ctx, cluster, req.ExpectedVersion); err != nil {
		return nil, mapQueryClusterRepoError(err)
	}
	return u.Get(ctx, cluster.ID)
}

func (u *SeoQueryClusterUsecase) Get(
	ctx context.Context,
	id uint64,
) (*seo_strategy.SeoQueryClusterAggregate, error) {
	if id == 0 {
		return nil, _errors.ReturnError(400, "query cluster id không hợp lệ")
	}

	cluster, err := u.repository.GetByID(ctx, id)
	if err != nil {
		return nil, mapQueryClusterRepoError(err)
	}
	intent, err := u.intentReader.GetByID(ctx, cluster.IntentHypothesisID)
	if err != nil {
		return nil, mapQueryClusterIntentError(err)
	}
	if intent == nil {
		return nil, _errors.ReturnError(404, "intent hypothesis của query cluster không tồn tại")
	}
	need, err := u.userNeedReader.GetByID(ctx, cluster.UserNeedID)
	if err != nil {
		return nil, mapQueryClusterUserNeedError(err)
	}
	if need == nil {
		return nil, _errors.ReturnError(404, "user need của query cluster không tồn tại")
	}
	members, err := u.repository.ListMembers(ctx, cluster.ID)
	if err != nil {
		return nil, mapQueryClusterRepoError(err)
	}
	intentEvidence, err := u.intentReader.ListLinkedEvidence(ctx, intent.ID)
	if err != nil {
		return nil, mapQueryClusterIntentError(err)
	}

	return &seo_strategy.SeoQueryClusterAggregate{
		Cluster:  cluster,
		Intent:   intent,
		UserNeed: need,
		Members:  members,
		Readiness: seo_strategy.EvaluateQueryClusterReadiness(
			need,
			intent,
			cluster,
			members,
			intentEvidence,
			u.currentTime(),
		),
	}, nil
}

func (u *SeoQueryClusterUsecase) List(
	ctx context.Context,
	req *dto.SeoQueryClusterListRequest,
) ([]seo_strategy.SeoQueryCluster, int64, error) {
	if req == nil {
		req = &dto.SeoQueryClusterListRequest{}
	}
	filter := repo.SeoQueryClusterFilter{
		ObjectiveID:        req.ObjectiveID,
		UserNeedID:         req.UserNeedID,
		IntentHypothesisID: req.IntentHypothesisID,
		Keyword:            strings.TrimSpace(req.Keyword),
		Page:               req.Page,
		Size:               req.Size,
	}
	if value := strings.TrimSpace(req.Status); value != "" {
		parsed, err := seo_strategy.ParseQueryClusterStatus(value)
		if err != nil {
			return nil, 0, _errors.ReturnError(400, err.Error())
		}
		filter.Status = parsed
	}
	if value := strings.TrimSpace(req.Confidence); value != "" {
		confidence := seo_strategy.ObjectiveConfidence(value)
		if !seo_strategy.IsValidObjectiveConfidence(confidence) {
			return nil, 0, _errors.ReturnError(400, "confidence không hợp lệ")
		}
		filter.Confidence = confidence
	}

	items, total, err := u.repository.List(ctx, filter)
	if err != nil {
		return nil, 0, mapQueryClusterRepoError(err)
	}
	return items, total, nil
}

func (u *SeoQueryClusterUsecase) UpsertMember(
	ctx context.Context,
	req *dto.UpsertSeoQueryClusterMemberRequest,
) (*seo_strategy.SeoQueryClusterAggregate, error) {
	if req == nil || req.ClusterID == 0 || req.QueryEvidenceID == 0 || req.ExpectedVersion == 0 {
		return nil, _errors.ReturnError(400, "query cluster member request không hợp lệ")
	}

	cluster, err := u.repository.GetByID(ctx, req.ClusterID)
	if err != nil {
		return nil, mapQueryClusterRepoError(err)
	}
	if cluster.LockVersion != req.ExpectedVersion {
		return nil, _errors.ReturnError(409, "query cluster đã được cập nhật; hãy tải lại dữ liệu")
	}
	if cluster.Status == seo_strategy.QueryClusterStatusArchived {
		return nil, _errors.ReturnError(409, "không thể sửa membership của query cluster đã archived")
	}

	evidence, err := u.queryEvidenceReader.GetByID(ctx, req.QueryEvidenceID)
	if err != nil {
		if errors.Is(err, repo.ErrSeoQueryEvidenceNotFound) {
			return nil, _errors.ReturnError(404, "seo query evidence không tồn tại")
		}
		return nil, err
	}
	if evidence == nil {
		return nil, _errors.ReturnError(404, "seo query evidence không tồn tại")
	}
	if evidence.ObjectiveID != cluster.ObjectiveID || evidence.UserNeedID != cluster.UserNeedID {
		return nil, _errors.ReturnError(409, "query cluster và query evidence không cùng user need/objective")
	}
	if !seo_strategy.IsQueryEvidenceCurrent(*evidence, u.currentTime()) {
		return nil, _errors.ReturnError(409, "chỉ Query Evidence active và còn hiệu lực mới được thêm vào cluster")
	}

	membershipType, err := seo_strategy.ParseQueryClusterMembershipType(req.MembershipType)
	if err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}
	intentEvidence, err := u.intentReader.ListLinkedEvidence(ctx, cluster.IntentHypothesisID)
	if err != nil {
		return nil, mapQueryClusterIntentError(err)
	}
	relation, ok := findIntentEvidenceRelation(intentEvidence, evidence.ID)
	if !ok {
		return nil, _errors.ReturnError(409, "query evidence phải được liên kết với intent hypothesis trước khi thêm vào cluster")
	}
	if err := seo_strategy.ValidateQueryClusterMembershipAgainstIntent(membershipType, relation); err != nil {
		return nil, _errors.ReturnError(422, err.Error())
	}

	now := u.currentTime()
	member := &seo_strategy.SeoQueryClusterMember{
		ObjectiveID:     cluster.ObjectiveID,
		UserNeedID:      cluster.UserNeedID,
		ClusterID:       cluster.ID,
		QueryEvidenceID: evidence.ID,
		MembershipType:  membershipType,
		Rationale:       strings.TrimSpace(req.Rationale),
		Status:          seo_strategy.QueryClusterMemberStatusActive,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if err := seo_strategy.ValidateQueryClusterMember(member); err != nil {
		return nil, _errors.ReturnError(400, err.Error())
	}
	if err := u.repository.UpsertMember(ctx, member, req.ExpectedVersion); err != nil {
		return nil, mapQueryClusterRepoError(err)
	}
	return u.Get(ctx, cluster.ID)
}

func (u *SeoQueryClusterUsecase) RemoveMember(
	ctx context.Context,
	req *dto.RemoveSeoQueryClusterMemberRequest,
) (*seo_strategy.SeoQueryClusterAggregate, error) {
	if req == nil || req.ClusterID == 0 || req.QueryEvidenceID == 0 || req.ExpectedVersion == 0 {
		return nil, _errors.ReturnError(400, "query cluster remove member request không hợp lệ")
	}
	cluster, err := u.repository.GetByID(ctx, req.ClusterID)
	if err != nil {
		return nil, mapQueryClusterRepoError(err)
	}
	if cluster.LockVersion != req.ExpectedVersion {
		return nil, _errors.ReturnError(409, "query cluster đã được cập nhật; hãy tải lại dữ liệu")
	}
	if cluster.Status == seo_strategy.QueryClusterStatusArchived {
		return nil, _errors.ReturnError(409, "không thể sửa membership của query cluster đã archived")
	}
	if err := u.repository.RemoveMember(ctx, cluster.ID, req.QueryEvidenceID, req.ExpectedVersion); err != nil {
		return nil, mapQueryClusterRepoError(err)
	}
	return u.Get(ctx, cluster.ID)
}

func (u *SeoQueryClusterUsecase) Activate(
	ctx context.Context,
	req *dto.SeoQueryClusterCommandRequest,
) (*seo_strategy.SeoQueryClusterAggregate, error) {
	if req == nil || req.ID == 0 || req.ExpectedVersion == 0 {
		return nil, _errors.ReturnError(400, "query cluster command không hợp lệ")
	}
	aggregate, err := u.Get(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if aggregate.Cluster.LockVersion != req.ExpectedVersion {
		return nil, _errors.ReturnError(409, "query cluster đã được cập nhật; hãy tải lại dữ liệu")
	}
	if err := seo_strategy.ValidateQueryClusterStatusTransition(
		aggregate.Cluster.Status,
		seo_strategy.QueryClusterStatusActive,
	); err != nil {
		return nil, _errors.ReturnError(409, err.Error())
	}
	if !aggregate.Readiness.Ready {
		codes := make([]string, 0, len(aggregate.Readiness.Blockers))
		for _, item := range aggregate.Readiness.Blockers {
			codes = append(codes, item.Code)
		}
		return nil, _errors.ReturnError(422, "query cluster chưa sẵn sàng: "+strings.Join(codes, ", "))
	}
	if err := u.repository.UpdateStatus(
		ctx,
		aggregate.Cluster.ID,
		seo_strategy.QueryClusterStatusActive,
		req.ExpectedVersion,
	); err != nil {
		return nil, mapQueryClusterRepoError(err)
	}
	return u.Get(ctx, aggregate.Cluster.ID)
}

func (u *SeoQueryClusterUsecase) Pause(
	ctx context.Context,
	req *dto.SeoQueryClusterCommandRequest,
) (*seo_strategy.SeoQueryClusterAggregate, error) {
	return u.changeStatus(ctx, req, seo_strategy.QueryClusterStatusPaused)
}

func (u *SeoQueryClusterUsecase) Archive(
	ctx context.Context,
	req *dto.SeoQueryClusterCommandRequest,
) (*seo_strategy.SeoQueryClusterAggregate, error) {
	return u.changeStatus(ctx, req, seo_strategy.QueryClusterStatusArchived)
}

func (u *SeoQueryClusterUsecase) changeStatus(
	ctx context.Context,
	req *dto.SeoQueryClusterCommandRequest,
	target seo_strategy.QueryClusterStatus,
) (*seo_strategy.SeoQueryClusterAggregate, error) {
	if req == nil || req.ID == 0 || req.ExpectedVersion == 0 {
		return nil, _errors.ReturnError(400, "query cluster command không hợp lệ")
	}
	current, err := u.repository.GetByID(ctx, req.ID)
	if err != nil {
		return nil, mapQueryClusterRepoError(err)
	}
	if current.LockVersion != req.ExpectedVersion {
		return nil, _errors.ReturnError(409, "query cluster đã được cập nhật; hãy tải lại dữ liệu")
	}
	if err := seo_strategy.ValidateQueryClusterStatusTransition(current.Status, target); err != nil {
		return nil, _errors.ReturnError(409, err.Error())
	}
	if err := u.repository.UpdateStatus(ctx, current.ID, target, req.ExpectedVersion); err != nil {
		return nil, mapQueryClusterRepoError(err)
	}
	return u.Get(ctx, current.ID)
}

func findIntentEvidenceRelation(
	items []seo_strategy.LinkedIntentQueryEvidence,
	queryEvidenceID uint64,
) (seo_strategy.IntentEvidenceRelation, bool) {
	for _, item := range items {
		if item.Evidence.ID == queryEvidenceID && item.Evidence.Status == seo_strategy.QueryEvidenceStatusActive {
			return item.Link.RelationType, true
		}
	}
	return "", false
}

func mapQueryClusterRepoError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, repo.ErrSeoQueryClusterNotFound):
		return _errors.ReturnError(404, "seo query cluster không tồn tại")
	case errors.Is(err, repo.ErrSeoQueryClusterConflict):
		return _errors.ReturnError(409, "query cluster đã được cập nhật; hãy tải lại dữ liệu")
	case errors.Is(err, repo.ErrSeoQueryClusterInvalid):
		return _errors.ReturnError(400, err.Error())
	case errors.Is(err, repo.ErrSeoQueryClusterMemberNotFound):
		return _errors.ReturnError(404, "query evidence chưa là thành viên active của query cluster")
	case errors.Is(err, repo.ErrSeoQueryClusterMembershipConflict):
		return _errors.ReturnError(409, "membership xung đột: cluster chỉ có một primary query và query chỉ thuộc một active answer set")
	case errors.Is(err, repo.ErrSeoQueryClusterMemberMismatch):
		return _errors.ReturnError(409, "query cluster và query evidence không cùng user need/objective")
	default:
		return err
	}
}

func mapQueryClusterIntentError(err error) error {
	if errors.Is(err, repo.ErrSeoIntentHypothesisNotFound) {
		return _errors.ReturnError(404, "seo intent hypothesis không tồn tại")
	}
	return err
}

func mapQueryClusterUserNeedError(err error) error {
	if errors.Is(err, repo.ErrSeoUserNeedNotFound) {
		return _errors.ReturnError(404, "seo user need không tồn tại")
	}
	return err
}
