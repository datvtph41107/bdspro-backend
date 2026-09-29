#!/usr/bin/env python3
from pathlib import Path
import re
import sys

ROOT = Path(__file__).resolve().parents[2]

OLD_ROOTS = (
    "gateway-service", "user-service", "auth-service", "organization-service",
    "bdspro-service", "crm-service", "tqd-service", "payment-service",
    "file-service", "notification-service", "hub-service", "assistant-service",
    "chat-service", "chat-v1-service", "social-service", "map-service",
    "relay-service", "search-service", "ai-service",
)

OLD_SHARED_PREFIXES = ("shared/common", "shared/code", "shared/protobuf", "shared/base", "shared/config", "shared/scripts", "shared/summary")

REQUIRED = (
    "api/gateway",
    "internal/user",
    "internal/authentication",
    "internal/organization-legacy",
    "internal/property",
    "internal/crm",
    "internal/planning",
    "internal/payment",
    "internal/file",
    "internal/notification",
    "internal/hub",
    "internal/assistant",
    "internal/chat",
    "internal/chat-legacy",
    "internal/social",
    "internal/map-legacy",
    "internal/realtime-relay",
    "internal/search",
    "internal/ai",
    "infrastructure/runtime",
    "infrastructure/base",
    "proto",
    "tools/development",
)

EXCLUDED = {
    Path("tools/development/deploy.sh"),  # protected exact-checksum compatibility artifact
    Path("tools/development/release.sh"), # legacy remote deployment compatibility
}

def operational_files():
    fixed = [
        ROOT / "go.work",
        ROOT / "compose.yaml",
        ROOT / "Makefile",
        ROOT / "bdspro.code-workspace",
        ROOT / ".gitignore",
    ]
    for p in fixed:
        if p.is_file():
            yield p

    for base in (ROOT / ".github", ROOT / ".vscode"):
        if base.exists():
            yield from (p for p in base.rglob("*") if p.is_file())

    for base in (ROOT / "api", ROOT / "internal", ROOT / "infrastructure", ROOT / "proto"):
        if not base.exists():
            continue
        for p in base.rglob("*"):
            if not p.is_file():
                continue
            if p.name == "go.mod" or p.name == "Makefile" or p.name.startswith("Dockerfile"):
                yield p

    for base in (ROOT / "tools" / "development", ROOT / "tools" / "scripts"):
        if not base.exists():
            continue
        for p in base.rglob("*"):
            if not p.is_file():
                continue
            rel = p.relative_to(ROOT)
            if rel in EXCLUDED or "legacy" in rel.parts or "manual" in rel.parts:
                continue
            if p.name == "Makefile" or p.name.startswith("Dockerfile") or p.suffix in {".mk", ".sh", ".py", ".yml", ".yaml"}:
                yield p

def stale_path(line: str, old: str) -> bool:
    escaped = re.escape(old)
    checks = (
        rf'(?<![A-Za-z0-9_-])(?:\.\./|\./)?{escaped}/',
        rf'^\s*working-directory:\s*{escaped}\s*$',
        rf'^\s*(?:COPY|ADD)\s+(?:--\S+\s+)*{escaped}(?:\s|$)',
        rf'\bcd\s+(?:\.\./|\./)?{escaped}(?:\s|$|&&|;)',
    )
    return any(re.search(pattern, line) for pattern in checks)

errors = []

for old in OLD_ROOTS:
    if (ROOT / old).exists():
        errors.append(f"legacy top-level source root still exists: {old}")
if (ROOT / "shared").exists():
    errors.append("legacy top-level source root still exists: shared")

for required in REQUIRED:
    if not (ROOT / required).exists():
        errors.append(f"required migrated source root missing: {required}")

seen = set()
for path in operational_files():
    path = path.resolve()
    if path in seen:
        continue
    seen.add(path)
    rel = path.relative_to(ROOT)
    try:
        lines = path.read_text(encoding="utf-8").splitlines()
    except (UnicodeDecodeError, OSError):
        continue
    for number, line in enumerate(lines, 1):
        for old in OLD_ROOTS:
            if stale_path(line, old):
                errors.append(f"{rel}:{number}: stale physical path {old!r}: {line.strip()}")
        for old in OLD_SHARED_PREFIXES:
            if re.search(rf'(?<![A-Za-z0-9_-])(?:\.\./|\./)?{re.escape(old)}(?:/|\s|$)', line):
                errors.append(f"{rel}:{number}: stale physical path {old!r}: {line.strip()}")

if errors:
    print("\n".join(errors), file=sys.stderr)
    raise SystemExit(1)

print(f"source topology verification PASS ({len(seen)} operational files checked)")
