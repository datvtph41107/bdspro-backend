package read

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"user/internal/models"
	"user/internal/modules/subscription/domain"
	"user/internal/modules/subscription/shadow"
)

const (
	ModeEnvironment          = "QHPRO_SUBSCRIPTION_READ_MODE"
	CanaryPercentEnvironment = "QHPRO_SUBSCRIPTION_CANARY_PERCENT"
	CanarySaltEnvironment    = "QHPRO_SUBSCRIPTION_CANARY_SALT"
	TimeoutEnvironment       = "QHPRO_SUBSCRIPTION_READ_TIMEOUT"
)

var ErrCatalogCutoverInvalid = errors.New("catalog subscription cutover is invalid")

type Mode string

const (
	ModeLegacy  Mode = "legacy"
	ModeCanary  Mode = "canary"
	ModeCatalog Mode = "catalog"
)

/** Logger is the minimal rollout log dependency. */
type Logger interface {
	Printf(format string, values ...any)
}

/** Config controls deterministic and reversible read routing. */
type Config struct {
	Mode          Mode
	CanaryPercent int
	CanarySalt    string
	Timeout       time.Duration
	Now           func() time.Time
	Logger        Logger
}

/** Router preserves the legacy API while selecting a verified source. */
// @bind: user/internal/modules/subscription/shadow.Observer
type Router struct {
	repository    shadow.Repository
	shadow        *shadow.Reader
	mode          Mode
	canaryPercent int
	canarySalt    string
	timeout       time.Duration
	now           func() time.Time
	logger        Logger
}

/** NewRouter creates an environment-controlled read router. */
func NewRouter(repository shadow.Repository) *Router {
	router := BuildRouterWithConfig(repository, Config{
		Mode:          parseMode(os.Getenv(ModeEnvironment)),
		CanaryPercent: parsePercent(os.Getenv(CanaryPercentEnvironment)),
		CanarySalt:    strings.TrimSpace(os.Getenv(CanarySaltEnvironment)),
		Timeout:       parseTimeout(os.Getenv(TimeoutEnvironment)),
		Now:           time.Now,
		Logger:        log.Default(),
	})
	// Task 4.4 observation remains independently default-off.
	router.shadow = shadow.NewReader(repository)
	return router
}

/** BuildRouterWithConfig creates a deterministic read router. */
func BuildRouterWithConfig(
	repository shadow.Repository,
	config Config,
) *Router {
	if config.Mode != ModeCanary && config.Mode != ModeCatalog {
		config.Mode = ModeLegacy
	}
	if config.CanaryPercent < 0 || config.CanaryPercent > 100 {
		config.CanaryPercent = 0
	}
	if config.Timeout <= 0 {
		config.Timeout = 200 * time.Millisecond
	}
	if config.Now == nil {
		config.Now = time.Now
	}
	if config.Logger == nil {
		config.Logger = log.Default()
	}
	return &Router{
		repository: repository,
		shadow: shadow.BuildReaderWithConfig(
			repository,
			shadow.Config{
				Mode:    shadow.ModeOff,
				Timeout: config.Timeout,
				Now:     config.Now,
				Logger:  config.Logger,
			},
		),
		mode:          config.Mode,
		canaryPercent: config.CanaryPercent,
		canarySalt:    strings.TrimSpace(config.CanarySalt),
		timeout:       config.Timeout,
		now:           config.Now,
		logger:        config.Logger,
	}
}

/** ObserveMembership keeps Task 4.4 observation in legacy mode. */
func (r *Router) ObserveMembership(ctx context.Context, membership *models.Membership) {
	if r == nil || r.mode != ModeLegacy || r.shadow == nil {
		return
	}
	r.shadow.ObserveMembership(ctx, membership)
}

