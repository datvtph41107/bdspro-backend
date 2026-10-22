package services

import "mdn4/modules/bdspro/enums"

const DealSettingEnumKey = "deal_setting"

// DealSettingRoleItems returns deal member roles shown in deal settings (RoleKeyMap deal rows 420–450).
func DealSettingRoleItems() []map[string]any {
	keys := []enums.RoleKey{
		enums.RoleKeyDealAdmin,
		enums.RoleKeyDealMember,
		enums.RoleKeyDealCustomer,
		enums.RoleKeyDealPartner,
	}
	items := make([]map[string]any, 0, len(keys))
	for _, k := range keys {
		items = append(items, map[string]any{
			"name":  enums.RoleKeyMap[k],
			"value": uint32(k),
		})
	}
	return items
}

// DealSettingEnumPayload is the `data` object for key deal_setting.
func DealSettingEnumPayload() map[string]any {
	return map[string]any{
		"roles":       DealSettingRoleItems(),
		"roleDefault": uint32(enums.RoleKeyDealMember),
	}
}

// DealSettingEnumListEntry is one element of GET /v2/bdspro/public/enums `data` array.
func DealSettingEnumListEntry() map[string]any {
	return map[string]any{
		"key":  DealSettingEnumKey,
		"data": DealSettingEnumPayload(),
	}
}
