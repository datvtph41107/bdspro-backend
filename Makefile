# BDSPro whole-repository developer API.
#
# Public workspace lifecycle:
#   make setup
#   make deps
#   make up | status | logs | smoke | test-e2e | down
#
# Daily application lifecycle belongs to its service:
#   make -C payment-service deps-up
#   make -C payment-service migrate
#   make -C payment-service rollback
#   make -C payment-service migration name=add_schema_change
#   make -C payment-service migration-version
#   make -C payment-service dev
#   make -C payment-service test
#
# Full verification/acceptance remains available to CI and release maintainers
# as `make accept`; it is not a prerequisite for opening an accepted clone.
# Implementations live in the repository engineering application below.
include shared/code/development/root.mk
