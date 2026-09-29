package access

import (
	"errors"
	"strings"
)

const MaxSubjectIDLength = 128

/**
 * SubjectType cho biết quota/access đang thuộc về profile hay organization.
 */
type SubjectType string

const (
	SubjectProfile      SubjectType = "profile"
	SubjectOrganization SubjectType = "organization"
)

/**
 * Subject xác định ai đang sở hữu quyền sử dụng hiện tại.
 */
type Subject struct {
	Type SubjectType
	ID   string
}

var ErrInvalidSubject = errors.New("access subject is invalid")

func (s Subject) IsValid() bool {
	if s.Type != SubjectProfile && s.Type != SubjectOrganization {
		return false
	}
	if s.ID == "" || len(s.ID) > MaxSubjectIDLength || strings.TrimSpace(s.ID) != s.ID {
		return false
	}
	for _, r := range s.ID {
		if (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '-' || r == '_' || r == '.' || r == ':' {
			continue
		}
		return false
	}
	return true
}
