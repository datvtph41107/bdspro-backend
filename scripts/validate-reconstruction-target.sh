#!/usr/bin/env bash
set -euo pipefail

fail() {
  echo "FAIL: $*" >&2
  exit 1
}

pass() {
  echo "PASS: $*"
}

require_file() {
  [[ -f "$1" ]] || fail "missing file: $1"
}

for file in \
  shared/common/businessoperation/code.go \
  shared/common/access/result.go \
  shared/common/requestgrpc/interceptor.go \
  shared/common/identitygrpc/interceptor.go \
  gateway-service/internal/grpcmetadata/from_http.go \
  user-service/internal/access/service.go \
  tqd-service/internal/report/service.go \
  tqd-service/internal/quota/service.go \
  tqd-service/infra/redisquota/reserve.go \
  tqd-service/internal/usage/event.go \
  tqd-service/internal/reportjob/worker.go \
  shared/protobuf/schema/user/access_internal.proto \
  shared/protobuf/schema/file/internal.proto \
  shared/protobuf/schema/buf.internal.gen.yaml \
  scripts/generate-internal-protobuf.sh \
  tqd-service/infra/db/migrations/v3/000009_report_quota_target.up.sql \
  tqd-service/infra/db/migrations/v3/000009_report_quota_target.down.sql; do
  require_file "$file"
done
pass "target source files exist"

core_paths=(
  shared/common/businessoperation
  shared/common/access
  user-service/internal/access
  tqd-service/internal/report
  tqd-service/internal/quota
  tqd-service/internal/usage
  tqd-service/internal/reportjob
  tqd-service/internal/quotacleanup
  tqd-service/internal/usagecheck
)

if grep -R -n -E '"net/http"|gin-gonic/gin' "${core_paths[@]}" --include='*.go'; then
  fail "HTTP dependency found below gateway boundary"
fi
pass "target business packages do not import HTTP/Gin"

if grep -R -n -E 'go-redis|redis/go-redis|gorm.io|google.golang.org/grpc' \
  user-service/internal/access tqd-service/internal/quota shared/common/access shared/common/businessoperation \
  --include='*.go'; then
  fail "DB/Redis/gRPC dependency leaked into core access/quota/operation packages"
fi
pass "core access/quota/operation packages stay independent from DB/Redis/gRPC"

if grep -n 'google.api.http' shared/protobuf/schema/user/access_internal.proto shared/protobuf/schema/file/internal.proto; then
  fail "internal gRPC proto exposes HTTP annotation"
fi
pass "internal service proto is gRPC-only"

if [[ -f shared/protobuf/types/user/access_internal.pb.gw.go || \
      -f shared/protobuf/types/file/internal.pb.gw.go ]]; then
  fail "internal gRPC proto generated an HTTP gateway handler"
fi
pass "internal service generated code has no HTTP gateway handler"

if grep -R -n -E 'x-api-key|planid|planfrom' gateway-service/internal/grpcmetadata --include='*.go' | grep -v '_test.go'; then
  fail "gateway target metadata forwards legacy commercial/raw credential fields"
fi
pass "gateway target metadata forwards canonical identity only"

grep -q 'usesTargetQuota(targetReportEnabled, info.FullMethod)' tqd-service/cmd/grpc.go \
  || fail "target report does not bypass legacy quota"
grep -q 'targetIngress(trustConfig, targetReportEnabled)' tqd-service/cmd/grpc.go \
  || fail "target TQD ingress is not wired"
grep -q 'skipTargetOperationRuntime(operationRuntime.UnaryServerInterceptor, targetReportEnabled)' tqd-service/cmd/grpc.go \
  || fail "legacy operation bridge is not feature-flag aware"
grep -q 'accessIngress(trustConfig)' user-service/cmd/grpc/main.go \
  || fail "target user access ingress is not wired"
pass "target RPCs bypass legacy identity/quota paths"

grep -q 'HasUsageForReservation(ctx, reservation.ID)' tqd-service/internal/quotacleanup/service.go \
  || fail "quota cleanup does not bind durable evidence to the exact reservation"
grep -q 'CompareAndSetUsed(' tqd-service/internal/usagecheck/service.go \
  || fail "usage reconciliation can overwrite a concurrent quota commit"
grep -q 'r.releaseStale(ctx)' tqd-service/internal/reportjob/runner.go \
  || fail "report runner does not release stale running jobs"
pass "background recovery is reservation-safe and concurrency-aware"

grep -q 'runtime.WithMetadata(_grpcmetadata.DuringReportMigration(' gateway-service/cmd/http/main.go \
  || fail "gateway does not use migration-safe HTTP-to-gRPC metadata adapter"
grep -q '_middleware.InjectGrpcMetadataContextMiddleware' gateway-service/cmd/http/main.go \
  || fail "gateway migration adapter does not preserve legacy metadata"
pass "gateway HTTP-to-gRPC boundary scopes target metadata to the migrated route"

if grep -R -n 'net/http' tqd-service/infra/filegrpc --include='*.go'; then
  fail "target TQD file client still uses HTTP"
fi
pass "target TQD file client is gRPC-only"

grep -q 'fileGRPCTargetEnabled()' file-service/main.go \
  || fail "internal file gRPC server is not wired behind its migration flag"
pass "internal file gRPC server is wired behind a default-off flag"

printf '\nTarget architecture static checks completed.\n'
