-- SEO Operating System V2.U typed source identity.
-- Administrative Unit numeric IDs originate from separate province/ward tables
-- and may overlap. ref_url stores the complete typed identity (province:<id> or
-- ward:<id>) and is therefore the stable uniqueness key for this page family.
CREATE UNIQUE INDEX IF NOT EXISTS uq_seo_domain_administrative_unit_typed_identity
ON seo_domain(ref_type, ref_source, ref_url)
WHERE deleted_at IS NULL
  AND ref_source = 'tqd.administrative_unit'
  AND ref_url IS NOT NULL
  AND btrim(ref_url) <> '';
