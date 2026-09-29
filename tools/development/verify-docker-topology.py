#!/usr/bin/env python3
from pathlib import Path
import re
import shlex
import sys

ROOT = Path(__file__).resolve().parents[2]

GO_DOCKERFILES = (
    "api/gateway/Dockerfile",
    "internal/authentication/Dockerfile",
    "internal/user/Dockerfile",
    "internal/organization-legacy/Dockerfile",
    "internal/payment/Dockerfile",
    "internal/planning/Dockerfile",
    "internal/notification/Dockerfile",
    "internal/file/Dockerfile",
    "internal/hub/Dockerfile",
    "internal/assistant/Dockerfile",
)

REPO_PREFIXES = ("api/", "internal/", "infrastructure/", "proto")
errors: list[str] = []

compose = ROOT / "compose.yaml"
if not compose.exists():
    errors.append("compose.yaml missing")
else:
    for number, line in enumerate(compose.read_text(encoding="utf-8").splitlines(), 1):
        for match in re.finditer(r"dockerfile:\s*([^,}\]\s]+)", line):
            raw = match.group(1).strip().strip("'\"")
            if not (ROOT / raw).is_file():
                errors.append(f"compose.yaml:{number}: dockerfile does not exist: {raw}")

for relative in GO_DOCKERFILES:
    path = ROOT / relative
    if not path.is_file():
        errors.append(f"missing Dockerfile: {relative}")
        continue

    stage = 0
    for number, raw in enumerate(path.read_text(encoding="utf-8").splitlines(), 1):
        line = raw.strip()
        if not line or line.startswith("#"):
            continue
        if line.upper().startswith("FROM "):
            stage += 1
            continue
        if stage != 1:
            continue

        if line.upper().startswith("WORKDIR "):
            workdir = line.split(None, 1)[1].strip()
            if workdir.startswith("/src/"):
                repo_rel = workdir[len("/src/"):]
                if not (ROOT / repo_rel).is_dir():
                    errors.append(f"{relative}:{number}: build WORKDIR has no repo source: {workdir}")
            continue

        if not line.upper().startswith("COPY ") or "--from=" in line:
            continue

        try:
            parts = shlex.split(line)
        except ValueError as exc:
            errors.append(f"{relative}:{number}: cannot parse COPY: {exc}")
            continue

        tokens = parts[1:]
        while tokens and tokens[0].startswith("--"):
            tokens.pop(0)
        if len(tokens) < 2:
            errors.append(f"{relative}:{number}: malformed COPY")
            continue

        sources = tokens[:-1]
        destination = tokens[-1]
        for source in sources:
            if any(ch in source for ch in "*?["):
                continue
            candidate = ROOT / source.rstrip("/")
            if not candidate.exists():
                errors.append(f"{relative}:{number}: COPY source missing: {source}")

        if len(sources) == 1:
            source = sources[0].rstrip("/")
            if source == "proto" or source.startswith(REPO_PREFIXES):
                if destination.startswith("./"):
                    expected = "./" + source
                    if destination.rstrip("/") != expected:
                        errors.append(
                            f"{relative}:{number}: build COPY must preserve repo topology: "
                            f"{source} -> {destination}, expected {expected}"
                        )

if errors:
    print("\n".join(errors), file=sys.stderr)
    raise SystemExit(1)

print(f"docker topology verification PASS ({len(GO_DOCKERFILES)} Go Dockerfiles + compose references)")
