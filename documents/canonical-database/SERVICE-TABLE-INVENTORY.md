# Service Table Inventory and Final Disposition

This is the full database inventory used for final architecture acceptance.

Legend:

- **CORE→TARGET** — semantic source exists but target representation is replaced by v1/v2 canonical core.
- **RETAIN** — bounded-context authority remains valid and separate.
- **OPERATIONAL** — durable operational/projection/audit state; not promoted to cross-domain business truth.
- **COMPAT→RETIRE** — compatibility/duplicated truth that remains during migration only.
- **RECONCILE** — legacy data must be profiled and classified before contraction.

---

## user-service

Canonical migration authority: `user-service/database/migrations/`.

### Identity / authentication

| Legacy table | Disposition | Final meaning |
|---|---|---|
| `auth_method` | CORE→TARGET | split into Account + password/federated/passkey credential facts |
| `auth_devices` | OPERATIONAL | device/security context; no Person/Organization authority |
| `user_otp` | OPERATIONAL | temporary auth challenge |
| `user_session` | CORE→TARGET | `account_sessions` |
| `user_status` | COMPAT→RETIRE | account/block/verification facts must not remain one status blob |
| `user_info` | COMPAT→RETIRE | profile/lock/role conflation |
| `user_pin` | OPERATIONAL | credential/security mechanism |
| `auth_config` | OPERATIONAL | security configuration |
| `admin_profiles` | RECONCILE | admin presentation + role/security duplication; split by fact |
| `user_profile` | CORE→TARGET | Person presentation migrates to `profiles`; authority fields removed |
| `profile_transfer` | OPERATIONAL | migration/transfer snapshot |
| `profile_deleted` | OPERATIONAL | historical/retention snapshot; not live Profile |
| `profile_media` | RETAIN | presentation/media reference |
| `kyc` | RETAIN | KYC workflow, not Person identity |

### IAM

| Legacy table | Disposition | Final meaning |
|---|---|---|
| `permissions` | CORE→TARGET | stable `permissions.code` authority |
| `roles` | CORE→TARGET | Organization-scoped roles only |
| `role_profiles` | CORE→TARGET | replaced by OrganizationMembership-scoped RoleAssignment |
| `role_permissions` | CORE→TARGET | `organization_role_permissions` |
| `role_groups` | OPERATIONAL/REVIEW | UI/catalog grouping of permissions; not business membership |
| `group_permissions` | OPERATIONAL/REVIEW | role-group permission grouping, not Group collaboration authority |
| `colors` | RETAIN | presentation taxonomy |
| `admin_access_control` | RETAIN | admin/security policy state |
| `admin_access_logs` | OPERATIONAL | audit evidence |

### Organization directory

| Legacy table | Disposition | Final meaning |
|---|---|---|
| `organizations` | CORE→TARGET | canonical Organization |
| `organization_members` | CORE→TARGET | split Invitation/Membership/Block/RoleAssignment/Ownership |
| legacy `owner_profile_id` column | COMPAT→RETIRE | replaced by `organization_ownerships` |
| legacy `code/type/status/verification/warning` columns | RECONCILE | not copied as canonical core |

### Commercial catalog/subscription

Retain as user-service commercial authority, with subject identity converging on Party rather than profile/organization type pairs.

- `catalog_products` — RETAIN
- `catalog_plans` — RETAIN
- `catalog_plan_versions` — RETAIN
- `catalog_plan_entitlements` — RETAIN
- `catalog_plan_operation_policies` — RETAIN
- `catalog_price_items` — RETAIN
- `catalog_subscriptions` — RETAIN, target logical subject = Party
- `catalog_subscription_events` — RETAIN
- `subscription_checkout_snapshots` — RETAIN snapshot/idempotency evidence
- `subscription_settlement_receipts` — RETAIN settlement evidence
- `pack_price_table` — COMPAT→RETIRE after catalog-v2 consumers switch

### Social/profile-support tables

These remain user/profile bounded-context data; they do not become Party/IAM authority.

