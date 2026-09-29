# Canonical golang-migrate mechanics for a database-owning service.
#
# Required by owner Makefile:
#   MIGRATION_DIR
#   DB_URL                         safe local fallback only
#
# Optional:
#   MIGRATION_URL_ENV              default MIGRATION_URL
#   MIGRATION_URL_FALLBACK_ENV     e.g. HUB_DATABASE_URL
#
# Runtime secrets/DSNs still come from environment. This file only avoids
# repeating the same migrate commands across every owner.
MIGRATION_DIR ?= migrate
MIGRATION_URL_ENV ?= MIGRATION_URL

ifeq ($(strip $(MIGRATION_URL_FALLBACK_ENV)),)
MIGRATION_DATABASE := $${$(MIGRATION_URL_ENV):-$(DB_URL)}
else
MIGRATION_DATABASE := $${$(MIGRATION_URL_ENV):-$${$(MIGRATION_URL_FALLBACK_ENV):-$(DB_URL)}}
endif

MIGRATION_NAME := $(strip $(if $(name),$(name),$(NAME)))

# Public database developer API.
migrate:
	@$(call with_env,migrate -path $(MIGRATION_DIR) -database "$(MIGRATION_DATABASE)" -verbose up)

rollback:
	@$(call with_env,migrate -path $(MIGRATION_DIR) -database "$(MIGRATION_DATABASE)" -verbose down 1)

migration-version:
	@$(call with_env,migrate -path $(MIGRATION_DIR) -database "$(MIGRATION_DATABASE)" version)

migration:
	@test -n "$(MIGRATION_NAME)" || { echo 'usage: make migration name=<schema_change>' >&2; exit 2; }
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(MIGRATION_NAME)

# Compatibility-only targets. They are intentionally not advertised by service
# help; keep them until external callers have migrated to the public API.
migrateup: migrate

migrateup1:
	@$(call with_env,migrate -path $(MIGRATION_DIR) -database "$(MIGRATION_DATABASE)" -verbose up 1)

migratedown:
	@$(call with_env,migrate -path $(MIGRATION_DIR) -database "$(MIGRATION_DATABASE)" -verbose down)

migratedown1: rollback

migrate-version: migration-version
new_migration: migration
migrate-up: migrate
migrate-down: rollback
migrate-create: migration

.PHONY: migrate rollback migration-version migration \
	migrateup migrateup1 migratedown migratedown1 migrate-version new_migration \
	migrate-up migrate-down migrate-create
