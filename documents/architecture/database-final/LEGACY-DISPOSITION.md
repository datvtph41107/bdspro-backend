# Legacy Database Disposition

Base source inspected: main @ b2a7a8167e9c8af2a47f1fde4f7b2a9d951611fd.

## Legend

- KEEP — already expresses an accepted durable fact; preserve owner and semantics.
- REPLACE — target fact exists but legacy representation is semantically wrong or duplicated.
- PROJECTION — useful read model/cache; never becomes write authority.
- OPERATIONAL — worker/inbox/outbox/config/runtime state; valid but not business identity.
- RETAIN-AS-EXTENSION — leaf feature is outside the core convergence and may remain under its current owner as long as it does not redefine a canonical fact.
- RETIRE — duplicate/dead/compatibility structure removed after cutover proof.

## User / identity

| Current table/pattern | Disposition | Target / reason |
| --- | --- | --- |
| user_profile | REPLACE | Person identity + Profile presentation; role/status/plan fields stop being authority |
| profile_transfer, profile_deleted | RETIRE after preservation/export | migration/deletion compatibility, not a second Person truth |
| auth_method | REPLACE | account_password_credentials / account_phone_authenticators / account_external_identities |
| auth_devices | REPLACE | account_devices |
| user_session | REPLACE | account_sessions |
| user_status | REPLACE | account_blocks plus verified authenticator/profile facts; no generic account status |
| user_info | REPLACE | Account + platform role assignments; no profile-level role authority |
| user_otp | OPERATIONAL | authentication challenge state only; never durable Person identity |
| user_pin | RETAIN-AS-EXTENSION | legacy authenticator until a dedicated PIN migration; never IAM authority |
| permissions | KEEP/REKEY | permissions.code becomes stable primitive |
| roles, role_profiles, role_permissions | REPLACE | split platform IAM from organization IAM |
| role_groups, group_permissions | RETAIN-AS-EXTENSION then review | permission grouping/catalog convenience, not assignment truth |
| admin_profiles | PROJECTION/REPLACE | Profile + system role assignment; admin-specific display fields may remain projection |
| admin_access_control, admin_access_logs | OPERATIONAL/AUDIT | keep under platform admin owner, not business identity |
| auth_config | OPERATIONAL | runtime/security config |
| block, follow, friend, bookmark_user | RETAIN-AS-EXTENSION | social graph leaf facts; must resolve endpoints to Person during identity migration |
| certification, kyc, profession | RETAIN-AS-EXTENSION | Person-scoped verification/credential facts; migrate Profile IDs to Person IDs |
| tb_contact | RETAIN-AS-EXTENSION / converge with CRM | do not let it become a second Person/CRM Contact authority |
| main_area, main_area_profile, purpose_use, purpose_use_profile, tags, tag_user, profile_media | RETAIN-AS-EXTENSION | presentation/taxonomy facts; migrate endpoints to Person/File |
| pack_price_table | RETIRE after commercial cutover | superseded by versioned catalog price items |
| catalog_products, catalog_plans, catalog_plan_versions | KEEP | mature catalog truth |
| catalog_plan_entitlements, catalog_plan_operation_policies, catalog_price_items | KEEP | mature immutable version detail |
| catalog_subscriptions, catalog_subscription_events | KEEP/REKEY | target subscriptions/events, Party subject |
| subscription_checkout_snapshots | KEEP/REKEY | immutable checkout snapshot, Party subject |
| subscription_settlement_receipts | KEEP/REKEY | idempotent settlement effect receipt |
| organizations, organization_members from user-service | REPLACE | canonical Organization model below |

## Organization / Deal

