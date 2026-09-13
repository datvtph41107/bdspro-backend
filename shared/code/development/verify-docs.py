#!/usr/bin/env python3
from __future__ import annotations

import re
from pathlib import Path

ROOT = Path(__file__).resolve().parents[3]
CANONICAL = [
    ROOT / 'README.md',
    ROOT / 'CONTRIBUTING.md',
    ROOT / 'BACKEND-ARCHITECTURE.md',
    ROOT / 'BACKEND-ACCEPTANCE.md',
    ROOT / 'documents' / 'DEVELOPER-OPERATING-GUIDE.md',
    ROOT / 'documents' / 'DEVELOPMENT-ARCHITECTURE-STANDARD.md',
    ROOT / 'documents' / 'SERVICE-MAP.md',
    ROOT / 'documents' / 'MIGRATIONS.md',
    ROOT / 'documents' / 'OPERATIONAL-MAPS.md',
    ROOT / 'documents' / 'FINAL-ACCEPTANCE-20260904.md',
    ROOT / 'documents' / 'IMPORT-CHECKLIST.md',
    ROOT / 'documents' / 'CHANGE-SUMMARY-20260904.md',
]
CANONICAL.extend(sorted(ROOT.glob('*-service/README.md')))

MIGRATION_DOC_SERVICES = (
    'user-service',
    'payment-service',
    'tqd-service',
    'notification-service',
    'file-service',
    'hub-service',
    'bdspro-service',
    'crm-service',
)
MIGRATION_AUTHORITY_SERVICES = (
    'user-service',
    'organization-service',
    'payment-service',
    'tqd-service',
    'notification-service',
    'file-service',
    'hub-service',
    'bdspro-service',
    'crm-service',
)
SERVICE_MAP_MIGRATION_SERVICES = (
    'user-service',
    'organization-service',
    'payment-service',
    'tqd-service',
    'notification-service',
    'file-service',
    'hub-service',
)
REQUIRED_MIGRATION_COMMANDS = (
    'make migrate',
    'make rollback',
    'make migration-version',
    'make migration name=<schema_change>',
)
LEGACY_MIGRATION_COMMANDS = (
    'make migrateup',
    'make migrateup1',
    'make migratedown',
    'make migratedown1',
    'make migrate-version',
    'make new_migration',
    'make migrate-up',
    'make migrate-down',
    'make migrate-create',
)

link_re = re.compile(r'(?<!!)\[[^\]]+\]\(([^)]+)\)')
migration_dir_re = re.compile(r'(?m)^MIGRATION_DIR\s*(?::=|\?=|=)\s*([^\s#]+)\s*$')
errors: list[str] = []
checked = 0
for doc in CANONICAL:
    if not doc.exists():
        errors.append(f'missing canonical doc: {doc.relative_to(ROOT)}')
        continue
    checked += 1
    text = doc.read_text(encoding='utf-8')
    for raw in link_re.findall(text):
        target = raw.strip().split()[0].strip('<>')
        if not target or target.startswith(('#', 'http://', 'https://', 'mailto:')):
            continue
        path_text = target.split('#', 1)[0]
        if not path_text:
            continue
        resolved = (doc.parent / path_text).resolve()
        try:
            resolved.relative_to(ROOT.resolve())
        except ValueError:
            errors.append(f'{doc.relative_to(ROOT)}: link escapes repo: {target}')
            continue
        if not resolved.exists():
            errors.append(f'{doc.relative_to(ROOT)}: broken link: {target}')

migration_dirs: dict[str, str] = {}
for service in MIGRATION_AUTHORITY_SERVICES:
    makefile = ROOT / service / 'Makefile'
    if not makefile.exists():
        errors.append(f'missing migration owner Makefile: {makefile.relative_to(ROOT)}')
        continue
    match = migration_dir_re.search(makefile.read_text(encoding='utf-8'))
    if not match:
        errors.append(f'{makefile.relative_to(ROOT)}: missing explicit MIGRATION_DIR owner')
        continue
    migration_dirs[service] = match.group(1).rstrip('/')

migration_guide = ROOT / 'documents' / 'MIGRATIONS.md'
if migration_guide.exists():
    migration_text = migration_guide.read_text(encoding='utf-8')
    for service, directory in migration_dirs.items():
        expected = f'{service}/{directory}'
        if expected not in migration_text:
            errors.append(f'{migration_guide.relative_to(ROOT)}: missing migration authority {expected}')
    for stale in ('bdspro-service/infra/db/migrations/v2', 'crm-service/infra/db/migrate_v2'):
        if stale in migration_text:
            errors.append(f'{migration_guide.relative_to(ROOT)}: stale migration authority: {stale}')

for service in MIGRATION_DOC_SERVICES:
    doc = ROOT / service / 'README.md'
    if not doc.exists():
        continue
    text = doc.read_text(encoding='utf-8')
    directory = migration_dirs.get(service)
    if directory and f'`{directory}/`' not in text:
        errors.append(f'{doc.relative_to(ROOT)}: migration path does not match Makefile owner {directory}/')
    for command in REQUIRED_MIGRATION_COMMANDS:
        if not re.search(rf'(?m)^{re.escape(command)}$', text):
            errors.append(f'{doc.relative_to(ROOT)}: missing canonical migration command: {command}')
    for legacy in LEGACY_MIGRATION_COMMANDS:
        if legacy in text:
            errors.append(f'{doc.relative_to(ROOT)}: advertises compatibility-only migration command: {legacy}')
    for stale in ('`migrate/`', 'infra/db/migrations/v2/', 'infra/db/migrate_v2/'):
        if stale in text:
            errors.append(f'{doc.relative_to(ROOT)}: stale migration path: {stale}')

service_map = ROOT / 'documents' / 'SERVICE-MAP.md'
if service_map.exists():
    lines = service_map.read_text(encoding='utf-8').splitlines()
    for service in SERVICE_MAP_MIGRATION_SERVICES:
        line = next((candidate for candidate in lines if candidate.startswith(f'| `{service}` |')), None)
        if line is None:
            errors.append(f'{service_map.relative_to(ROOT)}: missing service row for {service}')
            continue
        directory = migration_dirs.get(service)
        if directory and f'`{directory}/`' not in line:
            errors.append(f'{service_map.relative_to(ROOT)}: {service} migration path does not match Makefile owner {directory}/')

if errors:
    raise SystemExit('\n'.join(errors))
print(f'canonical markdown links + migration developer surface PASS ({checked} docs checked)')
