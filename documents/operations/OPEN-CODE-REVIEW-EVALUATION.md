# BDSPro Evaluation — Alibaba OpenCodeReview

Updated: 2026-09-18 Asia/Ho_Chi_Minh
Status: PILOT APPROVED / NOT YET INSTALLED / NOT AN ACCEPTANCE AUTHORITY

## Why evaluate it

BDSPro needs long-lived review capacity that scales without forcing the ChatGPT coordinator to manually inspect every changed line sequentially.

OpenCodeReview provides:
- deterministic Git diff/file selection plus LLM-agent review;
- exact line-level structured comments;
- commit/range/workspace review;
- preview mode without LLM calls;
- JSON/SARIF machine-readable output;
- resumable sessions;
- configurable review rules;
- concurrency;
- coding-agent and CI integrations.

This maps naturally to the BDSPro model of many read-only/proof workers around one source writer.

## Role in BDSPro

OpenCodeReview may become:

`candidate exact SHA -> OCR review lane -> coordinator triage -> deterministic acceptance gates`

It must never become:

`OCR says PASS -> source accepted`

OCR findings are advisory review evidence.

## Pilot placement

After a bounded writer slice produces a committed candidate:

1. deterministic local tests/audit/protected checks remain mandatory;
2. create/use detached proof worktree at candidate SHA;
3. OCR preview exact base SHA -> candidate SHA;
4. OCR review exact base SHA -> candidate SHA;
5. save JSON/SARIF/text output under `~/bdspro-ops/proofs/ocr/<candidate>/`;
6. coordinator classifies each finding: confirmed / false-positive / out-of-scope / follow-up;
7. any source fix goes back through the ONE WRITER lane;
8. candidate is re-committed/re-proved if source changes.

## Reproducibility/safety rules

- pin an exact OCR version for local/CI usage;
- record `ocr version` in evidence;
- use only coordinator-validated full 40-hex Git SHAs in ref arguments;
- never forward untrusted branch/ref text directly into OCR;
- keep LLM credentials outside Git/chat;
- do not enable auto-fix during pilot;
- do not use OCR output as a ratchet source by itself;
- do not use floating `@main` GitHub Action references;
- do not start with `pull_request_target` + repository secrets/write token;
- hosted OCR CI requires a separate security review and least-privilege design.

## BDSPro-specific rules opportunity

A future project rule file may teach OCR to look for:
- protected no-touch paths;
- canonical logger ownership/component/correlation rules;
- Error/Response FINAL owner boundaries;
- accidental direct HTTP/gRPC error projection in domain/usecase layers;
- transaction/idempotency/worker lifecycle mistakes;
- context cancellation/goroutine/resource lifecycle;
- generated protobuf/wire files that should not be hand-edited.

These are review prompts only. Existing deterministic checks remain authoritative.

## Adoption phases

A. local install/version/preview only;
B. one known historical candidate review for signal/noise calibration;
C. one live R5 candidate review;
D. custom BDSPro rule file if useful;
E. optional read-only hosted CI artifact lane;
F. only later consider PR comment publishing.

No phase may bypass exact-SHA proof or one-writer discipline.
