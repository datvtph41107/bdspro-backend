# Canonical Database Invariants

These invariants are part of the schema even when DBML cannot express the physical constraint directly.

## Identity and authentication

1. A Person is a Party; an Organization is a Party. Party identity is never inferred from Profile or Account.
2. Profile is descriptive/presentation state; it does not carry organization authority.
3. Account is digital identity for authentication. Closing Account does not delete Person.
4. External identity is unique by provider + provider_subject.
5. A current account block disables effective digital authority without rewriting account identity.
6. Platform system roles and Organization roles are separate assignment systems.

## Organization

1. Non-archived Organization has exactly one current Ownership row.
2. Ownership references a current, unblocked Membership of the same Organization.
3. Creator/founder is provenance, not a current authority fact.
4. A Person has at most one current Membership in one Organization.
5. Rejoin creates a new Membership episode.
6. Invitation and Membership are different records.
7. Organization Invitation target has exactly one form: person OR email OR phone.
8. Invitation terminal outcomes are immutable. Re-issuing an offer creates a new row.
9. Accepted Invitation produces exactly one Membership via originating_invitation_id.
10. Current Membership has exactly one current RoleAssignment.
11. RoleAssignment role belongs to the same Organization as Membership.
12. At most one current MembershipBlock exists per Membership.
13. Block does not end Membership or RoleAssignment; it only removes effective authority.
14. Ending Membership ends/removes all current dependent authority/structure relations in the same transaction.
15. Branch belongs to exactly one Organization.
16. Branch Assignment and manager Membership must belong to the same Organization as the Branch.
17. Closed Branch receives no new assignment, manager transfer or resource placement.

## Group

1. Group creator is not authority after creation unless a durable relation says so.
2. Current Group has exactly one current Lead.
3. Lead references a current Membership of that Group.
4. A Person has at most one current GroupMembership per Group.
5. Rejoin is a new GroupMembership.

## Deal

1. A Deal has exactly one primary context: Organization XOR Group.
2. Branch is permitted only for an Organization Deal and must belong to that Organization.
3. Deal lifecycle is PROCESSING until exactly one terminal outcome is set: COMPLETED or CANCELED.
4. Deal type BROKERAGE / INVEST / JOINT is classification, not lifecycle.
5. target_profit_minor is money and always paired with currency.
6. A Person has at most one current DealParticipation per Deal.
7. Invitation does not grant participation. Acceptance creates Participation atomically.
8. Deal Invitation is immutable after terminal outcome; re-invite creates a new Invitation.
9. Customer and Partner relationships reference a DealParticipation and may overlap.
10. Admin capacity is independent of Customer/Partner relationship.
11. Current Deal Lead references a current DealParticipation. Generic role update cannot create or remove Lead.
12. Commission terms reference Participation; they do not define participation identity.
13. Commission term has either percentage OR fixed_amount_minor+currency, never both.
14. Investment references Participation; investment progress is not copied as canonical boolean onto Participation.
15. Historical investments, costs and completed Deal relations are not deleted when Participation later ends.

## Property, Listing and Asset

1. Property identity is independent of Listing and Asset.
2. Listing is a market offer and has exactly one business context: Person XOR Organization XOR Group.
3. Listing Branch, when present, requires Organization context and same Organization.
4. Listing price change creates a new price fact; previous price remains historical.
5. Asset means an owned business asset tied to a Property.
6. Asset Ownership references Party and is temporal. At most one current exclusive owner is allowed until co-ownership is explicitly earned.
7. Property evidence and legal documents reference stable File IDs, never deployment filesystem paths.

## CRM, Chat and Social

1. CRM Customer is not a global identity. CRM opportunity references Person/Contact explicitly.
2. Chat has one membership authority: conversation_memberships. Legacy participants becomes compatibility/projection only.
3. Message sequence is unique within a Conversation.
4. Reaction identity is target-specific plus Person; no generic target_type/id is required for canonical post/comment reactions.
5. Post author is Person. Organization/Group are context, not author identity.

## Catalog / Subscription / Payment

1. Published PlanVersion is immutable. New commercial terms require a new version.
2. Checkout snapshot freezes plan version, currency, amount and terms checksum.
3. One command key in the same subject scope cannot create a second checkout/order with a different fingerprint.
4. Settlement is unique by provider + provider transaction ID.
5. Provider settlement evidence is never overwritten by fulfillment failure.
6. Fulfillment has one row per Order and uses lease + claim_version fencing.
7. Outbox event ID is unique and publication is retriable.
8. Subscription activation consumes settlement effect idempotently.
9. Wallet economic truth is immutable entries. Current balance is derived/materialized and must reconcile to entries.

## File and Notification

1. File identity uses stable storage key/reference; absolute host/container path is never durable authority.
2. File access is explicit and revocable.
3. Notification owns recipient/effect/delivery state only; Deal/Property/Auth audit does not live there as canonical truth.
4. Event inbox deduplicates producer event IDs before creating effects.
5. Delivery intent is unique for one Notification+Channel unless the product explicitly models multiple sends.

## Planning / Report / Quota

1. Hub owns the platform administrative province/ward catalog. TQD may keep domain jurisdiction and read projections but does not create a second platform authority.
2. Planning source facts and public/search projections are distinct.
3. Report command key is unique per subject.
4. Accepted Report, durable quota usage and ReportJob commit together.
5. quota_usage_events are immutable durable consumption facts.
6. quota_pools is current materialized usage state and must reconcile from usage facts.
7. ReportJob running claim requires locked_by/lease and monotonic claim_version.

## Concurrency lock order

- Organization governance mutation: lock Organization first, then current ownership/memberships in deterministic order.
- Branch governance mutation: lock Branch; lock Organization first only when an Organization-level invariant is also changed.
- Deal lead/participation acceptance: lock Deal or Invitation coordinator first, then dependent Participation rows deterministically.
- Subscription/payment workers: claim via database lease/fencing; do not rely on process-local mutexes.

## Authorization invariant

A relationship row never means allow by itself. Effective decision remains:

authenticated Account/Person → current context Membership/Participation → no blocking restriction → Role/Permission when applicable → resource relationship/scope → resource lifecycle/business state → policy → ALLOW or DENY.

Search/list queries must apply the same allowed scope before returning rows; frontend filtering is never authorization.