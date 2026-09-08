# Shared Make contract for direct service work. Keep this file data-only: root
# Make remains the orchestration owner, while service Makefiles remain usable
# from an IDE terminal without inheriting a stale shell GOROOT.
include ../shared/code/development/toolchain.versions

QHPRO_TOOLCHAIN_ROOT ?= $(HOME)/.cache/qhpro/toolchain
QHPRO_GOROOT := $(QHPRO_TOOLCHAIN_ROOT)/go$(GO_VERSION)
QHPRO_TOOL_BIN := $(QHPRO_TOOLCHAIN_ROOT)/bin
QHPRO_PROTOC_ROOT := $(QHPRO_TOOLCHAIN_ROOT)/protoc$(PROTOC_VERSION)

export GOROOT := $(QHPRO_GOROOT)
export GOTOOLCHAIN := local
export GOWORK := off
export PATH := $(QHPRO_TOOL_BIN):$(QHPRO_PROTOC_ROOT)/bin:$(QHPRO_GOROOT)/bin:$(PATH)
