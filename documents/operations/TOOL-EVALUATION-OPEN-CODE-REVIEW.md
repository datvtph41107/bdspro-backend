# BDSPro Tool Evaluation — Alibaba OpenCodeReview

Updated: 2026-09-18 Asia/Ho_Chi_Minh
Status: EVALUATED / PILOT AUTHORIZED / NOT AN ACCEPTANCE AUTHORITY

Tool:
- repository: `alibaba/open-code-review`
- CLI: `ocr`
- license: Apache-2.0
- current public project state reviewed on 2026-09-18

## Why it is relevant to BDSPro

OpenCodeReview provides:
- deterministic Git diff/file-selection pipeline;
- configurable per-path review rules;
- LLM-assisted line-level review;
- branch/range/commit/workspace review modes;
- machine-readable JSON and SARIF output;
- local CLI, agent plugins and CI integrations;
- support for OpenAI-compatible, Anthropic and other providers;
- preview mode that can inspect review scope without running the LLM.

This complements, but does not replace, BDSPro's acceptance model.

## Correct role in BDSPro

OpenCodeReview is an OPTIONAL REVIEW / PROOF PROJECTION.

It is NOT:
- source authority;
- architecture authority;
- replacement for tests;
- replacement for BDSPro audit/ratchets;
- replacement for protected-path/deploy invariants;
- replacement for exact-SHA hosted CI;
- permission to auto-fix source;
- a gate whose process exit code alone can declare acceptance.

Initial placement:

candidate exact SHA
  -> deterministic BDSPro proof lanes
  -> OpenCodeReview review lane
  -> coordinator reconciles findings
  -> fixes, if justified, go back through ONE WRITER
  -> new exact SHA
  -> all proof lanes rerun

## Why not make OCR a hard gate immediately

1. LLM review is probabilistic.
2. Findings can be false positives or miss project-specific invariants.
3. OpenCodeReview has active development and recent behavior/security issues.
4. Current upstream issue #1196 reports cases where code-comment failure can still yield exit 0 / complete coverage semantics, so process success alone is not sufficient admission evidence.
5. GitHub Action integration uses privileged GitHub permissions and documented pull_request_target patterns; BDSPro must threat-model and pin versions before enabling automated PR commenting.
6. BDSPro already has stronger deterministic acceptance gates for several concerns.

## Security / trust decision

Phase A pilot:
- LOCAL CLI only;
- READ-ONLY exact-SHA proof worktree;
- no GitHub token;
- no PR write permission;
- no auto-fix;
- no source mutation;
- no secrets committed to repository;
- run `ocr review --preview` first;
- pin an explicit OCR version after pilot validation.

Phase B, only after pilot:
- exact `base -> candidate` range review;
- JSON artifact stored under `~/bdspro-ops/proofs`;
- custom BDSPro rules/background context;
- coordinator classifies findings;
- OCR remains advisory unless a deterministic BDSPro policy promotes a specific machine-verifiable rule.

Phase C, optional future:
- GitHub Actions PR review;
- version-pinned action/package;
- minimal permissions;
- no unsafe fork-secret exposure;
- no acceptance gate based only on OCR exit code;
- artifacts retained for audit.

## BDSPro custom review context

OCR should be given BDSPro-specific background/rules such as:
- ONE TRUTH PER CONCERN -> ONE OWNER -> projections do not redefine truth;
- preserve protected paths;
- no competing outbound Error truth;
- logging/metrics are projections;
- one source writer / exact-SHA evidence discipline;
- no ratchet before zero proof;
- do not silently alter lifecycle/signal/exit semantics;
- no generated/protobuf/protected mutations unless explicitly authorized.

These rules improve review relevance but do not supersede durable architecture authority.

## Pilot plan

1. Complete local Phase 3C read-only inventory loop first.
2. In a separate USER-LANE, inspect local Node/npm and install OCR outside the repository.
3. Record installed OCR version and package checksum/source.
4. Run `ocr review --preview` on a known exact Git range without LLM/source mutation.
5. Verify selected files match the bounded candidate scope.
6. Configure provider only after the user explicitly chooses one and secrets handling is defined.
7. Run one advisory JSON review on a previously closed slice (Payment R5) as calibration.
8. Compare OCR findings with known source/proof reality.
9. Only then decide whether OCR becomes a standard proof lane for future candidate SHAs.

## Acceptance relationship

A future candidate is never accepted merely because OCR says no issues.

Deterministic proof and hosted exact-SHA CI remain mandatory.

OCR can:
- detect extra issues;
- request coordinator review;
- block progression if a credible finding is unresolved;
- never independently declare the slice CLOSED.

FINAL ACCEPTED = NO.
