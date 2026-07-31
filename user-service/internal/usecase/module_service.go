package usecase

import (
	_errors "common/errors"
	"context"
	"sync"

	"user/internal/domain"
	"user/internal/interface/repo"
)

// ModuleService quản lý cache modules trong memory
type ModuleService struct {
	moduleMap     map[string]*domain.RoleGroup // map theo ID
	mu            sync.RWMutex
	roleGroupRepo repo.RoleGroupRepository // inject repo để load modules
}

func NewModuleService(roleGroupRepo repo.RoleGroupRepository) *ModuleService {
	return &ModuleService{
		moduleMap:     make(map[string]*domain.RoleGroup),
		mu:            sync.RWMutex{},
		roleGroupRepo: roleGroupRepo,
	}
}

// LoadModules load tất cả modules từ database và cache vào memory
func (s *ModuleService) LoadModules(ctx context.Context) error {
	modules, err := s.roleGroupRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Clear existing maps
	s.moduleMap = make(map[string]*domain.RoleGroup)

	// Load modules vào maps
	for _, module := range modules {
		s.moduleMap[module.Code] = &module
	}

	return nil
}

func (s *ModuleService) GetModuleByCode(code string) (*domain.RoleGroup, error) {
	module, ok := s.moduleMap[code]
	if !ok {
		return nil, _errors.ReturnError(404, "Module not found")
	}
	return module, nil
}