| Current table/pattern | Disposition | Target / reason |
| --- | --- | --- |
| organizations.owner_profile_id / OwnerId | RETIRE | organization_ownerships is sole current owner fact |
| organization_members role_key + invited/active/suspended/removed | REPLACE | invitations + memberships + blocks + role assignments |
| organization_permissions | CONVERGE | global permission primitive lives once; organization role_permission references it |
| organization_roles / role_permissions | REPLACE | organization_roles + organization_role_permissions |
| organization_branches type branch/department/store | REPLACE | canonical Branch only; Department/Store require separate earned facts |
| organization_branch_members.user_id | REPLACE | organization_branch_assignments.membership_id |
| organization_branches.manager_id | REPLACE | manager_membership_id; no direct user authority pointer |
| organization_log_activities | REPLACE | organization_audit_events under Organization owner |
| groups | REPLACE in place | groups with explicit archived lifecycle |
| group_members role/status | REPLACE | group_memberships |
| Group.CreatedBy leader shortcut | RETIRE | group_leads; creator is provenance |
| group_settings config_key/config_value | RETAIN-AS-EXTENSION | EAV-like config is not canonical business identity; explicit settings are required before promoting a key |
| group_log_activities | REPLACE | group_audit_events |
| group_notifications | RETIRE/PROJECTION | notification effects belong to Notification owner |
| group_documents | REPLACE | group_files referencing File |
| group_chats | REPLACE | explicit Group↔Conversation context/link |
| bank_accounts owner_id/owner_type | REPLACE | explicit owner_party_id OR owner_group_id |
| deals.owner_id/owner_type | RETIRE | explicit Organization XOR Group context |
| deals.status | REPLACE | processing derived; completed_at/canceled_at express terminal lifecycle |
| deals.deal_type | KEEP | BROKERAGE / INVEST / JOINT is real business classification |
| deal_of_organization | RETIRE | duplicates Deal Organization context |
| deal_of_group | RETIRE | duplicates Deal Group context |
| deal_of_branch | RETIRE after branch backfill | direct secondary Branch scope on Deal |
| deal_members | REPLACE | invitation + participation + scoped relationships + authority + economics |
| deal_members.is_owner | RETIRE | broken shadow truth |
| deal_members.role_key | RETIRE | mixes Lead/Admin/Member/Customer/Partner |
| deal_members.member_type | RETIRE | same Person can occupy multiple relationships |
| deal_members.amount_commit | REPLACE | investment/commitment fact, not participation identity |
| deal_members.commission_* | REPLACE | deal_commission_terms |
| deal_members.done_investment | RETIRE | projection derived from deal_investments |
| deal_milestones | KEEP/REKEY | deal_milestones with timestamp-based completion |
| investments | REPLACE | deal_investments referencing DealParticipation and Money |
| attach_documents | RETIRE | owning domains reference File explicitly |
| internal_notes | KEEP/REKEY | deal_internal_notes |
| cost_types, deal_costs | KEEP/REKEY | explicit currency/minor units |
| deal_action / notification deal_history | CONVERGE | one Deal audit/event owner, not duplicated history stores |

## Property / listing / asset

| Current table/pattern | Disposition | Target / reason |
| --- | --- | --- |
| property_identify/property_lineage generation | RETIRE after convergence | duplicate Property generation; preserve external identity/history during backfill |
| newer property table | KEEP/REKEY | becomes properties canonical identity |
| property_land_info, property_building_info | KEEP/REKEY | one-to-one details |
| property_media | KEEP/REKEY | stable File references preferred over raw URL authority |
| tag/property_tag_link | KEEP/REKEY | property_feature_tags/property_tags |
| property_external_ref | KEEP/REKEY | explicit source_system + external_key |
| property_edvidence | KEEP/REKEY | property_evidence + File |
| property_relation generic Product/Asset relation | RETIRE | explicit Property→Asset and Listing.property_id |
| property_user | RETIRE | Asset ownership / Listing context are distinct facts |
| products | REPLACE | listings; Product is market offer, not Property identity |
| products.owner_id/owner_of | RETIRE | explicit Person/Organization/Group Listing context |
| product_price | KEEP/REKEY | listing_prices history |
| product_private | KEEP/REKEY | listing_private_terms |
| product_media | KEEP/REKEY | listing_media |
| product_user/product_organization | REPLACE | listing context/assignments; do not duplicate owner |
| assets | REPLACE in place | assets tied to Property |
| assets.owner_id/owner_of | RETIRE | asset_ownerships → Party |
| asset legal/document URL fields | REKEY | File references |
| projects/project_build/apartment/developer | CONVERGE | real_estate_projects/project_buildings/project_units; developer points to Organization Party |
| product_market and other query structs/views | PROJECTION | read model only |

