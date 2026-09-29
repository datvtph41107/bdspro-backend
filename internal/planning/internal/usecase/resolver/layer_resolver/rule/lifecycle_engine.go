package rule

import (
	"context"
	"fmt"
	"sort"
	"time"

	"tqd/internal/usecase/resolver/layer_resolver/types"
)

// LayerRepository — interface tối thiểu để LifecycleEngine query layer
type LayerRepository interface {
	GetByID(ctx context.Context, layerID uint64) (*LayerRecord, error)
	GetByFamilyID(ctx context.Context, familyID uint64) ([]*LayerRecord, error)
	GetStateTransitions(ctx context.Context, layerID uint64) ([]*StateTransitionRecord, error)
}

// LayerRecord — bản ghi layer từ DB
type LayerRecord struct {
	ID           uint64
	Name         string
	FamilyID     *uint64
	ParentID     *uint64
	ReplacedByID *uint64
	Version      int
	State        string
	CreatedAt    time.Time
	CreatedBy    string
}

// StateTransitionRecord — bản ghi chuyển trạng thái từ DB
type StateTransitionRecord struct {
	ID        uint64
	LayerID   uint64
	FromState string
	ToState   string
	Actor     string
	Reason    string
	LegalDoc  string
	CreatedAt time.Time
}

// LifecycleEngine — Engine quản lý vòng đời và version chain của layer
type LifecycleEngine struct {
	config    *types.LifecycleConfig
	layerRepo LayerRepository
}

// NewLifecycleEngine tạo LifecycleEngine
func NewLifecycleEngine(config *types.LifecycleConfig, layerRepo LayerRepository) *LifecycleEngine {
	if config == nil {
		config = &types.LifecycleConfig{
			AllowedTransitions:   types.DefaultAllowedTransitions(),
			AutoIncrementVersion: true,
		}
	}
	if config.AllowedTransitions == nil {
		config.AllowedTransitions = types.DefaultAllowedTransitions()
	}
	return &LifecycleEngine{config: config, layerRepo: layerRepo}
}

// Name trả về tên engine
func (e *LifecycleEngine) Name() string { return "LifecycleEngine" }

// Validate kiểm tra config hợp lệ
func (e *LifecycleEngine) Validate() error {
	if len(e.config.AllowedTransitions) == 0 {
		return fmt.Errorf("LifecycleEngine: allowed_transitions is empty")
	}
	return nil
}

// GetLineage truy vết toàn bộ cây phả hệ của 1 layer
func (e *LifecycleEngine) GetLineage(ctx context.Context, layerID uint64) (*types.LayerLineage, error) {
	layer, err := e.layerRepo.GetByID(ctx, layerID)
	if err != nil {
		return nil, fmt.Errorf("get layer %d: %w", layerID, err)
	}

	familyID := uint64(0)
	if layer.FamilyID != nil {
		familyID = *layer.FamilyID
	}

	// Lấy tất cả version trong cùng family
	records, err := e.layerRepo.GetByFamilyID(ctx, familyID)
	if err != nil {
		return nil, fmt.Errorf("get family %d: %w", familyID, err)
	}

	versions := make([]*types.LayerVersion, 0, len(records))
	for _, r := range records {
		transitions, _ := e.GetLifecycle(ctx, r.ID)

		v := &types.LayerVersion{
			Version:     r.Version,
			LayerID:     r.ID,
			ParentID:    r.ParentID,
			FamilyID:    familyID,
			State:       types.LifecycleState(r.State),
			Transitions: transitions,
			CreatedAt:   r.CreatedAt,
			CreatedBy:   r.CreatedBy,
		}
		versions = append(versions, v)
	}

	// Sắp xếp theo version
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].Version < versions[j].Version
	})

	// Đánh dấu latest
	if len(versions) > 0 {
		for i := range versions {
			versions[i].IsLatest = false
		}
		versions[len(versions)-1].IsLatest = true
	}

	return &types.LayerLineage{
		RootLayerID:  versions[0].LayerID,
		FamilyID:     familyID,
		Versions:     versions,
		CurrentID:    layerID,
		ReplacedByID: layer.ReplacedByID,
	}, nil
}

// GetLifecycle lấy toàn bộ lịch sử chuyển trạng thái của 1 layer
func (e *LifecycleEngine) GetLifecycle(ctx context.Context, layerID uint64) ([]types.LifecycleTransition, error) {
	records, err := e.layerRepo.GetStateTransitions(ctx, layerID)
	if err != nil {
		return nil, err
	}

	transitions := make([]types.LifecycleTransition, 0, len(records))
	for _, r := range records {
		transitions = append(transitions, types.LifecycleTransition{
			From:      types.LifecycleState(r.FromState),
			To:        types.LifecycleState(r.ToState),
			Timestamp: r.CreatedAt,
			Actor:     r.Actor,
			Reason:    r.Reason,
			LegalDoc:  r.LegalDoc,
		})
	}
	return transitions, nil
}

// ValidateTransition kiểm tra tính hợp lệ của 1 bước chuyển trạng thái
func (e *LifecycleEngine) ValidateTransition(from, to types.LifecycleState) error {
	allowed, ok := e.config.AllowedTransitions[from]
	if !ok {
		return fmt.Errorf("unknown state: %s", from)
	}
	for _, s := range allowed {
		if s == to {
			return nil
		}
	}
	return fmt.Errorf("invalid transition: %s → %s", from, to)
}

// DeriveState suy ra LifecycleState từ dữ liệu layer hiện có
func (e *LifecycleEngine) DeriveState(effectiveDate, expiryDate *time.Time, legalStatus uint32, replacedByID *uint64) types.LifecycleState {
	now := time.Now()

	if replacedByID != nil && *replacedByID > 0 {
		return types.StateReplaced
	}
	if expiryDate != nil && expiryDate.Before(now) {
		return types.StateExpired
	}
	if effectiveDate != nil && !effectiveDate.After(now) {
		return types.StateEffective
	}
	if effectiveDate != nil && effectiveDate.After(now) {
		return types.StateApproved
	}
	return types.StateDraft
}