- `block` — RETAIN relation
- `bookmark_user` — RETAIN relation
- `certification` — RETAIN workflow
- `tb_contact` — RECONCILE with CRM ownership; avoid two mutable contact truths
- `follow` — RETAIN relation
- `friend` — RETAIN relation
- `db_group` — RECONCILE; do not confuse with Organization-service Group
- `main_area` — RETAIN taxonomy
- `main_area_profile` — RETAIN relation
- `profession` — RETAIN taxonomy
- `purpose_use` — RETAIN taxonomy
- `purpose_use_profile` — RETAIN relation
- `tags` — RETAIN taxonomy
- `tag_user` — RETAIN relation

---

## organization-service

Current migration authority: `organization-service/migrate/`.

This service contains the largest semantic debt and is the main migration target of v0→v2.

### Duplicated Organization/IAM state

These become compatibility-only once user-service canonical Organization/IAM authority is switched:

- `colors` — COMPAT→RETIRE or read projection
- `business_domains` — RETAIN taxonomy if still required
- `organizations` — COMPAT→RETIRE
- `organization_business_domains` — RECONCILE; attach to canonical Organization if business use remains
- `organization_permissions` — COMPAT→RETIRE
- `organization_roles` — COMPAT→RETIRE
- `role_permissions` — COMPAT→RETIRE
- `organization_members` — COMPAT→RETIRE
- `organization_branches` — CORE→TARGET; target authority moves with Organization structure
- `organization_branch_members` — CORE→TARGET; replaced by `organization_branch_assignments`
- `organization_log_activities` — OPERATIONAL/AUDIT

### Group

- `groups` — CORE→TARGET
- `group_members` — CORE→TARGET; replaced by GroupMembership + GroupLeadership
- `group_settings` — RECONCILE; generic key/value is not canonical business model
- `group_log_activities` — OPERATIONAL/AUDIT
- `group_notifications` — COMPAT→RETIRE in favor of notification-service
- `group_documents` — CORE→TARGET explicit Group↔File relation
- `group_chats` — CORE→TARGET integration link

### Deal

- `deals` — CORE→TARGET
- `deal_products` — CORE→TARGET
- `deal_members` — CORE→TARGET; split Invitation/Participation/relationships/authority/financial terms
- `deal_milestones` — CORE→TARGET
- `investments` — CORE→TARGET
- `attach_documents` — COMPAT→RETIRE; generic owner/doc_owner rejected
- `internal_notes` — CORE→TARGET explicit Deal note
- `deal_of_organization` — COMPAT→RETIRE duplicate Deal context
- `deal_of_branch` — CORE→TARGET as direct optional Branch placement on Deal
- `deal_of_group` — COMPAT→RETIRE duplicate Deal context
- `cost_types` — CORE→TARGET
- `deal_costs` — CORE→TARGET with amount+currency
- `deal_action` — OPERATIONAL/RECONCILE transaction/action projection; not primary Deal truth

### Legacy bank account

- `bank_accounts` — CORE→TARGET but authority moves to payment-service; generic owner type/id rejected

---

## payment-service

Canonical migration authority: `payment-service/database/migrations/`.

The commerce recovery model is retained because it already separates durable concerns correctly.

### Commerce

- `commerce_orders` — RETAIN; target subject identity = Party
- `commerce_payment_attempts` — RETAIN
- `commerce_settlements` — RETAIN
- `commerce_fulfillments` — RETAIN
- `commerce_command_effects` — RETAIN idempotent effect evidence
- `commerce_outbox_events` — RETAIN integration outbox

### Money / wallet

- `banks` — RETAIN; target lifecycle removes soft-delete duplication
- `wallets` — RETAIN; target identity uses Person
- `wallet_transactions` — RETAIN
- `wallet_audit_logs` — OPERATIONAL/AUDIT
- `withdrawal_requests` — RETAIN workflow
- `payment_methods` — RETAIN
- `transaction_types` — RETAIN reference data
- `dashboard_metrics` — OPERATIONAL projection
- `payment_stats` — OPERATIONAL projection

