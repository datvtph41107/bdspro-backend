package usecase

import (
	"context"
	"sync"

	"user/internal/domain/access"
	"user/internal/interface/repo"
)

// RoleGroupRegistry keeps the role-group lookup used by role resolution.
type RoleGroupRegistry struct {
	moduleMap     map[string]*access.RoleGroup // map theo ID
	mu            sync.RWMutex
	roleGroupRepo repo.RoleGroupRepository // inject repo để load modules
}

func NewRoleGroupRegistry(roleGroupRepo repo.RoleGroupRepository) *RoleGroupRegistry {
	return &RoleGroupRegistry{
		moduleMap:     make(map[string]*access.RoleGroup),
		mu:            sync.RWMutex{},
		roleGroupRepo: roleGroupRepo,
	}
}

// Load refreshes the role-group lookup from its durable repository.
func (s *RoleGroupRegistry) Load(ctx context.Context) error {
	modules, err := s.roleGroupRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Clear existing maps
	s.moduleMap = make(map[string]*access.RoleGroup)

	// Load modules vào maps
	for _, module := range modules {
		s.moduleMap[module.Code] = &module
	}

	return nil
}

func (s *RoleGroupRegistry) GetByCode(code string) (*access.RoleGroup, error) {
	module, ok := s.moduleMap[code]
	if !ok {
		return nil, access.ErrRoleGroupNotFound
	}
	return module, nil
}
