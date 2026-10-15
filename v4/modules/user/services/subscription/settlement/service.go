package settlement

import "context"

type Service struct{ store Store }

func NewService(store Store) *Service { return &Service{store: store} }

func (s *Service) ApplySettlement(ctx context.Context, effect Effect) (Result, error) {
	if s == nil || s.store == nil {
		return Result{}, ErrCurrentState
	}
	fingerprint, err := effect.Fingerprint()
	if err != nil {
		return Result{}, err
	}
	result, err := s.store.ApplySettlement(ctx, Acceptance{Effect: effect, Fingerprint: fingerprint})
	if err != nil {
		return Result{}, err
	}
	switch result.Action {
	case ActionRejectedSameTier:
		return result, ErrSameTier
	case ActionRejectedDowngrade:
		return result, ErrDowngrade
	case ActionActivated, ActionUpgraded:
		return result, nil
	default:
		return Result{}, ErrCurrentState
	}
}