/** SelectMembership returns legacy data or a verified compatible projection. */
func (r *Router) SelectMembership(
	ctx context.Context,
	legacy *models.Membership,
) (*models.Membership, error) {
	if legacy == nil {
		return nil, fmt.Errorf("legacy membership is required")
	}
	if r == nil || r.mode == ModeLegacy {
		return copyMembership(legacy), nil
	}
	if r.repository == nil || r.shadow == nil {
		if r.mode == ModeCanary {
			return copyMembership(legacy), nil
		}
		return nil, fmt.Errorf("%w: repository is not configured", ErrCatalogCutoverInvalid)
	}
	if r.mode == ModeCanary && !r.inCanary(legacy.ProfileID) {
		return copyMembership(legacy), nil
	}

	readCtx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	if r.mode == ModeCatalog {
		mapping, err := r.repository.GetLegacyMembershipMapping(readCtx, legacy.ID)
		if errors.Is(err, shadow.ErrNotFound) {
			return copyMembership(legacy), nil
		}
		if err != nil {
			return nil, fmt.Errorf("%w: load cutover mapping: %v", ErrCatalogCutoverInvalid, err)
		}
		if mapping.Status != models.CatalogMigrationCutover {
			return copyMembership(legacy), nil
		}
	}

	entry := r.shadow.CompareMembership(readCtx, *legacy, r.now().UTC())
	if entry.Status != shadow.StatusMatched {
		if r.mode == ModeCanary {
			r.logger.Printf(
				"subscription-read mode=canary membership_id=%d status=fallback reason=%s detail=%q",
				legacy.ID,
				entry.Status,
				entry.Detail,
			)
			return copyMembership(legacy), nil
		}
		return nil, fmt.Errorf(
			"%w: membership_id=%d status=%s detail=%s",
			ErrCatalogCutoverInvalid,
			legacy.ID,
			entry.Status,
			entry.Detail,
		)
	}

	aggregate, err := r.repository.GetSubscription(readCtx, entry.SubscriptionID)
	if err != nil {
		if r.mode == ModeCanary {
			r.logger.Printf(
				"subscription-read mode=canary membership_id=%d status=fallback reason=subscription_reload detail=%q",
				legacy.ID,
				err.Error(),
			)
			return copyMembership(legacy), nil
		}
		return nil, fmt.Errorf("%w: reload subscription: %v", ErrCatalogCutoverInvalid, err)
	}
	if err := aggregate.Validate(); err != nil {
		if r.mode == ModeCanary {
			return copyMembership(legacy), nil
		}
		return nil, fmt.Errorf("%w: %v", ErrCatalogCutoverInvalid, err)
	}
	projected, err := projectLegacyCompatible(*legacy, aggregate)
	if err != nil {
		if r.mode == ModeCanary {
			r.logger.Printf(
				"subscription-read mode=canary membership_id=%d status=fallback reason=projection detail=%q",
				legacy.ID,
				err.Error(),
			)
			return copyMembership(legacy), nil
		}
		return nil, fmt.Errorf("%w: %v", ErrCatalogCutoverInvalid, err)
	}
	r.logger.Printf(
		"subscription-read mode=%s membership_id=%d subscription_id=%d status=catalog",
		r.mode,
		legacy.ID,
		aggregate.ID,
	)
	return projected, nil
}

func (r *Router) inCanary(profileID uint64) bool {
	if r.canaryPercent <= 0 || profileID == 0 || r.canarySalt == "" {
		return false
	}
	if r.canaryPercent >= 100 {
		return true
	}
	digest := sha256.Sum256([]byte(r.canarySalt + ":" + strconv.FormatUint(profileID, 10)))
	bucket := binary.BigEndian.Uint64(digest[:8]) % 100
	return int(bucket) < r.canaryPercent
}

/** ValidateConfig rejects accidental or unstable rollout settings. */
func ValidateConfig(config Config) error {
	if config.Mode != ModeLegacy && config.Mode != ModeCanary && config.Mode != ModeCatalog {
		return fmt.Errorf("unsupported subscription read mode %q", config.Mode)
	}
	if config.CanaryPercent < 0 || config.CanaryPercent > 100 {
		return fmt.Errorf("canary percent must be between 0 and 100")
	}
	if config.Mode == ModeCanary && config.CanaryPercent > 0 && strings.TrimSpace(config.CanarySalt) == "" {
		return fmt.Errorf("canary salt is required")
	}
	return nil
}

func projectLegacyCompatible(
	legacy models.Membership,
	aggregate domain.Aggregate,
) (*models.Membership, error) {
	return projectAggregate(legacy, aggregate)
}

func parseMode(value string) Mode {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case string(ModeCanary):
		return ModeCanary
	case string(ModeCatalog):
		return ModeCatalog
	default:
		return ModeLegacy
	}
}

func parsePercent(value string) int {
	if strings.TrimSpace(value) == "" {
		return 0
	}
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || parsed < 0 || parsed > 100 {
		return 0
	}
	return parsed
}

func parseTimeout(value string) time.Duration {
	if strings.TrimSpace(value) == "" {
		return 200 * time.Millisecond
	}
	parsed, err := time.ParseDuration(value)
	if err != nil || parsed <= 0 {
		return 200 * time.Millisecond
	}
	return parsed
}

func copyMembership(value *models.Membership) *models.Membership {
	if value == nil {
		return nil
	}
	copied := *value
	copied.CreatedAt = copyTime(value.CreatedAt)
	copied.UpdatedAt = copyTime(value.UpdatedAt)
	copied.FinishedDate = copyTime(value.FinishedDate)
	copied.PaidAt = copyTime(value.PaidAt)
	return &copied
}

func copyTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	copied := *value
	return &copied
}
