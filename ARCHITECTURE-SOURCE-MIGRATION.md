# BDSPro Source Architecture Migration V1

Branch: `architecture/source-migration-v1`
Baseline: `main@b2a7a8167e9c8af2a47f1fde4f7b2a9d951611fd`

## Purpose

This branch changes source topology without claiming new business facts.
The migration deliberately separates:

- API/edge runtime source;
- current business/runtime implementation areas;
- concrete infrastructure mechanics;
- protobuf contracts;
- engineering tools;
- integration proof and documentation.

It does **not** declare the current inner structure of large areas such as
`internal/user`, `internal/property`, or `internal/planning` to be the
final business-module decomposition. Those areas remain convergence targets
for fact-by-fact architecture work.

## Mapping

| Previous path | V1 path | Meaning |
| --- | --- | --- |
| gateway-service | api/gateway | external API/edge boundary |
| user-service | internal/user | current User authority host; decomposition pending |
| auth-service | internal/authentication | auth compatibility/read boundary |
| organization-service | internal/organization-legacy | legacy/migration organization source |
| bdspro-service | internal/property | property/BDSPro implementation area |
| crm-service | internal/crm | CRM implementation area |
| tqd-service | internal/planning | planning/report/quota/usage runtime area |
| payment-service | internal/payment | payment runtime/business owner |
| file-service | internal/file | file/artifact runtime |
| notification-service | internal/notification | notification/delivery runtime |
| hub-service | internal/hub | hub/config compatibility area |
| assistant-service | internal/assistant | assistant runtime |
| chat-service | internal/chat | current chat runtime |
| chat-v1-service | internal/chat-legacy | legacy chat runtime |
| social-service | internal/social | social runtime |
| map-service | internal/map-legacy | legacy map runtime |
| relay-service | internal/realtime-relay | realtime relay runtime |
| search-service | internal/search | search runtime/projection |
| ai-service | internal/ai | AI runtime |
| shared/base | infrastructure/base | shared technical base |
| shared/common | infrastructure/runtime | shared runtime mechanics; semantic debt still audited |
| shared/protobuf | proto | protobuf/wire contracts |
| shared/code | tools/development | engineering/build/development tooling |
| shared/config | infrastructure/config | shared runtime config |
| shared/scripts | tools/scripts | repository scripts |
| shared/summary | docs/shared-summary | historical summaries |
| nginx | infrastructure/nginx | edge infrastructure |
| documents | docs/reference | repository documentation |

## V1 rules

1. No business fact is renamed only to make the tree prettier.
2. Existing Go module import paths remain unchanged in V1.
3. Physical source paths are updated everywhere they are execution/build inputs.
4. `api/` means API/edge application surface, not business ownership.
5. `internal/` paths are current semantic/runtime areas, not proof of final bounded contexts.
6. `infrastructure/` contains concrete technical/runtime mechanics, not business owners.
7. `proto/` is wire-contract source, not durable business authority.
8. V2 may split large current areas only after fact/owner/invariant/lifecycle proof.
