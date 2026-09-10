# Shared Make contract for direct service work. Keep this file data-only: root
# Make remains the orchestration owner, while service Makefiles remain usable
# from an IDE terminal without inheriting a stale shell GOROOT.
include ../shared/code/development/toolchain.versions

QHPRO_TOOLCHAIN_ROOT ?= $(HOME)/.cache/qhpro/toolchain
QHPRO_GOROOT := $(QHPRO_TOOLCHAIN_ROOT)/go$(GO_VERSION)
QHPRO_TOOL_BIN := $(QHPRO_TOOLCHAIN_ROOT)/bin
QHPRO_PROTOC_ROOT := $(QHPRO_TOOLCHAIN_ROOT)/protoc$(PROTOC_VERSION)

export GOWORK := off

# Prefer the pinned cache toolchain when present; otherwise keep the host Go
# (air / IDE tasks on macOS often do not have ~/.cache/qhpro/toolchain yet).
ifeq ($(wildcard $(QHPRO_GOROOT)/bin/go),)
export PATH := $(QHPRO_TOOL_BIN):$(PATH)
else
export GOROOT := $(QHPRO_GOROOT)
export GOTOOLCHAIN := local
export PATH := $(QHPRO_TOOL_BIN):$(QHPRO_PROTOC_ROOT)/bin:$(QHPRO_GOROOT)/bin:$(PATH)
endif
