# Canonical local backing-infrastructure mechanics for a service.
#
# The including service owns the dependency declaration:
#
#   DEV_DEPENDENCIES := postgres rabbitmq
#
# This shared file only owns how those long-running dependencies are started.
# Schema migration is intentionally not part of deps-up: migrate jobs are
# one-shot state transitions and cannot share Docker Compose --wait semantics
# with healthy long-running infrastructure.

DEV_DEPENDENCIES ?=
REPOSITORY_ROOT ?= ..
DEVELOPMENT_COMPOSE_FILE ?= $(REPOSITORY_ROOT)/compose.yaml
DEVELOPMENT_ENV_FILE ?= $(REPOSITORY_ROOT)/.env
QHPRO_DOCKER_CONFIG ?= $(abspath $(REPOSITORY_ROOT)/.tmp/docker-public)

deps-up:
	@test -n "$(strip $(DEV_DEPENDENCIES))" || { echo 'DEV_DEPENDENCIES is empty for this service' >&2; exit 2; }
	@test -f "$(DEVELOPMENT_COMPOSE_FILE)" || { echo 'missing $(DEVELOPMENT_COMPOSE_FILE)' >&2; exit 2; }
	@test -f "$(DEVELOPMENT_ENV_FILE)" || { echo 'missing $(DEVELOPMENT_ENV_FILE)' >&2; exit 2; }
	@mkdir -p "$(QHPRO_DOCKER_CONFIG)"
	@DOCKER_CONFIG="$(QHPRO_DOCKER_CONFIG)" docker compose \
		--env-file "$(DEVELOPMENT_ENV_FILE)" \
		-f "$(DEVELOPMENT_COMPOSE_FILE)" \
		up -d --wait $(DEV_DEPENDENCIES)

deps-down:
	@test -n "$(strip $(DEV_DEPENDENCIES))" || { echo 'DEV_DEPENDENCIES is empty for this service' >&2; exit 2; }
	@mkdir -p "$(QHPRO_DOCKER_CONFIG)"
	@DOCKER_CONFIG="$(QHPRO_DOCKER_CONFIG)" docker compose \
		--env-file "$(DEVELOPMENT_ENV_FILE)" \
		-f "$(DEVELOPMENT_COMPOSE_FILE)" \
		stop $(DEV_DEPENDENCIES)

.PHONY: deps-up deps-down
