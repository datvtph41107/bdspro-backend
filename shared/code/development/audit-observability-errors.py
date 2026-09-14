#!/usr/bin/env python3
"""Inventory backend observability/error-response debt without changing source.

This is deliberately a measurement tool first. It emits stable TSV + JSON
artifacts and supports narrow zero-debt ratchets only after a migration scope
has a canonical owner and has actually reached zero.
"""

from __future__ import annotations

import argparse
import collections
import json
import re
from dataclasses import dataclass
from pathlib import Path


ROOT = Path(__file__).resolve().parents[3]
DEFAULT_OUTPUT = ROOT / ".tmp" / "observability-errors" / "inventory.tsv"
DEFAULT_SUMMARY = ROOT / ".tmp" / "observability-errors" / "summary.json"

SKIP_DIRS = {
    ".git",
    ".tmp",
    "vendor",
    "node_modules",
    "docs",
    "swagger-doc",
    "__pycache__",
}
SKIP_PREFIXES = (
    "shared/protobuf/",
)
SKIP_SUFFIXES = (
    "_test.go",
    ".pb.go",
    ".pb.gw.go",
    "wire_gen.go",
)

# Ratchets are intentionally narrow. Do not baseline arbitrary legacy debt here:
# add a scope only after its canonical migration has reached zero.
ZERO_RATCHETS = (
    ("go.direct_http_error_writer", "gateway-service"),
)


@dataclass(frozen=True)
class Rule:
    category: str
    severity: str
    suffixes: tuple[str, ...]
    pattern: re.Pattern[str]
    excluded_prefixes: tuple[str, ...] = ()
    required_path_fragments: tuple[str, ...] = ()


RULES = (
    Rule(
        "go.legacy_std_log",
        "debt",
        (".go",),
        re.compile(r"\blog\.(?:Print|Printf|Println|Fatal|Fatalf|Fatalln|Panic|Panicf|Panicln|New)\s*\("),
    ),
    Rule(
        "go.direct_slog_construction",
        "debt",
        (".go",),
        re.compile(r"\bslog\.(?:New|NewJSONHandler|NewTextHandler)\s*\("),
        excluded_prefixes=("shared/common/logging/",),
    ),
    Rule(
        "go.third_party_logger",
        "debt",
        (".go",),
        re.compile(r"(?:\bflogging\b|go\.uber\.org/zap|github\.com/rs/zerolog)"),
    ),
    Rule(
        "go.text_error_classification",
        "debt",
        (".go",),
        re.compile(r"strings\.(?:Contains|HasPrefix|HasSuffix)\s*\(\s*[A-Za-z0-9_\.]+\.Error\(\)"),
    ),
    Rule(
        "go.grpc_status_emitter",
        "review",
        (".go",),
        re.compile(r"\bstatus\.(?:Error|Errorf|New|Newf)\s*\("),
        excluded_prefixes=("shared/common/fault/",),
    ),
    Rule(
        "go.transport_error_in_domain_or_usecase",
        "debt",
        (".go",),
        re.compile(r"\bstatus\.(?:Error|Errorf|New|Newf)\s*\("),
        required_path_fragments=("/internal/usecase/", "/internal/domain/", "/domain/"),
    ),
    Rule(
        "go.direct_http_error_writer",
        "debt",
        (".go",),
        re.compile(r"(?:AbortWithStatusJSON\s*\(|\bhttp\.Error\s*\()"),
        excluded_prefixes=(
            "gateway-service/internal/httperror/",
            "gateway-service/internal/httpresponse/",
        ),
    ),
    Rule(
        "go.manual_json_response_review",
        "review",
        (".go",),
        re.compile(r"(?:json\.NewEncoder\s*\([^\n]*\)\.Encode\s*\(|\.JSON\s*\()"),
        excluded_prefixes=(
            "gateway-service/internal/httperror/",
            "gateway-service/internal/httpresponse/",
        ),
    ),
    Rule(
        "go.legacy_shared_error_response",
        "debt",
        (".go",),
        re.compile(r"(?:sharepb|sharedpb|commonpb)\.ErrorResponse\b"),
        excluded_prefixes=("gateway-service/internal/httperror/",),
    ),
    Rule(
        "go.process_print_review",
        "review",
        (".go",),
        re.compile(r"\bfmt\.(?:Print|Printf|Println)\s*\("),
    ),
    Rule(
        "python.logging_configuration",
        "review",
        (".py",),
        re.compile(r"\blogging\.(?:basicConfig|getLogger)\s*\("),
    ),
    Rule(
        "python.print_output_review",
        "review",
        (".py",),
        re.compile(r"(?<![A-Za-z0-9_])print\s*\("),
    ),
)


def should_scan(path: Path) -> bool:
    rel = path.relative_to(ROOT).as_posix()
    if any(part in SKIP_DIRS for part in path.relative_to(ROOT).parts):
        return False
    if rel.startswith(SKIP_PREFIXES):
        return False
    if rel.endswith(SKIP_SUFFIXES):
        return False
    return path.suffix in {".go", ".py"}


def owner_for(rel: str) -> str:
    first = rel.split("/", 1)[0]
    if first.endswith("-service"):
        return first
    if rel.startswith("shared/common/"):
        return "shared/common"
    if rel.startswith("shared/code/"):
        return "shared/code"
    return first


def rule_applies(rule: Rule, rel: str, path: Path) -> bool:
    if not any(path.name.endswith(suffix) for suffix in rule.suffixes):
        return False
    if rel.startswith(rule.excluded_prefixes):
        return False
    if rule.required_path_fragments and not any(fragment in f"/{rel}" for fragment in rule.required_path_fragments):
        return False
    return True