### New canonical payment fact

- `bank_accounts` — TARGET addition; owned by Party and Bank, replaces Organization-service generic bank account ownership.

---

## bdspro-service

Canonical migration authority: `bdspro-service/database/migrations/`.

### Property authority

- `property` — RETAIN canonical Property fact
- `property_system` — RETAIN physical partition
- `property_user` — RETAIN physical partition
- `property_system_default` — RETAIN physical/default partition
- `property_user_default` — RETAIN physical/default partition
- month partitions such as `property_system_202602`, `property_user_202602`, `property_system_202603`, `property_user_202603` — RETAIN physical storage partitions

Partition/source classification is storage/source provenance; it is not Party identity and must not leak into generic business ownership.

Deal product references are logical service-contract references to BDSPro Product/Property authority.

---

## tqd-service

Canonical migration authority: `tqd-service/database/migrations/`.

TQD remains a separate planning/quota/report bounded context. All below are retained unless explicitly noted.

### User/report/notification

- `user_subscriptions`
- `user_notifications`
- `reports`
- `report_events`
- `user_reported`
- `report_jobs`

### Planning layers/regions/labels

- `qh_layers`
- `qh_regions`
- `qh_labels`
- `qh_label_layers`
- `qh_layer_families`
- `qh_legends`
- `qh_land_use`
- `qh_layer_land_use`
- `qh_authority_issuring` — RETAIN but naming typo should be corrected only through safe migration
- `qh_planning_projects`
- `qh_planning_documents`
- `qh_planning_events`
- `qh_planning_relations`
- `qh_legal_documents`
- `qh_layer_legal_docs`
- `qh_zone_legal_docs`
- `qh_parcel_legal_docs`
- `qh_layer_legals`
- `qh_layer_resolver_configs`
- `qh_region_extends`
- `qh_region_import_error_logs`
- `qh_shape`
- `qh_direction`
- `qh_jurisdictions`
- `qh_parcel_direction`

### Parcel/POI/discovery

- `parcels`
- `qh_planning_land_use`
- `qh_parcel_info`
- `poi_categories`
- `pois`
- `poi_media`
- `poi_amenities`
- `amenities`
- `map_points`
- `search_index`
- `province_v2`
- `ward_v2`

### Directory/supporting data

- `contact_labels`
- `directory_categories`
- `directory_sources`
- `directory_suppliers`
- `open_hours`

### Audit/lifecycle

- `qh_audit_entries` — OPERATIONAL/AUDIT
- `qh_lifecycle_transitions` — RETAIN domain history

### AI/jobs

- `pro_ai_jobs` — OPERATIONAL process state

### User workflows

- `user_followed_planning_projects`
- `user_followed_parcels`
- `user_view_events` — OPERATIONAL/event
- `user_view_history` — projection/history

### Quota

- `quota_usage_events` — RETAIN durable usage fact
- `quota_pools` — RETAIN current projection/reservation pool
- related report/job state — RETAIN

TQD generic subject/profile identifiers are not allowed to redefine Party identity; integration must use canonical external IDs/contracts.

---

## crm-service

Canonical migration authority: `crm-service/database/migrations/`.

### SEO/domain registry

- `seo_domain`
- `seo_links`
- `seo_relative`
- `seo_trend`
- `seo_generation_log`
- `seo_search_performance_daily`
- `seo_url_inspection_snapshot`
- `seo_traffic_daily`
- `seo_keyword_target`
- `seo_ranking_snapshot`
- `seo_crawl_event`
- `seo_issue`

### Strategy/evidence/intent

- `seo_objective`
- `seo_objective_evidence`
- `seo_objective_metric`
- `seo_audience_segment`
- `seo_user_need`
- `seo_user_need_evidence`
- `seo_query_evidence`
- `seo_intent_hypothesis`
- `seo_intent_evidence_link`
- `seo_query_cluster`
- `seo_query_cluster_member`

### Content taxonomy/placement