## Payment

| Current table | Disposition |
| --- | --- |
| commerce_orders | KEEP/REKEY subject to Party |
| commerce_payment_attempts | KEEP |
| commerce_settlements | KEEP |
| commerce_fulfillments | KEEP |
| commerce_command_effects | KEEP OPERATIONAL-IDEMPOTENCY |
| commerce_outbox_events | KEEP OPERATIONAL |
| banks | KEEP/REKEY lifecycle |
| wallets.balance | PROJECTION | immutable wallet_entries becomes economic truth before wallet feature is authoritative |
| wallet_transactions | REPLACE | wallet_entries |
| wallet_audit_logs | RETIRE after ledger proof | redundant balance mutation log once ledger is truth |
| withdrawal_requests | KEEP/REKEY |
| payment_methods | KEEP/REKEY | organization_payment_methods |
| transaction_types | RETAIN-AS-EXTENSION | reference catalog if still used |
| dashboard_metrics, payment_stats | PROJECTION | never payment truth |

## Notification and File

| Current table/pattern | Disposition |
| --- | --- |
| notification owner_of/owner_id | REPLACE | notifications.recipient_party_id |
| notification_history, noti_histories, admin_histories | RETIRE/RETURN-TO-OWNER | generic history warehouse is not canonical domain audit |
| deal_history | RETIRE after Deal audit cutover | Deal owner keeps Deal audit |
| history_auth | MOVE/CONVERGE | authentication/security audit belongs to identity/security owner |
| property_histories | MOVE/CONVERGE | Property owner keeps Property audit |
| person_configs | REKEY | notification_preferences |
| account_warning_* | RETAIN-AS-EXTENSION | moderation/warning domain; must target Party/Account explicitly |
| payment_event_inbox | KEEP OPERATIONAL | event dedupe |
| event_notifications | REKEY | notifications |
| notification_delivery_intents | KEEP | claim-version delivery effect |
| file | KEEP/REKEY | files; absolute_path remains forbidden |
| file_access | REPLACE | file_grants with Party endpoint |
| owner_namespace/owner_key | KEEP OPERATIONAL-IDEMPOTENCY | create-effect identity, not business owner |

## CRM

| Current table/pattern | Disposition |
| --- | --- |
| customers | REPLACE/RENAME | crm_opportunities; current columns describe sales opportunity, not global Customer identity |
| admin_opportunity_events | KEEP/REKEY | crm_opportunity_events |
| contact/tag/pipeline/stage models | KEEP/REKEY | crm_contacts, crm_contact_tags, crm_pipelines, crm_pipeline_stages |
| appointments / appointment_participants | KEEP/REKEY | Person endpoints |
| support_tickets/notes/events | KEEP/REKEY | explicit Person and Organization endpoints |
| support ticket images | REPLACE | support_ticket_files → File |
| seo_domain | KEEP/RENAME | seo_pages |
| seo_links | KEEP |
| seo_relative | KEEP/RENAME | seo_relations |
| seo_generation_log | KEEP/RENAME | seo_generation_runs |
| seo strategy/intent/measurement tables dropped by 000024 | RETIRE confirmed | do not resurrect |
| payment_event_inbox, payment_crm_activities | KEEP PROJECTION | CRM consumer inbox/activity; never payment authority |
| user_saved_planning_news | RETAIN-AS-EXTENSION | Person→seo_page saved relation |

## Chat

| Current table/pattern | Disposition |
| --- | --- |
| conversations | KEEP/REKEY Person/context endpoints |
| conversation_memberships | KEEP as sole membership truth |
| participants | RETIRE | duplicate conversation membership |
| messages | KEEP/REKEY File/Person references and explicit lifecycle timestamps |
| message_reactions | KEEP |
| read_positions | KEEP/RENAME | conversation_read_positions |
| approval_requests | REPLACE | conversation_join_requests |
| chat_events | KEEP |
| chat_outbox | KEEP OPERATIONAL |
| chat_timelines | PROJECTION | timeline/read model, not second Message truth |
| background/presentation tables | RETAIN-AS-EXTENSION | UI state only |

