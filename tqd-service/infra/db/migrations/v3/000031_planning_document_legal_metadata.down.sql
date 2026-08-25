ALTER TABLE IF EXISTS qh_planning_documents
    DROP COLUMN IF EXISTS indicators,
    DROP COLUMN IF EXISTS authority,
    DROP COLUMN IF EXISTS sector,
    DROP COLUMN IF EXISTS field,
    DROP COLUMN IF EXISTS signer;
