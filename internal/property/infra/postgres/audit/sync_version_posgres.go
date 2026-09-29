package postgres_audit

import (
	"context"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type VersionSyncPosgres struct {
	db *gorm.DB
}

type DataVersion struct {
	ID          string
	SyncVersion uint64
	Payload     []byte
}

type ClientChange struct {
	ID              string
	Data            []byte // JSON bytes
	Operation       string // CREATE, UPDATE, DELETE
	BaseSyncVersion uint64
	ClientChangeID  string
}

type PushRepo struct {
	db *gorm.DB
}

func NewPushRepo(db *gorm.DB) *PushRepo {
	return &PushRepo{db: db}
}

func NewVersionRepo(db *gorm.DB) *VersionSyncPosgres {
	return &VersionSyncPosgres{db: db}
}

// Increment atomically increments the max_version for a resource and owner.
func (r *VersionSyncPosgres) Increment(
	ctx context.Context,
	resource string,
	ownerID uint64,
) (uint64, error) {
	query := `
		INSERT INTO sync_versions (resource, owner_id, max_version)
		VALUES (?, ?, 1)
		ON CONFLICT (resource, owner_id) DO UPDATE
		SET max_version = sync_versions.max_version + 1
		RETURNING max_version;
	`

	var version uint64
	err := r.db.WithContext(ctx).
		Raw(query, resource, ownerID).
		Scan(&version).Error

	if err != nil {
		return 0, fmt.Errorf("failed to increment version: %w", err)
	}

	return version, nil
}

// Get returns the current max_version for a resource and owner.
func (r *VersionSyncPosgres) Get(
	ctx context.Context,
	resource string,
	ownerID uint64,
) (uint64, error) {
	query := `
		SELECT COALESCE(max_version, 0)
		FROM sync_versions
		WHERE resource = ? AND owner_id = ?;
	`

	var version uint64
	err := r.db.WithContext(ctx).
		Raw(query, resource, ownerID).
		Scan(&version).Error

	if err != nil {
		return 0, err
	}

	return version, nil
}

// GetDelta returns rows with sync_version > sinceVersion, ordered by sync_version ASC.
// If cursorVersion > 0, it is used as an additional filter (for pagination).
func (r *VersionSyncPosgres) GetData(
	ctx context.Context,
	resource string,
	ownerID uint64,
	sinceVersion uint64,
	cursorVersion uint64,
	limit int,
) ([]DataVersion, error) {

	// whitelist resource → table
	var tableName string
	switch resource {
	case "products":
		tableName = "products"
	case "product_user":
		tableName = "product_user"
	default:
		return nil, fmt.Errorf("unsupported resource: %s", resource)
	}

	query := fmt.Sprintf(`
		SELECT id, sync_version, row_to_json(%s)::text AS payload
		FROM %s
		WHERE owner_id = ?
		  AND sync_version > ?
		  AND sync_version > ?
		ORDER BY sync_version ASC
		LIMIT ?;
	`, tableName, tableName)

	var rows []struct {
		ID          string
		SyncVersion uint64
		Payload     string
	}

	err := r.db.WithContext(ctx).
		Raw(query, ownerID, sinceVersion, cursorVersion, limit).
		Scan(&rows).Error

	if err != nil {
		return nil, fmt.Errorf("failed to query delta: %w", err)
	}

	results := make([]DataVersion, 0, len(rows))
	for _, r := range rows {
		results = append(results, DataVersion{
			ID:          r.ID,
			SyncVersion: r.SyncVersion,
			Payload:     []byte(r.Payload),
		})
	}

	return results, nil
}

func (r *PushRepo) ApplyChange(
	ctx context.Context,
	resource string,
	ownerID uint64,
	change ClientChange,
	newVersion uint64,
) error {
	// hiện tại support products
	if resource != "products" {
		return fmt.Errorf("unsupported resource: %s", resource)
	}

	// 1️⃣ lấy sync_version hiện tại
	var currentSyncVersion uint64
	err := r.db.WithContext(ctx).
		Raw(
			`SELECT sync_version FROM products WHERE id = ? AND owner_id = ?`,
			change.ID, ownerID,
		).
		Scan(&currentSyncVersion).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to get current sync version: %w", err)
	}

	// 2️⃣ detect conflict
	if change.Operation != "CREATE" && currentSyncVersion > change.BaseSyncVersion {
		return fmt.Errorf(
			"conflict: current version %d > base version %d",
			currentSyncVersion,
			change.BaseSyncVersion,
		)
	}

	switch change.Operation {
	case "CREATE":
		return r.applyCreate(ctx, ownerID, change, newVersion)
	case "UPDATE":
		return r.applyUpdate(ctx, ownerID, change, newVersion)
	case "DELETE":
		return r.applyDelete(ctx, ownerID, change, newVersion)
	default:
		return fmt.Errorf("unknown operation: %s", change.Operation)
	}
}

func (r *PushRepo) applyCreate(
	ctx context.Context,
	ownerID uint64,
	change ClientChange,
	newVersion uint64,
) error {

	var data map[string]interface{}
	if err := json.Unmarshal(change.Data, &data); err != nil {
		return fmt.Errorf("invalid data: %w", err)
	}

	query := `
		INSERT INTO products (
			id, owner_id, name, code, description, status, price, sync_version
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	return r.db.WithContext(ctx).Exec(
		query,
		change.ID,
		ownerID,
		data["name"],
		data["code"],
		data["description"],
		data["status"],
		data["price"],
		newVersion,
	).Error
}

func (r *PushRepo) applyUpdate(
	ctx context.Context,
	ownerID uint64,
	change ClientChange,
	newVersion uint64,
) error {

	var data map[string]interface{}
	if err := json.Unmarshal(change.Data, &data); err != nil {
		return err
	}

	query := `
		UPDATE products
		SET name = COALESCE(?, name),
		    code = COALESCE(?, code),
		    description = COALESCE(?, description),
		    status = COALESCE(?, status),
		    price = COALESCE(?, price),
		    sync_version = ?,
		    updated_at = now()
		WHERE id = ? AND owner_id = ?
	`

	return r.db.WithContext(ctx).Exec(
		query,
		data["name"],
		data["code"],
		data["description"],
		data["status"],
		data["price"],
		newVersion,
		change.ID,
		ownerID,
	).Error
}

func (r *PushRepo) applyDelete(
	ctx context.Context,
	ownerID uint64,
	change ClientChange,
	newVersion uint64,
) error {

	query := `
		UPDATE products
		SET deleted_at = now(),
		    sync_version = ?,
		    updated_at = now()
		WHERE id = ? AND owner_id = ?
	`

	return r.db.WithContext(ctx).
		Exec(query, newVersion, change.ID, ownerID).
		Error
}