## Social

| Current table/pattern | Disposition |
| --- | --- |
| news_feed owner_of/owner_id | REPLACE | posts.author_person_id + explicit Group/Organization context |
| news_feed_of_user | RETIRE | duplicate author/context truth |
| news_feed_of_group | RETIRE | duplicate Group context |
| news_feed_medias | KEEP/REKEY | post_media |
| comment | KEEP/REKEY | post_comments |
| tb_like target_id/target_type | REPLACE | post_reactions and comment_reactions |
| news_feed_shares | KEEP/REKEY | explicit target columns |
| friend_tag | KEEP/REKEY | post_friend_tags |
| report/report_reason | KEEP/REKEY | content_reports/report_reasons with explicit target columns |

## Hub / administrative catalog

| Current table/pattern | Disposition |
| --- | --- |
| provinces/districts/wards and province_v2/ward_v2 | CONVERGE | one Hub-owned administrative_provinces/administrative_wards authority; aliases/mapping preserved during cutover |
| applink | KEEP/RENAME | app_links |
| system_config | KEEP/RENAME | system_configs |
| faqs | KEEP |
| versions | KEEP/RENAME | app_versions |
| api_keys | KEEP/REKEY | store hash, never plain secret authority |
| tb_user_guide/user_guide_steps | RETAIN-AS-EXTENSION | product help content |
| tb_event_queue | OPERATIONAL | queue state |
| interactive_events, error_logs, update_data | OPERATIONAL/PROJECTION | telemetry/import state, not business identity |

## TQD / Planning

| Current table/pattern | Disposition |
| --- | --- |
| qh_planning_projects/documents | KEEP/REKEY | planning_projects/planning_documents |
| qh_layers, qh_layer_families, qh_regions, qh_labels, qh_label_layers | KEEP | mature planning model |
| qh_authority_issuring | KEEP/RENAME | planning_authorities |
| qh_legends, qh_land_use, qh_layer_land_use | KEEP |
| qh_legal_documents and legal links | KEEP | explicit planning legal facts |
| qh_jurisdictions | KEEP | planning-specific jurisdiction, not platform province authority |
| qh_layer_legals | CONVERGE | migrate into canonical planning_legal_documents links |
| qh_region_extends | KEEP/PROJECTION depending writer path | versioned imported region extension; do not make a second qh_regions authority |
| qh_region_import_error_logs | OPERATIONAL | import diagnostics |
| qh_layer_resolver_configs | OPERATIONAL | resolver config |
| qh_shape/qh_direction/qh_parcel_direction | RETAIN-AS-EXTENSION | planning geometry taxonomy |
| parcels/qh_parcel_info | CONVERGE | one Parcel identity plus planning detail extensions |
| map_points | KEEP | TQD MapPoint owner |
| amenities/open_hours/poi_media/poi_amenities/directory_* | RETAIN-AS-EXTENSION | directory/POI feature owner remains TQD until separately split |
| search_index | PROJECTION | searchable projection only |
| user_view_events | KEEP ANALYTICS FACT | immutable interaction event |
| user_view_history | PROJECTION | aggregate of view events |
| user_followed_parcels/projects | KEEP/REKEY | Person endpoints |
| user_reported | REPLACE/RENAME | reports |
| report_events | KEEP |
| report_jobs | KEEP OPERATIONAL | claim fencing required |
| quota_usage_events | KEEP | durable quota consumption |
| quota_pools | KEEP MATERIALIZED-STATE | reconciles from durable usage + entitlement |
| pro_ai_jobs | KEEP OPERATIONAL | AI job state |

## Services without relational canonical truth

- Gateway: ingress/verification/routing; no business database authority.
- Auth compatibility service: no second IAM database authority.
- Search: search engine projection only; source services remain truth.
- Relay: transport only.
- Assistant: provider/runtime boundary; no new business database authority.

## Contraction rule

A RETIRE decision is architectural, not permission to drop a table now. Every destructive migration waits for: complete backfill, zero unexplained compare drift, caller cutover, rollback window, and exact-SHA proof.