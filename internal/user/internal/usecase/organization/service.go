package organization

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strings"

	domain "user/internal/domain/organization"
)

const (
	PermissionView   = "ORGANIZATION_VIEW"
	PermissionManage = "ORGANIZATION_MANAGE"
)

var (
	ErrInvalidInput       = errors.New("invalid organization input")
	ErrNotFound           = errors.New("organization not found")
	ErrConflict           = errors.New("organization conflicts with existing data")
	ErrActiveSubscription = errors.New("organization has an active subscription")
)

var codePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]{1,63}$`)

type Repository interface {
	List(context.Context, domain.ListQuery) ([]domain.Organization, uint64, error)
	Get(context.Context, uint64) (domain.Organization, error)
	Create(context.Context, domain.Organization) (domain.Organization, error)
	Update(context.Context, domain.Update) (domain.Organization, error)
	Archive(context.Context, uint64, uint64) error
	CheckMembers(context.Context, uint64, []uint64) ([]domain.MemberCheck, error)
	ListForProfile(context.Context, uint64) ([]domain.Organization, error)
}

func (s *Service) CheckMembers(ctx context.Context, organizationID uint64, profileIDs []uint64) ([]domain.MemberCheck, error) {
	if s == nil || s.repository == nil {
		return nil, errors.New("organization repository is unavailable")
	}
	if organizationID == 0 {
		return nil, fmt.Errorf("%w: organization id is required", ErrInvalidInput)
	}
	if len(profileIDs) > 500 {
		return nil, fmt.Errorf("%w: at most 500 profiles may be checked", ErrInvalidInput)
	}
	seen := make(map[uint64]struct{}, len(profileIDs))
	unique := make([]uint64, 0, len(profileIDs))
	for _, profileID := range profileIDs {
		if profileID == 0 {
			continue
		}
		if _, exists := seen[profileID]; exists {
			continue
		}
		seen[profileID] = struct{}{}
		unique = append(unique, profileID)
	}
	return s.repository.CheckMembers(ctx, organizationID, unique)
}

func (s *Service) ListForProfile(ctx context.Context, profileID uint64) ([]domain.Organization, error) {
	if s == nil || s.repository == nil {
		return nil, errors.New("organization repository is unavailable")
	}
	if profileID == 0 {
		return nil, fmt.Errorf("%w: profile id is required", ErrInvalidInput)
	}
	return s.repository.ListForProfile(ctx, profileID)
}

type Service struct{ repository Repository }

func NewService(repository Repository) *Service { return &Service{repository: repository} }

func normalizeType(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "business", "doanh nghiệp", "doanh nghiep":
		return "business"
	case "team", "nhóm", "nhom":
		return "team"
	default:
		return strings.ToLower(strings.TrimSpace(value))
	}
}

func validateRequired(organization domain.Organization) error {
	if strings.TrimSpace(organization.Name) == "" {
		return fmt.Errorf("%w: name is required", ErrInvalidInput)
	}
	if !codePattern.MatchString(strings.TrimSpace(organization.Code)) {
		return fmt.Errorf("%w: code must be 2-64 letters, numbers, dot, dash or underscore", ErrInvalidInput)
	}
	if organization.Type != "business" && organization.Type != "team" && organization.Type != "other" {
		return fmt.Errorf("%w: unsupported organization type", ErrInvalidInput)
	}
	if value := strings.TrimSpace(organization.Email); value != "" {
		if _, err := mail.ParseAddress(value); err != nil {
			return fmt.Errorf("%w: invalid email", ErrInvalidInput)
		}
	}
	if organization.Status != "active" && organization.Status != "locked" {
		return fmt.Errorf("%w: unsupported organization status", ErrInvalidInput)
	}
	if organization.VerificationStatus != "unverified" && organization.VerificationStatus != "verifying" && organization.VerificationStatus != "verified" {
		return fmt.Errorf("%w: unsupported verification status", ErrInvalidInput)
	}
	switch organization.WarningLevel {
	case "none", "mild", "medium", "severe", "warned":
	default:
		return fmt.Errorf("%w: unsupported warning level", ErrInvalidInput)
	}
	return nil
}

func (s *Service) List(ctx context.Context, query domain.ListQuery) ([]domain.Organization, uint64, error) {
	if s == nil || s.repository == nil {
		return nil, 0, errors.New("organization repository is unavailable")
	}
	if query.Size == 0 {
		query.Size = 15
	}
	if query.Size > 100 {
		return nil, 0, fmt.Errorf("%w: size cannot exceed 100", ErrInvalidInput)
	}
	query.SearchText = strings.TrimSpace(query.SearchText)
	query.Type = normalizeType(query.Type)
	query.Status = strings.ToLower(strings.TrimSpace(query.Status))
	query.VerificationStatus = strings.ToLower(strings.TrimSpace(query.VerificationStatus))
	query.WarningLevel = strings.ToLower(strings.TrimSpace(query.WarningLevel))
	return s.repository.List(ctx, query)
}

func (s *Service) Get(ctx context.Context, id uint64) (domain.Organization, error) {
	if s == nil || s.repository == nil {
		return domain.Organization{}, errors.New("organization repository is unavailable")
	}
	if id == 0 {
		return domain.Organization{}, fmt.Errorf("%w: id is required", ErrInvalidInput)
	}
	return s.repository.Get(ctx, id)
}

func (s *Service) Create(ctx context.Context, organization domain.Organization) (domain.Organization, error) {
	if s == nil || s.repository == nil {
		return domain.Organization{}, errors.New("organization repository is unavailable")
	}
	organization.Code = strings.TrimSpace(organization.Code)
	organization.Name = strings.TrimSpace(organization.Name)
	organization.Type = normalizeType(organization.Type)
	organization.TaxCode = strings.TrimSpace(organization.TaxCode)
	organization.Email = strings.TrimSpace(organization.Email)
	organization.Status = "active"
	organization.VerificationStatus = "unverified"
	organization.WarningLevel = "none"
	if organization.OwnerProfileID == 0 || organization.CreatedBy == 0 {
		return domain.Organization{}, fmt.Errorf("%w: actor is required", ErrInvalidInput)
	}
	if err := validateRequired(organization); err != nil {
		return domain.Organization{}, err
	}
	return s.repository.Create(ctx, organization)
}

func (s *Service) Update(ctx context.Context, update domain.Update) (domain.Organization, error) {
	if s == nil || s.repository == nil {
		return domain.Organization{}, errors.New("organization repository is unavailable")
	}
	if update.ID == 0 || update.ActorID == 0 {
		return domain.Organization{}, fmt.Errorf("%w: id and actor are required", ErrInvalidInput)
	}
	current, err := s.repository.Get(ctx, update.ID)
	if err != nil {
		return domain.Organization{}, err
	}
	apply := func(target *string, value *string) {
		if value != nil {
			*target = strings.TrimSpace(*value)
		}
	}
	normalize := func(value **string) {
		if *value == nil {
			return
		}
		normalized := strings.TrimSpace(**value)
		*value = &normalized
	}
	for _, value := range []**string{
		&update.Code, &update.Name, &update.TaxCode, &update.Phone, &update.Email,
		&update.Address, &update.Website, &update.Description, &update.Status,
		&update.VerificationStatus, &update.WarningLevel,
	} {
		normalize(value)
	}
	if update.Type != nil {
		normalized := normalizeType(*update.Type)
		update.Type = &normalized
	}
	for _, value := range []**string{&update.Status, &update.VerificationStatus, &update.WarningLevel} {
		if *value != nil {
			normalized := strings.ToLower(**value)
			*value = &normalized
		}
	}
	apply(&current.Code, update.Code)
	apply(&current.Name, update.Name)
	if update.Type != nil {
		current.Type = normalizeType(*update.Type)
	}
	apply(&current.TaxCode, update.TaxCode)
	apply(&current.Phone, update.Phone)
	apply(&current.Email, update.Email)
	apply(&current.Address, update.Address)
	apply(&current.Website, update.Website)
	apply(&current.Description, update.Description)
	apply(&current.Status, update.Status)
	apply(&current.VerificationStatus, update.VerificationStatus)
	apply(&current.WarningLevel, update.WarningLevel)
	if err := validateRequired(current); err != nil {
		return domain.Organization{}, err
	}
	return s.repository.Update(ctx, update)
}

func (s *Service) Archive(ctx context.Context, id, actorID uint64) error {
	if s == nil || s.repository == nil {
		return errors.New("organization repository is unavailable")
	}
	if id == 0 || actorID == 0 {
		return fmt.Errorf("%w: id and actor are required", ErrInvalidInput)
	}
	return s.repository.Archive(ctx, id, actorID)
}
