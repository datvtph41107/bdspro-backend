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

migrateup:
	@$(call with_env,migrate -path $(MIGRATION_DIR) -database "$(MIGRATION_DATABASE)" -verbose up)

migrateup1:
	@$(call with_env,migrate -path $(MIGRATION_DIR) -database "$(MIGRATION_DATABASE)" -verbose up 1)

migratedown:
	@$(call with_env,migrate -path $(MIGRATION_DIR) -database "$(MIGRATION_DATABASE)" -verbose down)

migratedown1:
	@$(call with_env,migrate -path $(MIGRATION_DIR) -database "$(MIGRATION_DATABASE)" -verbose down 1)

migrate-version:
	@$(call with_env,migrate -path $(MIGRATION_DIR) -database "$(MIGRATION_DATABASE)" version)

new_migration:
	@test -n "$(name)" || { echo 'usage: make new_migration name=<schema_change>' >&2; exit 2; }
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $(name)

# Compatibility aliases used by older scripts and muscle memory.
migrate-up: migrateup
migrate-down: migratedown1
migrate-create: new_migration

.PHONY: migrateup migrateup1 migratedown migratedown1 migrate-version new_migration \
	migrate-up migrate-down migrate-create
