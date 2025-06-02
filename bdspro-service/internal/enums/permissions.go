package enums

type EPermission string

const (
	PermissionProductCreate           EPermission = "product_create"
	PermissionProductRead             EPermission = "product_read"
	PermissionProductUpdate           EPermission = "product_update"
	PermissionProductDelete           EPermission = "product_delete"
	PermissionProductArchive          EPermission = "product_archive"
	PermissionProductUpdateStatus     EPermission = "product_update_status"
	PermissionProductUpdateVisibility EPermission = "product_update_visibility"
	PermissionProductDeposite         EPermission = "product_deposite"
)
