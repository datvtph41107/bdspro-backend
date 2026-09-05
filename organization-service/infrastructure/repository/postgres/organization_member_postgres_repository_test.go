package postgres

import (
	"testing"

	"organization/internal/domain/entity"

	"github.com/stretchr/testify/assert"
)

func TestOrganizationMemberRoleCompatibility(t *testing.T) {
	t.Parallel()

	legacy := OrganizationMemberModelToEntity(&OrganizationMemberModel{Role: 9})
	assert.Equal(t, uint64(9), legacy.RoleId, "dữ liệu production legacy phải đọc role khi role_id chưa có")

	model := OrganizationMemberEntityToModel(&entity.OrganizationMember{RoleId: 11})
	assert.Equal(t, uint32(11), model.Role)
	assert.Equal(t, uint64(11), model.RoleId, "write mới phải đồng bộ hai cột trong giai đoạn compatibility")
}
