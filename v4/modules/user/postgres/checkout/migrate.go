package checkoutpostgres

// AutoMigrateModels công bố persistence model private cho đúng một schema
// owner ở user/db. Handler và process root không được gọi GORM trực tiếp.
func AutoMigrateModels() []any {
	return []any{&checkoutSnapshotRow{}}
}
