package enums

type GroupRoleKeyEnum uint32

const (
	GroupRoleKeyEnum_ADMIN  GroupRoleKeyEnum = 10
	GroupRoleKeyEnum_USER   GroupRoleKeyEnum = 20
	GroupRoleKeyEnum_ORG    GroupRoleKeyEnum = 30
	GroupRoleKeyEnum_CRM    GroupRoleKeyEnum = 40
	GroupRoleKeyEnum_GROUP  GroupRoleKeyEnum = 50
	GroupRoleKeyEnum_BRANCH GroupRoleKeyEnum = 60
	GroupRoleKeyEnum_DEAL   GroupRoleKeyEnum = 70
)

var GroupRoleKeyEnumMap = map[GroupRoleKeyEnum]string{
	GroupRoleKeyEnum_ADMIN:  "ADMIN",
	GroupRoleKeyEnum_USER:   "USER",
	GroupRoleKeyEnum_ORG:    "ORG",
	GroupRoleKeyEnum_CRM:    "CRM",
	GroupRoleKeyEnum_GROUP:  "GROUP",
	GroupRoleKeyEnum_BRANCH: "BRANCH",
	GroupRoleKeyEnum_DEAL:   "DEAL",
}

var GroupRoleKeyEnumNames = map[GroupRoleKeyEnum]string{
	GroupRoleKeyEnum_ADMIN:  "Quản trị hệ thống",
	GroupRoleKeyEnum_USER:   "Người dùng",
	GroupRoleKeyEnum_ORG:    "Tổ chức",
	GroupRoleKeyEnum_CRM:    "Quản lý CRM",
	GroupRoleKeyEnum_GROUP:  "Nhóm",
	GroupRoleKeyEnum_BRANCH: "Chi nhánh",
	GroupRoleKeyEnum_DEAL:   "Thương vụ",
}
