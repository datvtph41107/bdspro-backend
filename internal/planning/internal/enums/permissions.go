package enums

// Permission represents different types of permissions
type Permission uint32

const (
	// POI Permissions
	PermissionPOICreate Permission = iota
	PermissionPOIRead
	PermissionPOIUpdate
	PermissionPOIDelete
	PermissionPOIArchive
	PermissionPOIUpdateStatus
	PermissionPOIUpdateVisibility

	// POI Category Permissions
	PermissionPOICategoryCreate
	PermissionPOICategoryRead
	PermissionPOICategoryUpdate
	PermissionPOICategoryDelete

	// Amenity Permissions
	PermissionAmenityCreate
	PermissionAmenityRead
	PermissionAmenityUpdate
	PermissionAmenityDelete

	// Directory Category Permissions
	PermissionDirectoryCategoryCreate
	PermissionDirectoryCategoryRead
	PermissionDirectoryCategoryUpdate
	PermissionDirectoryCategoryDelete

	// POI Media Permissions
	PermissionPOIMediaCreate
	PermissionPOIMediaRead
	PermissionPOIMediaUpdate
	PermissionPOIMediaDelete

	// Open Hour Permissions
	PermissionOpenHourCreate
	PermissionOpenHourRead
	PermissionOpenHourUpdate
	PermissionOpenHourDelete

	// Contact Label Permissions
	PermissionContactLabelCreate
	PermissionContactLabelRead
	PermissionContactLabelUpdate
	PermissionContactLabelDelete

	// Directory Source Permissions
	PermissionDirectorySourceCreate
	PermissionDirectorySourceRead
	PermissionDirectorySourceUpdate
	PermissionDirectorySourceDelete

	// Directory Supplier Permissions
	PermissionDirectorySupplierCreate
	PermissionDirectorySupplierRead
	PermissionDirectorySupplierUpdate
	PermissionDirectorySupplierDelete
	PermissionDirectorySupplierList
	PermissionDirectorySupplierExport
	PermissionDirectorySupplierImport
)

// PermissionString returns the string representation of permission
func (p Permission) String() string {
	switch p {
	// POI Permissions
	case PermissionPOICreate:
		return "poi_create"
	case PermissionPOIRead:
		return "poi_read"
	case PermissionPOIUpdate:
		return "poi_update"
	case PermissionPOIDelete:
		return "poi_delete"
	case PermissionPOIArchive:
		return "poi_archive"
	case PermissionPOIUpdateStatus:
		return "poi_update_status"
	case PermissionPOIUpdateVisibility:
		return "poi_update_visibility"

	// POI Category Permissions
	case PermissionPOICategoryCreate:
		return "poi_category_create"
	case PermissionPOICategoryRead:
		return "poi_category_read"
	case PermissionPOICategoryUpdate:
		return "poi_category_update"
	case PermissionPOICategoryDelete:
		return "poi_category_delete"

	// Amenity Permissions
	case PermissionAmenityCreate:
		return "amenity_create"
	case PermissionAmenityRead:
		return "amenity_read"
	case PermissionAmenityUpdate:
		return "amenity_update"
	case PermissionAmenityDelete:
		return "amenity_delete"

	// POI Media Permissions
	case PermissionPOIMediaCreate:
		return "poi_media_create"
	case PermissionPOIMediaRead:
		return "poi_media_read"
	case PermissionPOIMediaUpdate:
		return "poi_media_update"
	case PermissionPOIMediaDelete:
		return "poi_media_delete"

	// Open Hour Permissions
	case PermissionOpenHourCreate:
		return "open_hour_create"
	case PermissionOpenHourRead:
		return "open_hour_read"
	case PermissionOpenHourUpdate:
		return "open_hour_update"
	case PermissionOpenHourDelete:
		return "open_hour_delete"

	// Contact Label Permissions
	case PermissionContactLabelCreate:
		return "contact_label_create"
	case PermissionContactLabelRead:
		return "contact_label_read"
	case PermissionContactLabelUpdate:
		return "contact_label_update"
	case PermissionContactLabelDelete:
		return "contact_label_delete"

	// Directory Category Permissions
	case PermissionDirectoryCategoryCreate:
		return "directory_category_create"
	case PermissionDirectoryCategoryRead:
		return "directory_category_read"
	case PermissionDirectoryCategoryUpdate:
		return "directory_category_update"
	case PermissionDirectoryCategoryDelete:
		return "directory_category_delete"

	// Directory Source Permissions
	case PermissionDirectorySourceCreate:
		return "directory_source_create"
	case PermissionDirectorySourceRead:
		return "directory_source_read"
	case PermissionDirectorySourceUpdate:
		return "directory_source_update"
	case PermissionDirectorySourceDelete:
		return "directory_source_delete"

	// Directory Supplier Permissions
	case PermissionDirectorySupplierCreate:
		return "directory_supplier_create"
	case PermissionDirectorySupplierRead:
		return "directory_supplier_read"
	case PermissionDirectorySupplierUpdate:
		return "directory_supplier_update"
	case PermissionDirectorySupplierDelete:
		return "directory_supplier_delete"
	case PermissionDirectorySupplierList:
		return "directory_supplier_list"
	case PermissionDirectorySupplierExport:
		return "directory_supplier_export"
	case PermissionDirectorySupplierImport:
		return "directory_supplier_import"

	default:
		return "unknown"
	}
}
