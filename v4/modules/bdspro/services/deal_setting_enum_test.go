package services

import (
	"testing"

	"mdn4/modules/bdspro/enums"
)

func TestDealSettingEnumPayload(t *testing.T) {
	payload := DealSettingEnumPayload()
	roles, ok := payload["roles"].([]map[string]any)
	if !ok || len(roles) != 4 {
		t.Fatalf("roles: got %#v", payload["roles"])
	}
	if payload["roleDefault"] != uint32(enums.RoleKeyDealMember) {
		t.Fatalf("roleDefault: got %v", payload["roleDefault"])
	}
	want := []struct {
		value uint32
		name  string
	}{
		{uint32(enums.RoleKeyDealAdmin), "Quản lý"},
		{uint32(enums.RoleKeyDealMember), "Thành viên"},
		{uint32(enums.RoleKeyDealCustomer), "Đối tác"},
		{uint32(enums.RoleKeyDealPartner), "Khách hàng"},
	}
	for i, w := range want {
		if roles[i]["value"] != w.value || roles[i]["name"] != w.name {
			t.Fatalf("roles[%d]: got %+v want %+v", i, roles[i], w)
		}
	}
}