- `seo_content_domain`
- `seo_content_category`
- `seo_content_category_assignment`
- `seo_content_placement`

### Measurement/control

- `seo_measurement_import_run`
- `seo_measurement_source_state`
- `seo_page_score_snapshot`
- `seo_runtime_event` — OPERATIONAL/event

### Support and admin opportunity

- `support_tickets`
- `support_ticket_notes`
- `support_ticket_events`
- `admin_opportunity_events`
- `customers` — RETAIN CRM process truth; explicitly **not** global Person identity

### Planning/news and payment integration

- `user_saved_planning_news`
- `payment_event_inbox` — OPERATIONAL inbox/idempotency
- `payment_crm_activities` — projection/activity

All above remain CRM-owned. The canonical Person may be referenced, but CRM Customer remains a scoped CRM relationship/process.

---

## notification-service

Canonical migration authority: `notification-service/database/migrations/`.

- `notification` — RETAIN
- `notification_history` — OPERATIONAL/history
- `noti_histories` — RECONCILE duplicate/legacy history naming
- `admin_histories` — OPERATIONAL/audit
- `deal_history` — OPERATIONAL/read history; Deal service remains truth owner
- `history_auth` — OPERATIONAL/audit
- `tb_history` — RECONCILE generic legacy history
- `person_configs` — RETAIN notification preferences
- `property_histories` — OPERATIONAL/projection
- `account_warning_templates` — RETAIN notification/security template
- `account_warnings` — RETAIN workflow/projection, not Account authority
- `account_warning_logs` — OPERATIONAL/audit
- `payment_event_inbox` — RETAIN inbox/idempotency
- `event_notifications` — RETAIN event→notification materialization
- `notification_delivery_intents` — RETAIN delivery workflow

---

## file-service

Canonical migration authority: `file-service/database/migrations/`.

- `file` — RETAIN File metadata authority
- `file_access` — RETAIN access relationship, but `access_id` semantics must be documented before it is treated as authorization truth

Canonical business domains reference File IDs explicitly through relation tables. File service does not infer owner semantics from generic resource types.

---

## hub-service

Canonical migration authority: `hub-service/database/migrations/`.

- `applink`
- `tb_event_queue` — OPERATIONAL
- `provinces`
- `districts`
- `wards`
- `province_v2`
- `ward_v2`
- `tb_user_guide`
- `user_guide_steps`
- `system_config`
- `faqs`
- `versions`
- `api_keys` — security/config authority for Hub only
- `interactive_events` — OPERATIONAL/event
- `error_logs` — OPERATIONAL
- `update_data` — OPERATIONAL

Hub lookup/location duplicates must not silently become the canonical planning/property geographic authority. Consumers must choose one contract for each use case.

---

# Cross-service ownership matrix

| Fact | Canonical write owner |
|---|---|
| Party / Person / Account / Profile | user-service |
| Authentication credential / session | user-service |
| Permission / Organization Role | user-service |
| Organization | user-service |
| Organization Invitation / Membership / Block | user-service |
| Organization Ownership | user-service |
| Organization Branch / assignment / manager | user-service target |
| Group / Group Membership / Leadership | organization-service |
| Deal / Invitation / Participation / Lead/Admin | organization-service |
| Deal financial terms / investments / costs | organization-service |
| Bank / Bank Account / Wallet / Payment | payment-service |
| Catalog / Subscription | user-service |
| Property | bdspro-service |
| Planning / quota / reports / POI | tqd-service |
| CRM / SEO / support | crm-service |
| Notification delivery/history | notification-service |
| File metadata/access | file-service |
| Platform lookup/help/config | hub-service |

# Acceptance rule

A table not in the v1 canonical core is **not automatically deleted**.

It is contracted only after:

1. every writer has moved to the canonical owner,
2. legacy rows are classified/backfilled,
3. old/new results are compared,
4. exact-SHA tests and runtime proof pass,
5. rollback no longer depends on the old write path,
6. the compatibility table/column has no remaining reader that treats it as authoritative.
