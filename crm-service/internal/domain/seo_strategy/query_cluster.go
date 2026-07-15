package seo_strategy

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type QueryClusterStatus string

const (
	QueryClusterStatusDraft    QueryClusterStatus = "draft"
	QueryClusterStatusActive   QueryClusterStatus = "active"
	QueryClusterStatusPaused   QueryClusterStatus = "paused"
	QueryClusterStatusArchived QueryClusterStatus = "archived"
)

type QueryClusterMembershipType string

const (
	QueryClusterMembershipPrimary  QueryClusterMembershipType = "primary"
	QueryClusterMembershipIncluded QueryClusterMembershipType = "included"
	QueryClusterMembershipExcluded QueryClusterMembershipType = "excluded"
)

type QueryClusterMemberStatus string

const (
	QueryClusterMemberStatusActive  QueryClusterMemberStatus = "active"
	QueryClusterMemberStatusRemoved QueryClusterMemberStatus = "removed"
)

type SeoQueryCluster struct {
	ID uint64

	ObjectiveID        uint64
	UserNeedID         uint64
	IntentHypothesisID uint64

	Code string
	Name string

	SearchPromise string
	InclusionRule string
	ExclusionRule string

	Status     QueryClusterStatus
	Confidence ObjectiveConfidence

	OwnerCode   string
	ReviewDueAt *time.Time
	LockVersion uint32

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type SeoQueryClusterMember struct {
	ObjectiveID uint64
	UserNeedID  uint64

	ClusterID       uint64
	QueryEvidenceID uint64

	MembershipType QueryClusterMembershipType
	Rationale      string
	Status         QueryClusterMemberStatus

	CreatedAt time.Time
	UpdatedAt time.Time
}

type LinkedQueryClusterMember struct {
	Member   SeoQueryClusterMember
	Evidence SeoQueryEvidence
}

type SeoQueryClusterAggregate struct {
	Cluster   *SeoQueryCluster
	Intent    *SeoIntentHypothesis
	UserNeed  *SeoUserNeed
	Members   []LinkedQueryClusterMember
	Readiness QueryClusterReadiness
}

// CreateSeoQueryClusterDraft là domain factory runtime, không phải Wire provider.
func CreateSeoQueryClusterDraft(
	objectiveID uint64,
	userNeedID uint64,
	intentHypothesisID uint64,
	code string,
	name string,
	now time.Time,
) (*SeoQueryCluster, error) {
	cluster := &SeoQueryCluster{
		ObjectiveID:        objectiveID,
		UserNeedID:         userNeedID,
		IntentHypothesisID: intentHypothesisID,
		Code:               NormalizeStrategyCode(code),
		Name:               strings.TrimSpace(name),
		Status:             QueryClusterStatusDraft,
		Confidence:         ObjectiveConfidenceExploratory,
		LockVersion:        1,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
	if err := ValidateQueryClusterIdentity(cluster); err != nil {
		return nil, err
	}
	return cluster, nil
}

func ValidateQueryClusterIdentity(cluster *SeoQueryCluster) error {
	if cluster == nil {
		return errors.New("query cluster không tồn tại")
	}
	if cluster.ObjectiveID == 0 {
		return errors.New("objectiveId không hợp lệ")
	}
	if cluster.UserNeedID == 0 {
		return errors.New("userNeedId không hợp lệ")
	}
	if cluster.IntentHypothesisID == 0 {
		return errors.New("intentHypothesisId không hợp lệ")
	}
	cluster.Code = NormalizeStrategyCode(cluster.Code)
	if err := ValidateStrategyCode(cluster.Code); err != nil {
		return err
	}
	cluster.Name = strings.TrimSpace(cluster.Name)
	if cluster.Name == "" {
		return errors.New("query cluster name không được để trống")
	}
	if len([]rune(cluster.Name)) > 500 {
		return errors.New("query cluster name không được dài quá 500 ký tự")
	}
	return nil
}

func ValidateQueryClusterDefinition(cluster *SeoQueryCluster) error {
	if err := ValidateQueryClusterIdentity(cluster); err != nil {
		return err
	}

	cluster.SearchPromise = strings.TrimSpace(cluster.SearchPromise)
	cluster.InclusionRule = strings.TrimSpace(cluster.InclusionRule)
	cluster.ExclusionRule = strings.TrimSpace(cluster.ExclusionRule)
	cluster.OwnerCode = strings.TrimSpace(cluster.OwnerCode)

	if len([]rune(cluster.SearchPromise)) > 6000 {
		return errors.New("search promise không được dài quá 6000 ký tự")
	}
	if len([]rune(cluster.InclusionRule)) > 6000 {
		return errors.New("inclusion rule không được dài quá 6000 ký tự")
	}
	if len([]rune(cluster.ExclusionRule)) > 6000 {
		return errors.New("exclusion rule không được dài quá 6000 ký tự")
	}
	if len([]rune(cluster.OwnerCode)) > 120 {
		return errors.New("owner code không được dài quá 120 ký tự")
	}
	if !IsValidQueryClusterStatus(cluster.Status) {
		return fmt.Errorf("query cluster status không hợp lệ: %s", cluster.Status)
	}
	if !IsValidObjectiveConfidence(cluster.Confidence) {
		return fmt.Errorf("query cluster confidence không hợp lệ: %s", cluster.Confidence)
	}
	if cluster.LockVersion == 0 {
		return errors.New("lock version phải lớn hơn 0")
	}
	return nil
}

func ValidateQueryClusterMember(member *SeoQueryClusterMember) error {
	if member == nil {
		return errors.New("query cluster member không tồn tại")
	}
	if member.ObjectiveID == 0 || member.UserNeedID == 0 ||
		member.ClusterID == 0 || member.QueryEvidenceID == 0 {
		return errors.New("query cluster member identity không hợp lệ")
	}
	if !IsValidQueryClusterMembershipType(member.MembershipType) {
		return fmt.Errorf("query cluster membership type không hợp lệ: %s", member.MembershipType)
	}
	if !IsValidQueryClusterMemberStatus(member.Status) {
		return fmt.Errorf("query cluster member status không hợp lệ: %s", member.Status)
	}
	member.Rationale = strings.TrimSpace(member.Rationale)
	if len([]rune(member.Rationale)) > 4000 {
		return errors.New("membership rationale không được dài quá 4000 ký tự")
	}
	if member.MembershipType != QueryClusterMembershipPrimary && member.Rationale == "" {
		return errors.New("included hoặc excluded query phải có membership rationale")
	}
	return nil
}

func IsValidQueryClusterStatus(value QueryClusterStatus) bool {
	switch value {
	case QueryClusterStatusDraft, QueryClusterStatusActive,
		QueryClusterStatusPaused, QueryClusterStatusArchived:
		return true
	default:
		return false
	}
}

func IsValidQueryClusterMembershipType(value QueryClusterMembershipType) bool {
	switch value {
	case QueryClusterMembershipPrimary,
		QueryClusterMembershipIncluded,
		QueryClusterMembershipExcluded:
		return true
	default:
		return false
	}
}

func IsValidQueryClusterMemberStatus(value QueryClusterMemberStatus) bool {
	switch value {
	case QueryClusterMemberStatusActive, QueryClusterMemberStatusRemoved:
		return true
	default:
		return false
	}
}

func ParseQueryClusterStatus(value string) (QueryClusterStatus, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return QueryClusterStatusDraft, nil
	}
	result := QueryClusterStatus(value)
	if !IsValidQueryClusterStatus(result) {
		return "", fmt.Errorf("query cluster status không hợp lệ: %s", value)
	}
	return result, nil
}

func ParseQueryClusterMembershipType(value string) (QueryClusterMembershipType, error) {
	result := QueryClusterMembershipType(strings.TrimSpace(value))
	if !IsValidQueryClusterMembershipType(result) {
		return "", fmt.Errorf("query cluster membership type không hợp lệ: %s", value)
	}
	return result, nil
}
