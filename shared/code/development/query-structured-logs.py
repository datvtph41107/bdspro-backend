#!/usr/bin/env python3
"""Read-only query over canonical retained BDSPro runtime.jsonl records."""

from __future__ import annotations

import argparse
import json
import sys
from dataclasses import dataclass
from pathlib import Path
from typing import Any, Iterable


@dataclass(frozen=True)
class Filters:
    request_id: str | None = None
    operation_id: str | None = None
    service: str | None = None
    event_name: str | None = None
    level: str | None = None
    component: str | None = None


def positive_int(value: str) -> int:
    parsed = int(value)
    if parsed <= 0:
        raise argparse.ArgumentTypeError("must be greater than zero")
    return parsed


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Query canonical retained runtime.jsonl records."
    )
    parser.add_argument("--root", required=True, type=Path)
    parser.add_argument("--request-id")
    parser.add_argument("--operation-id")
    parser.add_argument("--service")
    parser.add_argument("--event-name")
    parser.add_argument("--level")
    parser.add_argument("--component")
    parser.add_argument("--limit", type=positive_int, default=200)
    parser.add_argument("--format", choices=("human", "json"), default="human")
    return parser.parse_args()


def runtime_files(root: Path) -> list[Path]:
    files = list(root.glob("*/runs/*/runtime.jsonl"))
    files.extend(root.glob("*/current/runtime.jsonl"))
    return sorted(set(files))


def read_records(paths: Iterable[Path]) -> Iterable[tuple[dict[str, Any], Path, int]]:
    for path in paths:
        try:
            with path.open("r", encoding="utf-8") as handle:
                for line_number, raw_line in enumerate(handle, start=1):
                    line = raw_line.strip()
                    if not line:
                        continue
                    try:
                        record = json.loads(line)
                    except json.JSONDecodeError as exc:
                        print(
                            f"warning: skip malformed JSON {path}:{line_number}: {exc.msg}",
                            file=sys.stderr,
                        )
                        continue
                    if not isinstance(record, dict):
                        print(
                            f"warning: skip non-object JSON {path}:{line_number}",
                            file=sys.stderr,
                        )
                        continue
                    yield record, path, line_number
        except FileNotFoundError:
            # A run may disappear during local cleanup while a read-only query is scanning it.
            continue


def exact(record: dict[str, Any], key: str, expected: str | None) -> bool:
    return expected is None or str(record.get(key, "")) == expected


def matches(record: dict[str, Any], filters: Filters) -> bool:
    return (
        exact(record, "request_id", filters.request_id)
        and exact(record, "operation_id", filters.operation_id)
        and exact(record, "service.name", filters.service)
        and exact(record, "event_name", filters.event_name)
        and exact(record, "component", filters.component)
        and (
            filters.level is None
            or str(record.get("level", "")).upper() == filters.level.upper()
        )
    )


def query(
    root: Path,
    filters: Filters,
    limit: int,
) -> list[tuple[dict[str, Any], Path, int]]:
    found = [
        item
        for item in read_records(runtime_files(root))
        if matches(item[0], filters)
    ]
    found.sort(
        key=lambda item: (
            str(item[0].get("time", "")),
            str(item[0].get("service.name", "")),
            str(item[1]),
            item[2],
        )
    )
    return found[-limit:]


def run_id(path: Path) -> str:
    parts = path.parts
    if "runs" in parts:
        index = parts.index("runs")
        if index + 1 < len(parts):
            return parts[index + 1]
    return "current"


def human_line(record: dict[str, Any], path: Path) -> str:
    pieces = [
        str(record.get("time", "-")),
        str(record.get("level", "-")),
        str(record.get("service.name", "-")),
        f"run={run_id(path)}",
    ]
    for key in ("request_id", "operation_id", "event_name", "component"):
        value = record.get(key)
        if value not in (None, ""):
            pieces.append(f"{key}={value}")
    pieces.append(f"msg={record.get('msg', '')}")
    return " ".join(pieces)


def main() -> int:
    args = parse_args()
    filters = Filters(
        request_id=args.request_id,
        operation_id=args.operation_id,
        service=args.service,
        event_name=args.event_name,
        level=args.level,
        component=args.component,
    )
    for record, path, _ in query(args.root, filters, args.limit):
        if args.format == "json":
            print(json.dumps(record, ensure_ascii=False, separators=(",", ":")))
        else:
            print(human_line(record, path))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
