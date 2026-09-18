#!/usr/bin/env python3
from __future__ import annotations

import contextlib
import importlib.util
import io
import json
import sys
import tempfile

sys.dont_write_bytecode = True
import unittest
from pathlib import Path


SCRIPT = Path(__file__).with_name("query-structured-logs.py")
SPEC = importlib.util.spec_from_file_location("bdspro_query_structured_logs", SCRIPT)
MODULE = importlib.util.module_from_spec(SPEC)
sys.modules[SPEC.name] = MODULE
assert SPEC.loader is not None
SPEC.loader.exec_module(MODULE)


class StructuredLogQueryTest(unittest.TestCase):
    def write_runtime(
        self,
        root: Path,
        service: str,
        run: str,
        records: list[dict[str, object]],
    ) -> Path:
        path = root / service / "runs" / run / "runtime.jsonl"
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(
            "".join(json.dumps(record) + "\n" for record in records),
            encoding="utf-8",
        )
        return path

    def test_runtime_files_read_only_canonical_runtime(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            runtime = self.write_runtime(
                root,
                "payment-service",
                "run-1",
                [{"time": "2026-09-18T02:00:00Z", "msg": "canonical"}],
            )
            channel = runtime.parent / "channels" / "http.jsonl"
            channel.parent.mkdir(parents=True)
            channel.write_text(
                json.dumps({"time": "2026-09-18T02:00:01Z", "msg": "projection"})
                + "\n",
                encoding="utf-8",
            )
            self.assertEqual(MODULE.runtime_files(root), [runtime])

    def test_exact_filters_sort_and_limit(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            self.write_runtime(
                root,
                "payment-service",
                "run-b",
                [
                    {
                        "time": "2026-09-18T02:00:03Z",
                        "level": "ERROR",
                        "msg": "late",
                        "service.name": "payment-service",
                        "request_id": "req-1",
                        "operation_id": "op-1",
                        "component": "worker",
                    },
                    {
                        "time": "2026-09-18T02:00:01Z",
                        "level": "INFO",
                        "msg": "early",
                        "service.name": "payment-service",
                        "request_id": "req-1",
                        "operation_id": "op-1",
                        "event_name": "payment.started",
                    },
                ],
            )
            self.write_runtime(
                root,
                "user-service",
                "run-a",
                [
                    {
                        "time": "2026-09-18T02:00:02Z",
                        "level": "WARN",
                        "msg": "other",
                        "service.name": "user-service",
                        "request_id": "req-10",
                    }
                ],
            )

            filters = MODULE.Filters(
                request_id="req-1",
                service="payment-service",
            )
            found = MODULE.query(root, filters, limit=1)
            self.assertEqual(len(found), 1)
            self.assertEqual(found[0][0]["msg"], "late")

            exact_miss = MODULE.query(
                root,
                MODULE.Filters(request_id="req-10", service="payment-service"),
                limit=10,
            )
            self.assertEqual(exact_miss, [])

    def test_level_filter_is_case_insensitive_but_other_fields_are_exact(self) -> None:
        record = {
            "level": "WARN",
            "service.name": "user-service",
            "operation_id": "op-2",
            "component": "grpc",
        }
        self.assertTrue(
            MODULE.matches(
                record,
                MODULE.Filters(
                    level="warn",
                    service="user-service",
                    operation_id="op-2",
                    component="grpc",
                ),
            )
        )
        self.assertFalse(
            MODULE.matches(
                record,
                MODULE.Filters(service="user"),
            )
        )

    def test_malformed_tail_warns_and_valid_records_survive(self) -> None:
        with tempfile.TemporaryDirectory() as tmp:
            root = Path(tmp)
            runtime = self.write_runtime(
                root,
                "auth-service",
                "run-1",
                [
                    {
                        "time": "2026-09-18T02:00:00Z",
                        "level": "INFO",
                        "msg": "valid",
                        "service.name": "auth-service",
                        "request_id": "req-auth",
                    }
                ],
            )
            with runtime.open("a", encoding="utf-8") as handle:
                handle.write('{"time":')

            stderr = io.StringIO()
            with contextlib.redirect_stderr(stderr):
                found = MODULE.query(
                    root,
                    MODULE.Filters(request_id="req-auth"),
                    limit=10,
                )
            self.assertEqual([item[0]["msg"] for item in found], ["valid"])
            self.assertIn("warning: skip malformed JSON", stderr.getvalue())


if __name__ == "__main__":
    unittest.main()
