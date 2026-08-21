SELECT effect_key, fingerprint, order_id, subject_kind, subject_id,
       product_code, plan_code, plan_version_id, tier_rank,
       subscription_term_days, terms_checksum, occurred_at, action,
       subscription_id, created_at
  FROM commercial_subscription_effects
 LIMIT 0;
