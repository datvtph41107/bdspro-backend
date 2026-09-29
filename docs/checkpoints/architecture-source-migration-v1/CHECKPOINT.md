# BDSPro Architecture Source Migration V1 — Durable Checkpoint

Status: CLOSED AT V1 / PAUSED FOR OTHER WORK
Closed: 2026-09-29 (Asia/Bangkok)

## Canonical refs

Repository:
- datvtph41107/bdspro-backend

Live main at checkpoint close:
- branch: main
- commit: b2a7a8167e9c8af2a47f1fde4f7b2a9d951611fd
- tree: f4be525468190ea1d43b6123d72db30495754504

Architecture migration branch:
- branch: architecture/source-migration-v1
- commit: 5fa292903f74396a07b9a4eeb21540a76ddbe83c
- tree: 771f0d176c01e0b36b77780e316710441ac737e8
- parent: b2a7a8167e9c8af2a47f1fde4f7b2a9d951611fd

Local proof commit created during migration:
- commit: f2968b98340fcf276772a29d3dd2831984a1cd44
- tree: d607540161350fe977b9f7544e8f15efcdb807f1

Important:
- The remote migration tree is derived from the local proof commit.
- One intentional documentation difference exists: internal/user/README.md was sanitized during publication so the remote checkpoint does not embed development credential examples.
- No runtime/business code was changed for that sanitization.

Recovery Sprite:
- mcp-123-bdspro-architecture-migration-v1
- pre-migration checkpoint: v1
- committed V1 proof checkpoint: v2

## What V1 changed

V1 is a physical source-topology migration. It does NOT claim a finished business-domain decomposition.

Representative mapping:

- gateway-service -> api/gateway
- user-service -> internal/user
- auth-service -> internal/authentication
- organization-service -> internal/organization-legacy
- bdspro-service -> internal/property
- crm-service -> internal/crm
- tqd-service -> internal/planning
- payment-service -> internal/payment
- file-service -> internal/file
- notification-service -> internal/notification
- hub-service -> internal/hub
- assistant-service -> internal/assistant
- chat-service -> internal/chat
- chat-v1-service -> internal/chat-legacy
- social-service -> internal/social
- map-service -> internal/map-legacy
- relay-service -> internal/realtime-relay
- search-service -> internal/search
- ai-service -> internal/ai

Shared/repository mapping:

- shared/common -> infrastructure/runtime
- shared/base -> infrastructure/base
- shared/config -> infrastructure/config
- shared/protobuf -> proto
- shared/code -> tools/development
- shared/scripts -> tools/scripts
- shared/summary -> docs/shared-summary
- nginx -> infrastructure/nginx
- documents -> docs/reference

## Operational rule introduced by V1

Physical source path and runtime service identity are separate concepts.

Example:
- runtime identity may remain auth-service
- source path may be internal/authentication

Repository tooling resolves source paths explicitly instead of assuming:
- <service>-service

This prevents source organization from implicitly renaming network/runtime identity.

## What V1 deliberately did NOT do

V1 did NOT:

- merge the migration branch into main;
- redesign canonical business facts;
- finalize bounded contexts;
- split user into account/profile/authentication/authorization;
- redesign Organization/Membership facts;
- redesign databases;
- change canonical database ownership;
- create a new cross-service transaction model;
- rename runtime service identities or wire protocol service names;
- declare internal/user, internal/property, internal/planning, or infrastructure/runtime to be final architecture;
- start V2.

## Proof completed before checkpoint

Generation:
- make generate-backend: PASS
- protobuf generation: PASS
- Wire generation for the migrated paths: PASS

Core build:
- make build-backend: PASS across the core runtime chain.

Repository-module proof:
- make verify-repository-modules: PASS
- gate marker: __REPOSITORY_GATE_EXIT__=0
- included vet, test, race test, and build across the repository modules covered by the repository gate.

Additional gates:
- make verify-non-go-source: PASS
- make verify-docs: PASS
- tools/development/verify-source-layout.sh: PASS
- source/Docker topology checks: PASS
- staged git diff --check before V1 commit: PASS

## Known debt preserved intentionally

These are NOT silently considered solved:

1. internal/user is still a large current authority host and is not a final business boundary.
2. internal/planning still contains Planning/Report/Quota/Usage/Map-related concerns pending fact-driven decomposition.
3. internal/property remains a large current source area and requires business-capability reconstruction before deeper movement.
4. infrastructure/runtime is the migrated form of shared/common and still contains semantic debt; its name does not prove every package inside is infrastructure-safe.
5. internal/organization-legacy remains a legacy/migration source; V1 does not prove retirement.
6. A tracked Python __pycache__ / .pyc artifact exists under tools/development and is preserved in this checkpoint; clean it in a dedicated hygiene change, not by rewriting V1 history.
7. Database redesign is not started.
8. Architecture Learning business-fact theorems are not superseded by this source-topology branch.

## Worktree warning

After local commit f2968b9, later documentation verification activity caused additional working-tree documentation changes inside the recovery Sprite.

Those post-checkpoint working-tree changes are NOT part of this checkpoint.

For a future resume:
- prefer a fresh clone/worktree from the checkpoint branch or migration branch;
- or restore Sprite checkpoint v2 before inspecting/mutating local source;
- do not treat the later dirty Sprite worktree as canonical.

## Resume order

When this migration is resumed:

1. Read this CHECKPOINT.md.
2. Read CONTINUATION-PROMPT.md.
3. Fetch and verify live main.
4. Verify the checkpoint branch exact commit.
5. Verify architecture/source-migration-v1 exact commit.
6. Compare live source against this checkpoint before any source mutation.
7. Reconstruct the business fact/owner/invariant relevant to the chosen V2 slice.
8. Only then choose the minimum source/database/runtime change.

## Next work when explicitly resumed

V2 is NOT STARTED.

Candidate future work must be selected from evidence, not from folder aesthetics. Possible investigation areas include:
- internal/user semantic decomposition;
- internal/planning fact and lifecycle decomposition;
- internal/property capability decomposition;
- infrastructure/runtime semantic-debt cleanup;
- removal of tracked generated/cache artifacts;
- module/public-boundary enforcement;
- database ownership and schema redesign after business facts are proven.

No candidate above is pre-authorized by this checkpoint.

## Closure statement

Architecture Source Migration V1 is closed as an independently recoverable branch milestone.

main remains untouched.

Future work may proceed on other BDSPro priorities without relying on conversational memory to recover this migration.
