# Change Summary — Developer Operating Architecture

## Kept

- business ownership/correctness mechanisms hiện hữu;
- canonical Compose topology đã có trong baseline;
- durable payment/outbox/notification/report behavior;
- compatibility binaries còn cần proof retire;
- production deployment files/semantics.

## Simplified

- root command vocabulary;
- daily Compose rebuild behavior;
- duplicated service Make mechanics;
- migration commands;
- local DB fallback ports;
- repository root clutter;
- service documentation/navigation;
- workspace/task naming;
- active Cursor/editor rule được gom về một canonical source contract, rule cũ chuyển vào history.

## Corrected

- BDSPro migration v2 version 000003 pair/rollback;
- CRM Payment Inbox migration ownership;
- CRM ad-hoc migration file placement;
- empty CRM startup AutoMigrate call;
- source-level migration guard for all 9 owners.

## Explicitly not done

- no production deploy script rewrite;
- no production Compose rewrite;
- no registry/rollout/rollback redesign;
- no speculative removal of non-default services;
- no broad QHPRO runtime ENV rename without compatibility proof;
- no business package mass-renaming merely for architecture aesthetics.