def line_number(text: str, index: int) -> int:
    return text.count("\n", 0, index) + 1


def excerpt_at(text: str, line: int) -> str:
    lines = text.splitlines()
    if not lines or line < 1 or line > len(lines):
        return ""
    return " ".join(lines[line - 1].strip().split())[:240]


def match_is_comment(path: Path, text: str, index: int) -> bool:
    """Ignore comment-only historical examples while preserving source line numbers.

    This is intentionally conservative rather than a language parser. It removes
    the common false positives produced by commented-out Go/Python statements;
    executable-looking text in strings is still reviewable inventory.
    """
    line_start = text.rfind("\n", 0, index) + 1
    prefix = text[line_start:index].lstrip()

    if path.suffix == ".go":
        if prefix.startswith("//"):
            return True
        if text.rfind("/*", 0, index) > text.rfind("*/", 0, index):
            return True
    elif path.suffix == ".py" and prefix.startswith("#"):
        return True

    return False


def scan() -> list[dict[str, object]]:
    findings: list[dict[str, object]] = []
    for path in sorted(ROOT.rglob("*")):
        if not path.is_file() or not should_scan(path):
            continue
        rel = path.relative_to(ROOT).as_posix()
        try:
            text = path.read_text(encoding="utf-8")
        except UnicodeDecodeError:
            continue
        for rule in RULES:
            if not rule_applies(rule, rel, path):
                continue
            seen_lines: set[int] = set()
            for match in rule.pattern.finditer(text):
                if match_is_comment(path, text, match.start()):
                    continue
                line = line_number(text, match.start())
                if line in seen_lines:
                    continue
                seen_lines.add(line)
                findings.append(
                    {
                        "category": rule.category,
                        "severity": rule.severity,
                        "owner": owner_for(rel),
                        "path": rel,
                        "line": line,
                        "excerpt": excerpt_at(text, line),
                    }
                )
    return findings


def write_tsv(path: Path, findings: list[dict[str, object]]) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open("w", encoding="utf-8", newline="") as handle:
        handle.write("category\tseverity\towner\tpath\tline\texcerpt\n")
        for item in findings:
            excerpt = str(item["excerpt"]).replace("\t", " ").replace("\n", " ")
            handle.write(
                f'{item["category"]}\t{item["severity"]}\t{item["owner"]}\t'
                f'{item["path"]}\t{item["line"]}\t{excerpt}\n'
            )


def ratchet_counts(findings: list[dict[str, object]]) -> dict[str, int]:
    counts: dict[str, int] = {}
    for category, owner in ZERO_RATCHETS:
        key = f"{category}@{owner}"
        counts[key] = sum(
            1
            for item in findings
            if item["severity"] == "debt"
            and item["category"] == category
            and item["owner"] == owner
        )
    return counts


def build_summary(findings: list[dict[str, object]]) -> dict[str, object]:
    by_category = collections.Counter(str(item["category"]) for item in findings)
    debt_by_category = collections.Counter(
        str(item["category"]) for item in findings if item["severity"] == "debt"
    )
    by_owner = collections.Counter(str(item["owner"]) for item in findings)
    debt_by_owner = collections.Counter(
        str(item["owner"]) for item in findings if item["severity"] == "debt"
    )
    return {
        "schema_version": 1,
        "mode": "inventory",
        "total_findings": len(findings),
        "total_debt_findings": sum(debt_by_category.values()),
        "by_category": dict(sorted(by_category.items())),
        "debt_by_category": dict(sorted(debt_by_category.items())),
        "by_owner": dict(sorted(by_owner.items())),
        "debt_by_owner": dict(sorted(debt_by_owner.items())),
        "ratchets": ratchet_counts(findings),
        "policy": {
            "canonical_go_logging_api": "log/slog",
            "canonical_go_logging_owner": "shared/common/logging",
            "canonical_error_identity_owner": "shared/common/fault",
            "canonical_gateway_error_serializer": "gateway-service/internal/httpresponse.WriteProblem",
            "legacy_findings_fail_ci": False,
            "ratchets_fail_ci": True,
        },
    }


def print_summary(summary: dict[str, object]) -> None:
    print("Observability/error inventory")
    print(f'  findings: {summary["total_findings"]}')
    print(f'  debt:     {summary["total_debt_findings"]}')
    print("\nDebt by category:")
    for category, count in summary["debt_by_category"].items():
        print(f"  {category:42} {count:5}")
    print("\nDebt by owner:")
    for owner, count in summary["debt_by_owner"].items():
        print(f"  {owner:42} {count:5}")
    print("\nZero-debt ratchets:")
    for scope, count in summary["ratchets"].items():
        print(f"  {scope:60} {count:5}")


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--output", type=Path, default=DEFAULT_OUTPUT)
    parser.add_argument("--summary", type=Path, default=DEFAULT_SUMMARY)
    parser.add_argument("--fail-on-debt", action="store_true")
    parser.add_argument("--enforce-ratchets", action="store_true")
    args = parser.parse_args()

    findings = scan()
    summary = build_summary(findings)
    write_tsv(args.output, findings)
    args.summary.parent.mkdir(parents=True, exist_ok=True)
    args.summary.write_text(json.dumps(summary, indent=2, sort_keys=True) + "\n", encoding="utf-8")
    print_summary(summary)
    print(f"\nTSV:     {args.output.relative_to(ROOT)}")
    print(f"Summary: {args.summary.relative_to(ROOT)}")

    if args.enforce_ratchets:
        violations = {scope: count for scope, count in summary["ratchets"].items() if count}
        if violations:
            print("\nRatchet violation(s):")
            for scope, count in violations.items():
                print(f"  {scope}: {count}")
            return 1

    if args.fail_on_debt and summary["total_debt_findings"]:
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
