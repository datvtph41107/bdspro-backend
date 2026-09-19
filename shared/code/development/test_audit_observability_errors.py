#!/usr/bin/env python3

import importlib.util
import sys
import unittest
from pathlib import Path
from unittest import mock

SCRIPT = Path(__file__).with_name("audit-observability-errors.py")

spec = importlib.util.spec_from_file_location(
    "audit_observability_errors",
    SCRIPT,
)
audit = importlib.util.module_from_spec(spec)
sys.modules[spec.name] = audit

assert spec.loader is not None
spec.loader.exec_module(audit)

TEXT_RULE = next(
    rule
    for rule in audit.RULES
    if rule.category == "go.text_error_classification"
)


class TextErrorClassifierDetectorTest(unittest.TestCase):
    def classify(self, source: str) -> set[int]:
        return audit.text_error_classifier_lines(
            Path("fixture.go"),
            source,
            TEXT_RULE.pattern,
        )

    def test_direct_error_string_classifier(self):
        source = (
            "package fixture\n"
            "func f(err error) {\n"
            '\tif strings.Contains(err.Error(), "not found") {}\n'
            "}\n"
        )

        self.assertEqual(self.classify(source), {3})

    def test_alias_derived_from_error_string_is_classified(self):
        source = (
            "package fixture\n"
            "func f(err error) {\n"
            "\tmsg := err.Error()\n"
            '\tif strings.Contains(msg, "not found") {}\n'
            "}\n"
        )

        self.assertEqual(self.classify(source), {4})

    def test_reassignment_clears_error_string_taint(self):
        source = (
            "package fixture\n"
            "func f(err error) {\n"
            "\tmsg := err.Error()\n"
            '\tmsg = "safe presentation text"\n'
            '\tif strings.Contains(msg, "not found") {}\n'
            "}\n"
        )

        self.assertEqual(self.classify(source), set())

    def test_duplicate_classifiers_on_one_line_are_one_finding(self):
        source = (
            "package fixture\n"
            "func f(err error) {\n"
            "\tmsg := err.Error()\n"
            '\tif strings.Contains(msg, "required") || '
            'strings.Contains(msg, "invalid") {}\n'
            "}\n"
        )

        self.assertEqual(self.classify(source), {4})

    def test_plain_string_alias_is_not_error_classification(self):
        source = (
            "package fixture\n"
            "func f() {\n"
            '\tmsg := "not found"\n'
            '\tif strings.Contains(msg, "not found") {}\n'
            "}\n"
        )

        self.assertEqual(self.classify(source), set())

    def test_alias_remains_tainted_in_nested_scope(self):
        source = (
            "package fixture\n"
            "func f(err error, ok bool) {\n"
            "\tmsg := err.Error()\n"
            "\tif ok {\n"
            '\t\tif strings.Contains(msg, "not found") {}\n'
            "\t}\n"
            "}\n"
        )

        self.assertEqual(self.classify(source), {5})


    def test_tqd_text_classifier_zero_ratchet_is_registered(self):
        self.assertIn(
            (
                "go.text_error_classification",
                "tqd-service",
            ),
            audit.ZERO_RATCHETS,
        )

    def test_tqd_text_classifier_ratchet_enforces_regression(self):
        finding = {
            "category": "go.text_error_classification",
            "severity": "debt",
            "owner": "tqd-service",
            "path": "tqd-service/fixture.go",
            "line": 1,
            "excerpt": (
                'strings.Contains(err.Error(), "not found")'
            ),
        }

        output = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "ratchet-regression.tsv"
        )
        summary = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "ratchet-regression.json"
        )

        argv = [
            "audit-observability-errors.py",
            "--output",
            str(output),
            "--summary",
            str(summary),
            "--enforce-ratchets",
        ]

        with (
            mock.patch.object(
                audit,
                "scan",
                return_value=[finding],
            ),
            mock.patch.object(
                sys,
                "argv",
                argv,
            ),
            mock.patch("builtins.print"),
        ):
            try:
                self.assertEqual(
                    audit.main(),
                    1,
                )
            finally:
                output.unlink(missing_ok=True)
                summary.unlink(missing_ok=True)


    def test_shared_common_direct_http_writer_zero_ratchet_is_registered(
        self,
    ):
        self.assertIn(
            (
                "go.direct_http_error_writer",
                "shared/common",
            ),
            audit.ZERO_RATCHETS,
        )

    def test_shared_common_direct_http_writer_ratchet_enforces_regression(
        self,
    ):
        finding = {
            "category": "go.direct_http_error_writer",
            "severity": "debt",
            "owner": "shared/common",
            "path": (
                "shared/common/jwt/"
                "regression_fixture.go"
            ),
            "line": 1,
            "excerpt": (
                "c.AbortWithStatusJSON("
                "http.StatusUnauthorized, payload)"
            ),
        }

        output = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "shared-common-ratchet-regression.tsv"
        )

        summary = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "shared-common-ratchet-regression.json"
        )

        argv = [
            "audit-observability-errors.py",
            "--output",
            str(output),
            "--summary",
            str(summary),
            "--enforce-ratchets",
        ]

        with (
            mock.patch.object(
                audit,
                "scan",
                return_value=[finding],
            ),
            mock.patch.object(
                sys,
                "argv",
                argv,
            ),
            mock.patch(
                "builtins.print",
            ),
        ):
            try:
                self.assertEqual(
                    audit.main(),
                    1,
                )
            finally:
                output.unlink(
                    missing_ok=True,
                )
                summary.unlink(
                    missing_ok=True,
                )



    def test_user_service_legacy_shared_error_response_zero_ratchet_is_registered(
        self,
    ):
        self.assertIn(
            (
                "go.legacy_shared_error_response",
                "user-service",
            ),
            audit.ZERO_RATCHETS,
        )

    def test_user_service_legacy_shared_error_response_ratchet_enforces_regression(
        self,
    ):
        finding = {
            "category": "go.legacy_shared_error_response",
            "severity": "debt",
            "owner": "user-service",
            "path": "user-service/fixture.go",
            "line": 1,
            "excerpt": "&sharepb.ErrorResponse{Code: 404}",
        }

        output = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "user-service-legacy-error-response-ratchet.tsv"
        )

        summary = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "user-service-legacy-error-response-ratchet.json"
        )

        argv = [
            "audit-observability-errors.py",
            "--output",
            str(output),
            "--summary",
            str(summary),
            "--enforce-ratchets",
        ]

        with (
            mock.patch.object(
                audit,
                "scan",
                return_value=[finding],
            ),
            mock.patch.object(
                sys,
                "argv",
                argv,
            ),
            mock.patch(
                "builtins.print",
            ),
        ):
            try:
                self.assertEqual(
                    audit.main(),
                    1,
                )
            finally:
                output.unlink(
                    missing_ok=True,
                )
                summary.unlink(
                    missing_ok=True,
                )



    def test_bdspro_service_legacy_shared_error_response_zero_ratchet_is_registered(
        self,
    ):
        self.assertIn(
            (
                "go.legacy_shared_error_response",
                "bdspro-service",
            ),
            audit.ZERO_RATCHETS,
        )

    def test_bdspro_service_legacy_shared_error_response_ratchet_enforces_regression(
        self,
    ):
        finding = {
            "category": "go.legacy_shared_error_response",
            "severity": "debt",
            "owner": "bdspro-service",
            "path": "bdspro-service/fixture.go",
            "line": 1,
            "excerpt": "&sharepb.ErrorResponse{Code: 400}",
        }

        output = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "bdspro-service-legacy-error-response-ratchet.tsv"
        )

        summary = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "bdspro-service-legacy-error-response-ratchet.json"
        )

        argv = [
            "audit-observability-errors.py",
            "--output",
            str(output),
            "--summary",
            str(summary),
            "--enforce-ratchets",
        ]

        with (
            mock.patch.object(
                audit,
                "scan",
                return_value=[finding],
            ),
            mock.patch.object(
                sys,
                "argv",
                argv,
            ),
            mock.patch(
                "builtins.print",
            ),
        ):
            try:
                self.assertEqual(
                    audit.main(),
                    1,
                )
            finally:
                output.unlink(
                    missing_ok=True,
                )
                summary.unlink(
                    missing_ok=True,
                )



    def assert_legacy_shared_error_response_ratchet_enforces(
        self,
        owner: str,
    ):
        finding = {
            "category": "go.legacy_shared_error_response",
            "severity": "debt",
            "owner": owner,
            "path": f"{owner}/fixture.go",
            "line": 1,
            "excerpt": "&sharepb.ErrorResponse{Code: 401}",
        }

        slug = owner.replace("/", "-")

        output = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / f"{slug}-legacy-error-response-ratchet.tsv"
        )

        summary = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / f"{slug}-legacy-error-response-ratchet.json"
        )

        argv = [
            "audit-observability-errors.py",
            "--output",
            str(output),
            "--summary",
            str(summary),
            "--enforce-ratchets",
        ]

        with (
            mock.patch.object(
                audit,
                "scan",
                return_value=[finding],
            ),
            mock.patch.object(
                sys,
                "argv",
                argv,
            ),
            mock.patch(
                "builtins.print",
            ),
        ):
            try:
                self.assertEqual(
                    audit.main(),
                    1,
                )
            finally:
                output.unlink(
                    missing_ok=True,
                )
                summary.unlink(
                    missing_ok=True,
                )

    def test_chat_service_legacy_shared_error_response_zero_ratchet_is_registered(
        self,
    ):
        self.assertIn(
            (
                "go.legacy_shared_error_response",
                "chat-service",
            ),
            audit.ZERO_RATCHETS,
        )

    def test_chat_service_legacy_shared_error_response_ratchet_enforces_regression(
        self,
    ):
        self.assert_legacy_shared_error_response_ratchet_enforces(
            "chat-service",
        )

    def test_chat_v1_service_legacy_shared_error_response_zero_ratchet_is_registered(
        self,
    ):
        self.assertIn(
            (
                "go.legacy_shared_error_response",
                "chat-v1-service",
            ),
            audit.ZERO_RATCHETS,
        )

    def test_chat_v1_service_legacy_shared_error_response_ratchet_enforces_regression(
        self,
    ):
        self.assert_legacy_shared_error_response_ratchet_enforces(
            "chat-v1-service",
        )

    def test_relay_service_legacy_shared_error_response_zero_ratchet_is_registered(
        self,
    ):
        self.assertIn(
            (
                "go.legacy_shared_error_response",
                "relay-service",
            ),
            audit.ZERO_RATCHETS,
        )

    def test_relay_service_legacy_shared_error_response_ratchet_enforces_regression(
        self,
    ):
        self.assert_legacy_shared_error_response_ratchet_enforces(
            "relay-service",
        )


    def test_search_service_legacy_std_log_zero_ratchet_is_registered(self):
        self.assertIn(
            (
                "go.legacy_std_log",
                "search-service",
            ),
            audit.ZERO_RATCHETS,
        )

    def test_search_service_legacy_std_log_ratchet_enforces_regression(self):
        finding = {
            "category": "go.legacy_std_log",
            "severity": "debt",
            "owner": "search-service",
            "path": "search-service/main.go",
            "line": 1,
            "excerpt": 'log.Printf("search service failed: %v", err)',
        }

        output = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "search-service-legacy-std-log-ratchet.tsv"
        )
        summary = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "search-service-legacy-std-log-ratchet.json"
        )

        argv = [
            "audit-observability-errors.py",
            "--output",
            str(output),
            "--summary",
            str(summary),
            "--enforce-ratchets",
        ]

        with (
            mock.patch.object(
                audit,
                "scan",
                return_value=[finding],
            ),
            mock.patch.object(
                sys,
                "argv",
                argv,
            ),
            mock.patch("builtins.print"),
        ):
            try:
                self.assertEqual(
                    audit.main(),
                    1,
                )
            finally:
                output.unlink(missing_ok=True)
                summary.unlink(missing_ok=True)


    def test_file_service_legacy_std_log_zero_ratchet_is_registered(self):
        self.assertIn(
            (
                "go.legacy_std_log",
                "file-service",
            ),
            audit.ZERO_RATCHETS,
        )

    def test_file_service_legacy_std_log_ratchet_enforces_regression(self):
        finding = {
            "category": "go.legacy_std_log",
            "severity": "debt",
            "owner": "file-service",
            "path": "file-service/main.go",
            "line": 1,
            "excerpt": 'log.Printf("File-service stopped with error: %v", err)',
        }

        output = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "file-service-legacy-std-log-ratchet.tsv"
        )
        summary = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "file-service-legacy-std-log-ratchet.json"
        )

        argv = [
            "audit-observability-errors.py",
            "--output",
            str(output),
            "--summary",
            str(summary),
            "--enforce-ratchets",
        ]

        with (
            mock.patch.object(
                audit,
                "scan",
                return_value=[finding],
            ),
            mock.patch.object(
                sys,
                "argv",
                argv,
            ),
            mock.patch("builtins.print"),
        ):
            try:
                self.assertEqual(
                    audit.main(),
                    1,
                )
            finally:
                output.unlink(missing_ok=True)
                summary.unlink(missing_ok=True)


    def test_social_service_legacy_std_log_zero_ratchet_is_registered(self):
        self.assertIn(
            ("go.legacy_std_log", "social-service"),
            audit.ZERO_RATCHETS,
        )

    def test_social_service_legacy_std_log_ratchet_enforces_regression(self):
        finding = {
            "category": "go.legacy_std_log",
            "severity": "debt",
            "owner": "social-service",
            "path": "social-service/cmd/grpc/main.go",
            "line": 1,
            "excerpt": 'log.Printf("GRPC Listen: %v", port)',
        }

        output = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "social-service-legacy-std-log-ratchet.tsv"
        )
        summary = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "social-service-legacy-std-log-ratchet.json"
        )

        argv = [
            "audit-observability-errors.py",
            "--output",
            str(output),
            "--summary",
            str(summary),
            "--enforce-ratchets",
        ]

        with (
            mock.patch.object(
                audit,
                "scan",
                return_value=[finding],
            ),
            mock.patch.object(
                sys,
                "argv",
                argv,
            ),
            mock.patch("builtins.print"),
        ):
            try:
                self.assertEqual(
                    audit.main(),
                    1,
                )
            finally:
                output.unlink(missing_ok=True)
                summary.unlink(missing_ok=True)


    def test_assistant_service_legacy_std_log_zero_ratchet_is_registered(self):
        self.assertIn(
            ("go.legacy_std_log", "assistant-service"),
            audit.ZERO_RATCHETS,
        )

    def test_assistant_service_legacy_std_log_ratchet_enforces_regression(self):
        finding = {
            "category": "go.legacy_std_log",
            "severity": "debt",
            "owner": "assistant-service",
            "path": "assistant-service/cmd/grpc/main.go",
            "line": 1,
            "excerpt": 'log.Printf("Assistant Service is ready")',
        }

        output = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "assistant-service-legacy-std-log-ratchet.tsv"
        )
        summary = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "assistant-service-legacy-std-log-ratchet.json"
        )

        argv = [
            "audit-observability-errors.py",
            "--output",
            str(output),
            "--summary",
            str(summary),
            "--enforce-ratchets",
        ]

        with (
            mock.patch.object(
                audit,
                "scan",
                return_value=[finding],
            ),
            mock.patch.object(
                sys,
                "argv",
                argv,
            ),
            mock.patch("builtins.print"),
        ):
            try:
                self.assertEqual(
                    audit.main(),
                    1,
                )
            finally:
                output.unlink(missing_ok=True)
                summary.unlink(missing_ok=True)



    def test_notification_legacy_std_log_zero_ratchet_is_registered(self):
        self.assertIn(
            ("go.legacy_std_log", "notification-service"),
            audit.ZERO_RATCHETS,
        )

    def test_notification_legacy_std_log_ratchet_enforces_regression(self):
        finding = {
            "category": "go.legacy_std_log",
            "severity": "debt",
            "owner": "notification-service",
            "path": "notification-service/cmd/grpc_server.go",
            "line": 1,
            "excerpt": 'log.Printf("notification gRPC listening")',
        }

        output = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "notification-legacy-std-log-ratchet.tsv"
        )
        summary = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "notification-legacy-std-log-ratchet.json"
        )

        argv = [
            "audit-observability-errors.py",
            "--output",
            str(output),
            "--summary",
            str(summary),
            "--enforce-ratchets",
        ]

        with (
            mock.patch.object(
                audit,
                "scan",
                return_value=[finding],
            ),
            mock.patch.object(
                sys,
                "argv",
                argv,
            ),
            mock.patch("builtins.print"),
        ):
            try:
                self.assertEqual(
                    audit.main(),
                    1,
                )
            finally:
                output.unlink(missing_ok=True)
                summary.unlink(missing_ok=True)



    def test_chat_logging_zero_ratchets_are_registered(self):
        self.assertIn(
            ("go.legacy_std_log", "chat-service"),
            audit.ZERO_RATCHETS,
        )
        self.assertIn(
            ("go.third_party_logger", "chat-service"),
            audit.ZERO_RATCHETS,
        )

    def test_chat_legacy_std_log_ratchet_enforces_regression(self):
        finding = {
            "category": "go.legacy_std_log",
            "severity": "debt",
            "owner": "chat-service",
            "path": "chat-service/main.go",
            "line": 1,
            "excerpt": 'log.Printf("chat startup failed: %v", err)',
        }

        output = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "chat-legacy-std-log-ratchet.tsv"
        )
        summary = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "chat-legacy-std-log-ratchet.json"
        )

        argv = [
            "audit-observability-errors.py",
            "--output",
            str(output),
            "--summary",
            str(summary),
            "--enforce-ratchets",
        ]

        with (
            mock.patch.object(
                audit,
                "scan",
                return_value=[finding],
            ),
            mock.patch.object(
                sys,
                "argv",
                argv,
            ),
            mock.patch("builtins.print"),
        ):
            try:
                self.assertEqual(
                    audit.main(),
                    1,
                )
            finally:
                output.unlink(missing_ok=True)
                summary.unlink(missing_ok=True)

    def test_chat_third_party_logger_ratchet_enforces_regression(self):
        finding = {
            "category": "go.third_party_logger",
            "severity": "debt",
            "owner": "chat-service",
            "path": "chat-service/infra/server/http_server.go",
            "line": 1,
            "excerpt": (
                "github.com/hyperledger/"
                "fabric/common/flogging"
            ),
        }

        output = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "chat-third-party-logger-ratchet.tsv"
        )
        summary = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "chat-third-party-logger-ratchet.json"
        )

        argv = [
            "audit-observability-errors.py",
            "--output",
            str(output),
            "--summary",
            str(summary),
            "--enforce-ratchets",
        ]

        with (
            mock.patch.object(
                audit,
                "scan",
                return_value=[finding],
            ),
            mock.patch.object(
                sys,
                "argv",
                argv,
            ),
            mock.patch("builtins.print"),
        ):
            try:
                self.assertEqual(
                    audit.main(),
                    1,
                )
            finally:
                output.unlink(missing_ok=True)
                summary.unlink(missing_ok=True)


    def test_chat_v1_logging_zero_ratchets_are_registered(self):
        self.assertIn(
            ("go.legacy_std_log", "chat-v1-service"),
            audit.ZERO_RATCHETS,
        )
        self.assertIn(
            ("go.third_party_logger", "chat-v1-service"),
            audit.ZERO_RATCHETS,
        )

    def test_chat_v1_legacy_std_log_ratchet_enforces_regression(self):
        finding = {
            "category": "go.legacy_std_log",
            "severity": "debt",
            "owner": "chat-v1-service",
            "path": "chat-v1-service/main.go",
            "line": 1,
            "excerpt": 'log.Printf("chat-v1 startup failed: %v", err)',
        }

        output = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "chat-v1-legacy-std-log-ratchet.tsv"
        )

        summary = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "chat-v1-legacy-std-log-ratchet.json"
        )

        argv = [
            "audit-observability-errors.py",
            "--output",
            str(output),
            "--summary",
            str(summary),
            "--enforce-ratchets",
        ]

        with (
            mock.patch.object(
                audit,
                "scan",
                return_value=[finding],
            ),
            mock.patch.object(
                sys,
                "argv",
                argv,
            ),
            mock.patch(
                "builtins.print",
            ),
        ):
            try:
                self.assertEqual(
                    audit.main(),
                    1,
                )
            finally:
                output.unlink(
                    missing_ok=True,
                )
                summary.unlink(
                    missing_ok=True,
                )

    def test_chat_v1_third_party_logger_ratchet_enforces_regression(self):
        finding = {
            "category": "go.third_party_logger",
            "severity": "debt",
            "owner": "chat-v1-service",
            "path": "chat-v1-service/infrastructure/server/http_server.go",
            "line": 1,
            "excerpt": (
                "github.com/hyperledger/"
                "fabric/common/flogging"
            ),
        }

        output = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "chat-v1-third-party-logger-ratchet.tsv"
        )

        summary = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "chat-v1-third-party-logger-ratchet.json"
        )

        argv = [
            "audit-observability-errors.py",
            "--output",
            str(output),
            "--summary",
            str(summary),
            "--enforce-ratchets",
        ]

        with (
            mock.patch.object(
                audit,
                "scan",
                return_value=[finding],
            ),
            mock.patch.object(
                sys,
                "argv",
                argv,
            ),
            mock.patch(
                "builtins.print",
            ),
        ):
            try:
                self.assertEqual(
                    audit.main(),
                    1,
                )
            finally:
                output.unlink(
                    missing_ok=True,
                )
                summary.unlink(
                    missing_ok=True,
                )


    def test_relay_logging_zero_ratchets_are_registered(self):
        self.assertIn(
            ("go.legacy_std_log", "relay-service"),
            audit.ZERO_RATCHETS,
        )
        self.assertIn(
            ("go.third_party_logger", "relay-service"),
            audit.ZERO_RATCHETS,
        )

    def test_relay_legacy_std_log_ratchet_enforces_regression(self):
        finding = {
            "category": "go.legacy_std_log",
            "severity": "debt",
            "owner": "relay-service",
            "path": "relay-service/main.go",
            "line": 1,
            "excerpt": 'log.Printf("relay startup failed: %v", err)',
        }

        output = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "relay-legacy-std-log-ratchet.tsv"
        )
        summary = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "relay-legacy-std-log-ratchet.json"
        )

        argv = [
            "audit-observability-errors.py",
            "--output",
            str(output),
            "--summary",
            str(summary),
            "--enforce-ratchets",
        ]

        with (
            mock.patch.object(
                audit,
                "scan",
                return_value=[finding],
            ),
            mock.patch.object(
                sys,
                "argv",
                argv,
            ),
            mock.patch("builtins.print"),
        ):
            try:
                self.assertEqual(
                    audit.main(),
                    1,
                )
            finally:
                output.unlink(missing_ok=True)
                summary.unlink(missing_ok=True)

    def test_relay_third_party_logger_ratchet_enforces_regression(self):
        finding = {
            "category": "go.third_party_logger",
            "severity": "debt",
            "owner": "relay-service",
            "path": "relay-service/server/ws_server.go",
            "line": 1,
            "excerpt": (
                "github.com/hyperledger/"
                "fabric/common/flogging"
            ),
        }

        output = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "relay-third-party-logger-ratchet.tsv"
        )
        summary = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "relay-third-party-logger-ratchet.json"
        )

        argv = [
            "audit-observability-errors.py",
            "--output",
            str(output),
            "--summary",
            str(summary),
            "--enforce-ratchets",
        ]

        with (
            mock.patch.object(
                audit,
                "scan",
                return_value=[finding],
            ),
            mock.patch.object(
                sys,
                "argv",
                argv,
            ),
            mock.patch("builtins.print"),
        ):
            try:
                self.assertEqual(
                    audit.main(),
                    1,
                )
            finally:
                output.unlink(missing_ok=True)
                summary.unlink(missing_ok=True)



    def test_payment_logging_zero_ratchets_are_registered(self):
        self.assertIn(
            ("go.legacy_std_log", "payment-service"),
            audit.ZERO_RATCHETS,
        )
        self.assertIn(
            ("go.third_party_logger", "payment-service"),
            audit.ZERO_RATCHETS,
        )

    def test_payment_legacy_std_log_ratchet_enforces_regression(self):
        finding = {
            "category": "go.legacy_std_log",
            "severity": "debt",
            "owner": "payment-service",
            "path": "payment-service/worker/outbox.go",
            "line": 1,
            "excerpt": 'log.Printf("payment failed: %v", err)',
        }
        output = audit.ROOT / ".tmp" / "observability-errors" / "payment-std-log-ratchet.tsv"
        summary = audit.ROOT / ".tmp" / "observability-errors" / "payment-std-log-ratchet.json"
        argv = [
            "audit-observability-errors.py",
            "--output",
            str(output),
            "--summary",
            str(summary),
            "--enforce-ratchets",
        ]
        with (
            mock.patch.object(audit, "scan", return_value=[finding]),
            mock.patch.object(sys, "argv", argv),
            mock.patch("builtins.print"),
        ):
            try:
                self.assertEqual(audit.main(), 1)
            finally:
                output.unlink(missing_ok=True)
                summary.unlink(missing_ok=True)

    def test_payment_third_party_logger_ratchet_enforces_regression(self):
        finding = {
            "category": "go.third_party_logger",
            "severity": "debt",
            "owner": "payment-service",
            "path": "payment-service/cmd/grpc/runtime.go",
            "line": 1,
            "excerpt": 'github.com/hyperledger/fabric/common/flogging',
        }
        output = audit.ROOT / ".tmp" / "observability-errors" / "payment-third-party-ratchet.tsv"
        summary = audit.ROOT / ".tmp" / "observability-errors" / "payment-third-party-ratchet.json"
        argv = [
            "audit-observability-errors.py",
            "--output",
            str(output),
            "--summary",
            str(summary),
            "--enforce-ratchets",
        ]
        with (
            mock.patch.object(audit, "scan", return_value=[finding]),
            mock.patch.object(sys, "argv", argv),
            mock.patch("builtins.print"),
        ):
            try:
                self.assertEqual(audit.main(), 1)
            finally:
                output.unlink(missing_ok=True)
                summary.unlink(missing_ok=True)


    def test_canonical_http_error_serializer_owner_is_shared(self):
        summary = audit.build_summary([])

        self.assertEqual(
            summary["policy"]["canonical_http_error_serializer"],
            "shared/common/httpresponse.WriteProblem",
        )

        self.assertEqual(
            summary["policy"]["gateway_http_error_facade"],
            "gateway-service/internal/httpresponse.WriteProblem",
        )



    def test_user_service_third_party_logger_zero_ratchet_is_registered(self):
        self.assertIn(
            (
                "go.third_party_logger",
                "user-service",
            ),
            audit.ZERO_RATCHETS,
        )

        self.assertNotIn(
            (
                "go.legacy_std_log",
                "user-service",
            ),
            audit.ZERO_RATCHETS,
        )

    def test_user_service_third_party_logger_ratchet_enforces_regression(self):
        finding = {
            "category": "go.third_party_logger",
            "severity": "debt",
            "owner": "user-service",
            "path": "user-service/initial/startup.go",
            "line": 1,
            "excerpt": (
                "github.com/hyperledger/"
                "fabric/common/flogging"
            ),
        }

        output = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "user-third-party-logger-ratchet.tsv"
        )

        summary = (
            audit.ROOT
            / ".tmp"
            / "observability-errors"
            / "user-third-party-logger-ratchet.json"
        )

        argv = [
            "audit-observability-errors.py",
            "--output",
            str(output),
            "--summary",
            str(summary),
            "--enforce-ratchets",
        ]

        with (
            mock.patch.object(
                audit,
                "scan",
                return_value=[finding],
            ),
            mock.patch.object(
                sys,
                "argv",
                argv,
            ),
            mock.patch(
                "builtins.print",
            ),
        ):
            try:
                self.assertEqual(
                    audit.main(),
                    1,
                )
            finally:
                output.unlink(
                    missing_ok=True,
                )
                summary.unlink(
                    missing_ok=True,
                )



if __name__ == "__main__":
    unittest.main()
