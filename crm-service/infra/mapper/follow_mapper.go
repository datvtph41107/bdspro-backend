package mapper

import (
	"crm/internal/domain"
)

func ListFollowersToPb(followers []domain.Profile) []uint64 {
	followersPb := make([]uint64, len(followers))
	for i, follower := range followers {
		followersPb[i] = follower.ProfileId
	}
	return followersPb
}