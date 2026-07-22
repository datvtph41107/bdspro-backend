#!/usr/bin/env python3
from pathlib import Path
import re
import sys

root = Path(__file__).resolve().parents[2]
failures: list[str] = []

def read(relative: str) -> str:
    path = root / relative
    if not path.is_file():
        failures.append(f"missing file: {relative}")
        return ""
    return path.read_text(encoding="utf-8")

def require(relative: str, pattern: str, message: str) -> None:
    if not re.search(pattern, read(relative), re.MULTILINE):
        failures.append(f"{relative}: {message}")

def reject(relative: str, pattern: str, message: str) -> None:
    if re.search(pattern, read(relative), re.MULTILINE):
        failures.append(f"{relative}: {message}")

# RC2: one public Discovery contract, Gateway HTTP -> typed TQD gRPC.
require("shared/protobuf/schema/tqd/discovery.proto", r"service\s+DiscoveryService", "typed DiscoveryService is missing")
require("tqd-service/cmd/grpc.go", r"RegisterDiscoveryServiceServer", "TQD Discovery gRPC registration is missing")
require("gateway-service/service/discovery_grpc_bridge.go", r"DiscoveryServiceClient", "Gateway TQD gRPC bridge is missing")
require("gateway-service/cmd/http/main.go", r"/v2/tqd/discovery/", "TQD Discovery public route is not public")
reject("gateway-service/cmd/http/main.go", r'"/v2/discovery/', "legacy global Discovery route remains public")
reject("search-service/main.go", r'/v2/discovery', "search-service still exposes legacy Discovery federation")
for obsolete in [
    "gateway-service/service/global_search.go",
    "gateway-service/service/global_search_test.go",
    "search-service/handlers/discovery.go",
    "search-service/discovery/service.go",
    "gateway-service/cmd/http/main_v2.go",
]:
    if (root / obsolete).exists():
        failures.append(f"obsolete file remains: {obsolete}")

# SEO-M8: auditable imports, freshness, scoring, reconciliation, gRPC-only CRM transport.
for relative in [
    "crm-service/infra/db/migrate_v2/000013_seo_m8_measurement_control.up.sql",
    "crm-service/infra/db/migrate_v2/000013_seo_m8_measurement_control.down.sql",
    "crm-service/internal/usecase/seo_measurement_m8_usecase.go",
    "crm-service/infra/handler/seo_measurement_m8_runtime_handler.go",
    "shared/protobuf/types/crm/seo_measurement_runtime_transport.go",
    "gateway-service/service/seo_measurement_runtime.go",
]:
    read(relative)
require("crm-service/internal/usecase/seo_measurement_m8_usecase.go", r"stale_measurement_sources", "source freshness is not part of scoring")
require("crm-service/internal/usecase/seo_measurement_m8_usecase.go", r"measurementImportChecksum", "idempotent import checksum is missing")
require("crm-service/internal/repo/seo_measurement_repo.go", r"GetReusableImportRunByChecksum", "idempotent import repository boundary is missing")
require("crm-service/infra/db/migrate_v2/000013_seo_m8_measurement_control.up.sql", r"uq_seo_measurement_import_run_active_checksum", "idempotent import index is missing")
require("crm-service/infra/postgre/seo_measurement_postgres.go", r"ReconcileIssues", "feedback issue reconciliation is missing")
require("crm-service/infra/postgre/seo_measurement_postgres.go", r"issue_type IN \?", "reconciliation is not restricted to detector-managed issues")
require("gateway-service/service/seo_measurement_runtime.go", r"HandlePath\(http.MethodPost, \"/v2/crm/admin/seo/measurement/m8/import/\{source\}\"", "M8 import bridge is missing")
require("crm-service/cmd/grpc/main.go", r"RegisterSeoMeasurementRuntimeServiceServer", "CRM M8 gRPC service is not registered")

# Migration symmetry and source lineage columns.
up = read("crm-service/infra/db/migrate_v2/000013_seo_m8_measurement_control.up.sql")
down = read("crm-service/infra/db/migrate_v2/000013_seo_m8_measurement_control.down.sql")
for table in ["seo_measurement_import_run", "seo_measurement_source_state", "seo_page_score_snapshot"]:
    if f"CREATE TABLE IF NOT EXISTS {table}" not in up:
        failures.append(f"migration up missing table {table}")
    if f"DROP TABLE IF EXISTS {table}" not in down:
        failures.append(f"migration down missing table {table}")
for table in ["seo_search_performance_daily", "seo_url_inspection_snapshot", "seo_traffic_daily", "seo_ranking_snapshot", "seo_crawl_event"]:
    if f"ALTER TABLE {table} ADD COLUMN IF NOT EXISTS source_run_id" not in up:
        failures.append(f"migration missing source_run_id for {table}")

# Reject transport/config sources removed by RC2 in executable Go source.
for path in root.rglob("*.go"):
    if any(part in {"vendor", ".git"} for part in path.parts):
        continue
    text = path.read_text(encoding="utf-8", errors="replace")
    if "SEARCH_HTTP_URL" in text or "DISCOVERY_PROVIDER_URLS" in text:
        failures.append(f"legacy Discovery environment variable remains in {path.relative_to(root)}")

if failures:
    print(f"Backend RC2/SEO-M8 verification failed ({len(failures)})", file=sys.stderr)
    for failure in failures:
        print(f"- {failure}", file=sys.stderr)
    raise SystemExit(1)
print("Backend RC2/SEO-M8 verification passed")
